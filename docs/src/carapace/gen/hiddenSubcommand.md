# Hidden Subcommand

When [`carapace.Gen`](https://pkg.go.dev/github.com/rsteube/carapace#Gen) is invoked a hidden command (`_carapace`) is added to the root command unless it already exists. This handles completion script generation and [callbacks](../action/actionCallback.md).


## Completion

`SHELL` is optional and will be detected by parent process name.

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

