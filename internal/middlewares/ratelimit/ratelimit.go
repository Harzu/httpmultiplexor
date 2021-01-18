package ratelimit

import (
	"net/http"
	"sync/atomic"
)

type Middleware interface {
	RateLimitMiddleware(next http.Handler) http.Handler
}

type rateLimitMiddleware struct {
	limit     int
	rateCount int32
}

func NewRateLimitMiddleware(limit int) Middleware {
	return &rateLimitMiddleware{limit: limit}
}

func (m *rateLimitMiddleware) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.allow() {
			http.Error(w, "exceeded rate limit", http.StatusInternalServerError)
			return
		}

		atomic.AddInt32(&m.rateCount, 1)
		next.ServeHTTP(w, r)
		atomic.AddInt32(&m.rateCount, -1)
	})
}

func (m *rateLimitMiddleware) allow() bool {
	return m.rateCount < int32(m.limit)
}
