package main

import (
	"golang.org/x/net/html"
	"io"
)

const (
	anchorTag          = "a"
	anchorHrefProperty = "href"
)

func ExtractLinksFrom(htmlBody io.Reader) ([]string, error) {
	doc, err := html.Parse(htmlBody)
	if err != nil {
		return nil, err
	}

	var links []string

	var extractLinksFromNode func(n *html.Node)
	extractLinksFromNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == anchorTag {
			for _, a := range n.Attr {
				if a.Key == anchorHrefProperty {
					links = append(links, a.Val)
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
