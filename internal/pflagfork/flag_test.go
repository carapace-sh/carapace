package pflagfork

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestMatches(t *testing.T) {
	fs := pflag.NewFlagSet("match", pflag.PanicOnError)
	fs.Bool("bool", false, "")
	fs.String("string", "", "")
	fs.String("optarg", "", "")
	fs.String("optarg-custom", "", "")

	fs.Lookup("optarg").NoOptDefVal = " "
	fs.Lookup("optarg-custom").NoOptDefVal = " "
	fs.Lookup("optarg-custom").OptargDelimiter = ':'

	_test := func(flag, arg string, match bool) {
		if (Flag{fs.Lookup(flag)}.NameMatches(arg)) != match {
			t.Errorf("failed to match flag #%v with arg %#v", flag, arg)
		}
	}

	_test("bool", "--bool", true)
	_test("bool", "-bool", false)
	_test("bool", "--bool=", true)
	_test("bool", "--boolx", false)
	_test("string", "--string", true)
	_test("string", "-string", false)
	_test("string", "--string=", false)
	_test("string", "--stringx", false)
	_test("optarg", "--optarg", true)
	_test("optarg", "-optarg", false)
	_test("optarg", "--optarg:", false)
	_test("optarg", "--optarg=", true)
	_test("optarg", "--optarg:val", false)
	_test("optarg", "--optarg=val", true)
	_test("optarg-custom", "--optarg-custom", true)
	_test("optarg-custom", "-optarg-custom", false)
	_test("optarg-custom", "--optarg-custom:", true)
	_test("optarg-custom", "--optarg-custom=", false)
	_test("optarg-custom", "--optarg-custom:val", true)
	_test("optarg-custom", "--optarg-custom=val", false)
}
