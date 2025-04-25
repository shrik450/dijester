package models

import (
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
