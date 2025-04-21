package formatter

import (
	"strings"
	"testing"
)

func TestHTMLToMarkdown(t *testing.T) {
	// Test basic HTML to Markdown conversion
	htmlTests := []struct {
		name     string
		html     string
		contains []string // Looking for these substrings in the output
	}{
		{
			name: "Paragraph",
			html: "<p>Test paragraph</p>",
			contains: []string{
				"Test paragraph",
			},
		},
		{
			name: "Headers",
			html: "<h1>Header 1</h1><h2>Header 2</h2><h3>Header 3</h3>",
			contains: []string{
				"# Header 1",
				"## Header 2",
				"### Header 3",
			},
		},
		{
			name: "Emphasis",
			html: "<p><em>Italic</em> and <strong>Bold</strong></p>",
			contains: []string{
				"_Italic_",
				"**Bold**",
			},
		},
		{
			name: "Links",
			html: "<p><a href=\"https://example.com\">Example</a></p>",
			contains: []string{
				"[Example](https://example.com)",
			},
		},
		{
			name: "Lists",
			html: "<ul><li>Item 1</li><li>Item 2</li></ul>",
			contains: []string{
				"- Item 1",
				"- Item 2",
			},
		},
		{
			name: "Nested Content",
			html: "<div><h2>Nested Header</h2><p>Nested paragraph</p></div>",
			contains: []string{
				"## Nested Header",
				"Nested paragraph",
			},
		},
	}

	for _, tc := range htmlTests {
		t.Run(tc.name, func(t *testing.T) {
			result := HTMLToMarkdown(tc.html)
			for _, s := range tc.contains {
				if !strings.Contains(result, s) {
					t.Errorf("Expected result to contain '%s', but got:\n%s", s, result)
				}
			}
		})
	}
}

func TestHTMLToMarkdownComplex(t *testing.T) {
	// Test a more complex HTML document
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Test Document</title>
</head>
<body>
    <h1>Main Heading</h1>
    <p>This is a paragraph with <strong>bold text</strong> and <em>italic text</em>.</p>
    <p>Here's a <a href="https://example.com">link</a> to a website.</p>
    
    <h2>Lists</h2>
    <ul>
        <li>Unordered item 1</li>
        <li>Unordered item 2</li>
        <li>Unordered item 3 with <a href="https://example.org">link</a></li>
    </ul>
    
    <ol>
        <li>Ordered item 1</li>
        <li>Ordered item 2</li>
        <li>Ordered item 3</li>
    </ol>
    
    <h2>Code</h2>
    <pre><code>function example() {
    return "code block";
}</code></pre>
    
    <h2>Blockquote</h2>
    <blockquote>
        <p>This is a blockquote.</p>
        <p>It can contain multiple paragraphs.</p>
    </blockquote>
</body>
</html>
`

	markdown := HTMLToMarkdown(html)

	// Check for key elements in the output
	expectedContents := []string{
		"# Main Heading",
		"bold text",
		"italic text",
		"[link](https://example.com)",
		"- Unordered item",
		"1. Ordered item",
		"```",
		"function example()",
		"```",
		"This is a blockquote",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(markdown, expected) {
			t.Errorf("Expected markdown to contain '%s', but it doesn't.", expected)
		}
	}
}

func TestHTMLToMarkdownEdgeCases(t *testing.T) {
	edgeCases := []struct {
		name  string
		html  string
		check func(string) bool
	}{
		{
			name: "Empty HTML",
			html: "",
			check: func(markdown string) bool {
				return markdown == ""
			},
		},
		{
			name: "Invalid HTML",
			html: "<p>Unclosed paragraph",
			check: func(markdown string) bool {
				return strings.Contains(markdown, "Unclosed paragraph")
			},
		},
		{
			name: "HTML with script tags",
			html: "<p>Text</p><script>alert('hello');</script>",
			check: func(markdown string) bool {
				return strings.Contains(markdown, "Text") && !strings.Contains(markdown, "alert")
			},
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			result := HTMLToMarkdown(tc.html)
			if !tc.check(result) {
				t.Errorf("Failed edge case check for '%s', got: '%s'", tc.name, result)
			}
		})
	}
}