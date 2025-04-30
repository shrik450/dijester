package processor

import (
	"errors"
	"fmt"

	"github.com/shrik450/dijester/pkg/models"
)

// ErrContentProcessingFailed indicates that content processing failed.
var ErrContentProcessingFailed = errors.New("content processing failed")

type ProcessorConfig struct {
	// Processors is a list of processors to apply in the order defined
	Processors []string `toml:"processors"`

	// ProcessorConfigs contains optional configuration for each processor. If
	// a processor is not configured, it will use its default settings.
	ProcessorConfigs map[string]any `toml:"processor_configs"`
}

// Options contains configuration for content processors.
type Options struct {
	// Minimum content length to consider valid (in characters)
	MinContentLength int

	// Whether to include images in the processed content
	IncludeImages bool

	// Whether to include tables in the processed content
	IncludeTables bool

	// Whether to include videos in the processed content
	IncludeVideos bool

	// Maximum length for article content (0 means no limit)
	MaxContentLength int

	// Additional processor-specific options
	AdditionalOptions map[string]any
}

var defaultOptions = DefaultOptions()

// DefaultOptions returns the default processor options.
func DefaultOptions() Options {
	return Options{
		MinContentLength:  100,
		IncludeImages:     true,
		IncludeTables:     true,
		IncludeVideos:     true,
		MaxContentLength:  0,
		AdditionalOptions: map[string]any{},
	}
}

// Processor defines the interface for content processors.
type Processor interface {
	// Process processes the raw content and updates the article with processed content
	Process(article *models.Article, opts *Options) error

	// Name returns the name of this processor
	Name() string
}

var availableProcessors = []string{
	"readability",
	"sanitizer",
}

// List returns a list of available processor names.
func List() []string {
	processorNames := make([]string, len(availableProcessors))
	copy(processorNames, availableProcessors[:])
	return processorNames
}

// New returns a new instance of the specified processor.
func New(name string) (Processor, error) {
	switch name {
	case "readability":
		return NewReadabilityProcessor(), nil
	case "sanitizer":
		return NewSanitizerProcessor(), nil
	}

	return nil, fmt.Errorf("processor not found: %s", name)
}

func OptionsFromConfig(config map[string]any) (Options, error) {
	opts := DefaultOptions()

	if minLength, ok := config["min_content_length"].(int); ok {
		opts.MinContentLength = minLength
	}

	if maxLength, ok := config["max_content_length"].(int); ok {
		opts.MaxContentLength = maxLength
	}

	if includeImages, ok := config["include_images"].(bool); ok {
		opts.IncludeImages = includeImages
	}

	if includeTables, ok := config["include_tables"].(bool); ok {
		opts.IncludeTables = includeTables
	}

	if includeVideos, ok := config["include_videos"].(bool); ok {
		opts.IncludeVideos = includeVideos
	}

	if additionalOptions, ok := config["additional_options"].(map[string]any); ok {
		opts.AdditionalOptions = additionalOptions
	}

	return opts, nil
}

// InitializeProcessors initializes a list of processors based on the provided
// names and configurations. It returns a slice of Processor instances and
// their corresponding options. The options at index i correspond to the
// processor at index i in the processors slice.
func InitializeProcessors(cfg ProcessorConfig) ([]Processor, []Options, error) {
	processors := make([]Processor, len(cfg.Processors))
	options := make([]Options, len(cfg.Processors))

	for i, name := range cfg.Processors {
		processor, err := New(name)
		if err != nil {
			return nil, nil, fmt.Errorf("creating processor %s: %w", name, err)
		}
		processors[i] = processor

		config, ok := cfg.ProcessorConfigs[name].(map[string]any)
		if !ok {
			config = make(map[string]any)
		}

		opts, err := OptionsFromConfig(config)
		if err != nil {
			return nil, nil, fmt.Errorf("processing options for %s: %w", name, err)
		}

		processors[i] = processor
		options[i] = opts
	}

	return processors, options, nil
}
