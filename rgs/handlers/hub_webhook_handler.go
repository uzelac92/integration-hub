package handlers

import (
	"encoding/json"
	"net/http"
	"rgs/observability"
	"rgs/sqlc"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type HubWebhookHandler struct {
	queries *sqlc.Queries
}

func NewHubWebhookHandler(q *sqlc.Queries) *HubWebhookHandler {
	return &HubWebhookHandler{queries: q}
}

func (h *HubWebhookHandler) Receive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	eventID := r.Header.Get("X-Event-ID")
	if eventID == "" {
		eventID = uuid.New().String()
	}

	exists, err := h.queries.HubWebhookEventExists(ctx, eventID)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if exists {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status": "duplicate"}`))
		if err != nil {
			observability.Logger.Error("Failed to write response", zap.Error(err))
			return
		}
		return
	}

	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	payloadBytes, _ := json.Marshal(payload)
	err = h.queries.InsertHubWebhookEvent(ctx, sqlc.InsertHubWebhookEventParams{
		EventID: eventID,
		Payload: payloadBytes,
	})
	if err != nil {
		http.Error(w, "failed to store webhook", http.StatusInternalServerError)
		return
	}

	observability.Logger.Info("Webhook received from Integration Hub",
		zap.String("event_id", eventID),
		zap.Any("payload", payload),
	)

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"status": "received"}`))
	if err != nil {
		observability.Logger.Error("failed to write response", zap.Error(err))
		return
	}
}
