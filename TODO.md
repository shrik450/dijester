# Dijester Development Plan

## Step 1: Project Setup

- [x] Initialize Go module
- [x] Create basic directory structure
- [x] Write initial documentation

## Step 2: Core Infrastructure

- [x] Define core interfaces and types
  - [x] Source interface
  - [x] Content model
  - [x] Formatter interface
- [x] Implement configuration parsing
  - [x] Define TOML schema
  - [x] Create config loader
  - [x] Add validation

## Step 3: Source Implementation

- [x] Create HTTP fetcher utility
- [x] Implement RSS feed source
- [x] Add sample source implementations
  - [x] Generic RSS/Atom implementation
  - [x] Hacker News implementation

## Step 4: Content Processing

- [x] Implement content processor using go-readability
  - [x] Extract and normalize article content
  - [x] Extract metadata (title, author, summary)
  - [x] Add content sanitization (remove unwanted tags)
- [ ] Create HTML to plain text converter

## Step 5: Output Formats

- [x] Design digest format structure
- [x] Implement Markdown formatter
  - [x] Basic digest formatting
  - [x] HTML to Markdown conversion
- [ ] Create EPUB generator
- [x] Add formatting options

## Step 6: Command Line Interface

- [x] Design CLI flags and options
- [x] Implement configuration file loading
- [ ] Add output path handling
- [ ] Create progress reporting
- [ ] Implement formatter selection

## Step 7: Error Handling & Logging

- [ ] Implement robust error handling
- [x] Add basic logging infrastructure
- [ ] Create user-friendly error messages

## Step 8: Performance & Optimization

- [ ] Add concurrent fetching
- [ ] Implement caching
- [x] Add rate limiting for sources

## Step 9: Testing & Documentation

- [x] Write unit tests for core components
  - [x] Source implementations (RSS, HackerNews)
  - [x] HTTP fetcher and rate limiter
  - [x] Content processor
- [x] Create integration tests
  - [x] Content processor with realistic HTML
- [ ] Update documentation
- [ ] Add examples

## Step 10: Final Polish

- [ ] Conduct performance testing
- [x] Add sample configuration files
- [ ] Create user guide
- [ ] Package for distribution

