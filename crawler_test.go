package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
)

func TestCrawler_GetAllLinksFor_Success(t *testing.T) {
	var m sync.Mutex
	var linkA *url.URL
	var linkB *url.URL
	var linkE *url.URL

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			m.Lock()
			defer m.Unlock()
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
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))

	targetURL := makeURLFor(t, server.URL)

	crawler := NewCrawler(nil)

	var linksForTargetURLs []*LinksByTargetURL
	onTargetURLProcessed := func(linksForTargetURL *LinksByTargetURL) {
		m.Lock()
		defer m.Unlock()
		linksForTargetURLs = append(linksForTargetURLs, linksForTargetURL)
	}

	var errs []error
	onError := func(err error) {
		m.Lock()
		defer m.Unlock()
		errs = append(errs, err)
	}

	crawler.GetAllLinksFor(context.Background(), targetURL, onTargetURLProcessed, onError)

	assert.Empty(t, errs)
	assert.Len(t, linksForTargetURLs, 4)
	for _, linksForTargetURL := range linksForTargetURLs {
		if linksForTargetURL.targetURL == targetURL {
			assert.Contains(t, linksForTargetURL.links, linkA)
			assert.Contains(t, linksForTargetURL.links, linkB)
			assert.Contains(t, linksForTargetURL.links, targetURL.ResolveReference(linkE))
		} else {
			assert.Empty(t, linksForTargetURL.links)
		}
	}
}

func TestCrawler_GetAllLinksFor_TwoPagesSuccess(t *testing.T) {
	var m sync.Mutex
	var linkA *url.URL
	var linkB *url.URL

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Lock()
		defer m.Unlock()

		var htmlContent string
		if r.URL.Path == "/" {
			linkA = makeURLFor(t, fmt.Sprintf("http://%s/abc", r.Host))
			htmlContent = fmt.Sprintf(`<a href="%s"/>`, linkA)
		} else if r.URL.Path == "/abc" {
			linkB = makeURLFor(t, fmt.Sprintf("http://%s/cba", r.Host))
			htmlContent = fmt.Sprintf(`<a href="%s"/>`, linkB)
		} else {
			htmlContent = ""
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(htmlContent))
		assert.NoError(t, err)
	}))

	targetURL := makeURLFor(t, server.URL)

	crawler := NewCrawler(nil)
	var linksForTargetURLs []*LinksByTargetURL
	onTargetURLProcessed := func(linksForTargetURL *LinksByTargetURL) {
		m.Lock()
		defer m.Unlock()
		linksForTargetURLs = append(linksForTargetURLs, linksForTargetURL)
	}

	var errs []error
	onError := func(err error) {
		m.Lock()
		defer m.Unlock()
		errs = append(errs, err)
	}

	crawler.GetAllLinksFor(context.Background(), targetURL, onTargetURLProcessed, onError)

	assert.Empty(t, errs)
	assert.Len(t, linksForTargetURLs, 3)
	for _, linksForTargetURL := range linksForTargetURLs {
		if linksForTargetURL.targetURL.String() == targetURL.String() {
			assert.Contains(t, linksForTargetURL.links, linkA)
		} else if linksForTargetURL.targetURL.String() == linkA.String() {
			assert.Contains(t, linksForTargetURL.links, linkB)
		} else {
			assert.Empty(t, linksForTargetURL.links)
		}
	}
}

func TestCrawler_GetAllLinksFor_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGatewayTimeout)
	}))

	targetURL := makeURLFor(t, server.URL)

	crawler := NewCrawler(&CrawlerParams{httpClient: &http.Client{Timeout: time.Nanosecond}})
	var linksForTargetURLs []*LinksByTargetURL
	onTargetURLProcessed := func(linksForTargetURL *LinksByTargetURL) {
		linksForTargetURLs = append(linksForTargetURLs, linksForTargetURL)
	}

	var errs []error
	onError := func(err error) {
		errs = append(errs, err)
	}

	crawler.GetAllLinksFor(context.Background(), targetURL, onTargetURLProcessed, onError)

	var crawlerError *CrawlerError
	errors.As(errs[0], &crawlerError)

	assert.Equal(t, targetURL, crawlerError.targetURL)
	assert.Error(t, crawlerError)
}
