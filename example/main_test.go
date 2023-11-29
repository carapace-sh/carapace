package main

import (
	"testing"

	"github.com/rsteube/carapace"
)

// func doComplete(t *testing.T, shell string, cmdline string, contained ...string) {
// 	t.Run(cmdline, func(t *testing.T) {
// 		t.Parallel()
// 		var cmd *exec.Cmd

// 		switch shell {
// 		case "bash":
// 			cmd = exec.Command("./_test/invoke_bash", cmdline)
// 		case "elvish":
// 			replacer := strings.NewReplacer(
// 				`\ `, "' '", // TODO simple escape fix
// 			)
// 			cmd = exec.Command("./_test/invoke_elvish", replacer.Replace(cmdline))
// 		case "fish":
// 			cmd = exec.Command("./_test/invoke_fish", cmdline)
// 		case "oil":
// 			cmd = exec.Command("./_test/invoke_oil", cmdline)
// 		case "powershell":
// 			replacer := strings.NewReplacer(
// 				",", "`,",
// 				`\ `, "` ", // TODO simple escape fix
// 			)
// 			cmd = exec.Command("./_test/invoke_powershell", replacer.Replace(cmdline))
// 		case "xonsh":
// 			replacer := strings.NewReplacer(
// 				`\ `, `" "`, // TODO simple escape fix
// 			)
// 			cmd = exec.Command("./_test/invoke_xonsh", replacer.Replace(cmdline))
// 		case "zsh":
// 			cmd = exec.Command("./_test/invoke_zsh", cmdline)
// 		}

// 		var stderr bytes.Buffer
// 		cmd.Stderr = &stderr
// 		cmd.Env = append(os.Environ(), "HOME=/tmp/carapace-fakehome")
// 		if output, err := cmd.Output(); err != nil {
// 			t.Fatal(err.Error() + "\n" + stderr.String())
// 		} else {
// 			o := string(output)
// 			for _, s := range contained {
// 				if !strings.Contains(o, s) {
// 					t.Fatalf("\n%v\nshould contain: %v", o, s)
// 				}
// 			}
// 		}
// 	})
// }

// var tests = map[string]string{
// 	`example -ap `:               "action",
// 	`example action `:            "p",
// 	`example action positional`:  "positional1",
// 	`example action positional1`: "positional1 with space",
// 	//`example action "positional1 `: "positional1 with space", // TODO this test does not yet work with bash as it's missing quote handling in the snippet
// 	//`example action --`:                                            "--values_described", // weird: causes regex match in expect/xonsh not to work
// 	//`example action -`:                                             "-o", // weird: causes regex match in expect/xonsh not to work
// 	`example special --optarg `:                                    "p",
// 	`example special --optarg positional`:                          "positional1",
// 	`example special --optar`:                                      "--optarg",
// 	`example special --optarg=`:                                    "optarg",
// 	`example special -o`:                                           "count flag",
// 	`example special -oc`:                                          "count flag",
// 	`example special -o `:                                          "p",
// 	`example special -o pos`:                                       "positional",
// 	`example multiparts `:                                          "VALUE",
// 	`example multiparts -`:                                         "-c",
// 	`example multiparts --`:                                        "--comma",
// 	`example multiparts --at `:                                     "first",
// 	`example multiparts --at first@`:                               "second",
// 	`example multiparts --at first@third\ with\ space@`:            "second",
// 	`example multiparts --colon `:                                  "first",
// 	`example multiparts --colon first:`:                            "second",
// 	`example multiparts --colon first:third\ with\ space:`:         "second",
// 	`example multiparts --comma `:                                  "first",
// 	`example multiparts --comma first,`:                            "second",
// 	`example multiparts --comma first,third\ with\ space,`:         "second", // TODO escape space correctly for each shell
// 	`example multiparts --dot `:                                    "first",
// 	`example multiparts --dot first.`:                              "second",
// 	`example multiparts --dot first.third\ with\ space.`:           "second",
// 	`example multiparts --dotdotdot `:                              "first",
// 	`example multiparts --dotdotdot first...`:                      "second",
// 	`example multiparts --dotdotdot first...third\ with\ space...`: "second",
// 	`example multiparts --equals `:                                 "first",
// 	`example multiparts --equals first=`:                           "second",
// 	`example multiparts --equals first=third\ with\ space=`:        "second",
// 	`example multiparts --slash `:                                  "first",
// 	`example multiparts --slash first/`:                            "second",
// 	`example multiparts --slash first/third\ with\ space/`:         "second",
// 	`example multiparts --none `:                                   "a",
// 	`example multiparts --none a`:                                  "b",
// 	`example multiparts --none ab`:                                 "c",
// 	`example multiparts VALUE=`:                                    "one",
// 	`example multiparts VALUE=one,`:                                "DIRECTORY",
// 	`example multiparts VALUE=one,DIRECTORY=`:                      "/",

