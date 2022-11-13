package pflagfork

import (
	"reflect"

	"github.com/spf13/pflag"
)

type flagSet struct {
	*pflag.FlagSet
}

func (f flagSet) IsPosix() bool {
	if method := reflect.ValueOf(f.FlagSet).MethodByName("IsPosix"); method.IsValid() {
		if values := method.Call([]reflect.Value{}); len(values) == 1 && values[0].Kind() == reflect.Bool {
			return values[0].Bool()
		}
	}
	return true
}

func FlagSet(f *pflag.FlagSet) *flagSet {
	return &flagSet{FlagSet: f}
}
