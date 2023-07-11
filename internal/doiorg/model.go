package doiorg

// Some™ models used by doi.org's REST API.
//
// For more information about the particular fields see:
//	https://www.doi.org/the-identifier/resources/factsheets/doi-resolution-documentation

const (
	// URL is the doi.org URL
	URL string = "doi.org"
	// API is the path to doi.org's REST API.
	API string = "api/handles"
)

// doi.org's REST API query parameters.
const (
	QueryTypeNone string = "type=none"
	QueryPretty   string = "pretty=true"
)
