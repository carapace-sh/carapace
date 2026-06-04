# Bash Startup Files and Configuration

In-depth reference for which files bash reads at startup, how shell options work, and where to install completions.

## Startup File Decision Flow

### Interactive Login Shell (or `--login` flag)

1. `/etc/profile` — system-wide (first)
2. **First found** of:
   - `~/.bash_profile`
   - `~/.bash_login`
   - `~/.profile`
3. On exit: `~/.bash_logout` (if it exists)

### Interactive Non-Login Shell

- Reads `~/.bashrc`
- `--norc` disables this
- `--rcfile file` uses a custom file instead of `~/.bashrc`

### Non-Interactive Shell (script execution)

- Looks for `BASH_ENV` environment variable
- If set, reads and executes that file: `if [ -n "$BASH_ENV" ]; then . "$BASH_ENV"; fi`
- PATH is **not** used to locate the file

### Invoked as `sh`

- **Login shell**: Reads `/etc/profile` then `~/.profile`
- **Interactive `sh`**: Looks for `ENV` variable and reads that file
- **Non-interactive `sh`**: Does not read any startup files
- `--rcfile` has no effect when invoked as `sh`

### POSIX Mode (`--posix` flag)

- Interactive shells expand `ENV` variable and read that file
- No other startup files are read

### Remote Shell Daemon (rshd/sshd)

- Reads `~/.bashrc` if it exists
- Does **not** read if invoked as `sh`

### Elevated Privileges (unequal UID/GID)

- No startup files are read unless `-p` flag is supplied
- Shell functions from environment are ignored
- Certain environment variables are ignored: `SHELLOPTS`, `BASHOPTS`, `CDPATH`, `GLOBIGNORE`

## File Purposes

| File | When Read | Typical Contents |
|------|-----------|------------------|
| `/etc/profile` | Login shells | System-wide PATH, umask, environment |
| `~/.bash_profile` | Bash login shells | User login initializations, source `.bashrc` |
| `~/.bash_login` | Bash login (if no `.bash_profile`) | Fallback for `.bash_profile` |
| `~/.profile` | sh login / Bash (if no `.bash_profile`/`.bash_login`) | Bourne-compatible environment |
| `~/.bashrc` | Interactive non-login shells | Aliases, functions, shell options, completions |
| `/etc/bash.bashrc` | Interactive shells (Debian-specific) | System-wide interactive settings |
| `~/.bash_logout` | Login shell exit | Cleanup, logout messages |

## Common Patterns

### Typical `~/.bash_profile`

```bash
# Source .bashrc for interactive shells
if [ -f ~/.bashrc ]; then
    . ~/.bashrc
fi

# Login-specific initializations
export PATH="$HOME/bin:$PATH"
```

This pattern ensures `.bashrc` is sourced for both login and non-login interactive shells.

### Typical `~/.bashrc`

```bash
# Exit if not interactive
[[ $- != *i* ]] && return

# Environment
export PATH="$HOME/bin:$PATH"
export EDITOR=vim

# Aliases
alias ll='ls -la'
alias la='ls -A'

# Shell options
shopt -s autocd
shopt -s cdspell
shopt -s checkwinsize
shopt -s histappend

# History
export HISTSIZE=10000
export HISTCONTROL=ignoredups:erasedups
export HISTIGNORE="ls:cd:exit"

# Prompt
PS1='\u@\h:\w\$ '

# Completions
if [ -f /usr/share/bash-completion/bash_completion ]; then
    . /usr/share/bash-completion/bash_completion
elif [ -f /etc/bash_completion ]; then
    . /etc/bash_completion
fi
```

## Pitfalls

1. **Order matters**: If `~/.bash_profile` exists, `~/.profile` is **ignored**
2. **Missing `.bashrc` sourcing**: Login shells won't get aliases/functions unless `.bash_profile` sources `.bashrc`
3. **BASH_ENV vs ENV**: `BASH_ENV` is for non-interactive bash; `ENV` is for POSIX mode and `sh` invocation
4. **Scripts ignore interactive configs**: Non-interactive shells only read `BASH_ENV`, not `.bashrc`
5. **SSH sessions**: SSH typically invokes login shells (`.bash_profile`), but `ssh host command` is non-interactive
6. **Error handling**: If a startup file exists but cannot be read, bash reports an error

## Shell Options

### `set` Options

