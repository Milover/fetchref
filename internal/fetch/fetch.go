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
	"github.com/Milover/fetchpaper/internal/crossref"
	"github.com/Milover/fetchpaper/internal/metainfo"
	"go.uber.org/ratelimit"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

var (
	// GlobalReqTimeout is the global HTTP request timeout.
	GlobalReqTimeout = 3 * time.Second

	// GlobalRateLimiter is the global outgoing HTTP request limiter.
	GlobalRateLimiter = ratelimit.New(50)

	// CitationFormat is the citation output format.
	CitationFormat = crossref.BibTeX

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
	g := new(errgroup.Group)
	var wg sync.WaitGroup
	wg.Add(len(dois))

	for _, d := range dois {
		a := article.Article{Doi: d}
		a.GeneratorFunc(article.SnakeCaseGenerator)

		// GET info from Sci-Hub
		g.Go(func() error {
			defer wg.Done()
			if err := doInfoRequest(&a); err != nil {
				log.Printf("%v: %v", a.Doi, err)
				return err
			}
			ch <- &a
			return nil
		})

		// GET citation from Crossref
		// WARNING: the public pool is often more responsive than the polite one
		// WARNING: doCrossrefRequest should take from 'ch' even though it's
		// not necessary as per the current implementation.
		g.Go(func() error {
			if err := doCrossrefCitationRequest(&a, CitationFormat); err != nil {
				log.Printf("%v: %v", a.Doi, err)
				return err
			}
			return nil
		})

		// GET article PDF
		g.Go(func() error {
			if a, ok := <-ch; ok {
				if err := doDownloadRequest(a); err != nil {
					log.Printf("%v: %v", a.Doi, err)
					return err
				}
			}
			return nil
		})
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	if err := g.Wait(); err != nil {
		return fmt.Errorf("errors occurred")
	}
	return nil
}

// sendGetRequest sends a GET request to the specified URL.
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
	defer res.Body.Close()

	body, err := html.Parse(res.Body)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// extract title and url from parsed HTML
	ses := []htmlSelectorExtractor{
		newHse(selectTitleNode, extractTitle),
		newHse(selectURLNode, extractURL),
	}
	getFromHTML(body, ses)

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
// if successful. On a failed request, another Sci-Hub mirror is chosen,
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
func doDownloadRequest(a *article.Article) error {
	ctx, cncl := context.WithTimeout(context.Background(), GlobalReqTimeout)
	defer cncl()

	res, err := sendGetRequest(ctx, a.Url.String())
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	out, err := os.Create(a.GenerateFileName() + ".pdf")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if _, err := io.Copy(out, res.Body); err != nil {
		os.Remove(out.Name())
		return fmt.Errorf("%w", err)
	}

	return nil
}

// doCrossrefCitationRequest requests the article citation
func doCrossrefCitationRequest(
	a *article.Article,
	c crossref.ContentType,
) error {
	if err := c.IsValid(); err != nil {
		return fmt.Errorf("%w", err)
	}

	u := &url.URL{
		Scheme: "https",
		Host:   crossref.API,
		Path:   crossref.Works,
	}
	u = u.JoinPath(url.PathEscape(a.Doi), c.Endpoint())

	ctx, cncl := context.WithTimeout(context.Background(), GlobalReqTimeout)
	defer cncl()

	res, err := sendGetRequest(ctx, u.String())
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer res.Body.Close()

	if a.Citation, err = io.ReadAll(res.Body); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
