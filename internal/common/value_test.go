package common

import (
	"sort"
	"testing"
)

func TestTrimmedDescription(t *testing.T) {
	r := RawValue{
		Description: "123456789012345678901234567890123456789012345678901234567890123456789012345678901",
	}
	if r.TrimmedDescription() != "12345678901234567890123456789012345678901234567890123456789012345678901234567..." {
		t.Error("description should be trimmed to 80 characters")
	}
}

func TestRawValuesFrom(t *testing.T) {
	v := RawValuesFrom("first", "second")
	if !equalRawValues(v[0], RawValue{
		Value:       "first",
		Display:     "first",
		Description: "",
	}) {
		t.Fail()
	}

	if !equalRawValues(v[1], RawValue{
		Value:       "second",
		Display:     "second",
		Description: "",
	}) {
		t.Fail()
	}
}

func TestFilterPrefix(t *testing.T) {
	v := RawValuesFrom("first", "second").FilterPrefix("sec")
	if len(v) != 1 && !equalRawValues(v[0], RawValue{
		Value:       "second",
		Display:     "second",
		Description: "",
	}) {
		t.Fail()
	}
}

func equalRawValues(a, b RawValue) bool {
	return a.Value == b.Value && a.Display == b.Display && a.Description == b.Description
}

func TestSort(t *testing.T) {
	r := RawValuesFrom("3", "2", "1")
	sort.Sort(ByValue(r))
	if r[0].Value != "1" {
		t.Fail()
	}

	r = RawValuesFrom("3", "2", "1")
	sort.Sort(ByDisplay(r))
	if r[0].Value != "1" {
		t.Fail()
	}
}
