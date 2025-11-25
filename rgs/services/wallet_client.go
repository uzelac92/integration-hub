package services

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"rgs/observability"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type WalletClient struct {
	baseURL string
	secret  string
	client  *http.Client
}

type walletRequest struct {
	PlayerID    string `json:"playerId"`
	AmountCents int64  `json:"amountCents"`
	Currency    string `json:"currency"`
	RefID       string `json:"refId"`
}

func NewWalletClient(baseURL, secret string) *WalletClient {
	return &WalletClient{
		baseURL: baseURL,
		secret:  secret,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (w *WalletClient) sign(body []byte, timestamp string) string {
	mac := hmac.New(sha256.New, []byte(w.secret))
	mac.Write(body)
	mac.Write([]byte(timestamp))
	return hex.EncodeToString(mac.Sum(nil))
}

func (w *WalletClient) call(ctx context.Context, path string, playerID int32, amount float64, requestID string) (bool, error) {
	reqBody := walletRequest{
		PlayerID:    fmt.Sprintf("p%d", playerID),
		AmountCents: int64(amount * 100),
		Currency:    "EUR",
		RefID:       requestID,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return false, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signature := w.sign(data, timestamp)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", w.baseURL+path, bytes.NewBuffer(data))
	if err != nil {
		return false, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Idempotency-Key", requestID)
	httpReq.Header.Set("X-Timestamp", timestamp)
	httpReq.Header.Set("X-Signature", signature)

	resp, err := w.client.Do(httpReq)
	if err != nil {
		observability.Logger.Error("wallet http request failed", zap.Error(err))
		return false, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			observability.Logger.Error("failed to close wallet call", zap.Error(err))
		}
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Errorf("hub wallet error: %s", resp.Status)
	}

	return true, nil
}

func (w *WalletClient) Debit(ctx context.Context, playerID int32, amount float64, requestID string) (bool, error) {
	return w.call(ctx, "/wallet/debit", playerID, amount, requestID)
}

func (w *WalletClient) Credit(ctx context.Context, playerID int32, amount float64, requestID string) (bool, error) {
	return w.call(ctx, "/wallet/credit", playerID, amount, requestID)
}
