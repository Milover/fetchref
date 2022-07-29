package fetch

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Milover/fetchpaper/internal/article"
	"golang.org/x/net/html"
)

// A list of Sci-Hub mirrors.
var mirrors = []string{
	"sci-hub.se",
	"sci-hub.st",
	"sci-hub.ru",
}

// Fetch downloads articles from Sci-Hub from a list of supplied DOIs.
func Fetch(dois []string) error {

	ch := make(chan *article.Article, len(dois))
	sync := make(chan bool, len(dois))
	e := false

	for _, d := range dois {
		a := article.Article{Doi: d}
		a.GeneratorFunc(article.SnakeCaseGenerator)

		go func() {
			err := doInfoRequest(&a)
			if err != nil {
				ch <- nil
				e = true
				log.Println("%v", err)
			}
			ch <- &a
			sync <- true
			if len(sync) == len(dois) {
				close(ch)
			}
		}()
	}

	for {
		select {
		case a, ok := <-ch:
			if !ok {
				if e {
					return fmt.Errorf("error")
				}
				return nil
			}
			if a != nil {
				go func() {
					if err := doDownloadRequest(*a); err != nil {
						e = true
						log.Println("%v", err)
					}
				}()
			}
		}
	}

	return nil
}

// sendGetRequest sends a GET request to 'Sci-Hub/DOI'.
// The request is sent to a different Sci-Hub mirror if the request times out.
// An error is returned if a valid response cannot be obtained.
func sendGetRequest(ctx context.Context, url string) (*http.Response, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, fmt.Errorf("%w", err)
	}
	if res.StatusCode > 399 {
		return res, fmt.Errorf("%s", res.Status)
	}

	return res, nil
}

func infoRequestFromMirror(
	ctx context.Context,
	a article.Article,
) (*http.Response, error) {
	u := &url.URL{
		Scheme: "https",
		Path:   a.Doi,
	}

	for _, m := range mirrors {
		u.Host = m

		res, err := sendGetRequest(ctx, u.String())
		if err != nil {
			var e url.Error
			if errors.Is(err, &e) {
				if ue := err.(*url.Error); ue.Timeout() {
					log.Println("connection timed out: %v", err)
					continue
				}
			} else {
				return res, fmt.Errorf("%w", err)
			}
		}

		return res, err
	}

	panic("never reached")
}

// doInfoRequest extracts the article title and URL from a HTML response.
func doInfoRequest(a *article.Article) error {

	ctx, cncl := context.WithTimeout(context.Background(), time.Second*3)
	defer cncl()

	res, err := infoRequestFromMirror(ctx, *a)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	body, err := html.Parse(res.Body)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	ses := []htmlSelectorExtractor{
		newHse(selectTitleNode, extractTitle),
		newHse(selectUrlNode, extractUrl),
	}
	getFromHtml(body, ses)

	a.Title = ses[0].data.String()
	a.Url, err = url.Parse(ses[1].data.String())
	if err != nil {
		fmt.Errorf("%w", err)
	}
	a.Url.Scheme = "https"
	if len(a.Url.Host) == 0 {
		a.Url.Host = res.Request.URL.Host
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

	ctx, cncl := context.WithTimeout(context.Background(), time.Second*3)
	defer cncl()

	res, err := sendGetRequest(ctx, a.Url.String())
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	if _, err := io.Copy(out, res.Body); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
