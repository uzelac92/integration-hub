package main

import (
	"integration-hub/internal/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	hhtp "integration-hub/internal/http"
)

func main() {
	db := storage.Connect()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := hhtp.NewHandler(db.Queries)
	r.Mount("/wallet", h.Router())

	// Health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			log.Println("Error writing health response:", err)
			return
		}
	})

	log.Println("Integration Hub running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
