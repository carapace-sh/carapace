# Tcsh Startup Files and Configuration

In-depth reference for which files tcsh reads at startup, how shell variables and environment variables work, and where to install completions.

## Startup File Decision Flow

### Login Shell

A login shell is one where the first argument (argv[0]) starts with `-`, or the shell is invoked with `-l` as the only flag.

**Execution order:**

1. `/etc/csh.cshrc` — system-wide, read by every shell
2. `/etc/csh.login` — system-wide, login shells only
3. `~/.tcshrc` — user config, read by every shell (falls back to `~/.cshrc` if not found)
4. `~/.history` — login shells only, if `savehist` is set
5. `~/.login` — login shells only
6. `~/.cshdirs` — login shells only, if `savedirs` is set

**On logout (login shells only):**

1. `/etc/csh.logout` — system-wide
2. `~/.logout` — user config

### Non-Login Interactive Shell

1. `/etc/csh.cshrc` — system-wide
2. `~/.tcshrc` (or `~/.cshrc` fallback) — user config

### Non-Interactive Shell (script execution)

1. `/etc/csh.cshrc` — system-wide
2. `~/.tcshrc` (or `~/.cshrc` fallback) — user config

Unlike bash's `BASH_ENV`, tcsh always reads `~/.tcshrc` for non-interactive shells. There is no way to skip this except with the `-f` flag.

### The `-f` Flag

The `-f` flag causes tcsh to start faster by not reading any startup files:

```tcsh
tcsh -f -c "command"
```

This is useful for scripts that should not be affected by user configuration.

### The `-l` Flag

The `-l` flag makes tcsh behave as a login shell:

```tcsh
tcsh -l
```

### The `-m` Flag

The `-m` flag causes tcsh to load `~/.tcshrc` even if it does not belong to the effective user:

```tcsh
tcsh -m
```

### System-Specific File Locations

| System | cshrc | login | logout |
|--------|-------|-------|--------|
| Default | `/etc/csh.cshrc` | `/etc/csh.login` | `/etc/csh.logout` |
| Solaris 2.x | `/etc/.cshrc` | `/etc/.login` | `/etc/.logout` |
| ConvexOS, Stellix | `/etc/cshrc` | `/etc/login` | `/etc/logout` |
| NeXT | `/etc/cshrc.std` | `/etc/login.std` | `/etc/logout.std` |
| A/UX, AMIX, Cray, IRIX | `/etc/cshrc` | (login/logout in tcsh) | — |

## The `~/.tcshrc` vs `~/.cshrc` Fallback

If `~/.tcshrc` does not exist, tcsh falls back to `~/.cshrc`. This provides backward compatibility with csh.

### Compatibility Strategy

For configurations that must work with both csh and tcsh:

1. **Use only `~/.cshrc`** with tcsh-specific checks:

```tcsh
if ($?tcsh) then
  # tcsh-specific settings
  set complete = enhance
  set edit
  bindkey -e
endif
```

2. **Use both files** where `~/.tcshrc` sources `~/.cshrc`:

```tcsh
# In ~/.tcshrc:
source ~/.cshrc
# tcsh-specific additions follow
```

## Where to Install Completions

Completions defined with the `complete` builtin should be placed in `~/.tcshrc` (or a file sourced from it):

```tcsh
# In ~/.tcshrc:
complete cd 'p/1/d/'
complete man 'p/*/c/'
complete find 'n/-name/f/' 'n/-user/u/' 'c/-/(name newer user group type)/' 'p/*/d/'
```

### System-Wide Completions

System-wide completions can be placed in `/etc/csh.cshrc`, but this affects all users. A better approach is to create a separate file:

```tcsh
# In /etc/csh.cshrc:
if (-r /etc/tcsh.completions) then
  source /etc/tcsh.completions
endif
```

### The `complete.tcsh` Example File

tcsh ships with an example completions file (~1277 lines) typically installed at:

- `/usr/share/doc/tcsh/examples/complete.gz` (Debian)
- `/usr/share/tcsh/complete.tcsh` (some systems)

To use it:

```tcsh
# In ~/.tcshrc:
if (-r /usr/share/tcsh/complete.tcsh) then
  source /usr/share/tcsh/complete.tcsh
endif
```

