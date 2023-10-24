package main

import (
	"context"
	"fmt"
	"github.com/avast/retry-go/v4"
	"net/http"
	"net/url"
	"sync"
)

type Crawler struct {
	httpClient    *http.Client
	pageVisited   map[string]bool
	workerPool    *WorkerPool
	m             sync.Mutex
	retryAttempts uint
}

type CrawlerParams struct {
	httpClient      *http.Client
	numberOfWorkers int
	retryAttempts   uint
}

func NewCrawler(params *CrawlerParams) *Crawler {
	return &Crawler{
		httpClient:    params.httpClient,
		pageVisited:   make(map[string]bool),
		workerPool:    NewWorkerPool(params.numberOfWorkers),
		retryAttempts: params.retryAttempts,
	}
}

func (c *Crawler) GetAllLinksFor(
	ctx context.Context,
	targetURL *url.URL,
	onTargetURLProcessed func(*LinksByTargetURL),
	onError func(error),
) {
	c.MarkPageAsVisited(targetURL)
	c.workerPool.AddTask(targetURL)

	c.workerPool.ProcessTasks(func(nextTargetURL interface{}) {
		linksForTargetURL, err := c.GetLinksForTargetURL(ctx, nextTargetURL.(*url.URL))
		if err != nil {
			onError(err)
			return
		}
		onTargetURLProcessed(linksForTargetURL)

		for _, l := range linksForTargetURL.links {
			if ok := c.MarkPageAsVisited(l); ok {
				c.workerPool.AddTask(l)
			}
		}
	})
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

	var response *http.Response
	err = retry.Do(func() error {
		response, err = c.httpClient.Do(request)
		return err
	}, retry.Context(ctx), retry.Attempts(c.retryAttempts), retry.LastErrorOnly(true))
	if err != nil {
		return nil, &CrawlerError{
			err:       fmt.Errorf("failed to make the request: %w", err),
			targetURL: targetURL,
		}
	}

	// I decided to not check if the Status Code from the response is in the range of
	// 2XX as some pages return links even when the response is not success (e.g. https://monzo.com/non-existent-page/)

	return &LinksByTargetURL{
		targetURL: targetURL,
		links:     FilterURLsBySubdomain(targetURL, ExtractLinksFrom(response.Body)),
	}, nil
}

func (c *Crawler) MarkPageAsVisited(targetURL *url.URL) bool {
	c.m.Lock()
	defer c.m.Unlock()
	_, pageAlreadyVisited := c.pageVisited[targetURL.Host+targetURL.Path]
	if pageAlreadyVisited {
		return false
	}

	c.pageVisited[targetURL.Host+targetURL.Path] = true
	return true
}
