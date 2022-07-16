package fetch

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// A list of Sci-Hub mirrors.
var mirrors = []string{
	"sci-hub.se",
	"sci-hub.st",
	"sci-hub.ru",
}

// The Article holds data needed to download and write the article to disc.
// The DOI is used to fetch the title and PDF URL from Sci-Hub.
type Article struct {
	Doi   string
	Title string
	Url   string
}

// Fetch downloads articles from Sci-Hub from a list of supplied DOIs.
func Fetch(dois []string) error {

	articles := make([]Article, len(dois))
	for i, d := range dois {
		articles[i] = Article{Doi: d}
	}

	for _, a := range articles {
		res, err := sendRequest(&a)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		err = processRequest(res, &a)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		fmt.Printf("title: %q\n", a.Title)
		fmt.Printf("link: %q\n", a.Url)
	}

	return nil
}

// TODO: reimplement with timeouts and mirror switching
// sendRequest sends a GET request to 'Sci-Hub/DOI'.
// The request is sent to a different Sci-Hub mirror if the request times out.
// An error is returned if a valid response cannot be obtained.
func sendRequest(a *Article) (*http.Response, error) {
	res, err := http.Get("https://" + mirrors[0] + "/" + url.QueryEscape(a.Doi))
	if err != nil {
		return res, fmt.Errorf("%w", err)
	}

	return res, nil
}

// processRequest extracts the article title and URL from a HTML response.
func processRequest(res *http.Response, a *Article) error {
	if res.StatusCode > 399 {
		fmt.Errorf("%s", res.Status)
	}

	body, err := html.Parse(res.Body)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	// NOTE:
	//	With this approach we're evaluating all nodes each time we want to
	//	extract a piece of data. Ideally we would want to extract all data
	//	during the initial (and final) node evaluation.
	a.Title = getFromHtml(
		body,
		func(n *html.Node) bool {
			return n.Type == html.ElementNode && n.Data == "i"
		},
		func(n *html.Node) string {
			b := &bytes.Buffer{}
			b.WriteString(n.FirstChild.Data)
			ss := strings.SplitN(b.String(), ".", -1)

			var s string
			for i := 0; i < len(ss) - 2; i++ {
				if len(s) != 0 {
					s += "."
				}
				s += ss[i]
			}
			return s
		},
	)
	a.Url = getFromHtml(
		body,
		func(n *html.Node) bool {
			return n.Type == html.ElementNode && n.Data == "button"
		},
		func(n *html.Node) string {
			for _, atr := range n.Attr {
				if atr.Key == "onclick" {
					s := strings.TrimPrefix(atr.Val, "location.href='")
					s = strings.TrimSuffix(s, "?download=true'")
					return s
				}
			}
			return ""
		},
	)

	return nil
}

// getFromHtmlData walks an HTML tree and extracts data.
// If the current node in the HTML tree is selected by the 'selector', then
// data is extracted from the node by the 'exctractor', otherwise another node
// is selected. Both children and sibling nodes are walked.
func getFromHtml(
	n *html.Node,
	selector func(*html.Node) bool,
	extractor func(*html.Node) string,
) string {
	var s string
	if selector(n) {
		s += extractor(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s += getFromHtml(c, selector, extractor)
	}

	return s
}
