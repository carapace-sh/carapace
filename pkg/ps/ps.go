package ps

import (
	"os"
	"strings"

	"github.com/mitchellh/go-ps"
)

func DetermineShell() string {
	process, err := ps.FindProcess(os.Getpid())
	for {
		if process, err = ps.FindProcess(process.PPid()); err != nil || process == nil {
			return ""
		}

		switch strings.SplitN(strings.TrimSuffix(process.Executable(), ".exe"), "-", 2)[0] {
		case "bash":
			return "bash"
		case "elvish":
			return "elvish"
		case "fish":
			return "fish"
		case "ion":
			return "ion"
		case "nu":
			return "nushell"
		case "oil":
			return "oil"
		case "osh":
			return "oil"
		case "powershell":
			return "powershell"
		case "pwsh":
			return "powershell"
		case "tcsh":
			return "tcsh"
		case "xonsh":
			return "xonsh"
		case "zsh":
			return "zsh"
		}
	}
}
