# DashCompletion

The first `--` (dash) argument that is not a flag argumet disables further flag parsing.
In carapace the positional arguments that follow it are completed using the following functions.
The Context is also updated to only contain the arguments after the dash.

[`DashCompletion`] defines completion for positional arguments after dash using a list of [actions](../action.md).


```go
carapace.Gen(rootCmd).DashCompletion(
    carapace.ActionValues("a", "b", "c"),
    // ...
)
```

[`DashAnyCompletion`] defines completion for any positional argument after dash not already defined.

```go
carapace.Gen(rootCmd).DashAnyCompletion(
    carapace.ActionFiles(""),
)
```


## Complete using different command

[`DashAnyCompletion`] can be combined with [`ActionInvoke`] or [`ActionImport`] to complete the arguments after dash with a different completer.
E.g. for `gh repo fork` which allows additional git flags after dash.

```sh
gh repo fork [<repository>] [-- <gitflags>...] [flags]
```

Here the context is updated so that only flags are completed and the `--branch` flag can complete remote branches.
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

[`ActionInvoke`]:../action/actionInvoke.md
[`ActionImport`]:../action/actionImport.md
[`DashCompletion`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.DashCompletion
[`DashAnyCompletion`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.DashAnyCompletion
