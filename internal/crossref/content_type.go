package crossref

import "fmt"

// ContentType represents a content return type supported by Crossref's API.
type ContentType int

// Citation content return types supported by Crossref's API.
// For more information see: https://citation.crosscite.org/docs.html
const (
	BibTeX ContentType = iota
	CiteprocJSON
	RDFXML
	RDFTurtle
	RIS
	SchemaorgJSON
	TextCitation
	UnixrefXML
	UnixsdXML
)

// Endpoint returns the Crossref API content type endpoint path.
// This is appended to the '/works/{doi}' endpoint, which, when requested,
// returns the citation formatted in the requested type.
//
// For more information see: https://citation.crosscite.org/docs.html
func (c ContentType) Endpoint() string {
	return endpoints[c]
}

// FileExtension returns the file extension for the ContentType.
func (c ContentType) FileExtension() string {
	return extensions[c]
}

// IsValid checks if the content type value is valid.
func (c ContentType) IsValid() error {
	if (int(c) < 0) || (int(c) >= len(names)) {
		return fmt.Errorf("invalid content type")
	}
	return nil
}

// Set sets the value of the content type based on the provided
// content type name.
func (c *ContentType) Set(name string) error {
	for i, n := range names {
		if name == n {
			*c = ContentType(i)
			return nil
		}
	}

	return fmt.Errorf("unknown content type: %v", name)
}

// String returns the ContentType (name) as a user-friendly string.
func (c ContentType) String() string {
	return names[c]
}

// Type returns the type used by ContentType.Set.
func (c ContentType) Type() string {
	return "string"
}

var (
	// names are the user-friendly ContentType names.
	names = [...]string{
		"bibtex",
		"citeprocjson",
		"rdfxml",
		"rdfturtle",
		"ris",
		"schemaorgjson",
		"textcitation",
		"unixrefxml",
		"unixsdxml",
	}
	// extensions are the ContentType file extensions.
	extensions = [...]string{
		".bib",
		".json",
		".rdf",
		".ttl",
		".ris",
		".jsonld",
		".txt",
		".xml",
		".xml",
	}
	// endpoints are Crossref API endpoint paths for the different ContentTypes.
	endpoints = [...]string{
		"transform/application/x-bibtex",
		"transform/application/vnd.citationstyles.csl+json",
		"transform/application/rdf+xml",
		"transform/text/turtle",
		"transform/application/x-research-info-systems",
		"transform/application/vnd.schemaorg.ld+json",
		"transform/text/x-bibliography",
		"transform/application/vnd.crossref.unixref+xml",
		"transform/application/vnd.crossref.unixsd+xml",
	}
)
