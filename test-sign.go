package main

import (
	"fmt"
	"integration-hub/internal/pkg/hmac"
	"time"
)

func main() {
	body := `{"playerId":"p1","amountCents":1000,"currency":"EUR","refId":"abc"}`
	secret := "testsecret123"
	ts := fmt.Sprintf("%d", time.Now().Unix())

	sig := hmac.Sign(secret, []byte(body), ts)

	fmt.Println("Body: ", string(body))
	fmt.Println("X-Timestamp:", ts)
	fmt.Println("X-Signature:", sig)
}