| Option | Flag | Effect |
|--------|------|--------|
| `allexport` | `-a` | Mark variables for export on assignment |
| `braceexpand` | `-B` | Enable brace expansion |
| `emacs` | | Use emacs editing mode |
| `errexit` | `-e` | Exit on command failure |
| `errtrace` | `-E` | ERR trap inherited by functions |
| `functrace` | `-T` | DEBUG/RETURN traps inherited by functions |
| `hashall` | `-h` | Hash commands on first lookup |
| `histexpand` | `-H` | Enable `!` history expansion |
| `history` | | Enable command history |
| `ignoreeof` | | Don't exit on `^D` (require `exit`) |
| `keyword` | `-k` | All assignment statements in environment |
| `monitor` | `-m` | Enable job control |
| `noclobber` | `-C` | Don't overwrite existing files with `>` |
| `noexec` | `-n` | Read commands but don't execute |
| `noglob` | `-f` | Disable filename expansion |
| `notify` | `-b` | Report terminated background jobs immediately |
| `nounset` | `-u` | Error on unset variable reference |
| `onecmd` | `-t` | Exit after one command |
| `physical` | `-P` | Don't resolve symlinks in `cd`/`pwd` |
| `pipefail` | | Pipeline fails if any command fails |
| `posix` | | POSIX mode |
| `privileged` | `-p` | Don't process startup files |
| `verbose` | `-v` | Echo input lines |
| `vi` | | Use vi editing mode |
| `xtrace` | `-x` | Trace command execution |

### `shopt` Options

| Option | Default | Effect |
|--------|---------|--------|
| `autocd` | off | `cd` by typing directory name |
| `cdable_vars` | off | `cd` to variable name if directory doesn't exist |
| `cdspell` | off | Auto-correct minor `cd` typos |
| `checkwinsize` | on | Update LINES/COLUMNS after each command |
| `cmdhist` | on | Save multi-line commands as single history entry |
| `compat31` | off | Compatibility with bash 3.1 |
| `compat32` | off | Compatibility with bash 3.2 |
| `compat40` | off | Compatibility with bash 4.0 |
| `dirspell` | off | Auto-correct directory name typos in completion |
| `dotglob` | off | Include dotfiles in glob expansion |
| `execfail` | off | Non-interactive shell doesn't exit on exec failure |
| `expand_aliases` | on (interactive) | Expand aliases |
| `extdebug` | off | Extended debugging (affects DEBUG trap) |
| `extglob` | off | Extended glob patterns (`?()`, `*()`, `+()`, `@()`, `!()`) |
| `extquote` | on | Allow `$'...'` and `$"..."` in expansions |
| `failglob` | off | Error if glob pattern doesn't match |
| `force_fignore` | on | FIGNORE applies even if only match |
| `globasciiranges` | on | Glob ranges use ASCII order |
| `globstar` | off | `**` matches directories recursively |
| `gnu_errfmt` | off | GNU-style error messages |
| `histappend` | off | Append to (not overwrite) history file |
| `histreedit` | off | Re-edit failed history expansion |
| `histverify` | off | Don't execute history expansion immediately |
| `hostcomplete` | on | Complete hostnames containing `@` |
| `huponexit` | off | Send SIGHUP to all jobs on login shell exit |
| `inherit_errexit` | off | Command substitution inherits `errexit` |
| `interactive_comments` | on | Allow `#` comments in interactive shell |
| `lastpipe` | off | Last pipeline element runs in current shell (no job control) |
| `lithist` | off | Save multi-line commands with embedded newlines |
| `localvar_inherit` | off | Local variables inherit value from outer scope |
| `login_shell` | (set by shell) | Shell is a login shell |
| `mailwarn` | off | Warn if mail file was accessed since last check |
| `no_empty_cmd_completion` | off | Don't try completion on empty line |
| `nocaseglob` | off | Case-insensitive glob matching |
| `nocasematch` | off | Case-insensitive pattern matching in `case`/`[[` |
| `nullglob` | off | Unmatched globs expand to nothing |
| `progcomp` | on | Enable programmable completion |
| `progcomp_alias` | off | Try compspec for alias expansion result |
| `promptvars` | on | Expand variables in prompt strings |
| `restricted_shell` | (set by shell) | Shell is restricted |
| `shift_verbose` | off | `shift` prints error when shifting past end |
| `sourcepath` | on | `source` uses PATH to find files |
| `xpg_echo` | off | `echo` expands escape sequences by default |

