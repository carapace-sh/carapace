package carapace

import (
	"os"
	"path/filepath"
	"strings"
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
		return strings.TrimSuffix(filepath.Dir(wd("")), "/") + "/" + s
	}
	return strings.TrimSuffix(filepath.Dir(wd("")), "/") + "/"
}

func TestContextAbs(t *testing.T) {
	testCases := []struct {
		Dir      string
		Path     string
		Expected string
	}{
		{"/", "file", "/file"},
		{"", "file", wd("file")},
		{"", "../", parent("")},
		{"", "../file", parent("file")},
		{"", "~/file", home("file")},
		{"/", "~/file", home("file")},
		{"/", "file", "/file"},
		{"/dir", "file", "/dir/file"},
		{"/dir", "./.file", "/dir/.file"},
		{"", "/dir/", "/dir/"},
		{"/dir/", "", "/dir/"},
		{"~/", "file", home("file")},
		{"", "/", "/"},
		{"", ".hidden", wd(".hidden")},
		{"", "./", wd("") + "/"},
		{"", "", wd("") + "/"},
		{"", ".", wd("") + "/" + "."},
		{"", "~", home("")},
		{"", "~/file", home("file")},
	}

	for _, tc := range testCases {
		actual, err := Context{Dir: tc.Dir}.Abs(tc.Path)
		if err != nil {
			t.Error(err.Error())
		}
		if tc.Expected != actual {
			t.Errorf("context: '%v' arg: '%v' expected: '%v' was: '%v'", tc.Dir, tc.Path, tc.Expected, actual)
		}
	}
}

func TestEnv(t *testing.T) {
	c := Context{}
	if c.Getenv("example") != "" {
		t.Fail()
	}
	if v, exist := c.LookupEnv("example"); v != "" || exist {
		t.Fail()
	}

	c.Setenv("example", "value")
	if c.Getenv("example") != "value" {
		t.Fail()
	}
	if v, exist := c.LookupEnv("example"); v != "value" || !exist {
		t.Fail()
	}

	c.Setenv("example", "newvalue")
	if c.Getenv("example") != "newvalue" {
		t.Fail()
	}
	if v, exist := c.LookupEnv("example"); v != "newvalue" || !exist {
		t.Fail()
	}
}

func TestEnvsubst(t *testing.T) {
	c := Context{}
	c.Setenv("REPLACE", "me")
	for s, expected := range map[string]string{
		"BEFORE${REPLACE}AFTER":  "BEFOREmeAFTER",
		"BEFORE${!REPLACE}AFTER": "BEFORE${REPLACE}AFTER",

		"BEFORE${REPLACE:-default}AFTER":  "BEFOREmeAFTER",
		"BEFORE${!REPLACE:-default}AFTER": "BEFORE${REPLACE:-default}AFTER",

		"BEFORE${REPLACE/me/you}AFTER":  "BEFOREyouAFTER",
		"BEFORE${!REPLACE/me/you}AFTER": "BEFORE${REPLACE/me/you}AFTER",

		// TODO support curly brackets
		// "BEFORE${REPLACE/me/with\\}curly}AFTER":  "BEFOREwith}curlyAFTER",
		// "BEFORE${!REPLACE/me/with\\}curly}AFTER": "BEFORE${REPLACE/me/with\\}curly}AFTER",

		"BEFORE${UNSET:-default}AFTER":  "BEFOREdefaultAFTER",
		"BEFORE${!UNSET:-default}AFTER": "BEFORE${UNSET:-default}AFTER",
	} {
		t.Run(s, func(t *testing.T) {
			actual, err := c.Envsubst(s)
			if err != nil {
				t.Fatal(err)
			}
			if actual != expected {
				t.Fatalf("invalid replacement\nexpected: %#v\nactual  : %#v", expected, actual)
			}
		})
	}
}
