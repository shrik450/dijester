# Dijester Architecture

This document outlines the architecture of Dijester, a UNIX-y utility for generating EPUB/Markdown digests from news sources.

## Overview

Dijester follows a modular architecture focused on extensibility, simplicity, and the UNIX philosophy of doing one thing well. The application transforms content from various sources into a unified digest format for comfortable reading.

## Core Components

### 1. Configuration Layer

- **Config Parser**: Processes TOML configuration files defining news sources
- **Source Registry**: Maintains registry of available source implementations
- **Settings Management**: Handles global application settings

### 2. Source Layer

- **Source Interface**: Common interface implemented by all news sources
- **Built-in Sources**: Pre-defined implementations for common news sites
- **Custom Sources**: Extension points for user-defined sources

### 3. Fetcher Layer

- **Content Fetcher**: Retrieves raw content from configured sources
- **Rate Limiter**: Ensures responsible access to external services
- **Caching**: Prevents duplicate fetches and improves performance

### 4. Processor Layer

- **Content Normalizer**: Transforms diverse content into unified format
- **Metadata Extractor**: Pulls titles, authors, dates from raw content
- **Content Cleaner**: Removes ads, unnecessary formatting, etc.

### 5. Output Layer

- **Digest Generator**: Combines processed content into a single digest
- **Formatter Interface**: Common interface for output formats
- **EPUB/Markdown Formatters**: Implementations for supported outputs

## Data Flow

1. Load configuration from TOML file
2. Initialize configured sources
3. Fetch content from each source
4. Process and normalize content
5. Generate and format the final digest
6. Save digest to specified location

## Extension Points

Dijester is designed to be extended via code. Key extension points:

- **New Sources**: Add custom source implementations
- **Custom Processors**: Extend content processing pipeline
- **Output Formats**: Implement additional output formatters

This architecture enables the core goals of extensibility and simplicity while maintaining a focused approach to the single responsibility of generating content digests.