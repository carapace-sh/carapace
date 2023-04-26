package pflagfork

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestIsShorthandSeries(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.PanicOnError)
	fs.BoolP("bool", "b", false, "")
	fs.CountP("count", "c", "")
	fs.StringP("optarg", "o", "", "")
	fs.StringP("optarg-custom", "x", "", "")
	fs.StringP("string", "s", "", "")

	fs.Lookup("optarg").NoOptDefVal = " "
	fs.Lookup("optarg-custom").NoOptDefVal = " "
	fs.Lookup("optarg-custom").OptargDelimiter = ':'

	_test := func(arg string, match bool) {
		if (FlagSet{fs}).IsShorthandSeries(arg) != match {
			t.Errorf("failed to match #%v with %#v", arg, match)
		}
	}

	_test("-a", false)
	_test("-b", true)
	_test("-c", true)
	_test("-o", true)
	_test("-s", false)

	_test("-o=", false)
	_test("-o:", true)
	_test("-o:", true)
	_test("-o", true)

	_test("-bc", true)
	_test("--a", false)
	_test("--b", false)
	_test("--ab", false)
}
