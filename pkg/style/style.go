// Package style provide display coloring
package style

import (
	"strings"

	"github.com/rsteube/carapace/internal/lscolors"
)

var (
	Default string = "default"

	Black   string = "black"
	Red     string = "red"
	Green   string = "green"
	Yellow  string = "yellow"
	Blue    string = "blue"
	Magenta string = "magenta"
	Cyan    string = "cyan"
	White   string = "white"

	BrightBlack   string = "bright-black"
	BrightRed     string = "bright-red"
	BrightGreen   string = "bright-green"
	BrightYellow  string = "bright-yellow"
	BrightBlue    string = "bright-blue"
	BrightMagenta string = "bright-magenta"
	BrightCyan    string = "bright-cyan"

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

var ansi = map[string]string{
	"30": Black,
	"31": Red,
	"32": Green,
	"33": Yellow,
	"34": Blue,
	"35": Magenta,
	"36": Cyan,
	"37": White,

	"90": BrightBlack,
	"40": BgBlack,
	"41": BgRed,
	"42": BgGreen,
	"43": BgYellow,
	"44": BgBlue,
	"45": BgMagenta,
	"46": BgCyan,
	"47": BgWhite,

	"100": BgBrightBlack,
	"101": BgBrightRed,
	"102": BgBrightGreen,
	"103": BgBrightYellow,
	"104": BgBrightBlue,
	"105": BgBrightMagenta,
	"106": BgBrightCyan,
	"107": BgBrightWhite,

	"01": Bold,
	"02": Dim,
	"03": Italic,
	"04": Underlined,
	"05": Blink,
	"07": Inverse,
}

// Of combines different styles
func Of(s ...string) string {
	return strings.Join(s, " ")
}

// ForPath returns the style for given path
func ForPath(path string) string {
	if ansiStyle := lscolors.GetColorist().GetStyle(path); ansiStyle != "" {
		styles := make([]string, 0)
		for _, code := range strings.Split(ansiStyle, ";") {
			if style, ok := ansi[code]; ok {
				styles = append(styles, style)
			}
		}
		return Of(styles...)
	}
	return Default
}
