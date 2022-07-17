package fetch

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode"

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
	Doi      string
	Title    string
	FileName string
	Url      string
}

// FileName generates and caches the file name of the article.
// All punctuation, spaces and control codes are replaces by '_'s, which are
// squeezed, and a '.pdf' extension is added.
func (a *Article) GenerateFileName() {
	if len(a.FileName) != 0 || len(a.Title) == 0 {
		return
	}

	var b strings.Builder
	var cache rune

	b.Grow(len(a.Title))
	for _, r := range a.Title {
		if unicode.In(r, unicode.P, unicode.Z, unicode.Cc) {
			if cache == '_' {
				continue
			}
			cache = '_'
		} else {
			cache = unicode.ToLower(r)
		}
		b.WriteRune(cache)
	}
	b.WriteString(".pdf")

	a.FileName = b.String()
}

// Reset resets all article data.
func (a *Article) Reset() {
	*a = Article{}
}

// htmlSelectorExtractor holds selector and extractor functions which work
// on HTML tree nodes, and a pointer to a data string.
type htmlSelectorExtractor struct {
	selector  func(*html.Node) bool
	extractor func(*html.Node) string

	data *string
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
		a.GenerateFileName()
		fmt.Printf("file: %q\n", a.FileName)
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

	// define selectors/extractors
	ses := []htmlSelectorExtractor{
		{
			selector: func(n *html.Node) bool {
				return n.Type == html.ElementNode && n.Data == "i"
			},
			extractor: func(n *html.Node) string {
				var b strings.Builder
				bufSize := 150 // educated guess

				b.Grow(bufSize)
				b.WriteString(n.FirstChild.Data)
				ss := strings.SplitAfterN(b.String(), ".", -1)

				b.Reset()
				b.Grow(bufSize)
				for i := 0; i < len(ss)-2; i++ {
					b.WriteString(ss[i])
				}
				return strings.TrimSuffix(b.String(), ".")
			},
			data: &a.Title,
		},
		{
			selector: func(n *html.Node) bool {
				return n.Type == html.ElementNode && n.Data == "button"
			},
			extractor: func(n *html.Node) string {
				for _, atr := range n.Attr {
					if atr.Key == "onclick" {
						s := strings.TrimPrefix(atr.Val, "location.href='")
						s = strings.TrimSuffix(s, "?download=true'")
						return s
					}
				}
				return ""
			},
			data: &a.Url,
		},
	}
	getFromHtml(body, ses)

	return nil
}

// getFromHtmlData walks an HTML tree and extracts data.
// If the current node in the HTML tree is selected by the 'selector', then
// data is extracted from the node by the 'exctractor', otherwise another node
// is selected. Both children and sibling nodes are walked.
func getFromHtml(n *html.Node, es []htmlSelectorExtractor) {
	for _, e := range es {
		if e.selector(n) {
			*e.data += e.extractor(n)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getFromHtml(c, es)
	}
}
