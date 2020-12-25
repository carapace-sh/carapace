package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func doComplete(t *testing.T, shell string, cmdline string, contained ...string) {
	t.Run(cmdline, func(t *testing.T) {
		t.Parallel()
		var cmd *exec.Cmd

		switch shell {
		case "bash":
			cmd = exec.Command("./_test/invoke_bash", cmdline)
		case "elvish":
			cmd = exec.Command("./_test/invoke_elvish", cmdline)
		case "fish":
			cmd = exec.Command("./_test/invoke_fish", cmdline)
		case "powershell":
			cmd = exec.Command("./_test/invoke_powershell", strings.Replace(cmdline, ",", "`,", -1))
		case "xonsh":
			cmd = exec.Command("./_test/invoke_xonsh", cmdline)
		case "zsh":
			cmd = exec.Command("./_test/invoke_zsh", cmdline)
		}

		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if output, err := cmd.Output(); err != nil {
			t.Fatal(err.Error() + "\n" + stderr.String())
		} else {
			o := string(output)
			for _, s := range contained {
				if !strings.Contains(o, s) {
					t.Fatalf("\n%v\nshould contain: %v", o, s)
				}
			}
		}
	})
}

var tests = map[string]string{
	`example action `:                         "positional1",
	`example action --`:                       "--values_described",
	`example action --optarg `:                "positional1",
	`example action --optarg`:                 "--optarg",
	`example action --optarg=`:                "blue",
	`example condition `:                      "ERR",
	`example condition --required `:           "valid",
	`example condition --required invalid `:   "ERR",
	`example condition --required valid `:     "fulfilled",
	`example callback `:                       "callback1",
	`example callback --callback `:            "cb1",
	`example multiparts `:                     "VALUE",
	`example multiparts VALUE=`:               "one",
	`example multiparts VALUE=one,`:           "DIRECTORY",
	`example multiparts VALUE=one,DIRECTORY=`: "/",
}

func TestBash(t *testing.T) {
	if err := exec.Command("bash", "--version").Run(); err != nil {
		t.Skip("skipping bash")
	}
	for cmdline, text := range tests {
		doComplete(t, "bash", cmdline, text)
	}
}

func TestElvish(t *testing.T) {
	if err := exec.Command("elvish", "--version").Run(); err != nil {
		t.Skip("skipping elvish")
	}
	for cmdline, text := range tests {
		doComplete(t, "elvish", cmdline, text)
	}
}

func TestFish(t *testing.T) {
	if err := exec.Command("fish", "--version").Run(); err != nil {
		t.Skip("skipping fish")
	}
	for cmdline, text := range tests {
		doComplete(t, "fish", cmdline, text)
	}
}

func TestXonsh(t *testing.T) {
	if err := exec.Command("xonsh", "--version").Run(); err != nil {
		t.Skip("skipping xonsh")
	}
	for cmdline, text := range tests {
		doComplete(t, "xonsh", cmdline, text)
	}
}

func TestPowershell(t *testing.T) {
	if err := exec.Command("pwsh", "--version").Run(); err != nil {
		t.Skip("skipping powershell")
	}
	for cmdline, text := range tests {
		doComplete(t, "powershell", cmdline, text)
	}
}

func TestZsh(t *testing.T) {
	if err := exec.Command("zsh", "--version").Run(); err != nil {
		t.Skip("skipping zsh")
	}
	for cmdline, text := range tests {
		doComplete(t, "zsh", cmdline, text)
	}
}
