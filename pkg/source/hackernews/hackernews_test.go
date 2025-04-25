package hackernews

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/shrik450/dijester/pkg/fetcher"
)

// MockHTTPClient implements the HTTPClient interface for testing
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

func TestHackerNewsSource_Configure(t *testing.T) {
	source := New()

	validConfig := map[string]interface{}{
		"name":         "Custom HN",
		"max_articles": 15,
		"min_score":    200,
		"page":         "new",
	}

	err := source.Configure(validConfig)
	if err != nil {
		t.Errorf("Configure with valid config returned error: %v", err)
	}

	if source.name != "Custom HN" {
		t.Errorf("Expected name 'Custom HN', got '%s'", source.name)
	}

	if source.maxArticles != 15 {
		t.Errorf("Expected maxArticles 15, got %d", source.maxArticles)
	}

	if source.minScore != 200 {
		t.Errorf("Expected minScore 200, got %d", source.minScore)
	}

	if source.pageType != NewPage {
		t.Errorf("Expected pageType 'new', got '%s'", source.pageType)
	}

	// Test with past page
	source = New()
	err = source.Configure(map[string]interface{}{
		"page": "past",
	})
	if err != nil {
		t.Errorf("Configure with past page returned error: %v", err)
	}

	if source.pageType != PastPage {
		t.Errorf("Expected pageType 'past', got '%s'", source.pageType)
	}

	// Test with active page
	source = New()
	err = source.Configure(map[string]interface{}{
		"page": "active",
	})
	if err != nil {
		t.Errorf("Configure with active page returned error: %v", err)
	}

	if source.pageType != ActivePage {
		t.Errorf("Expected pageType 'active', got '%s'", source.pageType)
	}

	// Test with jobs page
	source = New()
	err = source.Configure(map[string]interface{}{
		"page": "jobs",
	})
	if err != nil {
		t.Errorf("Configure with jobs page returned error: %v", err)
	}

	if source.pageType != JobsPage {
		t.Errorf("Expected pageType 'jobs', got '%s'", source.pageType)
	}

	// Test with invalid page (should default to frontpage)
	source = New()
	err = source.Configure(map[string]interface{}{
		"page": "invalid",
	})
	if err != nil {
		t.Errorf("Configure with invalid page returned error: %v", err)
	}

	if source.pageType != FrontPage {
		t.Errorf("Expected pageType to default to 'frontpage', got '%s'", source.pageType)
	}

	// Test default values
	source = New()
	err = source.Configure(map[string]interface{}{})
	if err != nil {
		t.Errorf("Configure with empty config returned error: %v", err)
	}

	if source.name != "hackernews" {
		t.Errorf("Expected default name 'hackernews', got '%s'", source.name)
	}

	if source.maxArticles != 30 {
		t.Errorf("Expected default maxArticles 30, got %d", source.maxArticles)
	}

	if source.minScore != 10 {
		t.Errorf("Expected default minScore 10, got %d", source.minScore)
	}

	if source.pageType != FrontPage {
		t.Errorf("Expected default pageType 'frontpage', got '%s'", source.pageType)
	}
}

func TestHackerNewsSource_ParseHNPage(t *testing.T) {
	htmlContent := `
<html>
  <body>
    <table class="itemlist">
      <tr class="athing" id="100">
        <td class="title">
          <span class="titleline">
            <a href="https://example.com/item100">High Score Item</a>
          </span>
        </td>
      </tr>
      <tr>
        <td>
          <span class="score">150 points</span> by
          <a class="hnuser">user1</a>
          <span class="age">1 hour ago</span> |
          <a>10 comments</a>
        </td>
      </tr>
      <tr class="athing" id="101">
        <td class="title">
          <span class="titleline">
            <a href="item?id=101">Text Post</a>
          </span>
        </td>
      </tr>
      <tr>
        <td>
          <span class="score">50 points</span> by
          <a class="hnuser">user2</a>
          <span class="age">2 hours ago</span> |
          <a>5 comments</a>
        </td>
      </tr>
    </table>
  </body>
</html>
`
	stories, err := parseHNPage(htmlContent)
	if err != nil {
		t.Fatalf("parseHNPage returned error: %v", err)
	}

	if len(stories) != 2 {
		t.Fatalf("Expected 2 stories, got %d", len(stories))
	}

	story := stories[0]
	if story.ID != 100 {
		t.Errorf("Expected ID 100, got %d", story.ID)
	}

	if story.Title != "High Score Item" {
		t.Errorf("Expected title 'High Score Item', got '%s'", story.Title)
	}

	if story.URL != "https://example.com/item100" {
		t.Errorf("Expected URL 'https://example.com/item100', got '%s'", story.URL)
	}

	if story.By != "user1" {
		t.Errorf("Expected author 'user1', got '%s'", story.By)
	}

	if story.Score != 150 {
		t.Errorf("Expected score 150, got %d", story.Score)
	}

	if story.Comments != 10 {
		t.Errorf("Expected 10 comments, got %d", story.Comments)
	}

	// Test relative URL handling
	story = stories[1]
	if story.URL != "https://news.ycombinator.com/item?id=101" {
		t.Errorf("Expected URL 'https://news.ycombinator.com/item?id=101', got '%s'", story.URL)
	}
}

