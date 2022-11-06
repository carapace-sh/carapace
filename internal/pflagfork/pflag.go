package pflagfork

import (
	"reflect"

	"github.com/spf13/pflag"
)

// style defines how flags are represented.
type style int

const (
	Default         style = iota // default behaviour
	ShorthandOnly                // only the shorthand should be used
	NameAsShorthand              // non-posix style where the name is also added as shorthand (single `-` prefix)
)

func Style(flag *pflag.Flag) style {
	if field := reflect.ValueOf(flag).Elem().FieldByName("Style"); field.IsValid() && field.Kind() == reflect.Int {
		return style(field.Int())
	}
	return Default
}

func IsPosix(flagset *pflag.FlagSet) bool {
	if method := reflect.ValueOf(flagset).MethodByName("IsPosix"); method.IsValid() {
		if values := method.Call([]reflect.Value{}); len(values) == 1 && values[0].Kind() == reflect.Bool {
			return values[0].Bool()
		}
	}
	return true
}
