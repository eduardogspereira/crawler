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

func TestCrawler_GetAllLinksFor_Success(t *testing.T) {
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

	targetURL := makeURLFor(t, server.URL)

	crawler := NewCrawler(http.DefaultClient)
	linksForTargetURLs, errorForTargetURLs := crawler.GetAllLinksFor(context.Background(), targetURL)

	assert.Nil(t, errorForTargetURLs)
	assert.Len(t, linksForTargetURLs, 4)
	assert.Contains(t, linksForTargetURLs[0].links, linkA)
	assert.Contains(t, linksForTargetURLs[0].links, linkB)
	assert.Contains(t, linksForTargetURLs[0].links, targetURL.ResolveReference(linkE))
}

func TestCrawler_GetAllLinksFor_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGatewayTimeout)
	}))

	targetURL := makeURLFor(t, server.URL)

	crawler := NewCrawler(&http.Client{Timeout: time.Nanosecond})
	_, errorForTargetURLs := crawler.GetAllLinksFor(context.Background(), targetURL)

	assert.Equal(t, targetURL, errorForTargetURLs[0].targetURL)
	assert.Error(t, errorForTargetURLs[0].err)
}
