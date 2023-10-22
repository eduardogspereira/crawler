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

func ExtractLinksFrom(htmlBody io.Reader) ([]*url.URL, error) {
	doc, err := html.Parse(htmlBody)
	if err != nil {
		return nil, err
	}

	var links []*url.URL

	var extractLinksFromNode func(n *html.Node)
	extractLinksFromNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == anchorTag {
			for _, a := range n.Attr {
				if a.Key == anchorHrefProperty {
					u, err := url.Parse(a.Val)
					if err != nil {
						continue
					}

					links = append(links, u)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractLinksFromNode(c)
		}
	}
	extractLinksFromNode(doc)

	return links, nil
}
