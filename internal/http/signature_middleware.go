package http

import (
	"bytes"
	"integration-hub/internal/pkg/hmac"
	"io"
	"net/http"
	"strconv"
	"time"
)

const skew = 5 * time.Minute

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func SignatureMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			signature := r.Header.Get("X-Signature")
			timestamp := r.Header.Get("X-Timestamp")

			if signature == "" || timestamp == "" {
				http.Error(w, "missing signature headers", http.StatusUnauthorized)
				return
			}

			ts, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				http.Error(w, "invalid timestamp", http.StatusBadRequest)
				return
			}
			now := time.Now().Unix()
			if abs(now-ts) > int64(skew.Seconds()) {
				http.Error(w, "timestamp out of range", http.StatusUnauthorized)
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "cannot read body", http.StatusBadRequest)
				return
			}
			errClose := r.Body.Close()
			if errClose != nil {
				http.Error(w, "cannot read body", http.StatusBadRequest)
				return
			}

			// Restore body
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			if !hmac.VerifySignature(secret, body, timestamp, signature) {
				http.Error(w, "signature invalid", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
