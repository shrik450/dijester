package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/shrik450/dijester/pkg/config"
	"github.com/shrik450/dijester/pkg/fetcher"
	"github.com/shrik450/dijester/pkg/source"
	"github.com/shrik450/dijester/pkg/source/hackernews"
	"github.com/shrik450/dijester/pkg/source/rss"
)

// Create source factories map
var sourceFactories = map[string]source.Factory{
	"rss": func(f interface{}) source.Source {
		if httpFetcher, ok := f.(*fetcher.HTTPFetcher); ok {
			return rss.New(httpFetcher)
		}
		return nil
	},
	"hackernews": func(f interface{}) source.Source {
		if httpFetcher, ok := f.(*fetcher.HTTPFetcher); ok {
			return hackernews.New(httpFetcher)
		}
		return nil
	},
}

func main() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to config file")
	testSource := flag.String("test-source", "", "Test a specific source (rss or hackernews)")
	flag.Parse()

	// Set up logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Dijester starting up...")

	if *testSource != "" {
		testSourceImplementation(*testSource)
		return
	}

	if *configPath == "" {
		log.Fatal("No config file specified. Use --config to specify a config file.")
	}

	// Load configuration
	cfg, err := config.LoadFile(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Printf("Loaded configuration from %s", *configPath)

	// Initialize source registry
	registry := source.NewRegistry()
	source.RegisterDefaultSources(registry, cfg, sourceFactories)

	// Initialize active sources
	activeSources, err := source.InitializeSources(registry, cfg, sourceFactories)
	if err != nil {
		log.Fatalf("Error initializing sources: %v", err)
	}
	log.Printf("Initialized %d active sources", len(activeSources))

	// Test by fetching from each source
	ctx := context.Background()
	for _, src := range activeSources {
		log.Printf("Fetching from source: %s", src.Name())
		articles, err := src.Fetch(ctx)
		if err != nil {
			log.Printf("Error fetching from %s: %v", src.Name(), err)
			continue
		}
		log.Printf("Fetched %d articles from %s", len(articles), src.Name())

		// Print first article title for each source
		if len(articles) > 0 {
			log.Printf("First article: %s", articles[0].Title)
		}
	}

	log.Println("Dijester completed successfully")
}

func testSourceImplementation(sourceType string) {
	log.Printf("Testing %s source implementation...", sourceType)

	// Create HTTP fetcher
	httpFetcher := fetcher.NewHTTPFetcher(
		fetcher.WithUserAgent("Dijester Test/1.0"),
		fetcher.WithTimeout(30*time.Second),
	)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var src source.Source
	var err error

	// Get factory for the source type
	factory, ok := sourceFactories[sourceType]
	if !ok {
		log.Fatalf("Unknown source type: %s", sourceType)
	}

	// Create and configure the source
	src = factory(httpFetcher)

	switch sourceType {
	case "rss":
		err = src.Configure(map[string]interface{}{
			"url":          "https://news.ycombinator.com/rss",
			"max_articles": 5,
		})
	case "hackernews":
		err = src.Configure(map[string]interface{}{
			"max_articles": 5,
			"min_score":    50,
		})
	}

	if err != nil {
		log.Fatalf("Error configuring source: %v", err)
	}

	// Fetch articles
	log.Printf("Fetching articles from %s...", src.Name())
	articles, err := src.Fetch(ctx)
	if err != nil {
		log.Fatalf("Error fetching articles: %v", err)
	}

	log.Printf("Fetched %d articles from %s", len(articles), src.Name())

	// Print articles as JSON
	for i, article := range articles {
		articleJSON, err := json.MarshalIndent(article, "", "  ")
		if err != nil {
			log.Printf("Error marshaling article to JSON: %v", err)
			continue
		}

		fmt.Printf("--- Article %d ---\n", i+1)
		fmt.Println(string(articleJSON))
		fmt.Println()
	}
}
