// Package source defines the interface for content sources.
package source

import (
	"time"

	"github.com/shrik450/dijester/pkg/config"
	"github.com/shrik450/dijester/pkg/fetcher"
)

// Factory creates a new instance of a source.
type Factory func(fetcher interface{}) Source

// RegisterDefaultSources registers the default source implementations.
func RegisterDefaultSources(registry *Registry, cfg *config.Config, factories map[string]Factory) {
	// Create HTTP fetcher
	httpFetcher := fetcher.NewHTTPFetcher(
		fetcher.WithUserAgent(cfg.Global.UserAgent),
		fetcher.WithTimeout(cfg.Global.Timeout),
	)

	// Create rate-limited fetcher
	limiter := fetcher.NewLimitedFetcher(httpFetcher, 1*time.Second)

	// Register sources using factories
	for name, factory := range factories {
		if name == "hackernews" {
			registry.Register(factory(limiter.Fetcher))
		} else {
			registry.Register(factory(httpFetcher))
		}
	}
}

// InitializeSources configures and returns active sources from config.
func InitializeSources(
	registry *Registry,
	cfg *config.Config,
	factories map[string]Factory,
) ([]Source, error) {
	activeSources := make([]Source, 0, len(cfg.Sources))

	// Create a standard HTTP fetcher for configuring new sources
	httpFetcher := fetcher.NewHTTPFetcher()

	for name, sourceCfg := range cfg.Sources {
		// Skip disabled sources
		if !sourceCfg.Enabled {
			continue
		}

		// Get source implementation
		impl, ok := registry.Get(sourceCfg.Type)
		if !ok {
			// Try to get by custom name
			impl, ok = registry.Get(name)
			if !ok {
				continue // Skip unknown sources
			}
		}

		// Create a new instance of the source using the factory
		factory, ok := factories[impl.Name()]
		if !ok {
			continue // Skip unknown source types
		}

		source := factory(httpFetcher)

		// Configure the source
		if err := source.Configure(sourceCfg.Options); err != nil {
			return nil, err
		}

		activeSources = append(activeSources, source)
	}

	return activeSources, nil
}
