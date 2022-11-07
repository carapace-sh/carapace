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

func (f Flag) OptargDelimiter() string {
	if field := reflect.ValueOf(f.Flag).Elem().FieldByName("OptargDelimiter"); field.IsValid() && field.Kind() == reflect.Int32 {
		return string(rune(field.Int()))
	}
	return "="
}

func (f Flag) IsRepeatable() bool {
	if strings.Contains(f.Value.Type(), "Slice") ||
		strings.Contains(f.Value.Type(), "Array") ||
		f.Value.Type() == "count" {
		return true
	}
	return false
}

func (f Flag) Split(arg string) (prefix, optarg string) {
	if delimiter := f.OptargDelimiter(); delimiter != "" {
		splitted := strings.SplitN(arg, f.OptargDelimiter(), 2)
		return splitted[0] + delimiter, splitted[1]
	}
	// TODO handle ""
	panic("TODO")
}

func (f Flag) Matches(arg string, posix bool) bool {
	if !strings.HasPrefix(arg, "-") { // not a flag
		return false
	}

	switch {

	case strings.HasPrefix(arg, "--"):
		name := strings.TrimPrefix(arg, "--")
		name = strings.SplitN(name, f.OptargDelimiter(), 2)[0]

		switch f.Style() {
		case ShorthandOnly, NameAsShorthand:
			return false
		default:
			return name == f.Name
		}

	case !posix:
		name := strings.TrimPrefix(arg, "-")
		name = strings.SplitN(name, f.OptargDelimiter(), 2)[0]

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
