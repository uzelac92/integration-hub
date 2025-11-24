package http

import (
	"integration-hub/internal/storage"
	"log"
	"net/http"
)

type responseRecorder struct {
	http.ResponseWriter
	body []byte
}

func IdempotencyMiddleware(store *storage.IdempotencyStore) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				http.Error(w, "missing Idempotency-Key", http.StatusBadRequest)
				return
			}

			if data, found, _ := store.Get(key); found {
				w.Header().Set("Content-Type", "application/json")
				_, err := w.Write(data)
				if err != nil {
					log.Printf("failed to write cached response: %v", err)
				}
				return
			}

			rec := &responseRecorder{
				ResponseWriter: w,
				body:           []byte{},
			}
			next.ServeHTTP(rec, r)

			err := store.Save(key, rec.body)
			if err != nil {
				log.Printf("failed to save idempotency record: %v", err)
			}
		})
	}
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}
