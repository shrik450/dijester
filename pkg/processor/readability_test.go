package processor

import (
	"strings"
	"testing"

	"github.com/shrik450/dijester/pkg/models"
)

func TestReadabilityProcessor_Process_SimpleHTML(t *testing.T) {
	p := NewReadabilityProcessor()
	article := &models.Article{
		URL:     "https://example.com",
		Content: `<html><head><title>Page Title</title></head><body><article><h1>Test Article</h1><p>This is a test paragraph</p></article><div class="noise">Noise content</div></body></html>`,
	}
	err := p.Process(article, &defaultOptions)
	if err != nil {
		t.Fatalf("Failed to process article: %v", err)
	}

	// Check that the content was extracted (we don't check exact HTML structure as readability may format it differently)
	if !contains(article.Content, "Test Article") {
		t.Error("Processed content should retain article heading text")
	}
	if !contains(article.Content, "This is a test paragraph") {
		t.Error("Processed content should retain article paragraph text")
	}

	// Readability often uses the page title as the article title, especially for simple HTML
	if article.Title != "Page Title" && article.Title != "Test Article" {
		t.Errorf(
			"Expected title to be either 'Page Title' or 'Test Article', got '%s'",
			article.Title,
		)
	}
}

func TestReadabilityProcessor_Process_WithOptions(t *testing.T) {
	p := NewReadabilityProcessor()
	htmlContent := `<html><body>
		<h1>Test Article</h1>
		<p>This is the main content.</p>
		<img src="test.jpg" alt="Test Image">
		<table><tr><td>Table content</td></tr></table>
	</body></html>`

	// Test with images disabled
	t.Run("NoImages", func(t *testing.T) {
		article := &models.Article{
			URL:     "https://example.com",
			Content: htmlContent,
		}
		opts := DefaultOptions()
		opts.IncludeImages = false

		err := p.Process(article, &opts)
		if err != nil {
			t.Fatalf("Failed to process article: %v", err)
		}

		if contains(article.Content, "<img") {
			t.Error("Processed content should not contain images when IncludeImages=false")
		}
	})

	// Test with tables disabled
	t.Run("NoTables", func(t *testing.T) {
		article := &models.Article{
			URL:     "https://example.com",
			Content: htmlContent,
		}
		opts := DefaultOptions()
		opts.IncludeTables = false

		err := p.Process(article, &opts)
		if err != nil {
			t.Fatalf("Failed to process article: %v", err)
		}

		if contains(article.Content, "<table") {
			t.Error("Processed content should not contain tables when IncludeTables=false")
		}
	})

	// Test with minimum content length
	t.Run("MinContentLength", func(t *testing.T) {
		article := &models.Article{
			URL:     "https://example.com",
			Content: "<html><body><p>Short</p></body></html>",
		}
		opts := DefaultOptions()
		opts.MinContentLength = 100

		err := p.Process(article, &opts)
		if err != ErrContentProcessingFailed {
			t.Errorf(
				"Expected ErrContentProcessingFailed for content shorter than minimum, got %v",
				err,
			)
		}
	})

	// Test with maximum content length
	t.Run("MaxContentLength", func(t *testing.T) {
		article := &models.Article{
			URL:     "https://example.com",
			Content: "<html><body><p>This is a longer paragraph that should get truncated.</p></body></html>",
		}
		opts := DefaultOptions()
		opts.MaxContentLength = 20

		err := p.Process(article, &opts)
		if err != nil {
			t.Fatalf("Failed to process article: %v", err)
		}

		if len(article.Content) > 20 {
			t.Errorf(
				"Processed content length should be limited to MaxContentLength, got %d",
				len(article.Content),
			)
		}
	})
}

func TestReadabilityProcessor_Process_ExtractsMetadata(t *testing.T) {
	p := NewReadabilityProcessor()
	htmlContent := `<html>
		<head>
			<title>Page Title</title>
			<meta name="author" content="Test Author">
			<meta name="description" content="Test Description">
		</head>
		<body>
			<h1>Article Headline</h1>
			<div class="byline">By Test Author</div>
			<p>This is the article content.</p>
		</body>
	</html>`

	article := &models.Article{
		URL:     "https://example.com",
		Content: htmlContent,
	}

	err := p.Process(article, &defaultOptions)
	if err != nil {
		t.Fatalf("Failed to process article: %v", err)
	}

	// Check if title was extracted - readability typically uses the page title
	expectedTitle := "Page Title"
	if article.Title == "" {
		t.Error("Failed to extract any title from the article")
	} else if article.Title != expectedTitle {
		t.Errorf("Expected title '%s', got: '%s'", expectedTitle, article.Title)
	}

	// Check that the author was extracted
	expectedAuthor := "Test Author"
	if article.Author == "" {
		t.Error("Failed to extract author information from the article")
	} else if article.Author != expectedAuthor {
		t.Errorf("Expected author '%s', got: '%s'", expectedAuthor, article.Author)
	}

	// Check that the summary was extracted
	expectedSummary := "Test Description"
	if article.Summary == "" {
		t.Error("Failed to extract summary information from the article")
	} else if article.Summary != expectedSummary {
		t.Errorf("Expected summary '%s', got: '%s'", expectedSummary, article.Summary)
	}
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
