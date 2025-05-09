package processor

import (
	"strings"
	"testing"

	"github.com/shrik450/dijester/pkg/models"
	"github.com/shrik450/dijester/pkg/processor/tests"
)

func TestReadabilityProcessor_RealWorldHTML(t *testing.T) {
	p := NewReadabilityProcessor()
	article := &models.Article{
		URL:     "https://example.com/test-article",
		Content: tests.SampleArticleHTML,
	}

	err := p.Process(article, &defaultOptions)
	if err != nil {
		t.Fatalf("Failed to process article: %v", err)
	}

	// Check if the title was extracted
	expectedTitle := "Test Article Page"
	if article.Title == "" {
		t.Error("Failed to extract any title from the article")
	} else if article.Title != expectedTitle {
		t.Errorf("Expected title '%s', got: '%s'", expectedTitle, article.Title)
	}

	// Check that the author was extracted
	expectedAuthor := "John Doe"
	if article.Author == "" {
		t.Error("Failed to extract author information from the article")
	} else if article.Author != expectedAuthor {
		t.Errorf("Expected author '%s', got: '%s'", expectedAuthor, article.Author)
	}

	// Check that the main content was extracted
	elements := []string{
		"Content extraction is a critical part",
		"Why Content Extraction Matters",
		"Density-Based Approaches",
		"DOM-Based Approaches",
		"Conclusion",
	}

	for _, element := range elements {
		if !strings.Contains(article.Content, element) {
			t.Errorf("Expected processed content to contain '%s', but it doesn't", element)
		}
	}

	// Check that navigation and footer elements were removed
	navigationElements := []string{
		"<nav>",
		"Home</a>",
		"About Us",
		"Categories",
		"cookie-notice",
	}

	for _, navElement := range navigationElements {
		if strings.Contains(article.Content, navElement) {
			t.Errorf(
				"Processed content should not contain navigation element '%s', but it does",
				navElement,
			)
		}
	}

	// Test with different options
	t.Run("NoImagesNoTables", func(t *testing.T) {
		article := &models.Article{
			URL:     "https://example.com/test-article",
			Content: tests.SampleArticleHTML,
		}

		opts := DefaultOptions()
		opts.IncludeImages = false
		opts.IncludeTables = false

		err := p.Process(article, &opts)
		if err != nil {
			t.Fatalf("Failed to process article: %v", err)
		}

		if strings.Contains(article.Content, "<img") {
			t.Error("Processed content should not contain images when IncludeImages=false")
		}

		if strings.Contains(article.Content, "<table") {
			t.Error("Processed content should not contain tables when IncludeTables=false")
		}
	})
}
