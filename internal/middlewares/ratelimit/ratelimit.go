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
	return &rateLimitMiddleware{
		limit: limit,
	}
}

func (m *rateLimitMiddleware) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&m.rateCount) > int32(m.limit) {
			http.Error(w, "exceeded rate limit", http.StatusTooManyRequests)
			return
		}

		atomic.AddInt32(&m.rateCount, 1)
		next.ServeHTTP(w, r)
		atomic.AddInt32(&m.rateCount, -1)
	})
}
