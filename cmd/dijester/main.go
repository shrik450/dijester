package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/shrik450/dijester/pkg/config"
	"github.com/shrik450/dijester/pkg/constants"
	"github.com/shrik450/dijester/pkg/fetcher"
	"github.com/shrik450/dijester/pkg/formatter"
	"github.com/shrik450/dijester/pkg/models"
	"github.com/shrik450/dijester/pkg/processor"
	"github.com/shrik450/dijester/pkg/source"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version and exit")
	configPath := flag.String("config", "", "Path to config file")
	outputDir := flag.String(
		"output-dir",
		"",
		"Path to output directory. If the config specifies a relative path, it will be relative to this directory.",
	)
	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s\n", constants.VERSION)
		return
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Dijester starting up...")

	if *configPath == "" {
		log.Fatal("No config file specified. Use -config to specify a config file.")
	}

	cfg, err := config.LoadFile(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	log.Printf("Loaded configuration from %s", *configPath)

	now := time.Now()
	formattedTitle, err := executeTmpl(cfg.Digest.Title, now)
	if err != nil {
		log.Printf("Error formatting title: %v, using original", err)
		formattedTitle = cfg.Digest.Title
	}

	digest := &models.Digest{
		Title:       formattedTitle,
		GeneratedAt: now,
		Articles:    make([]*models.Article, 0),
	}

	globalFetcher := fetcher.FromConfig(cfg.FetcherConfig)
	globalProcs, globalProcsOpts, err := processor.InitializeProcessors(cfg.ProcessorConfig)
	if err != nil {
		log.Fatalf("Error initializing global processors: %v", err)
	}

	outputFormatter, err := formatter.New(cfg.Digest.Format)
	if err != nil {
		log.Fatalf("Error initializing formatter: %v", err)
	}
	fmtOpts := formatter.OptionsFromConfig(cfg.Formatting, globalFetcher)

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

		if len(srcCfg.WordDenylist) > 0 {
			originalCount := len(articles)
			articles = source.FilterArticlesByWordDenylist(articles, srcCfg.WordDenylist)
			log.Printf(
				"Kept %d/%d articles from %s after word denylist filtering",
				len(articles),
				originalCount,
				src.Name(),
			)
		}

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

	if len(digest.Articles) == 0 {
		log.Println("No articles found. Exiting.")
		return
	}

	log.Printf("Fetched %d articles from all sources", len(digest.Articles))

	if cfg.Digest.DedupByURL {
		log.Println("Deduplicating articles by URL")
		originalCount := len(digest.Articles)
		digest.Articles = models.DeduplicateArticlesByURL(digest.Articles)
		log.Printf(
			"Kept %d/%d articles after URL deduplication",
			len(digest.Articles),
			originalCount,
		)
	}

	if len(cfg.Digest.SortBy) > 0 {
		log.Println("Sorting articles by configured properties")
		models.SortArticles(digest.Articles, cfg.Digest.SortBy)
		log.Println("Articles sorted successfully")
	}

	var finalOutputPath string

	formattedPath, err := executeTmpl(cfg.Digest.OutputPath, now)
	if err != nil {
		log.Printf("Error formatting output path: %v, using original", err)
		finalOutputPath = cfg.Digest.OutputPath
	} else {
		finalOutputPath = formattedPath
	}

	if !filepath.IsAbs(finalOutputPath) && *outputDir != "" {
		finalOutputPath = filepath.Join(*outputDir, finalOutputPath)
	}

	oDir := filepath.Dir(finalOutputPath)
	if err := os.MkdirAll(oDir, 0o755); err != nil {
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

type templateData struct {
	Now      time.Time
	Year     int
	Month    string
	MonthNum int
	Day      int
	Date     string
	Time     string
	TStamp   string
	DTLong   string
}

func executeTmpl(text string, t time.Time) (string, error) {
	if !strings.Contains(text, "{{") {
		return text, nil
	}

	data := templateData{
		Now:      t,
		Year:     t.Year(),
		Month:    t.Month().String(),
		MonthNum: int(t.Month()),
		Day:      t.Day(),
		Date:     t.Format(time.DateOnly),
		Time:     t.Format(time.TimeOnly),
		TStamp:   t.Format(time.RFC3339),
		DTLong:   t.Format(time.RubyDate),
	}

	tmpl, err := template.New("template").Parse(text)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return result.String(), nil
}
