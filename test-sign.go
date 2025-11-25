package main

import (
	"fmt"
	"integration-hub/internal/pkg/hmac"
	"time"
)

func main() {
	body := `{"playerId":"p1","amountCents":100,"currency":"EUR","refId":"abc"}`
	secret := "my-secret-key"
	ts := fmt.Sprintf("%d", time.Now().Unix())

	sig := hmac.Sign(secret, []byte(body), ts)

	fmt.Println("X-Timestamp:", ts)
	fmt.Println("X-Signature:", sig)
}
