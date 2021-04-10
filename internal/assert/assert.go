package assert

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
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
		t.Error("\n" + string(output))
	}
}
