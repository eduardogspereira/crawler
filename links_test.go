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
	linkA, err := url.Parse("https://a.com")
	assert.NoError(t, err)
	linkB, err := url.Parse("https://b.com")
	assert.NoError(t, err)
	linkC, err := url.Parse("https://c.com")
	assert.NoError(t, err)
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
	linkA, err := url.Parse("https://a.com")
	assert.NoError(t, err)
	linkB, err := url.Parse("https://b.com")
	assert.NoError(t, err)
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
	linkA, err := url.Parse("/a")
	assert.NoError(t, err)
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
	linkA, err := url.Parse("/a#section-51")
	assert.NoError(t, err)
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
