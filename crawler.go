package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Crawler struct {
	httpClient *http.Client

	pageVisited map[string]bool
}

func NewCrawler(httpClient *http.Client) *Crawler {
	return &Crawler{
		httpClient:  httpClient,
		pageVisited: make(map[string]bool),
	}
}

func (c *Crawler) GetAllLinksFor(ctx context.Context, url *url.URL) {
	workerPool := NewWorkerPool(100) //!TODO: set number of workers as a setting
	workerPool.AddTask(url)
}

func (c *Crawler) GetLinksFromURL(ctx context.Context, url *url.URL) ([]*url.URL, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	return FilterURLsBySubdomain(url, ExtractLinksFrom(response.Body)), nil
}
