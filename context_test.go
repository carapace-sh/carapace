package carapace

import (
	"os"
	"path/filepath"
	"testing"
)

func wd(s string) string {
	if wd, _ := os.Getwd(); s != "" {
		return wd + "/" + s
	} else {
		return wd
	}
}

func home(s string) string {
	if hd, _ := os.UserHomeDir(); s != "" {
		return hd + "/" + s
	} else {
		return hd
	}
}

func parent(s string) string {
	if s != "" {
		return filepath.Dir(wd("")) + "/" + s
	}
	return filepath.Dir(wd("")) + "/"
}

func TestContextAbs(t *testing.T) {
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
		"", "./", wd("")+"/",
		"", "", wd("")+"/",
		"", ".", wd("")+"/"+".",
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
