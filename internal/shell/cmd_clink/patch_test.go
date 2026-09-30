package cmd_clink

import (
	"reflect"
	"testing"
)

func TestPatch(t *testing.T) {
	t.Setenv("CARAPACE_COMPLINE", `example C:\path\to "a b" a^b 'single'`)
	got, err := Patch([]string{"cmd-clink"})
	if err != nil {
		t.Fatal(err.Error())
	}
	expected := []string{"cmd-clink", "example", `C:\path\to`, "a b", "ab", "'single'"}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("got %#v, want %#v", got, expected)
	}
}
