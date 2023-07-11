package fetch

import (
	"context"
	"encoding/json"
	"errors"
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

	// CitationFileName is the citation filename base (w/o extension).
	CitationFileName = "citations"

	// NoUserAgent controls weather to omit the User-Agent header in
	// HTTP requests.
	NoUserAgent = false

	// A list of Sci-Hub mirrors.
	mirrors = []string{
		"sci-hub.se",
		"sci-hub.st",
		"sci-hub.ru",
	}
)

// Fetch downloads articles from Sci-Hub and/or citations from Crossref,
// from a list of supplied DOIs.
func Fetch(dois []string) error {
	if len(dois) == 0 {
		return nil
	}
	articles := make([]article.Article, 0, len(dois))

	for i := range dois {
		articles = append(articles, article.Article{Doi: dois[i]})
		a := &articles[i]
		// FIXME: the generator should be configurable
		a.GeneratorFunc(article.SnakeCaseGenerator)
	}

	g := new(errgroup.Group)
	g.Go(func() error {
		return FetchArticles(articles)
	})
	g.Go(func() error {
		return FetchCitations(articles)
	})
	return g.Wait()
}

// FIXME: make this simpler, the synchronization is unnecessarily complex
func FetchArticles(articles []article.Article) error {
	ch := make(chan *article.Article, len(articles))
	g := new(errgroup.Group)
	var wg sync.WaitGroup // so we know when to close the channel
	wg.Add(len(articles))

	for i := range articles {
		a := &articles[i]

		// GET article info
		g.Go(func() error {
			defer wg.Done()
			eg := new(errgroup.Group)

			// GET PDF download link from Sci-Hub
			eg.Go(func() error {
				if err := reqSciHubInfo(a); err != nil {
					log.Printf("%v: %v", a.Doi, err)
					return err
				}
				return nil
			})
			// GET metadata from Crossref, and set the article title
			eg.Go(func() error {
				meta, err := reqCrossrefMeta(a)
				if err != nil {
					log.Printf("%v: %v", a.Doi, err)
					return err
				}
				// WARNING: is it ok to assume that the first item is the one we want?
				a.Title = meta.Message.Title[0]
				if len(a.Title) == 0 {
					a.Title = a.Doi
					log.Printf("%v: could not set title", a.Doi)
					return fmt.Errorf("%v: could not set title", a.Doi)
				}
				return nil
			})
			err := eg.Wait()
			ch <- a
			return err
		})
		// GET article PDF
		// FIXME: move this into get article info at the end, and remove
		// the channel and the WaitGroup
		g.Go(func() error {
			if a, ok := <-ch; ok {
				if err := reqDownload(a); err != nil {
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

	return g.Wait()
}

func FetchCitations(articles []article.Article) error {
	g := new(errgroup.Group)

	for i := range articles {
		a := &articles[i]

		g.Go(func() error {
			if err := reqCrossrefCitation(a); err != nil {
				log.Printf("%v: %v", a.Doi, err)
				return err
			}
			return nil
		})
	}
	err := g.Wait()
	err = errors.Join(err, writeCitations(articles))
	return err
}

// writeCitations writes all citations to a file.
func writeCitations(articles []article.Article) error {
	out, err := os.Create(CitationFileName + CitationFormat.FileExtension())
	if err != nil {
		panic(err)
	}
	defer out.Close()
	for _, a := range articles {
		if len(a.Citation) == 0 {
			continue
		}
		if a.Citation[len(a.Citation)-1] != '\n' {
			a.Citation = append(a.Citation, '\n')
		}
		if _, err := out.Write(a.Citation); err != nil {
			return err
		}
	}
	return nil
}

// sendGetRequest sends a GET request to the specified URL.
// An error is returned if a valid response cannot be obtained.
func sendGetRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if !NoUserAgent {
		req.Header.Set("User-Agent", metainfo.HTTPUserAgent)
	}

	// sendGetRequest can be called from multiple threads
	GlobalRateLimiter.Take()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, err
	}
	if res.StatusCode > 399 {
		return res, fmt.Errorf("%s", res.Status)
	}
	//fmt.Println(res.Header.Get("x-api-pool"))

	return res, nil
}

// reqSciHubMirrorInfo requests article info from a Sci-Hub mirror
// and parses the article title and download URL from the response HTML.
func reqSciHubMirrorInfo(a *article.Article, mirror string) error {
	u := &url.URL{
		Scheme: "https",
		Host:   mirror,
		Path:   a.Doi,
	}

	ctx, cncl := context.WithTimeout(context.Background(), GlobalReqTimeout)
	defer cncl()

	res, err := sendGetRequest(ctx, u.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := html.Parse(res.Body)
	if err != nil {
		return err
	}

	// extract url from parsed HTML
	hse := newHse(selectURLNode, extractURL)
	getFromHTML(body, hse)
	if hse.data.Len() == 0 {
		return fmt.Errorf("could not extract article URL from HTML")
	}

	// set and check title/body, clean up if necessary
	//	a.Title = hse.data.String()
	//	if len(a.Title) == 0 {
	//		a.Title = a.Doi
	//		log.Printf("%v: could not extract title", a.Doi)
	//	}

	a.Url, err = url.Parse(hse.data.String())
	if err != nil {
		return err
	}
	a.Url.Scheme = "https"
	if len(a.Url.Host) == 0 {
		a.Url.Host = res.Request.URL.Host
	}

	return nil
}

// reqSciHubInfo requests article info from Sci-Hub, and updates the article
// if successful. On a failed request, another Sci-Hub mirror is chosen,
// until all mirrors have been exhausted.
func reqSciHubInfo(a *article.Article) error {
	for i, m := range mirrors {
		if err := reqSciHubMirrorInfo(a, m); err != nil {
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
func reqDownload(a *article.Article) error {
	ctx, cncl := context.WithTimeout(context.Background(), GlobalReqTimeout)
	defer cncl()

	if a.Url == nil {
		return fmt.Errorf("could not download article, URL empty")
	}
	res, err := sendGetRequest(ctx, a.Url.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	out, err := os.Create(a.GenerateFileName() + ".pdf")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if _, err := io.Copy(out, res.Body); err != nil {
		os.Remove(out.Name())
		return err
	}

	return nil
}

// reqCrossrefCitation requests the article citation from Crossref
func reqCrossrefCitation(a *article.Article) error {
	u := &url.URL{
		Scheme: "https",
		Host:   crossref.API,
		Path:   crossref.Works,
	}
	u = u.JoinPath(url.PathEscape(a.Doi), CitationFormat.Endpoint())

	ctx, cncl := context.WithTimeout(context.Background(), GlobalReqTimeout)
	defer cncl()

	res, err := sendGetRequest(ctx, u.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if a.Citation, err = io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

// reqCrossrefMeta requests the article metadata from Crossref
func reqCrossrefMeta(a *article.Article) (crossref.WorkMessage, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   crossref.API,
		Path:   crossref.Works,
	}
	u = u.JoinPath(url.PathEscape(a.Doi))

	ctx, cncl := context.WithTimeout(context.Background(), GlobalReqTimeout)
	defer cncl()

	res, err := sendGetRequest(ctx, u.String())
	if err != nil {
		return crossref.WorkMessage{}, err
	}
	defer res.Body.Close()

	var msg crossref.WorkMessage
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return crossref.WorkMessage{}, err
	}
	err = json.Unmarshal(b, &msg)
	if err != nil {
		return crossref.WorkMessage{}, err
	}

	return msg, nil
}
