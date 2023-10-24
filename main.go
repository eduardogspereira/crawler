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
		retryAttempts:   params.numberOfRetries,
	}
	crawler := NewCrawler(crawlerParams)

	onTargetURLProcessed := func(linksForTargetURL *LinksByTargetURL) {
		log.Printf("URL -> %s: LINKS -> %s\n", linksForTargetURL.targetURL, linksForTargetURL.links)
	}

	var errs []error
	onError := func(err error) {
		errs = append(errs, err)
	}

	crawler.GetAllLinksFor(context.Background(), params.targetURL, onTargetURLProcessed, onError)

	for _, err = range errs {
		log.Println(err)
	}
}

type parameters struct {
	numberOfWorkers int
	timeout         time.Duration
	targetURL       *url.URL
	numberOfRetries uint
}

func parseCommandLineFlags() (*parameters, error) {
	workers := pflag.IntP("workers", "w", 100, "Number of workers")
	timeout := pflag.IntP("timeout", "t", 30, "HTTP timeout (seconds)")
	targetURL := pflag.StringP("url", "u", "", "Target URL")
	retries := pflag.UintP("retries", "r", 3, "Number of task retries")

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
		numberOfRetries: *retries,
	}, nil
}
