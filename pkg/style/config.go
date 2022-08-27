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
}

func init() {
	Register("carapace", &Carapace)
}
