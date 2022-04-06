// Package style provide display coloring
package style

import (
	"strings"

	"github.com/rsteube/carapace/third_party/github.com/elves/elvish/pkg/ui"
)

var (
	Default string = ""

	Black   string = "black"
	Red     string = "red"
	Green   string = "green"
	Yellow  string = "yellow"
	Blue    string = "blue"
	Magenta string = "magenta"
	Cyan    string = "cyan"
	White   string = "white"
	Gray    string = Of(Dim, White)

	BrightBlack   string = "bright-black"
	BrightRed     string = "bright-red"
	BrightGreen   string = "bright-green"
	BrightYellow  string = "bright-yellow"
	BrightBlue    string = "bright-blue"
	BrightMagenta string = "bright-magenta"
	BrightCyan    string = "bright-cyan"
	BrightWhite   string = "bright-white"

	BgBlack   string = "bg-black"
	BgRed     string = "bg-red"
	BgGreen   string = "bg-green"
	BgYellow  string = "bg-yellow"
	BgBlue    string = "bg-blue"
	BgMagenta string = "bg-magenta"
	BgCyan    string = "bg-cyan"
	BgWhite   string = "bg-white"

	BgBrightBlack   string = "bg-bright-black"
	BgBrightRed     string = "bg-bright-red"
	BgBrightGreen   string = "bg-bright-green"
	BgBrightYellow  string = "bg-bright-yellow"
	BgBrightBlue    string = "bg-bright-blue"
	BgBrightMagenta string = "bg-bright-magenta"
	BgBrightCyan    string = "bg-bright-cyan"
	BgBrightWhite   string = "bg-bright-white"

	Bold       string = "bold"
	Dim        string = "dim"
	Italic     string = "italic"
	Underlined string = "underlined"
	Blink      string = "blink"
	Inverse    string = "inverse"
)

// Of combines different styles
func Of(s ...string) string { return strings.Join(s, " ") }

// XTerm256Color returns a color from the xterm 256-color palette.
func XTerm256Color(i uint8) string { return ui.XTerm256Color(i).String() }

// TrueColor returns a 24-bit true color.
func TrueColor(r, g, b uint8) string { return ui.TrueColor(r, g, b).String() }

// SGR returns the SGR sequence for given style
func SGR(s string) string { return parseStyle(s).SGR() }

func parseStyle(s string) ui.Style {
	stylings := make([]ui.Styling, 0)
	for _, word := range strings.Split(s, " ") {
		if styling := ui.ParseStyling(word); styling != nil {
			stylings = append(stylings, styling)
		}
	}
	return ui.ApplyStyling(ui.Style{}, stylings...)
}
