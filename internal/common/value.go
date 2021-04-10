package common

import "strings"

// RawValue represents a completion candidate
type RawValue struct {
	Value       string
	Display     string
	Description string
}

// RawValues is an alias for []RawValue
type RawValues []RawValue

// RawValuesFrom creates RawValues from given values
func RawValuesFrom(values ...string) RawValues {
	rawValues := make([]RawValue, len(values))
	for index, val := range values {
		rawValues[index] = RawValue{Value: val, Display: val}
	}
	return rawValues
}

// FilterPrefix filters values with given prefix
func (r RawValues) FilterPrefix(prefix string) RawValues {
	filtered := make(RawValues, 0)
	for _, r := range r {
		if strings.HasPrefix(r.Value, prefix) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// ByValue alias to filter by value
type ByValue []RawValue

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Less(i, j int) bool { return a[i].Value < a[j].Value }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Filter filters values with given prefix
func (a ByValue) Filter(prefix string) []RawValue {
	filtered := make([]RawValue, 0, len(a))
	for _, v := range a {
		if strings.HasPrefix(v.Value, prefix) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// ByDisplay alias to filter by display
type ByDisplay []RawValue

func (a ByDisplay) Len() int           { return len(a) }
func (a ByDisplay) Less(i, j int) bool { return a[i].Display < a[j].Display }
func (a ByDisplay) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
