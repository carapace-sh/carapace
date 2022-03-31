// Package style provide display coloring
package style

import (
	"mime"
	"os"
	"strings"
)

// TODO full style support: e.g. "bg-bright-black bright-cyan bold underlined"
// https://elv.sh/ref/builtin.html#styled
var (
	Default string = "default"
	Red     string = "red"
	Green   string = "green"
	Yellow  string = "yellow"
	Blue    string = "blue"
	Magenta string = "magenta"
	Cyan    string = "cyan"

	BrightBlack string = "bright-black"
	//BrightRed     Color = "bright-red"
	//BrightGreen   Color = "bright-green"
	//BrightYellow  Color = "bright-yellow"
	//BrightBlue    Color = "bright-blue"
	//BrightMagenta Color = "bright-magenta"
	//BrightCyan    Color = "bright-cyan"
	//BrightWhite   Color = "bright-white"
)

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
