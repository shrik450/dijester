package config

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"

	"github.com/shrik450/dijester/pkg/formatter"
)

// SourceConfig contains configuration for a single source.
type SourceConfig struct {
	// Type identifies the source implementation to use
	Type string `toml:"type"`

	// Name is an optional custom name for this source
	Name string `toml:"name"`

	// Enabled determines if this source should be processed
	Enabled bool `toml:"enabled"`

	// MaxArticles limits the number of articles to include from this source
	MaxArticles int `toml:"max_articles"`

	// Options contains source-specific configuration
	Options map[string]interface{} `toml:"options"`
}

// Config represents the application configuration.
type Config struct {
	// Digest contains configuration for the generated digest
	Digest struct {
		// Title is the title of the generated digest
		Title string `toml:"title"`

		// Format specifies the output format (markdown or epub)
		Format formatter.Format `toml:"format"`

		// OutputPath is where the digest will be saved
		OutputPath string `toml:"output_path"`
	} `toml:"digest"`

	// Global contains application-wide settings
	Global struct {
		// ConcurrentFetches limits how many sources to fetch concurrently
		ConcurrentFetches int `toml:"concurrent_fetches"`

		// UserAgent is the user agent string to use for HTTP requests
		UserAgent string `toml:"user_agent"`

		// Timeout for HTTP requests
		Timeout time.Duration `toml:"timeout"`
	} `toml:"global"`

	// Sources is a map of source configurations
	Sources map[string]SourceConfig `toml:"sources"`

	// Formatting contains formatting options
	Formatting formatter.Options
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

	// Set defaults
	config.Global.ConcurrentFetches = 5
	config.Global.UserAgent = "Dijester/1.0"
	config.Global.Timeout = 30 * time.Second
	config.Formatting.IncludeSummary = true

	// Parse TOML
	if _, err := toml.Decode(data, config); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// Validate
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate ensures the configuration is valid.
func (c *Config) Validate() error {
	if c.Digest.Title == "" {
		return fmt.Errorf("digest title is required")
	}

	if c.Digest.OutputPath == "" {
		return fmt.Errorf("output path is required")
	}

	if c.Digest.Format != formatter.FormatMarkdown && c.Digest.Format != formatter.FormatEPUB {
		return fmt.Errorf("unsupported format: %s", c.Digest.Format)
	}

	if len(c.Sources) == 0 {
		return fmt.Errorf("at least one source must be configured")
	}

	for name, source := range c.Sources {
		if source.Type == "" {
			return fmt.Errorf("source '%s' is missing type", name)
		}
	}

	return nil
}
