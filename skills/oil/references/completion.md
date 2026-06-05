# Oil Shell Programmable Completion

In-depth reference for Oil shell's (OSH/YSH) programmable completion system — the builtins, variables, registration, quoting model, and how it differs from bash.

## The Completion Flow

When the user presses TAB, Oil:

1. Reads the current line buffer and cursor position from GNU Readline
2. Parses the command line up to the cursor using the **OSH parser** (not an ad-hoc parser)
3. Checks for shell language completions first (variables, tilde, redirects)
4. Statically evaluates words to build `partial_argv` (the `COMP_ARGV` array)
5. Looks up the completion spec for the command
6. Invokes the registered function, command, or action
7. Reads `COMPREPLY` for completion candidates
8. Passes candidates to the completion display (minimal or nice UI)

### Spec Lookup Order

Oil searches for a completion spec in this order:

1. **Exact match** — `lookup[argv0]`
2. **Basename match** — `lookup[basename(argv0)]` (e.g., `/usr/bin/git` → `git`)
3. **Glob pattern match** — linear search through registered patterns (Oil extension, not in bash)
4. **Alias expansion** — if the command is an alias, check the expansion result
5. **Fallback** — `__fallback` pseudo-command (registered with `complete -D`)
6. **First-word** — `__first` pseudo-command for empty line completion (registered with `complete -E`)

### Match Generation Order

Within a spec, matches are generated in this order:

1. Actions from `-A` options (e.g., `-A file`, `-A command`)
2. Filename patterns from `-G globpat`
3. Words from `-W wordlist` (split by IFS, then expanded)
4. Shell function via `-F function` (populates `COMPREPLY`)
5. External command via `-C command` (stdout becomes completions)

After generation:

6. `-X filterpat` removes matching completions (or non-matching if `!` prefix)
7. `-P prefix` and `-S suffix` are applied
8. Fallback options are tried if no matches: `-o dirnames`, `-o plusdirs`, `-o bashdefault`, `-o default`

## COMP_* Variables

When a completion function is invoked, Oil sets these variables:

| Variable | Type | Status | Description |
|----------|------|--------|-------------|
| `COMP_ARGV` | Array | **Preferred** | Partial command arguments to complete (OSH extension) |
| `COMPREPLY` | Array | Required | User-defined completion functions fill this with candidates |
| `COMP_WORDS` | Array | Discouraged | Words split by `:` and `=` for bash compatibility |
| `COMP_CWORD` | Integer | Discouraged | Index into `COMP_WORDS` of the cursor word |
| `COMP_LINE` | String | Discouraged | The entire command line as a string |
| `COMP_POINT` | Integer | Discouraged | Cursor position within `COMP_LINE` |
| `COMP_WORDBREAKS` | String | Discouraged | Characters that word splitting is performed on |

### COMP_ARGV vs COMP_WORDS

The fundamental difference: `COMP_ARGV` contains **argv entries** (parsed arguments), while `COMP_WORDS` contains **words split at `:` and `=`** for bash compatibility.

```bash
# For a command like: git push origin:main --opt=val
# COMP_ARGV:  (git push origin:main --opt=val)  — 4 entries
# COMP_WORDS: (git push origin main --opt val)  — 6 entries (split at : and =)
```

New completion scripts should use `COMP_ARGV`. The `compadjust` builtin bridges the gap for bash-completion scripts that expect `COMP_WORDS`-style splitting.

### COMP_WORDBREAKS

Default value: `" \t\n\"'@><=;|&(:`. These characters split the "word being completed" in `COMP_WORDS`. Unlike bash, Oil **discourages** use of `COMP_WORDBREAKS` — the `compadjust` builtin handles word breaking explicitly.

## The `complete` Builtin

```bash
complete [-abcdefgjksuv] [-o comp-option] [-DEI] [-A action]
         [-G globpat] [-W wordlist] [-F function] [-C command]
         [-X filterpat] [-P prefix] [-S suffix] name [name ...]
complete -pr [-DEI] [name ...]
```

### Action Options (`-A action` / short flags)

