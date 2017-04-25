package wns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingFiles(t *testing.T) {
	files := []string{"a", "b", "c", "d"}
	tests := []struct {
		pos      string
		expected []string
	}{
		{
			pos:      "unknown",
			expected: []string{"a", "b", "c", "d"},
		},
		{
			pos:      "a",
			expected: []string{"b", "c", "d"},
		},
		{
			pos:      "b",
			expected: []string{"c", "d"},
		},
		{
			pos:      "c",
			expected: []string{"d"},
		},
		{
			pos:      "d",
			expected: []string{},
		},
	}

	for _, test := range tests {
		actual := missingFiles(test.pos, files)
		assert.Equal(t, test.expected, actual)
	}
}
