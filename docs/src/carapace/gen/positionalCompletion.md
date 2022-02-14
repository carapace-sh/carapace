# PositionalCompletion

[`PositionalCompletion`] defines completion for positional arguments using a list of [actions](../action.md).


```go
carapace.Gen(rootCmd).PositionalCompletion(
    carapace.ActionValues("a", "b", "c"),
    // ...
)
```

[`PositionalAnyCompletion`] defines completion for any positional argument not already defined.

```go
carapace.Gen(rootCmd).PositionalAnyCompletion(
    carapace.ActionFiles(),
)
```

[`PositionalCompletion`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.PositionalCompletion
[`PositionalAnyCompletion`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.PositionalAnyCompletion
