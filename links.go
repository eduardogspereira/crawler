package main

import (
	"golang.org/x/net/html"
	"io"
	"net/url"
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
			if token.Data == anchorTag {
				for _, attr := range token.Attr {
					if attr.Key == anchorHrefProperty {
						u, err := url.Parse(attr.Val)
						if err != nil {
							continue
						}
						links = append(links, u)
					}
				}
			}
		}
	}
}

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
