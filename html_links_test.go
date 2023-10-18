package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestGetLinksFrom_Success(t *testing.T) {
	linkA := "https://a.com"
	linkB := "https://b.com"
	linkC := "https://c.com"
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
	links, err := GetLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Contains(t, links, linkA)
	assert.Contains(t, links, linkB)
	assert.Contains(t, links, linkC)
}

func TestGetLinksFrom_NoLinks_Success(t *testing.T) {
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
		<body></body>
	</html>
`)

	r := strings.NewReader(htmlContent)
	links, err := GetLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Empty(t, links)
}

func TestGetLinksFrom_InvalidHTML_Success(t *testing.T) {
	linkA := "https://a.com"
	linkB := "https://b.com"
	htmlContent := fmt.Sprintf(`
	<a href="%s"/>
	>>>aDFSAFAS.>>>asd<a href="%s"/>fsa>!23213
`, linkA, linkB)

	r := strings.NewReader(htmlContent)
	links, err := GetLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Contains(t, links, linkA)
	assert.Contains(t, links, linkB)
}

func TestGetLinksFrom_LinksWithNoHref_Success(t *testing.T) {
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
	links, err := GetLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Empty(t, links)
}

func TestGetLinksFrom_RelativeLinks_Success(t *testing.T) {
	linkA := "/a"
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
		<body>
			<a href="%s"/>
		</body>
	</html>
`, linkA)
	r := strings.NewReader(htmlContent)
	links, err := GetLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Contains(t, links, linkA)
}

func TestGetLinksFrom_LinksWithAnchor_Success(t *testing.T) {
	linkA := "/a#section-51"
	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
		<body>
			<a href="%s"/>
		</body>
	</html>
`, linkA)
	r := strings.NewReader(htmlContent)
	links, err := GetLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Contains(t, links, linkA)
}

func TestGetLinksFrom_CSSAsset_Success(t *testing.T) {
	htmlContent := `
	.flagContent {
		width: auto;
		white-space: nowrap;
		display: inline-block;
		font-size: 0
	}
`
	r := strings.NewReader(htmlContent)
	links, err := GetLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Empty(t, links)
}

func TestGetLinksFrom_JSAsset_Success(t *testing.T) {
	htmlContent := `
	const random = () => 'random'
`
	r := strings.NewReader(htmlContent)
	links, err := GetLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Empty(t, links)
}

func TestGetLinksFrom_TextAsset_Success(t *testing.T) {
	htmlContent := "random text document"
	r := strings.NewReader(htmlContent)
	links, err := GetLinksFrom(io.NopCloser(r))
	assert.NoError(t, err)
	assert.Empty(t, links)
}
