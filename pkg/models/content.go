package models

import (
	"reflect"
	"sort"
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
func ParseSortField(fieldStr string) SortField {
	field := SortField{
		Direction: "asc",
	}

	if strings.Contains(fieldStr, ":") {
		parts := strings.Split(fieldStr, ":")
		field.Name = parts[0]
		if len(parts) > 1 && (parts[1] == "asc" || parts[1] == "desc") {
			field.Direction = parts[1]
		}
	} else {
		field.Name = fieldStr
	}

	return field
}

// SortArticles sorts a slice of articles based on the provided sort fields.
// Each sort field can be in the format "FieldName" or "FieldName:direction"
// where direction is either "asc" or "desc".
func SortArticles(articles []*Article, sortFields []string) {
	if len(articles) <= 1 || len(sortFields) == 0 {
		return
	}

	parsedFields := make([]SortField, 0, len(sortFields))
	for _, field := range sortFields {
		parsedFields = append(parsedFields, ParseSortField(field))
	}

	compareArticles := func(i, j int) bool {
		for _, field := range parsedFields {
			iVal := reflect.ValueOf(*articles[i]).FieldByName(field.Name)
			jVal := reflect.ValueOf(*articles[j]).FieldByName(field.Name)

			if !iVal.IsValid() || !jVal.IsValid() {
				continue
			}

			var isLess bool

			switch iVal.Kind() {
			case reflect.String:
				isLess = iVal.String() < jVal.String()
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				isLess = iVal.Int() < jVal.Int()
			case reflect.Float32, reflect.Float64:
				isLess = iVal.Float() < jVal.Float()
			case reflect.Struct:
				// Special handling for time.Time
				if iVal.Type() == reflect.TypeOf(time.Time{}) {
					iTime := iVal.Interface().(time.Time)
					jTime := jVal.Interface().(time.Time)
					isLess = iTime.Before(jTime)
				} else {
					continue
				}
			default:
				continue
			}

			// If we have different values, return based on sort direction
			if iVal.Interface() != jVal.Interface() {
				if field.Direction == "desc" {
					return !isLess
				}
				return isLess
			}
			// If values are equal, continue to next sort field
		}

		// If all fields compared are equal
		return false
	}

	sort.Slice(articles, compareArticles)
}
