package models

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"time"
)

// Article represents a single piece of content from a source.
type Article struct {
	// Title of the article
	Title string

	// Author of the article, if available
	Author string

	// PublishedAt is the publication date of the article
	PublishedAt time.Time

	// URL is the original source URL of the article
	URL string

	// Content is the full text content of the article
	Content string

	// Summary is a short summary or description of the article
	Summary string

	// SourceName identifies which source this article came from
	SourceName string

	// Tags are optional categorization tags for the article
	Tags []string

	// Metadata contains any additional source-specific metadata
	Metadata map[string]any
}

// Digest represents a collection of articles ready for formatting.
type Digest struct {
	// Title of the digest
	Title string

	// GeneratedAt is when this digest was created
	GeneratedAt time.Time

	// Articles is the collection of articles in this digest
	Articles []*Article

	// Metadata contains any additional digest-specific metadata
	Metadata map[string]interface{}
}

var allowedSortFields = []string{
	"Title",
	"Author",
	"PublishedAt",
	"URL",
	"Content",
	"Summary",
	"SourceName",
}

// DeduplicateArticlesByURL removes duplicate articles with the same URL.
// Returns a new slice containing only unique articles, keeping the first
// occurrence of each URL.
func DeduplicateArticlesByURL(articles []*Article) []*Article {
	if len(articles) == 0 {
		return articles
	}

	urlMap := make(map[string]bool)
	deduped := make([]*Article, 0, len(articles))

	for _, article := range articles {
		if !urlMap[article.URL] {
			urlMap[article.URL] = true
			deduped = append(deduped, article)
		}
	}

	return deduped
}

// SortField represents a field to sort articles by, along with its direction.
type SortField struct {
	// Name is the name of the field to sort by
	Name string

	// Direction is either "asc" or "desc"
	Direction string
}

// ParseSortField parses a sort field string in the format "FieldName:direction".
// If no direction is specified, defaults to "asc".
// Returns a SortField struct.
func ParseSortField(fieldStr string) (SortField, error) {
	field := SortField{
		Direction: "asc",
	}

	if strings.Contains(fieldStr, ":") {
		parts := strings.Split(fieldStr, ":")
		field.Name = parts[0]
		if len(parts) > 2 {
			return SortField{}, fmt.Errorf("invalid sort field format: %s", fieldStr)
		}
		if len(parts) == 2 {
			field.Direction = parts[1]
		}
	} else {
		field.Name = fieldStr
	}

	if field.Name == "" || !slices.Contains(allowedSortFields, field.Name) {
		return SortField{}, fmt.Errorf("invalid sort field: %s", field.Name)
	}

	if field.Direction != "asc" && field.Direction != "desc" {
		return SortField{}, fmt.Errorf("invalid sort direction: %s", field.Direction)
	}

	return field, nil
}

// SortArticles sorts a slice of articles based on the provided sort fields.
// Each sort field can be in the format "FieldName" or "FieldName:direction"
// where direction is either "asc" or "desc".
func SortArticles(articles []*Article, sortFields []string) error {
	if len(articles) <= 1 || len(sortFields) == 0 {
		return nil
	}

	parsedFields := make([]SortField, len(sortFields))
	for i, fieldStr := range sortFields {
		field, err := ParseSortField(fieldStr)
		if err != nil {
			return fmt.Errorf("error parsing sort field %s: %w", fieldStr, err)
		}
		parsedFields[i] = field
	}

	compareArticles := func(a, b *Article) int {
		for _, field := range parsedFields {
			cmpVal := cmpByField(a, b, field.Name)

			if cmpVal != 0 {
				if field.Direction == "desc" {
					return -cmpVal
				}
				return cmpVal
			}
		}

		return 0
	}

	slices.SortStableFunc(articles, compareArticles)

	return nil
}

func cmpByField(a, b *Article, fieldName string) int {
	switch fieldName {
	case "Title":
		return cmp.Compare(a.Title, b.Title)
	case "Author":
		return cmp.Compare(a.Author, b.Author)
	case "PublishedAt":
		return cmp.Compare(a.PublishedAt.Unix(), b.PublishedAt.Unix())
	case "URL":
		return cmp.Compare(a.URL, b.URL)
	case "Content":
		return cmp.Compare(a.Content, b.Content)
	case "Summary":
		return cmp.Compare(a.Summary, b.Summary)
	case "SourceName":
		return cmp.Compare(a.SourceName, b.SourceName)
	default:
		return 0
	}
}
