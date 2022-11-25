package fetch

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/Milover/fetchpaper/internal/article"
	"github.com/Milover/fetchpaper/internal/metainfo"
	"go.uber.org/ratelimit"
	"golang.org/x/net/html"
)

var (
	// Global HTTP request timeout.
	GlobalReqTimeout  = 3 * time.Second
	GlobalRateLimiter = ratelimit.New(50)

	// A list of Sci-Hub mirrors.
	mirrors = []string{
		"sci-hub.se",
		"sci-hub.st",
		"sci-hub.ru",
	}
)

// Fetch downloads articles from Sci-Hub from a list of supplied DOIs.
func Fetch(dois []string) error {

	if len(dois) == 0 {
		return nil
	}

	ch := make(chan *article.Article, len(dois))
	syn := make(chan bool, len(dois))
	good := true

	for _, d := range dois {
		a := article.Article{Doi: d}
		a.GeneratorFunc(article.SnakeCaseGenerator)

		go func() {
			if err := doInfoRequest(&a); err != nil {
				ch <- nil
				good = false
				log.Printf("%v: %v", a.Doi, err)
			} else {
				ch <- &a
			}
			syn <- true
			if len(syn) == len(dois) {
				close(ch)
			}
		}()
	}

	var wg sync.WaitGroup
	for {
		a, ok := <-ch
		if !ok {
			wg.Wait()
			if !good {
				return fmt.Errorf("errors occured")
			}
			return nil
		}
		if a != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := doDownloadRequest(*a); err != nil {
					good = false
					log.Printf("%v: %v", a.Doi, err)
				}
			}()
		}
	}
}

// sendGetRequest sends a GET request to 'Sci-Hub/DOI'.
// The request is sent to a different Sci-Hub mirror if the request times out.
// An error is returned if a valid response cannot be obtained.
func sendGetRequest(ctx context.Context, url string) (*http.Response, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	req.Header.Set("User-Agent", metainfo.HTTPUserAgent)

	// sendGetRequest can be called from multiple threads
	GlobalRateLimiter.Take()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, fmt.Errorf("%w", err)
	}
	if res.StatusCode > 399 {
		return res, fmt.Errorf("%s", res.Status)
	}

	return res, nil
}

// doInfoRequestFromMirror requests article info from a Sci-Hub mirror
// and parses the article title and download URL from the response HTML.
func doInfoRequestFromMirror(a *article.Article, mirror string) error {
	u := &url.URL{
		Scheme: "https",
		Host:   mirror,
		Path:   a.Doi,
	}
	ctx, cncl := context.WithTimeout(context.Background(), GlobalReqTimeout)
	defer cncl()

	res, err := sendGetRequest(ctx, u.String())
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	body, err := html.Parse(res.Body)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	// extract title and url from parsed HTML
	ses := []htmlSelectorExtractor{
		newHse(selectTitleNode, extractTitle),
		newHse(selectUrlNode, extractUrl),
	}
	getFromHtml(body, ses)

	// set and check title/body, clean up if necessary
	a.Title = ses[0].data.String()
	if len(a.Title) == 0 {
		a.Title = a.Doi
		log.Printf("%v: could not extract title", a.Doi)
	}

	a.Url, err = url.Parse(ses[1].data.String())
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	a.Url.Scheme = "https"
	if len(a.Url.Host) == 0 {
		a.Url.Host = res.Request.URL.Host
	}

	return nil
}

// doInfoRequest requests article info from Sci-Hub, and updates the article
// if successfull. On a failed request, another Sci-Hub mirror is chosen,
// until all mirrors have been exhausted.
func doInfoRequest(a *article.Article) error {
	for i, m := range mirrors {
		if err := doInfoRequestFromMirror(a, m); err != nil {
			log.Printf("%v: %v", a.Doi, err)
			// fail if there are no more mirrors to try
			if i == len(mirrors)-1 {
				return fmt.Errorf("could not get article info")
			}
		} else {
			break
		}
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

	ctx, cncl := context.WithTimeout(context.Background(), GlobalReqTimeout)
	defer cncl()

	res, err := sendGetRequest(ctx, a.Url.String())
	if err != nil {
		os.Remove(out.Name())
		return fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	if _, err := io.Copy(out, res.Body); err != nil {
		os.Remove(out.Name())
		return fmt.Errorf("%w", err)
	}

	return nil
}
