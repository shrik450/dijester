package fetcher

import (
	"context"
	"time"

	"github.com/shrik450/dijester/pkg/constants"
)

type FetcherConfig struct {
	// UserAgent is the user agent string to use for HTTP requests
	UserAgent string `toml:"user_agent"`
	// Timeout for HTTP requests
	Timeout time.Duration `toml:"timeout"`
	// RateLimit is the time in seconds to wait between requests
	RateLimit float64 `toml:"rate_limit"`
}

var (
	defaultUserAgent = "Dijester/" + constants.VERSION
	defaultTimeout   = 10 * time.Second
)

type Fetcher interface {
	// FetchURLAsString fetches the content of a URL as a string.
	FetchURLAsString(ctx context.Context, url string) (string, error)
	// FetchURL fetches the content of a URL as a byte slice.
	FetchURL(ctx context.Context, url string) ([]byte, error)
}

func FromConfig(cfg FetcherConfig) Fetcher {
	opts := make([]HTTPFetcherOption, 0)

	if cfg.Timeout > 0 {
		opts = append(opts, WithTimeout(cfg.Timeout))
	} else {
		opts = append(opts, WithTimeout(defaultTimeout))
	}

	if cfg.UserAgent != "" {
		opts = append(opts, WithUserAgent(cfg.UserAgent))
	} else {
		opts = append(opts, WithUserAgent(defaultUserAgent))
	}

	fcr := NewHTTPFetcher(opts...)

	if cfg.RateLimit > 0 {
		dur := time.Duration(cfg.RateLimit * float64(time.Second))
		return NewLimitedFetcher(fcr, dur)
	}

	return fcr
}
