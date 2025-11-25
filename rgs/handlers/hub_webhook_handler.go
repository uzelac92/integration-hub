package handlers

import (
	"encoding/json"
	"net/http"
	"rgs/observability"

	"go.uber.org/zap"
)

type HubWebhookHandler struct{}

func NewHubWebhookHandler() *HubWebhookHandler {
	return &HubWebhookHandler{}
}

func (h *HubWebhookHandler) Receive(w http.ResponseWriter, r *http.Request) {
	var payload map[string]any

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid webhook payload", http.StatusBadRequest)
		return
	}

	observability.Logger.Info("Webhook received from Integration Hub",
		zap.Any("payload", payload),
	)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status": "received"}`))
}