// 	// `example flag --Bool `:           "",
// 	// `example flag --BoolSlice `:      "",
// 	`example flag --BytesBase64 `: "M",
// 	`example flag --BytesHex `:    "0",
// 	// `example flag --Count `:          "",
// 	`example flag --Duration `:       "1h",
// 	`example flag --DurationSlice `:  "1h",
// 	`example flag --Float32P `:       "1",
// 	`example flag --Float32Slice `:   "1",
// 	`example flag --Float64P `:       "1",
// 	`example flag --Float64Slice `:   "1",
// 	`example flag --Int16 `:          "1",
// 	`example flag --Int32 `:          "1",
// 	`example flag --Int32Slice `:     "1",
// 	`example flag --Int64 `:          "1",
// 	`example flag --Int64Slice `:     "1",
// 	`example flag --Int8 `:           "1",
// 	`example flag --Int `:            "1",
// 	`example flag --IntSlice `:       "1",
// 	`example flag --IPMask `:         "0",
// 	`example flag --IP `:             "0",
// 	`example flag --IPNet `:          "0",
// 	`example flag --IPSlice `:        "0",
// 	`example flag --StringArray `:    "1",
// 	`example flag --String `:         "1",
// 	`example flag --StringSlice `:    "1",
// 	`example flag --StringToInt64 `:  "a",
// 	`example flag --StringToInt `:    "a",
// 	`example flag --StringToString `: "a",
// 	`example flag --Uint16 `:         "1",
// 	`example flag --Uint32 `:         "1",
// 	`example flag --Uint64 `:         "1",
// 	`example flag --Uint8 `:          "1",
// 	`example flag --Uint `:           "1",
// 	`example flag --UintSlice `:      "1",

// 	// `example flag --optarg --Bool=`:           "",
// 	// `example flag --optarg --BoolSlice=`:      "",
// 	`example flag --optarg --BytesBase64=`: "M",
// 	`example flag --optarg --BytesHex=`:    "0",
// 	// `example flag --optarg --Count=`:          "",
// 	`example flag --optarg --Duration=`:       "1h",
// 	`example flag --optarg --DurationSlice=`:  "1h",
// 	`example flag --optarg --Float32P=`:       "1",
// 	`example flag --optarg --Float32Slice=`:   "1",
// 	`example flag --optarg --Float64P=`:       "1",
// 	`example flag --optarg --Float64Slice=`:   "1",
// 	`example flag --optarg --Int16=`:          "1",
// 	`example flag --optarg --Int32=`:          "1",
// 	`example flag --optarg --Int32Slice=`:     "1",
// 	`example flag --optarg --Int64=`:          "1",
// 	`example flag --optarg --Int64Slice=`:     "1",
// 	`example flag --optarg --Int8=`:           "1",
// 	`example flag --optarg --Int=`:            "1",
// 	`example flag --optarg --IntSlice=`:       "1",
// 	`example flag --optarg --IPMask=`:         "0",
// 	`example flag --optarg --IP=`:             "0",
// 	`example flag --optarg --IPNet=`:          "0",
// 	`example flag --optarg --IPSlice=`:        "0",
// 	`example flag --optarg --StringArray=`:    "1",
// 	`example flag --optarg --String=`:         "1",
// 	`example flag --optarg --StringSlice=`:    "1",
// 	`example flag --optarg --StringToInt64=`:  "a",
// 	`example flag --optarg --StringToInt=`:    "a",
// 	`example flag --optarg --StringToString=`: "a",
// 	`example flag --optarg --Uint16=`:         "1",
// 	`example flag --optarg --Uint32=`:         "1",
// 	`example flag --optarg --Uint64=`:         "1",
// 	`example flag --optarg --Uint8=`:          "1",
// 	`example flag --optarg --Uint=`:           "1",
// 	`example flag --optarg --UintSlice=`:      "1",
// }

