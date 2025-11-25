package main

import (
	"log"
	"net/http"
)

func main() {
	store := NewStore()
	service := NewService(store)
	handler := NewHandler(service)

	r := NewRouter(handler)

	log.Println("Operator mock running on :9000")
	log.Fatal(http.ListenAndServe(":9000", r))
}
