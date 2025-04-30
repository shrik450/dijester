package formatter

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-shiori/go-epub"
	"golang.org/x/net/html"

	"github.com/shrik450/dijester/pkg/constants"
	"github.com/shrik450/dijester/pkg/fetcher"
	"github.com/shrik450/dijester/pkg/models"
)

// EPUBFormatter implements the Formatter interface for EPUB output.
type EPUBFormatter struct{}

// NewEPUBFormatter creates a new EPUB formatter.
func NewEPUBFormatter() *EPUBFormatter {
	return &EPUBFormatter{}
}

// Format writes the formatted digest to the provided writer.
func (f *EPUBFormatter) Format(w io.Writer, digest *models.Digest, opts *Options) error {
	if digest == nil {
		return fmt.Errorf("digest cannot be nil")
	}

	if opts == nil {
		opt := DefaultOptions()
		opts = &opt
	}

	fetcher := opts.Fetcher
	if fetcher == nil && opts.StoreImages {
		return fmt.Errorf("store_images is true but no fetcher was provided")
	}

	e, err := epub.NewEpub(digest.Title)
	if err != nil {
		return fmt.Errorf("error creating epub: %w", err)
	}
	if client, ok := opts.AdditionalOptions["httpClient"].(*http.Client); ok {
		e.Client = client
	}

	e.SetAuthor("Dijester v" + constants.VERSION)
	e.SetDescription(
		fmt.Sprintf("Digest generated on %s", digest.GeneratedAt.Format(time.RFC1123)),
	)

	tocHTML := f.generateTOC(digest)
	_, err = e.AddSection(tocHTML, "Table of Contents", "", "")
	if err != nil {
		return fmt.Errorf("error adding table of contents: %w", err)
	}

	tmpDir, err := os.MkdirTemp("", "dijester-epub-")
	if err != nil {
		return fmt.Errorf("error creating temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	for i, article := range digest.Articles {
		sectionHTML := f.generateArticleHTML(e, article, opts, tmpDir, fetcher)

		_, err = e.AddSection(sectionHTML, article.Title, fmt.Sprintf("article-%d", i+1), "")
		if err != nil {
			return fmt.Errorf("error adding article %d: %w", i+1, err)
		}
	}

	_, err = e.WriteTo(w)
	if err != nil {
		return fmt.Errorf("error writing EPUB to writer: %w", err)
	}

	return nil
}

var tocTemplate = `
<html>
<head>
	<title>{{.Title}}</title>
</head>
<body>
	<h1>{{.Title}}</h1>
	<p>Generated on: {{.GeneratedAt}}</p>
	<h2>Contents</h2>
	<ol>
		{{range $i, $a := .Articles}}
			<li><a href="article-{{$i}}.xhtml">{{$a.Title}}</a></li>
		{{end}}
	</ol>
</body>
</html>
`

// generateTOC creates the HTML table of contents.
func (f *EPUBFormatter) generateTOC(digest *models.Digest) string {
	tmpl := template.Must(template.New("toc").Parse(tocTemplate))
	var sb strings.Builder

	tmpl.Execute(&sb, map[string]any{
		"Title":       digest.Title,
		"GeneratedAt": digest.GeneratedAt.Format(time.RFC1123),
		"Articles":    digest.Articles,
	})

	return sb.String()
}

var articleTemplate = `
<h1>{{.Title}}</h1>

<div style="background-color: #f5f5f5; border: 1px solid #e0e0e0; padding: 10px; margin-bottom: 20px; font-size: 0.9em; color: #555; border-radius: 4px;">
	<div style="margin-bottom: 5px;">
		{{if .Author}}<span style="margin-right: 15px;">By {{.Author}}</span>{{end}}
		{{if .PublishedAt}}<span style="margin-right: 15px;">at {{.PublishedAt}}</span>{{end}}
	</div>
	<div>
		<span style="margin-right: 15px;"><a href="{{.URL}}" style="color: #1a73e8; text-decoration: none;">Link</a></span>
		{{if .Tags}}<span><strong>Tags:</strong> {{.Tags}}</span>{{end}}
	</div>
</div>

{{if .IncludeSummary}}
	{{if .Summary}}
		<div style="background-color: #f8f9fa; border-left: 4px solid #1a73e8; padding: 10px 15px; margin-bottom: 20px;">
			<h2 style="margin-top: 0; font-size: 1.2em;">Summary</h2>
			<p>{{.Summary}}</p>
		</div>
	{{end}}
{{end}}

<div class="article-content">
	{{.Content}}
</div>

{{if .IncludeMetadata}}
	{{if .Metadata}}
		<div style="margin-top: 30px; border-top: 1px solid #e0e0e0; padding-top: 15px;">
			<h2>Metadata</h2>
			<ul>
				{{range $key, $value := .Metadata}}
					<li><strong>{{$key}}:</strong> {{$value}}</li>
				{{end}}
			</ul>
		</div>
	{{end}}
{{end}}
`

// generateArticleHTML creates the HTML for a single article.
func (f *EPUBFormatter) generateArticleHTML(
	e *epub.Epub,
	article *models.Article,
	opts *Options,
	tmpDir string,
	fetcher fetcher.Fetcher,
) string {
	tmpl := template.Must(template.New("article").Parse(articleTemplate))

	if opts.StoreImages {
		embedImages(e, article, tmpDir, fetcher)
	}

	var sb strings.Builder
	tmpl.Execute(&sb, map[string]any{
		"Title":           article.Title,
		"Content":         template.HTML(article.Content),
		"Author":          article.Author,
		"PublishedAt":     article.PublishedAt.Format(time.RFC1123),
		"URL":             article.URL,
		"SourceName":      article.SourceName,
		"Tags":            strings.Join(article.Tags, ", "),
		"IncludeSummary":  opts.IncludeSummary,
		"Summary":         article.Summary,
		"IncludeMetadata": opts.IncludeMetadata,
		"Metadata":        article.Metadata,
	})

	return sb.String()
}

func embedImages(e *epub.Epub, article *models.Article, tmpDir string, fetcher fetcher.Fetcher) {
	node, err := html.Parse(strings.NewReader(article.Content))
	if err != nil {
		log.Printf("error parsing HTML: %s", err)
		return
	}
	articleUrl, err := url.Parse(article.URL)
	if err != nil {
		log.Printf("error parsing article URL: %s", err)
		return
	}

	ctx := context.Background()

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for i := range n.Attr {
				attr := &n.Attr[i]
				if attr.Key == "src" {
					beforeURL := attr.Val
					if strings.HasPrefix(beforeURL, "data:image/") {
						continue
					}

					resolvedURL, err := resolveURL(articleUrl, attr.Val)
					if err != nil {
						log.Printf("error resolving URL: %s", err)
						continue
					}

					tmpFile, err := os.CreateTemp(tmpDir, "img-")
					if err != nil {
						log.Printf("error creating temp file: %s", err)
						continue
					}
					defer tmpFile.Close()
					// tmpFile shouldn't be deleted here, as the epub only
					// "grabs" the file when it writes the full epub. It will
					// be deleted when the entire `tmpDir` (managed by the
					// calling function) is deleted.

					err = fetcher.StreamURL(ctx, resolvedURL, tmpFile)
					if err != nil {
						log.Printf("error fetching image: %s", err)
						continue
					}

					newURL, err := e.AddImage(tmpFile.Name(), "")
					if err != nil {
						log.Printf("error adding image to EPUB: %s", err)
						continue
					}

					if newURL != "" {
						attr.Val = newURL
					}

					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(node)

	var buf strings.Builder
	if err := html.Render(&buf, node); err != nil {
		log.Printf("error rendering HTML: %s", err)
		return
	}

	newContent := buf.String()
	newContent = strings.TrimPrefix(newContent, "<html><head></head><body>")
	newContent = strings.TrimSuffix(newContent, "</body></html>")

	article.Content = newContent
}

// resolveURL resolves a potentially relative URL against a base URL.
func resolveURL(baseURL *url.URL, urlStr string) (string, error) {
	if strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://") {
		return urlStr, nil
	}

	if strings.HasPrefix(urlStr, "data:") {
		return urlStr, nil
	}

	relURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	resolvedURL := baseURL.ResolveReference(relURL)
	return resolvedURL.String(), nil
}
