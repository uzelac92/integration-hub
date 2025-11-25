package main

import (
	walletHttp "integration-hub/walletmock/internal/http"
	"integration-hub/walletmock/internal/wallet"
	"log"
	"net/http"
)

func main() {
	store := wallet.NewStore()
	service := wallet.NewService(store)
	handler := walletHttp.NewHandler(service)

	r := walletHttp.NewRouter(handler)

	log.Println("Operator mock running on :9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
