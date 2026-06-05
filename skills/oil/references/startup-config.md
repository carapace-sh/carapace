# Oil Shell Startup and Configuration

In-depth reference for Oil shell's startup files, configuration, environment variables, shell options, and how OSH and YSH differ in their initialization.

## Startup File Decision Flow

### Interactive OSH Shell

1. `~/.config/oils/oshrc` — the **only** initialization file
2. All files in `~/.config/oils/oshrc.d/` — additional RC files

### Interactive YSH Shell

1. `~/.config/oils/yshrc` — the **only** initialization file
2. All files in `~/.config/oils/yshrc.d/` — additional RC files

### Non-Interactive Shell

- No startup files are read by default
- Use `--eval FILE` to evaluate a file on startup

### Disabling Startup Files

| Flag | Effect |
|------|--------|
| `--norc` | Disable startup files |
| `--rcfile /dev/null` | Explicitly disable the startup file |
| `--rcdir /dev/null` | Explicitly disable the RC directory |

### Alternative RC File/Directory

| Flag | Effect |
|------|--------|
| `--rcfile FILE` | Use a custom RC file instead of the default |
| `--rcdir DIR` | Use a custom RC directory instead of the default |

## Key Difference from Bash

Oil **intentionally avoids** the confusing multi-file initialization sequence of bash:

| Shell | Startup Files |
|-------|---------------|
| bash (login) | `/etc/profile` → `~/.bash_profile` OR `~/.bash_login` OR `~/.profile` |
| bash (non-login) | `~/.bashrc` |
| bash (non-interactive) | `$BASH_ENV` |
| OSH | `~/.config/oils/oshrc` + `~/.config/oils/oshrc.d/*` |
| YSH | `~/.config/oils/yshrc` + `~/.config/oils/yshrc.d/*` |

If you want bash's startup files, simply `source` them in your `oshrc`:

```bash
# ~/.config/oils/oshrc
source /etc/profile
source ~/.bashrc
```

## Configuration Directory

```bash
mkdir -p ~/.config/oils       # for oshrc and yshrc
mkdir -p ~/.local/share/oils  # for osh_history
```

### Symlink Convenience

Symlink `~/.config/oils/oshrc` to `~/.oshrc` for convenience:

```bash
ln -s ~/.config/oils/oshrc ~/.oshrc
```

## File Purpose

| File | When Read | Typical Contents |
|------|-----------|------------------|
| `~/.config/oils/oshrc` | Interactive OSH startup | Aliases, functions, shell options, completions, prompt |
| `~/.config/oils/oshrc.d/*` | Interactive OSH startup | Drop-in configuration snippets |
| `~/.config/oils/yshrc` | Interactive YSH startup | YSH functions, options, prompt |
| `~/.config/oils/yshrc.d/*` | Interactive YSH startup | Drop-in configuration snippets |
| `~/.local/share/oils/osh_history` | N/A (written by shell) | Command history |

## Typical oshrc

```bash
# Exit if not interactive (defensive)
[[ $- != *i* ]] && return

# Source bash compatibility (optional)
# source /etc/profile

# Environment
export PATH="$HOME/bin:$PATH"
export EDITOR=vim

# Aliases
alias ll='ls -la'
alias la='ls -A'

# Shell options
shopt -s strict:all    # Stricter error handling
shopt -s nullglob      # Unmatched globs expand to nothing
shopt -s dotglob       # Include dotfiles in globs

# History
export HISTSIZE=10000
export HISTCONTROL=ignoredups:erasedups

# Prompt
PS1='\u@\h:\w\$ '

# Completion UI
OILS_COMP_UI=nice      # or 'minimal'

# Completions
source /usr/share/bash-completion/bash_completion 2>/dev/null

# Vi mode (optional)
# set -o vi
```

## Typical yshrc

```bash
# YSH prompt
func renderPrompt(io) {
  return (io.promptVal('s') ++ ' ')
}

# Shell options
shopt --set ysh:upgrade   # Enable YSH features

# Completion UI
OILS_COMP_UI=nice
```

## YSH vs OSH at Startup

