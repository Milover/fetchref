package crossref

import "encoding/json"

// Some™ models used by Crossref's REST API.
//
// For more information about the particular fields see:
//	https://api.crossref.org/swagger-ui/index.html
//	https://github.com/CrossRef/rest-api-doc/blob/master/api_format.md
//
// TODO: All of this should probably be autogenerated somehow,
// probably from Crossref's REST API documentation page or
// an example-JSON-response.

const (
	// API is the URL of Crossref's REST API.
	API string = "api.crossref.org"
)

// Endpoint is a type alias for the endpoint string.
type Endpoint = string

// Crossref's REST API endpoints.
const (
	Funders  Endpoint = "funders"
	Journals Endpoint = "journals"
	Licences Endpoint = "licenses"
	Members  Endpoint = "members"
	Prefixes Endpoint = "prefixes"
	Types    Endpoint = "types"
	Works    Endpoint = "works"
)

// Affiliation holds the name of an affiliated institution.
type Affiliation struct {
	Name string `json:"name"`
}

// Author holds contributor-related data (author, editor, reviewer etc.).
type Author struct {
	ORCID              string        `json:"ORCID"`
	Suffix             string        `json:"suffix"`
	Given              string        `json:"given"`
	Family             string        `json:"family"`
	Affiliation        []Affiliation `json:"affiliation"`
	Name               string        `json:"name"`
	AuthenticatedORCID bool          `json:"authenticated-orcid"`
	Prefix             string        `json:"prefix"`
	Sequence           string        `json:"sequence"`
}

// Date holds the full information about a date in several formats.
type Date struct {
	// DateParts is an ordered array of year, month and day.
	DateParts [][]int `json:"date-parts"`
	// DateTime is the ISO 8601 formatted date time.
	DateTime string `json:"date-time"`
	// TimeStamp is the number of seconds since the UNIX epoch.
	Timestamp int `json:"timestamp"`
}

// DateParts holds partial information about a date.
type DateParts struct {
	// DateParts is an ordered array of year, month and day.
	DateParts [][]int `json:"date-parts"`
}

// Reference holds basic data about a work which references another work.
type Reference struct {
	ISSN               string `json:"issn"`
	StandardsBody      string `json:"standards-body"`
	Issue              string `json:"issue"`
	Key                string `json:"key"`
	SeriesTitle        string `json:"series-title"`
	ISBNType           string `json:"isbn-type"`
	DOIAssertedBy      string `json:"doi-asserted-by"`
	FirstPage          string `json:"first-page"`
	ISBN               string `json:"isbn"`
	DOI                string `json:"doi"`
	Component          string `json:"component"`
	ArticleTitle       string `json:"article-title"`
	VolumeTitle        string `json:"volume-title"`
	Volume             string `json:"volume"`
	Author             string `json:"author"`
	StandardDesignator string `json:"standard-designator"`
	Year               string `json:"year"`
	Unstructured       string `json:"unstructured"`
	Edition            string `json:"edition"`
	JournalTitle       string `json:"journal-title"`
	ISSNType           string `json:"issn-type"`
}

