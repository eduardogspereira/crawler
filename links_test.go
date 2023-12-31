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
	links := ExtractLinksFrom(io.NopCloser(r))
	assert.Contains(t, links, linkA)
	assert.Contains(t, links, linkB)
	assert.Contains(t, links, linkC)
}

func TestExtractLinksFrom_MalformedPage_Success(t *testing.T) {
	linkA := makeURLFor(t, "https://a.com")
	linkB := makeURLFor(t, "https://b.com")
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
			<a href="%s"/>
			<div>
				<p>Unclosed <b>tag</p>
			</div>
		<body><a href="%s"/></body>
	</html>
`, linkA, linkB)

	r := strings.NewReader(htmlContent)
	links := ExtractLinksFrom(io.NopCloser(r))
	assert.Contains(t, links, linkA)
	assert.Contains(t, links, linkB)
}

func TestExtractLinksFrom_NoLinks_Success(t *testing.T) {
	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<body></body>
		</html>`)

	r := strings.NewReader(htmlContent)
	links := ExtractLinksFrom(io.NopCloser(r))
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
	links := ExtractLinksFrom(io.NopCloser(r))
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
		</html>`)
	r := strings.NewReader(htmlContent)
	links := ExtractLinksFrom(io.NopCloser(r))
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
		</html>`, linkA)
	r := strings.NewReader(htmlContent)
	links := ExtractLinksFrom(io.NopCloser(r))
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
		</html>`, linkA)
	r := strings.NewReader(htmlContent)
	links := ExtractLinksFrom(io.NopCloser(r))
	assert.Contains(t, links, linkA)
}

func TestExtractLinksFrom_CSSAsset_Success(t *testing.T) {
	cssContent := `
	.flagContent {
		width: auto;
		white-space: nowrap;
		display: inline-block;
		font-size: 0
	}`
	r := strings.NewReader(cssContent)
	links := ExtractLinksFrom(io.NopCloser(r))
	assert.Empty(t, links)
}

func TestExtractLinksFrom_JSAsset_Success(t *testing.T) {
	jsContent := `const random = () => 'random'`
	r := strings.NewReader(jsContent)
	links := ExtractLinksFrom(io.NopCloser(r))
	assert.Empty(t, links)
}

func TestExtractLinksFrom_TextAsset_Success(t *testing.T) {
	textContent := "random text document"
	r := strings.NewReader(textContent)
	links := ExtractLinksFrom(io.NopCloser(r))
	assert.Empty(t, links)
}

func TestFilterURLsBySubdomain_Success(t *testing.T) {
	linkA := makeURLFor(t, "https://abc.com/path-a")
	linkB := makeURLFor(t, "https://bca.com/path-d")
	linkC := makeURLFor(t, "https://abc.com/path-c")

	startURL := makeURLFor(t, "https://abc.com/path-j")

	links := FilterURLsBySubdomain(startURL, []*url.URL{linkA, linkB, linkC})
	assert.Contains(t, links, linkA)
	assert.NotContains(t, links, linkB)
	assert.Contains(t, links, linkC)
}

func TestFilterURLsBySubdomain_RelativeLinks_Success(t *testing.T) {
	linkA := makeURLFor(t, "https://abc.com/path-a")
	linkB := makeURLFor(t, "https://bca.com/path-d")
	linkC := makeURLFor(t, "/path-c")

	startURL := makeURLFor(t, "https://abc.com")

	links := FilterURLsBySubdomain(startURL, []*url.URL{linkA, linkB, linkC})
	assert.Contains(t, links, linkA)
	assert.NotContains(t, links, linkB)
	assert.Contains(t, links, startURL.ResolveReference(linkC))
}

func TestFilterURLsBySubdomain_RemoveNonHTTPProtocols_Success(t *testing.T) {
	linkA := makeURLFor(t, "https://abc.com/path-a")
	linkB := makeURLFor(t, "http://abc.com/path-b")
	linkC := makeURLFor(t, "/path-c")
	linkD := makeURLFor(t, "mailto:indiana-jones@abc.co")
	linkE := makeURLFor(t, "ftp://abc.com/path-d")
	linkF := makeURLFor(t, "//abc.com/path-b")

	startURL := makeURLFor(t, "https://abc.com")

	links := FilterURLsBySubdomain(startURL, []*url.URL{linkA, linkB, linkC, linkD, linkE, linkF})
	assert.Contains(t, links, linkA)
	assert.Contains(t, links, linkB)
	assert.Contains(t, links, startURL.ResolveReference(linkC))
	assert.NotContains(t, links, linkD)
	assert.NotContains(t, links, linkE)
	assert.Contains(t, links, startURL.ResolveReference(linkF))
}

func makeURLFor(t *testing.T, rawURL string) *url.URL {
	link, err := url.Parse(rawURL)
	assert.NoError(t, err)
	return link
}
