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
		"%v/%v (%v)",
		Project,
		Version,
		Url,
	)
}
