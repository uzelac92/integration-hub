package http

import (
	"encoding/json"
	"errors"
	"integration-hub/internal/operator"
	"log"
	"net/http"
	"strings"
	"time"
)

type WalletRequest struct {
	PlayerID    string `json:"playerId"`
	AmountCents int64  `json:"amountCents"`
	Currency    string `json:"currency"`
	RefID       string `json:"refId"`
	Meta        any    `json:"meta,omitempty"`
}

type WalletResponse struct {
	Status       string `json:"status"`
	BalanceCents int64  `json:"balanceCents"`
	Reason       string `json:"reason,omitempty"`
}

func (h *Handler) Debit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validateWalletReq(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	log.Printf("[DEBIT] player=%s amount=%d status=%s duration=%s",
		req.PlayerID, req.AmountCents, resp.Status, time.Since(start))
}

func (h *Handler) Credit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validateWalletReq(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	log.Printf("[CREDIT] player=%s amount=%d status=%s duration=%s",
		req.PlayerID, req.AmountCents, resp.Status, time.Since(start))
}

func validateWalletReq(req WalletRequest) error {
	if req.PlayerID == "" {
		return errors.New("missing playerId")
	}
	if req.RefID == "" {
		return errors.New("missing refId")
	}
	if req.AmountCents <= 0 {
		return errors.New("amountCents must be > 0")
	}
	if req.Currency == "" {
		return errors.New("missing currency")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Println("error writing response:", err)
		return
	}
}
