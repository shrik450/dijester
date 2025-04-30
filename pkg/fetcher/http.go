package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient is an interface for HTTP clients.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPFetcher provides utilities for fetching content via HTTP.
type HTTPFetcher struct {
	client    HTTPClient
	userAgent string
	timeout   time.Duration
}

// HTTPFetcherOption configures an HTTPFetcher.
type HTTPFetcherOption func(*HTTPFetcher)

// WithUserAgent sets the User-Agent header for requests.
func WithUserAgent(userAgent string) HTTPFetcherOption {
	return func(f *HTTPFetcher) {
		f.userAgent = userAgent
	}
}

// WithTimeout sets the timeout for HTTP requests.
func WithTimeout(timeout time.Duration) HTTPFetcherOption {
	return func(f *HTTPFetcher) {
		f.timeout = timeout
	}
}

// WithClient sets a custom HTTP client.
func WithClient(client HTTPClient) HTTPFetcherOption {
	return func(f *HTTPFetcher) {
		f.client = client
	}
}

// NewHTTPFetcher creates a new HTTP fetcher with the given options.
func NewHTTPFetcher(opts ...HTTPFetcherOption) *HTTPFetcher {
	f := &HTTPFetcher{}

	for _, opt := range opts {
		opt(f)
	}

	if f.client == nil {
		f.client = &http.Client{
			Timeout: f.timeout,
		}
	}

	return f
}

// FetchURL retrieves content from a URL.
func (f *HTTPFetcher) FetchURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", f.userAgent)

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // Limit to 1MB
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return body, nil
}

// FetchURLAsString retrieves content from a URL as a string.
func (f *HTTPFetcher) FetchURLAsString(ctx context.Context, url string) (string, error) {
	body, err := f.FetchURL(ctx, url)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// StreamURL streams content from a URL to the provided writer.
func (f *HTTPFetcher) StreamURL(ctx context.Context, url string, writer io.Writer) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", f.userAgent)

	resp, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("copying response body: %w", err)
	}

	return nil
}
