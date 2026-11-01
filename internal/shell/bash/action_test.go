package bash

import (
	"testing"

	"github.com/carapace-sh/carapace/internal/common"
)

func TestActionRawValuesPartialCompletion(t *testing.T) {
	tests := []struct {
		name         string
		currentWord  string
		values       common.RawValues
		wordbreakPre string
		expected     string
	}{
		{
			name:        "collapse when common value prefix extends current word",
			currentWord: "~steam/f",
			values: common.RawValuesFrom(
				"~steam/file1", "~steam/file2",
			),
			expected: "true\001~steam/file",
		},
		{
			name:        "drop candidates not extending the current word",
			currentWord: "~steam/D",
			values: common.RawValuesFrom(
				"~steam/Desktop/", "~steam/Documents/", "~steam/Downloads/", "~steam/d1/", "~steam/dev/",
			),
			expected: "true\001~steam/Desktop/\n~steam/Documents/\n~steam/Downloads/",
		},
		{
			name:        "keep candidates when none extend the current word",
			currentWord: "~",
			values: common.RawValuesFrom(
				"a", "b",
			),
			expected: "false\001a\nb",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wordbreakPrefix, compType = "", "" // reset package state possibly changed by Patch
			values := make(common.RawValues, len(test.values))
			copy(values, test.values)
			if actual := ActionRawValues(test.currentWord, common.Meta{}, values); actual != test.expected {
				t.Errorf("expected [%v], got [%v]", test.expected, actual)
			}
		})
	}
}
