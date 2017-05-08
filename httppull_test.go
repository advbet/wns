package wns

import (
	"encoding/xml"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckForErrs(t *testing.T) {
	tests := []struct {
		msg      string
		xml      string
		expected error
	}{
		{
			msg: "test error tag",
			xml: `<?xml version="1.0" encoding="ISO-8859-1"?><error>Too frequent download. (From IP: 127.0.0.1, 3 seconds ago)</error>`,
			expected: &APIError{
				Type: ErrTypeTooFrequent,
				Err:  "Too frequent download. (From IP: 127.0.0.1, 3 seconds ago)",
			},
		},
		{
			msg: "test error-message tag",
			xml: `<?xml version="1.0" encoding="ISO-8859-1"?><error-message>There are no files ready for transfer at the moment.</error-message>`,
			expected: &APIError{
				Type: ErrTypeNoNew,
				Err:  "There are no files ready for transfer at the moment.",
			},
		},
		{
			msg:      "test random tag",
			xml:      `<?xml version="1.0" encoding="ISO-8859-1"?><random>Some random content</random>`,
			expected: nil,
		},
		{
			msg:      "test empty xml",
			xml:      "",
			expected: io.EOF,
		},
		{
			msg:      "test malformed xml",
			xml:      `<?xml`,
			expected: &xml.SyntaxError{Msg: "unexpected EOF", Line: 1},
		},
	}

	for _, test := range tests {
		actual := checkForErr(strings.NewReader(test.xml))
		assert.Equal(t, test.expected, actual, test.msg)
	}
}
