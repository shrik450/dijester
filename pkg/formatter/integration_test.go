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

// TestProcessAndFormatMarkdown tests the end-to-end flow from HTML processing to Markdown formatting.
func TestProcessAndFormatMarkdown(t *testing.T) {
	digest := prepareTestDigest(t)

	// Format the digest using the Markdown formatter
	f := NewMarkdownFormatter()
	var buf bytes.Buffer
	err := f.Format(&buf, digest, nil)
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
			t.Errorf(
				"Formatted output should not contain navigation element '%s', but it does",
				navElement,
			)
		}
	}
}

// TestProcessAndFormatEPUB tests the end-to-end flow from HTML processing to EPUB formatting.
func TestProcessAndFormatEPUB(t *testing.T) {
	digest := prepareTestDigest(t)

	// Format the digest using the EPUB formatter
	f := NewEPUBFormatter()
	var buf bytes.Buffer
	err := f.Format(&buf, digest, nil)
	if err != nil {
		t.Fatalf("Failed to format digest: %v", err)
	}

	// Verify EPUB was generated (should be a binary file with EPUB signature)
	if buf.Len() == 0 {
		t.Error("Expected non-empty buffer")
	}

	// Check for EPUB signature (PK zip header)
	if buf.Bytes()[0] != 0x50 || buf.Bytes()[1] != 0x4B {
		t.Error("Output does not appear to be a valid EPUB (zip) file")
	}
}

// prepareTestDigest creates a test digest with processed content for use in formatter tests.
func prepareTestDigest(t *testing.T) *models.Digest {
	// Create a processor
	p := processor.NewReadabilityProcessor()

	// Create an article with raw HTML
	article := &models.Article{
		URL:     "https://example.com/test-article",
		Content: tests.SampleArticleHTML, // Using the sample HTML from processor tests
		Title:   "",                      // Let processor extract title
	}

	opts := processor.DefaultOptions()
	err := p.Process(article, &opts)
	if err != nil {
		t.Fatalf("Failed to process article: %v", err)
	}

	// Verify the processor extracted the expected data
	if article.Title != "Test Article Page" {
		t.Errorf("Expected title 'Test Article Page', got '%s'", article.Title)
	}

	// Now create a digest with the processed article
	return &models.Digest{
		Title:       "Test Digest",
		GeneratedAt: time.Now(),
		Articles:    []*models.Article{article},
	}
}
