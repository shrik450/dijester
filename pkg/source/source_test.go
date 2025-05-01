package source

import (
	"strings"
	"testing"

	"github.com/shrik450/dijester/pkg/models"
)

func TestFilterArticlesByWordDenylist(t *testing.T) {
	tests := []struct {
		name     string
		articles []*models.Article
		denylist []string
		want     int
	}{
		{
			name:     "empty article list",
			articles: []*models.Article{},
			denylist: []string{"test", "example"},
			want:     0,
		},
		{
			name: "empty denylist",
			articles: []*models.Article{
				{Title: "Test Article", Content: "This is a test article", Summary: "Test summary"},
				{
					Title:   "Example Article",
					Content: "This is an example article",
					Summary: "Example summary",
				},
			},
			denylist: []string{},
			want:     2,
		},
		{
			name: "no matches",
			articles: []*models.Article{
				{
					Title:   "Hello Article",
					Content: "This is a good article",
					Summary: "Good summary",
				},
				{
					Title:   "World Article",
					Content: "This is a nice article",
					Summary: "Nice summary",
				},
			},
			denylist: []string{"test", "example", "bad"},
			want:     2,
		},
		{
			name: "match in title",
			articles: []*models.Article{
				{Title: "Test Article", Content: "This is a good article", Summary: "Good summary"},
				{
					Title:   "World Article",
					Content: "This is a nice article",
					Summary: "Nice summary",
				},
			},
			denylist: []string{"test", "example"},
			want:     1,
		},
		{
			name: "match in content",
			articles: []*models.Article{
				{
					Title:   "Hello Article",
					Content: "This is a test article",
					Summary: "Good summary",
				},
				{
					Title:   "World Article",
					Content: "This is a nice article",
					Summary: "Nice summary",
				},
			},
			denylist: []string{"test", "example"},
			want:     1,
		},
		{
			name: "match in summary",
			articles: []*models.Article{
				{
					Title:   "Hello Article",
					Content: "This is a good article",
					Summary: "Test summary",
				},
				{
					Title:   "World Article",
					Content: "This is a nice article",
					Summary: "Nice summary",
				},
			},
			denylist: []string{"test", "example"},
			want:     1,
		},
		{
			name: "case insensitive match",
			articles: []*models.Article{
				{
					Title:   "Hello Article",
					Content: "This is a TEST article",
					Summary: "Good summary",
				},
				{
					Title:   "World Article",
					Content: "This is a nice article",
					Summary: "Nice summary",
				},
			},
			denylist: []string{"test", "example"},
			want:     1,
		},
		{
			name: "multiple matches same article",
			articles: []*models.Article{
				{
					Title:   "Test Article",
					Content: "This is a test example",
					Summary: "Example summary",
				},
				{
					Title:   "World Article",
					Content: "This is a nice article",
					Summary: "Nice summary",
				},
			},
			denylist: []string{"test", "example"},
			want:     1,
		},
		{
			name: "all articles filtered",
			articles: []*models.Article{
				{Title: "Test Article", Content: "This is a test article", Summary: "Test summary"},
				{
					Title:   "Example Article",
					Content: "This is an example article",
					Summary: "Example summary",
				},
			},
			denylist: []string{"test", "example"},
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := FilterArticlesByWordDenylist(tt.articles, tt.denylist)

			if len(filtered) != tt.want {
				t.Errorf(
					"FilterArticlesByWordDenylist() returned %d articles, want %d",
					len(filtered),
					tt.want,
				)
			}

			// Verify no articles contain denylisted words
			for _, article := range filtered {
				for _, denyWord := range tt.denylist {
					if containsDenylistedWord(article, denyWord) {
						t.Errorf(
							"FilterArticlesByWordDenylist() returned article with denylisted word: %s",
							denyWord,
						)
					}
				}
			}
		})
	}
}

func containsDenylistedWord(article *models.Article, denyWord string) bool {
	denyWord = strings.ToLower(denyWord)
	return strings.Contains(strings.ToLower(article.Title), denyWord) ||
		strings.Contains(strings.ToLower(article.Content), denyWord) ||
		strings.Contains(strings.ToLower(article.Summary), denyWord)
}
