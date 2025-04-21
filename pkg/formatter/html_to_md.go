package formatter

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
)

// HTMLToMarkdown converts HTML content to Markdown format using
// the html-to-markdown library.
func HTMLToMarkdown(htmlContent string) string {
	// Create a new converter
	converter := md.NewConverter("", true, nil)

	// Add GitHub flavored markdown rules
	converter.Use(plugin.GitHubFlavored())

	// Configure options - no need to set additional options

	// Convert HTML to Markdown
	markdown, err := converter.ConvertString(htmlContent)
	if err != nil {
		// If conversion fails, return a sanitized version of the original content
		return htmlContent
	}

	return markdown
}
