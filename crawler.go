package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type Crawler struct {
	linksByTargetURLs []*LinksByTargetURL
	errs              []error

	httpClient  *http.Client
	pageVisited map[string]bool
	workerPool  *WorkerPool
	m           sync.Mutex
}

type CrawlerParams struct {
	httpClient      *http.Client
	numberOfWorkers int
}

func NewCrawler(params *CrawlerParams) *Crawler {
	httpClient := http.DefaultClient
	if params != nil && params.httpClient != nil {
		httpClient = params.httpClient
	}

	numberOfWorkers := 100
	if params != nil && params.numberOfWorkers != 0 {
		numberOfWorkers = params.numberOfWorkers
	}

	return &Crawler{
		httpClient:  httpClient,
		pageVisited: make(map[string]bool),
		workerPool:  NewWorkerPool(numberOfWorkers),
	}
}

func (c *Crawler) GetAllLinksFor(ctx context.Context, targetURL *url.URL) ([]*LinksByTargetURL, []error) {
	c.MarkPageAsVisited(targetURL)
	c.workerPool.AddTask(targetURL)

	c.workerPool.ProcessTasks(func(nextTask interface{}) {
		nextTargetURL := nextTask.(*url.URL)
		linksForTargetURL, err := c.GetLinksForTargetURL(ctx, nextTargetURL)

		c.m.Lock()
		defer c.m.Unlock()
		if err != nil {
			c.errs = append(c.errs, err)
			return
		}
		c.linksByTargetURLs = append(c.linksByTargetURLs, linksForTargetURL)

		for _, l := range linksForTargetURL.links {
			if ok := c.MarkPageAsVisited(l); ok {
				c.workerPool.AddTask(l)
			}
		}
	})

	return c.linksByTargetURLs, c.errs
}

type LinksByTargetURL struct {
	links     []*url.URL
	targetURL *url.URL
}

type CrawlerError struct {
	targetURL *url.URL
	err       error
}

func (c CrawlerError) Error() string {
	return fmt.Sprintf("failed to extract links from %s: %s", c.targetURL.String(), c.err.Error())
}

func (c *Crawler) GetLinksForTargetURL(ctx context.Context, targetURL *url.URL) (*LinksByTargetURL, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL.String(), nil)
	if err != nil {
		return nil, &CrawlerError{
			err:       fmt.Errorf("failed to build request: %w", err),
			targetURL: targetURL,
		}
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, &CrawlerError{
			err:       fmt.Errorf("failed to build request: %w", err),
			targetURL: targetURL,
		}
	}

	return &LinksByTargetURL{
		targetURL: targetURL,
		links:     FilterURLsBySubdomain(targetURL, ExtractLinksFrom(response.Body)),
	}, nil
}

func (c *Crawler) MarkPageAsVisited(targetURL *url.URL) bool {
	_, pageAlreadyVisited := c.pageVisited[targetURL.Host+targetURL.Path]
	if pageAlreadyVisited {
		return false
	}

	c.pageVisited[targetURL.Host+targetURL.Path] = true
	return true
}
