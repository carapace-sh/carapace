package common

import (
	"bytes"
	"io"
	"os"
)

func CaptureStdout(f func() error) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	err := f()

	w.Close()
	out := <-outC
	os.Stdout = old

	return out, err
}
