---
name: carapace-dev-shell
description: >
  Use when implementing or debugging carapace's per-shell output formatting and snippet
  generation. Covers the 12 supported shells, per-shell quoting/escaping rules, snippet
  format, value formatting, nospace handling, style rendering, message display, and shell-specific
  features (zsh zstyle, bash COMP_TYPE, elvish complex-candidate, nushell style format).
  Triggers on: "carapace shell output", "shell formatting", "shell snippet", "shell quoting",
  "bash completion script", "zsh completion script", "fish completion script", "nushell completion",
  "per-shell formatting", "shell internals".
user-invocable: true
---

# Carapace Library: Per-Shell Output Formatting

Reference for [carapace](https://github.com/carapace-sh/carapace)'s shell-specific completion output — how `RawValues` and `Meta` are formatted for each of the 12 supported shells.

## Supported Shells

| Shell | Snippet name | Package |
|-------|-------------|---------|
| Bash | `bash` | `internal/shell/bash/` |
| Bash BLE | `bash-ble` | `internal/shell/bash_ble/` |
| Cmd (Clink) | `cmd-clink` | `internal/shell/cmd_clink/` |
| Elvish | `elvish` | `internal/shell/elvish/` |
| Fish | `fish` | `internal/shell/fish/` |
| Ion | `ion` | `internal/shell/ion/` |
| Nushell | `nushell` | `internal/shell/nushell/` |
| Oil | `oil` | `internal/shell/oil/` |
| PowerShell | `powershell` | `internal/shell/powershell/` |
| Tcsh | `tcsh` | `internal/shell/tcsh/` |
| Xonsh | `xonsh` | `internal/shell/xonsh/` |
| Zsh | `zsh` | `internal/shell/zsh/` |
| Export (JSON) | `export` | `internal/shell/export/` |

## Dispatch: `shell.Value()`

`shell.Value(shell, value, meta, values)` is called after traversal and invocation. Before dispatching to the per-shell formatter, it applies:

1. **Color disable**: if `env.ColorDisabled()`, strips styles and sets fallback styles
2. **Prefix filtering**: `values.FilterPrefix(value)` unless `CARAPACE_UNFILTERED` is set
3. **Flag merging**: merges "shorthand flags"/"longhand flags" tags into "flags" — implicit for zsh, explicit via `CARAPACE_MERGEFLAGS`
4. **Message integration**: for shells without native message support (not elvish/export/zsh), messages are injected as synthetic `RawValue` entries (styled with `style.Carapace.Error`)
5. **Nospace propagation**: if messages exist or `CARAPACE_NOSPACE` is set, add `*` to nospace set
6. **Sort + dedup**: `sort.Sort(ByDisplay(values))` → clear UIDs → `values.Unique()`

## Snippet Generation

Each shell package has a `Snippet(cmd *cobra.Command) string` function that generates the completion integration script. The snippet:

- Registers the command with the shell's completion system
- Sets up the callback to invoke `carapace` with the correct arguments
- Handles retry logic (appending `''`, `'`, `"` for open-quote scenarios)

## Per-Shell Output Formats

### Bash

- **Format**: `\001` separates nospace flag from completion data. Values are shell-escaped.
- **Quoting**: Two modes — `escapingReplacer` (unquoted) and `escapingQuotedReplacer` (double-quoted) with auto-detection via `requiresQuoting()`.
- **COMP_TYPE handling**: On successive tabs (`COMP_TYPE_LIST_SUCCESSIVE_TABS`), shows display+description instead of values.
- **Partial completion workaround**: When all displays share a prefix, bash inserts it losing formatting. Workaround collapses to common prefix.
- **Snippet**: `complete -F` function, exports `COMP_LINE`/`COMP_POINT`/`COMP_TYPE`/`COMP_WORDBREAKS`.

### Zsh

- **Format**: `zstyle\001message\001data` where data is tag-grouped with `\002`/`\003` delimiters for `_describe`.
- **Quoting**: 5-state machine (DEFAULT, QUOTING_ESCAPING, QUOTING, FULL_QUOTING_ESCAPING, FULL_QUOTING) based on `shlex.Split` analysis. Each state applies different escaping and matching closing quote.
- **zstyle**: Generates `zstyle` color rules using `(#b)` pattern matching. Two patterns per value. Disabled for >500 values (performance).
- **Named directories**: Parses `CARAPACE_ZSH_HASH_DIRS` (output of `hash -d`) to support `~nameddir/` path completion.
- **Snippet**: `#compdef`, passes `CARAPACE_COMPLINE` and `CARAPACE_ZSH_HASH_DIRS`.

### Fish

- **Format**: Simplest — `value\tdescription` tab-separated, newline-delimited.
- **Snippet**: `complete -c 'cmd' -f -a '(_completion)' -r` with `(commandline -cp)` and open-quote retry.
- **Optarg detection**: If `--flag=value`, styles only the `value` part via `style.ForPath`.

### Elvish

- **Format**: JSON with `complexCandidate` struct including `CodeSuffix` (space or empty for nospace). Uses elvish's `styled` builtin.
- **Message display**: Uses `edit:notify` for error messages.
- **Snippet**: Sets `edit:completion:arg-completer[cmd]`, uses `from-json` to parse output. Adds `.exe` completer on Windows.
- **Suppresses usage** when values exist (avoids notification spam).

### Nushell

- **Format**: JSON array of `{value, display, description, style}`. Values with special chars are double-quoted.
- **Nospace**: Handled by omitting trailing space.
- **Style mapping**: Full 256-color + attributes to nushell format (`fg`, `bg`, `attr` with `l/b/d/i/r/u` flags).
- **Snippet**: Closure `let cmd_completer = {|spans| ...}` piped through `from json`.

### PowerShell

- **Format**: JSON `CompletionResult` with `CompletionText`, `ListItemText`, `ToolTip`.
- **Styles**: SGR escape codes in `ListItemText` colors. `ToolTip` used for tooltip mode via `CARAPACE_TOOLTIP`.
- **Empty strings**: Replaced with space (PowerShell `CompletionResult` rejects empties).
- **Snippet**: `Register-ArgumentCompleter -Native` with cursor-aware `$commandAst` parsing. Strips single quotes from args. Registers both `cmd` and `cmd.exe`.

### Oil

- **Format**: Uses `\001` as inline nospace indicator appended to values. Description shown in parentheses for multi-value lists.
- **Snippet**: Similar to bash but simpler — uses `COMP_LINE`/`COMP_POINT`, `mapfile`, `\001` as nospace indicator.

### Xonsh

- **Format**: JSON array of `{Value, Display, Description, Style}`. Special chars trigger single-quoted values; backslash triggers raw strings (`r'...'`).
- **Style mapping**: Converts to xonsh format (`bg:ansiblack fg:ansired bold italic underline blink reverse`). Full 256-color hex/ansi mapping.
- **Snippet**: Python `@contextual_command_completer` decorator with `RichCompletion` objects.

### Tcsh

- **Format**: Values on separate lines. `nospace` indicator using tcsh's built-in mechanism.
- **Snippet**: `complete` command setup with cursor position handling.

### Ion

- **Format**: Simple value listing. Minimal shell support.

### Cmd (Clink)

- **Format**: Lua-compatible completion output via `CARAPACE_COMPLINE` environment variable.
- **Patching**: Re-lexes `CARAPACE_COMPLINE` via `carapace-shlex` to handle redirects like bash.

### Export (JSON)

- **Format**: `{"Version": "...", "Meta": {...}, "Values": [...]}` — used by `ActionImport()` for internal bridge communication.

## Nospace Handling

Each shell handles "no trailing space" differently:

| Shell | Mechanism |
|-------|-----------|
| bash | `compopt -o nospace` when nospace matches |
| zsh | `CodeSuffix` in completion entry (empty = no space) |
| fish | Natively handles via tab-separated format |
| elvish | `CodeSuffix` in `complexCandidate` (empty = no space) |
| nushell | Omits trailing space in JSON output |
| powershell | Strips trailing space from `CompletionText` |
| oil | `\001` suffix indicates nospace |
| xonsh | `append_closing_quote=False` in `RichCompletion` |
| tcsh | Built-in nospace support |

## Related Skills

- **carapace-setup** — installation and shell integration (user-facing)
- **carapace-dev-traverse** — the completion engine that produces Actions before formatting
- **carapace-dev-style** — how styles are resolved before shell rendering