# Elvish Startup and Configuration

In-depth reference for elvish's startup process, configuration directory, `rc.elv`, module system, package manager, and runtime configuration.

## Configuration Directory

Elvish uses the XDG Base Directory specification:

| Path | Purpose |
|------|---------|
| `~/.config/elvish/` | Configuration directory (or `$XDG_CONFIG_HOME/elvish/`) |
| `~/.config/elvish/rc.elv` | Startup script |
| `~/.config/elvish/lib/` | User-defined module search directory |
| `~/.local/share/elvish/` | Data directory (or `$XDG_DATA_HOME/elvish/`) |
| `~/.local/share/elvish/db.elv` | Persistent store database (history, location, etc.) |

### Legacy Directory (Removed in 0.21.0)

The `~/.elvish/` directory was deprecated in 0.16.0 and removed in 0.21.0. Elvish no longer reads from this location.

## The `rc.elv` Startup Script

`~/.config/elvish/rc.elv` is executed when elvish starts in interactive mode. It is the primary place for:

- Setting up keybindings
- Registering arg-completers
- Customizing prompts
- Loading modules
- Defining helper functions
- Setting environment variables

### Example `rc.elv`

```elvish
# Load modules
use str
use re
use path

# Customize prompt
set edit:prompt = { tilde-abbr $pwd; styled '> ' cyan }

# Right prompt with git branch
set edit:rprompt = {
    var branch = (git rev-parse --abbrev-ref HEAD 2>/dev/null)
    if (not-eq $branch '') {
        styled ' '$branch magenta
    }
}

# Keybindings
set edit:insert:binding[Ctrl-F] = { edit:complete-filename }
set edit:insert:binding[Alt-l] = { edit:location:start }

# Completion for custom commands
set edit:completion:arg-completer[myapp] = {|@args|
    myapp _carapace elvish (all $args) | from-json | each {|completion|
        put $completion[Candidates] | all (one) | peach {|c|
            edit:complex-candidate $c[Value] &display=(styled $c[Display] $c[Style]) &code-suffix=$c[CodeSuffix]
        }
    }
}

# Abbreviations
set edit:abbr[gs] = 'git status'
set edit:abbr[gc] = 'git commit'

# Case-insensitive matching for arguments
set edit:completion:matcher[argument] = {|seed| edit:match-prefix $seed &ignore-case=$true }

# After-command hook
set edit:after-command = [ $@edit:after-command {|m|
    # Track command duration
} ]
```

### Startup Order

1. Elvish binary starts
2. XDG directories are resolved
3. Persistent store (`db.elv`) is opened
4. `rc.elv` is executed (if it exists and elvish is in interactive mode)
5. Editor is initialized (prompts, keybindings, completions)
6. REPL loop begins

## Module System

### Module Search Directories

Elvish searches for modules in these directories (in order):

1. `~/.config/elvish/lib/` (or `$XDG_CONFIG_HOME/elvish/lib/`)
2. Directories listed in `$E:ELVISH_LIB_PATH` (colon-separated, since 0.20.0)
3. Bundled standard library modules (pre-defined modules)

### Importing Modules

```elvish
use str                    # import str module
use str s                  # import with alias 's'
use ./my-module            # relative import
use ../shared/utils       # relative import from parent
```

### Writing Custom Modules

A module is a `.elv` file in a module search directory:

```
~/.config/elvish/lib/
  my-module.elv            # use my-module
  utils/
    strings.elv            # use utils/strings
```

**`my-module.elv`:**

```elvish
# Module: my-module
# This file is loaded when 'use my-module' is called

fn greet {|name|
    echo "Hello, "$name"!"
}

fn list-files {|dir|
    put $dir/*[nomatch-ok]
}

# Exported via namespace — all top-level functions and variables
# are accessible as my-module:greet, my-module:list-files
```

### Module Caching

Modules are loaded once and cached. Subsequent `use` calls return the cached namespace. This means:

- Module code only executes once
- Changes to module files require restarting elvish
- Circular imports are detected and reported as errors

### Re-importing

```elvish
use-mod my-module  # outputs the namespace without importing into current scope
```

## Elvish Package Manager (epm)

`epm` is elvish's built-in package manager for installing third-party modules.

### Installing Packages

```elvish
use epm
epm:install github.com/user/package
```

### Package Directory

Installed packages are stored in:

```
~/.config/elvish/lib/
  github.com/
    user/
      package/
        *.elv files
```

### Package Registry