`bin/ysh` is the same as `bin/osh` with the `ysh:all` option group set:

```bash
osh -o ysh:all -c 'echo hi'  # Same as YSH
```

| Aspect | OSH | YSH |
|--------|-----|-----|
| Startup file | `~/.config/oils/oshrc` | `~/.config/oils/yshrc` |
| Startup directory | `~/.config/oils/oshrc.d/` | `~/.config/oils/yshrc.d/` |
| Default strictness | Compatible with bash | `strict:all` + `ysh:all` |
| Variable declaration | `local`, `readonly` | `var`, `const`, `setvar` |
| Arithmetic | `$(( ))`, `(( ))` | Python-like expressions |
| Control flow | `then fi`, `do done` | `{ }` braces |
| Prompt | `PS1` with escapes | `renderPrompt()` function |

## Command Line Options

### POSIX-Standard Flags

```bash
osh -o errexit -c 'false'    # Enable errexit
ysh -n myfile.ysh            # Parse check mode
ysh +o errexit -c 'false'   # Disable errexit
```

### Oil-Specific Flags

| Flag | Description |
|------|-------------|
| `--eval FILE` | Evaluate a file on startup (like `source`) |
| `--eval-pure FILE` | Like `--eval` but disallows I/O ("pure mode") |
| `--norc` | Disable startup files |
| `--rcfile FILE` | Use custom RC file |
| `--rcdir DIR` | Use custom RC directory |
| `--headless` | Start in headless mode (UI decoupled) |
| `--debug-file FILE` | Print internal debug logs to file or FIFO |
| `--xtrace-to-debug-file` | Send `set -o xtrace` output to debug file |
| `-n` | Parse but don't execute; print AST |
| `--ast-format FMT` | AST output format (`text` or `text-abbrev`) |
| `--location-str STR` | Use string for error message display |
| `--location-start-line N` | Line number offset for error messages |
| `--tool TOOL` | Run a tool instead of shell (`cat-em` or `syntax-tree`) |

### Multiple --eval Files

```bash
ysh --eval one.ysh --eval two.ysh -c 'echo hi'  # Run 2 files first
```

## Environment Variables

| Variable | Purpose |
|----------|---------|
| `OILS_HIJACK_SHEBANG` | Path to a shell to use instead of the shebang-specified shell |
| `OILS_COMP_UI` | Completion UI: `minimal` (default) or `nice` |
| `OSH_DEBUG_DIR` | Directory for per-process debug logs (`$PID-osh.log`) |
| `PS1`, `PS2`, `PS4` | Prompt strings (may be unintentionally inherited) |

### OILS_HIJACK_SHEBANG

Makes all recursively-run scripts use OSH rather than `/bin/sh`:

```bash
OILS_HIJACK_SHEBANG=osh osh myscript.sh
```

This is useful for testing whether scripts work under OSH without modifying their shebangs.

### PS1 Inheritance

If running OSH from bash or zsh, `PS1` may be unintentionally inherited:

```bash
PS1='' bin/osh    # Clear PS1 before starting OSH
```

### Locale Issues

On some distros (e.g., Arch Linux), `$LANG` may not get set without `/etc/profile`:

```bash
# Add to oshrc
source /etc/profile
```

## Shell Options

### POSIX Shell Options

| Option | Flag | Effect |
|--------|------|--------|
| `errexit` | `-e` | Exit on command failure |
| `nounset` | `-u` | Error on unset variable reference |
| `pipefail` | | Pipeline fails if any command fails |
| `noglob` | `-f` | Disable filename expansion |
| `noclobber` | `-C` | Don't overwrite existing files with `>` |
| `xtrace` | `-x` | Trace command execution |
| `errtrace` | `-E` | ERR trap inherited by functions |
| `emacs` | | Use emacs editing mode |
| `vi` | | Use vi editing mode |

### Bash-Compatible Options

| Option | Effect |
|--------|--------|
| `inherit_errexit` | Command substitution inherits `errexit` |
| `nullglob` | Unmatched globs expand to nothing |
| `failglob` | Error if glob doesn't match |
| `dotglob` | Include dotfiles in glob expansion |
| `expand_aliases` | Expand aliases (on by default in interactive) |
| `lastpipe` | Last pipeline part runs in current shell (default in OSH) |
| `progcomp` | Enable programmable completion |

