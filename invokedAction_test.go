package carapace

import (
	"strings"
	"testing"

	"github.com/rsteube/carapace/pkg/style"
)

func TestToMultiParts(t *testing.T) {
	_test := func(cv, expected string, delimiter ...string) {
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
		if actual := a.Invoke(Context{CallbackValue: cv}).value("export", cv); !strings.Contains(actual, expected) {
			t.Errorf("expected '%v' in '%v' for '%v'", expected, actual, cv)
		}
	}

	_test("A/a:1", `{"Value":"A/a:1","Display":"1","Description":"one","Style":"green"}`, "/", ":")
	_test("A/a:1", `{"Value":"A/a:1","Display":"1","Description":"one","Style":"green"}`, ":", "/")
	_test("A/a:1", `{"Value":"A/a:1","Display":"a:1","Description":"one","Style":"green"}`, "/")
	_test("A", `{"Value":"A/","Display":"A/"}`, "/", ":")
	_test("A", `{"Value":"A/","Display":"A/"}`, "/")
	_test("", `{"Value":"A/","Display":"A/"}`, "/")
	_test("A/", `{"Value":"A/a:","Display":"a:"}`, "/", ":")
	_test("A/", `{"Value":"A/a:1","Display":"a:1","Description":"one","Style":"green"}`, "/")
	_test("B/", `{"Value":"B/c/","Display":"c/","Description":"withsuffix","Style":"underlined"}`, "/")
	_test("B/c:5", `{"Value":"B/c:5:2/","Display":"c:5:2/"}`, "/")
	_test("B/c:5", `{"Value":"B/c:5:","Display":"5:"}`, "/", ":")
	_test("B/c:5", `{"Value":"B/c:5:","Display":"5:"}`, ":", "/")

	_test("C/d/1", `{"Value":"C/d/1()2","Display":"1()2","Description":"withbrackets","Style":"yellow"}`, "/")
}
