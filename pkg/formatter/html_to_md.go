package formatter

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/JohannesKaufmann/html-to-markdown/plugin"
)

// HTMLToMarkdown converts HTML content to Markdown format using
// the html-to-markdown library.
func HTMLToMarkdown(htmlContent string) string {
	converter := md.NewConverter("", true, nil)

	converter.Use(plugin.GitHubFlavored())

	markdown, err := converter.ConvertString(htmlContent)
	if err != nil {
		return htmlContent
	}

	return markdown
}
