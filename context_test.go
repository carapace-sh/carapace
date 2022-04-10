package carapace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestContextAbs(t *testing.T) {
	wd := func(s string) string {
		if true {
			return "/home/rsteube/Documents/development/github/carapace/" + s
		}

		if wd, _ := os.Getwd(); s != "" {
			return wd + "/" + s
		} else {
			return wd
		}
	}

	home := func(s string) string {
		if true {
			return "/home/rsteube/" + s
		}

		if hd, _ := os.UserHomeDir(); s != "" {
			return hd + "/" + s
		} else {
			return hd
		}
	}

	parent := func(s string) string {
		if true {
			return "/home/rsteube/Documents/development/github/" + s
		}

		if s != "" {
			return filepath.Dir(wd("")) + "/" + s
		}
		return filepath.Dir(wd("")) + "/"
	}

	tests := append([]string{},
		"/", "file", "/file",
		"", "file", wd("file"),
		"", "../", parent(""),
		"", "../file", parent("file"),
		"", "~/file", home("file"),
		"/", "~/file", home("file"),
		"/", "file", "/file",
		"/dir", "file", "/dir/file",
		"/dir", "./.file", "/dir/.file",
		"", "/dir/", "/dir/",
		"/dir/", "", "/dir/",
		"~", "file", home("file"),
		"", "/", "/",
		"", ".hidden", wd(".hidden"),
		"", "./", wd(""),
		"", "", wd(""),
		"", ".", wd("")+".",
	)

	for index := 0; index < len(tests); index += 3 {
		actual, err := Context{Dir: tests[index]}.Abs(tests[index+1])
		if err != nil {
			t.Error(err.Error())
		}
		if expected := tests[index+2]; expected != actual {
			t.Errorf("context: '%v' arg: '%v' expected: '%v' was: '%v'", tests[index], tests[index+1], expected, actual)
		}
	}

}
