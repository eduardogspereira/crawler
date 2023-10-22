package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetLinksFromURL_Success(t *testing.T) {
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
	links, err := crawler.GetLinksFromURL(context.Background(), requestURL)

	assert.NoError(t, err)
	assert.Contains(t, links, linkA)
	assert.Contains(t, links, linkB)
	assert.Contains(t, links, requestURL.ResolveReference(linkE))
}
