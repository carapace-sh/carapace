package style

import (
	"github.com/rsteube/carapace/internal/config"
)

// Register a style configuration
//
//	var Carapace = struct {
//		Value       string `desc:"default style for values"`
//		Description string `desc:"default style for descriptions"`
//	}{
//		Value:       Default,
//		Description: Gray,
//	}
//
//	func init() {
//		Register("carapace", &Carapace)
//	}
func Register(name string, i interface{}) { config.RegisterStyle(name, i) }

// Set a style
//
//	Set("carapace.Value", "bold magenta")
func Set(key, value string) error { return config.SetStyle(key, value) }

var Carapace = struct {
	Value       string `desc:"default style for values"`
	Description string `desc:"default style for descriptions"`
	Error       string `desc:"default style for errors"`

	KeywordAmbiguous string `desc:"keyword describing a ambiguous state"`
	KeywordNegative  string `desc:"keyword describing a negative state"`
	KeywordPositive  string `desc:"keyword describing a positive state"`
	KeywordUnknown   string `desc:"keyword describing an unknown state"`

	LogLevelTrace    string `desc:"LogLevel TRACE"`
	LogLevelDebug    string `desc:"LogLevel DEBUG"`
	LogLevelInfo     string `desc:"LogLevel INFO"`
	LogLevelWarning  string `desc:"LogLevel WARNING"`
	LogLevelError    string `desc:"LogLevel ERROR"`
	LogLevelCritical string `desc:"LogLevel CRITICAL"`
	LogLevelFatal    string `desc:"LogLevel FATAL"`

	H1 string `desc:"Highlight 1"`
	H2 string `desc:"Highlight 2"`
	H3 string `desc:"Highlight 3"`
	H4 string `desc:"Highlight 4"`
	H5 string `desc:"Highlight 5"`
	H6 string `desc:"Highlight 6"`

	H7  string `desc:"Highlight 7"`
	H8  string `desc:"Highlight 8"`
	H9  string `desc:"Highlight 9"`
	H10 string `desc:"Highlight 10"`

	H11 string `desc:"Highlight 11"`
	H12 string `desc:"Highlight 12"`
}{
	Value:       Default,
	Description: Gray,
	Error:       Of(Bold, Red),

	KeywordAmbiguous: Yellow,
	KeywordNegative:  Red,
	KeywordPositive:  Green,
	KeywordUnknown:   Gray,

	LogLevelTrace:    Blue,
	LogLevelDebug:    Gray,
	LogLevelInfo:     Green,
	LogLevelWarning:  Yellow,
	LogLevelError:    Magenta,
	LogLevelCritical: Red,
	LogLevelFatal:    Cyan,

	H1: Blue,
	H2: Yellow,
	H3: Magenta,
	H4: Cyan,
	H5: Green,

	H6:  Of(Blue, Dim),
	H7:  Of(Yellow, Dim),
	H8:  Of(Magenta, Dim),
	H9:  Of(Cyan, Dim),
	H10: Of(Green, Dim),

	H11: Bold,
	H12: Of(Bold, Dim),
}

func init() {
	Register("carapace", &Carapace)
}
