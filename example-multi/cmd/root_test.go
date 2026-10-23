package cmd

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/assert"
	"github.com/carapace-sh/carapace/pkg/sandbox"
	"github.com/carapace-sh/carapace/pkg/style"
)

func testScript(t *testing.T, shell string, file string) {
	if content, err := os.ReadFile(file); err != nil {
		t.Fatal("failed to read fixture file")
	} else {
		rootCmd.InitDefaultHelpCmd()
		s, _ := carapace.Gen(rootCmd).Snippet(shell)
		expected := strings.ReplaceAll(string(content), "\r\n", "\n") // autocrlf on windows checkouts
		if runtime.GOOS == "windows" && expected != s+"\n" {
			t.Skipf("fixture does not cover windows-specific snippet: %v", shell)
		}
		assert.Equal(t, expected, s+"\n")
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

func TestTcsh(t *testing.T) {
	testScript(t, "tcsh", "./_test/tcsh.sh")
}

func TestXonsh(t *testing.T) {
	testScript(t, "xonsh", "./_test/xonsh.py")
}

func TestZsh(t *testing.T) {
	testScript(t, "zsh", "./_test/zsh.sh")
}

func TestIdentify(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-multi")(func(s *sandbox.Sandbox) {
		s.Run("identify", "-").
			Expect(carapace.Batch(
				carapace.ActionStyledValuesDescribed(
					"-f", "image format", style.Blue,
					"-h", "help for identify", style.Default,
					"-v", "verbose output", style.Default,
				).Tag("shorthand flags"),
				carapace.ActionStyledValuesDescribed(
					"--format", "image format", style.Blue,
					"--help", "help for identify", style.Default,
					"--verbose", "verbose output", style.Default,
				).Tag("longhand flags"),
			).ToA().NoSpace('.'))

		s.Run("identify", "--format", "").
			Expect(carapace.ActionValues("png", "jpeg", "gif", "tiff", "bmp").
				Usage("image format"))
	})
}

func TestConvert(t *testing.T) {
	sandbox.Package(t, "github.com/carapace-sh/carapace/example-multi")(func(s *sandbox.Sandbox) {
		s.Run("convert", "-").
			Expect(carapace.Batch(
				carapace.ActionStyledValuesDescribed(
					"-o", "output format", style.Blue,
					"-q", "output quality", style.Blue,
					"-h", "help for convert", style.Default,
				).Tag("shorthand flags"),
				carapace.ActionStyledValuesDescribed(
					"--output", "output format", style.Blue,
					"--quality", "output quality", style.Blue,
					"--help", "help for convert", style.Default,
				).Tag("longhand flags"),
			).ToA().NoSpace('.'))

		s.Run("convert", "--output", "").
			Expect(carapace.ActionValues("png", "jpeg", "gif", "tiff", "bmp").
				Usage("output format"))

		s.Run("convert", "--quality", "").
			Expect(carapace.ActionValues("1", "50", "75", "90", "100").
				Usage("output quality"))
	})
}
