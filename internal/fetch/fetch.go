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
	"github.com/Milover/fetchpaper/internal/doiorg"
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

	// CiteFormat is the citation output format.
	CiteFormat = crossref.BibTeX

	// CiteFileName is the citation filename base (w/o extension).
	CiteFileName = "citations"

	// CiteAppend controls whether the citation file will be
	// appended to or overwritten.
	CiteAppend = false

	// CiteSeparate controls whether the citations will be written to
	// separate files
	CiteSeparate = false

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

// FetchMode signifies the mode in which the Fetch function will run.
type FetchMode int

// Available modes of the Fetch function.
const (
	DefaultMode FetchMode = iota
	SourceMode
	CiteMode
)

// Fetch downloads articles from Sci-Hub and/or citations from Crossref,
// from a list of supplied DOIs.
func Fetch(mode FetchMode, dois []string) error {
	if len(dois) == 0 {
		return nil
	}
	dois = CheckDOIs(dois)
	articles := make([]article.Article, 0, len(dois))
	for i := range dois {
		articles = append(articles, article.Article{Doi: dois[i]})
		a := &articles[i]
		// FIXME: the generator should be configurable
		a.GeneratorFunc(article.SnakeCaseGenerator)
	}
	// fetch article metadata
	// FIXME: this could be done better:
	// we're waiting for all metadata requests to finish before requesting
	// citations, while in reality, both citation and download link fetching
	// can be done concurrently with the meta data request.
	// Hence, if we had a better way of naming the individual citation files
	// (if writing to individual citation files, we need the article names),
	// or a better way of synchronizing between different requests,
	// we could make this faster.
	//
	// Maybe instead of batching a set of requests for all articles,
	// we could batch all requests for each article, this would make
	// correctly ordering/parallelizing requests easier.
	if err := fetchMetas(articles); err != nil {
		log.Println("errors occurred during metadata fetch")
	}

	// fetch articles and citations
	g := new(errgroup.Group)
	if mode != CiteMode {
		g.Go(func() error {
			return fetchArticles(articles)
		})
	}
	if mode != SourceMode {
		g.Go(func() error {
			return fetchCitations(articles)
		})
	}
	return g.Wait()
}

// CheckDOIs is a function which takes a list of DOIs and returns only the
// (valid) ones registered at doi.org.
func CheckDOIs(dois []string) []string {
	ch := make(chan string, len(dois))
	var wg sync.WaitGroup
	wg.Add(len(dois))
	for i := range dois {
		d := dois[i]
		go func() {
			defer wg.Done()
			if err := logErr(d, CheckDOI(d)); err == nil {
				ch <- d
			}
		}()
	}
	wg.Wait()
	close(ch)

	// select valid DOIs
	valid := make([]string, 0, len(ch))
	for doi := range ch {
		valid = append(valid, doi)
	}
	return valid
}

// CheckDOI checks if a doi is valid (registered) by querying doi.org
// for a good response.
func CheckDOI(doi string) error {
	u := &url.URL{
		Scheme: "https",
		Host:   doiorg.URL,
		Path:   doiorg.API,
	}
	// set the query
	query := url.Values{}
	query.Add(doiorg.QueryKeyType, doiorg.QueryValType)
	u.RawQuery = query.Encode()

	u = u.JoinPath(url.PathEscape(doi))

	ctx, cncl := context.WithTimeout(context.Background(), GlobalReqTimeout)
	defer cncl()

	_, err := sendGetRequest(ctx, u.String())
	return err
}

// WARNING: assumes that articles have DOIs set.
func fetchMetas(articles []article.Article) error {
	g := new(errgroup.Group)

	for i := range articles {
		a := &articles[i]

		// GET article info
		g.Go(func() error {
			// GET metadata from Crossref, and set the article title
			meta, err := reqCrossrefMeta(a)
			if err != nil {
				return logErr(a.Doi, err)
			}
			// WARNING: is it ok to assume that the first item is the one we want?
			a.Title = meta.Message.Title[0]
			if len(a.Title) == 0 {
				a.Title = a.Doi
				log.Printf("%v: could not set title", a.Doi)
			}
			return nil
		})
	}
	return g.Wait()
}

// WARNING: assumes that articles have DOIs and generator functions set.
func fetchArticles(articles []article.Article) error {
	g := new(errgroup.Group)

	for i := range articles {
		a := &articles[i]

		// GET article info
		g.Go(func() error {
			// GET PDF download link from Sci-Hub
			if err := logErr(a.Doi, reqSciHubInfo(a)); err != nil {
				return err
			}
			// GET article PDF
			return logErr(a.Doi, reqDownload(a))
		})
	}
	return g.Wait()
}

// WARNING: assumes that articles have DOIs set.
func fetchCitations(articles []article.Article) error {
	g := new(errgroup.Group)

	for i := range articles {
		a := &articles[i]

		g.Go(func() error {
			return logErr(a.Doi, reqCrossrefCitation(a))
		})
	}
	err := g.Wait()
	err = errors.Join(err, writeCitations(articles))
	return err
}

// logErr is a helper function which logs err, if it is not nil, and an
// associated DOI, and returns err.
func logErr(doi string, err error) error {
	if err != nil {
		log.Printf("%v: %v", doi, err)
	}
	return err
}

func openCiteFile(filename string) (*os.File, error) {
	flag := os.O_RDWR | os.O_CREATE
	if CiteAppend {
		flag |= os.O_APPEND
	} else {
		flag |= os.O_TRUNC
	}
	return os.OpenFile(filename, flag, 0666)
}

// writeCitations writes all citations to a file.
func writeCitations(articles []article.Article) error {
	var out *os.File
	var err error
	for i, a := range articles {
		if len(a.Citation) == 0 {
			continue
		}
		// precaution and output readability
		if a.Citation[len(a.Citation)-1] != '\n' {
			a.Citation = append(a.Citation, '\n')
		}

		if CiteSeparate {
			out, err = openCiteFile(a.GenerateFileName() + CiteFormat.Extension())
		} else if i == 0 { // open/close only once
			out, err = openCiteFile(CiteFileName + CiteFormat.Extension())
			defer out.Close()
		}
		if err != nil {
			return err
		}

		if _, err := out.Write(a.Citation); err != nil {
			return err
		}
		if CiteSeparate {
			out.Close()
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
		return res, fmt.Errorf("%s: %s", res.Request.URL, res.Status)
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
	u = u.JoinPath(url.PathEscape(a.Doi), CiteFormat.Endpoint())

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
