package zsh

import (
	"testing"
)

func TestActionCallback(t *testing.T) {
	a := ActionCallback(func(args []string) Action {
		return ActionMessage("ActionCallback test")
	}).finalize("someId")

	if a.Value != ` eval \$(${os_args[1]} _zsh_completion 'someId' ${os_args:1})` {
		t.Error(highlight(a.Value))
	}
}

func TestActionPathFiles(t *testing.T) {
	a := ActionPathFiles("").finalize("someId")

	if a.Value != `_path_files` {
		t.Error(highlight(a.Value))
	}
}

func TestActionPathFilesWithPattern(t *testing.T) {
	a := ActionPathFiles("*.go").finalize("someId")

	if a.Value != `_path_files -g '*.go'` {
		t.Error(highlight(a.Value))
	}
}
