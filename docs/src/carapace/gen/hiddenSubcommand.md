# Hidden Subcommand

When [`carapace.Gen`](https://pkg.go.dev/github.com/rsteube/carapace#Gen) is invoked a hidden command (`_carapace`) is added to the root command unless it already exists. This handles completion script generation and [callbacks](../action/actionCallback.md).


## Lazy Completion
```sh
command _carapace
```

## Full Completion

```sh
command _carapace [SHELL]
```

## State
```sh
command _carapace [SHELL] state command subcommand [options] ...
```

## Callback

```sh
command _carapace [SHELL] [UID] command subcommand [options] ...
```

