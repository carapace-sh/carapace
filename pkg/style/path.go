package style

import (
	"github.com/rsteube/carapace/third_party/github.com/elves/elvish/pkg/cli/lscolors"
	"github.com/rsteube/carapace/third_party/github.com/elves/elvish/pkg/ui"
)

// ForPath returns the style for given path
//
//	/tmp/locally/reachable/file.txt
func ForPath(path string) string { return fromSGR(lscolors.GetColorist().GetStyle(path)) }

// ForPath returns the style for given path by extension only
//
//	/tmp/non/existing/file.txt
func ForPathExt(path string) string { return fromSGR(lscolors.GetColorist().GetStyleExt(path)) }

func fromSGR(sgr string) string {
	s := ui.StyleFromSGR(sgr)
	result := []string{"fg-default", "bg-default"}
	if s.Foreground != nil {
		result = append(result, s.Foreground.String())
	}
	if s.Background != nil {
		result = append(result, "bg-"+s.Background.String())
	}
	if s.Bold {
		result = append(result, Bold)
	}
	if s.Dim {
		result = append(result, Dim)
	}
	if s.Italic {
		result = append(result, Italic)
	}
	if s.Underlined {
		result = append(result, Underlined)
	}
	if s.Blink {
		result = append(result, Blink)
	}
	if s.Inverse {
		result = append(result, Inverse)
	}
	return Of(result...)
}
