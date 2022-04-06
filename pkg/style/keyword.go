package style

import "strings"

var keywords = map[string]string{
	"yes": Green,
	"no":  Red,

	"true":  Green,
	"false": Red,

	"on":  Green,
	"off": Red,

	"all":  Green,
	"none": Red,

	"always": Green,
	"auto":   Yellow,
	"never":  Red,

	"start":      Red,
	"started":    Red,
	"running":    Yellow,
	"inprogress": Yellow,
	"pause":      Yellow,
	"paused":     Yellow,
	"stop":       Red,
	"stopped":    Red,

	"onsuccess": Green,
	"onfailure": Red,
	"onerror":   Red,

	"success": Green,
	"unknown": Gray,
	"warn":    Yellow,
	"error":   Red,
	"failed":  Red,

	"nonblock": Yellow,
	"block":    Red,

	"ondemand": Yellow,
}

var keywordReplacer = strings.NewReplacer(
	"-", "",
	"_", "",
)

// ForKeyword returns the style for given keyword
func ForKeyword(s string) string {
	if _style, ok := keywords[keywordReplacer.Replace(strings.ToLower(s))]; ok {
		return _style
	}
	return Default
}
