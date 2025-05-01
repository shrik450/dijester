package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/shrik450/dijester/pkg/fetcher"
	"github.com/shrik450/dijester/pkg/processor"
	"github.com/shrik450/dijester/pkg/source"
)

// Config represents the application configuration.
type Config struct {
	// Digest contains configuration for the generated digest
	Digest struct {
		// Title is the title of the generated digest
		Title string `toml:"title"`

		// Format specifies the output format (markdown or epub)
		Format string `toml:"format"`

		// OutputPath is where the digest will be saved
		OutputPath string `toml:"output_path"`

		// DedupByURL determines whether to deduplicate articles by URL
		DedupByURL bool `toml:"dedup_by_url"`

		// SortBy contains a list of article properties to sort by
		SortBy []string `toml:"sort_by"`
	} `toml:"digest"`

	// Sources is a map of source configurations
	Sources map[string]source.SourceConfig `toml:"sources"`

	// ProcessorConfig contains configuration for the processors that will be
	// applied to all sources unless overridden in the source config
	ProcessorConfig processor.ProcessorConfig `toml:"global_processors"`

	// FetcherConfig contains configuration for the fetcher that will be used
	// to retrieve articles from sources unless overridden in the source config
	FetcherConfig fetcher.FetcherConfig `toml:"global_fetcher"`

	// Formatting contains formatting options
	Formatting map[string]any `toml:"formatting"`
}

// LoadFile loads configuration from a TOML file.
func LoadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	return Load(string(data))
}

// Load parses configuration from a TOML string.
func Load(data string) (*Config, error) {
	config := &Config{}

	if _, err := toml.Decode(data, config); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return config, nil
}
