# ActionInvoke

[`ActionInvoke`] invokes a different command and parses it's output using [`ActionImport`].
It does so by updating `os.Args` so that [Export] is called with the current context as arguments.

E.g. in [`gh repo fork`] which allows additional git flags after dash.

```sh
gh repo fork [<repository>] [-- <gitflags>...] [flags]
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
		return carapace.ActionInvoke(git.Execute).Invoke(c).ToA()
	}),
)
```

[![asciicast](https://asciinema.org/a/468206.svg)](https://asciinema.org/a/468206)

[`ActionInvoke`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionInvoke
[`ActionImport`]:../action/actionImport.md
[Export]:../export.md
[`gh repo fork`]:https://github.com/rsteube/carapace-bin/blob/84717177317a9c9b1aa0d150d25d1b5c12cf9422/completers/gh_completer/cmd/repo_fork.go
