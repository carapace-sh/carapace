# PositionalAnyCompletion

[`PositionalAnyCompletion`] defines completion for any positional argument not already defined.

```go
carapace.Gen(rootCmd).PositionalAnyCompletion(
    carapace.ActionValues("posAny", "positionalAny"),
)
```

[`PositionalAnyCompletion`]:https://pkg.go.dev/github.com/carapace-sh/carapace#Carapace.PositionalAnyCompletion