The default registry is at [github.com/elves/epm-reg](https://github.com/elves/epm-reg). Packages can also be installed directly from GitHub repositories.

### epm Commands

| Command | Description |
|---------|-------------|
| `epm:install $pkg` | Install a package |
| `epm:uninstall $pkg` | Uninstall a package |
| `epm:update $pkg` | Update a package |
| `epm:update` | Update all installed packages |
| `epm:list` | List installed packages |
| `epm:installed $pkg` | Check if package is installed |

## Standard Library Modules

### `str` — String Manipulation

```elvish
use str
str:contains $s $substr     # substring test
str:has-prefix $s $prefix   # prefix test
str:has-suffix $s $suffix   # suffix test
str:join $sep $list          # join list with separator
str:split $sep $s            # split string
str:replace $old $new $s    # replace
str:to-lower $s              # lowercase
str:to-upper $s              # uppercase
str:trim $s $cutset          # trim characters
str:trim-space $s            # trim whitespace
str:fields $s                # split by whitespace
str:repeat $n $s             # repeat string
```

### `re` — Regular Expressions

```elvish
use re
re:match $pattern $s         # boolean match
re:find $pattern $s          # find all matches
re:replace $pattern $repl $s # replace
re:split $pattern $s          # split
re:awk $pattern {|fields...| body }  # awk-like field extraction
```

### `math` — Mathematical Functions

```elvish
use math
math:abs $x          # absolute value
math:ceil $x         # ceiling
math:floor $x        # floor
math:round $x        # round
math:sqrt $x         # square root
math:pow $x $y       # power
math:log $x          # natural log
math:max $a $b       # maximum
math:min $a $b       # minimum
# Constants: $math:e, $math:pi
```

### `path` — Path Manipulation

```elvish
use path
path:join $base $rel         # join path components
path:separator               # OS path separator (/ or \)
path:list-separator          # OS list separator (: or ;)
path:dev-tty                 # path to TTY device
path:dev-null                # path to null device
```

### `platform` — Platform Information

```elvish
use platform
$platform:os          # operating system
$platform:arch        # architecture
$platform:is-unix     # boolean
$platform:is-windows  # boolean
platform:hostname     # hostname
```

### `os` — OS Operations

```elvish
use os
os:mkdir-all $path    # create directory tree
os:symlink $old $new  # create symbolic link
os:rename $old $new   # rename file
```

### `file` — File Operations

```elvish
use file
file:open $path           # open file for reading
file:open-output $path    # open file for writing
file:seek $file $offset   # seek in file
file:tell $file           # get position in file
```

### `unix` — Unix-Specific

```elvish
use unix
$unix:rlimits       # resource limits map
$unix:umask         # current umask
```

### `md` — Markdown Rendering

```elvish
use md
md:show $markdown   # render markdown in terminal
```

### `doc` — Documentation Access

```elvish
use doc
doc:find $keyword    # find documentation
doc:show $topic      # show documentation
doc:source           # documentation source
```

### `runtime` — Runtime Paths

```elvish
use runtime
# Exposes paths important for the elvish runtime
```

## Environment Variables

### Accessing Environment Variables

```elvish
echo $E:PATH         # read PATH
set E:EDITOR = vim   # set EDITOR
```

The `E:` namespace provides read/write access to environment variables. `$E:PATH` and `$paths` are kept in sync.

### The `$paths` Variable

```elvish
$paths               # list of PATH directories (kept in sync with $E:PATH)
set paths = [~/bin $@paths]  # prepend to PATH
```

### Elvish-Specific Environment Variables

| Variable | Purpose |
|----------|---------|
| `$E:ELVISH_LIB_PATH` | Additional module search directories (colon-separated) |
| `$E:NO_COLOR` | Disable colors if set and non-empty (since 0.20.0) |
| `$E:LSCOLORS` | File coloring configuration (macOS) |
| `$E:LS_COLORS` | File coloring configuration (Linux) |
| `$E:COLORTERM` | Used to detect TrueColor support |

## Persistent Store

The persistent store (`db.elv`) is a BoltDB file that stores:

- **Command history** — all commands entered in interactive mode
- **Directory history** — directories visited (used by location mode)
- **Shared variables** — (removed in 0.19.1)

### Store Location

```
~/.local/share/elvish/db.elv
```

### History Configuration

```elvish
# Control what gets added to history
set edit:add-cmd-filters = [ $@edit:add-cmd-filters {|cmd|
    # Don't add commands starting with space
    not-eq $cmd[0] ' '
} ]
```

## Integration with External Completion Frameworks

### Carapace Integration

Add to `rc.elv`:

```elvish
# Register carapace for a specific command
set edit:completion:arg-completer[myapp] = {|@arg|
    myapp _carapace elvish (all $arg) | from-json | each {|completion|
        put $completion[Messages] | all (one) | each {|m|
            edit:notify (styled "error: " red)$m
        }
        if (not-eq $completion[Usage] "") {
            edit:notify (styled "usage: " $completion[DescriptionStyle])$completion[Usage]
        }
        put $completion[Candidates] | all (one) | peach {|c|
            if (eq $c[Description] "") {
                edit:complex-candidate $c[Value] &display=(styled $c[Display] $c[Style]) &code-suffix=$c[CodeSuffix]
            } else {
                edit:complex-candidate $c[Value] &display=(styled $c[Display] $c[Style])(styled " " $completion[DescriptionStyle]" bg-default")(styled "("$c[Description]")" $completion[DescriptionStyle]) &code-suffix=$c[CodeSuffix]
            }
        }
    }
}
```

### Carapace-Bin Integration

For carapace-bin (which provides completions for 500+ commands), the snippet can be evaluated:

```elvish
# Evaluate the carapace snippet for all commands
eval (carapace _carapace elvish | slurp)
```

This registers arg-completers for all commands that carapace-bin supports.

## References

- [Elvish Get/Install](https://elv.sh/get/) — installation instructions
- [Elvish Learn: Tour](https://elv.sh/learn/tour.html) — startup script section
- [Elvish Language: Modules](https://elv.sh/ref/language.html#namespaces-and-modules) — module specification
- [Elvish Builtin Functions](https://elv.sh/ref/builtin.html) — builtin commands
- [epm Repository](https://github.com/elves/epm) — package manager

## Related Skills

- For language fundamentals (value types, expressions, pipelines), see [references/language.md](language.md).
- For completion system details (arg-completer, complex-candidate), see [references/completion.md](completion.md).
- For carapace-specific elvish integration, see the **carapace-dev** skill → `references/shell-elvish.md`.
