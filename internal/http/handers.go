package http

import (
	"encoding/json"
	"integration-hub/internal/operator"
	"log"
	"net/http"
	"strings"
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

	if req.PlayerID == "" || req.RefID == "" {
		http.Error(w, "missing playerId or refId", http.StatusBadRequest)
		return
	}
	if req.AmountCents <= 0 {
		http.Error(w, "amountCents must be > 0", http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		http.Error(w, "missing currency", http.StatusBadRequest)
		return
	}

	opReq := operator.WithdrawRequest{
		Amount:   req.AmountCents,
		Currency: req.Currency,
		RefID:    req.RefID,
	}

	opResp, err := h.operator.Withdraw(req.PlayerID, opReq)
	if err != nil {
		http.Error(w, "operator error: "+err.Error(), http.StatusBadGateway)
		return
	}

	resp := WalletResponse{
		Status:       strings.ToUpper(opResp.Status), // OK | REJECTED
		BalanceCents: opResp.Balance,
		Reason:       opResp.ErrorMessage,
	}

	writeJSON(w, resp)
}

func (h *Handler) Credit(w http.ResponseWriter, r *http.Request) {
	var req WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.PlayerID == "" || req.RefID == "" {
		http.Error(w, "missing playerId or refId", http.StatusBadRequest)
		return
	}
	if req.AmountCents <= 0 {
		http.Error(w, "amountCents must be > 0", http.StatusBadRequest)
		return
	}
	if req.Currency == "" {
		http.Error(w, "missing currency", http.StatusBadRequest)
		return
	}

	opReq := operator.DepositRequest{
		Amount:   req.AmountCents,
		Currency: req.Currency,
		RefID:    req.RefID,
	}

	opResp, err := h.operator.Deposit(req.PlayerID, opReq)
	if err != nil {
		http.Error(w, "operator error: "+err.Error(), http.StatusBadGateway)
		return
	}

	resp := WalletResponse{
		Status:       strings.ToUpper(opResp.Status),
		BalanceCents: opResp.Balance,
		Reason:       opResp.ErrorMessage,
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
