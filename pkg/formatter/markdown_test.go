package formatter

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/shrik450/dijester/pkg/models"
)

func TestMarkdownFormatter_SupportedFormat(t *testing.T) {
	f := NewMarkdownFormatter()
	if f.SupportedFormat() != FormatMarkdown {
		t.Errorf("Expected format to be 'markdown', got '%s'", f.SupportedFormat())
	}
}

func TestMarkdownFormatter_Format(t *testing.T) {
	f := NewMarkdownFormatter()
	var buf bytes.Buffer

	// Create test digest
	digest := &models.Digest{
		Title:       "Test Digest",
		GeneratedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Articles: []*models.Article{
			{
				Title:       "Article 1",
				Author:      "Author 1",
				PublishedAt: time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC),
				URL:         "https://example.com/1",
				Content:     "<p>This is the content of article 1</p>",
				Summary:     "Summary of article 1",
				SourceName:  "Test Source",
				Tags:        []string{"tag1", "tag2"},
				Metadata: map[string]interface{}{
					"score": 10,
				},
			},
			{
				Title:      "Article 2",
				URL:        "https://example.com/2",
				Content:    "<p>This is the content of article 2</p>",
				SourceName: "Test Source",
			},
		},
	}

	// Test with default options
	err := f.Format(&buf, digest, nil)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	result := buf.String()

	// Verify basic content
	expectedTitleLine := "# Test Digest"
	if !strings.Contains(result, expectedTitleLine) {
		t.Errorf("Output should contain '%s'", expectedTitleLine)
	}

	expectedGeneratedLine := "Generated on: Sun, 01 Jan 2023 12:00:00 UTC"
	if !strings.Contains(result, expectedGeneratedLine) {
		t.Errorf("Output should contain '%s'", expectedGeneratedLine)
	}

	// Verify table of contents
	expectedTOCLine1 := "1. [Article 1](#article-1)"
	if !strings.Contains(result, expectedTOCLine1) {
		t.Errorf("Output should contain '%s'", expectedTOCLine1)
	}

	// Verify first article content
	expectedArticleTitle := "## Article 1"
	if !strings.Contains(result, expectedArticleTitle) {
		t.Errorf("Output should contain '%s'", expectedArticleTitle)
	}

	expectedAuthorLine := "**Author:** Author 1"
	if !strings.Contains(result, expectedAuthorLine) {
		t.Errorf("Output should contain '%s'", expectedAuthorLine)
	}

	expectedTagsLine := "**Tags:** tag1, tag2"
	if !strings.Contains(result, expectedTagsLine) {
		t.Errorf("Output should contain '%s'", expectedTagsLine)
	}

	// Test with metadata enabled
	buf.Reset()
	err = f.Format(&buf, digest, &Options{IncludeSummary: true, IncludeMetadata: true})
	if err != nil {
		t.Fatalf("Format returned error with metadata enabled: %v", err)
	}

	resultWithMetadata := buf.String()
	expectedMetadataLine := "### Metadata"
	if !strings.Contains(resultWithMetadata, expectedMetadataLine) {
		t.Errorf("Output with metadata enabled should contain '%s'", expectedMetadataLine)
	}

	// Test with nil digest
	buf.Reset()
	err = f.Format(&buf, nil, nil)
	if err == nil {
		t.Error("Format should return error with nil digest")
	}
}

func TestMarkdownFormatter_Format_NoSummary(t *testing.T) {
	f := NewMarkdownFormatter()
	var buf bytes.Buffer

	digest := &models.Digest{
		Title:       "Test Digest",
		GeneratedAt: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Articles: []*models.Article{
			{
				Title:      "Article 1",
				Content:    "<p>This is the content of article 1</p>",
				Summary:    "Summary of article 1",
				SourceName: "Test Source",
			},
		},
	}

	// Test with summary disabled
	opts := &Options{IncludeSummary: false}
	err := f.Format(&buf, digest, opts)
	if err != nil {
		t.Fatalf("Format returned error: %v", err)
	}

	result := buf.String()
	notExpectedLine := "### Summary"
	if strings.Contains(result, notExpectedLine) {
		t.Errorf("Output should not contain '%s' when summaries are disabled", notExpectedLine)
	}
}
