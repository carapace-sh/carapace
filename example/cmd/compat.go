package cmd

import (
	"fmt"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var compatCmd = &cobra.Command{
	Use:   "compat",
	Short: "",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(compatCmd).Standalone()

	compatCmd.Flags()

	compatCmd.Flags().String("error", "", "ShellCompDirectiveError")
	compatCmd.Flags().String("nospace", "", "ShellCompDirectiveNoSpace")
	compatCmd.Flags().String("nofilecomp", "", "ShellCompDirectiveNoFileComp")
	compatCmd.Flags().String("filterfileext", "", "ShellCompDirectiveFilterFileExt")
	compatCmd.Flags().String("filterdirs", "", "ShellCompDirectiveFilterDirs")
	compatCmd.Flags().String("filterdirs-chdir", "", "ShellCompDirectiveFilterDirs")
	compatCmd.Flags().String("keeporder", "", "ShellCompDirectiveKeepOrder")
	compatCmd.Flags().String("default", "", "ShellCompDirectiveDefault")

	compatCmd.Flags().String("unset", "", "no completions defined")
	compatCmd.PersistentFlags().String("persistent-compat", "", "persistent flag defined with cobra")

	rootCmd.AddCommand(compatCmd)

	_ = compatCmd.RegisterFlagCompletionFunc("error", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveError
	})
	_ = compatCmd.RegisterFlagCompletionFunc("nospace", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"one", "two"}, cobra.ShellCompDirectiveNoSpace
	})
	_ = compatCmd.RegisterFlagCompletionFunc("nofilecomp", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	})

	_ = compatCmd.RegisterFlagCompletionFunc("filterfileext", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"mod", "sum"}, cobra.ShellCompDirectiveFilterFileExt
	})

	_ = compatCmd.RegisterFlagCompletionFunc("filterdirs", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveFilterDirs
	})

	_ = compatCmd.RegisterFlagCompletionFunc("filterdirs-chdir", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"subdir"}, cobra.ShellCompDirectiveFilterDirs
	})

	_ = compatCmd.RegisterFlagCompletionFunc("keeporder", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"one", "three", "two"}, cobra.ShellCompDirectiveKeepOrder
	})

	_ = compatCmd.RegisterFlagCompletionFunc("default", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return nil, cobra.ShellCompDirectiveDefault
	})

	_ = compatCmd.RegisterFlagCompletionFunc("persistent-compat", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{
			fmt.Sprintf("args: %#v toComplete: %#v", args, toComplete),
			"alternative",
		}, cobra.ShellCompDirectiveNoFileComp
	})

	compatCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		switch len(args) {
		case 0:
			return []string{"p1", "positional1"}, cobra.ShellCompDirectiveDefault
		case 1:
			return nil, cobra.ShellCompDirectiveDefault
		case 2:
			return []string{
				fmt.Sprintf("args: %#v toComplete: %#v", args, toComplete),
				"alternative",
			}, cobra.ShellCompDirectiveNoFileComp
		default:
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
	}
}
