package pflagfork

import (
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

func TestLookupPosixShorthandArg(t *testing.T) {
	_test := func(arg, name string, args ...string) {
		t.Run(arg, func(t *testing.T) {
			if args == nil {
				args = []string{}
			}

			fs := &FlagSet{pflag.NewFlagSet("test", pflag.PanicOnError)}

			fs.BoolP("bool", "b", false, "")
			fs.CountP("count", "c", "")
			fs.StringP("string", "s", "", "")

			f := fs.lookupPosixShorthandArg(arg)
			if f == nil || f.Name != name {
				t.Fatalf("should be " + name)
			}

			if !reflect.DeepEqual(f.Args, args) {
				t.Fatalf("args dont match %v: actual: %#v expected: %#v", arg, f.Args, args)
			}

		})
	}

	_test("-b=", "bool", "")
	_test("-b=t", "bool", "t")
	_test("-b=true", "bool", "true")
	_test("-ccb", "bool")
	_test("-ccb=", "bool", "")
	_test("-ccb=t", "bool", "t")
	_test("-ccb=true", "bool", "true")
	_test("-ccbs=val1", "string", "val1")
	_test("-ccbsval1", "string", "val1")
}
