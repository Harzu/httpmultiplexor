package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"httpmultiplexor/internal/pagefetcher"
	"io/ioutil"
	"net/http"
	"time"
)

const urlsLimit = 20

func (h *handlers) getPages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "failed method, need POST method", http.StatusBadRequest)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed read body: %s", err.Error()), http.StatusBadRequest)
		return
	}

	var urls []string
	if err := json.Unmarshal(bodyBytes, &urls); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse body urls: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if len(urls) > urlsLimit {
		http.Error(w,"exciting urls limit", http.StatusBadRequest)
	}

	requestTimeout := time.Duration(len(urls) + 10) * time.Second
	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	fetcher := pagefetcher.NewPageFetcher(urls, h.config.ClientRequestRateLimit)
	res, err := fetcher.Fetch(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed fetch: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if err := h.writeResponse(w, res); err != nil {
		http.Error(w, fmt.Sprintf("failed write response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
