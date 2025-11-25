package main

import (
	"fmt"
	"integration-hub/config"
	"integration-hub/internal/operator"
	"integration-hub/internal/storage"
	"integration-hub/internal/webhook"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	httpInternal "integration-hub/internal/http"
)

func main() {
	cfg := config.LoadConfig()
	port := fmt.Sprintf(":%s", cfg.Port)

	db := storage.Connect(cfg)
	opClient := operator.NewClient(cfg.WalletUrl)

	dispatcher := webhook.NewDispatcher(db.Queries, fmt.Sprintf("%s/webhook/hub", cfg.RgsUrl))
	dispatcher.Start()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := httpInternal.NewHandler(db.Queries, opClient)
	r.Mount("/", h.Router())

	// Health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			log.Println("Error writing health response:", err)
			return
		}
	})

	log.Println("Integration Hub running on ", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}
