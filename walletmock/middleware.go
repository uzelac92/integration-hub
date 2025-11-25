package main

import (
	"math/rand/v2"
	"net/http"
	"sync"
	"time"
)

const maxRequestsPerMin = 60

var (
	mu          sync.Mutex
	reqCount    = 0
	windowStart = time.Now()
)

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mu.Lock()
		now := time.Now()

		if now.Sub(windowStart) >= time.Minute {
			reqCount = 0
			windowStart = now
		}

		if reqCount >= maxRequestsPerMin {
			mu.Unlock()
			w.Header().Set("Retry-After", "1")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		reqCount++
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

var rng = rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixNano()<<1)))

func RandomFailures(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if rng.IntN(100) < 10 {
			http.Error(w, "random 500", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
