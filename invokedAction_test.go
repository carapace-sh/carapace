package carapace

import (
	"strings"
	"testing"

	"github.com/rsteube/carapace/pkg/style"
)

func TestToMultiParts(t *testing.T) {
	_test := func(value, expected string, delimiter ...string) {
		a := ActionStyledValuesDescribed(
			"A/a:1", "one", style.Green,
			"A/a:2", "two", style.BgBlue,
			"A/b:3", "three", style.Red,
			"B/b:4", "four", style.Magenta,
			"B/c:5:1", "fiftyone", style.Black,
			"B/c:5:2/z", "fiftytwo", style.Yellow,
			"B/c/", "withsuffix", style.Underlined,
			"C/d/1()2", "withbrackets", style.Yellow,
		)
		a = a.Invoke(Context{}).ToMultiPartsA(delimiter...)
		if actual := a.Invoke(Context{Value: value}).value("export", value); !strings.Contains(actual, expected) {
			t.Errorf("expected '%v' in '%v' for '%v'", expected, actual, value)
		}
	}

	_test("A/a:1", `{"value":"A/a:1","display":"1","description":"one","style":"green"}`, "/", ":")
	_test("A/a:1", `{"value":"A/a:1","display":"1","description":"one","style":"green"}`, ":", "/")
	_test("A/a:1", `{"value":"A/a:1","display":"a:1","description":"one","style":"green"}`, "/")
	_test("A", `{"value":"A/","display":"A/"}`, "/", ":")
	_test("A", `{"value":"A/","display":"A/"}`, "/")
	_test("", `{"value":"A/","display":"A/"}`, "/")
	_test("A/", `{"value":"A/a:","display":"a:"}`, "/", ":")
	_test("A/", `{"value":"A/a:1","display":"a:1","description":"one","style":"green"}`, "/")
	_test("B/c:5", `{"value":"B/c:5:2/","display":"c:5:2/"}`, "/")
	_test("B/c:5", `{"value":"B/c:5:","display":"5:"}`, "/", ":")
	_test("B/c:5", `{"value":"B/c:5:","display":"5:"}`, ":", "/")

	_test("C/d/1", `{"value":"C/d/1()2","display":"1()2","description":"withbrackets","style":"yellow"}`, "/")
}
