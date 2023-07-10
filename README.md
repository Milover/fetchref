# fetchpaper

A simple, ~~hacky~~ command line utility for fetching article PDFs from
Sci-Hub and formatted citations from Crossref from supplied DOIs.

## TODO

- [x] Sci-Hub returns HTTP 200 even if the article is unavailable, so sometimes
      we download garbage and we need to check for this somehow
    - if we can't extract the article URL, report and don't attempt download
- [ ] add more functionality for managing citations
	- [ ] add flag for enabling new citations to be appended to an
		existing `citations` file
	- [ ] add flag for enabling writing citations to individual files
		- *some formats only make sense as individual files, e.g. RIS*
- [ ] add book-fetch by ISBN from LibGen
	- fetch all available versions of the book from:
	  `libgen.is/json.php?isbn=<isbn>&fields='*'`
	- pick the best match by available file format and size, language, No. pages,
		and year
	- download from Cloudflare:
		`https://cloudflare-ipfs.com/ipfs/<ipfs_cid>?filename=<md5>.<extension>`
	- also check available LibGen mirrors
	- rough JSON endpoint [API reference][libgen_api]
- [ ] separate the citation fetching and resource download to separate subcommands
- [ ] rename the project to something more suitable like `getref`, `fetchref` or
	something
- [ ] add better logging ([logrus][logrus]) and instrumentation
- [ ] add unit/integration tests


[logrus]: https://pkg.go.dev/github.com/sirupsen/logrus
[libgen_api]: http://faq.fyicenter.com/1231_What_Is_Library_Genesis_API.html
