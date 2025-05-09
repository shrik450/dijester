# Configuring Dijester

Dijester uses TOML configuration files to control its behavior. This document explains all configuration options.

## Configuration File Structure

Dijester's configuration file is divided into several sections:

- `digest`: Controls the output digest format and naming
- `global`: Sets global fetch parameters
- `sources`: Defines content sources to collect from
- `global_processors`: Configures content processing
- `formatting`: Sets formatting options for the output

## Digest Settings

The `digest` section controls the output format and metadata:

```toml
[digest]
format = "epub"  # Can be "epub" or "markdown"
output_path = "digest-{{.DTLong}}.epub"  # Supports Go templates for the generated time
title = "My Daily Digest - {{.Date}}"  # Same as above
dedup_by_url = true  # Remove duplicate articles with the same URL
sort_by = ["PublishedAt:desc", "Title:asc"]  # Sort articles by field(s) with direction
```

### Available Date Templates

- `{{.Year}}`: 4-digit year (e.g., 2025)
- `{{.Month}}`: Month name (e.g., April)
- `{{.MonthNum}}`: Month number (e.g., 04)
- `{{.Day}}`: Day of month (e.g., 25)
- `{{.Date}}`: Short date format (e.g., 2025-04-25)
- `{{.TStamp}}`: RFC3339 timestamp (e.g., 2025-04-25T15:04:05Z)
- `{{.DTLong}}`: Long date-time format with timestamp (e.g., Mon Jan 02
    15:04:05 -0700 2006)


## Global Settings

The `global` section controls general fetch behavior:

```toml
[global]
concurrent_fetches = 3  # Number of articles to fetch in parallel, not currently implemented
timeout = "30s"  # HTTP request timeout
user_agent = "Dijester/1.0"  # User agent for HTTP requests
```

## Source Configuration

The `sources` section defines where dijester fetches content from:

```toml
[sources.NAME]
enabled = true  # Whether this source is active
max_articles = 10  # Maximum articles to include from this source
type = "TYPE"  # Source type (e.g., "hackernews", "rss")
word_denylist = ["spam", "unwanted"]  # Filter out articles containing these words

[sources.NAME.options]
# Source-specific options
```

### Available Source Types

#### Hacker News Source

```toml
[sources.hackernews]
type = "hackernews"
enabled = true
max_articles = 10

[sources.hackernews.options]
min_score = 100  # Minimum score to include an article
page = "front_page"  # Can be "front_page", "new", "best", "ask", "show", or "past"
show_dead = false  # Whether to include dead posts
```

#### RSS Source

```toml
[sources.example_rss]
type = "rss"
enabled = true
max_articles = 10

[sources.example_rss.options]
url = "https://example.com/feed.xml"  # URL of the RSS feed
include_content = true  # Whether to include the content from the RSS feed
```

## Content Processing

The `global_processors` section configures how dijester processes all articles
unless overridden for a source:

```toml
[global_processors]
processors = ["readability", "sanitizer"]  # List of processors to apply

[global_processors.processor_configs.readability]
min_content_length = 100  # Minimum content length to process
max_content_length = 0  # Maximum content length (0 means no limit)
include_images = true  # Whether to include images
include_tables = true  # Whether to include tables
include_videos = true  # Whether to include videos
```

### Available Processors

- `readability`: Extracts the main content from a webpage, removing navigation,
   ads, etc.
- `sanitizer`: Cleans up HTML content to remove unwanted tags and attributes.
   If you are outputting EPUB, you should always have this processor at the end
   of the pipeline.

## Formatting Options

The `formatting` section controls how the output is formatted:

```toml
[formatting]
include_metadata = true  # Whether to include source metadata
```

## Article Filtering and Sorting

Dijester provides several ways to filter and sort articles:

### URL Deduplication

When multiple sources provide the same content (or a source has duplicate entries), you can enable URL-based deduplication:

```toml
[digest]
dedup_by_url = true  # Remove duplicate articles with the same URL
```

This will keep only the first occurrence of each unique URL in the final digest.

### Word Denylist Filtering

Each source can have a list of words that will cause articles to be filtered
out if they appear in the title, content, or summary:

```toml
[sources.example]
word_denylist = ["crypto", "nft", "bitcoin"]
```

The filtering is case-insensitive. Any article containing any of these words will be excluded from the digest.

### Article Sorting

You can sort articles based on their properties:

```toml
[digest]
sort_by = ["SourceName", "PublishedAt:desc", "Title"]
```

Each sort field consists of a property name and an optional direction (`asc` or `desc`). If no direction is specified, ascending order is used by default.

Available properties for sorting:
- `Title`: Article title
- `Author`: Article author
- `PublishedAt`: Publication date
- `URL`: Article URL
- `SourceName`: Name of the source

Multiple sort fields are applied in order, with later fields used as tie-breakers.

If no sorts are specified, there is no guarantee on the order of articles in the output.

## Override Configurations

You can override global fetch and processor settings for specific sources:

```toml
[sources.custom_source.fetcher_config]
timeout = "60s"  # Overrides the global timeout
rate_limit = 5.0  # Wait 5 seconds between requests

[sources.custom_source.processor_config]
processors = ["readability", "sanitizer"]  # Different processor pipeline for this source
```

## Complete Example

For complete examples, see the sample configuration files in the dijester repository:
- `example-config-epub.toml`: Example configuration for EPUB output
- `example-config-md.toml`: Example configuration for Markdown output
