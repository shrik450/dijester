package processor

import (
	"bytes"
	"errors"
	"net/url"
	"strings"

	"github.com/go-shiori/go-readability"
	"golang.org/x/net/html"

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

// removeHTMLTags uses proper HTML parsing to remove all instances of a given tag from HTML content.
func removeHTMLTags(content, tagName string) string {
	// Parse the HTML document
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		// If parsing fails, return the original content
		return content
	}

	// Create a buffer for the modified HTML
	var buf bytes.Buffer

	// Remove tags by traversing the DOM
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		// If this is the tag we want to remove, skip it entirely
		if n.Type == html.ElementNode && strings.EqualFold(n.Data, tagName) {
			return
		}

		// For all other nodes, render them
		if n.Type == html.ElementNode {
			buf.WriteByte('<')
			buf.WriteString(n.Data)
			for _, attr := range n.Attr {
				buf.WriteByte(' ')
				buf.WriteString(attr.Key)
				buf.WriteString(`="`)
				buf.WriteString(html.EscapeString(attr.Val))
				buf.WriteByte('"')
			}
			buf.WriteByte('>')

			// Process all children
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				traverse(c)
			}

			// Close the tag
			buf.WriteString("</")
			buf.WriteString(n.Data)
			buf.WriteByte('>')
		} else if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		} else if n.Type == html.CommentNode {
			buf.WriteString("<!--")
			buf.WriteString(n.Data)
			buf.WriteString("-->")
		} else if n.Type == html.DoctypeNode {
			buf.WriteString("<!DOCTYPE ")
			buf.WriteString(n.Data)
			buf.WriteByte('>')
		}
	}

	// Process all top-level nodes
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		traverse(c)
	}

	result := buf.String()
	return result
}
