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

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

type MockFetcher struct {
	*fetcher.HTTPFetcher
	MockFetchURLAsString func(ctx context.Context, url string) (string, error)
}

func (m *MockFetcher) FetchURLAsString(ctx context.Context, url string) (string, error) {
	return m.MockFetchURLAsString(ctx, url)
}

func TestHackerNewsSource_Configure(t *testing.T) {
	source := New(nil)

	validConfig := map[string]interface{}{
		"name":         "Custom HN",
		"max_articles": 15,
		"min_score":    200,
		"categories":   []interface{}{"front_page", "new"},
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

	if len(source.categories) != 2 || source.categories[0] != "front_page" ||
		source.categories[1] != "new" {
		t.Errorf("Expected categories ['front_page', 'new'], got %v", source.categories)
	}

	source = New(nil)
	err = source.Configure(map[string]interface{}{})
	if err != nil {
		t.Errorf("Configure with empty config returned error: %v", err)
	}

	if source.name != "hackernews" {
		t.Errorf("Expected default name 'hackernews', got '%s'", source.name)
	}

	if source.maxArticles != 10 {
		t.Errorf("Expected default maxArticles 10, got %d", source.maxArticles)
	}

	if source.minScore != 100 {
		t.Errorf("Expected default minScore 100, got %d", source.minScore)
	}

	if len(source.categories) != 1 || source.categories[0] != "front_page" {
		t.Errorf("Expected default categories ['front_page'], got %v", source.categories)
	}
}

func TestHackerNewsSource_Fetch(t *testing.T) {
	responseMap := make(map[string]string)

	topStories := []int{100, 101, 102, 103, 104}
	topStoriesJSON, _ := json.Marshal(topStories)
	responseMap[topStoriesURL] = string(topStoriesJSON)

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
		"High Score Item",
		150,
		"user1",
		false,
		false,
	)
	responseMap[fmt.Sprintf(itemURLFormat, 101)] = mockItem(
		101,
		"Low Score Item",
		50,
		"user2",
		false,
		false,
	)
	responseMap[fmt.Sprintf(itemURLFormat, 102)] = mockItem(
		102,
		"Dead Item",
		120,
		"user3",
		true,
		false,
	)
	responseMap[fmt.Sprintf(itemURLFormat, 103)] = mockItem(
		103,
		"Deleted Item",
		130,
		"user4",
		false,
		true,
	)
	responseMap[fmt.Sprintf(itemURLFormat, 104)] = mockItem(
		104,
		"Non-story Item",
		140,
		"user5",
		false,
		false,
	)

	var item104 HNItem
	json.Unmarshal([]byte(responseMap[fmt.Sprintf(itemURLFormat, 104)]), &item104)
	item104.Type = "comment"
	item104JSON, _ := json.Marshal(item104)
	responseMap[fmt.Sprintf(itemURLFormat, 104)] = string(item104JSON)

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

	source := New(mockFetcher)
	err := source.Configure(map[string]interface{}{
		"name":         "Test HN",
		"max_articles": 10,
		"min_score":    100,
	})
	if err != nil {
		t.Fatalf("Failed to configure source: %v", err)
	}

	ctx := context.Background()
	articles, err := source.Fetch(ctx)
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}

	if len(articles) != 1 {
		t.Fatalf("Expected 1 article, got %d", len(articles))
	}

	article := articles[0]
	if article.Title != "High Score Item" {
		t.Errorf("Expected title 'High Score Item', got '%s'", article.Title)
	}

	if article.URL != "https://example.com/item100" {
		t.Errorf("Expected URL 'https://example.com/item100', got '%s'", article.URL)
	}

	if article.Author != "user1" {
		t.Errorf("Expected author 'user1', got '%s'", article.Author)
	}

	if article.SourceName != "Test HN" {
		t.Errorf("Expected source name 'Test HN', got '%s'", article.SourceName)
	}

	source = New(mockFetcher)
	err = source.Configure(map[string]interface{}{
		"min_score": 40,
	})
	if err != nil {
		t.Fatalf("Failed to configure source: %v", err)
	}

	articles, err = source.Fetch(ctx)
	if err != nil {
		t.Fatalf("Fetch returned error: %v", err)
	}

	if len(articles) != 2 {
		t.Fatalf("Expected 2 articles with lower min_score, got %d", len(articles))
	}

	if score, ok := articles[0].Metadata["score"].(int); !ok || score != 150 {
		t.Errorf("Expected metadata to contain score 150, got %v", articles[0].Metadata["score"])
	}

	if comments, ok := articles[0].Metadata["comments"].(int); !ok || comments != 10 {
		t.Errorf(
			"Expected metadata to contain comments 10, got %v",
			articles[0].Metadata["comments"],
		)
	}
}
