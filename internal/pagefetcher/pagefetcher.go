package pagefetcher

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type PageFetcher interface {
	Fetch(ctx context.Context) ([]FetchResult, error)
}

type pageFetcher struct {
	urls                   []string
	errChan                chan error
	resultChan             chan FetchResult
	clientRequestRateLimit int
}

func NewPageFetcher(urls []string, clientRequestRateLimit int) PageFetcher {
	return &pageFetcher{
		urls:                   urls,
		errChan:                make(chan error, clientRequestRateLimit),
		resultChan:             make(chan FetchResult, len(urls)),
		clientRequestRateLimit: clientRequestRateLimit,
	}
}

func (f *pageFetcher) Fetch(ctx context.Context) ([]FetchResult, error) {
	result := make([]FetchResult, 0, len(f.urls))
	semaphoreChan := make(chan struct{}, f.clientRequestRateLimit)

	go func() {
		for _, url := range f.urls {
			semaphoreChan <- struct{}{}
			go f.fetch(ctx, url)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Printf(
				"fetch urls handle canceling, urls: %d, processed: %d, unprocessed: %d",
				len(f.urls), len(result), len(f.urls) - len(result),
			)

			return nil, ctx.Err()
		case err := <-f.errChan:
			return nil, err
		case fetchRes := <-f.resultChan:
			log.Printf("result by: %s processed", fetchRes.Url)

			result = append(result, fetchRes)
			if len(result) == len(f.urls) {
				return result, nil
			}

			<-semaphoreChan
		}
	}
}

func (f *pageFetcher) fetch(ctx context.Context, url string) {
	log.Printf("start process url: %s", url)

	done := make(chan struct{}, 1)
	client := &http.Client{Timeout: 1 * time.Second}

	go func() {
		resp, err := client.Get(url)
		if err != nil {
			f.errChan <- fmt.Errorf("failed to send request: %w", err)
			return
		}
		defer func() {
			if closeErr := resp.Body.Close(); err == nil && closeErr != nil {
				f.errChan <- closeErr
			}
		}()

		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			f.errChan <- fmt.Errorf("failed to read response bytes: %w", err)
			return
		}

		f.resultChan <- FetchResult{
			Url:     url,
			Payload: string(respBytes),
		}

		done <- struct{}{}
		close(done)
	}()

	select {
	case <-ctx.Done():
		client.CloseIdleConnections()
		return
	case <-done:
		return
	}
}

type FetchResult struct {
	Url     string `json:"url"`
	Payload string `json:"payload"`
}