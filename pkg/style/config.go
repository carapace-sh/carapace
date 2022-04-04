package style

import (
	"github.com/rsteube/carapace/internal/config"
)

// Register a style configuration
//   var Carapace = struct {
//   	Value       string `desc:"default style for values"`
//   	Description string `desc:"default style for descriptions"`
//   }{
//   	Value:       Default,
//   	Description: Gray,
//   }
//
//   func init() {
//   	Register("carapace", &Carapace)
//   }
func Register(name string, i interface{}) { config.RegisterStyle(name, i) }

// Set a style
//   Set("carapace.Value", "bold magenta")
func Set(key, value string) error { return config.SetStyle(key, value) }

var Carapace = struct {
	Value       string `desc:"default style for values"`
	Description string `desc:"default style for descriptions"`
}{
	Value:       Default,
	Description: Gray,
}

func init() {
	Register("carapace", &Carapace)
}
