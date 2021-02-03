package ratelimit

import (
	"net/http"
	"sync"
)

type Middleware interface {
	RateLimitMiddleware(next http.Handler) http.Handler
}

type rateLimitMiddleware struct {
	mu        *sync.RWMutex
	limit     int
	rateCount int32
}

func NewRateLimitMiddleware(limit int) Middleware {
	return &rateLimitMiddleware{
		mu:    &sync.RWMutex{},
		limit: limit,
	}
}

func (m *rateLimitMiddleware) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		if m.rateCount > int32(m.limit) {
			http.Error(w, "exceeded rate limit", http.StatusInternalServerError)
			return
		}

		m.rateCount++
		next.ServeHTTP(w, r)
		m.rateCount--
	})
}
