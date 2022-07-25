package fetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/Milover/fetchpaper/internal/article"
	"golang.org/x/net/html"
)

// A list of Sci-Hub mirrors.
var mirrors = []string{
	"sci-hub.se",
	"sci-hub.st",
	"sci-hub.ru",
}

// htmlSelectorExtractor holds selector and extractor functions which work
// on HTML tree nodes, and a pointer to a data string.
type htmlSelectorExtractor struct {
	selector  func(*html.Node) bool
	extractor func(*html.Node) string

	data *string
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

// Fetch downloads articles from Sci-Hub from a list of supplied DOIs.
func Fetch(dois []string) error {

	articles := make([]article.Article, len(dois))
	for i, d := range dois {
		articles[i] = article.Article{Doi: d}
		articles[i].GeneratorFunc(article.SnakeCaseGenerator)
	}

	for _, a := range articles {
		err := doInfoRequest(&a)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		if err = doDownloadRequest(a); err != nil {
			return fmt.Errorf("%w", err)
		}

		// TODO: remove
		fmt.Printf("title: %q\n", a.Title)
		fmt.Printf("link: %q\n", a.Url.String())
		fmt.Printf("file: %q\n", a.GenerateFileName())
	}

	return nil
}

// sendGetRequest sends a GET request to 'Sci-Hub/DOI'.
// The request is sent to a different Sci-Hub mirror if the request times out.
// An error is returned if a valid response cannot be obtained.
// TODO: reimplement with timeouts
func sendGetRequest(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return res, fmt.Errorf("%w", err)
	}
	if res.StatusCode > 399 {
		return res, fmt.Errorf("%s", res.Status)
	}

	return res, nil
}

// processRequest extracts the article title and URL from a HTML response.
func doInfoRequest(a *article.Article) error {
	// TODO: add mirror switching
	m := mirrors[0]

	reqUrl := &url.URL{
		Scheme: "https",
		Host:   m,
		Path:   a.Doi,
	}

	res, err := sendGetRequest(reqUrl.String())
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	body, err := html.Parse(res.Body)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	// define selectors/extractors
	var articleUri string
	ses := []htmlSelectorExtractor{
		{
			selector: func(n *html.Node) bool {
				return n.Type == html.ElementNode && n.Data == "i"
			},
			// extract title
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
			// extract download link
			extractor: func(n *html.Node) string {
				for _, atr := range n.Attr {
					if atr.Key == "onclick" {
						re := regexp.MustCompile(`location.href='(.*)'`)
						match := re.FindStringSubmatch(atr.Val)
						return match[len(match)-1]
					}
				}
				return ""
			},
			data: &articleUri,
		},
	}
	getFromHtml(body, ses)

	// finalize the url
	a.Url, err = url.Parse(articleUri)
	if err != nil {
		fmt.Errorf("%w", err)
	}
	a.Url.Scheme = "https"
	if len(a.Url.Host) == 0 {
		a.Url.Host = m
	}

	return nil
}

// downloadArticle downloads the article and writes a PDF to disc. The download
// URL and the file name are retrieved from the Article.
func doDownloadRequest(a article.Article) error {
	out, err := os.Create(a.GenerateFileName() + ".pdf")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	res, err := sendGetRequest(a.Url.String())
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	if _, err := io.Copy(out, res.Body); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
