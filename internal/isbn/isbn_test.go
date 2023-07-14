package isbn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type isbnTest struct {
	Name   string
	Input  string
	Output bool
}

var isbnTests = []isbnTest{
	{
		Name:   "good-isbn13",
		Input:  "9780136091813",
		Output: true,
	},
	{
		Name:   "good-isbn10",
		Input:  "0136091814",
		Output: true,
	},
	{
		Name:   "good-isbn10-x",
		Input:  "123456789X",
		Output: true,
	},
	{
		Name:   "good-isbn13-dash",
		Input:  "978-0-306-40615-7",
		Output: true,
	},
	{
		Name:   "good-isbn10-dash",
		Input:  "0-306-40615-2",
		Output: true,
	},
	{
		Name:   "good-isbn10-x-dash",
		Input:  "1-234-56789-X",
		Output: true,
	},
	{
		Name:   "bad-isbn13-checksum",
		Input:  "9780136091817",
		Output: false,
	},
	{
		Name:   "bad-isbn10-checksum",
		Input:  "0136091812",
		Output: false,
	},
	{
		Name:   "bad-isbn10-x-checksum",
		Input:  "013609181Y",
		Output: false,
	},
	{
		Name:   "bad-isbn13-string",
		Input:  "aaaaaaaaaaaa",
		Output: false,
	},
	{
		Name:   "bad-isbn13-rune",
		Input:  "9780136a91813",
		Output: false,
	},
	{
		Name:   "bad-isbn10-string",
		Input:  "aaaaaaaaaa",
		Output: false,
	},
	{
		Name:   "bad-isbn10-rune",
		Input:  "0136091a14",
		Output: false,
	},
	{
		Name:   "bad-cake",
		Input:  "cake",
		Output: false,
	},
}

func TestValid(t *testing.T) {
	for _, tt := range isbnTests {
		t.Run(tt.Name, func(t *testing.T) {
			out := Valid(tt.Input)
			assert.Equal(t, tt.Output, out)
		})
	}
}
