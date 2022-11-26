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

type flag struct {
	*pflag.Flag
}

func (f flag) Style() style {
	if field := reflect.ValueOf(f.Flag).Elem().FieldByName("Style"); field.IsValid() && field.Kind() == reflect.Int {
		return style(field.Int())
	}
	return Default
}

func Flag(f *pflag.Flag) *flag {
	return &flag{Flag: f}
}