| Action | Short | Description |
|--------|-------|-------------|
| `alias` | `-a` | Alias names |
| `arrayvar` | | Array variable names |
| `binding` | | Readline key binding names |
| `builtin` | `-b` | Shell builtin commands |
| `command` | `-c` | Command names from PATH |
| `directory` | `-d` | Directory names |
| `disabled` | | Disabled shell builtins |
| `enabled` | | Enabled shell builtins |
| `export` | `-e` | Exported shell variables |
| `file` | `-f` | File and directory names |
| `function` | | Shell function names |
| `group` | `-g` | Group names |
| `helptopic` | | Help topics |
| `hostname` | | Hostnames |
| `job` | `-j` | Job names |
| `keyword` | `-k` | Shell reserved words |
| `running` | | Running jobs |
| `service` | `-s` | Service names |
| `setopt` | | Valid `set -o` arguments |
| `shopt` | | `shopt` option names |
| `signal` | | Signal names |
| `stopped` | | Stopped jobs |
| `user` | `-u` | User names |
| `variable` | `-v` | All shell variables |

### Completion Options (`-o comp-option`)

| Option | Description |
|--------|-------------|
| `bashdefault` | Fall back to Bash default completions if no matches |
| `default` | Use Readline's default filename completion if no matches |
| `dirnames` | Attempt directory name completion if no matches |
| `filenames` | Tell Readline completions are filenames (add slash to dirs, quote special chars) |
| `fullquote` | Quote all completed words |
| `noquote` | Don't quote completed words |
| `nosort` | Don't sort completions alphabetically |
| `nospace` | Don't append space after completion |
| `plusdirs` | After generating matches, also add directory completions |

### Special Registration Options

| Option | Description |
|--------|-------------|
| `-D` | Apply to "default" command completion (commands with no spec) |
| `-E` | Apply to "empty" command completion (blank line) |
| `-I` | Apply to initial word (command name completion) |
| `-p` | Print existing compspecs |
| `-r` | Remove compspec for specified names |

### Glob Pattern Support (Oil Extension)

Oil supports glob patterns for command names in `complete`, unlike bash:

```bash
complete -F _git_completion 'git*'   # matches git, gitk, git-status, etc.
```

Pattern matching uses `fnmatch` — linear search through registered patterns.

### `-F function` Details

When the function is called:

- `$1` — name of the command whose arguments are being completed
- `$2` — word being completed (from `COMP_ARGV`, not raw command string)
- `$3` — word preceding the word being completed

The function populates `COMPREPLY`. When it finishes, Oil reads `COMPREPLY` for candidates.

**Exit status 124**: If the function returns 124, Oil restarts the completion from the beginning. This enables dynamic loading — the function can source a completion file and return 124 to re-trigger with the newly registered spec. This is compatible with bash's behavior.

## The `compgen` Builtin

```bash
compgen [-V varname] [option] [word]
```

Generates completion matches and writes them to stdout (or stores in `varname` with `-V`). Accepts the same options as `complete` except `-p`, `-r`, `-D`, `-E`, `-I`.

Can also be used **outside a completion function** in scripts, unlike some bash implementations.

### Common Usage

```bash
compgen -W "start stop status" -- "$cur"     # Word list
compgen -f -- "$cur"                          # Filenames
compgen -d -- "$cur"                          # Directories
compgen -A command -- "$cur"                  # Commands
compgen -A variable -- "$cur"                # Variables
```

## The `compopt` Builtin

```bash
compopt [-o option] [+o option] [-DEI] [name]
```

Modifies completion options dynamically — either for a named command or for the currently-executing completion (if no name given).

- `-o option` — enable option
- `+o option` — disable option
- Only works when `comp_state.currently_completing` is True (within a completion function)
- Returns error if called outside a completion function

## The `compadjust` Builtin (OSH Extension)

```bash
compadjust [-F func] [-a words_var] [-c cur_var] [-p prev_var] [-w cword_var] [-s split_var] -- delimiter...
```

Adjusts `COMP_ARGV` according to specified delimiters, and optionally sets variables:

| Variable | Description |
|----------|-------------|
| `cur` | Current word being completed |
| `prev` | Word preceding the current word |
| `words` | Array of words (after word-break splitting) |
| `cword` | Index of the current word in `words` |
| `split` | Whether word breaking occurred |

### How It Works