func TestHackerNewsSource_Fetch(t *testing.T) {
	responseMap := make(map[string]string)

	// Mock HTML content for different page types
	frontpageHTML := `
<html>
  <body>
    <table class="itemlist">
      <tr class="athing" id="100">
        <td class="title">
          <span class="titleline">
            <a href="https://example.com/item100">Front Page Item</a>
          </span>
        </td>
      </tr>
      <tr>
        <td>
          <span class="score">150 points</span> by
          <a class="hnuser">user1</a>
          <span class="age">1 hour ago</span> |
          <a>10 comments</a>
        </td>
      </tr>
    </table>
  </body>
</html>
`
	newpageHTML := `
<html>
  <body>
    <table class="itemlist">
      <tr class="athing" id="101">
        <td class="title">
          <span class="titleline">
            <a href="https://example.com/item101">New Page Item</a>
          </span>
        </td>
      </tr>
      <tr>
        <td>
          <span class="score">50 points</span> by
          <a class="hnuser">user2</a>
          <span class="age">2 hours ago</span> |
          <a>5 comments</a>
        </td>
      </tr>
    </table>
  </body>
</html>
`
	pastpageHTML := `
<html>
  <body>
    <table class="itemlist">
      <tr class="athing" id="102">
        <td class="title">
          <span class="titleline">
            <a href="https://example.com/item102">Past Page Item</a>
          </span>
        </td>
      </tr>
      <tr>
        <td>
          <span class="score">120 points</span> by
          <a class="hnuser">user3</a>
          <span class="age">3 hours ago</span> |
          <a>8 comments</a>
        </td>
      </tr>
    </table>
  </body>
</html>
`
	activepageHTML := `
<html>
  <body>
    <table class="itemlist">
      <tr class="athing" id="103">
        <td class="title">
          <span class="titleline">
            <a href="https://example.com/item103">Active Page Item</a>
          </span>
        </td>
      </tr>
      <tr>
        <td>
          <span class="score">130 points</span> by
          <a class="hnuser">user4</a>
          <span class="age">4 hours ago</span> |
          <a>15 comments</a>
        </td>
      </tr>
    </table>
  </body>
</html>
`
	jobspageHTML := `
<html>
  <body>
    <table class="itemlist">
      <tr class="athing" id="104">
        <td class="title">
          <span class="titleline">
            <a href="https://example.com/item104">Jobs Page Item</a>
          </span>
        </td>
      </tr>
      <tr>
        <td>
          <span class="age">5 hours ago</span>
        </td>
      </tr>
    </table>
  </body>
</html>
`

	// Add page HTML responses
	responseMap[HNPageURLs[FrontPage]] = frontpageHTML
	responseMap[HNPageURLs[NewPage]] = newpageHTML
	responseMap[HNPageURLs[PastPage]] = pastpageHTML
	responseMap[HNPageURLs[ActivePage]] = activepageHTML
	responseMap[HNPageURLs[JobsPage]] = jobspageHTML

	// Mock API responses for item details
	mockItem := func(id int, title string, score int, by string, isDead bool, isDeleted bool) string {
		item := HNItem{
			ID:          id,
			Title:       title,
			URL:         fmt.Sprintf("https://example.com/item%d", id),
			Text:        fmt.Sprintf("Content of item %d", id),
			By:          by,
			Score:       score,
			Time:        time.Now().Unix(),
			Type:        "story",
			Kids:        []int{1000, 1001},
			Dead:        isDead,
			Deleted:     isDeleted,
			Descendants: 10,
		}
		itemJSON, _ := json.Marshal(item)
		return string(itemJSON)
	}

	responseMap[fmt.Sprintf(itemURLFormat, 100)] = mockItem(
		100,
		"Front Page API Item",
		150,
		"user1",
		false,
		false,
	)
	responseMap[fmt.Sprintf(itemURLFormat, 101)] = mockItem(
		101,
		"New Page API Item",
		50,
		"user2",
		false,
		false,
	)
	responseMap[fmt.Sprintf(itemURLFormat, 102)] = mockItem(
		102,
		"Past Page API Item",
		120,
		"user3",
		true,
		false,
	)
	responseMap[fmt.Sprintf(itemURLFormat, 103)] = mockItem(
		103,
		"Active Page API Item",
		130,
		"user4",
		false,
		true,
	)
	responseMap[fmt.Sprintf(itemURLFormat, 104)] = mockItem(
		104,
		"Jobs Page API Item",
		140,
		"user5",
		false,
		false,
	)

	// Setup the mock fetcher
	mockFetcher := &MockFetcher{
		HTTPFetcher: fetcher.NewHTTPFetcher(),
		MockFetchURLAsString: func(ctx context.Context, url string) (string, error) {
			if response, ok := responseMap[url]; ok {
				return response, nil
			}
			if url == "https://example.com/item100" || url == "https://example.com/item101" {
				return "<html><body><article><h1>Article Content</h1><p>This is the full article content.</p></article></body></html>", nil
			}
			return "", fmt.Errorf("unexpected URL: %s", url)
		},
	}

	// Test frontpage fetch
	source := New()
	err := source.Configure(map[string]any{
		"name":         "Test HN",
		"max_articles": 10,
		"min_score":    10,
		"page":         "frontpage",
	})
	if err != nil {
		t.Fatalf("Failed to configure source: %v", err)
	}

	ctx := context.Background()
	articles, err := source.Fetch(ctx, mockFetcher)
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}

	if len(articles) != 1 {
		t.Fatalf("Expected 1 articles, got %d", len(articles))
	}

	article := articles[0]
	if article.Title != "Front Page API Item" {
		t.Errorf("Expected title 'Front Page API Item', got '%s'", article.Title)
	}

	// Test newpage fetch
	source = New()
	err = source.Configure(map[string]any{
		"page":      "new",
		"min_score": 10,
	})
	if err != nil {
		t.Fatalf("Failed to configure source: %v", err)
	}

	articles, err = source.Fetch(ctx, mockFetcher)
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}

	if len(articles) != 1 {
		t.Fatalf("Expected 1 articles, got %d", len(articles))
	}

	article = articles[0]
	if article.Title != "New Page API Item" {
		t.Errorf("Expected title 'New Page API Item', got '%s'", article.Title)
	}

	// Test past page fetch
	source = New()
	err = source.Configure(map[string]any{
		"page":      "past",
		"min_score": 10,
		"show_dead": true,
	})
	if err != nil {
		t.Fatalf("Failed to configure source: %v", err)
	}

	articles, err = source.Fetch(ctx, mockFetcher)
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}

	if len(articles) != 1 {
		t.Fatalf("Expected 1 article, got %d", len(articles))
	}

	article = articles[0]
	if article.Title != "Past Page API Item" {
		t.Errorf("Expected title 'Past Page API Item', got '%s'", article.Title)
	}

	// Test active page fetch
	source = New()
	err = source.Configure(map[string]any{
		"page":         "active",
		"min_score":    10,
		"show_deleted": true,
	})
	if err != nil {
		t.Fatalf("Failed to configure source: %v", err)
	}

	articles, err = source.Fetch(ctx, mockFetcher)
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}

	if len(articles) != 1 {
		t.Fatalf("Expected 1 article, got %d", len(articles))
	}

	article = articles[0]
	if article.Title != "Active Page API Item" {
		t.Errorf("Expected title 'Active Page API Item', got '%s'", article.Title)
	}

	// Test jobs page fetch
	source = New()
	err = source.Configure(map[string]any{
		"page":      "jobs",
		"min_score": 0, // Jobs don't have scores
	})
	if err != nil {
		t.Fatalf("Failed to configure source: %v", err)
	}

	articles, err = source.Fetch(ctx, mockFetcher)
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}

	if len(articles) != 1 {
		t.Fatalf("Expected 1 article, got %d", len(articles))
	}

	article = articles[0]
	if article.Title != "Jobs Page API Item" {
		t.Errorf("Expected title 'Jobs Page API Item', got '%s'", article.Title)
	}
}
