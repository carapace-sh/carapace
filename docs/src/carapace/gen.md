# Gen

Calling [`Gen`](https://pkg.go.dev/github.com/rsteube/carapace#Gen) on the root command is sufficient to enable completion script generation using the [Hidden Subcommand](#hidden-subcommand).

```go
import (
    "github.com/rsteube/carapace"
)

carapace.Gen(rootCmd)
```

Additionally invoke [`carapace.Test`](https://pkg.go.dev/github.com/rsteube/carapace#Test) in a [test](https://golang.org/doc/tutorial/add-a-test) to verify configuration during build time.
```go
func TestCarapace(t *testing.T) {
    carapace.Test(t)
}
```

## Hidden Subcommand

When [`Gen`](https://pkg.go.dev/github.com/rsteube/carapace#Gen) is invoked a hidden subcommand (`_carapace`) is added. This handles completion script generation and [callbacks](./defaultActions/actionCallback.md).


### Completion

`SHELL` is optional and will be detected by parent process name.

```sh
command _carapace [SHELL]
```

```sh
# bash
source <(command _carapace)

# elvish
eval (command _carapace | slurp)

# fish
command _carapace | source

# nushell (update config.nu according to output)
command _carapace nushell

# oil
source <(command _carapace)

# powershell
Set-PSReadLineOption -Colors @{ "Selection" = "`e[7m" }
Set-PSReadlineKeyHandler -Key Tab -Function MenuComplete
command _carapace | Out-String | Invoke-Expression

# tcsh
set autolist
eval `command _carapace tcsh`

# xonsh
COMPLETIONS_CONFIRM=True
exec($(command _carapace))

# zsh
source <(command _carapace)
```

> Directly sourcing multiple completions in your shell init script increases startup time [considerably](https://medium.com/@jzelinskie/please-dont-ship-binaries-with-shell-completion-as-commands-a8b1bcb8a0d0). See [lazycomplete](https://github.com/rsteube/lazycomplete) for a solution to this problem.
