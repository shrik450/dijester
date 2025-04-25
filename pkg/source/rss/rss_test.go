package rss

import (
	"context"
	"net/http"
	"testing"

	"github.com/shrik450/dijester/pkg/fetcher"
)

// MockHTTPClient for testing
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// MockFetcher extends HTTPFetcher for testing
type MockFetcher struct {
	*fetcher.HTTPFetcher
	MockFetchURLAsString func(ctx context.Context, url string) (string, error)
}

// FetchURLAsString overrides the HTTPFetcher method
func (m *MockFetcher) FetchURLAsString(ctx context.Context, url string) (string, error) {
	return m.MockFetchURLAsString(ctx, url)
}

func TestRSSSource_Configure(t *testing.T) {
	source := New()

	// Test valid configuration
	validConfig := map[string]interface{}{
		"url":          "https://example.com/feed.xml",
		"name":         "Test Feed",
		"max_articles": 5,
	}

	err := source.Configure(validConfig)
	if err != nil {
		t.Errorf("Configure with valid config returned error: %v", err)
	}

	if source.url != "https://example.com/feed.xml" {
		t.Errorf("Expected URL 'https://example.com/feed.xml', got '%s'", source.url)
	}

	if source.name != "Test Feed" {
		t.Errorf("Expected name 'Test Feed', got '%s'", source.name)
	}

	if source.maxArticles != 5 {
		t.Errorf("Expected maxArticles 5, got %d", source.maxArticles)
	}

	// Test missing URL
	invalidConfig := map[string]interface{}{
		"name": "Test Feed",
	}

	err = source.Configure(invalidConfig)
	if err == nil {
		t.Error("Configure without URL should return error")
	}

	// Test default values
	minimalConfig := map[string]interface{}{
		"url": "https://example.com/feed.xml",
	}

	source = New()
	err = source.Configure(minimalConfig)
	if err != nil {
		t.Errorf("Configure with minimal config returned error: %v", err)
	}

	if source.name != "rss" {
		t.Errorf("Expected default name 'rss', got '%s'", source.name)
	}

	if source.maxArticles != 10 {
		t.Errorf("Expected default maxArticles 10, got %d", source.maxArticles)
	}
}

func TestRSSSource_Fetch(t *testing.T) {
	// Sample RSS feed content
	sampleFeed := `
<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Sample Feed</title>
    <link>https://example.com</link>
    <description>A sample RSS feed for testing</description>
    <item>
      <title>Article 1</title>
      <link>https://example.com/article1</link>
      <description>Description of article 1</description>
      <pubDate>Mon, 01 Jan 2023 12:00:00 GMT</pubDate>
      <author>author1@example.com (Author One)</author>
      <category>Category1</category>
      <category>Tag1</category>
    </item>
    <item>
      <title>Article 2</title>
      <link>https://example.com/article2</link>
      <description>Description of article 2</description>
      <pubDate>Tue, 02 Jan 2023 12:00:00 GMT</pubDate>
    </item>
    <item>
      <title>Article 3 (no content)</title>
      <link>https://example.com/article3</link>
      <pubDate>Wed, 03 Jan 2023 12:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>
`

	// Create a mock HTTP client
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       http.NoBody,
			}, nil
		},
	}

	// Create a mock fetcher that returns our sample feed
	mockFetcher := &MockFetcher{
		HTTPFetcher: fetcher.NewHTTPFetcher(
			fetcher.WithClient(mockClient),
		),
		MockFetchURLAsString: func(ctx context.Context, url string) (string, error) {
			return sampleFeed, nil
		},
	}

	// Create and configure our source
	source := New()
	err := source.Configure(map[string]interface{}{
		"url":          "https://example.com/feed.xml",
		"name":         "Test Feed",
		"max_articles": 3,
	})
	if err != nil {
		t.Fatalf("Failed to configure source: %v", err)
	}

	// Fetch articles
	ctx := context.Background()
	articles, err := source.Fetch(ctx, mockFetcher)
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}

	// Verify the number of articles
	// Note: We expect 2 articles because the third article has no content or description
	if len(articles) != 2 {
		t.Fatalf("Expected 2 articles, got %d", len(articles))
	}

	// Verify the first article's content
	article := articles[0]
	if article.Title != "Article 1" {
		t.Errorf("Expected title 'Article 1', got '%s'", article.Title)
	}

	if article.URL != "https://example.com/article1" {
		t.Errorf("Expected URL 'https://example.com/article1', got '%s'", article.URL)
	}

	if article.Content != "Description of article 1" {
		t.Errorf("Expected content 'Description of article 1', got '%s'", article.Content)
	}

	// In our implementation, when we use description as content, summary is empty
	if article.Summary != "" {
		t.Errorf(
			"Expected empty summary when content is from description, got '%s'",
			article.Summary,
		)
	}

	if article.Author != "Author One" {
		t.Errorf("Expected author 'Author One', got '%s'", article.Author)
	}

	if article.SourceName != "Test Feed" {
		t.Errorf("Expected source name 'Test Feed', got '%s'", article.SourceName)
	}

	if len(article.Tags) != 2 || article.Tags[0] != "Category1" || article.Tags[1] != "Tag1" {
		t.Errorf("Expected tags ['Category1', 'Tag1'], got %v", article.Tags)
	}

	// Verify article with no content uses description as content
	article = articles[1]
	if article.Content != "Description of article 2" {
		t.Errorf("Expected content to fall back to description, got '%s'", article.Content)
	}
	if article.Summary != "" {
		t.Errorf(
			"Expected empty summary when content falls back to description, got '%s'",
			article.Summary,
		)
	}

	// Test max articles limit
	source = New()
	err = source.Configure(map[string]interface{}{
		"url":          "https://example.com/feed.xml",
		"max_articles": 1,
	})
	if err != nil {
		t.Fatalf("Failed to configure source: %v", err)
	}

	articles, err = source.Fetch(ctx, mockFetcher)
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}

	if len(articles) != 1 {
		t.Errorf("Expected 1 article with max_articles=1, got %d", len(articles))
	}
}
