// Package assert provides test helpers
package assert

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// Equal calls t.Error if given strings are not equal
func Equal(t *testing.T, expected string, actual string) {
	if expected != actual {
		// TODO handle err
		expectedFile, _ := os.CreateTemp(os.TempDir(), "carapace_test")
		actualFile, _ := os.CreateTemp(os.TempDir(), "carapace_test")

		_ = os.WriteFile(expectedFile.Name(), []byte(expected), os.ModePerm)
		_ = os.WriteFile(actualFile.Name(), []byte(actual), os.ModePerm)
		output, _ := exec.Command("diff", "--color=always", expectedFile.Name(), actualFile.Name()).Output()
		_, file, line, _ := runtime.Caller(2)
		t.Errorf("%v:%v:\n%v", filepath.Base(file), line, string(output))
	}
}
