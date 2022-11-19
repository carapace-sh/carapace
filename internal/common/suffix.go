package common

import (
	"sort"
	"strings"
)

type SuffixMatcher string

func (sm SuffixMatcher) Add(s string) SuffixMatcher {
	if strings.Contains(string(sm), "*") || strings.Contains(s, "*") {
		return SuffixMatcher("*")
	}

	unique := []rune(sm)
	for _, r := range []rune(s) {
		if !strings.Contains(string(sm), string(r)) {
			unique = append(unique, r)
		}
	}
	sort.Sort(ByRune(unique))
	return SuffixMatcher(unique)
}

func (sm SuffixMatcher) Matches(s string) bool {
	for _, r := range []rune(sm) {
		if r == '*' || strings.HasSuffix(s, string(r)) {
			return true
		}
	}
	return false
}

type ByRune []rune

func (r ByRune) Len() int           { return len(r) }
func (r ByRune) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByRune) Less(i, j int) bool { return r[i] < r[j] }
