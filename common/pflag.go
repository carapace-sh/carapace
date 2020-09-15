package common

import (
	"reflect"

	"github.com/spf13/pflag"
)

// uses reflection to check for pflag.Flag.Shorthandonly to support both spf13/pflag
// and cornfeedhobos shorthand change (needed for carapace-bin)
// won't be necessary if https://github.com/spf13/pflag/pull/256 should ever be merged
func IsShorthandOnly(flag *pflag.Flag) bool {
	ValueIface := reflect.ValueOf(flag)
	Field := ValueIface.Elem().FieldByName("ShorthandOnly")
	if !Field.IsValid() {
		return false
	} else {
		return Field.Bool()
	}
}
