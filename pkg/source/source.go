package source

import (
	"context"
	"fmt"

	"github.com/shrik450/dijester/pkg/fetcher"
	"github.com/shrik450/dijester/pkg/models"
	"github.com/shrik450/dijester/pkg/processor"
	"github.com/shrik450/dijester/pkg/source/hackernews"
	"github.com/shrik450/dijester/pkg/source/rss"
)

// SourceConfig contains configuration for a single source.
type SourceConfig struct {
	// Type identifies the source implementation to use
	Type string `toml:"type"`

	// Enabled determines if this source should be processed
	Enabled bool `toml:"enabled"`

	// MaxArticles limits the number of articles to include from this source
	MaxArticles int `toml:"max_articles"`

	// FetcherConfig contains configuration for the fetcher
	FetcherConfig *fetcher.FetcherConfig `toml:"fetcher_config"`

	// ProcessorConfig contains configuration for the processor
	ProcessorConfig *processor.ProcessorConfig `toml:"processor_config"`

	// Options contains source-specific configuration
	Options map[string]any `toml:"options"`
}

// Source defines the interface that all content sources must implement.
type Source interface {
	// Name returns a unique identifier for this source
	Name() string

	// Fetch retrieves articles from the source
	Fetch(ctx context.Context, fetcher fetcher.Fetcher) ([]*models.Article, error)

	// Configure sets up the source with configuration parameters
	Configure(config map[string]any) error
}

var availableSources = [...]string{
	"hackernews",
	"rss",
}

// List returns a list of available source names.
func List() []string {
	sourceNames := make([]string, len(availableSources))
	copy(sourceNames, availableSources[:])
	return sourceNames
}

// New returns a new instance of the specified source.
func New(name string) (Source, error) {
	switch name {
	case "hackernews":
		return hackernews.New(), nil
	case "rss":
		return rss.New(), nil
	}

	return nil, fmt.Errorf("source not found: %s", name)
}
