package hackernews

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/shrik450/dijester/pkg/fetcher"
	"github.com/shrik450/dijester/pkg/models"
)

const (
	apiBaseURL    = "https://hacker-news.firebaseio.com/v0"
	itemURLFormat = apiBaseURL + "/item/%d.json"
	hnBaseURL     = "https://news.ycombinator.com"
)

// PageType represents available HN page types
type PageType string

const (
	FrontPage  PageType = "frontpage"
	NewPage    PageType = "new"
	PastPage   PageType = "past"
	ActivePage PageType = "active"
	JobsPage   PageType = "jobs"
)

// HNPage maps page types to their respective URLs
var HNPageURLs = map[PageType]string{
	FrontPage:  hnBaseURL,
	NewPage:    hnBaseURL + "/newest",
	PastPage:   hnBaseURL + "/front",
	ActivePage: hnBaseURL + "/active",
	JobsPage:   hnBaseURL + "/jobs",
}

// Source implements a Hacker News source
type Source struct {
	name        string
	maxArticles int
	minScore    int
	showDead    bool
	showDeleted bool
	pageType    PageType
}

// New creates a new Hacker News source with default settings
func New() *Source {
	return &Source{
		name:        "hackernews",
		maxArticles: 30,
		minScore:    10,
		pageType:    FrontPage,
	}
}

// Name returns the source name
func (s *Source) Name() string {
	return s.name
}

// Configure sets up the source with the provided configuration
func (s *Source) Configure(config map[string]any) error {
	if name, ok := config["name"].(string); ok && name != "" {
		s.name = name
	}

	if max, ok := config["max_articles"].(int); ok && max > 0 {
		s.maxArticles = max
	}

	if score, ok := config["min_score"].(int); ok && score > 0 {
		s.minScore = score
	}

	if showDead, ok := config["show_dead"].(bool); ok {
		s.showDead = showDead
	}

	if showDeleted, ok := config["show_deleted"].(bool); ok {
		s.showDeleted = showDeleted
	}

	if page, ok := config["page"].(string); ok && page != "" {
		switch strings.ToLower(page) {
		case "frontpage", "front":
			s.pageType = FrontPage
		case "new", "newest":
			s.pageType = NewPage
		case "past":
			s.pageType = PastPage
		case "active":
			s.pageType = ActivePage
		case "jobs":
			s.pageType = JobsPage
		default:
			log.Printf("Unknown HN page type '%s', defaulting to frontpage", page)
			s.pageType = FrontPage
		}
	}

	return nil
}

// HNItem represents a Hacker News API item
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

// StoryItem represents a story parsed from the HN HTML page
type StoryItem struct {
	ID        int
	Title     string
	URL       string
	By        string
	Score     int
	Comments  int
	TimeAgo   string
	Timestamp time.Time
}

// parseHNPage parses the HTML content of a Hacker News page to extract stories
func parseHNPage(content string) ([]*StoryItem, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("parsing HN page: %w", err)
	}

	stories := make([]*StoryItem, 0)
	doc.Find("tr.athing").Each(func(i int, s *goquery.Selection) {
		story := &StoryItem{}
		idStr, exists := s.Attr("id")
		if !exists {
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return
		}
		story.ID = id

		titleElement := s.Find("td.title > span.titleline > a").First()
		story.Title = strings.TrimSpace(titleElement.Text())
		story.URL, _ = titleElement.Attr("href")

		if story.URL != "" && !strings.HasPrefix(story.URL, "http") &&
			!strings.HasPrefix(story.URL, "//") {
			story.URL = hnBaseURL + "/" + strings.TrimPrefix(story.URL, "/")
		}

		// The metadata is in the next row
		metaRow := s.Next()
		if metaRow.Length() == 0 {
			return
		}

		scoreText := metaRow.Find("span.score").Text()
		if scoreText != "" {
			re := regexp.MustCompile(`(\d+)`)
			matches := re.FindStringSubmatch(scoreText)
			if len(matches) > 1 {
				story.Score, _ = strconv.Atoi(matches[1])
			}
		}

		story.By = strings.TrimSpace(metaRow.Find("a.hnuser").Text())

		commentLink := metaRow.Find("a").FilterFunction(func(i int, s *goquery.Selection) bool {
			return strings.Contains(s.Text(), "comment") || strings.Contains(s.Text(), "discuss")
		})
		if commentLink.Length() > 0 {
			commentText := commentLink.Text()
			re := regexp.MustCompile(`(\d+)`)
			matches := re.FindStringSubmatch(commentText)
			if len(matches) > 1 {
				story.Comments, _ = strconv.Atoi(matches[1])
			}
		}

		story.TimeAgo = strings.TrimSpace(metaRow.Find("span.age").Text())

		stories = append(stories, story)
	})

	return stories, nil
}

// Fetch retrieves articles from Hacker News
func (s *Source) Fetch(ctx context.Context, fetcher fetcher.Fetcher) ([]*models.Article, error) {
	pageURL, ok := HNPageURLs[s.pageType]
	if !ok {
		pageURL = HNPageURLs[FrontPage]
	}

	content, err := fetcher.FetchURLAsString(ctx, pageURL)
	if err != nil {
		return nil, fmt.Errorf("fetching HN page %s: %w", s.pageType, err)
	}

	stories, err := parseHNPage(content)
	if err != nil {
		return nil, fmt.Errorf("parsing HN page: %w", err)
	}

	articles := make([]*models.Article, 0, s.maxArticles)
	for _, story := range stories {
		if len(articles) >= s.maxArticles {
			break
		}

		if s.pageType != JobsPage && story.Score < s.minScore {
			continue
		}

		itemURL := fmt.Sprintf(itemURLFormat, story.ID)
		itemContent, err := fetcher.FetchURLAsString(ctx, itemURL)

		var item HNItem
		var useAPIItem bool

		if err == nil {
			if err := json.Unmarshal([]byte(itemContent), &item); err == nil {
				useAPIItem = true

				if !s.showDead && item.Dead {
					continue
				}

				if !s.showDeleted && item.Deleted {
					continue
				}
			}
		}

		article := &models.Article{
			Title:      story.Title,
			Author:     story.By,
			URL:        story.URL,
			SourceName: s.name,
			Metadata: map[string]any{
				"score":        story.Score,
				"comments":     story.Comments,
				"id":           story.ID,
				"comments_url": fmt.Sprintf("https://news.ycombinator.com/item?id=%d", story.ID),
			},
		}

		if useAPIItem {
			article.Title = item.Title
			article.Author = item.By
			article.PublishedAt = time.Unix(item.Time, 0)
			article.URL = item.URL
			article.Content = item.Text
			article.Metadata["score"] = item.Score
			article.Metadata["comments"] = item.Descendants
		}

		if article.URL == "" {
			article.URL = article.Metadata["comments_url"].(string)
		} else if article.Content == "" && article.URL != "" {
			articleContent, err := fetcher.FetchURLAsString(ctx, article.URL)
			if err == nil && articleContent != "" {
				article.Content = articleContent
			}
		}

		article.Summary = fmt.Sprintf("%d points, %d comments",
			article.Metadata["score"].(int), article.Metadata["comments"].(int))

		articles = append(articles, article)
	}

	return articles, nil
}
