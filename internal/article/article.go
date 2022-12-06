package article

import (
	"net/url"
	"strings"
	"unicode"
)

// The Article holds data needed to download and write the article to disc.
// The DOI is used to fetch the title and PDF URL from Sci-Hub.
type Article struct {
	Doi      string
	Title    string
	Url      *url.URL // PDF download link
	Citation []byte

	// fileNameGenerator generates a file name for the article
	fileNameGenerator func(Article) string

	// fileName is the (local) name of the article file
	fileName string
}

// Reset resets all article data.
func (a *Article) Reset() {
	*a = Article{}
}

// GeneratorFunc assigns a new generator
func (a *Article) GeneratorFunc(fileNameGenerator func(Article) string) {
	a.fileNameGenerator = fileNameGenerator
}

// GenerateFileName generates and caches the file name of the article.
func (a *Article) GenerateFileName() string {
	if a.fileNameGenerator == nil {
		panic("file name generator not set")
	}

	if len(a.fileName) == 0 {
		a.fileName = a.fileNameGenerator(*a)
	}
	return a.fileName
}

// SnakeCaseGenerator is generates a snake-case file name from the Article
// title. All punctuation, spaces and control codes are replaced by '_'s, which
// are squeezed.
// No extension is added to the file name.
func SnakeCaseGenerator(a Article) string {
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

	return b.String()
}
