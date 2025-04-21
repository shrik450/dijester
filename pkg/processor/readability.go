package processor

import (
	"bytes"
	"errors"
	"net/url"
	"strings"

	"github.com/go-shiori/go-readability"

	"github.com/shrik450/dijester/pkg/models"
)

// ReadabilityProcessor uses go-readability to extract the main content from HTML.
type ReadabilityProcessor struct{}

// NewReadabilityProcessor creates a new readability-based content processor.
func NewReadabilityProcessor() *ReadabilityProcessor {
	return &ReadabilityProcessor{}
}

// Name returns the name of this processor.
func (p *ReadabilityProcessor) Name() string {
	return "readability"
}

// Process extracts the main content from an article's HTML content.
func (p *ReadabilityProcessor) Process(article *models.Article, opts *Options) error {
	if article == nil {
		return errors.New("article cannot be nil")
	}

	if article.Content == "" {
		return errors.New("article content cannot be empty")
	}

	// Parse the content using readability
	parser := readability.NewParser()

	// Set parser options based on our options
	if opts != nil {
		if v, ok := opts.AdditionalOptions["classesToPreserve"]; ok {
			if classes, ok := v.([]string); ok {
				parser.ClassesToPreserve = classes
			}
		}
	}

	// Parse URL for readability
	articleURL, err := url.Parse(article.URL)
	if err != nil {
		return err
	}

	// Parse the content
	result, err := parser.Parse(strings.NewReader(article.Content), articleURL)
	if err != nil {
		return err
	}

	// Apply length constraints if specified
	content := result.Content
	if opts != nil {
		if opts.MinContentLength > 0 && len(content) < opts.MinContentLength {
			return ErrContentProcessingFailed
		}

		if opts.MaxContentLength > 0 && len(content) > opts.MaxContentLength {
			content = content[:opts.MaxContentLength]
		}

		// Filter out images if not wanted
		if !opts.IncludeImages {
			content = removeHTMLTags(content, "img")
		}

		// Filter out tables if not wanted
		if !opts.IncludeTables {
			content = removeHTMLTags(content, "table")
		}
	}

	// Update the article with processed content
	article.Content = content

	// Extract or update article metadata if not already set
	if article.Title == "" {
		article.Title = result.Title
	}

	if article.Author == "" && len(result.Byline) > 0 {
		article.Author = result.Byline
	}

	if article.Summary == "" && len(result.Excerpt) > 0 {
		article.Summary = result.Excerpt
	}

	return nil
}

// removeHTMLTags removes specified HTML tags from content.
// This is a simplified implementation - in production, use a proper HTML parser.
func removeHTMLTags(content, tagName string) string {
	openTag := "<" + tagName
	closeTag := "</" + tagName + ">"

	var result bytes.Buffer
	remainder := content

	for {
		openIndex := strings.Index(remainder, openTag)
		if openIndex == -1 {
			// No more tags to remove
			result.WriteString(remainder)
			break
		}

		// Write everything before the tag
		result.WriteString(remainder[:openIndex])

		// Find the closing tag
		remainder = remainder[openIndex:]
		closeIndex := strings.Index(remainder, closeTag)
		if closeIndex == -1 {
			// No closing tag found, just remove the opening tag
			remainder = strings.Replace(remainder, openTag, "", 1)
		} else {
			// Skip to after the closing tag
			remainder = remainder[closeIndex+len(closeTag):]
		}
	}

	return result.String()
}
