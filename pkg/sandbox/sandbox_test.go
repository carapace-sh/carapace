package sandbox

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
)

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
