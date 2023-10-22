package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/url"
	"strings"
	"testing"
)

func TestExtractLinksFrom_Success(t *testing.T) {
	linkA := makeURLFor(t, "https://a.com")
	linkB := makeURLFor(t, "https://b.com")
	linkC := makeURLFor(t, "https://c.com")
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
	</html>
`, linkA, linkB, linkC)

	r := strings.NewReader(htmlContent)
	links, err := ExtractLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Contains(t, links, linkA)
	assert.Contains(t, links, linkB)
	assert.Contains(t, links, linkC)
}

func TestExtractLinksFrom_NoLinks_Success(t *testing.T) {
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
		<body></body>
	</html>
`)

	r := strings.NewReader(htmlContent)
	links, err := ExtractLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Empty(t, links)
}

func TestExtractLinksFrom_InvalidHTML_Success(t *testing.T) {
	linkA := makeURLFor(t, "https://a.com")
	linkB := makeURLFor(t, "https://b.com")
	htmlContent := fmt.Sprintf(`
	<a href="%s"/>
	>>>aDFSAFAS.>>>asd<a href="%s"/>fsa>!23213
`, linkA, linkB)

	r := strings.NewReader(htmlContent)
	links, err := ExtractLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Contains(t, links, linkA)
	assert.Contains(t, links, linkB)
}

func TestExtractLinksFrom_LinksWithNoHref_Success(t *testing.T) {
	htmlContent := fmt.Sprintf(`
	<a/>
	<!DOCTYPE html>
	<html>
			<div>
				<div><a/></div>
			</div>
		<body>
			<a/>
		</body>
	</html>
`)
	r := strings.NewReader(htmlContent)
	links, err := ExtractLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Empty(t, links)
}

func TestExtractLinksFrom_RelativeLinks_Success(t *testing.T) {
	linkA := makeURLFor(t, "/a")
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
		<body>
			<a href="%s"/>
		</body>
	</html>
`, linkA)
	r := strings.NewReader(htmlContent)
	links, err := ExtractLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Contains(t, links, linkA)
}

func TestExtractLinksFrom_LinksWithAnchor_Success(t *testing.T) {
	linkA := makeURLFor(t, "/a#section-51")
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
		<body>
			<a href="%s"/>
		</body>
	</html>
`, linkA)
	r := strings.NewReader(htmlContent)
	links, err := ExtractLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Contains(t, links, linkA)
}

func TestExtractLinksFrom_CSSAsset_Success(t *testing.T) {
	htmlContent := `
	.flagContent {
		width: auto;
		white-space: nowrap;
		display: inline-block;
		font-size: 0
	}
`
	r := strings.NewReader(htmlContent)
	links, err := ExtractLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Empty(t, links)
}

func TestExtractLinksFrom_JSAsset_Success(t *testing.T) {
	htmlContent := `
	const random = () => 'random'
`
	r := strings.NewReader(htmlContent)
	links, err := ExtractLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Empty(t, links)
}

func TestExtractLinksFrom_TextAsset_Success(t *testing.T) {
	htmlContent := "random text document"
	r := strings.NewReader(htmlContent)
	links, err := ExtractLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Empty(t, links)
}

func TestFilterURLsBySubdomain_Success(t *testing.T) {
	linkA := makeURLFor(t, "https://abc.com/path-a")
	linkB := makeURLFor(t, "https://bca.com/path-d")
	linkC := makeURLFor(t, "https://abc.com/path-c")

	startURL := makeURLFor(t, "https://abc.com/path-j")

	links := FilterURLsBySubdomain(startURL, []*url.URL{linkA, linkB, linkC})
	assert.Contains(t, links, linkA)
	assert.Contains(t, links, linkC)
}

func TestFilterURLsBySubdomain_RelativeLinks_Success(t *testing.T) {
	linkA := makeURLFor(t, "https://abc.com/path-a")
	linkB := makeURLFor(t, "https://bca.com/path-d")
	linkC := makeURLFor(t, "/path-c")

	startURL := makeURLFor(t, "https://abc.com")

	links := FilterURLsBySubdomain(startURL, []*url.URL{linkA, linkB, linkC})
	assert.Contains(t, links, linkA)
	assert.Contains(t, links, startURL.ResolveReference(linkC))
}

func makeURLFor(t *testing.T, rawURL string) *url.URL {
	link, err := url.Parse(rawURL)
	assert.NoError(t, err)
	return link
}
