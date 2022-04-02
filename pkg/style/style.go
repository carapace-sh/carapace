// Package style provide display coloring
package style

import (
	"fmt"
	"regexp"
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
	Gray    string = Of(Dim, White)

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

// Color256 returns style for 256 color
func Color256(i int) string {
	if i < 0 || i > 255 {
		return Default
	}
	return fmt.Sprintf("color%d", i)
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

// FormatAnsi formats given string with given style as ansi escape sequence
func FormatAnsi(s, _style string) string {
	reColor256 := regexp.MustCompile(`^color(?P<number>\d+)$`)

	result := make([]string, 0)
	for _, word := range strings.Split(_style, " ") {
		switch word {
		case Red:
			result = append(result, "\033[31m")
		case Green:
			result = append(result, "\033[32m")
		case Yellow:
			result = append(result, "\033[33m")
		case Blue:
			result = append(result, "\033[34m")
		case Magenta:
			result = append(result, "\033[35m")
		case Cyan:
			result = append(result, "\033[36m")

		case BrightBlack:
			result = append(result, "\033[90m")
		case BrightRed:
			result = append(result, "\033[91m")
		case BrightGreen:
			result = append(result, "\033[92m")
		case BrightYellow:
			result = append(result, "\033[93m")
		case BrightBlue:
			result = append(result, "\033[94m")
		case BrightMagenta:
			result = append(result, "\033[95m")
		case BrightCyan:
			result = append(result, "\033[96m")

		case BgRed:
			result = append(result, "\033[41m")
		case BgGreen:
			result = append(result, "\033[42m")
		case BgYellow:
			result = append(result, "\033[43m")
		case BgBlue:
			result = append(result, "\033[44m")
		case BgMagenta:
			result = append(result, "\033[45m")
		case BgCyan:
			result = append(result, "\033[46m")

		case BgBrightBlack:
			result = append(result, "\033[100m")
		case BgBrightRed:
			result = append(result, "\033[101m")
		case BgBrightGreen:
			result = append(result, "\033[102m")
		case BgBrightYellow:
			result = append(result, "\033[103m")
		case BgBrightBlue:
			result = append(result, "\033[104m")
		case BgBrightMagenta:
			result = append(result, "\033[105m")
		case BgBrightCyan:
			result = append(result, "\033[106m")

		case Bold:
			result = append(result, "\033[1m")
		case Dim:
			result = append(result, "\033[2m")
		case Italic:
			result = append(result, "\033[3m")
		case Underlined:
			result = append(result, "\033[4m")
		case Blink:
			result = append(result, "\033[5m")
		case Inverse:
			result = append(result, "\033[7m")
		default:
			if reColor256.MatchString(word) {
				result = append(result, fmt.Sprintf("\033[38;5;%vm", reColor256.FindStringSubmatch(word)[1]))
			}
		}
	}
	result = append(result, s)
	result = append(result, "\033[0m")
	return strings.Join(result, "")
}
