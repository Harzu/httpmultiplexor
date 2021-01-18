package middlewares

import (
	"httpmultiplexor/internal/config"
	"httpmultiplexor/internal/middlewares/ratelimit"
)

type MiddlewareProcessor struct {
	RateLimitMiddleware ratelimit.Middleware
}

func NewMiddlewareProcessor(cfg *config.Config) *MiddlewareProcessor {
	return &MiddlewareProcessor{
		RateLimitMiddleware: ratelimit.NewRateLimitMiddleware(cfg.HandleRateLimit),
	}
}