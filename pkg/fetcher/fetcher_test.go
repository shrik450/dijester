package fetcher

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockHTTPClient implements the HTTPClient interface for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestHTTPFetcher_WithOptions(t *testing.T) {
	// Test that options are correctly applied
	customUserAgent := "CustomAgent/1.0"
	customTimeout := 15 * time.Second

	f := NewHTTPFetcher(
		WithUserAgent(customUserAgent),
		WithTimeout(customTimeout),
	)

	if f.userAgent != customUserAgent {
		t.Errorf("Expected user agent '%s', got '%s'", customUserAgent, f.userAgent)
	}

	if f.timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, f.timeout)
	}
}

func TestHTTPFetcher_FetchURL(t *testing.T) {
	// Create a test server
	testContent := "test content"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user agent is set correctly
		if r.Header.Get("User-Agent") != "TestAgent/1.0" {
			t.Errorf("Expected User-Agent 'TestAgent/1.0', got '%s'", r.Header.Get("User-Agent"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testContent))
	}))
	defer server.Close()

	// Create fetcher with custom client
	f := NewHTTPFetcher(
		WithUserAgent("TestAgent/1.0"),
		WithClient(server.Client()),
	)

	// Fetch content
	ctx := context.Background()
	content, err := f.FetchURL(ctx, server.URL)
	if err != nil {
		t.Fatalf("FetchURL returned error: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Expected content '%s', got '%s'", testContent, string(content))
	}

	// Test error handling for non-200 response
	errorServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}),
	)
	defer errorServer.Close()

	f = NewHTTPFetcher(WithClient(errorServer.Client()))
	_, err = f.FetchURL(ctx, errorServer.URL)

	if err == nil {
		t.Error("Expected error for non-200 response, got nil")
	}
}

func TestRateLimiter_Wait(t *testing.T) {
	// Create a rate limiter with a 100ms interval
	interval := 100 * time.Millisecond
	limiter := NewRateLimiter(interval)

	// First request to a domain should not wait
	ctx := context.Background()
	domain := "example.com"
	url := "https://" + domain + "/path"

	start := time.Now()
	err := limiter.Wait(ctx, url)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("First Wait returned error: %v", err)
	}

	if elapsed > 10*time.Millisecond {
		t.Errorf("First request should not wait, but waited %v", elapsed)
	}

	// Second request to the same domain should wait
	start = time.Now()
	err = limiter.Wait(ctx, url)
	elapsed = time.Since(start)

	if err != nil {
		t.Fatalf("Second Wait returned error: %v", err)
	}

	if elapsed < interval {
		t.Errorf("Second request should wait at least %v, but only waited %v", interval, elapsed)
	}

	// Request to a different domain should not wait
	differentURL := "https://different.com/path"
	start = time.Now()
	err = limiter.Wait(ctx, differentURL)
	elapsed = time.Since(start)

	if err != nil {
		t.Fatalf("Different domain Wait returned error: %v", err)
	}

	if elapsed > 10*time.Millisecond {
		t.Errorf("Different domain request should not wait, but waited %v", elapsed)
	}
}

func TestLimitedFetcher(t *testing.T) {
	// Create a mock HTTP client that counts requests
	requestCount := 0
	mock := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			requestCount++
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       http.NoBody,
			}, nil
		},
	}

	// Create a base fetcher with the mock client
	baseF := NewHTTPFetcher(WithClient(mock))

	// Create a rate-limited fetcher with a long interval (to test actual waiting)
	interval := 200 * time.Millisecond
	f := NewLimitedFetcher(baseF, interval)

	// Make two quick requests to the same domain
	ctx := context.Background()
	url := "https://example.com/path"

	start := time.Now()

	// First request should go through immediately
	_, err := f.FetchURL(ctx, url)
	if err != nil {
		t.Fatalf("First request returned error: %v", err)
	}

	// Second request should be rate-limited
	_, err = f.FetchURL(ctx, url)
	if err != nil {
		t.Fatalf("Second request returned error: %v", err)
	}

	elapsed := time.Since(start)

	// Verify that we waited at least the interval
	if elapsed < interval {
		t.Errorf(
			"Expected to wait at least %v between requests, but only waited %v",
			interval,
			elapsed,
		)
	}

	// Verify that both requests were made
	if requestCount != 2 {
		t.Errorf("Expected 2 requests to be made, got %d", requestCount)
	}
}
