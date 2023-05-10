package metainfo

import "fmt"

var (
	Project       string = "_"
	Version       string = "_"
	Url           string = "_"
	Maintainer    string = "_"
	HTTPUserAgent string = "_"
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
