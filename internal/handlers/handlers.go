package handlers

import (
	"encoding/json"
	"fmt"
	"httpmultiplexor/internal/config"
	"httpmultiplexor/internal/middlewares"
	"net/http"
)

type Handlers interface {
	Register()
}

type handlers struct {
	config      *config.Config
	middlewares *middlewares.MiddlewareProcessor
}

func New(cfg *config.Config) Handlers {
	return &handlers{
		config:      cfg,
		middlewares: middlewares.NewMiddlewareProcessor(cfg),
	}
}

func (h *handlers) Register() {
	http.Handle("/pages", h.middlewares.RateLimitMiddleware.RateLimitMiddleware(http.HandlerFunc(h.getPages)))
}

func (h *handlers) writeResponse(w http.ResponseWriter, res interface{}) error {
	resBytes, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed response marshal: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resBytes); err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}

	return nil
}