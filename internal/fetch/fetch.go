package fetch

import (
	"fmt"
	"net/http"

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

func Fetch(dois []string) error {

	// initialize a list of articles
	articles := make([]Article, len(dois))
	for i, d := range dois {
		articles[i] = Article{Doi: d}
	}

	// for each article
	//		send a GET to one of the mirrors
	//		parse html
	//			if found(title, url)
	//				download and write to disc
	//			else
	//				return an error
	for _, a := range articles {
		// FIXME: should try mirrors in order
		res, err := http.Get("https://" + mirrors[0] + "/" + a.Doi)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		body, err := html.Parse(res.Body)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "button" {
				for _, a := range n.Attr {
					if a.Key == "onclick" {
						fmt.Println(a.Val)
						break
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(body)
	}

	return nil
}
