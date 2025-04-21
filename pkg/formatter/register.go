package formatter

// RegisterDefaultFormatters registers all the standard formatters.
func RegisterDefaultFormatters(registry *Registry) {
	registry.Register(NewMarkdownFormatter())
}

// DefaultOptions returns default formatter options.
func DefaultOptions() *Options {
	return &Options{
		IncludeSummary:    true,
		IncludeMetadata:   false,
		AdditionalOptions: make(map[string]interface{}),
	}
}