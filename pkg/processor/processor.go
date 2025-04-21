package processor

import (
	"errors"

	"github.com/shrik450/dijester/pkg/models"
)

// ErrContentProcessingFailed indicates that content processing failed.
var ErrContentProcessingFailed = errors.New("content processing failed")

// Options contains configuration for content processors.
type Options struct {
	// Minimum content length to consider valid (in characters)
	MinContentLength int

	// Whether to include images in the processed content
	IncludeImages bool

	// Whether to include tables in the processed content
	IncludeTables bool

	// Maximum length for article content (0 means no limit)
	MaxContentLength int

	// Additional processor-specific options
	AdditionalOptions map[string]interface{}
}

// DefaultOptions returns the default processor options.
func DefaultOptions() *Options {
	return &Options{
		MinContentLength:  100,
		IncludeImages:     true,
		IncludeTables:     true,
		MaxContentLength:  0,
		AdditionalOptions: map[string]interface{}{},
	}
}

// Processor defines the interface for content processors.
type Processor interface {
	// Process processes the raw content and updates the article with processed content
	Process(article *models.Article, opts *Options) error

	// Name returns the name of this processor
	Name() string
}

// Registry maintains a collection of available processor implementations.
type Registry struct {
	processors map[string]Processor
}

// NewRegistry creates a new processor registry.
func NewRegistry() *Registry {
	return &Registry{
		processors: make(map[string]Processor),
	}
}

// Register adds a processor to the registry.
func (r *Registry) Register(processor Processor) {
	r.processors[processor.Name()] = processor
}

// Get retrieves a processor by name.
func (r *Registry) Get(name string) (Processor, bool) {
	processor, ok := r.processors[name]
	return processor, ok
}
