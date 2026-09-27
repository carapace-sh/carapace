# Carapace Library: _carapace CLI

Reference for [carapace](https://github.com/carapace-sh/carapace)'s hidden `_carapace` subcommand and its subcommands in `command.go`.

## Registration

`carapace.Gen(cmd)` calls `addCompletionCommand(targetCmd)` which adds the `_carapace` subcommand to any command that doesn't already have one:

```go
func addCompletionCommand(targetCmd *cobra.Command) {
    for _, c := range targetCmd.Commands() {
        if c.Name() == "_carapace" {
            return
        }
    }
    // add _carapace subcommand
}
```

The subcommand is `Hidden: true` and uses `DisableFlagParsing: true` with `UnknownFlags: true` to handle unknown flags gracefully during completion.

## _carapace Entry Point

```go
carapaceCmd := &cobra.Command{
    Use:    "_carapace",
    Hidden: true,
    Run: func(cmd *cobra.Command, args []string) {
        LOG.Print(strings.Repeat("-", 80))
        LOG.Printf("%#v", os.Args)

        if len(args) > 2 && strings.HasPrefix(args[2], "_") {
            cmd.Hidden = false
        }

        if !cmd.HasParent() {
            panic("missing parent command")
        }

        parentCmd := cmd.Parent()
        if parentCmd.Annotations[annotation_standalone] == "true" {
            parentCmd.RemoveCommand(cmd) // don't complete local _carapace in standalone mode
        }

        if s, err := complete(parentCmd, args); err != nil {
            fmt.Fprintln(io.MultiWriter(parentCmd.OutOrStderr(), LOG.Writer()), err.Error())
        } else {
            fmt.Fprintln(io.MultiWriter(parentCmd.OutOrStdout(), LOG.Writer()), s)
        }
    },
}
```

Arguments to `_carapace`:
- `args[0]` = shell name (e.g., `"bash"`, `"zsh"`)
- `args[1]` = root command name (e.g., `"git"`)
- `args[2+]` = completion arguments (e.g., `"branch"`, `"--force"`)

## Dispatch to complete()

`complete(parentCmd, args)` is the main dispatch function (in `complete.go`). It handles:
- 0 args: output snippet for auto-detected shell
- 1 arg: output snippet for named shell
- 2+ args: perform actual completion

## _carapace spec

```go
specCmd := &cobra.Command{
    Use: "spec",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Fprint(cmd.OutOrStdout(), spec.Spec(targetCmd))
    },
}
carapaceCmd.AddCommand(specCmd)
```

Running `myapp _carapace spec` outputs YAML spec for the entire command tree.

## _carapace style

```go
styleCmd := &cobra.Command{
    Use:  "style",
    Args: cobra.ExactArgs(1),
    Run:  func(cmd *cobra.Command, args []string) {},
}
carapaceCmd.AddCommand(styleCmd)

styleSetCmd := &cobra.Command{
    Use:  "set",
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        for _, arg := range args {
            if splitted := strings.SplitN(arg, "=", 2); len(splitted) == 2 {
                if err := style.Set(splitted[0], splitted[1]); err != nil {
                    fmt.Fprint(cmd.ErrOrStderr(), err.Error())
                }
            }
        }
    },
}
styleCmd.AddCommand(styleSetCmd)
Carapace{styleSetCmd}.PositionalAnyCompletion(
    ActionStyleConfig(),
)
```

`_carapace style set` updates `styles.json` on disk. `_carapace style set key=value [key=value...]`.

## Positional Completion on _carapace

The `_carapace` command itself has completion for its arguments:

```go
Carapace{carapaceCmd}.PositionalCompletion(
    ActionStyledValues("bash", "zsh", "fish", ...),  // shell names
    ActionValues(targetCmd.Root().Name()),           // root command name
)
Carapace{carapaceCmd}.PositionalAnyCompletion(
    ActionCallback(func(c Context) Action {
        args := []string{"_carapace", "export", ""}
        args = append(args, c.Args[2:]...)
        args = append(args, c.Value)

        executable, err := os.Executable()
        if err != nil {
            return ActionMessage(err.Error())
        }
        return ActionExecCommand(executable, args...)(func(output []byte) Action {
            if string(output) == "" {
                return ActionValues()
            }
            return ActionImport(output)
        })
    }),
)
```

The `PositionalAnyCompletion` re-invokes the binary with the `export` shell to get completions for arbitrary arguments. This is the mechanism that enables `_carapace export` to work.

## Standalone Mode

```go
if parentCmd.Annotations[annotation_standalone] == "true" {
    parentCmd.RemoveCommand(cmd)
}
```

When a command is in standalone mode (`annotation_standalone = "true"`), the `_carapace` subcommand is removed from the command tree. This prevents completion of the hidden carapace command itself.

## Multi-Completer Support

A single binary can provide completions for multiple commands (e.g. `magick`, `identify`, `convert` in ImageMagick). This is configured via `Gen()` options defined in `carapace.go`:

```go
carapace.Gen(rootCmd,
    carapace.WithSubcommands(identifyCmd, convertCmd),
    carapace.WithDefault("magick"),
    carapace.WithSnippetFuncs(map[string]string{
        "bash": "_magick_custom_func() { ... }",
    }),
)
```

### Options

| Option | Purpose |
|--------|---------|
| `WithSubcommands(cmds ...*cobra.Command)` | Register additional commands as independent completers. The binary name is automatically included as a pseudo-subcommand for self-completion. |
| `WithDefault(name string)` | Set the shell completer function name (defaults to executable name). No-op without `WithSubcommands`. |
| `WithSnippetFuncs(funcs map[string]string)` | Inject custom shell code into snippets (keyed by shell name). Multiple calls accumulate per shell in order. |

### Snippet Generation

When `WithSubcommands` is set, `Snippet()` delegates to `multi.Snippet()` (`internal/shell/multi/snippet.go`) instead of `shell.Snippet()`. The multi-completer snippet registers all subcommand completers at once in a single shell script.

Multi-completer snippets are supported by all shells except `cmd-clink` and `ion`.

### Execute() Arg Rewriting

`Execute()` intercepts `os.Args` for multi-completer routing via `rewriteArgs()`. The rewriting handles three patterns:

1. **Root `_carapace` invocation** — `binary _carapace <shell> <args...>` is rewritten to route to the correct subcommand based on the completion target. If the target is the binary itself (pseudo-subcommand), it's handled at root level.

2. **Bridge re-invocation** — `binary _carapace export "" <subcommand> <args...>` (from `bridge.ActionCarapace`) is rewritten to `binary <subcommand> _carapace export "" <args...>` so the subcommand's own completion handles it.

3. **Subcommand snippet request** — `binary <subcommand> _carapace <shell>` produces a single-command snippet via `multi.SingleSnippetOrEmpty()`.

### Storage Entry Fields

The `entry` struct in `storage.go` gains these fields when multi-completer options are applied:

| Field | Type | Purpose |
|-------|------|---------|
| `subcommands` | `[]*cobra.Command` | The registered subcommands |
| `subcommandNames` | `[]string` | Binary name + subcommand names (for snippet generation) |
| `defaultName` | `string` | Shell completer function name |
| `snippetFuncs` | `map[string][]string` | Custom shell code per shell name |

## Annotation Constants

Defined in `carapace.go`:
- `annotation_standalone` — marks a command as standalone (no `_carapace` subcommand)
- `annotation_skip` — skips the command in certain operations

## Re-invocation Flow

```
Shell invokes: myapp _carapace bash git branch --force
    │
    ├─ complete(git, ["bash", "git", "branch", "--force"])
    │    └─ traverse(git, ["branch", "--force"])
    │         └─ storage.getFlag(git, "force")
    │              └─ action.Invoke(context).value("bash")
    │
    └─ myapp _carapace export git branch --force --force
         └─ ActionImport(json) → re-invokes the export shell
```

The export shell re-invocation is used for subcommand completion to resolve nested command trees.

## Gotchas

- **Hidden subcommand visible with underscore prefix**: If `args[2]` starts with `_`, `cmd.Hidden = false` temporarily makes `_carapace` visible. This allows completion for `_carapace` itself.
- **JSON export via ActionExecCommand**: The `PositionalAnyCompletion` on `_carapace` uses `ActionExecCommand` to re-invoke the binary with `export` shell and parse the JSON result back via `ActionImport`. This is the only path for nested subcommand completion.
- **complete() log output**: `LOG.Printf("%#v", os.Args)` prints all invocations when `CARAPACE_LOG` is set, including the full argument list.
- **Standalone removes _carapace**: In standalone mode, the `_carapace` subcommand is removed after the parent command is set up, so it doesn't appear in the completion candidates for the root command.

## Related Skills

- **references/traverse.md** — complete() dispatch and traverse()
- **references/action.md** — ActionImport for parsing export JSON
- **references/export.md** — the JSON wire format
- **references/spec.md** — spec.Spec used by `_carapace spec`