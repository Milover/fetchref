package fetch

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// htmlSelectorExtractor holds selector and extractor functions which work
// on HTML tree nodes, and a pointer to a data string.
type htmlSelectorExtractor struct {
	selector  func(*html.Node) bool
	extractor func(*html.Node) string

	data *strings.Builder
}

func newHse(
	s func(*html.Node) bool,
	e func(*html.Node) string,
) htmlSelectorExtractor {
	return htmlSelectorExtractor{
		selector:  s,
		extractor: e,
		data:      &strings.Builder{},
	}
}

// getFromHTML walks an HTML tree and extracts data.
// If the current node in the HTML tree is selected by the 'selector', then
// data is extracted from the node by the 'exctractor', otherwise another node
// is selected. Both children and sibling nodes are walked.
func getFromHTML(n *html.Node, es []htmlSelectorExtractor) {
	for _, e := range es {
		if e.selector(n) {
			e.data.WriteString(e.extractor(n))
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getFromHTML(c, es)
	}
}

// selectTitleNode selects the HTML node containing the article (download) URL.
func selectURLNode(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "button"
}

// extractUrl extracts the article (download) URL from a HTML body.
func extractURL(n *html.Node) string {
	for _, atr := range n.Attr {
		if atr.Key == "onclick" {
			re := regexp.MustCompile(`location.href='(.*)'`)
			match := re.FindStringSubmatch(atr.Val)
			return match[len(match)-1]
		}
	}
	return ""
}
