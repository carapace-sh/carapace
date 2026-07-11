# Capturing Zsh Completions Externally

In-depth reference for the **completion capture** technique — spawning a headless zsh via `zpty`, intercepting `compadd` to capture matches into arrays instead of inserting them, and parsing the output externally. This is the mechanism used by [carapace-bridge](https://github.com/carapace-sh/carapace-bridge) (vendoring [Valodim/zsh-capture-completion](https://github.com/Valodim/zsh-capture-completion)) to bridge zsh-native completions into carapace.

## Overview

The normal completion flow inserts matches into the edit buffer via ZLE. The capture technique instead:

1. Spawns an interactive zsh inside a pseudo-terminal (`zpty`)
2. Sources an init script that overrides `compadd` to **capture** matches into arrays (`-A`/`-D`) rather than insert them
3. Sends the command line + TAB key to trigger completion
4. Reads the captured matches back from the pty
5. Parses them on the consumer side (e.g., Go code)

This allows an external process to obtain zsh completion candidates as plain text, without running inside the completion system itself.

## The zpty Spawn

```zsh
zmodload zsh/zpty
zpty z zsh -f -i          # spawn shell: -f (no rcs), -i (interactive)
```

`zpty` creates a pseudo-terminal pair. The spawned shell thinks it's interactive (required for ZLE and completion), but `-f` ensures no user startup files load. All communication happens through `zpty -w` (write) and `zpty -r` (read).

### Why Interactive Mode Is Required

The completion system only operates when ZLE is active, and ZLE only activates in interactive shells. A non-interactive shell (`zsh -c`) cannot run `compadd` or trigger TAB completion. The `zpty` trick gives us an interactive shell whose I/O we fully control.

## The Init Script (Sourced Into the Pty Shell)

The init script is written to a process substitution file and sourced via `zpty -w z source $1`. It uses `setopt rcquotes` in the outer script so that doubled single quotes (`''`) inside the heredoc become literal single quotes after expansion.

### Prompt Suppression

```zsh
PROMPT=
```

No prompt — the consumer parses raw output, so any prompt text would corrupt the parse.

### Completion System Initialization

```zsh
oldfpath="$fpath"
autoload -U compinit && compinit -d "${CARAPACE_BRIDGE_CONFIG_HOME:-$HOME/.config}/carapace/bridge/zsh/.zcompdump_capture"
source "${CARAPACE_BRIDGE_CONFIG_HOME:-$HOME/.config}/carapace/bridge/zsh/.zshrc"
[[ "$oldfpath" != "$fpath" ]] && compinit # second call to adopt any changes to fpath
```

- `autoload -U compinit` — the `-U` flag prevents alias expansion when loading, which is important in environments where `compinit` may be aliased.
- `compinit -d <path>` — uses a dedicated dump file to avoid clobbering the user's `~/.zcompdump`.
- Sourcing `.zshrc` allows user customization (adding fpath entries, defining completion functions).
- The conditional second `compinit` re-scans if `.zshrc` added directories to `$fpath`. Without it, completion functions added via `fpath` wouldn't be registered.

### Key Bindings

```zsh
bindkey '^M' undefined    # Enter — prevent command execution
bindkey '^J' undefined    # Ctrl-J (Linefeed) — prevent command execution
bindkey '^I' complete-word  # TAB — trigger completion
```

Binding Enter and Linefeed to `undefined` ensures that even if a newline leaks into the pty, no command runs. TAB is bound to `complete-word` to trigger the completion system.

### Null-Byte Framing

```zsh
null-line () {
    echo -E - $'\0'
}
compprefuncs=( null-line )
comppostfuncs=( null-line exit )
```

`compprefuncs` and `comppostfuncs` are arrays of functions called before and after completion output. The null-byte sentinel delimits the start and end of completion output, making parsing reliable — null bytes cannot appear in completion text.

The `exit` in `comppostfuncs` terminates the pty shell after one completion attempt, causing the `zpty -r` read loop to end.

### Zstyle Configuration

```zsh
zstyle ':completion:*' list-grouped false
zstyle ':completion:*' insert-tab false
zstyle ':completion:*' list-separator ''
```

- `list-grouped false` — don't group matches by tag in the output; each match appears on its own line.
- `insert-tab false` — don't insert a literal tab when completing on an empty line (would corrupt output).
- `list-separator ''` — no separator between display text and descriptions, simplifying parsing on the consumer side.

## The compadd Override (Core Hook)

This is the heart of the capture mechanism. The init script defines a function named `compadd` that shadows the builtin:

### Delegation for Capture Calls

```zsh
compadd () {
    if [[ ${@[1,(i)(-|--)]} == *-(O|A|D)\ * ]]; then
        builtin compadd "$@"
        return $?
    fi
```

If the caller already passed `-O`, `-A`, or `-D` (the array-capture options), delegate to `builtin compadd` and return. This prevents infinite recursion and respects callers that use these flags directly. The pattern `${@[1,(i)(-|--)]}` slices from argument 1 up to the first `-` or `--`, checking if any of the capture flags appear.

### Description Extraction

```zsh
    if (( $@[(I)-d] )); then
        __tmp=${@[$[${@[(i)-d]}+1]]}
        if [[ $__tmp == \(* ]]; then
            eval "__dscr=$__tmp"
        else
            __dscr=( "${(@P)__tmp}" )
        fi
    fi
```

The `-d` option provides display/description strings. It can be either:
- An inline `()` array literal — extracted via `eval`
- A parameter name — expanded with `${(@P)__tmp}` (indirect parameter expansion)

`zparseopts` is not used here because combined option parameters like `-default-` confuse it (the comment notes this).

### Match Capture via -A and -D

```zsh
    builtin compadd -A __hits -D __dscr "$@"
```

This is the key trick. By injecting `-A __hits -D __dscr` into the `builtin compadd` call:
- `-A array` — matching strings generated by the completion code are stored in `__hits` (after match specs are applied)
- `-D array` — non-matching words are removed from `__dscr`

The matches are **captured** into arrays rather than being inserted into the edit buffer. zsh's matching logic (prefix matching, matcher specs, etc.) is applied for free.

### Prefix/Suffix Extraction

```zsh
    typeset -A apre hpre hsuf asuf
    zparseopts -E P:=apre p:=hpre S:=asuf s:=hsuf
```

Extracts the four prefix/suffix options from the original compadd arguments:
- `-P` (prefix) → `apre`
- `-p` (hidden prefix) → `hpre`
- `-S` (suffix) → `asuf`
- `-s` (hidden suffix) → `hsuf`

**Why `typeset -A` (associative array) matters**: With `:=` in `zparseopts`, an associative array stores `key=value` pairs (e.g., `apre[-P]=hello`). When expanded unquoted (`$apre`), only the values are emitted — so `$apre` yields just the prefix string, not the option flag. A regular array (`typeset -a`) would store `(-P hello)` as two elements, and `$apre` would expand to both.

The `-E` flag tells `zparseopts` to parse only the option specs and leave remaining arguments in place — it does not consume `$@`.

### Directory Suffix Detection

```zsh
    if [[ -z $hsuf && "${${@//-default-/}% -# *}" == *-[[:alnum:]]#f* ]]; then
        dirsuf=1
    fi
```

Heuristically detects if `-f` (filename mode) was passed to `compadd`. If so, directories get a trailing `/` appended. The pattern strips `-default-` pseudo-options and checks for `f` in combined option strings. This is a "half-assed" emulation of compadd's built-in directory suffix behavior.

### Output Format

```zsh
    for i in {1..$#__hits}; do
        (( dirsuf )) && [[ -d $__hits[$i] ]] && dsuf=/ || dsuf=
        (( $#__dscr >= $i )) && dscr=" -- ${${__dscr[$i]}##$__hits[$i] #}" || dscr=
        echo -E - $IPREFIX$apre$hpre$__hits[$i]$dsuf$hsuf$asuf$dscr
    done
```

Each match is echoed as one line:
- `$IPREFIX` — initial prefix (ignored for matching but inserted)
- `$apre$hpre` — prefix and hidden-prefix from `-P`/`-p`
- `$__hits[$i]` — the match itself
- `$dsuf` — directory suffix (`/`) if applicable
- `$hsuf$asuf` — hidden-suffix and suffix from `-s`/`-S`
- `$dscr` — description, formatted as ` -- <description>` (with the match text stripped from the front of the display string)

`echo -E -` outputs the string literally — no escape interpretation (`-E`) and no dash-flag parsing (`-`), preventing mangling of values that start with `-`.

## The Read Loop (Consumer Side)

```zsh
zpty -w z "$*"$'\t'

integer tog=0
while zpty -r z; do :; done | while IFS= read -r line; do
    if [[ $line == *$'\0\r' ]]; then
        (( tog++ )) && return 0 || continue
    fi
    (( tog )) && echo -E - $line
done

return 0
```

1. `zpty -w z "$*"$'\t'` — writes the command line arguments followed by a TAB character into the pty, triggering completion.
2. `while zpty -r z; do :; done` — drains the pty output, piping it linewise to the inner loop.
3. The toggle (`tog`) tracks the two null-byte markers: the first sets `tog=1` (start of completions), the second exits with `return 0`.
4. Lines between the markers are echoed to stdout for the external consumer.
5. `return 0` — exits cleanly. (The original upstream used `return 2`, which was a bug — see below.)

### The Exit Code Bug (Fixed in carapace-bridge)

The original upstream script ended with `return 2`. Since `return 0` inside the piped `while` loop only exits the subshell (not the script), `return 2` always executed unconditionally. This meant the script always exited with code 2.

When used with carapace's `ActionExecCommand` (which treats non-zero exit codes as errors and skips the output callback), this silently broke the zsh bridge — the captured completions were never parsed. The fix changes `return 2` to `return 0`.

## Consumer-Side Parsing (Go)

The carapace-bridge Go code in `pkg/actions/bridge/zsh.go` consumes the output:

1. Splits output on `\r\n` (pty output uses CRLF line endings)
2. Strips ANSI escape codes (some completions emit color codes)
3. Unquotes zsh quoting (backslash-escaped special characters like `\<`, `\>`, `\{`, etc.)
4. Splits each line on ` -- ` to separate the display value from the description
5. Returns `carapace.ActionValuesDescribed` with the parsed pairs

The unquoter handles the zsh quoting that `compadd` applies to completion candidates with special characters.

## Key Mechanisms Summary

| Mechanism | Role |
|-----------|------|
| `zpty` | Spawn headless interactive zsh for completion |
| `bindkey '^M' undefined` | Prevent command execution in pty |
| `compprefuncs`/`comppostfuncs` | Null-byte framing around completion output |
| `compadd` override | Intercept all match generation |
| `builtin compadd -A -D` | Capture matches into arrays instead of inserting |
| `zparseopts -E` with `typeset -A` | Extract prefix/suffix options, expand to values only |
| `echo -E -` | Literal output, no escape interpretation |
| `zstyle list-grouped false` | One match per line, no grouping |
| `zstyle list-separator ''` | No separator between display and description |
| `setopt rcquotes` | Allow `''` inside single-quoted heredoc for literal quotes |
| `autoload -U compinit` | Load compinit without alias expansion |
| Conditional second `compinit` | Adopt fpath changes from sourced `.zshrc` |

## References

- [Valodim/zsh-capture-completion](https://github.com/Valodim/zsh-capture-completion) — original upstream
- [carapace-bridge zsh bridge](https://github.com/carapace-sh/carapace-bridge/blob/main/pkg/actions/bridge/zsh.go) — Go consumer of capture output
- [zsh Manual: zsh/zpty](https://zsh.sourceforge.io/Doc/Release/Zsh-Modules.html#The-zsh_002fzpty-Module) — pseudo-terminal module
- [zsh Manual: compadd](https://zsh.sourceforge.io/Doc/Release/Completion-Widgets.html#Completion-Builtin-Commands) — `-A`, `-D`, `-O` capture options
- [zsh Manual: zparseopts](https://zsh.sourceforge.io/Doc/Release/Zsh-Modules.html#The-zsh_002fzutil-Module) — option parsing utility
