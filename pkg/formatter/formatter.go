package formatter

import (
	"io"

	"github.com/shrik450/dijester/pkg/models"
)

// Format represents an output format type.
type Format string

const (
	// FormatMarkdown represents the Markdown output format.
	FormatMarkdown Format = "markdown"

	// FormatEPUB represents the EPUB output format.
	FormatEPUB Format = "epub"
)

// Options contains configuration options for formatters.
type Options struct {
	// IncludeSummary toggles whether article summaries are included
	IncludeSummary bool

	// IncludeMetadata toggles whether metadata is included
	IncludeMetadata bool

	// CustomTemplate is an optional custom template to use for formatting
	CustomTemplate string

	// AdditionalOptions contains format-specific options
	AdditionalOptions map[string]interface{}
}

// Formatter defines the interface for digest output formatters.
type Formatter interface {
	// Format writes the formatted digest to the provided writer
	Format(w io.Writer, digest *models.Digest, opts *Options) error

	// SupportedFormat returns the format this formatter supports
	SupportedFormat() Format
}

// Registry maintains a collection of available formatter implementations.
type Registry struct {
	formatters map[Format]Formatter
}

// NewRegistry creates a new formatter registry.
func NewRegistry() *Registry {
	return &Registry{
		formatters: make(map[Format]Formatter),
	}
}

// Register adds a formatter to the registry.
func (r *Registry) Register(formatter Formatter) {
	r.formatters[formatter.SupportedFormat()] = formatter
}

// Get retrieves a formatter by format.
func (r *Registry) Get(format Format) (Formatter, bool) {
	formatter, ok := r.formatters[format]
	return formatter, ok
}
