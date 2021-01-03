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
			replacer := strings.NewReplacer(
				`\ `, "' '", // TODO simple escape fix
			)
			cmd = exec.Command("./_test/invoke_elvish", replacer.Replace(cmdline))
		case "fish":
			cmd = exec.Command("./_test/invoke_fish", cmdline)
		case "oil":
			cmd = exec.Command("./_test/invoke_oil", cmdline)
		case "powershell":
			replacer := strings.NewReplacer(
				",", "`,",
				`\ `, "` ", // TODO simple escape fix
			)
			cmd = exec.Command("./_test/invoke_powershell", replacer.Replace(cmdline))
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
	`example -ap `:                                                 "action",
	`example action `:                                              "p",
	`example action positional`:                                    "positional1",
	`example action --`:                                            "--values_described",
	`example action --optarg `:                                     "p",
	`example action --optarg positional`:                           "positional1",
	`example action --optarg`:                                      "--optarg",
	`example action --optarg=`:                                     "blue",
	`example action -o `:                                           "p",
	`example action -o positional`:                                 "positional1",
	`example action -o`:                                            "-o",
	`example action -o=`:                                           "blue",
	`example action -fgo=`:                                         "blue",
	`example action -fgo= `:                                        "p",
	`example condition `:                                           "ERR",
	`example condition --required `:                                "valid",
	`example condition --required invalid `:                        "ERR",
	`example condition --required valid `:                          "fulfilled",
	`example callback `:                                            "callback",
	`example callback callback`:                                    "callback1",
	`example callback --callback `:                                 "cb",
	`example callback --callback cb`:                               "cb1",
	`example multiparts `:                                          "VALUE",
	`example multiparts --at `:                                     "first",
	`example multiparts --at first@`:                               "second",
	`example multiparts --at first@third\ with\ space@`:            "second",
	`example multiparts --colon `:                                  "first",
	`example multiparts --colon first:`:                            "second",
	`example multiparts --colon first:third\ with\ space:`:         "second",
	`example multiparts --comma `:                                  "first",
	`example multiparts --comma first,`:                            "second",
	`example multiparts --comma first,third\ with\ space,`:         "second", // TODO escape space correctly for each shell
	`example multiparts --dot `:                                    "first",
	`example multiparts --dot first.`:                              "second",
	`example multiparts --dot first.third\ with\ space.`:           "second",
	`example multiparts --dotdotdot `:                              "first",
	`example multiparts --dotdotdot first...`:                      "second",
	`example multiparts --dotdotdot first...third\ with\ space...`: "second",
	`example multiparts --equals `:                                 "first",
	`example multiparts --equals first=`:                           "second",
	`example multiparts --equals first=third\ with\ space=`:        "second",
	`example multiparts --slash `:                                  "first",
	`example multiparts --slash first/`:                            "second",
	`example multiparts --slash first/third\ with\ space/`:         "second",
	`example multiparts --none `:                                   "a",
	`example multiparts --none a`:                                  "b",
	`example multiparts --none ab`:                                 "c",
	`example multiparts VALUE=`:                                    "one",
	`example multiparts VALUE=one,`:                                "DIRECTORY",
	`example multiparts VALUE=one,DIRECTORY=`:                      "/",
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

func TestOil(t *testing.T) {
	if err := exec.Command("oil", "--version").Run(); err != nil {
		t.Skip("skipping oil")
	}
	for cmdline, text := range tests {
		doComplete(t, "oil", cmdline, text)
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
