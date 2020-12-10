package assert

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func Highlight(s string) string {
	cmd := exec.Command("bat", "-l", "bash", "--color", "always", "--wrap", "never")
	cmd.Stdin = strings.NewReader(s)
	if output, err := cmd.Output(); err != nil {
		return s // no bat command
	} else {
		return string(output)
	}
}

func Equal(t *testing.T, expected string, actual string) {
	if expected == actual {
		t.Log("\n" + Highlight(actual))
	} else {
		expectedFile, _ := ioutil.TempFile(os.TempDir(), "carapace_test")
		actualFile, _ := ioutil.TempFile(os.TempDir(), "carapace_test")

		ioutil.WriteFile(expectedFile.Name(), []byte(expected), os.ModePerm)
		ioutil.WriteFile(actualFile.Name(), []byte(actual), os.ModePerm)
		output, _ := exec.Command("diff", "--color=always", expectedFile.Name(), actualFile.Name()).Output()
		t.Error("\n" + string(output))
	}
}
