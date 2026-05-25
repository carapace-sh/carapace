package pflagfork

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// mode defines how flags are represented.
type mode int

const (
	Default         mode = iota // default behaviour
	ShorthandOnly               // only the shorthand should be used
	NameAsShorthand             // non-posix mode where the name is also added as shorthand (single `-` prefix)
)

type Flag struct {
	*pflag.Flag
	Prefix string
	Args   []string
}

func (f Flag) Nargs() int {
	if field := reflect.ValueOf(f.Flag).Elem().FieldByName("Nargs"); field.IsValid() && field.Kind() == reflect.Int {
		return int(field.Int())
	}
	return 0
}

func (f Flag) acceptsNext() bool {
	return f.ArgumentStyle() == 0 || f.ArgumentStyle()&1 != 0
}

func (f Flag) acceptsDelimited() bool {
	return f.ArgumentStyle() == 0 || f.ArgumentStyle()&2 != 0
}

func (f Flag) acceptsAttached() bool {
	return f.ArgumentStyle() == 0 || f.ArgumentStyle()&4 != 0
}

func (f Flag) OptargDelimiter() rune {
	if field := reflect.ValueOf(f.Flag).Elem().FieldByName("OptargDelimiter"); field.IsValid() && field.Kind() == reflect.Int32 {
		return (rune(field.Int()))
	}
	return '='
}

func (f Flag) ArgumentStyle() uint {
	if field := reflect.ValueOf(f.Flag).Elem().FieldByName("ArgumentStyle"); field.IsValid() && field.Kind() == reflect.Uint {
		return uint(field.Uint())
	}
	return 0 // 0 means accept all variants (backward compatible)
}

func (f Flag) GetMode() mode {
	if field := reflect.ValueOf(f.Flag).Elem().FieldByName("Mode"); field.IsValid() && field.Kind() == reflect.Int {
		return mode(field.Int())
	}
	return Default
}

func (f Flag) IsRepeatable() bool {
	if strings.Contains(f.Value.Type(), "Slice") ||
		strings.Contains(f.Value.Type(), "Array") ||
		f.Value.Type() == "count" {
		return true
	}
	return false
}

func (f Flag) TakesValue() bool {
	switch f.Value.Type() {
	case "bool", "boolSlice", "count":
		return false
	default:
		return true
	}
}

func (f Flag) IsOptarg() bool {
	return f.NoOptDefVal != ""
}

func (f Flag) Style() string {
	switch {
	case !f.TakesValue():
		return style.Carapace.FlagNoArg
	case f.IsOptarg():
		return style.Carapace.FlagOptArg
	case f.Nargs() != 0:
		return style.Carapace.FlagMultiArg
	default:
		return style.Carapace.FlagArg
	}
}

func (f Flag) Required() bool {
	if annotation := f.Annotations[cobra.BashCompOneRequiredFlag]; len(annotation) == 1 && annotation[0] == "true" {
		return true
	}
	return false
}

func (f Flag) Definition() string {
	var definition string
	switch f.GetMode() {
	case ShorthandOnly:
		definition = fmt.Sprintf("-%v", f.Shorthand)
	case NameAsShorthand:
		definition = fmt.Sprintf("-%v, -%v", f.Shorthand, f.Name)
	default:
		switch f.Shorthand {
		case "":
			definition = fmt.Sprintf("--%v", f.Name)
		default:
			definition = fmt.Sprintf("-%v, --%v", f.Shorthand, f.Name)
		}
	}

	if f.Hidden {
		definition += "&"
	}

	if f.Required() {
		definition += "!"
	}

	if f.IsRepeatable() {
		definition += "*"
	}

	switch {
	case f.IsOptarg():
		switch f.Value.Type() {
		case "bool", "boolSlice", "count":
		default:
			definition += "?"
		}
	case f.TakesValue():
		definition += "="
	}

	return definition
}

func (f Flag) Consumes(arg string) bool {
	switch {
	case f.Flag == nil:
		return false
	case !f.TakesValue():
		return false
	case f.IsOptarg():
		return false
	case len(f.Args) == 0:
		return true
	case f.Nargs() > 1 && len(f.Args) < f.Nargs():
		return true
	case f.Nargs() < 0 && !strings.HasPrefix(arg, "-"):
		return true
	default:
		return false
	}
}
