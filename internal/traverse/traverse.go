package traverse

import "github.com/spf13/cobra"

func Traverse(cmd *cobra.Command, args []string) (*cobra.Command, error) {
	return traverse(cmd, args[:len(args)-1])
}

func traverse(cmd *cobra.Command, args []string) (*cobra.Command, error) {
	preRun(cmd, args)

	if err := cmd.ParseFlags(args); err != nil { // TODO filter errors
		return nil, err
	} else {
		args = cmd.Flags().Args()
	}

	return cmd, nil
}

func preRun(cmd *cobra.Command, args []string) {
	if subcommand, _, err := cmd.Find([]string{"_carapace"}); err != nil {
		if subcommand.PreRun != nil {
			subcommand.PreRun(cmd, args)
		}
	}
}