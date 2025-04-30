# Dijester Changelog

## v0.2.0 (2025-04-30)

- Fix inconsistent image loading in the EPUB formatter.
- Add sanitizer processor for sanitizing HTML content. If you are outputting
  EPUB, you should always have this processor at the end of the pipeline.
- Article content HTML is limited to 1MB by default to prevent excessive memory
  usage. This is only limiting the HTML content, and so this shouldn't affect
  most articles.


## v0.1.0 (2025-04-25)

The initial release of Dijester!

This release includes:

- Sources: Hacker News, RSS/Atom
- Formatters: Markdown, EPUB
- Processors: Readability
- Configuration via TOML