// var testsIntegratedMessage = map[string]string{
// 	`example action -z`:          "unknown",
// 	`example action --fail `:     "unknown",
// 	`example --fail acti`:        "unknown",
// 	`example -z acti`:            "unknown",
// 	`example flag -o=`:           "unknown", // seems shorthand flag should not accept optional arguments and `=` is seen as another flag
// 	`example action --callback `: "values flag is not set",
// }

// func TestBash(t *testing.T) {
// 	if err := exec.Command("bash", "--version").Run(); err != nil {
// 		t.Skip("skipping bash")
// 	}
// 	for cmdline, text := range tests {
// 		doComplete(t, "bash", cmdline, text)
// 	}
// 	for cmdline, text := range testsIntegratedMessage {
// 		doComplete(t, "bash", cmdline, text)
// 	}
// }

// func TestElvish(t *testing.T) {
// 	if err := exec.Command("elvish", "--version").Run(); err != nil {
// 		t.Skip("skipping elvish")
// 	}
// 	for cmdline, text := range tests {
// 		doComplete(t, "elvish", cmdline, text)
// 	}
// }

// func TestFish(t *testing.T) {
// 	if err := exec.Command("fish", "--version").Run(); err != nil {
// 		t.Skip("skipping fish")
// 	}
// 	for cmdline, text := range tests {
// 		doComplete(t, "fish", cmdline, text)
// 	}
// 	for cmdline, text := range testsIntegratedMessage {
// 		doComplete(t, "fish", cmdline, text)
// 	}
// }

// func TestXonsh(t *testing.T) {
// 	if err := exec.Command("xonsh", "--version").Run(); err != nil {
// 		t.Skip("skipping xonsh")
// 	}
// 	for cmdline, text := range tests {
// 		doComplete(t, "xonsh", cmdline, text)
// 	}
// 	for cmdline, text := range testsIntegratedMessage {
// 		doComplete(t, "xonsh", cmdline, text)
// 	}
// }

// func TestOil(t *testing.T) {
// 	if err := exec.Command("oil", "--version").Run(); err != nil {
// 		t.Skip("skipping oil")
// 	}
// 	for cmdline, text := range tests {
// 		doComplete(t, "oil", cmdline, text)
// 	}
// 	for cmdline, text := range testsIntegratedMessage {
// 		doComplete(t, "oil", cmdline, text)
// 	}
// }

//func TestPowershell(t *testing.T) {
//	if err := exec.Command("pwsh", "--version").Run(); err != nil {
//		t.Skip("skipping powershell")
//	}
//	for cmdline, text := range tests {
//		doComplete(t, "powershell", cmdline, text)
//	}
//	for cmdline, text := range testsIntegratedMessage {
//		doComplete(t, "powershell", cmdline, text)
//	}
//}

// func TestZsh(t *testing.T) {
// 	if err := exec.Command("zsh", "--version").Run(); err != nil {
// 		t.Skip("skipping zsh")
// 	}
// 	for cmdline, text := range tests {
// 		doComplete(t, "zsh", cmdline, text)
// 	}
// }

func TestCarapace(t *testing.T) {
	carapace.Test(t)
}
