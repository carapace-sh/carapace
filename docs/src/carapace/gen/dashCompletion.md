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

[`DashAnyCompletion`] can be combined with [`ActionExecute`] or [`ActionImport`] to complete the arguments after dash with a different completer.

[`ActionExecute`]:../action/actionExecute.md
[`ActionImport`]:../action/actionImport.md
[`DashCompletion`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.DashCompletion
[`DashAnyCompletion`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.DashAnyCompletion
