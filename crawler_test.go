package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestTaskGetLinksForURL_Success(t *testing.T) {
	var linkA *url.URL
	var linkB *url.URL
	var linkE *url.URL

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		linkA = makeURLFor(t, fmt.Sprintf("http://%s/abc", r.Host))
		linkB = makeURLFor(t, fmt.Sprintf("http://%s/bca", r.Host))
		linkC := makeURLFor(t, "https://community.monzo.com/cab")
		linkE = makeURLFor(t, "/bca/abce/indiana/jones")
		htmlContent := fmt.Sprintf(`
			<a href="%s"/>
			<!DOCTYPE html>
			<html>
					<div>
						<div><a href="%s"/></div>
					</div>
				<body>
					<a href="%s"/>
				</body>
				<a href="%s"/>
			</html>
		`, linkA, linkB, linkC, linkE)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(htmlContent))
		assert.NoError(t, err)
	}))

	requestURL := makeURLFor(t, server.URL)

	crawler := NewCrawler(http.DefaultClient)

	failedURLs := make(chan *url.URL)
	linksForURLResults := make(chan *LinksForURLResult)

	go crawler.TaskGetLinksForURL(context.Background(), requestURL, linksForURLResults, failedURLs)

	results := <-linksForURLResults

	close(failedURLs)
	close(linksForURLResults)

	assert.Contains(t, results.links, linkA)
	assert.Contains(t, results.links, linkB)
	assert.Contains(t, results.links, requestURL.ResolveReference(linkE))
}

func TestTaskGetLinksForURL_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGatewayTimeout)
	}))

	requestURL := makeURLFor(t, server.URL)

	crawler := NewCrawler(&http.Client{Timeout: time.Nanosecond})

	failedURLs := make(chan *url.URL)
	linksForURLResults := make(chan *LinksForURLResult)

	go crawler.TaskGetLinksForURL(context.Background(), requestURL, linksForURLResults, failedURLs)

	failedURL := <-failedURLs

	close(failedURLs)
	close(linksForURLResults)

	assert.Equal(t, requestURL, failedURL)
}
