package hackernews

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shrik450/dijester/pkg/models"
	"github.com/shrik450/dijester/pkg/processor"
	"github.com/shrik450/dijester/pkg/source"
)

const (
	apiBaseURL    = "https://hacker-news.firebaseio.com/v0"
	topStoriesURL = apiBaseURL + "/topstories.json"
	itemURLFormat = apiBaseURL + "/item/%d.json"
)

// APIFetcher defines the interface for fetching from the API
type APIFetcher interface {
	FetchURLAsString(ctx context.Context, url string) (string, error)
}

// Source implements a Hacker News source.
type Source struct {
	name        string
	maxArticles int
	minScore    int
	categories  []string
	fetcher     APIFetcher
	processor   processor.Processor
	procOptions *processor.Options
}

// ensure Source implements the source.Source interface
var _ source.Source = (*Source)(nil)

// New creates a new Hacker News source.
func New(fetcher APIFetcher) *Source {
	return &Source{
		name:        "hackernews",
		maxArticles: 10,
		minScore:    100,
		categories:  []string{"front_page"},
		fetcher:     fetcher,
		processor:   processor.NewReadabilityProcessor(),
		procOptions: processor.DefaultOptions(),
	}
}

// Name returns the source name.
func (s *Source) Name() string {
	return s.name
}

// Configure sets up the source with the provided configuration.
func (s *Source) Configure(config map[string]interface{}) error {
	if name, ok := config["name"].(string); ok && name != "" {
		s.name = name
	}

	if max, ok := config["max_articles"].(int); ok && max > 0 {
		s.maxArticles = max
	}

	if score, ok := config["min_score"].(int); ok && score > 0 {
		s.minScore = score
	}

	if cats, ok := config["categories"].([]interface{}); ok && len(cats) > 0 {
		s.categories = make([]string, 0, len(cats))
		for _, cat := range cats {
			if catStr, ok := cat.(string); ok {
				s.categories = append(s.categories, catStr)
			}
		}
	}

	return nil
}

// HNItem represents a Hacker News API item.
type HNItem struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Text        string `json:"text"`
	By          string `json:"by"`
	Score       int    `json:"score"`
	Time        int64  `json:"time"`
	Type        string `json:"type"`
	Kids        []int  `json:"kids"`
	Dead        bool   `json:"dead"`
	Deleted     bool   `json:"deleted"`
	Descendants int    `json:"descendants"`
}

// Fetch retrieves articles from Hacker News.
func (s *Source) Fetch(ctx context.Context) ([]*models.Article, error) {
	content, err := s.fetcher.FetchURLAsString(ctx, topStoriesURL)
	if err != nil {
		return nil, fmt.Errorf("fetching top stories: %w", err)
	}

	var storyIDs []int
	if err := json.Unmarshal([]byte(content), &storyIDs); err != nil {
		return nil, fmt.Errorf("parsing story IDs: %w", err)
	}

	if len(storyIDs) > s.maxArticles*2 {
		storyIDs = storyIDs[:s.maxArticles*2]
	}

	articles := make([]*models.Article, 0, s.maxArticles)
	for _, id := range storyIDs {
		if len(articles) >= s.maxArticles {
			break
		}

		itemURL := fmt.Sprintf(itemURLFormat, id)
		content, err := s.fetcher.FetchURLAsString(ctx, itemURL)
		if err != nil {
			continue
		}

		var item HNItem
		if err := json.Unmarshal([]byte(content), &item); err != nil {
			continue
		}

		if item.Deleted || item.Dead {
			continue
		}

		if item.Score < s.minScore {
			continue
		}

		if item.Type != "story" {
			continue
		}

		article := &models.Article{
			Title:       item.Title,
			Author:      item.By,
			PublishedAt: time.Unix(item.Time, 0),
			URL:         item.URL,
			Content:     item.Text,
			SourceName:  s.name,
			Metadata: map[string]interface{}{
				"score":        item.Score,
				"comments":     item.Descendants,
				"id":           item.ID,
				"comments_url": fmt.Sprintf("https://news.ycombinator.com/item?id=%d", item.ID),
			},
		}

		if article.URL == "" {
			article.URL = article.Metadata["comments_url"].(string)
		} else {
			articleContent, err := s.fetcher.FetchURLAsString(ctx, article.URL)
			if err == nil {
				article.Content = articleContent

				if err := s.processor.Process(article, s.procOptions); err != nil {
					if article.Content == "" && item.Text != "" {
						article.Content = item.Text
					}
				}
			} else if item.Text != "" {
				article.Content = item.Text
			}
		}

		article.Summary = fmt.Sprintf("%d points, %d comments",
			item.Score, item.Descendants)

		articles = append(articles, article)
	}

	return articles, nil
}
