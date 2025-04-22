package rss

import (
	"context"
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/shrik450/dijester/pkg/models"
	"github.com/shrik450/dijester/pkg/source"
)

// FeedFetcher defines the interface for fetching feeds
type FeedFetcher interface {
	FetchURLAsString(ctx context.Context, url string) (string, error)
}

// Source implements an RSS/Atom feed source.
type Source struct {
	name        string
	url         string
	maxArticles int
	fetcher     FeedFetcher
	parser      *gofeed.Parser
}

// ensure Source implements the source.Source interface
var _ source.Source = (*Source)(nil)

// New creates a new RSS source.
func New(fetcher FeedFetcher) *Source {
	return &Source{
		name:        "rss",
		maxArticles: 10,
		fetcher:     fetcher,
		parser:      gofeed.NewParser(),
	}
}

// Name returns the source name.
func (s *Source) Name() string {
	return s.name
}

// Configure sets up the source with the provided configuration.
func (s *Source) Configure(config map[string]interface{}) error {
	// Set custom name if provided
	if name, ok := config["name"].(string); ok && name != "" {
		s.name = name
	}

	url, ok := config["url"].(string)
	if !ok || url == "" {
		return fmt.Errorf("rss source requires a 'url' configuration value")
	}
	s.url = url

	if max, ok := config["max_articles"].(int); ok && max > 0 {
		s.maxArticles = max
	}

	return nil
}

// Fetch retrieves articles from the RSS feed.
func (s *Source) Fetch(ctx context.Context) ([]*models.Article, error) {
	content, err := s.fetcher.FetchURLAsString(ctx, s.url)
	if err != nil {
		return nil, fmt.Errorf("fetching RSS feed: %w", err)
	}

	feed, err := s.parser.ParseString(content)
	if err != nil {
		return nil, fmt.Errorf("parsing RSS feed: %w", err)
	}

	articles := make([]*models.Article, 0, len(feed.Items))
	for _, item := range feed.Items {
		// Skip items without content
		if item.Content == "" && item.Description == "" {
			continue
		}

		// Use published date if available, otherwise use updated date
		var publishedAt time.Time
		if item.PublishedParsed != nil {
			publishedAt = *item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			publishedAt = *item.UpdatedParsed
		} else {
			publishedAt = time.Now()
		}

		content := item.Content
		summary := item.Description

		if content == "" {
			content = item.Description
			summary = ""
		}

		author := ""
		if item.Author != nil {
			author = item.Author.Name
		}

		article := &models.Article{
			Title:       item.Title,
			Author:      author,
			PublishedAt: publishedAt,
			URL:         item.Link,
			Content:     content,
			Summary:     summary,
			SourceName:  s.name,
			Tags:        item.Categories,
			Metadata:    make(map[string]interface{}),
		}

		articles = append(articles, article)

		if len(articles) >= s.maxArticles {
			break
		}
	}

	return articles, nil
}
