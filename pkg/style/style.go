// Package style provide display coloring
package style

import (
	"mime"
	"os"
	"strings"
)

var (
	Default string = "default"
	Red     string = "red"
	Green   string = "green"
	Yellow  string = "yellow"
	Blue    string = "blue"
	Magenta string = "magenta"
	Cyan    string = "cyan"

	BrightBlack   string = "bright-black"
	BrightRed     string = "bright-red"
	BrightGreen   string = "bright-green"
	BrightYellow  string = "bright-yellow"
	BrightBlue    string = "bright-blue"
	BrightMagenta string = "bright-magenta"
	BrightCyan    string = "bright-cyan"

	BgRed     string = "bg-red"
	BgGreen   string = "bg-green"
	BgYellow  string = "bg-yellow"
	BgBlue    string = "bg-blue"
	BgMagenta string = "bg-magenta"
	BgCyan    string = "bg-cyan"

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
func Of(s ...string) string {
	return strings.Join(s, " ")
}

// ForPath returns the style for given path
func ForPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return Default
		}
		path = strings.Replace(path, "~/", home+"/", 1)
	}

	stat, err := os.Lstat(path)
	if err != nil {
		return Default
	}
	if stat.IsDir() {
		return Blue
	}

	if stat.Mode()&os.ModeSymlink == os.ModeSymlink {
		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			return Cyan
		}
		return Yellow
	} else if stat.Mode()&0111 == 0111 { // any executable
		return Green
	}

	if stat.Mode()&os.ModeSymlink == os.ModeSymlink {
		return Yellow
	}

	if index := strings.LastIndex(path, "."); index != -1 {
		if mime := mime.TypeByExtension(path[index:]); strings.HasPrefix(mime, "image") {
			return Magenta
		} else if strings.HasPrefix(mime, "application") {
			return Red
		}
	}
	return Default
}
