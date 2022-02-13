// Package assert provides test helpers
package assert

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	exec "golang.org/x/sys/execabs"
)

// Equal calls t.Error if given strings are not equal
func Equal(t *testing.T, expected string, actual string) {
	if expected != actual {
		// TODO handle err
		expectedFile, _ := ioutil.TempFile(os.TempDir(), "carapace_test")
		actualFile, _ := ioutil.TempFile(os.TempDir(), "carapace_test")

		_ = ioutil.WriteFile(expectedFile.Name(), []byte(expected), os.ModePerm)
		_ = ioutil.WriteFile(actualFile.Name(), []byte(actual), os.ModePerm)
		output, _ := exec.Command("diff", "--color=always", expectedFile.Name(), actualFile.Name()).Output()
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("%v:%v:\n%v", filepath.Base(file), line, string(output))
	}
}
