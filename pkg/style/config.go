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

type carapace struct {
	Value       string `desc:"default style for values" tag:"core styles"`
	Description string `desc:"default style for descriptions" tag:"core styles"`
	Error       string `desc:"default style for errors" tag:"core styles"`

	KeywordAmbiguous string `desc:"keyword describing a ambiguous state" tag:"keyword styles"`
	KeywordNegative  string `desc:"keyword describing a negative state" tag:"keyword styles"`
	KeywordPositive  string `desc:"keyword describing a positive state" tag:"keyword styles"`
	KeywordUnknown   string `desc:"keyword describing an unknown state" tag:"keyword styles"`

	LogLevelTrace    string `desc:"LogLevel TRACE" tag:"loglevel styles"`
	LogLevelDebug    string `desc:"LogLevel DEBUG" tag:"loglevel styles"`
	LogLevelInfo     string `desc:"LogLevel INFO" tag:"loglevel styles"`
	LogLevelWarning  string `desc:"LogLevel WARNING" tag:"loglevel styles"`
	LogLevelError    string `desc:"LogLevel ERROR" tag:"loglevel styles"`
	LogLevelCritical string `desc:"LogLevel CRITICAL" tag:"loglevel styles"`
	LogLevelFatal    string `desc:"LogLevel FATAL" tag:"loglevel styles"`

	Highlight1 string `desc:"Highlight 1" tag:"highlight styles"`
	Highlight2 string `desc:"Highlight 2" tag:"highlight styles"`
	Highlight3 string `desc:"Highlight 3" tag:"highlight styles"`
	Highlight4 string `desc:"Highlight 4" tag:"highlight styles"`
	Highlight5 string `desc:"Highlight 5" tag:"highlight styles"`
	Highlight6 string `desc:"Highlight 6" tag:"highlight styles"`

	Highlight7  string `desc:"Highlight 7" tag:"highlight styles"`
	Highlight8  string `desc:"Highlight 8" tag:"highlight styles"`
	Highlight9  string `desc:"Highlight 9" tag:"highlight styles"`
	Highlight10 string `desc:"Highlight 10" tag:"highlight styles"`

	Highlight11 string `desc:"Highlight 11" tag:"highlight styles"`
	Highlight12 string `desc:"Highlight 12" tag:"highlight styles"`
}

var Carapace = carapace{
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

	Highlight1: Blue,
	Highlight2: Yellow,
	Highlight3: Magenta,
	Highlight4: Cyan,
	Highlight5: Green,

	Highlight6:  Of(Blue, Dim),
	Highlight7:  Of(Yellow, Dim),
	Highlight8:  Of(Magenta, Dim),
	Highlight9:  Of(Cyan, Dim),
	Highlight10: Of(Green, Dim),

	Highlight11: Bold,
	Highlight12: Of(Bold, Dim),
}

// Highlight returns the style for given level (0..n)
func (c carapace) Highlight(level int) string {
	switch level {
	case 0:
		return c.Highlight1
	case 1:
		return c.Highlight2
	case 2:
		return c.Highlight3
	case 3:
		return c.Highlight4
	case 4:
		return c.Highlight5
	case 5:
		return c.Highlight6
	case 6:
		return c.Highlight7
	case 7:
		return c.Highlight8
	case 8:
		return c.Highlight9
	case 9:
		return c.Highlight10
	case 10:
		return c.Highlight11
	case 11:
		return c.Highlight12
	default:
		return Default
	}

}

func init() {
	Register("carapace", &Carapace)
}
