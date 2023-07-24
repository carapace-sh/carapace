package pflagfork

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/spf13/pflag"
)

type FlagSet struct {
	*pflag.FlagSet
}

func (f FlagSet) IsInterspersed() bool {
	if fv := reflect.ValueOf(f.FlagSet).Elem().FieldByName("interspersed"); fv.IsValid() {
		return fv.Bool()
	}
	return false
}

func (f FlagSet) IsPosix() bool {
	if method := reflect.ValueOf(f.FlagSet).MethodByName("IsPosix"); method.IsValid() {
		if values := method.Call([]reflect.Value{}); len(values) == 1 && values[0].Kind() == reflect.Bool {
			return values[0].Bool()
		}
	}
	return true
}

func (f FlagSet) IsShorthandSeries(arg string) bool {
	re := regexp.MustCompile("^-(?P<shorthand>[^-=]+)")
	return re.MatchString(arg) && f.IsPosix()
}

func (f FlagSet) IsMutuallyExclusive(flag *pflag.Flag) bool {
	if groups, ok := flag.Annotations["cobra_annotation_mutually_exclusive"]; ok {
		for _, group := range groups {
			for _, name := range strings.Split(group, " ") {
				if other := f.Lookup(name); other != nil && other.Changed {
					return true
				}
			}
		}
	}
	return false
}

func (f *FlagSet) VisitAll(fn func(*Flag)) {
	f.FlagSet.VisitAll(func(flag *pflag.Flag) {
		fn(&Flag{Flag: flag, Args: []string{}})
	})

}

func (fs FlagSet) LookupArg(arg string) (result *Flag) {
	isPosix := fs.IsPosix()

	if isPosix && !strings.HasPrefix(arg, "--") {
		return fs.lookupPosixShorthandArg(arg)
	}

	fs.VisitAll(func(f *Flag) {
		if result != nil {
			return
		}

		if f.Matches(arg, isPosix) {
			result = f
		}
	})
	return
}

func (fs FlagSet) ShorthandLookup(name string) *Flag {
	if f := fs.FlagSet.ShorthandLookup(name); f != nil {
		return &Flag{
			Flag: f,
			Args: []string{},
		}
	}
	return nil
}

func (fs FlagSet) lookupPosixShorthandArg(arg string) *Flag {
	if !strings.HasPrefix(arg, "-") || !fs.IsPosix() || len(arg) < 2 {
		return nil
	}

	for index, r := range arg[1:] {
		index += 1
		flag := fs.ShorthandLookup(string(r))

		switch {
		case flag == nil || len(arg) == index+1:
			flag.Prefix = arg
			return flag
		case arg[index+1] == byte(flag.OptargDelimiter()) && len(arg) > index+2:
			flag.Prefix = arg[:index+2]
			flag.Args = []string{arg[index+2:]}
			return flag
		case arg[index+1] == byte(flag.OptargDelimiter()):
			flag.Prefix = arg[:index+2]
			flag.Args = []string{""}
			return flag
		case !flag.IsOptarg() && len(arg) > index+1:
			flag.Prefix = arg[:index+1]
			flag.Args = []string{arg[index+1:]}
			return flag
		}
	}
	return nil
}
