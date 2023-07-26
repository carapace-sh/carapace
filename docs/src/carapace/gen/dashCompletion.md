# DashCompletion

The first `--` (dash) argument that is not a flag argumet disables further flag parsing.
In carapace the positional arguments that follow it are completed using the following functions.
The Context is also updated to only contain the arguments after the dash.

[`DashCompletion`] defines completion for positional arguments after dash using a list of [actions](../action.md).


```go
carapace.Gen(rootCmd).DashCompletion(
    carapace.ActionValues("d1", "dash1"),
    carapace.ActionValues("d2", "dash2"),
)
```

[`DashCompletion`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.DashCompletion