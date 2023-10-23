package main

import (
	"context"
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

func (c *Crawler) GetAllLinksFor(ctx context.Context, targetURL *url.URL) {
	workerPool := NewWorkerPool(100) //!TODO: set number of workers as a setting
	workerPool.AddTask(targetURL)
}

type LinksForURLResult struct {
	links     []*url.URL
	targetURL *url.URL
}

func (c *Crawler) TaskGetLinksForURL(ctx context.Context, targetURL *url.URL, linksForURLResults chan<- *LinksForURLResult, failedURLs chan<- *url.URL) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL.String(), nil)
	if err != nil {
		failedURLs <- targetURL
		return
		//return nil, fmt.Errorf("failed to build request: %w", err) //!TODO: ADD DEBUG LOG
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		failedURLs <- targetURL
		return
		//return nil, fmt.Errorf("failed to do request: %w", err) //!TODO: ADD DEBUG LOG
	}

	linksForURLResults <- &LinksForURLResult{
		targetURL: targetURL,
		links:     FilterURLsBySubdomain(targetURL, ExtractLinksFrom(response.Body)),
	}
}
