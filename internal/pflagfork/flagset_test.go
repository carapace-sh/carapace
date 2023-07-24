package pflagfork

import (
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

func TestLookupPosixShorthandArg(t *testing.T) {
	_test := func(arg, name, prefix string, args ...string) {
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

			if f.Prefix != prefix {
				t.Fatalf("prefix doesnt match actual: %#v, expected: %#v", f.Prefix, prefix)
			}

			if !reflect.DeepEqual(f.Args, args) {
				t.Fatalf("args dont match %v: actual: %#v expected: %#v", arg, f.Args, args)
			}

		})
	}

	_test("-b=", "bool", "-b=", "")
	_test("-b=t", "bool", "-b=", "t")
	_test("-b=true", "bool", "-b=", "true")
	_test("-ccb", "bool", "-ccb")
	_test("-ccb=", "bool", "-ccb=", "")
	_test("-ccb=t", "bool", "-ccb=", "t")
	_test("-ccb=true", "bool", "-ccb=", "true")
	_test("-ccbs=val1", "string", "-ccbs=", "val1")
	_test("-ccbsval1", "string", "-ccbs", "val1")
}
