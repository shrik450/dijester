package formatter

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/shrik450/dijester/pkg/fetcher"
	"github.com/shrik450/dijester/pkg/models"
)

// Helper function for comparing values
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestEPUBFormatter_Format(t *testing.T) {
	// Create test data
	now := time.Now()
	digest := &models.Digest{
		Title:       "Test Digest",
		GeneratedAt: now,
		Articles: []*models.Article{
			{
				Title:       "Article 1",
				Author:      "Author 1",
				PublishedAt: now.Add(-24 * time.Hour),
				URL:         "https://example.com/article1",
				Content:     "<p>This is the content of article 1.</p>",
				Summary:     "Summary of article 1",
				SourceName:  "Example Source",
				Tags:        []string{"tag1", "tag2"},
				Metadata: map[string]interface{}{
					"key1": "value1",
				},
			},
		},
	}

	formatter := NewEPUBFormatter()
	buf := &bytes.Buffer{}

	// Test with nil digest
	err := formatter.Format(buf, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "digest cannot be nil") {
		t.Errorf("Expected error for nil digest, got %v", err)
	}

	// Test with valid digest
	buf.Reset()
	opts := DefaultOptions()
	err = formatter.Format(buf, digest, &opts)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
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

func TestEPUBFormatter_StoreImages(t *testing.T) {
	// Create a test HTTP server to serve test images
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve a simple 1x1 pixel JPEG
		if r.URL.Path == "/image1.jpg" {
			w.Header().Set("Content-Type", "image/jpeg")
			// This is a minimal valid JPEG file (1x1 pixel)
			jpeg1x1Pixel := []byte{
				0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01, 0x01,
				0x00, 0x48, 0x00, 0x48, 0x00, 0x00, 0xFF, 0xDB, 0x00, 0x43, 0x00, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xC2, 0x00, 0x0B, 0x08, 0x00, 0x01, 0x00,
				0x01, 0x01, 0x01, 0x11, 0x00, 0xFF, 0xC4, 0x00, 0x14, 0x00, 0x01, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF,
				0xDA, 0x00, 0x08, 0x01, 0x01, 0x00, 0x00, 0x3F, 0x00, 0x37, 0xFF, 0xD9,
			}
			w.Write(jpeg1x1Pixel)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer testServer.Close()

	// Create test data with images and videos, using the test server URL for images
	now := time.Now()
	digest := &models.Digest{
		Title:       "Test Digest with Media",
		GeneratedAt: now,
		Articles: []*models.Article{
			{
				Title:       "Article with Media",
				Author:      "Author 1",
				PublishedAt: now.Add(-24 * time.Hour),
				URL:         testServer.URL,
				Content: `
					<p>This is a test article with images and videos.</p>
					<img src="` + testServer.URL + `/image1.jpg" alt="Test Image 1">
					<p>Some text between media.</p>
					<video controls>
						<source src="` + testServer.URL + `/video.mp4" type="video/mp4">
					</video>
					<iframe src="https://www.youtube.com/embed/dQw4w9WgXcQ" width="560" height="315" frameborder="0"></iframe>
				`,
				Summary:    "Summary of article with media",
				SourceName: "Example Source",
			},
		},
	}

	formatter := NewEPUBFormatter()
	buf := &bytes.Buffer{}

	// Test with StoreImages enabled
	opts := &Options{
		IncludeSummary:  true,
		IncludeMetadata: false,
		StoreImages:     true,
		Fetcher:         fetcher.NewHTTPFetcher(),
	}

	err := formatter.Format(buf, digest, opts)
	if err != nil {
		t.Fatalf("Failed to create EPUB: %v", err)
	}

	// Verify EPUB was generated
	if buf.Len() == 0 {
		t.Error("Expected non-empty buffer")
	}

	// Check for EPUB signature (PK zip header)
	if buf.Bytes()[0] != 0x50 || buf.Bytes()[1] != 0x4B {
		t.Error("Output does not appear to be a valid EPUB (zip) file")
	}

	// Check if the EPUB contains the image by examining the zip file contents
	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("Failed to open EPUB as zip: %v", err)
	}

	// Check for image file
	var hasImage bool
	var hasImageReference bool
	var hasImageContent bool
	var articleContent string
	var imageSize int
	var foundFiles []string
	var imageFiles []string
	var packageContent string

	// Examine EPUB content structure
	t.Log("== EPUB Content Analysis ==")
	for _, file := range zipReader.File {
		foundFiles = append(foundFiles, file.Name)
		t.Logf("- File: %s (size: %d bytes)", file.Name, file.UncompressedSize64)

		// Track image files specifically
		if strings.Contains(file.Name, "/images/") {
			imageFiles = append(imageFiles, file.Name)
			hasImage = true

			// Verify image content
			rc, err := file.Open()
			if err != nil {
				t.Fatalf("Failed to open image file: %v", err)
			}
			imgData, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				t.Fatalf("Failed to read image file: %v", err)
			}

			imageSize = len(imgData)
			t.Logf("  Image data size: %d bytes", imageSize)

			// Check for JPEG signature (FF D8 FF)
			if len(imgData) >= 3 && imgData[0] == 0xFF && imgData[1] == 0xD8 && imgData[2] == 0xFF {
				hasImageContent = true
				t.Logf("  Image verified as valid JPEG")
			} else {
				t.Errorf("  Image does not have valid JPEG signature, first bytes: %v", imgData[:min(len(imgData), 10)])
			}
		}

		// Examine package.opf to verify image is in manifest
		if strings.HasSuffix(file.Name, "package.opf") {
			rc, err := file.Open()
			if err != nil {
				t.Fatalf("Failed to open package file: %v", err)
			}
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				t.Fatalf("Failed to read package file: %v", err)
			}
			packageContent = string(content)
			t.Logf("== Package.opf Content ==")
			t.Logf("%s", packageContent)
		}

		// Get article content to verify image reference
		if strings.Contains(file.Name, "/xhtml/article-1") {
			rc, err := file.Open()
			if err != nil {
				t.Fatalf("Failed to open article file: %v", err)
			}
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				t.Fatalf("Failed to read article file: %v", err)
			}
			articleContent = string(content)

			// Log the full article content for debugging
			t.Logf("== Article-1 Content ==")
			t.Logf("%s", articleContent)

			// Check for any images reference in the article HTML
			if strings.Contains(articleContent, "../images/") {
				hasImageReference = true
			}
		}
	}

	// Verification results
	if !hasImage {
		t.Error("EPUB does not contain image file")
		t.Logf("Files in EPUB: %v", foundFiles)
	} else {
		t.Logf("Found image files: %v", imageFiles)
	}

	if !hasImageReference {
		t.Error("Article HTML does not reference the image correctly")
	}

	if !hasImageContent {
		t.Error("Image file does not contain valid image data")
	} else {
		t.Logf("Image contains valid JPEG data (%d bytes)", imageSize)
	}

	// Check if image is properly registered in the OPF manifest
	if packageContent != "" {
		if !strings.Contains(packageContent, "media-type=\"image/jpeg\"") {
			t.Error("Image/jpeg media type not found in package manifest")
		}
	}
}

func TestResolveURL(t *testing.T) {
	baseURLStr := "https://example.com/articles/news/"
	baseURL, _ := url.Parse(baseURLStr)

	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{"Absolute URL", "https://other.com/image.jpg", "https://other.com/image.jpg", false},
		{
			"Relative URL",
			"images/photo.jpg",
			"https://example.com/articles/news/images/photo.jpg",
			false,
		},
		{"Root relative URL", "/images/photo.jpg", "https://example.com/images/photo.jpg", false},
		{"Data URL", "data:image/png;base64,ABC123", "data:image/png;base64,ABC123", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := resolveURL(baseURL, tc.input)
			if (err != nil) != tc.hasError {
				t.Errorf("resolveURL() error = %v, wantErr %v", err, tc.hasError)
				return
			}
			if result != tc.expected {
				t.Errorf("resolveURL() = %v, want %v", result, tc.expected)
			}
		})
	}
}
