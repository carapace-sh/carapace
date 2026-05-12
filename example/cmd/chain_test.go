package cmd

import (
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
)

func TestShorthandChain(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("chain", "-b").
			Expect(carapace.ActionImport([]byte(`{
  "nospace": "co",
  "values": [
    {"value": "-b", "display": "-b", "tag": "shorthand flags"},
    {"value": "-bc", "display": "c", "tag": "shorthand flags"},
    {"value": "-bo", "display": "o", "style": "yellow", "tag": "shorthand flags"},
    {"value": "-bv", "display": "v", "style": "blue", "tag": "shorthand flags"}
  ]
}`)))

		s.Run("chain", "-bc").
			Expect(carapace.ActionImport([]byte(`{
  "nospace": "o",
  "values": [
    {"value": "-bc", "display": "-bc", "tag": "shorthand flags"},
    {"value": "-bco", "display": "o", "style": "yellow", "tag": "shorthand flags"},
    {"value": "-bcv", "display": "v", "style": "blue", "tag": "shorthand flags"}
  ]
}`)))

		s.Run("chain", "-bcc").
			Expect(carapace.ActionImport([]byte(`{
  "nospace": "o",
  "values": [
    {"value": "-bcc", "display": "-bcc", "tag": "shorthand flags"},
    {"value": "-bcco", "display": "o", "style": "yellow", "tag": "shorthand flags"},
    {"value": "-bccv", "display": "v", "style": "blue", "tag": "shorthand flags"}
  ]
}`)))

		s.Run("chain", "-bcco").
			Expect(carapace.ActionImport([]byte(`{
  "nospace": "c",
  "values": [
    {"value": "-bcco", "display": "-bcco", "style": "yellow", "tag": "shorthand flags"},
    {"value": "-bccoc", "display": "c", "tag": "shorthand flags"},
    {"value": "-bccov", "display": "v", "style": "blue", "tag": "shorthand flags"}
  ]
}`)))

		s.Run("chain", "-bcco", "").
			Expect(carapace.ActionValues(
				"p1",
				"positional1",
			))

		s.Run("chain", "-bcco=").
			Expect(carapace.ActionValues(
				"opt1",
				"opt2",
			).Prefix("-bcco="))

		s.Run("chain", "-bccv", "").
			Expect(carapace.ActionValues(
				"val1",
				"val2",
			))

		s.Run("chain", "-bccv=").
			Expect(carapace.ActionValues(
				"val1",
				"val2",
			).Prefix("-bccv="))

		s.Run("chain", "-bccv", "val1", "-c").
			Expect(carapace.ActionImport([]byte(`{
  "nospace": "o",
  "values": [
    {"value": "-c", "display": "-c", "tag": "shorthand flags"},
    {"value": "-co", "display": "o", "style": "yellow", "tag": "shorthand flags"}
  ]
}`)))
	})
}