### Oils-Specific Options

| Option | Default | Effect |
|--------|---------|--------|
| `dashglob` | on (OSH), off (YSH) | Include files starting with `-` in globs |
| `strict_arith` | on | Strings that don't look like integers cause fatal errors |
| `strict_argv` | off | Empty `argv` arrays disallowed |
| `strict_array` | off | No implicit string↔array conversions |
| `strict_control_flow` | off | `break`/`continue` outside loops are fatal |
| `strict_errexit` | off | Warns when errors would be lost |
| `strict_word_eval` | off | More word evaluation errors are fatal |
| `strict_nameref` | off | Invalid variable names in namerefs are fatal |
| `strict_tilde` | off | Failed tilde expansions cause errors |
| `strict_glob` | off | Parse the sublanguage more strictly |
| `eval_unsafe_arith` | off | Allow dynamically parsed `a[$(echo 42)]` |
| `sigpipe_status_ok` | off | Status 141 (SIGPIPE) → 0 in pipelines |
| `simple_word_eval` | off | No implicit splitting, static globbing (YSH) |
| `xtrace_rich` | off | Hierarchical and process tracing |
| `rewrite_extern` | on (non-interactive) | Transparent rewriting of external commands to builtins |

### Option Groups

| Group | Includes | Purpose |
|-------|----------|---------|
| `strict:all` | All `strict_*` options | Disallow problematic shell constructs |
| `ysh:upgrade` | `parse_at`, `parse_brace`, `parse_paren`, `parse_proc`, `simple_word_eval`, `xtrace_rich`, etc. | Enable YSH features while maintaining compatibility |
| `ysh:all` | `strict:all` + `ysh:upgrade` + additional | Full YSH language |

### Setting Option Groups

```bash
shopt -s strict:all       # Enable all strict options
shopt -s ysh:upgrade      # Enable YSH features
shopt -s ysh:all          # Full YSH language
```

## Debugging

### Debug File

```bash
# Create a FIFO debug file
mkfifo _tmp/debug
osh --debug-file _tmp/debug

# In another window
cat _tmp/debug
```

### Per-Process Debug Logs

```bash
OSH_DEBUG_DIR=/tmp/osh-debug osh
# Creates /tmp/osh-debug/$PID-osh.log for every shell process
```

### Xtrace to Debug File

```bash
osh --debug-file _tmp/debug --xtrace-to-debug-file
# set -o xtrace output goes to debug file instead of stderr
```

### AST Dump

```bash
osh -n -c 'ls | wc -l'                    # Pretty-print AST
osh -n --ast-format text -c 'ls | wc -l'  # Unabridged format
```

## Pitfalls

1. **PS1 inheritance**: Running OSH from bash/zsh may inherit `PS1` — clear it first: `PS1='' osh`
2. **Locale not set**: On some distros, add `source /etc/profile` to `oshrc`
3. **Programs using `eval`**: Some programs (starship, zoxide) use `eval` in ways that may need compatibility adjustments — see [OSH Compatibility Tips](https://github.com/oils-for-unix/oils/wiki/OSH-Compatibility-Tips)
4. **No login shell distinction**: OSH doesn't have separate login/non-login startup files — everything goes through `oshrc`
5. **OILS_COMP_UI checked once**: The completion UI variable is only checked at shell initialization, not dynamically

## References

- [Oil Shell Getting Started](https://oils.pub/release/latest/doc/getting-started.html)
- [Oil Shell Front End Reference](https://oils.pub/release/latest/doc/ref/chap-front-end.html)
- [Oil Shell Options Reference](https://oils.pub/release/latest/doc/ref/chap-option.html)
- [YSH vs. Shell](https://oils.pub/release/latest/doc/ysh-vs-shell.html)
- [Known Differences](https://oils.pub/release/latest/doc/known-differences.html)

## Related Skills

- **bash skill → references/startup.md** — bash startup files and configuration
