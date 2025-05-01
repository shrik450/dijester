package models

import (
	"testing"
	"time"
)

func TestDeduplicateArticlesByURL(t *testing.T) {
	tests := []struct {
		name     string
		articles []*Article
		want     []*Article
	}{
		{
			name:     "empty slice",
			articles: []*Article{},
			want:     []*Article{},
		},
		{
			name: "no duplicates",
			articles: []*Article{
				{URL: "https://example.com/1", Title: "Article 1"},
				{URL: "https://example.com/2", Title: "Article 2"},
				{URL: "https://example.com/3", Title: "Article 3"},
			},
			want: []*Article{
				{URL: "https://example.com/1", Title: "Article 1"},
				{URL: "https://example.com/2", Title: "Article 2"},
				{URL: "https://example.com/3", Title: "Article 3"},
			},
		},
		{
			name: "with duplicates",
			articles: []*Article{
				{URL: "https://example.com/1", Title: "Article 1"},
				{URL: "https://example.com/2", Title: "Article 2"},
				{URL: "https://example.com/1", Title: "Article 1 - Duplicate"},
				{URL: "https://example.com/3", Title: "Article 3"},
				{URL: "https://example.com/2", Title: "Article 2 - Duplicate"},
			},
			want: []*Article{
				{URL: "https://example.com/1", Title: "Article 1"},
				{URL: "https://example.com/2", Title: "Article 2"},
				{URL: "https://example.com/3", Title: "Article 3"},
			},
		},
		{
			name: "all duplicates",
			articles: []*Article{
				{URL: "https://example.com/1", Title: "Article 1"},
				{URL: "https://example.com/1", Title: "Article 1 - Duplicate 1"},
				{URL: "https://example.com/1", Title: "Article 1 - Duplicate 2"},
			},
			want: []*Article{
				{URL: "https://example.com/1", Title: "Article 1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeduplicateArticlesByURL(tt.articles)

			if len(got) != len(tt.want) {
				t.Errorf(
					"DeduplicateArticlesByURL() returned %d articles, want %d",
					len(got),
					len(tt.want),
				)
				return
			}

			// Check that all URLs in the result are unique
			urlMap := make(map[string]bool)
			for _, article := range got {
				if urlMap[article.URL] {
					t.Errorf("DeduplicateArticlesByURL() returned duplicate URL: %s", article.URL)
				}
				urlMap[article.URL] = true
			}

			// Check that all expected URLs are present
			for _, wantArticle := range tt.want {
				found := false
				for _, gotArticle := range got {
					if gotArticle.URL == wantArticle.URL {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("DeduplicateArticlesByURL() is missing URL: %s", wantArticle.URL)
				}
			}
		})
	}
}

func TestParseSortField(t *testing.T) {
	tests := []struct {
		name     string
		fieldStr string
		want     SortField
	}{
		{
			name:     "field only",
			fieldStr: "Title",
			want:     SortField{Name: "Title", Direction: "asc"},
		},
		{
			name:     "field with asc direction",
			fieldStr: "Title:asc",
			want:     SortField{Name: "Title", Direction: "asc"},
		},
		{
			name:     "field with desc direction",
			fieldStr: "PublishedAt:desc",
			want:     SortField{Name: "PublishedAt", Direction: "desc"},
		},
		{
			name:     "field with invalid direction defaults to asc",
			fieldStr: "Content:invalid",
			want:     SortField{Name: "Content", Direction: "asc"},
		},
		{
			name:     "multiple colons take first part as field name",
			fieldStr: "URL:desc:extra",
			want:     SortField{Name: "URL", Direction: "desc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseSortField(tt.fieldStr)
			if got.Name != tt.want.Name || got.Direction != tt.want.Direction {
				t.Errorf("ParseSortField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortArticles(t *testing.T) {
	now := time.Now()
	oneHourAgo := now.Add(-1 * time.Hour)
	twoHoursAgo := now.Add(-2 * time.Hour)

	tests := []struct {
		name       string
		articles   []*Article
		sortFields []string
		want       []*Article
	}{
		{
			name:       "empty slice",
			articles:   []*Article{},
			sortFields: []string{"Title"},
			want:       []*Article{},
		},
		{
			name: "sort by title ascending",
			articles: []*Article{
				{Title: "C Article", URL: "https://example.com/3"},
				{Title: "A Article", URL: "https://example.com/1"},
				{Title: "B Article", URL: "https://example.com/2"},
			},
			sortFields: []string{"Title"},
			want: []*Article{
				{Title: "A Article", URL: "https://example.com/1"},
				{Title: "B Article", URL: "https://example.com/2"},
				{Title: "C Article", URL: "https://example.com/3"},
			},
		},
		{
			name: "sort by title descending",
			articles: []*Article{
				{Title: "C Article", URL: "https://example.com/3"},
				{Title: "A Article", URL: "https://example.com/1"},
				{Title: "B Article", URL: "https://example.com/2"},
			},
			sortFields: []string{"Title:desc"},
			want: []*Article{
				{Title: "C Article", URL: "https://example.com/3"},
				{Title: "B Article", URL: "https://example.com/2"},
				{Title: "A Article", URL: "https://example.com/1"},
			},
		},
		{
			name: "sort by publishedAt",
			articles: []*Article{
				{Title: "Article 1", PublishedAt: now, URL: "https://example.com/1"},
				{Title: "Article 2", PublishedAt: twoHoursAgo, URL: "https://example.com/2"},
				{Title: "Article 3", PublishedAt: oneHourAgo, URL: "https://example.com/3"},
			},
			sortFields: []string{"PublishedAt"},
			want: []*Article{
				{Title: "Article 2", PublishedAt: twoHoursAgo, URL: "https://example.com/2"},
				{Title: "Article 3", PublishedAt: oneHourAgo, URL: "https://example.com/3"},
				{Title: "Article 1", PublishedAt: now, URL: "https://example.com/1"},
			},
		},
		{
			name: "multiple sort fields",
			articles: []*Article{
				{Title: "B Article", Author: "Author 1", URL: "https://example.com/1"},
				{Title: "A Article", Author: "Author 2", URL: "https://example.com/2"},
				{Title: "A Article", Author: "Author 1", URL: "https://example.com/3"},
			},
			sortFields: []string{"Title", "Author"},
			want: []*Article{
				{Title: "A Article", Author: "Author 1", URL: "https://example.com/3"},
				{Title: "A Article", Author: "Author 2", URL: "https://example.com/2"},
				{Title: "B Article", Author: "Author 1", URL: "https://example.com/1"},
			},
		},
		{
			name: "multiple sort fields with directions",
			articles: []*Article{
				{Title: "B Article", Author: "Author 1", URL: "https://example.com/1"},
				{Title: "A Article", Author: "Author 2", URL: "https://example.com/2"},
				{Title: "A Article", Author: "Author 1", URL: "https://example.com/3"},
			},
			sortFields: []string{"Title:asc", "Author:desc"},
			want: []*Article{
				{Title: "A Article", Author: "Author 2", URL: "https://example.com/2"},
				{Title: "A Article", Author: "Author 1", URL: "https://example.com/3"},
				{Title: "B Article", Author: "Author 1", URL: "https://example.com/1"},
			},
		},
		{
			name: "invalid field is ignored",
			articles: []*Article{
				{Title: "C Article", URL: "https://example.com/3"},
				{Title: "A Article", URL: "https://example.com/1"},
				{Title: "B Article", URL: "https://example.com/2"},
			},
			sortFields: []string{"NonExistentField", "Title"},
			want: []*Article{
				{Title: "A Article", URL: "https://example.com/1"},
				{Title: "B Article", URL: "https://example.com/2"},
				{Title: "C Article", URL: "https://example.com/3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the input articles to avoid modifying the test data
			articlesCopy := make([]*Article, len(tt.articles))
			for i, article := range tt.articles {
				articleCopy := *article
				articlesCopy[i] = &articleCopy
			}

			SortArticles(articlesCopy, tt.sortFields)

			if len(articlesCopy) != len(tt.want) {
				t.Errorf(
					"SortArticles() resulted in %d articles, want %d",
					len(articlesCopy),
					len(tt.want),
				)
				return
			}

			for i := range articlesCopy {
				if articlesCopy[i].URL != tt.want[i].URL {
					t.Errorf(
						"SortArticles() at index %d: got URL %s, want %s",
						i,
						articlesCopy[i].URL,
						tt.want[i].URL,
					)
				}
				if articlesCopy[i].Title != tt.want[i].Title {
					t.Errorf(
						"SortArticles() at index %d: got Title %s, want %s",
						i,
						articlesCopy[i].Title,
						tt.want[i].Title,
					)
				}
			}
		})
	}
}
