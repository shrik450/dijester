package processor

import (
	"bytes"
	"errors"
	"log"
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
		return nil
	}

	parser := readability.NewParser()

	if opts != nil {
		if v, ok := opts.AdditionalOptions["classesToPreserve"]; ok {
			if classes, ok := v.([]string); ok {
				parser.ClassesToPreserve = classes
			}
		}
	}

	articleURL, err := url.Parse(article.URL)
	if err != nil {
		return err
	}

	result, err := parser.Parse(strings.NewReader(article.Content), articleURL)
	if err != nil {
		return err
	}

	content := result.Content
	if opts != nil {
		if opts.MinContentLength > 0 && len(content) < opts.MinContentLength {
			return ErrContentProcessingFailed
		}

		if opts.MaxContentLength > 0 && len(content) > opts.MaxContentLength {
			content = content[:opts.MaxContentLength]
		}

		if !opts.IncludeImages {
			content, err = removeHTMLTags(content, "img")
			if err != nil {
				log.Printf("error removing images: %v; using original content", err)
			}
		}

		if !opts.IncludeTables {
			content, err = removeHTMLTags(content, "table")
			if err != nil {
				log.Printf("error removing tables: %v; using original content", err)
			}
		}

		if !opts.IncludeVideos {
			content, err = removeHTMLTags(content, "video", "iframe")
			if err != nil {
				log.Printf("error removing videos: %v; using original content", err)
			}
		}
	}

	article.Content = content

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

// removeHTMLTags removes all occurrences of the specified HTML tags from the content.
func removeHTMLTags(content string, tagNames ...string) (string, error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return content, err
	}

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			removed := false

			for _, tagName := range tagNames {
				if c.Type == html.ElementNode && strings.EqualFold(c.Data, tagName) {
					n.RemoveChild(c)
					break
				}
			}

			if removed {
				continue
			}

			traverse(c)
		}
	}

	traverse(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return content, err
	}

	result := buf.String()
	return result, nil
}