The `AdjustArg` state machine splits words at break characters:

```
Input:  "origin:main" with break_chars [':', '=']
Output: ['origin', 'main']
```

This is derived from a cleanup of the `bash-completion` project's `_init_completion` function. It makes it easier to run existing bash-completion scripts on OSH.

### Usage with bash-completion

```bash
_my_completion() {
    compadjust -- cur prev words cword  # sets $cur, $prev, $words, $cword
    # Now $cur, $prev, etc. work like in bash-completion
    case $prev in
        --output) compgen -f -- "$cur" ;;
        *) COMPREPLY=($(compgen -W "start stop" -- "$cur")) ;;
    esac
}
complete -F _my_completion mycmd
```

## The `compexport` Builtin (OSH Extension)

```bash
compexport -c 'command string' [-b begin] [-e end] [-f format]
```

Completes an **entire shell command string**, not just individual arguments. This is unique to Oil and enables completing shell language constructs.

### Examples

```bash
compexport -c 'echo $H'     # completes variables like $HOME
compexport -c 'ha'          # completes builtins like 'hay' and external commands
```

### Output Formats

| Format | Description |
|--------|-------------|
| `jlines` | JSON Lines — one JSON-encoded string per line |
| `tsv8` | Tab-separated values with J8 encoding |

The `-b` and `-e` flags specify cursor begin and end positions within the command string.

## The Quoting Model: Key Difference from Bash

The most fundamental difference between Oil and bash completion is the **quoting responsibility**:

> "The completion API is modeled after the bash completion API. However, an incompatibility is that it deals with `argv` entries and not command strings. OSH moves the **responsibility for quoting** into the shell. Completion plugins should not do it."

### What This Means

| Aspect | Bash | Oil |
|--------|------|-----|
| Input to function | Raw command string (`COMP_LINE`) | Parsed argv entries (`COMP_ARGV`) |
| Quoting/dequoting | Plugin must handle | Shell handles internally |
| Word splitting | Plugin must re-parse | Already split into argv |
| `$2` value | Last segment after `COMP_WORDBREAKS` | Full argv entry |

### Practical Implications

1. **No need for `_comp_dequote`** — Oil handles dequoting before passing to the function
2. **No need for `_comp_quote`** — Oil handles quoting when inserting completions
3. **`COMPREPLY` entries are argv** — just set the values, don't quote them
4. **`COMP_ARGV` is already parsed** — no need to re-parse `COMP_LINE`

### Example

```bash
# Bash: must handle quoting manually
_my_func() {
    local cur
    _init_completion -n = || return
    # $cur is dequoted, but COMPREPLY must be quoted for filenames
    COMPREPLY=($(compgen -f -- "$cur"))
    # Readline handles quoting, but plugin must set -o filenames
}

# Oil: quoting is handled by the shell
_my_func() {
    compadjust -- cur prev words cword
    # $cur is already an argv entry, no dequoting needed
    COMPREPLY=($(compgen -f -- "$cur"))
    # Shell handles quoting when inserting into the line
}
```

## Shell Language Completion

Unlike bash, Oil provides **first-class completion for the shell language itself** using its parser:

| Construct | Example | Completion |
|-----------|---------|------------|
| Variable names | `echo $HOM` | Completes `$HOME`, `$HOSTNAME`, etc. |
| Braced variables | `echo ${HOM` | Completes `${HOME}`, `${HOSTNAME}`, etc. |
| Tilde expansion | `cd ~us` | Completes `~user` from passwd |
| Redirect targets | `cat > ` | Completes filenames |
| Commands in blocks | `if tr` | Completes `true` (works in `if`/`while`/`for`/`case`) |
| Command substitution | `$(gre` | Completes `grep` (bash-completion fails here) |

### Parser-as-Library Architecture

Oil uses its **parser as a library** for completion, unlike bash which uses ad-hoc parsers:

| Feature | Bash | Oil |
|---------|------|-----|
| History expansion | Partial/incorrect parser | Full parser — correctly handles `${x:-a b c}` |
| Variable completion | `echo $HOM<TAB>` works | Same, plus `${HOM<TAB>` also works |
| Command completion | Limited to simple statements | Works in `if`, `while`, `for`, `case` blocks |
| Completions in `$(...)` | Syntax error in bash-completion | Works correctly |

