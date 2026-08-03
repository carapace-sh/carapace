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
				t.Fatal("should be " + name)
			}

			if f.ArgPrefix != prefix {
				t.Fatalf("prefix doesnt match actual: %#v, expected: %#v", f.ArgPrefix, prefix)
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

func TestArgumentStyleAcceptance(t *testing.T) {
	// Test AcceptNext only accepts next argument style
	t.Run("AcceptNext_rejects_delimited", func(t *testing.T) {
		fs := &FlagSet{pflag.NewFlagSet("test", pflag.ContinueOnError)}
		fs.String("accept-next", "", "Accept next only")
		fs.Lookup("accept-next").ArgumentStyle = pflag.AcceptNext

		err := fs.Parse([]string{"--accept-next=value"})
		if err == nil {
			t.Error("expected error for delimited style, got nil")
		}
	})

	// Test AcceptDelimited only accepts delimited argument style
	t.Run("AcceptDelimited_rejects_next", func(t *testing.T) {
		fs := &FlagSet{pflag.NewFlagSet("test", pflag.ContinueOnError)}
		fs.String("accept-delimited", "", "Accept delimited only")
		fs.Lookup("accept-delimited").ArgumentStyle = pflag.AcceptDelimited

		err := fs.Parse([]string{"--accept-delimited", "value"})
		if err == nil {
			t.Error("expected error for next style, got nil")
		}
	})

	// Test AcceptDelimited accepts delimited argument style
	t.Run("AcceptDelimited_accepts_delimited", func(t *testing.T) {
		fs := &FlagSet{pflag.NewFlagSet("test", pflag.ContinueOnError)}
		fs.String("accept-delimited", "", "Accept delimited only")
		fs.Lookup("accept-delimited").ArgumentStyle = pflag.AcceptDelimited

		err := fs.Parse([]string{"--accept-delimited=value"})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if fs.Lookup("accept-delimited").Value.String() != "value" {
			t.Errorf("expected 'value', got %s", fs.Lookup("accept-delimited").Value.String())
		}
	})

	// Test AcceptDelimited|AcceptNext accepts both
	t.Run("AcceptDelimitedOrNext_accepts_both", func(t *testing.T) {
		fs := &FlagSet{pflag.NewFlagSet("test", pflag.ContinueOnError)}
		fs.String("accept-both", "", "Accept both")
		fs.Lookup("accept-both").ArgumentStyle = pflag.AcceptDelimited | pflag.AcceptNext

		err := fs.Parse([]string{"--accept-both=value"})
		if err != nil {
			t.Errorf("expected no error for delimited, got %v", err)
		}

		fs = &FlagSet{pflag.NewFlagSet("test", pflag.ContinueOnError)}
		fs.String("accept-both", "", "Accept both")
		fs.Lookup("accept-both").ArgumentStyle = pflag.AcceptDelimited | pflag.AcceptNext
		err = fs.Parse([]string{"--accept-both", "value"})
		if err != nil {
			t.Errorf("expected no error for next, got %v", err)
		}
	})

	// Test shorthand with AcceptNext
	t.Run("AcceptNext_shorthand_rejects_delimited", func(t *testing.T) {
		fs := &FlagSet{pflag.NewFlagSet("test", pflag.ContinueOnError)}
		fs.StringP("accept-next-s", "n", "", "Accept next only")
		fs.Lookup("accept-next-s").ArgumentStyle = pflag.AcceptNext

		err := fs.Parse([]string{"-n=value"})
		if err == nil {
			t.Error("expected error for delimited style, got nil")
		}
	})

	// Test shorthand with AcceptDelimited
	t.Run("AcceptDelimited_shorthand_accepts_delimited", func(t *testing.T) {
		fs := &FlagSet{pflag.NewFlagSet("test", pflag.ContinueOnError)}
		fs.StringP("accept-delimited-s", "d", "", "Accept delimited only")
		fs.Lookup("accept-delimited-s").ArgumentStyle = pflag.AcceptDelimited

		err := fs.Parse([]string{"-d=value"})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// Test shorthand with AcceptNext accepts next
	t.Run("AcceptNext_shorthand_accepts_next", func(t *testing.T) {
		fs := &FlagSet{pflag.NewFlagSet("test", pflag.ContinueOnError)}
		fs.StringP("accept-next-s", "n", "", "Accept next only")
		fs.Lookup("accept-next-s").ArgumentStyle = pflag.AcceptNext

		err := fs.Parse([]string{"-n", "value"})
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	// Test AcceptNext|AcceptAttached accepts next but not delimited
	// Note: Attached style works for shorthand (-nvalue) not longhand (--flagvalue)
	t.Run("AcceptNextOrAttached_accepts_next", func(t *testing.T) {
		fs := &FlagSet{pflag.NewFlagSet("test", pflag.ContinueOnError)}
		fs.String("accept-next-or-attached", "", "Accept next or attached")
		fs.Lookup("accept-next-or-attached").ArgumentStyle = pflag.AcceptNext | pflag.AcceptAttached

		// Accepts next style
		err := fs.Parse([]string{"--accept-next-or-attached", "value"})
		if err != nil {
			t.Errorf("expected no error for next, got %v", err)
		}

		// Rejects delimited style
		fs = &FlagSet{pflag.NewFlagSet("test", pflag.ContinueOnError)}
		fs.String("accept-next-or-attached", "", "Accept next or attached")
		fs.Lookup("accept-next-or-attached").ArgumentStyle = pflag.AcceptNext | pflag.AcceptAttached
		err = fs.Parse([]string{"--accept-next-or-attached=value"})
		if err == nil {
			t.Error("expected error for delimited style, got nil")
		}
	})
}

func TestLookupNonPosixShorthandArgOptargNoDelim(t *testing.T) {
	_test := func(arg, name, prefix string, args ...string) {
		t.Run(arg, func(t *testing.T) {
			if args == nil {
				args = []string{}
			}

			fs := &FlagSet{pflag.NewFlagSet("test", pflag.PanicOnError)}
			fs.StringS("r", "r", "", "recurse")
			fs.StringS("ai", "ai", "", "include archive filenames")
			fs.Lookup("r").NoOptDefVal = " "
			fs.Lookup("r").OptargDelimiter = -1
			fs.Lookup("ai").NoOptDefVal = " "
			fs.Lookup("ai").OptargDelimiter = -1

			f := fs.lookupNonPosixShorthandArg(arg)
			if f == nil || f.Name != name {
				t.Fatal("should be " + name)
			}

			if f.ArgPrefix != prefix {
				t.Fatalf("prefix doesnt match actual: %#v, expected: %#v", f.ArgPrefix, prefix)
			}

			if !reflect.DeepEqual(f.Args, args) {
				t.Fatalf("args dont match %v: actual: %#v expected: %#v", arg, f.Args, args)
			}
		})
	}

	// single-char shorthand with attached value
	_test("-r-", "r", "-r", "-")
	_test("-r0", "r", "-r", "0")
	_test("-r", "r", "-r")

	// multi-char shorthand with attached value
	_test("-aifoo", "ai", "-ai", "foo")
	_test("-ai", "ai", "-ai")
}
