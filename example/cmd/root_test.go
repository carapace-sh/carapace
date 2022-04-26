package cmd

import (
	"os"
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/internal/assert"
)

func testScript(t *testing.T, shell string, file string) {
	if content, err := os.ReadFile(file); err != nil {
		t.Fatal("failed to read fixture file")
	} else {
		rootCmd.InitDefaultHelpCmd()
		s, _ := carapace.Gen(rootCmd).Snippet(shell)
		assert.Equal(t, string(content), s+"\n")
	}
}

func TestBash(t *testing.T) {
	testScript(t, "bash", "./_test/bash.sh")
}

func TestBashBle(t *testing.T) {
	testScript(t, "bash-ble", "./_test/bash-ble.sh")
}

func TestElvish(t *testing.T) {
	testScript(t, "elvish", "./_test/elvish.elv")
}

func TestFish(t *testing.T) {
	testScript(t, "fish", "./_test/fish.fish")
}

func TestNushell(t *testing.T) {
	testScript(t, "nushell", "./_test/nushell.nu")
}

func TestOil(t *testing.T) {
	testScript(t, "oil", "./_test/oil.sh")
}

func TestPowershell(t *testing.T) {
	testScript(t, "powershell", "./_test/powershell.ps1")
}

func TestXonsh(t *testing.T) {
	testScript(t, "xonsh", "./_test/xonsh.py")
}

func TestZsh(t *testing.T) {
	testScript(t, "zsh", "./_test/zsh.sh")
}
