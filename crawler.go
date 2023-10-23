package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

type Crawler struct {
	httpClient *http.Client

	linksForTargetURLResult []*LinksForTargetURL
	errorForTargetURLResult []*ErrorForTargetURL
	pageVisited             map[string]bool
	workerPool              *WorkerPool
	m                       sync.Mutex
}

func NewCrawler(httpClient *http.Client) *Crawler {
	return &Crawler{
		httpClient:  httpClient,
		pageVisited: make(map[string]bool),
		workerPool:  NewWorkerPool(100), //!TODO: set number of workers as a setting
	}
}

func (c *Crawler) GetAllLinksFor(ctx context.Context, targetURL *url.URL) ([]*LinksForTargetURL, []*ErrorForTargetURL) {
	c.MarkPageAsVisited(targetURL)
	c.workerPool.AddTask(targetURL)

	c.workerPool.ProcessTasks(func(nextTask interface{}) {
		nextTargetURL := nextTask.(*url.URL)
		linksForTargetURL, errorForTargetURL := c.TaskGetLinksForURL(ctx, nextTargetURL)
		c.m.Lock()
		defer c.m.Unlock()
		if errorForTargetURL != nil {
			c.errorForTargetURLResult = append(c.errorForTargetURLResult, errorForTargetURL)
			return
		}
		c.linksForTargetURLResult = append(c.linksForTargetURLResult, linksForTargetURL)
	})

	return c.linksForTargetURLResult, c.errorForTargetURLResult
}

type LinksForTargetURL struct {
	links     []*url.URL
	targetURL *url.URL
}

// !TODO: verify if it makes more sense to build a struct error
type ErrorForTargetURL struct {
	err       error
	targetURL *url.URL
}

func (c *Crawler) TaskGetLinksForURL(ctx context.Context, targetURL *url.URL) (*LinksForTargetURL, *ErrorForTargetURL) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL.String(), nil)
	if err != nil {
		return nil, &ErrorForTargetURL{
			err:       fmt.Errorf("failed to build request: %w", err),
			targetURL: targetURL,
		}
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, &ErrorForTargetURL{
			err:       fmt.Errorf("failed to build request: %w", err),
			targetURL: targetURL,
		}
	}

	linksForTargetURL := FilterURLsBySubdomain(targetURL, ExtractLinksFrom(response.Body))
	for _, l := range linksForTargetURL {
		if !c.PageAlreadyVisited(l) {
			c.MarkPageAsVisited(l)
			c.workerPool.AddTask(l)
		}
	}

	return &LinksForTargetURL{
		targetURL: targetURL,
		links:     linksForTargetURL,
	}, nil
}

func (c *Crawler) MarkPageAsVisited(targetURL *url.URL) {
	c.m.Lock()
	defer c.m.Unlock()
	c.pageVisited[targetURL.Host+targetURL.Path] = true
}

func (c *Crawler) PageAlreadyVisited(targetURL *url.URL) bool {
	c.m.Lock()
	defer c.m.Unlock()
	_, exists := c.pageVisited[targetURL.Host+targetURL.Path]
	return exists
}
