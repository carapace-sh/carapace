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

func TestLookPathAbsolute(t *testing.T) {
	c := Context{Dir: "/tmp"}
	path, err := c.LookPath("/bin/sh")
	if err != nil {
		t.Fatal(err)
	}
	if path != "/bin/sh" {
		t.Fatalf("expected /bin/sh, got %v", path)
	}
}

func TestLookPathFromPATH(t *testing.T) {
	c := Context{Dir: "/tmp"}
	c.Setenv("PATH", "/usr/local/bin:/usr/bin:/bin")
	_, err := c.LookPath("sh")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLookPathCustomPATH(t *testing.T) {
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "myapp")
	if err := os.WriteFile(binPath, []byte("#!/bin/sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	c := Context{Dir: "/tmp"}
	c.Setenv("PATH", tmpDir)
	path, err := c.LookPath("myapp")
	if err != nil {
		t.Fatal(err)
	}
	if path != binPath {
		t.Fatalf("expected %v, got %v", binPath, path)
	}
}

func TestLookPathNotFound(t *testing.T) {
	c := Context{Dir: "/tmp"}
	c.Setenv("PATH", "/nonexistent")
	_, err := c.LookPath("doesnotexist")
	if err == nil {
		t.Fatal("expected error for nonexistent executable")
	}
}

func TestLookPathEmptyPATH(t *testing.T) {
	c := Context{Dir: "/tmp"}
	c.Setenv("PATH", "")
	_, err := c.LookPath("sh")
	if err == nil {
		t.Fatal("expected error when PATH is empty")
	}
}

func TestLookPathRelativeWithSlash(t *testing.T) {
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "mybin")
	if err := os.WriteFile(binPath, []byte("#!/bin/sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	c := Context{Dir: tmpDir}
	path, err := c.LookPath("./mybin")
	if err != nil {
		t.Fatal(err)
	}
	if path != "./mybin" {
		t.Fatalf("expected ./mybin, got %v", path)
	}
}

func TestLookPathUsesContextEnv(t *testing.T) {
	c := Context{Dir: "/tmp"}
	c.Env = []string{"PATH=/custom/bin:/bin"}
	path, err := c.LookPath("sh")
	if err != nil {
		t.Fatal(err)
	}
	if path != "/bin/sh" {
		t.Fatalf("expected /bin/sh, got %v", path)
	}
}

func TestLookPathRespectsDir(t *testing.T) {
	tmpDir := t.TempDir()
	localBin := filepath.Join(tmpDir, "localbin")
	if err := os.WriteFile(localBin, []byte("#!/bin/sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	c := Context{Dir: tmpDir, Env: []string{"PATH=."}}
	path, err := c.LookPath("localbin")
	if err != nil {
		t.Fatal(err)
	}
	if path != localBin {
		t.Fatalf("expected %v, got %v", localBin, path)
	}
}
