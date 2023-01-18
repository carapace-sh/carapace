package pflagfork

import (
	"reflect"
	"strings"

	"github.com/spf13/pflag"
)

// style defines how flags are represented.
type style int

const (
	Default         style = iota // default behaviour
	ShorthandOnly                // only the shorthand should be used
	NameAsShorthand              // non-posix style where the name is also added as shorthand (single `-` prefix)
)

type Flag struct {
	*pflag.Flag
}

func (f Flag) Style() style {
	if field := reflect.ValueOf(f.Flag).Elem().FieldByName("Style"); field.IsValid() && field.Kind() == reflect.Int {
		return style(field.Int())
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

func (f Flag) Matches(arg string, posix bool) bool {
	if !strings.HasPrefix(arg, "-") { // not a flag
		return false
	}

	switch {

	case strings.HasPrefix(arg, "--"):
		name := strings.TrimPrefix(arg, "--")
		name = strings.SplitN(name, "=", 2)[0]

		switch f.Style() {
		case ShorthandOnly, NameAsShorthand:
			return false
		default:
			return name == f.Name
		}

	case !posix:
		name := strings.TrimPrefix(arg, "-")
		name = strings.SplitN(name, "=", 2)[0]

		if name == "" {
			return false
		}

		switch f.Style() {
		case ShorthandOnly:
			return name == f.Shorthand
		default:
			return name == f.Name || name == f.Shorthand
		}

	default:
		if f.Shorthand != "" {
			return strings.HasSuffix(arg, f.Shorthand)
		}
		return false
	}
}

// func lookupFlag(arg string) (flag *pflag.Flag) {

// 	if strings.HasPrefix(arg, "--") {
// 		flag = cmd.Flags().Lookup(nameOrShorthand)
// 	} else if strings.HasPrefix(arg, "-") && len(nameOrShorthand) > 0 {
// 		if (pflagfork.FlagSet{FlagSet: cmd.Flags()}).IsPosix() {
// 			flag = cmd.Flags().ShorthandLookup(string(nameOrShorthand[len(nameOrShorthand)-1]))
// 		} else {
// 			flag = cmd.Flags().ShorthandLookup(nameOrShorthand)
// 		}
// 	}
// 	return
// }

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
