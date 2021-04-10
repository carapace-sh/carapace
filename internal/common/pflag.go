package common

import (
	"reflect"

	"github.com/spf13/pflag"
)

// IsShorthandOnly uses reflection to check for pflag.Flag.Shorthandonly to support both spf13/pflag
// and cornfeedhobos shorthand change (needed for carapace-bin)
// won't be necessary if https://github.com/spf13/pflag/pull/256 should ever be merged
func IsShorthandOnly(flag *pflag.Flag) (b bool) {
	if field := reflect.ValueOf(flag).Elem().FieldByName("ShorthandOnly"); field.IsValid() {
		b = field.Bool()
	}
	return
}
