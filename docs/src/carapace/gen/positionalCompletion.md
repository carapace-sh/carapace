# PositionalCompletion

[`PositionalCompletion`] defines completion for positional arguments.


```go
carapace.Gen(rootCmd).PositionalCompletion(
    carapace.ActionValues("pos1", "positional1"),
    carapace.ActionValues("pos2", "positional2"),
)
```

[`PositionalCompletion`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.PositionalCompletion