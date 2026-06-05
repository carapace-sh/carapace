# Carapace Library: Per-Shell Output Formatting

Reference for [carapace](https://github.com/carapace-sh/carapace)'s shell-specific completion output — how `RawValues` and `Meta` are formatted for each of the 12 supported shells, and how the shells differ.

## Supported Shells

| Shell | Snippet name | Package | Deep-dive skill |
|-------|-------------|---------|-----------------|
| Bash | `bash` | `internal/shell/bash/` | `references/shell-bash.md` |
| Bash BLE | `bash-ble` | `internal/shell/bash_ble/` | `references/shell-bash-ble.md` |
| Oil | `oil` | `internal/shell/oil/` | `references/shell-oil.md` |
| Cmd (Clink) | `cmd-clink` | `internal/shell/cmd_clink/` | **cmd-clink** skill |
| Elvish | `elvish` | `internal/shell/elvish/` | `references/shell-elvish.md` |
| Fish | `fish` | `internal/shell/fish/` | `references/shell-fish.md` |
| Ion | `ion` | `internal/shell/ion/` | — |
| Nushell | `nushell` | `internal/shell/nushell/` | `references/shell-nushell.md` |
| PowerShell | `powershell` | `internal/shell/powershell/` | `references/shell-powershell.md` |
| Tcsh | `tcsh` | `internal/shell/tcsh/` | **tcsh** skill |
| Xonsh | `xonsh` | `internal/shell/xonsh/` | `references/shell-xonsh.md` |
| Zsh | `zsh` | `internal/shell/zsh/` | `references/shell-zsh.md` |
| Export (JSON) | `export` | `internal/shell/export/` | — |

## Shared Dispatch: `shell.Value()`

`shell.Value(shell, value, meta, values)` is called after traversal and invocation. Before dispatching to the per-shell formatter, it applies these steps (in order):

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

## Cross-Shell Feature Comparison

### Primary Shells

| Feature | Bash | Bash BLE | Oil | Zsh | Fish | Elvish | Nushell | Xonsh | PowerShell |
|---------|------|----------|-----|-----|------|--------|---------|-------|------------|
| **Argument acquisition** | `COMP_LINE`/`COMP_POINT` env vars | `COMP_LINE`/`COMP_POINT` env vars | `COMP_LINE`/`COMP_POINT` env vars | `CARAPACE_COMPLINE` env var + `words` array | `commandline -cp` + `xargs` | `(all $arg)` rest param | `...$spans` spread | `CommandContext.args` + `sub_proc_get_output` | `$commandAst.CommandElements` |
| **Registration** | `complete -o noquote -F func cmd` | `complete -F func_ble cmd` | `complete -F func cmd` | `compdef func cmd` + `compquote` guard | `complete -c 'cmd' -f -a '(func)' -r` | `edit:completion:arg-completer[cmd]` | Closure `let cmd_completer = {\|spans\| ...}` | `@contextual_command_completer` + `add_one_completer` | `Register-ArgumentCompleter -Native` |
| **Output format** | `\001`-delimited nospace+values | Tab-delimited value/display with `\x1c` separators | Values with `\001` nospace indicator | `\001`-delimited zstyle+message+tag groups | Tab-separated `value\tdesc\n` | JSON `completion` with `complexCandidate` array | JSON `{value,display,description,style}` | JSON `{Value,Display,Description,Style}` | JSON `CompletionResult` array |
| **Nospace** | Global `compopt -o nospace` | Per-candidate suffix field | Inline `\001` indicator | Per-candidate space suffix in `_describe` values | **Not supported** | Per-candidate `CodeSuffix` | Trailing space in `value` field | Space baked into `Value` field | Trailing space in `CompletionText` |
| **Style/color** | Not supported | Not supported (carapace doesn't emit styles) | Not supported | `zstyle list-colors` with `(#b)` patterns | Not supported | `styled` builtin via `ParseStyling` | Full 256-color + attr `{fg,bg,attr}` record | `bg: fg:` format with `ansi` prefix | SGR escape codes in `ListItemText` |
| **Messages** | Integrated as styled values | Integrated as values | Integrated as values | Native `_message -r` | Integrated as `ERR` values | Native `edit:notify` | Integrated as `ERR` values | Integrated as styled values | Integrated as styled `ListItemText` |
| **Go-side patching** | `bash.Patch()` (redirects + wordbreaks) | None | None | None | None | None | `nushell.Patch()` (open quotes) | None | None |
| **Open-quote retry** | 3-stage snippet retry (`''`, `'"`, `"`) | `sed` + xargs | `sed` + xargs | 3-stage snippet retry (`''`, `'"`, `"`) | 3-stage snippet retry (`''`, `'"`, `"`) | Not needed (parser handles) | Not needed (parser handles) | `fix_prefix` strips quotes | Not needed (AST handles) |
| **Description display** | COMP_TYPE=63 only | Always (tab-delimited) | Inline `(description)` | Right column via `_describe` | Pager column | Right column with `(description)` | Right column | Right column | `ListItemText` or `ToolTip` |
| **Tag grouping** | Not supported | Not supported | Not supported | Native via `_describe -t` | Not supported | Not supported | Not supported | Not supported | Not supported |
| **Flag merging** | Explicit only (`CARAPACE_MERGEFLAGS`) | Explicit only | Explicit only | **Implicit** (always merges) | Explicit only | Explicit only | Explicit only | Explicit only | Explicit only |
| **Named directories** | Not supported | Not supported | Not supported | Supported via `hash -d` | Not supported | Not supported | Not supported | Not supported | Not supported |
| **Extra env vars** | `COMP_LINE`, `COMP_POINT`, `COMP_TYPE`, `COMP_WORDBREAKS` | Same as bash | `COMP_LINE`, `COMP_POINT` | `CARAPACE_COMPLINE`, `CARAPACE_ZSH_HASH_DIRS` | None | None | None | None | None |

### Secondary Shells

| Feature | Tcsh | Ion | Clink | Export |
|---------|------|-----|-------|--------|
| **Output format** | Values on lines, backslash-escaped | JSON `{Value, Display}` | Tab-delimited `value\tdisplay\tdesc\tsuffix` | JSON `{"Version","Meta","Values"}` |
| **Nospace** | Built-in support | Trailing space in `Value` | Per-candidate `appendchar` field | N/A (JSON styles) |
| **Style/color** | Not supported | Not supported | Not supported | JSON styles |
| **Go-side patching** | None | None | `cmd_clink.Patch()` (redirects via CARAPACE_COMPLINE) | None |
| **Snippet** | `complete` with `$COMMAND_LINE` | Empty (no snippet) | Lua function with `CARAPACE_COMPLINE` | N/A (internal use) |

## Nospace Handling Comparison

Each shell handles "no trailing space" differently:

| Shell | Mechanism | Per-candidate? |
|-------|-----------|---------------|
| Bash | `compopt -o nospace` when nospace matches | No (global) |
| Bash BLE | Suffix field in `\x1c`-delimited candidate | Yes |
| Oil | `\001` suffix indicates nospace | Yes (inline) |
| Zsh | Space suffix in `_describe` values (empty = no space) | Yes |
| Fish | **Not supported** — tab-separated format has no nospace mechanism | N/A |
| Elvish | `CodeSuffix` in `complexCandidate` (empty = no space) | Yes |
| Nushell | Trailing space in `value` field (absent = no space) | Yes |
| PowerShell | Trailing space in `CompletionText` (absent = no space) | Yes |
| Xonsh | Space baked into `Value` field (absent = no space) | Yes |
| Tcsh | Built-in nospace support | Yes |
| Clink | Per-candidate `appendchar` field | Yes |

## Message Handling Comparison

| Shell | Approach | Native support? |
|-------|----------|----------------|
| Zsh | `_message -r` builtin | Yes |
| Elvish | `edit:notify` command | Yes |
| Export | JSON `Messages` field | Yes |
| All others | Integrated as synthetic `ERR` values via `meta.Messages.Integrate()` | No |

## Go-Side Patching Comparison

Only three shells require argument patching before traversal:

| Shell | Patch function | What it does |
|-------|---------------|-------------|
| Bash | `bash.Patch()` | Re-lexes `COMP_LINE` with carapace-shlex; strips redirects; extracts wordbreak prefix and COMP_TYPE as package globals; unsets COMP_* env vars |
| Nushell | `nushell.Patch()` | Strips leading quotes (`"`, `'`) and backticks from args using `shlex.Split` |
| Clink | `cmd_clink.Patch()` | Re-lexes `CARAPACE_COMPLINE` with carapace-shlex; strips redirects |

All other shells pass arguments directly to `traverse()` without patching.

## Quoting Complexity Comparison

| Shell | Quoting approach | Complexity |
|-------|-----------------|------------|
| Bash | Two modes: unquoted (`escapingReplacer`, backslash-escape) and double-quoted (`escapingQuotedReplacer`) | Medium |
| Zsh | 5-state machine (DEFAULT, QUOTING_ESCAPING, QUOTING, FULL_QUOTING_ESCAPING, FULL_QUOTING) | High |
| Fish | None (tab-separated, no quoting in output) | None |
| Elvish | None (JSON output, elvish handles quoting on insertion) | None |
| Nushell | Double-quote metacharacters; tilde kept outside quotes | Low |
| Xonsh | Python single-quoted `'value'` or raw strings `r'value'` | Low |
| PowerShell | Single-quote wrapping for special chars | Low |
| Oil | None (raw values) | None |
| Tcsh | Backslash-escape via `quoter` replacer | Medium |

## Related Skills

- **references/shell-bash.md** — bash integration deep dive
- **references/shell-bash-ble.md** — bash BLE integration deep dive
- **references/shell-oil.md** — oil integration deep dive
- **references/shell-elvish.md** — elvish integration deep dive
- **references/shell-fish.md** — fish integration deep dive
- **references/shell-nushell.md** — nushell integration deep dive
- **references/shell-powershell.md** — PowerShell integration deep dive
- **references/shell-xonsh.md** — xonsh integration deep dive
- **references/shell-zsh.md** — zsh integration deep dive
- **carapace-bin skill** (carapace-bin repo) — installation and shell integration (user-facing)
- **references/traverse.md** — the completion engine that produces Actions before formatting
- **references/style.md** — how styles are resolved before shell rendering