// Work holds all metadata about a particular work (article, book, dataset etc.).
type Work struct {
	Institution         WorkInstitution     `json:"institution"`
	Indexed             Date                `json:"indexed"`
	Posted              DateParts           `json:"posted"`
	PublisherLocation   string              `json:"publisher-location"`
	UpdateTo            []WorkUpdate        `json:"update-to"`
	StandardsBody       []WorkStandardsBody `json:"standards-body"`
	EditionNumber       string              `json:"edition-number"`
	GroupTitle          []string            `json:"group-title"`
	ReferenceCount      int                 `json:"reference-count"`
	Publisher           string              `json:"publisher"`
	Issue               string              `json:"issue"`
	ISBNType            []WorkISSNType      `json:"isbn-type"`
	License             []WorkLicense       `json:"license"`
	Founder             []WorkFunder        `json:"founder"`
	ContentDomain       WorkDomain          `json:"content-domain"`
	Chair               []Author            `json:"chair"`
	ShortContainerTitle string              `json:"shortcontainer-title"`
	Accepted            DateParts           `json:"accepted"`
	ContentUpdated      DateParts           `json:"content-updated"`
	PublishedPrint      DateParts           `json:"published-print"`
	Abstract            string              `json:"abstract"`
	DOI                 string              `json:"DOI"`
	Type                string              `json:"type"`
	Created             Date                `json:"created"`
	Approved            DateParts           `json:"approved"`
	Page                string              `json:"page"`
	UpdatePolicy        string              `json:"update-policy"`
	Source              string              `json:"source"`
	IsReferencedByCount int                 `json:"is-referenced-by-count"`
	Title               []string            `json:"title"`
	Prefix              string              `json:"prefix"`
	Volume              string              `json:"volume"`
	ClinicalTrialNumber []WorkClinicalTrial `json:"clinical-trial-number"`
	Author              []Author            `json:"author"`
	Member              string              `json:"member"`
	ContentCreated      DateParts           `json:"content-created"`
	PublishedOnline     DateParts           `json:"published-online"`
	Reference           []Reference         `json:"reference"`
	ContainerTitle      []string            `json:"container-title"`
	Review              WorkReview          `json:"review"`
	OriginalTitle       []string            `json:"original-title"`
	Language            string              `json:"language"`
	Link                []WorkLink          `json:"link"`
	Deposited           Date                `json:"deposited"`
	Score               int                 `json:"score"`
	Degree              string              `json:"degree"`
	Subtitle            []string            `json:"subtitle"`
	Translator          []Author            `json:"translator"`
	FreeToRead          WorkFreeToRead      `json:"free-to-read"`
	Editor              []Author            `json:"editor"`
	ComponentNumber     string              `json:"component-number"`
	ShortTitle          []string            `json:"short-title"`
	Issued              DateParts           `json:"issued"`
	ISBN                []string            `json:"ISBN"`
	ReferencesCount     int                 `json:"references-count"`
	PartNumber          string              `json:"part-number"`
	JournalIssue        WorkJournalIssue    `json:"journal-issue"`
	AlternativeID       []string            `json:"alternative-id"`
	URL                 string              `json:"URL"`
	Archive             []string            `json:"archive"`
	Relation            json.RawMessage     `json:"relation"` // WARNING: unused
	ISSN                []string            `json:"ISSN"`
	ISSNType            []WorkISSNType      `json:"issn-type"`
	Subject             []string            `json:"subject"`
	PublishedOther      DateParts           `json:"published-other"`
	Published           DateParts           `json:"published"`
	Assertion           []WorkAssertion     `json:"assertion"`
	Subtype             string              `json:"subtype"`
	ArticleNumber       string              `json:"article-number"`
}

// WorkAssertion holds information about various assertions related to the work.
type WorkAssertion struct {
	Group       WorksMessageMessageItemsAssertionGroup       `json:"group"`
	Explanation WorksMessageMessageItemsAssertionExplanation `json:"explanation"`
	Name        string                                       `json:"name"`
	Value       string                                       `json:"value"`
	URL         string                                       `json:"URL"`
	Order       int                                          `json:"order"`
	Label       string                                       `json:"label"` // WARNING: undocumented
}

// WorkClinicalTrial holds data about the clinical trial related to the work.
type WorkClinicalTrial struct {
	ClinicalTrialNumber string `json:"clinical-trial-number"`
	Registry            string `json:"registry"` // DOI of the registry
	Type                string `json:"type"`
}

// WorkDomain holds information about domains that support Crossmark for the work.
type WorkDomain struct {
	Domain               []string `json:"domain"`
	CrossmarkRestriction bool     `json:"crossmark-restriction"`
}

// WorkFreeToRead (presumably) holds information about open-access availability
// of the work.
// NOTE: This field is possibly unused.
type WorkFreeToRead struct {
	StartDate DateParts `json:"start-date"`
	EndDate   DateParts `json:"end-date"`
}

