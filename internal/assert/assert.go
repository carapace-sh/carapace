// Package assert provides test helpers
package assert

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type T interface {
	Cleanup(func())
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Helper()
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Setenv(key, value string)
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	TempDir() string
}

// Equal calls t.Error if given strings are not equal.
func Equal(t T, expected string, actual string) {
	if expected != actual {
		// TODO handle err
		tempDir := t.TempDir()
		expctdFile := fmt.Sprintf("%v/expctd", tempDir)
		actualFile := fmt.Sprintf("%v/actual", tempDir)

		if err := os.WriteFile(expctdFile, []byte(expected), os.ModePerm); err != nil {
			t.Fatal(err.Error())
		}

		if err := os.WriteFile(actualFile, []byte(actual), os.ModePerm); err != nil {
			t.Fatal(err.Error())
		}
		output, _ := exec.Command("diff", "--color=always", "--context=5", expctdFile, actualFile).Output()
		_, file, line, _ := runtime.Caller(2)
		t.Errorf("%v:%v:\n%v", filepath.Base(file), line, string(output))
	}
}
