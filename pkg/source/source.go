package source

import (
	"context"

	"github.com/shrik450/dijester/pkg/models"
)

// Source defines the interface that all content sources must implement.
type Source interface {
	// Name returns a unique identifier for this source
	Name() string

	// Description returns a human-readable description of this source
	Description() string

	// Fetch retrieves articles from the source
	Fetch(ctx context.Context) ([]*models.Article, error)

	// Configure sets up the source with configuration parameters
	Configure(config map[string]interface{}) error
}

// Registry maintains a collection of available source implementations.
type Registry struct {
	sources map[string]Source
}

// NewRegistry creates a new source registry.
func NewRegistry() *Registry {
	return &Registry{
		sources: make(map[string]Source),
	}
}

// Register adds a source to the registry.
func (r *Registry) Register(source Source) {
	r.sources[source.Name()] = source
}

// Get retrieves a source by name.
func (r *Registry) Get(name string) (Source, bool) {
	source, ok := r.sources[name]
	return source, ok
}

// List returns all registered sources.
func (r *Registry) List() []Source {
	var result []Source
	for _, source := range r.sources {
		result = append(result, source)
	}
	return result
}
