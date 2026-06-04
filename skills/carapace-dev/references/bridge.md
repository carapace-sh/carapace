# Carapace-bridge: Implementing & Debugging Bridges

Reference for implementing bridges in [carapace-bridge](https://github.com/carapace-sh/carapace-bridge) — the library that delegates shell completion to external completion frameworks.

## Architecture

Three layers work together:

1. **Bridge Actions** (`pkg/actions/bridge/`) — Individual bridge implementations returning `carapace.Action`
2. **Bridge Discovery** (`pkg/bridges/`) — Runtime detection of which commands have completions in each shell
3. **Choice System** (`pkg/choices/`) — Persistent user overrides pinning a specific bridge

## Bridge Function Signature

All bridge functions follow the same pattern:

```go
func ActionXxx(command ...string) carapace.Action
```

The variadic `command` arg: first element is the command name, rest are subcommands. Register in `bridgeActions` map in `bridge.go`.

## The `actionCommand` Wrapper

Every bridge function is wrapped with `actionCommand(command...)` (`carapace.go:54`) which handles the "no command provided" case:

- `len(command) > 0`: calls inner function directly
- `command` is empty: creates a temporary `cobra.Command` that offers `ActionExecutables()` + `ActionFiles()` for the first positional arg, then delegates with `.Shift(1)` for subsequent args

This enables `bridge.ActionCobra()()` style usage where the command name is the first argument.

## Resolution Flow: `ActionBridge(command...)`

1. **Choice lookup**: Check `~/.config/carapace/choices/<name>` file. If `choice.Group == "" || choice.Group == "bridge"`, use `choice.Variant` directly.
2. **Environment fallback**: Iterate `CARAPACE_BRIDGES` env var. For each bridge, check if command is in that bridge's discovered completers. First match wins.
3. **No match**: Return `ActionValues()` (empty).

## Adding a New Bridge

1. Create `<name>.go` in `pkg/actions/bridge/`
2. Implement `func Action<Name>(command ...string) carapace.Action` using `actionCommand` wrapper
3. Add entry to `bridgeActions` map in `bridge.go`
4. Optionally add to `candidates` list in `detect.go` for auto-detection
5. Add builtins to `pkg/bridges/<name>_builtins.go` if it's a shell bridge

## Bridge Implementation Patterns

### Invoking the External Completion

| Bridge type | Invocation pattern |
|-------------|-------------------|
| Cobra | `command[0] __complete <subcommands> <args> <value>` |
| Carapace | `command[0] _carapace export "" <subcommands> <args> <value>` |
| CarapaceBin | `carapace <command> export "" <subcommands> <args> <value>` |
| Argcomplete v3 | `_ARGCOMPLETE_STDOUT_FILENAME=<tmpfile> _ARGCOMPLETE=1 <command> <args>` |
| Click | `_<CMD>_COMPLETE=zsh_complete <command>` |
| Yargs | `<command> --get-yargs-completions <args>` |
| Clap | `<command> complete --index <N> --type 9 -- <args>` |
| Kingpin | `<command> --completion-bash <args>` |
| Urfavecli v3 | `<command> <subcommands> --generate-shell-completion` |
| Bash | `bash --rcfile <config> -i <script>` with `COMP_LINE` set |
| Zsh | `zsh --no-rcs -c <script> -- <command> <args>` |
| Fish | `complete --do-complete=<quoted args>` |
| PowerShell | `[CommandCompletion]::CompleteInput()` via profile script |

### Output Parsing Conventions

| Format | Bridges using it | Parsing |
|--------|-----------------|---------|
| Tab-separated `value\tdescription` | Cobra, Fish | Split on `\t` |
| Newline-separated values | Kingpin, Complete, Clap | Split on `\n` |
| `name:description` per line | Urfavecli v3, Yargs | Split on first `:` |
| `value (description)` per line | Bash | Split on `(` |
| `value -- description` per line | Zsh | Split on ` -- ` |
| JSON | PowerShell, Inshellisense, Kitten | `json.Unmarshal` |
| Groups of 3 lines (type/value/desc) | Click | State machine |
| Native JSON export | Carapace, CarapaceBin | `ActionImport()` |

### NoSpace Detection

Most bridges apply `NoSpace('/=@:.,')` by default. Some use heuristic detection:

| Pattern | Used by |
|---------|---------|
| Always `NoSpace('/=@:.,')` | Bash, Zsh, Fish, Urfavecli |
| Check if any value ends with `/=@:.,` (excluding last char) | Kingpin, Complete, Clap |
| Dynamically from `CompletionText` trailing space | PowerShell |
| Check for `=`, `/`, `,` suffix in results | Argcomplete |
| Omit trailing space per value | Nushell, Elvish, Xonsh |

### Fallback Behavior

Most bridges fall back to `ActionFiles()` when the external completer returns empty results. This ensures the user still gets useful completions even when the bridge source has nothing.

## Cobra Bridge: Directive Handling

The cobra bridge parses `ShellCompDirective` from the last output line (`:N` format):

| Directive | Action |
|-----------|--------|
| `ShellCompDirectiveFilterDirs` | `ActionDirectories()` |
| `ShellCompDirectiveFilterFileExt` | `ActionFiles(extensions...)` |
| `ShellCompDirectiveNoFileComp` + <3 lines | Fallback to `ActionFiles()` |
| `ShellCompDirectiveError` | Empty `ActionValues()` |
| `ShellCompDirectiveNoSpace` | `.NoSpace()` |

## Argcomplete: Optarg Handling

When completing `--flag=value`:
1. Split on `=`
2. Add flag as a separate arg
3. Set `current=""` to get all flag values
4. Re-add `=` prefix via `.Prefix()`

For partial flags (`--par`), set `current="--"` to get all flags. For partial positionals, set `current=""`.

## Shell Bridge Config Directories

Each shell bridge creates a custom config directory on first use via `ensureExists()`:

| Shell | Config path | Purpose |
|-------|------------|---------|
| bash | `~/.config/carapace/bridge/bash/.bashrc` | Custom bash rc for bridge invocation |
| zsh | `~/.config/carapace/bridge/zsh/.zshrc` | Custom zsh rc |
| fish | `~/.config/carapace/bridge/fish/config.fish` | Custom fish config |
| powershell | `~/.config/carapace/bridge/powershell/Microsoft.PowerShell_profile.ps1` | Custom PS profile |

## Auto-Detection: `Detect(cmd)`

`Detect()` probes a command to determine which completion framework it uses (experimental):

1. Verify command exists via `exec.LookPath`
2. Iterate `candidates` in fixed order: `carapace`, `cobra`, `complete`, `kingpin`, `urfavecli`, `urfavecli_v1`, `argcomplete`, `argcomplete_v1`, `click`, `yargs`
3. For each candidate, create a temp dir, set up a throwaway cobra command with the bridge action as `PositionalAnyCompletion`, call `x.Complete()` to trigger completion with cursor on `-`
4. Check if any returned value is `-h`, `-help`, or `--help`
5. First match wins → return `{Name, Action}` tuple

**Limitation**: Detection relies on the command outputting help flags when completing `-`. If a framework doesn't suggest `-h`/`--help`, detection fails.

## Bridge Discovery: `pkg/bridges/`

Shell bridges discover which commands have completions at runtime:

| Shell | Discovery method | Cache |
|-------|-----------------|-------|
| bash | Scan bash completion directories | 24h JSON in XDG cache |
| zsh | Invoke zsh to list completion functions (embedded `zsh.sh`) | 24h JSON |
| fish | Scan fish `fish_complete_path` directories | 24h JSON |
| inshellisense | `inshellisense specs list` | 24h JSON |

Each shell has a builtin blacklist (e.g., `bash_builtins.go`, `zsh_builtins.go`) to exclude shell internals.

`Bridges()` aggregates: only includes bridges listed in `CARAPACE_BRIDGES` env var. First bridge to claim a command wins (no overwrites).

## Testing

No unit tests in carapace-bridge. Integration testing via Docker Compose in `.docker/`:

| Service | Tests bridge for |
|---------|-----------------|
| click | Python Click framework |
| cobra | Go cobra framework |
| kingpin | Go kingpin framework |
| argcomplete | Python argcomplete |
| urfavecli | Go urfave/cli |
| yargs | Node.js yargs |
| inshellisense | Microsoft inshellisense |
| complete | Go posener/complete |

## Related Skills

- **carapace skill** (references/choice.md) — resolution priority, choices, CARAPACE_BRIDGES (user-facing)
- **carapace skill** (references/integrate.md) — using bridges in cobra CLIs (Split, SplitP, ActionCarapaceBin)
