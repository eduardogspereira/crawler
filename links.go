package main

import (
	"io"
	"net/url"

	"golang.org/x/net/html"
)

const (
	anchorTag          = "a"
	anchorHrefProperty = "href"
)

func ExtractLinksFrom(htmlBody io.Reader) []*url.URL {
	var links []*url.URL

	tokenizer := html.NewTokenizer(htmlBody)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data != anchorTag {
				continue
			}
			for _, attr := range token.Attr {
				if attr.Key != anchorHrefProperty {
					continue
				}

				u, err := url.Parse(attr.Val)
				if err != nil {
					continue
				}
				links = append(links, u)
			}
		}
	}
}

// Currently this function doesn't care about what is the file extension of the
// link. Potentially, one improvement that I could do for this part is:
//   - Add a list of extensions that I want to be removed (e.g. *.pdf, *.txt, *.jpg, etc.)
//   - Add functions to extract links for a given extension (e.g. 'func ExtractLinksFromHTML',
//     'func ExtractLinksFromPDF', 'func ExtractLinksFromTXT', etc.)
func FilterURLsBySubdomain(domain *url.URL, links []*url.URL) []*url.URL {
	var filteredURLs []*url.URL

	for _, link := range links {
		if link.Scheme != "" && link.Scheme != "http" && link.Scheme != "https" {
			continue
		}

		if link.Host == "" {
			filteredURLs = append(filteredURLs, domain.ResolveReference(link))
			continue
		}

		if domain.Host == link.Host {
			if link.Scheme == "" {
				filteredURLs = append(filteredURLs, domain.ResolveReference(link))
				continue
			}

			filteredURLs = append(filteredURLs, link)
		}
	}

	return filteredURLs
}
