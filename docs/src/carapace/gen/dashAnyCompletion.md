# DashAnyCompletion

[`DashAnyCompletion`] defines completion for any positional arguments after `--` (dash) not already defined.

```go
carapace.Gen(rootCmd).DashAnyCompletion(
    carapace.ActionValues("dAny", "dashAny"),
)
```

[`DashAnyCompletion`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.DashAnyCompletion
