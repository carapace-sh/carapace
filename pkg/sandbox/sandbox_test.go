package sandbox

import (
	"os"
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func TestPreInvoke(t *testing.T) {
	Command(t, func() *cobra.Command {
		rootCmd := &cobra.Command{}
		rootCmd.CompletionOptions.DisableDefaultCmd = true
		rootCmd.SetHelpCommand(nil)
		carapace.Gen(rootCmd).PreInvoke(func(cmd *cobra.Command, flag *pflag.Flag, action carapace.Action) carapace.Action {
			return action.Chdir("subdir1")
		})
		carapace.Gen(rootCmd).PositionalCompletion(
			carapace.ActionFiles(),
		)

		subCmd := &cobra.Command{
			Use: "sub",
		}
		carapace.Gen(subCmd).PreInvoke(func(cmd *cobra.Command, flag *pflag.Flag, action carapace.Action) carapace.Action {
			return action.Chdir("subdir2")
		})
		carapace.Gen(subCmd).PositionalCompletion(
			carapace.ActionFiles(),
		)
		rootCmd.AddCommand(subCmd)

		return rootCmd
	})(func(s *Sandbox) {
		s.Files(
			"subdir1/file1.txt", "",
			"subdir1/subdir2/file2.txt", "",
		)

		s.Run("").
			Expect(carapace.ActionValues(
				"file1.txt",
				"subdir2/").
				StyleF(style.ForPath).
				NoSpace('/').
				Tag("files").
				Chdir("subdir1"))

		s.Run("sub", "").
			Expect(carapace.ActionValues("file2.txt").
				StyleF(style.ForPath).
				NoSpace('/').
				Tag("files").
				Chdir("subdir1/subdir2/"))
	})
}

func TestPreRun(t *testing.T) {
	Command(t, func() *cobra.Command {
		rootCmd := &cobra.Command{}
		carapace.Gen(rootCmd).PreRun(func(cmd *cobra.Command, args []string) {
			cmd.Flags().Bool("root", false, "root flag")
		})

		subCmd := &cobra.Command{
			Use: "sub",
		}
		carapace.Gen(subCmd).PreRun(func(cmd *cobra.Command, args []string) {
			cmd.Flags().String("sub", "", "sub flag")

			carapace.Gen(cmd).FlagCompletion(carapace.ActionMap{
				"sub": carapace.ActionValues(cmd.Parent().Flag("root").Value.String()),
			})
		})
		rootCmd.AddCommand(subCmd)

		return rootCmd
	})(func(s *Sandbox) {
		s.Run("--root").
			Expect(carapace.ActionValuesDescribed(
				"--root", "root flag").
				NoSpace('.').
				Tag("flags"))

		s.Run("--root", "sub", "--").
			Expect(carapace.ActionStyledValuesDescribed(
				"--help", "help for sub", style.Default,
				"--sub", "sub flag", style.Blue).
				NoSpace('.').
				Tag("flags"))

		s.Run("--root", "sub", "--sub", "").
			Expect(carapace.ActionValues("true").
				Usage("sub flag"))

		s.Run("sub", "--sub", "").
			Expect(carapace.ActionValues("false").
				Usage("sub flag"))
	})
}

func TestEnv(t *testing.T) {
	Command(t, func() *cobra.Command {
		rootCmd := &cobra.Command{}
		rootCmd.CompletionOptions.DisableDefaultCmd = true
		rootCmd.SetHelpCommand(nil)
		carapace.Gen(rootCmd).PositionalCompletion(
			carapace.ActionCallback(func(c carapace.Context) carapace.Action {
				return carapace.ActionValues(c.Getenv("LS_COLORS"))
			}),
		)
		return rootCmd
	})(func(s *Sandbox) {
		s.Run("").
			Expect(carapace.ActionValues(os.Getenv("LS_COLORS")))
	})
}