## Shell Variables

### The `set` Builtin

```tcsh
set                    # Print all variables
set name               # Set to null string
set name = word        # Set to single word
set name = (wordlist)  # Set to list of words
set name[index] = word # Set specific element
set -r name            # Make read-only
set -r name = word     # Set and make read-only
```

### The `unset` Builtin

```tcsh
unset pattern    # Remove variables matching pattern
```

Removes all variables whose names match the glob pattern (except read-only variables).

### Key Shell Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `argv` | Shell's argument list | — |
| `cwd` | Current working directory | — |
| `home` | Home directory | From environment |
| `path` | Command search path | From environment |
| `prompt` | Primary prompt string | `%# ` |
| `prompt2` | Continuation prompt | `%R? ` |
| `prompt3` | Correction prompt | `CORRECT>%R (y\|n\|e\|a)? ` |
| `status` | Exit status of last command | — |
| `history` | Number of commands to save | — |
| `savehist` | Save history on logout | — |
| `histfile` | History file location | `~/.history` |
| `histdup` | Handle duplicates: `all`, `prev`, `erase` | — |
| `histlit` | Use literal (unexpanded) form | — |
| `edit` | Enable command-line editor | Set (interactive) |
| `vimode` | Use vi-style editing | Unset |
| `wordchars` | Non-alphanumeric word characters | `*?_-.[]~=` (emacs) |
| `ignoreeof` | Prevent exit on `^D` | — |
| `noclobber` | Prevent accidental file overwrite | — |
| `noglob` | Disable filename substitution | — |
| `nonomatch` | Leave unmatched patterns unchanged | — |
| `verbose` | Echo input after history substitution | — |
| `echo` | Echo commands before execution | — |
| `addsuffix` | Add `/` or space after completion | Set |
| `autolist` | List choices on failed completion | — |
| `autoexpand` | Expand history before completion | — |
| `autocorrect` | Spell-correct before completion | — |
| `complete` | `enhance` for smart completion | — |
| `correct` | Spelling correction mode | — |
| `fignore` | File suffixes to ignore | — |
| `listmax` | Max items before prompting | — |
| `listmaxrows` | Max rows before prompting | — |
| `matchbeep` | When to beep on completion | `ambiguous` |
| `nobeep` | Disable all beeping | — |
| `visiblebell` | Flash instead of beep | — |
| `color` | Enable color in `ls-F` | — |
| `nostat` | Directories to skip stat(2) | — |
| `recexact` | Complete on shortest unique match | — |
| `recognize_only_executables` | List only executables | — |
| `autologout` | Auto-logout after N minutes | 60 (login) |
| `mail` | Check for new mail | — |
| `notify` | Immediate job status notification | — |
| `time` | Auto-time threshold | — |
| `watch` | Report user login/logout | — |
| `tcsh` | Version string (e.g., `6.24.16`) | — |
| `version` | Compile-time options | — |

### Variable Substitution Syntax

| Syntax | Description |
|--------|-------------|
| `$name` | Value of variable |
| `${name}` | Value of variable (unambiguous) |
| `$?name` | `1` if set, `0` if not |
| `$#name` | Number of words in variable |
| `$name[selector]` | Select words by index or range |
| `$$` | Process ID of parent shell |
| `$!` | Process ID of last background process |
| `$<` | Read a line from stdin |
| `$0` | Name of current script/shell |

## Environment Variables

### The `setenv` Builtin

```tcsh
setenv              # Print all environment variables
setenv NAME value   # Set environment variable
setenv NAME         # Set to null string
```

### The `unsetenv` Builtin

```tcsh
unsetenv pattern    # Remove environment variables matching pattern
```

### Key Environment Variables

| Variable | Description |
|----------|-------------|
| `HOME` | User's home directory |
| `PATH` | Command search path (synced with `path` variable) |
| `TERM` | Terminal type |
| `USER` | Username |
| `SHELL` | Path to current shell |
| `PWD` | Current working directory |
| `COMMAND_LINE` | Set during completion backtick execution |
| `HOSTTYPE` | System type (set by tcsh) |
| `VENDOR` | System vendor (set by tcsh) |
| `OSTYPE` | OS type (set by tcsh) |
| `MACHTYPE` | Machine type (set by tcsh) |
| `GROUP` | User's group (set by tcsh) |
| `HOST` | Hostname (set by tcsh) |

