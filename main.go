package main

import (
	"context"
	"errors"
	"github.com/spf13/pflag"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	params, err := parseCommandLineFlags()
	if err != nil {
		log.Fatal(err)
	}

	crawlerParams := &CrawlerParams{
		httpClient:      &http.Client{Timeout: params.timeout},
		numberOfWorkers: params.numberOfWorkers,
	}
	crawler := NewCrawler(crawlerParams)
	linksByTargetURLs, errs := crawler.GetAllLinksFor(context.Background(), params.targetURL)
	log.Println(errs)
	log.Println(linksByTargetURLs)
}

type parameters struct {
	numberOfWorkers int
	timeout         time.Duration
	targetURL       *url.URL
}

func parseCommandLineFlags() (*parameters, error) {
	workers := pflag.IntP("workers", "w", 100, "Number of workers")
	timeout := pflag.IntP("timeout", "t", 30, "HTTP timeout (seconds)")
	targetURL := pflag.StringP("url", "u", "", "Target URL")

	pflag.Parse()

	if *targetURL == "" {
		return nil, errors.New("url parameters is required")
	}

	u, err := url.Parse(*targetURL)
	if err != nil {
		return nil, err
	}

	return &parameters{
		targetURL:       u,
		timeout:         time.Duration(*timeout) * time.Second,
		numberOfWorkers: *workers,
	}, nil
}