The parser produces a "trail" of words and tokens on incomplete input, which the completion system uses to determine what to complete.

## Completion UI: OILS_COMP_UI

The `OILS_COMP_UI` variable controls which completion display to use. It is checked **once at shell initialization**.

| Value | Description |
|-------|-------------|
| `minimal` | Approximates default GNU Readline behavior. No color, no terminal width dependency. Useful for browser builds. |
| `nice` | Fancy pager with horizontal scrolling prompt (no wrapping). Shows descriptions. Uses ANSI escape codes for cursor control. |

### Nice UI Features

- **Line erasure** — remembers how many lines were displayed and erases them before showing new completions
- **Horizontal scrolling** — `set horizontal-scroll-mode on` prevents line wrapping
- **Description display** — shows flag/builtin descriptions in yellow
- **Common prefix stripping** — strips common prefixes according to Oil's rules (not Readline's)
- **Multi-press detection** — tracks repeated TAB presses via hash of all matches, showing more lines progressively
- **Right-side info** — can show additional info on the right with reverse video

### Setting the UI

```bash
# In oshrc
OILS_COMP_UI=nice
```

## Headless Mode for Completion

Oil supports `--headless` mode where the UI is decoupled from the shell:

```bash
osh --headless
```

| UI Process Handles | Shell Process Handles |
|---|---|
| Auto-completion (using Oil for parsing) | Parsing and evaluating the language |
| History retrieval | Maintaining state (options, variables) |
| Cancelling commands in progress | |

Communication happens over a Unix domain socket using the **FANOS** protocol (File descriptors And Netstrings Over Sockets). Alternative UIs should use Oil for parsing rather than trying to parse shell syntax themselves.

## Exit Status 124 Retry Mechanism

When a completion function returns exit status 124, Oil restarts the completion from the beginning. This enables **dynamic loading**:

1. Completion function is called
2. Function calls `complete` to register new specs (e.g., lazy-loading)
3. Function returns 124
4. Oil checks if the command's spec was changed
5. If changed, re-fetches the spec and retries completion

```python
# Internal implementation (simplified)
status = self.cmd_ev.RunFuncForCompletion(self.func, argv)
commands_changed = self.comp_lookup.GetCommandsChanged()

if status == 124:
    cmd = os_path.basename(comp.first)
    if cmd in commands_changed:
        raise _RetryCompletion()
```

This is compatible with bash's exit-124 protocol used by bash-completion for lazy loading.

## Differences from Bash: Summary

| Feature | Bash | Oil |
|---------|------|-----|
| Primary completion variable | `COMP_WORDS` | `COMP_ARGV` (preferred) |
| Quoting responsibility | Plugin must handle | Shell handles internally |
| Argument processing | Command strings | argv entries |
| Shell language completion | Limited/ad-hoc | First-class via parser |
| `compadjust` builtin | Not available | OSH extension for bash-compat |
| `compexport` builtin | Not available | OSH extension for full command completion |
| Glob pattern specs | Not supported | Supported via `fnmatch` |
| Headless mode | Not available | Unique feature via FANOS |
| Completion UI | Readline only | `minimal` or `nice` (OILS_COMP_UI) |
| History expansion | Buggy on complex expressions | Accurate (uses full parser) |
| `compgen` outside function | Varies | Supported |

## References

- [Oil Shell Completion Documentation](https://oils.pub/release/latest/doc/completion.html)
- [Oil Shell Builtin Commands Reference](https://oils.pub/release/latest/doc/ref/chap-builtin-cmd.html)
- [Oil Shell Special Variables Reference](https://oils.pub/release/latest/doc/ref/chap-special-var.html)
- [Known Differences Between OSH and Other Shells](https://oils.pub/release/latest/doc/known-differences.html)
- [Oil Shell Blog: History and Completion](https://www.oilshell.org/blog/2020/01/history-and-completion.html)

## Related Skills

- **bash skill → references/completion.md** — bash programmable completion (shares `complete`/`compgen`/`compopt` builtins)
- **carapace-dev skill → references/shell-oil.md** — carapace-specific Oil integration (snippet, value formatting, nospace)
