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

	if s, err := c.Envsubst("start${example}end"); s != "startend" || err != nil {
		t.Fail()
	}

	if s, err := c.Envsubst("start${example:-default}end"); s != "startdefaultend" || err != nil {
		t.Fail()
	}

	c.Setenv("example", "value")
	if s, err := c.Envsubst("start${example}end"); s != "startvalueend" || err != nil {
		t.Fail()
	}

	if s, err := c.Envsubst("start${example:-default}end"); s != "startvalueend" || err != nil {
		t.Fail()
	}
}
