package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ReconciliationResponse struct {
	Transactions []Transaction `json:"transactions"`
}

func (h *Handler) Reconciliation(w http.ResponseWriter, _ *http.Request) {
	resp := ReconciliationResponse{
		Transactions: h.wallet.store.Transactions.Items,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("json encode error:", err)
	}
}
