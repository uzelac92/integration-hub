package http

import (
	"encoding/json"
	"integration-hub/internal/storage/db"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type WebhookPayload map[string]any

func (h *Handler) OperatorWebhook(w http.ResponseWriter, r *http.Request) {
	var payload WebhookPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid webhook payload", http.StatusBadRequest)
		return
	}

	eventID := r.Header.Get("X-Event-ID")
	if eventID == "" {
		eventID = uuid.New().String()
	}

	payloadBytes, _ := json.Marshal(payload)

	err := h.queries.InsertWebhookOutbox(r.Context(), db.InsertWebhookOutboxParams{
		EventID: eventID,
		Payload: payloadBytes,
	})
	if err != nil {
		log.Println("error inserting webhook outbox", err)
	}

	w.WriteHeader(http.StatusOK)
	_, errWrite := w.Write([]byte(`{"status":"received"}`))
	if errWrite != nil {
		log.Println("Error writing response: ", errWrite)
		return
	}
}
