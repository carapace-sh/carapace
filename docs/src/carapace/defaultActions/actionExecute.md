# ActionExecute

[`ActionExecute`] invokes an internal command and parses it's output using [`ActionImport`].
It does so by updating the commands args so that [Export] is called with the current context as arguments.

```go
var executeCmd = &cobra.Command{
	Use:                "execute",
	Short:              "execute example",
	DisableFlagParsing: true,
	Run:                func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(executeCmd)

	cmd := &cobra.Command{
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {},
	}

	cmd.Flags().Bool("test", false, "")

	carapace.Gen(cmd).PositionalCompletion(
		carapace.ActionValues("one", "two"),
		carapace.ActionValues("three", "four"),
	)

	carapace.Gen(executeCmd).PositionalAnyCompletion(
		carapace.ActionExecute(cmd),
	)
}
```

E.g. in [`gh repo clone`] which allows additional git flags after dash.

```sh
gh repo clone [<repository>] [-- <gitflags>...] [flags]
```

Here, the context is updated so that only flags are completed and the `--branch` flag can list remote branches.
```go
carapace.Gen(repo_cloneCmd).DashAnyCompletion(
	carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		repo := ""
		if args := repo_cloneCmd.Flags().Args(); len(args) > 0 {
			repo = fmt.Sprintf("https://github.com/%v.git", args[0])
		}
		c.Args = append([]string{"clone", repo, ""}, c.Args...)
		return git.ActionExecute().Invoke(c).ToA()
	}),
)
```

[![asciicast](https://asciinema.org/a/468206.svg)](https://asciinema.org/a/468206)

[`ActionExecute`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionExecute
[`ActionImport`]:../action/actionImport.md
[Export]:../export.md
[`gh repo clone`]:https://github.com/rsteube/carapace-bin/blob/84717177317a9c9b1aa0d150d25d1b5c12cf9422/completers/gh_completer/cmd/repo_clone.go
