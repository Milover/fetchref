package libgen

// Some™ models used by libgen's REST API.
//
// For more information about the particular fields see:
//	http://faq.fyicenter.com/1231_What_Is_Library_Genesis_API.html
//	https://libgen.lc/json.php

const (
	// Mirrors is a list of libgen mirror URLs
	Mirrors = []string{
		"libgen.is",
		"libgen.rs",
		"libgen.st",
	}

	// API is the path to libgen's JSON API endpoint.
	API string = "json.php"
)

// doi.org's REST API query parameters.
const (
	QueryKeyFields string = "fields"
	QueryValFields string = "*"
	QueryKeyISBN   string = "isbn"
)

// Message holds all data from libgen's JSON API response.
type Message struct {
	Id               string `json:"id"`
	Title            string `json:"title"`
	VolumeInfo       string `json:"volumeinfo"`
	Series           string `json:"series"`
	Periodical       string `json:"periodical"`
	Author           string `json:"author"`
	Year             string `json:"year"`
	Edition          string `json:"edition"`
	Publisher        string `json:"publisher"`
	City             string `json:"city"`
	Pages            string `json:"pages"`
	Language         string `json:"language"`
	Topic            string `json:"topic"`
	Library          string `json:"library"`
	Issue            string `json:"issue"`
	Identifier       string `json:"identifier"`
	ISSN             string `json:"issn"`
	ASIN             string `json:"asin"`
	UDC              string `json:"udc"`
	LBC              string `json:"lbc"`
	DDC              string `json:"ddc"`
	LCC              string `json:"lcc"`
	DOI              string `json:"doi"`
	GooglebookID     string `json:"googlebookid"`
	OpenLibraryID    string `json:"openlibraryid"`
	Commentary       string `json:"commentary"`
	DPI              string `json:"dpi"`
	Color            string `json:"color"`
	Cleaned          string `json:"cleaned"`
	Orientation      string `json:"orientation"`
	Paginated        string `json:"paginated"`
	Scanned          string `json:"scanned"`
	Bookmarked       string `json:"bookmarked"`
	Searchable       string `json:"searchable"`
	FileSize         string `json:"filesize"`
	Extension        string `json:"extension"`
	MD5              string `json:"md5"`
	Generic          string `json:"generic"`
	Visible          string `json:"visible"`
	Locator          string `json:"locator"`
	Local            string `json:"local"`
	TimeAdded        string `json:"timeadded"`
	TimeLastModified string `json:"timelastmodified"`
	CoverURL         string `json:"coverurl"`
	IdentifierWodash string `json:"identifierwodash"`
	Tags             string `json:"tags"`
	PagesInFile      string `json:"pagesinfile"`
	Descr            string `json:"descr"`
	TOC              string `json:"toc"`
	SHA1             string `json:"sha1"`
	SHA256           string `json:"sha256"`
	CRC32            string `json:"crc32"`
	EdonKey          string `json:"edonkey"`
	AICH             string `json:"aich"`
	TTH              string `json:"tth"`
	IpfsCID          string `json:"ipfs_cid"`
	BTIH             string `json:"btih"`
	Torrent          string `json:"torrent"`
}
