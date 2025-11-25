package http

import (
	"integration-hub/internal/operator"
	"integration-hub/internal/storage"
	"integration-hub/internal/storage/db"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	store    *storage.IdempotencyStore
	operator *operator.Client
	queries  *db.Queries
}

func NewHandler(q *db.Queries, op *operator.Client) *Handler {
	return &Handler{
		store:    storage.NewIdempotencyStore(q),
		operator: op,
		queries:  q,
	}
}

func (h *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(IdempotencyMiddleware(h.store))
		r.Use(SignatureMiddleware("testsecret123"))

		r.Post("/wallet/debit", h.Debit)
		r.Post("/wallet/credit", h.Credit)
	})

	r.Post("/webhook/operator", h.OperatorWebhook)

	return r
}
