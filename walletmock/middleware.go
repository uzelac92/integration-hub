package main

import (
	"math/rand"
	"net/http"
	"time"
)

const maxRequestsPerMin = 60

var (
	requestCount = 0
	lastReset    = time.Now()
)

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if time.Since(lastReset) > time.Minute {
			requestCount = 0
			lastReset = time.Now()
		}

		requestCount++
		if requestCount > maxRequestsPerMin {
			w.Header().Set("Retry-After", "2")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RandomFailures(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.Intn(10) < 2 {
			http.Error(w, "operator internal error", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}