// WorkFunder holds information about a funding organization of the work.
type WorkFunder struct {
	Name          string   `json:"name"`
	DOI           string   `json:"DOI"` // Open Funder Registry DOI
	DOIAssertedBy string   `json:"doi-asserted-by"`
	Award         []string `json:"award"`
}

// WorkInstitution holds information about the institution related to the work.
type WorkInstitution struct {
	Name       string   `json:"name"`
	Place      []string `json:"place"`
	Department []string `json:"department"`
	Acronym    []string `json:"acronym"`
}

// WorkISSNType holds ISSN information related to the work.
type WorkISSNType struct {
	// Type is the ISSN type: 'eissn', 'pissn' or 'lissn'
	Type  string `json:"type"`
	Value string `json:"value"`
}

// WorkJournalIssue is the journal issue in which the work was published.
type WorkJournalIssue struct {
	Issue string `json:"issue"`
}

// WorkLicense holds data about the license of the work.
type WorkLicense struct {
	// URL is a link to the web page describing the license.
	URL string `json:"URL"`
	// Start is the date on which this license begins to take effect
	Start Date `json:"start"`
	// DelayInDays is the number of days between the publication date of the
	// work and the start date of this license.
	DelayInDays int `json:"delay-in-days"`
	// ContentVersion is one of 'vor' (version of record), 'am' (accepted
	// manuscript), 'tdm' (text and data mining) or 'unspecified'.
	ContentVersion string `json:"content-version"`
}

// WorkLink holds information about the URL of the full-text location of the work.
type WorkLink struct {
	// URL is the direct link to a full-text download location.
	URL string `json:"URL"`
	// ContentType is the content type (or MIME type) of the full-text object.
	ContentType string `json:"content-type"`
	// ContentVersion is one of 'vor' (version of record), 'am' (accepted
	// manuscript) or 'unspecified'.
	ContentVersion string `json:"content-version"`
	// IntendedApplication is one of 'text-mining', 'similarity-checking'
	// or 'unspecified'.
	IntendedApplication string `json:"intended-application"`
}

// WorkMessage is the return type of the '/works' and '/works/{doi}' endpoints.
type WorkMessage struct {
	Status         string `json:"status"`
	MessageType    string `json:"message-type"`
	MessageVersion string `json:"message-version"`
	Message        Work   `json:"message"`
}

// WorkReview holds peer review metadata.
type WorkReview struct {
	// Type is one of 'major-revision', 'minor-revision', 'reject',
	// 'reject-with-resubmit' or 'accept'.
	Type string `json:"type"`
	// Stage is one of 'pre-publication' or 'post-publication'.
	Stage string `json:"stage"`
	// Recommendation is one of 'referee-report', 'editor-report',
	// 'author-comment', 'community-comment' or 'aggregate'.
	Recommendation             string `json:"recommendation"`
	RunningNumber              string `json:"running-number"`
	RevisionRound              string `json:"revision-round"`
	Language                   string `json:"language"`
	CompetingInterestStatement string `json:"competing-interest-statement"`
}

// WorkStandardsBody holds information about the standards body
// related to the work.
type WorkStandardsBody struct {
	Name    string   `json:"name"`
	Acronym []string `json:"acronym"`
}

// WorkUpdate holds information about updates of the work.
type WorkUpdate struct {
	// Label is a display-friendly label for the update type.
	Label string `json:"label"`
	// DOI is the DOI of the updated work.
	DOI string `json:"DOI"`
	// Type is the type of update, e.g., 'retraction' or 'correction'.
	Type string `json:"type"`
	// Updated is the date on which the update was published.
	Updated Date `json:"updated"`
}

// WorksMessageMessageItemsAssertionGroup holds data about the group to which
// the assertion related to the work belongs.
type WorksMessageMessageItemsAssertionGroup struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

// WorksMessageMessageItemsAssertionExplanation holds an explanation of the
// assertion related to the work.
type WorksMessageMessageItemsAssertionExplanation struct {
	URL string `json:"URL"`
}
