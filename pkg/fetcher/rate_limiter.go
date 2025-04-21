package fetcher

import (
	"context"
	"net/url"
	"sync"
	"time"
)

// RateLimiter limits the rate of requests to a domain.
type RateLimiter struct {
	mu           sync.Mutex
	lastRequests map[string]time.Time
	minInterval  time.Duration
}

// NewRateLimiter creates a new rate limiter with the given minimum interval between requests.
func NewRateLimiter(minInterval time.Duration) *RateLimiter {
	return &RateLimiter{
		lastRequests: make(map[string]time.Time),
		minInterval:  minInterval,
	}
}

// Wait blocks until the rate limit allows a request to the given URL.
func (r *RateLimiter) Wait(ctx context.Context, rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	domain := parsed.Hostname()

	r.mu.Lock()
	lastRequest, ok := r.lastRequests[domain]
	now := time.Now()
	r.lastRequests[domain] = now
	r.mu.Unlock()

	if !ok {
		// First request to this domain
		return nil
	}

	elapsed := now.Sub(lastRequest)
	if elapsed >= r.minInterval {
		// Enough time has passed
		return nil
	}

	// Need to wait
	waitTime := r.minInterval - elapsed

	select {
	case <-time.After(waitTime):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// LimitedFetcher wraps an HTTP fetcher with rate limiting.
type LimitedFetcher struct {
	Fetcher *HTTPFetcher
	limiter *RateLimiter
}

// NewLimitedFetcher creates a new rate-limited HTTP fetcher.
func NewLimitedFetcher(fetcher *HTTPFetcher, minInterval time.Duration) *LimitedFetcher {
	return &LimitedFetcher{
		Fetcher: fetcher,
		limiter: NewRateLimiter(minInterval),
	}
}

// FetchURL fetches a URL with rate limiting.
func (f *LimitedFetcher) FetchURL(ctx context.Context, url string) ([]byte, error) {
	if err := f.limiter.Wait(ctx, url); err != nil {
		return nil, err
	}

	return f.Fetcher.FetchURL(ctx, url)
}

// FetchURLAsString fetches a URL as a string with rate limiting.
func (f *LimitedFetcher) FetchURLAsString(ctx context.Context, url string) (string, error) {
	if err := f.limiter.Wait(ctx, url); err != nil {
		return "", err
	}

	return f.Fetcher.FetchURLAsString(ctx, url)
}
