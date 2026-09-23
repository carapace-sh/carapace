// Package ps provides shell determination by process name
package ps

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/carapace-sh/carapace/third_party/github.com/mitchellh/go-ps"
)

// DetermineShell determines shell by parent process name.
func DetermineShell() string {
	process, err := ps.FindProcess(os.Getpid())
	if err != nil {
		return ""
	}
	for {
		if process, err = ps.FindProcess(process.PPid()); err != nil || process == nil {
			return ""
		}

		executable := resolveExecutable(process.Pid(), process.Executable())
		switch strings.SplitN(strings.TrimSuffix(executable, ".exe"), "-", 2)[0] {
		case "bash":
			if isBLE() {
				return "bash-ble"
			}
			return "bash"
		case "cmd":
			return "cmd-clink"
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
		default:
			if strings.Contains(executable, "xonsh-wrapped") { // nix packaged version
				return "xonsh"
			}
		}
	}
}

// resolveExecutable returns the binary name for a process. On Termux (Android 10+),
// execve is rewritten by termux-exec to run binaries through the system linker
// (/system/bin/linker64) to bypass app data exec restrictions, so the kernel comm
// shows "linker"/"linker64" instead of the binary name. The real binary path is
// recoverable as the first entry of /proc/<pid>/cmdline.
func resolveExecutable(pid int, executable string) string {
	switch executable {
	case "linker", "linker64":
		data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
		if err != nil {
			return executable
		}
		if name := strings.SplitN(string(data), "\x00", 2)[0]; name != "" {
			return filepath.Base(name)
		}
	}
	return executable
}

func isBLE() bool {
	bleEnvs := []string{
		"_bleopt_connect_tty",
		"_ble_util_fdlist_cloexec",
		"_ble_util_fd_null",
		"_ble_util_fd_stderr",
		"_ble_util_fd_stdin",
		"_ble_util_fd_stdout",
		"_ble_util_fdvars_export",
		"_ble_util_fd_zero",
	}
	for _, e := range bleEnvs {
		if _, ok := os.LookupEnv(e); ok {
			return true
		}
	}
	return false
}
