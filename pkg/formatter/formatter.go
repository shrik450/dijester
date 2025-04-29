package formatter

import (
	"fmt"
	"io"

	"github.com/shrik450/dijester/pkg/fetcher"
	"github.com/shrik450/dijester/pkg/models"
)

// Options contains configuration options for formatters.
type Options struct {
	// IncludeSummary toggles whether article summaries are included
	IncludeSummary bool

	// IncludeMetadata toggles whether metadata is included
	IncludeMetadata bool

	// StoreImages toggles whether to store images locally. Doesn't apply to
	// Markdown.
	StoreImages bool

	// Fetcher is a fetcher used to fetch images and other resources if
	// required.
	Fetcher fetcher.Fetcher

	// AdditionalOptions contains format-specific options
	AdditionalOptions map[string]any
}

func DefaultOptions() Options {
	return Options{
		IncludeSummary:    true,
		IncludeMetadata:   false,
		StoreImages:       false,
		AdditionalOptions: make(map[string]any),
	}
}

// Formatter defines the interface for digest output formatters.
type Formatter interface {
	// Format writes the formatted digest to the provided writer
	Format(w io.Writer, digest *models.Digest, opts *Options) error
}

var availableFormatters = []string{
	"markdown",
	"epub",
}

// List returns a list of available formatter names.
func List() []string {
	formatterNames := make([]string, len(availableFormatters))
	copy(formatterNames, availableFormatters[:])
	return formatterNames
}

// New returns a new instance of the specified formatter.
func New(name string) (Formatter, error) {
	switch name {
	case "markdown":
		return NewMarkdownFormatter(), nil
	case "epub":
		return NewEPUBFormatter(), nil
	}

	return nil, fmt.Errorf("formatter not found: %s", name)
}

// OptionsFromConfig converts a configuration map to Options.
func OptionsFromConfig(config map[string]any, fetcher fetcher.Fetcher) Options {
	opts := DefaultOptions()

	if includeSummary, ok := config["include_summary"].(bool); ok {
		opts.IncludeSummary = includeSummary
	}
	if includeMetadata, ok := config["include_metadata"].(bool); ok {
		opts.IncludeMetadata = includeMetadata
	}
	if storeImages, ok := config["store_images"].(bool); ok {
		opts.StoreImages = storeImages
	}

	opts.Fetcher = fetcher

	return opts
}
