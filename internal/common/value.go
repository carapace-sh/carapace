package common

import "strings"

type RawValue struct {
	Value       string
	Display     string
	Description string
}

func RawValuesFrom(values ...string) []RawValue {
	rawValues := make([]RawValue, len(values))
	for index, val := range values {
		rawValues[index] = RawValue{Value: val, Display: val}
	}
	return rawValues
}

type ByDisplay []RawValue

func (a ByDisplay) Len() int           { return len(a) }
func (a ByDisplay) Less(i, j int) bool { return a[i].Display < a[j].Display }
func (a ByDisplay) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByValues []RawValue

func (r ByValues) Filter(prefix string) []RawValue {
	filtered := make([]RawValue, 0, len(r))
	for _, v := range r {
		if strings.HasPrefix(v.Value, prefix) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}
