package cmd

import (
	"os"
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/internal/assert"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
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

func TestRoot(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace/example")(func(s *sandbox.Sandbox) {
		s.Run("").
			Expect(carapace.Batch(
				carapace.ActionValuesDescribed(
					"action", "action example",
					"alias", "action example",
					"flag", "flag example",
				).Style(style.Blue).Tag("main commands"),
				carapace.ActionValuesDescribed(
					"modifier", "modifier example",
				).Style(style.Yellow).Tag("modifier commands"),
				carapace.ActionValuesDescribed(
					"injection", "just trying to break things",
				).Style(style.Magenta).Tag("test commands"),
				carapace.ActionValuesDescribed(
					"completion", "Generate the autocompletion script for the specified shell",
					"help", "Help about any command",
					"multiparts", "multiparts example",
					"special", "",
				).Tag("other commands"),
			).ToA())

		s.Run("a").
			Expect(carapace.ActionStyledValuesDescribed(
				"action", "action example", style.Blue,
				"alias", "action example", style.Blue,
			).Tag("main commands"))

		s.Run("action").
			Expect(carapace.ActionStyledValuesDescribed(
				"action", "action example", style.Blue,
			).Tag("main commands"))

		s.Run("-").
			Expect(carapace.ActionStyledValuesDescribed(
				"--array", "multiflag", style.Blue,
				"-a", "multiflag", style.Blue,
				"--persistentFlag", "Help message for persistentFlag", style.Yellow,
				"--persistentFlag2", "Help message for persistentFlag2", style.Blue,
				"-p", "Help message for persistentFlag", style.Yellow,
				"--toggle", "Help message for toggle", style.Default,
				"-t", "Help message for toggle", style.Default,
			).NoSpace('.').Tag("flags"))

		s.Run("--").
			Expect(carapace.ActionStyledValuesDescribed(
				"--array", "multiflag", style.Blue,
				"--persistentFlag", "Help message for persistentFlag", style.Yellow,
				"--persistentFlag2", "Help message for persistentFlag2", style.Blue,
				"--toggle", "Help message for toggle", style.Default,
			).NoSpace('.').Tag("flags"))

		s.Run("--a").
			Expect(carapace.ActionStyledValuesDescribed(
				"--array", "multiflag", style.Blue,
			).NoSpace('.').Tag("flags"))

		s.Run("--array").
			Expect(carapace.ActionStyledValuesDescribed(
				"--array", "multiflag", style.Blue,
			).NoSpace('.').Tag("flags"))

		s.Run("--array", "", "--a").
			Expect(carapace.ActionStyledValuesDescribed(
				"--array", "multiflag", style.Blue,
			).NoSpace('.').Tag("flags"))

		s.Run("-a", "", "--a").
			Expect(carapace.ActionStyledValuesDescribed(
				"--array", "multiflag", style.Blue,
			).NoSpace('.').Tag("flags"))

		s.Run("--toggle=").
			Expect(carapace.ActionStyledValues(
				"false", style.Red,
				"true", style.Green,
			).Prefix("--toggle="))
	})
}
