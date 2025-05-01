package processor

import (
	"errors"

	"github.com/microcosm-cc/bluemonday"

	"github.com/shrik450/dijester/pkg/models"
)

// SanitizerProcessor is a processor that sanitizes HTML content.
type SanitizerProcessor struct{}

// NewSanitizerProcessor creates a new instance of SanitizerProcessor.
func NewSanitizerProcessor() *SanitizerProcessor {
	return &SanitizerProcessor{}
}

func (s *SanitizerProcessor) Name() string {
	return "sanitizer"
}

// Process sanitizes the HTML content of the article.
func (p *SanitizerProcessor) Process(article *models.Article, opts *Options) error {
	if article == nil {
		return errors.New("article cannot be nil")
	}

	policy := bluemonday.UGCPolicy()
	policy.AllowDataURIImages()

	if article.Content == "" {
		return nil
	}

	sanitizedContent := policy.Sanitize(article.Content)

	article.Content = sanitizedContent
	return nil
}