### The `path` Variable vs `PATH` Environment Variable

tcsh automatically synchronizes the `path` shell variable and the `PATH` environment variable:

- Setting `path` updates `PATH`
- Setting `PATH` with `setenv` updates `path`
- The `path` variable is a list: `set path = (/usr/local/bin /usr/bin /bin)`

### The `COMMAND_LINE` Environment Variable

This variable is set only during completion backtick execution. It contains the full command line being completed. See [references/completion.md](references/completion.md) for details.

## The `@` Builtin (Arithmetic)

```tcsh
@                    # Print all variables
@ name = expr        # Assign arithmetic expression
@ name[index] = expr # Assign to array element
@ name++             # Increment
@ name--             # Decrement
```

Operators: `+`, `-`, `*`, `/`, `%`, `==`, `!=`, `>`, `<`, `>=`, `<=`, `!`, `~`, `&`, `|`, `^`, `<<`, `>>`, `&&`, `||`, `?:` (ternary).

Note: `<`, `>`, `&`, and `|` must be quoted to prevent shell interpretation:

```tcsh
@ x = 5 '<' 10
@ flags = $flags '|' 0x02
```

## The `source` Builtin

```tcsh
source [-h] name [args ...]
```

Reads and executes commands from `name` in the current shell context. Arguments are available in `argv`. The `-h` flag places commands in the history list.

Common usage for loading completions:

```tcsh
source ~/.tcshrc      # Reload configuration
source ~/completions.tcsh  # Load completion definitions
```

## The `rehash` Builtin

```tcsh
rehash
```

Recomputes the internal hash table of command locations. Needed if a new command is added to a directory in `path` while `autorehash` is not set.

## The `hashstat` Builtin

```tcsh
hashstat
```

Prints statistics on the effectiveness of the internal hash table for command lookup.

## The `unhash` Builtin

```tcsh
unhash
```

Disables the internal hash table, forcing a PATH search for every command.

## History Configuration

### Key Variables

| Variable | Description |
|----------|-------------|
| `history` | Number of commands to keep in memory |
| `savehist` | Number of commands to save to `histfile` on logout |
| `histfile` | File to save/restore history (default: `~/.history`) |
| `histdup` | Handle duplicates: `all` (keep), `prev` (keep if not previous), `erase` (remove old) |
| `histlit` | Store literal (unexpanded) form of commands |
| `histchars` | Change history characters (default: `!^#`) |

### Example Configuration

```tcsh
set history = 1000
set savehist = (1000 merge)
set histfile = ~/.tcsh_history
set histdup = erase
```

The `merge` keyword in `savehist` merges the current session's history with the existing history file rather than overwriting it.

## Directory Stack Configuration

| Variable | Description |
|-----------|------------|
| `savedirs` | Save directory stack on logout |
| `dirsfile` | File to save/restore directory stack |
| `dirstack` | Current directory stack (array) |
| `pushdtohome` | `pushd` with no args pushes home |
| `pushdsilent` | Don't print stack on pushd/popd |
| `dunique` | Don't push duplicates onto stack |

## Terminal Configuration

| Variable | Description |
|-----------|------------|
| `term` | Terminal type (usually set in `~/.login`) |
| `tty` | Tty name (set by tcsh) |
| `termcap` | Terminal capabilities (set by tcsh) |

## Login Shell Detection

The `$?loginsh` variable is `1` if the shell is a login shell:

```tcsh
if ($?loginsh) then
  echo "This is a login shell"
endif
```

## References

- tcsh(1) man page — startup, shell variables, environment variables
- `sh.init.c` in the tcsh source — initialization code
- `complete.tcsh` in the tcsh source — example completions

## Related Skills

- [references/completion.md](references/completion.md) — completion-related variables and where to install completions
- [references/editor.md](references/editor.md) — editor-related variables (edit, vimode, wordchars)
- [references/execution.md](references/execution.md) — how startup files are executed
