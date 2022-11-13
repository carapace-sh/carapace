// Package common code
package common

import (
	"strings"

	"github.com/rsteube/carapace/pkg/style"
)

// RawValue represents a completion candidate.
type RawValue struct {
	// All shells
	Value       string
	Display     string
	Description string
	Style       string

	// Shell-specific (ZSH) structuring
	Group           string
	Tag             string
	SuffixRemovable string // Used for = in -f=<val>, or commas in list, etc
}

// A tag can be used by some shells to classify and manipulate their completion values.
type Tag string

const (
	Command Tag = "command"
	Flag    Tag = "option"
	Value   Tag = "value"
)

// TrimmedDescription returns the trimmed description.
func (r RawValue) TrimmedDescription() string {
	maxLength := 80
	description := strings.SplitN(r.Description, "\n", 2)[0]
	description = strings.TrimSpace(description)
	if len([]rune(description)) > maxLength {
		description = string([]rune(description)[:maxLength-3]) + "..."
	}
	return description
}

// RawValues is an alias for []RawValue.
type RawValues []RawValue

// RawValuesFrom creates RawValues from given values.
func RawValuesFrom(values ...string) RawValues {
	rawValues := make([]RawValue, len(values))
	for index, val := range values {
		rawValues[index] = RawValue{Value: val, Display: val, Style: style.Default}
	}
	return rawValues
}

// FilterPrefix filters values with given prefix.
func (r RawValues) FilterPrefix(prefix string) RawValues {
	filtered := make(RawValues, 0)
	for _, r := range r {
		if strings.HasPrefix(r.Value, prefix) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// ByValue alias to filter by value.
type ByValue []RawValue

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// ByDisplay alias to filter by display.
type ByDisplay []RawValue

func (a ByDisplay) Len() int           { return len(a) }
func (a ByDisplay) Less(i, j int) bool { return a[i].Display < a[j].Display }
func (a ByDisplay) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
