package metainfo

import "fmt"

var (
	Project       string
	Version       string
	Url           string
	Maintainer    string
	HTTPUserAgent string
)

func init() {
	HTTPUserAgent = fmt.Sprintf(
		"%v/%v (%v; mailto:%v)",
		Project,
		Version,
		Url,
		Maintainer,
	)
}
