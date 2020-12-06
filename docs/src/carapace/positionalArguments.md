# Positional Arguments

```sh
command subcommand [arg1] [arg2] ...
```

> A command can have only either arguments or subcommands. A mix of these is **not** supported.

## Completion

Completion for positional arguments can be configured with [`PositionalCompletion`](https://pkg.go.dev/github.com/rsteube/carapace#Carapace.PositionalCompletion) and a list of [Actions](./action.md).

```go
carapace.Gen(callbackCmd).PositionalCompletion(
    carapace.ActionValues("a", "b", "c"),
    // ...
)
```

It is also possible to define completion for any positional argument not already explicitly defined using [`PositionalAnyCompletion`](https://pkg.go.dev/github.com/rsteube/carapace#Carapace.PositionalAnyCompletion).

```go
carapace.Gen(callbackCmd).PositionalCompletion(
    carapace.ActionValues("a", "b", "c"),
    // ...
)
```
