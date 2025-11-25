package http

import "github.com/go-chi/chi/v5"

func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(RateLimit)
	r.Use(RandomFailures)

	r.Post("/v2/players/{playerID}/withdraw", h.Withdraw)
	r.Post("/v2/players/{playerID}/deposit", h.Deposit)

	return r
}
