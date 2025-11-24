package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type WalletRequest struct {
	PlayerID    string `json:"playerId"`
	AmountCents int64  `json:"amountCents"`
	Currency    string `json:"currency"`
	RefID       string `json:"refId"`
	Meta        any    `json:"meta,omitempty"`
}

type WalletResponse struct {
	Status       string `json:"status"` // OK | REJECTED
	BalanceCents int64  `json:"balanceCents"`
	Reason       string `json:"reason,omitempty"`
}

func (h *Handler) Debit(w http.ResponseWriter, r *http.Request) {
	var req WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	resp := WalletResponse{
		Status:       "OK",
		BalanceCents: 100000,
	}

	writeJSON(w, resp)
}

func (h *Handler) Credit(w http.ResponseWriter, r *http.Request) {
	var req WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	resp := WalletResponse{
		Status:       "OK",
		BalanceCents: 100000,
	}

	writeJSON(w, resp)
}

func writeJSON(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Println("error writing response:", err)
		return
	}
}
