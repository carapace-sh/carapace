// Package ps provides shell determination by process name
package ps

import (
	"os"
	"strings"

	"github.com/mitchellh/go-ps"
)

// DetermineShell determines shell by parent process name
func DetermineShell() string {
	process, err := ps.FindProcess(os.Getpid())
	if err != nil {
		return ""
	}
	for {
		if process, err = ps.FindProcess(process.PPid()); err != nil || process == nil {
			return ""
		}

		switch strings.SplitN(strings.TrimSuffix(process.Executable(), ".exe"), "-", 2)[0] {
		case "bash":
			if _, ok := os.LookupEnv("BLE_ATTACHED"); ok { // TODO env is not exported
				return "bash-ble"
			}
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
