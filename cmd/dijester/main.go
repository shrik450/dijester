package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/shrik450/dijester/pkg/config"
	"github.com/shrik450/dijester/pkg/fetcher"
	"github.com/shrik450/dijester/pkg/formatter"
	"github.com/shrik450/dijester/pkg/models"
	"github.com/shrik450/dijester/pkg/processor"
	"github.com/shrik450/dijester/pkg/source"
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	outputPath := flag.String("output", "", "Path to output file, overrides config")
	format := flag.String("format", "", "Output format (markdown or epub), overrides config")
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Dijester starting up...")

	if *configPath == "" {
		log.Fatal("No config file specified. Use --config to specify a config file.")
	}

	cfg, err := config.LoadFile(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Printf("Loaded configuration from %s", *configPath)

	var outputFormat string
	if *format != "" {
		outputFormat = *format
	} else {
		outputFormat = cfg.Digest.Format
	}

	outputFormatter, err := formatter.New(outputFormat)
	if err != nil {
		log.Fatalf("Error initializing formatter: %v", err)
	}
	fmtOpts := formatter.OptionsFromConfig(cfg.Formatting)

	digest := &models.Digest{
		Title:       cfg.Digest.Title,
		GeneratedAt: time.Now(),
		Articles:    make([]*models.Article, 0),
	}

	globalFetcher := fetcher.FromConfig(cfg.FetcherConfig)
	globalProcs, globalProcsOpts, err := processor.InitializeProcessors(cfg.ProcessorConfig)
	if err != nil {
		log.Fatalf("Error initializing global processors: %v", err)
	}

	ctx := context.Background()

	for srcName, srcCfg := range cfg.Sources {
		if !srcCfg.Enabled {
			log.Println("Skipping disabled source: ", srcName)
			continue
		}

		log.Println("Fetching from source: ", srcName)
		src, err := source.New(srcCfg.Type)
		if err != nil {
			log.Printf("Error initializing source %s: %v", srcName, err)
			continue
		}
		err = src.Configure(srcCfg.Options)
		if err != nil {
			log.Printf("Error configuring source %s: %v", srcName, err)
			continue
		}

		var srcFetcher fetcher.Fetcher
		if srcCfg.FetcherConfig != nil {
			srcFetcher = fetcher.FromConfig(*srcCfg.FetcherConfig)
		} else {
			srcFetcher = globalFetcher
		}

		var srcProcs []processor.Processor
		var srcProcsOpts []processor.Options
		if srcCfg.ProcessorConfig != nil {
			srcProcs, srcProcsOpts, err = processor.InitializeProcessors(*srcCfg.ProcessorConfig)
			if err != nil {
				log.Printf("Error initializing processors for source %s: %v", srcName, err)
				continue
			}
		} else {
			srcProcs = globalProcs
			srcProcsOpts = globalProcsOpts
		}

		articles, err := src.Fetch(ctx, srcFetcher)
		if err != nil {
			log.Printf("Error fetching from %s: %v", src.Name(), err)
			continue
		}
		log.Printf("Fetched %d articles from %s", len(articles), src.Name())

		for _, article := range articles {
			for procI, proc := range srcProcs {
				if err := proc.Process(article, &srcProcsOpts[procI]); err != nil {
					log.Printf("Error processing article with %s: %v", proc.Name(), err)
					continue
				}
			}
		}

		digest.Articles = append(digest.Articles, articles...)
	}

	var finalOutputPath string

	if *outputPath != "" {
		finalOutputPath = *outputPath
	} else {
		finalOutputPath = cfg.Digest.OutputPath
	}

	outputDir := filepath.Dir(finalOutputPath)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	file, err := os.Create(finalOutputPath)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer file.Close()

	log.Printf("Writing digest with %d articles to %s", len(digest.Articles), finalOutputPath)
	if err := outputFormatter.Format(file, digest, &fmtOpts); err != nil {
		log.Fatalf("Error formatting digest: %v", err)
	}

	log.Println("Dijester completed successfully")
}
