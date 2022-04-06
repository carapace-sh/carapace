package style

import "strings"

var keywords = map[string]*string{
	"yes": &Carapace.KeywordPositive,
	"no":  &Carapace.KeywordNegative,

	"true":  &Carapace.KeywordPositive,
	"false": &Carapace.KeywordNegative,

	"on":  &Carapace.KeywordPositive,
	"off": &Carapace.KeywordNegative,

	"all":  &Carapace.KeywordPositive,
	"none": &Carapace.KeywordNegative,

	"always": &Carapace.KeywordPositive,
	"auto":   &Carapace.KeywordAmbiguous,
	"never":  &Carapace.KeywordNegative,

	"start":      &Carapace.KeywordNegative,
	"started":    &Carapace.KeywordNegative,
	"running":    &Carapace.KeywordAmbiguous,
	"inprogress": &Carapace.KeywordAmbiguous,
	"pause":      &Carapace.KeywordAmbiguous,
	"paused":     &Carapace.KeywordAmbiguous,
	"stop":       &Carapace.KeywordNegative,
	"stopped":    &Carapace.KeywordNegative,

	"onsuccess": &Carapace.KeywordPositive,
	"onfailure": &Carapace.KeywordNegative,
	"onerror":   &Carapace.KeywordNegative,

	"success": &Carapace.KeywordPositive,
	"unknown": &Carapace.KeywordUnknown,
	"warn":    &Carapace.KeywordAmbiguous,
	"error":   &Carapace.KeywordNegative,
	"failed":  &Carapace.KeywordNegative,

	"nonblock": &Carapace.KeywordAmbiguous,
	"block":    &Carapace.KeywordNegative,

	"ondemand": &Carapace.KeywordAmbiguous,
}

var keywordReplacer = strings.NewReplacer(
	"-", "",
	"_", "",
)

// ForKeyword returns the style for given keyword
func ForKeyword(s string) string {
	if _style, ok := keywords[keywordReplacer.Replace(strings.ToLower(s))]; ok {
		return *_style
	}
	return Default
}
