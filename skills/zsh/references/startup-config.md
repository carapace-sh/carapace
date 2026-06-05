# Zsh Startup Files and Configuration

In-depth reference for zsh's startup file order, the `fpath` mechanism, function autoloading, `compinit` initialization, and how to install custom completion functions.

## Startup File Decision Flow

### Interactive Login Shell (or `--login` flag)

1. `/etc/zshenv` — system-wide (always read, cannot be overridden)
2. `$ZDOTDIR/.zshenv` (or `$HOME/.zshenv` if `ZDOTDIR` unset) — always read
3. `/etc/zprofile` — login shell only
4. `$ZDOTDIR/.zprofile` — login shell only
5. `/etc/zshrc` — interactive shell only
6. `$ZDOTDIR/.zshrc` — interactive shell only
7. `/etc/zlogin` — login shell only
8. `$ZDOTDIR/.zlogin` — login shell only

**At logout (login shells only):**

1. `$ZDOTDIR/.zlogout`
2. `/etc/zlogout`

### Interactive Non-Login Shell

1. `/etc/zshenv`
2. `$ZDOTDIR/.zshenv`
3. `/etc/zshrc`
4. `$ZDOTDIR/.zshrc`

### Non-Interactive Shell (script execution)

1. `/etc/zshenv`
2. `$ZDOTDIR/.zshenv`

### Mode Combinations

| Mode | .zshenv | .zprofile | .zshrc | .zlogin |
|------|---------|-----------|--------|---------|
| Login + Interactive | ✓ | ✓ | ✓ | ✓ |
| Login + Non-interactive | ✓ | ✓ | ✗ | ✓ |
| Non-login + Interactive | ✓ | ✗ | ✓ | ✗ |
| Non-login + Non-interactive | ✓ | ✗ | ✗ | ✗ |

**Key insight**: `.zshenv` is executed for ALL shell invocations. Only set essential environment variables there — avoid expensive operations.

### The RCS and GLOBAL_RCS Options

```zsh
setopt no_rcs         # Skip ALL startup files
setopt no_global_rcs  # Skip /etc/ (system-wide) files only
```

### ZDOTDIR

`ZDOTDIR` controls where zsh looks for user startup files. If unset, defaults to `$HOME`:

```zsh
export ZDOTDIR=$HOME/.config/zsh
```

This allows moving all zsh config to `~/.config/zsh/` for a cleaner home directory.

## The fpath Array

`fpath` (or `FPATH`) is an array of directories searched for autoloadable functions. It is the primary mechanism for adding custom completion functions.

### Default fpath

Typically includes:

```
/usr/local/share/zsh/site-functions
/usr/share/zsh/site-functions
/usr/share/zsh/5.9/functions
```

### Adding Custom Directories

```zsh
# Prepend custom directory (takes priority)
fpath=(~/my-completions $fpath)

# Remove duplicates
typeset -U fpath
```

**Important**: `fpath` must be modified **before** calling `compinit`, because `compinit` scans `fpath` directories for completion functions.

### Completion Directory Structure

Zsh's built-in completion functions are organized in subdirectories:

| Directory | Contents |
|-----------|----------|
| `Base/` | Core functions and widgets |
| `Base/Core/` | Essential completion functions |
| `Base/Utility/` | Utility functions (`_describe`, `_arguments`, etc.) |
| `Base/Completion/` | Completer functions (`_complete`, `_approximate`, etc.) |
| `Zsh/` | Shell builtin completions |
| `Unix/` | External command completions |
| `Unix/Command/` | Specific command completions |
| `X/` | X11-related completions |
| `Linux/`, `BSD/`, `AIX/`, etc. | Platform-specific completions |

## Function Autoloading

### The autoload Builtin

```zsh
autoload fn           # Mark function for autoloading
autoload -Uz fn       # Recommended: suppress alias expansion, zsh-style
autoload +X fn        # Load function definition without calling it
```

| Option | Description |
|--------|-------------|
| `-U` | Suppress alias expansion when loading (recommended) |
| `-z` | Force zsh-style loading (default) |
| `-k` / `-KSH_AUTOLOAD` | ksh-style: file content executed directly |
| `+X` | Load definition without executing |
| `-X` | Define as autoloaded stub |

### How Autoloading Works

1. When an autoloaded function is first called, zsh searches `fpath` for a file matching the function name
2. The file content is loaded as the function definition
3. Subsequent calls use the loaded definition

### Loading All Functions in a Directory

```zsh
autoload -Uz ${fpath[1]}/*(:t)
```

The `(:t)` glob qualifier extracts just the filename (tail), stripping the directory path.

## compinit — Completion System Initialization

### Basic Setup

```zsh
autoload -Uz compinit && compinit
```

### What compinit Does

1. Scans all directories in `fpath` for files starting with `_`
2. Reads the first line of each file for `#compdef` directives
3. Builds the `_comps` associative array (command → function mapping)
4. Creates the `_main_complete` dispatcher
5. Sets up keybindings for completion widgets
6. Writes the dump file (`.zcompdump`) for faster future invocations

### compinit Options

| Option | Description |
|--------|-------------|
| `-D` | Disable dump file creation |
| `-d file` | Specify dump file location |
| `-C` | Skip security checks, use cached dump (fastest) |
| `-u` | Use insecure files without asking |
| `-i` | Silently ignore insecure files |
| `-n` | Skip function existence check |

