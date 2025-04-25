package formatter

import (
	"fmt"
	"io"

	"github.com/shrik450/dijester/pkg/models"
)

// Options contains configuration options for formatters.
type Options struct {
	// IncludeSummary toggles whether article summaries are included
	IncludeSummary bool

	// IncludeMetadata toggles whether metadata is included
	IncludeMetadata bool

	// CustomTemplate is an optional custom template to use for formatting
	CustomTemplate string

	// StoreImages toggles whether to store images locally. Doesn't apply to
	// Markdown.
	StoreImages bool

	// AdditionalOptions contains format-specific options
	AdditionalOptions map[string]any
}

func DefaultOptions() Options {
	return Options{
		IncludeSummary:    true,
		IncludeMetadata:   false,
		StoreImages:       true,
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
func OptionsFromConfig(config map[string]any) Options {
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
	if customTemplate, ok := config["custom_template"].(string); ok {
		opts.CustomTemplate = customTemplate
	}

	return opts
}
