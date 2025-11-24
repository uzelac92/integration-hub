package http

import (
	"integration-hub/internal/storage"
	"integration-hub/internal/storage/db"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store *storage.IdempotencyStore
}

func NewHandler(q *db.Queries) *Handler {
	return &Handler{
		store: storage.NewIdempotencyStore(q),
	}
}

func (h *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(IdempotencyMiddleware(h.store))

	r.Use(SignatureMiddleware("my-secret-key"))

	r.Post("/debit", h.Debit)
	r.Post("/credit", h.Credit)

	return r
}
