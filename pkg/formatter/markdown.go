package formatter

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/shrik450/dijester/pkg/models"
)

// MarkdownFormatter implements the Formatter interface for Markdown output.
type MarkdownFormatter struct{}

// NewMarkdownFormatter creates a new Markdown formatter.
func NewMarkdownFormatter() *MarkdownFormatter {
	return &MarkdownFormatter{}
}

// SupportedFormat returns the format this formatter supports.
func (f *MarkdownFormatter) SupportedFormat() Format {
	return FormatMarkdown
}

// Format writes the formatted digest to the provided writer.
func (f *MarkdownFormatter) Format(w io.Writer, digest *models.Digest, opts *Options) error {
	if digest == nil {
		return fmt.Errorf("digest cannot be nil")
	}

	if opts == nil {
		opts = &Options{
			IncludeSummary:  true,
			IncludeMetadata: false,
		}
	}

	fmt.Fprintf(w, "# %s\n\n", digest.Title)
	fmt.Fprintf(w, "Generated on: %s\n\n", digest.GeneratedAt.Format(time.RFC1123))
	fmt.Fprintf(w, "## Contents\n\n")

	for i, article := range digest.Articles {
		fmt.Fprintf(w, "%d. [%s](#article-%d)\n", i+1, article.Title, i+1)
	}
	fmt.Fprintln(w, "")

	for i, article := range digest.Articles {
		fmt.Fprintf(w, "<a id=\"article-%d\"></a>\n", i+1)
		fmt.Fprintf(w, "## %s\n\n", article.Title)

		if article.Author != "" {
			fmt.Fprintf(w, "**Author:** %s  \n", article.Author)
		}
		if !article.PublishedAt.IsZero() {
			fmt.Fprintf(w, "**Published:** %s  \n", article.PublishedAt.Format(time.RFC1123))
		}
		fmt.Fprintf(w, "**Source:** [%s](%s)  \n\n", article.SourceName, article.URL)

		if len(article.Tags) > 0 {
			fmt.Fprintf(w, "**Tags:** %s  \n\n", strings.Join(article.Tags, ", "))
		}

		if opts.IncludeSummary && article.Summary != "" {
			fmt.Fprintf(w, "### Summary\n\n%s\n\n", article.Summary)
		}

		fmt.Fprintf(w, "### Content\n\n%s\n\n", HTMLToMarkdown(article.Content))

		if opts.IncludeMetadata && len(article.Metadata) > 0 {
			fmt.Fprintf(w, "### Metadata\n\n")
			for key, value := range article.Metadata {
				fmt.Fprintf(w, "- **%s:** %v\n", key, value)
			}
			fmt.Fprintln(w, "")
		}

		if i < len(digest.Articles)-1 {
			fmt.Fprintf(w, "---\n\n")
		}
	}

	return nil
}