## Completion Installation Locations

### Where to Install Completions

| Location | Scope | Priority |
|----------|-------|----------|
| `~/.local/share/bash-completion/completions/` | User | 1 (highest) |
| `$BASH_COMPLETION_USER_DIR/completions/` | User | 1 |
| `/usr/share/bash-completion/completions/` | System | 2 |
| `<prefix>/share/bash-completion/completions/` | Package | 3 |
| `$XDG_DATA_DIRS/bash-completion/completions/` | System | 4 |
| `/etc/bash_completion.d/` | Legacy | 5 (lowest) |

### When Completions Are Loaded

- **bash-completion package**: Loaded on demand via `_completion_loader` (default compspec with `-D`)
- **Sourced in `.bashrc`**: Loaded at shell startup
- **Sourced in `.bash_profile`**: Loaded at login

### Recommended Installation

For user-specific completions:

```bash
mkdir -p ~/.local/share/bash-completion/completions
cp mytool-completion.bash ~/.local/share/bash-completion/completions/mytool
```

For system-wide:

```bash
cp mytool-completion.bash /usr/share/bash-completion/completions/mytool
```

The file should be named after the command (no `.bash` extension required, but both work).

## Environment Variables

### Key Bash Variables

| Variable | Description |
|----------|-------------|
| `BASH_ENV` | File to source for non-interactive shells |
| `BASH_XTRACEFD` | File descriptor for `set -x` output |
| `CDPATH` | Search path for `cd` |
| `FIGNORE` | Suffixes to ignore in completion (`:.` separated) |
| `GLOBIGNORE` | Patterns to exclude from globbing |
| `HISTFILE` | History file path (default: `~/.bash_history`) |
| `HISTSIZE` | Max history entries in memory |
| `HISTFILESIZE` | Max history entries on disk |
| `HISTCONTROL` | History filtering: `ignoredups`, `ignorespace`, `erasedups`, `ignoreboth` |
| `HISTIGNORE` | Patterns to exclude from history |
| `HISTTIMEFORMAT` | Timestamp format for history entries |
| `HOSTFILE` | File for hostname completion (like `/etc/hosts`) |
| `INPUTRC` | Override default `.inputrc` location |
| `LANG` / `LC_ALL` / `LC_*` | Locale settings |
| `MAIL` / `MAILPATH` | Mail checking paths |
| `PROMPT_COMMAND` | Command(s) to run before displaying PS1 |
| `PROMPT_DIRTRIM` | Trim directory depth in `\w` prompt |
| `PS1` / `PS2` / `PS3` / `PS4` | Prompt strings |
| `SHELLOPTS` | Enabled `set` options (colon-separated) |
| `TIMEFORMAT` | Format for `time` output |
| `TMOUT` | Auto-logout timeout (seconds, 0 = disabled) |

### Carapace-Relevant Environment Variables

| Variable | Description |
|----------|-------------|
| `CARAPACE_LOG` | Enable carapace debug logging |
| `CARAPACE_SANDBOX` | JSON mock context (set by sandbox tests) |
| `CARAPACE_LENIENT` | Allow unknown flags |
| `CARAPACE_MATCH` | Set to `CASE_INSENSITIVE` for case-insensitive matching |
| `CARAPACE_NOSPACE` | Additional nospace suffixes |
| `CARAPACE_UNFILTERED` | Skip prefix filtering |
| `CARAPACE_EXPERIMENTAL` | Enable experimental features |
| `CARAPACE_HIDDEN` | Show hidden commands/flags |
| `CARAPACE_TOOLTIP` | Enable tooltip style |
| `CARAPACE_DESCRIPTION_LENGTH` | Max description length (default 80) |
| `NO_COLOR` / `CLICOLOR=0` | Disable colors |

## References

- [GNU Bash Manual: Bash Startup Files](https://www.gnu.org/software/bash/manual/html_node/Bash-Startup-Files.html)
- [GNU Bash Manual: The Set Builtin](https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html)
- [GNU Bash Manual: The Shopt Builtin](https://www.gnu.org/software/bash/manual/html_node/The-Shopt-Builtin.html)
- [GNU Bash Manual: Bash Variables](https://www.gnu.org/software/bash/manual/html_node/Bash-Variables.html)

## Related Skills

- **references/completion.md** — where to install completions and how they're loaded
- **references/execution.md** — how startup files affect the execution environment
