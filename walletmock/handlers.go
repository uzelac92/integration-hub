package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	wallet *Service
}

func NewHandler(w *Service) *Handler {
	return &Handler{wallet: w}
}

type WalletReq struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	RefID    string `json:"refId"`
}

type WalletResp struct {
	Status  string `json:"status"`
	Balance int64  `json:"balance"`
	Error   string `json:"error,omitempty"`
}

func (h *Handler) Withdraw(wr http.ResponseWriter, r *http.Request) {
	playerID := chi.URLParam(r, "playerID")

	var req WalletReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(wr, "invalid json", http.StatusBadRequest)
		return
	}

	balance, err := h.wallet.Withdraw(playerID, req.Amount)
	if err != nil {
		writeJSON(wr, WalletResp{
			Status:  "REJECTED",
			Balance: balance,
			Error:   err.Error(),
		})
		return
	}

	writeJSON(wr, WalletResp{
		Status:  "OK",
		Balance: balance,
	})
}

func (h *Handler) Deposit(wr http.ResponseWriter, r *http.Request) {
	playerID := chi.URLParam(r, "playerID")

	var req WalletReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(wr, "invalid json", http.StatusBadRequest)
		return
	}

	balance := h.wallet.Deposit(playerID, req.Amount)

	writeJSON(wr, WalletResp{
		Status:  "OK",
		Balance: balance,
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Println("json encode error:", err)
		return
	}
}
