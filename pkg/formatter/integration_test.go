package formatter

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/shrik450/dijester/pkg/models"
	"github.com/shrik450/dijester/pkg/processor"
	"github.com/shrik450/dijester/pkg/processor/tests"
)

// TestProcessAndFormat tests the end-to-end flow from HTML processing to Markdown formatting.
func TestProcessAndFormat(t *testing.T) {
	// Create a processor
	p := processor.NewReadabilityProcessor()

	// Create an article with raw HTML
	article := &models.Article{
		URL:     "https://example.com/test-article",
		Content: tests.SampleArticleHTML, // Using the sample HTML from processor tests
		Title:   "",                      // Let processor extract title
	}

	// Process the article
	err := p.Process(article, processor.DefaultOptions())
	if err != nil {
		t.Fatalf("Failed to process article: %v", err)
	}

	// Verify the processor extracted the expected data
	if article.Title != "Test Article Page" {
		t.Errorf("Expected title 'Test Article Page', got '%s'", article.Title)
	}

	// Now create a digest with the processed article
	digest := &models.Digest{
		Title:       "Test Digest",
		GeneratedAt: time.Now(),
		Articles:    []*models.Article{article},
	}

	// Format the digest using the Markdown formatter
	f := NewMarkdownFormatter()
	var buf bytes.Buffer
	err = f.Format(&buf, digest, DefaultOptions())
	if err != nil {
		t.Fatalf("Failed to format digest: %v", err)
	}

	// Check the formatted output
	output := buf.String()

	// Verify key elements are present in the formatted output
	expectedElements := []string{
		"# Test Digest",
		"## Test Article Page",
		"Content extraction is a critical part",
		"Why Content Extraction Matters",
	}

	for _, element := range expectedElements {
		if !strings.Contains(output, element) {
			t.Errorf("Expected formatted output to contain '%s', but it doesn't", element)
		}
	}

	// Verify that navigation elements were removed during processing
	// and don't appear in the final output
	navigationElements := []string{
		"Home</a>",
		"About Us",
		"cookie-notice",
	}

	for _, navElement := range navigationElements {
		if strings.Contains(output, navElement) {
			t.Errorf("Formatted output should not contain navigation element '%s', but it does", navElement)
		}
	}
}