### Security Model

`compinit` checks that completion files and their parent directories:
- Are owned by root or the current user
- Are not in world-writable or group-writable directories

Insecure files are reported and skipped unless `-u` or `-i` is given.

Audit and fix insecure files:

```zsh
compaudit | xargs chmod g-w
```

### The Dump File (zcompdump)

`compinit` caches its mapping in a dump file for faster subsequent invocations.

Default location: `$ZDOTDIR/.zcompdump` or `$HOME/.zcompdump`

Optimized initialization pattern:

```zsh
autoload -Uz compinit
# Only do a full check once per day
if [[ -n ~/.zcompdump(#qN.mh+24) ]]; then
  compinit
else
  compinit -C
fi
```

The `(#qN.mh+24)` glob qualifier matches the dump file only if it was modified more than 24 hours ago.

Force a rebuild:

```zsh
rm ~/.zcompdump* && compinit
```

### The #compdef Directive

Completion script files declare their target commands via a first-line directive:

```zsh
#compdef git                    # Complete for git command
#compdef git stg                # Complete for git and stg
#compdef -p git*                # Pattern-based (tried before exact match)
#compdef -P git*                # Pattern-based (tried after exact match)
#compdef -k _complete complete-word ^I  # Key binding
#compdef -K function name style keyseqs  # Multiple widgets
#autoload [options]             # Function for autoloading
```

`compinit` scans for these directives when building the `_comps` mapping. Files without a recognized directive are ignored by the completion system.

### The _compdir Parameter

If fewer than 20 `_` files are found in `fpath`, `compinit` tries to add directories using the `_compdir` parameter:

```zsh
_compdir=/usr/share/zsh/functions compinit
```

## compinstall — Interactive Configuration

`compinstall` is an interactive tool that helps configure the completion system. It inserts code into `.zshrc`:

```zsh
autoload -Uz compinstall && compinstall
```

It handles:
- Setting `completer` style
- Configuring `matcher-list`
- Setting `menu` style
- Configuring colors (`list-colors`)
- Setting `group-name` and `format` styles
- Enabling/disabling caching

`compinstall` can be rerun without losing existing zstyle commands.

## bashcompinit — Bash Compatibility

`bashcompinit` provides bash-style `complete` and `compgen` functions for compatibility:

```zsh
autoload -Uz bashcompinit && bashcompinit
```

After loading, bash-style completion registration works:

```zsh
complete -F _myfunc mycommand
```

This is useful for porting bash completions to zsh without rewriting them.

## Installing Custom Completion Functions

### Method 1: Add to fpath

```zsh
# In .zshrc, before compinit:
fpath=(~/.zsh/completions $fpath)
autoload -Uz compinit && compinit
```

Place completion files in `~/.zsh/completions/` with names starting with `_`:

```
~/.zsh/completions/
  _myapp
  _otherapp
```

### Method 2: Direct compdef

```zsh
# In .zshrc, after compinit:
compdef _myapp_completion myapp

_myapp_completion() {
    _describe 'subcommand' subcmds
}
```

### Method 3: Source the File

```zsh
# In .zshrc, after compinit:
source ~/.zsh/completions/_myapp
```

This immediately defines the function and registers it.

### Method 4: Plugin Manager

With a plugin manager (e.g., zinit, znap, antibody):

```zsh
# zinit example
zinit ice wait lucid
zinit snippet OMZP::git
```

Plugin managers handle `fpath` manipulation and `compinit` calls automatically.

### Reloading After Changes

After editing a completion function:

```zsh
unfunction _mycmd && autoload -U _mycmd
```

Or to rebuild the entire completion cache:

```zsh
rm ~/.zcompdump* && exec zsh
```

## Common Configuration Patterns

### Minimal Fast Setup

```zsh
autoload -Uz compinit && compinit -C
```

### Full Setup with Caching and Security

```zsh
autoload -Uz compinit

# Cache for 24 hours
if [[ -n ~/.zcompdump(#qN.mh+24) ]]; then
  compinit -u
else
  compinit -C -u
fi
```

### XDG-Compatible Setup

```zsh
export ZDOTDIR=$HOME/.config/zsh
# In $ZDOTDIR/.zshrc:
fpath=($ZDOTDIR/completions $fpath)
autoload -Uz compinit && compinit -d $XDG_CACHE_HOME/zsh/zcompdump
```

## References

### Documentation

- [Zsh Manual: Startup Files](https://zsh.sourceforge.io/Doc/Release/Files.html) — official startup file reference
- [Zsh Manual: Functions](https://zsh.sourceforge.io/Doc/Release/Functions.html) — function autoloading
- [Zsh Manual: Completion System](https://zsh.sourceforge.io/Doc/Release/Completion-System.html) — compinit and compdef
- [zsh(1) man page](https://linux.die.net/man/1/zsh) — main zsh reference
- [zshcompsys(1) man page](https://linux.die.net/man/1/zshcompsys) — completion system reference

### Tutorials

- [Zsh Completions HOWTO](https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org) — writing and installing completions
- [Solving compinit's Insecure Files Warning](https://dev.to/manojspace/solving-zsh-compinits-insecure-files-warning-34pg) — security model explained
- [Zsh Autoloading Functions](https://blog.augustfeng.app/articles/zsh-autoloading-functions/) — autoload internals
