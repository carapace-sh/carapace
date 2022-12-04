package common

import "testing"

func TestSuffixMatcherAdd(t *testing.T) {
	sm := SuffixMatcher{""}

	sm.Add('*')
	if sm.string != "*" {
		t.Errorf(`should be "*" [was: "%v"]`, sm)
	}
}
