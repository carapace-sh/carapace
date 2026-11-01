> experimental

# Multi Completer

A **multi completer** is a single binary that provides completion for multiple independent commands (e.g. `magick identify`, `magick convert`). This is useful for tools that bundle several subcommands as separate executables via symlinks or wrappers.

## Setup

Use [`WithSubcommands`](./gen.md#withsubcommands) when calling [`Gen`](./gen.md) on the root command, and call [`Execute`](./gen.md#execute) instead of `cmd.Execute()`:

```go
// root.go
package cmd

import (
    "github.com/carapace-sh/carapace"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "mytool"}

func init() {
    rootCmd.AddCommand(identifyCmd)
    rootCmd.AddCommand(convertCmd)

    carapace.Gen(rootCmd, carapace.WithSubcommands(identifyCmd, convertCmd))
}

func Execute() error {
    return carapace.Gen(rootCmd).Execute()
}
```

Each subcommand is defined in its own file:

```go
// identify.go
package cmd

import (
    "github.com/carapace-sh/carapace"
    "github.com/spf13/cobra"
)

var identifyCmd = &cobra.Command{
    Use:   "identify",
    Short: "identify image format",
    Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
    identifyCmd.Flags().StringP("format", "f", "", "image format")

    carapace.Gen(identifyCmd).FlagCompletion(carapace.ActionMap{
        "format": carapace.ActionValues("png", "jpeg", "gif"),
    })
    carapace.Gen(identifyCmd).PositionalCompletion(
        carapace.ActionFiles(".png", ".jpg"),
    )
}
```

```go
// convert.go
package cmd

import (
    "github.com/carapace-sh/carapace"
    "github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
    Use:   "convert",
    Short: "convert image format",
    Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
    convertCmd.Flags().StringP("output", "o", "", "output format")
    convertCmd.Flags().IntP("quality", "q", 0, "output quality")

    carapace.Gen(convertCmd).FlagCompletion(carapace.ActionMap{
        "output":  carapace.ActionValues("png", "jpeg"),
        "quality": carapace.ActionValues("1", "50", "75", "90", "100"),
    })
    carapace.Gen(convertCmd).PositionalCompletion(
        carapace.ActionFiles(".png", ".jpg"),
    )
}
```

## How It Works

### Snippet Generation

When `WithSubcommands` is set, [`Snippet()`](./gen/snippet.md) produces a **multi-completer snippet** that registers all subcommand completers at once. This replaces the normal single-command snippet. The binary name is automatically included as a pseudo-subcommand for self-completion.

```sh
# Generate snippet for all completers at once
mytool _carapace bash
```

Individual subcommand snippets are also available:

```sh
mytool identify _carapace bash
```

### Arg Rewriting

[`Execute()`](./gen.md#execute) intercepts `os.Args` to route completion callbacks to the correct subcommand. For example, a bridge call like:

```
mytool _carapace export "" identify -verbose image.png
```

is rewritten to:

```
mytool identify _carapace export "" -verbose image.png
```

For single completers (no `WithSubcommands`), `Execute()` simply calls `cmd.Execute()`.

### Default Subcommand

When the routing logic encounters an unknown subcommand name, it falls back to the default. By default this is the first subcommand listed in `WithSubcommands`. Override with [`WithDefault`](./gen.md#withdefault):

```go
carapace.Gen(rootCmd,
    carapace.WithSubcommands(identifyCmd, convertCmd),
    carapace.WithDefault("convert"),
)
```

## Supported Shells

Multi-completer snippets are supported for:

- bash
- bash-ble
- elvish
- fish
- nushell
- oil
- powershell
- tcsh
- xonsh
- zsh

`cmd-clink` is not supported for multi-completer snippets.

## WithSnippetFuncs

[`WithSnippetFuncs`](./gen.md#withsnippetfuncs) injects custom shell code into the generated snippet. This works for both single and multi completers. Multiple calls accumulate per shell in order.

```go
carapace.Gen(rootCmd, carapace.WithSnippetFuncs(map[string]string{
    "bash": "# custom bash initialization\n",
    "zsh":  "# custom zsh initialization\n",
}))
```
