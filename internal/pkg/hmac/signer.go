package hmac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Sign(secret string, body []byte, timestamp string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	h.Write([]byte(timestamp))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifySignature(secret string, body []byte, timestamp string, signature string) bool {
	expected := Sign(secret, body, timestamp)
	return hmac.Equal([]byte(signature), []byte(expected))
}
