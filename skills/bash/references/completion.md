# Bash Programmable Completion

In-depth reference for bash's programmable completion system — the builtins, variables, registration, and the bash-completion project's helper framework.

## The Completion Flow

When the user presses TAB (or another completion key), bash:

1. Identifies the command word and searches for a matching compspec (completion specification)
2. Sets `COMP_*` variables in the completion function's environment
3. Invokes the registered function or command
4. Reads `COMPREPLY` for completion candidates
5. Passes candidates to Readline for display

### Compspec Search Order

Bash searches for a compspec in this order:

1. **Full pathname** — if the command contains a `/`, the full path is checked first
2. **Portion after final slash** — e.g., `/usr/bin/git` → check for `git`
3. **Default compspec** (`-D`) — for commands with no compspec
4. **Alias expansion result** — if the command is an alias, check the expansion

### Match Generation Order

Within a compspec, matches are generated in this order:

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

When a completion function is invoked, bash sets these variables:

| Variable | Type | Description |
|----------|------|-------------|
| `COMP_WORDS` | Array | All words on the command line (split by whitespace) |
| `COMP_CWORD` | Integer | Index in `COMP_WORDS` of the word containing the cursor |
| `COMP_LINE` | String | The entire command line as a single string |
| `COMP_POINT` | Integer | Cursor position within `COMP_LINE` (0-based index) |
| `COMP_TYPE` | Integer | Type of completion attempt (see table below) |
| `COMP_KEY` | Integer | Key pressed to invoke completion (decimal ASCII) |
| `COMP_WORDBREAKS` | String | Characters that Readline treats as word breaks |

### COMP_TYPE Values

| Value | Character | Meaning |
|-------|-----------|---------|
| `9` | TAB | Normal completion (first TAB press) |
| `33` | `!` | Listing alternatives on partial word |
| `37` | `%` | Menu completion |
| `63` | `?` | Listing after successive TABs |
| `64` | `@` | Listing if word not unmodified |

### COMP_WORDBREAKS

Default value: `" \t\n\"'@><=;|&(:`". These characters split the "word being completed" — bash only passes the **last segment** (after the most recent word-break character) as `$2` to the completion function.

This has major implications:

- **Colon (`:`)** — breaks `origin/main`, `user@host:port`, `File::Basename`
- **Equals (`=`)** — breaks `--option=value` into `--option` and `value`
- **At sign (`@`)** — breaks `user@host`

Completion functions must reassemble the full word by re-parsing `COMP_LINE` or using helpers like `_comp__reassemble_words`.

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
| `helptopic` | | Help topics (`help` builtin) |
| `hostname` | | Hostnames from `HOSTFILE` |
| `job` | `-j` | Job names (if job control active) |
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
| `bashdefault` | Fall back to Bash default completions if compspec generates no matches |
| `default` | Use Readline's default filename completion if no matches |
| `dirnames` | Attempt directory name completion if no matches |
| `filenames` | Tell Readline completions are filenames (add slash to dirs, quote special chars, suppress trailing spaces) |
| `fullquote` | Quote all completed words (not just filenames) |
| `noquote` | Don't quote completed words |
| `nosort` | Don't sort completions alphabetically (Bash 4.4+) |
| `nospace` | Don't append space after completion |
| `plusdirs` | After generating matches, also add directory completions |

### Special Registration Options

| Option | Description |
|--------|-------------|
| `-D` | Apply to "default" command completion (commands with no compspec) |
| `-E` | Apply to "empty" command completion (blank line) |
| `-I` | Apply to initial word (command name completion) |
| `-p` | Print existing compspecs (for reuse) |
| `-r` | Remove compspec for specified names |

Priority: `-D` > `-E` > `-I`. If any of these are supplied, name arguments are ignored.

### Other Options

| Option | Description |
|--------|-------------|
| `-C command` | Execute command in subshell; stdout becomes completions |
| `-F function` | Execute shell function; read `COMPREPLY` for results |
| `-G globpat` | Expand filename pattern for completions |
| `-W wordlist` | Split wordlist by IFS, expand each word, match against current word |
| `-X filterpat` | Remove completions matching pattern; `!` negates |
| `-P prefix` | Prepend prefix to each completion |
| `-S suffix` | Append suffix to each completion |

### `-F function` Details

When the function is called:

- `$1` — name of the command whose arguments are being completed
- `$2` — word being completed
- `$3` — word preceding the word being completed

The function populates `COMPREPLY`. When it finishes, bash reads `COMPREPLY` for candidates.

**Exit status 124**: If the function returns 124, bash restarts the completion from the beginning. This enables dynamic loading — the function can source a completion file and return 124 to re-trigger with the newly registered compspec.

## The `compgen` Builtin

```bash
compgen [-V varname] [option] [word]
```

Generates completion matches and writes them to stdout (or stores in `varname` with `-V`). Accepts the same options as `complete` except `-p`, `-r`, `-D`, `-E`, `-I`.

### Common Usage

```bash
compgen -W "start stop status" -- "$cur"     # Word list
compgen -f -- "$cur"                          # Filenames
compgen -d -- "$cur"                          # Directories
compgen -A command -- "$cur"                  # Commands
compgen -A variable -- "$cur"                # Variables
compgen -A function -- "$cur"                # Functions
compgen -A user -- "$cur"                    # Users
compgen -A hostname -- "$cur"                # Hostnames
compgen -W "prod dev" -P "--env=" -- "$cur"  # With prefix
compgen -f -X '!*.txt' -- "$cur"            # Only .txt files
```

The `--` separates options from the word argument. Always use `--` to prevent the word from being interpreted as an option.

## The `compopt` Builtin

```bash
compopt [-o option] [+o option] [-DEI] [name]
```

Modifies completion options dynamically — either for a named command or for the currently-executing completion (if no name given).

- `-o option` — enable option
- `+o option` — disable option
- `-D` / `-E` / `-I` — apply to default/empty/initial completion

This is the primary way to control nospace, filenames, nosort, etc. from within a completion function:

```bash
_mycomp() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    COMPREPLY=($(compgen -f -- "$cur"))
    compopt -o filenames   # Tell Readline these are filenames
    compopt -o nospace      # Don't append trailing space
}
```

## Writing Completion Functions

### Basic Pattern

```bash
_mytool_completion() {
    local cur prev words cword
    _init_completion -n = || return

    # Handle option arguments
    case "$prev" in
        --mode|-m)
            COMPREPLY=($(compgen -W "verbose quiet debug" -- "$cur"))
            return
            ;;
        --output|-o)
            COMPREPLY=($(compgen -d -- "$cur"))
            compopt -o filenames
            return
            ;;
    esac

    # Handle options
    case "$cur" in
        -*)
            COMPREPLY=($(compgen -W "--help --version --mode --output" -- "$cur"))
            ;;
    esac
}
complete -F _mytool_completion mytool
```

### Without `_init_completion`

```bash
_mycomp() {
    local cur="${COMP_WORDS[COMP_CWORD]}"
    local prev="${COMP_WORDS[COMP_CWORD-1]}"
    # ... logic ...
    COMPREPLY=($(compgen -W "opt1 opt2" -- "$cur"))
}
complete -F _mycomp mycmd
```

### Handling Subcommands

```bash
_mytool() {
    local cur prev words cword
    _init_completion || return

    local subcmd="${words[1]}"

    # If no subcommand yet, complete subcommand names
    if [[ $cword -eq 1 ]]; then
        COMPREPLY=($(compgen -W "build deploy status" -- "$cur"))
        return
    fi

    # Dispatch to subcommand handlers
    case "$subcmd" in
        build)  _mytool_build ;;
        deploy) _mytool_deploy ;;
        status) _mytool_status ;;
    esac
}
```

### Handling `--` Delimiter

After `--`, all arguments are positional (not options):

```bash
_mytool() {
    local cur prev words cword
    _init_completion || return

    # Check if -- has been seen
    local i=1
    for ((i=1; i<cword; i++)); do
        if [[ "${words[i]}" == "--" ]]; then
            # After --, only positional args
            COMPREPLY=($(compgen -f -- "$cur"))
            compopt -o filenames
            return
        fi
    done

    # Normal option/argument handling
    # ...
}
```

### Dynamic Option Parsing

Parse `--help` output to generate option completions:

```bash
# Using bash-completion helpers
_parse_help "$1" --help       # Parse --help output for options
_parse_usage "$1" --help      # Parse usage strings

# Manual approach
COMPREPLY=($(compgen -W "$($1 --help 2>&1 | \
    sed -ne 's/.*\(--[a-z][a-z-]*\).*/\1/p' | sort -u)" -- "$cur"))
```

## The bash-completion Project

The [bash-completion](https://github.com/scop/bash-completion) project provides a framework of helper functions and hundreds of pre-built completions.

### Directory Structure

```
/usr/share/bash-completion/
├── bash_completion              # Main script with helper functions
└── completions/                 # Individual command completions
    ├── git
    ├── ssh
    ├── tar
    └── ...

~/.local/share/bash-completion/
└── completions/                 # User-specific completions

/etc/bash_completion.d/         # Legacy compat directory
```

### Search Order for Completion Files

1. `BASH_COMPLETION_USER_DIR/completions/` (user)
2. Main `bash_completion` directory `completions/` (system)
3. Command's installation prefix `share/bash-completion/completions/`
4. `XDG_DATA_DIRS` `bash-completion/completions/`

### Dynamic Loading (`_completion_loader`)

Completions are loaded on demand, not at shell startup:

```bash
_completion_loader() {
    . "/usr/share/bash-completion/completions/$1" >/dev/null 2>&1 && return 124
}
complete -D -F _completion_loader -o bashdefault -o default
```

The `-D` flag registers this as the default handler. When a command with no compspec is completed:

1. `_completion_loader` sources the matching file
2. Returns 124 → bash restarts completion with the newly registered compspec
3. The actual completion function runs

### Naming Conventions

| Pattern | Purpose |
|---------|---------|
| `_comp_cmd_COMMAND` | Main completion function for a command |
| `_comp_cmd_COMMAND__HELPER` | Command-specific helper |
| `_comp_xfunc_COMMAND_ACTION` | Cross-functional helper (usable by other commands) |
| `_comp_compgen_*` | Completion generator functions |
| `_comp_*` | Core helper functions |

### Core Helper Functions

| Function | Purpose |
|----------|---------|
| `_comp_initialize` | Initialize `cur`, `prev`, `words`, `cword` (replaces `_init_completion`) |
| `_init_completion` | Legacy wrapper around `_comp_initialize` |
| `_comp_compgen` | Primary completion generator (calls compgen or generator functions) |
| `_comp_compgen_split` | Split text and generate completions (safer than `compgen -W "$(cmd)"`) |
| `_comp_compgen_set` | Directly set completion array with words |
| `_comp_split` | Split string into array (safe IFS handling) |
| `_comp_quote` | Shell-quote argument (result in `REPLY`) |
| `_comp_dequote` | Safely expand quoted word (result in `REPLY`) |
| `_comp__reassemble_words` | Reassemble words excluding chars from COMP_WORDBREAKS |
| `_comp_expand_glob` | Expand glob pattern in controlled environment |
| `_comp_looks_like_path` | Check if argument looks like a path (`/`, `.`, `~`) |
| `_comp_have_command` | Check if command exists in PATH |
| `_comp_userland` | Detect OS userland (GNU, BSD) |
| `_filedir` | File completion (handles spaces in filenames) |
| `_parse_help` | Parse `--help` output for options |
| `_parse_usage` | Parse usage strings for options |
| `_comp_count_args` | Count arguments (respects `--` delimiter) |
| `_comp_upvars` | Assign variables one scope above caller |

### `_comp_compgen` Options

| Option | Description |
|--------|-------------|
| `-a` | Append to COMPREPLY instead of replacing |
| `-v arr` | Store results in array `arr` instead of COMPREPLY |
| `-c cur` | Use `cur` as prefix filter |
| `-R` | Raw output without filtering |
| `-C dir` | Evaluate in specified directory |
| `-P prefix` | Prepend prefix to completions |
| `-F sep` | Set separator for word splitting |
| `-x cmd` | Call exported generator `_comp_xfunc_CMD_compgen_NAME` |
| `-i cmd` | Call internal generator `_comp_cmd_CMD__compgen_NAME` |

### Legacy to Modern Function Migration

| Deprecated | Modern |
|------------|--------|
| `_init_completion` | `_comp_initialize` |
| `_filedir` | `_comp_compgen -a filedir` |
| `_get_cword` | `_comp_get_words cur` |
| `_count_args` | `_comp_count_args` |
| `have` | `_comp_have_command` |

## Controlling the Completion Menu

### Menu Completion (Cycle Through Matches)

By default, TAB shows all matches or inserts the common prefix. To cycle through matches one at a time:

```bash
# In .inputrc
TAB: menu-complete
"\e[Z": menu-complete-backward   # Shift-Tab

# Or via bind
bind 'TAB:menu-complete'
bind '"\e[Z": menu-complete-backward'
```

### Readline Variables Affecting Completion Display

| Variable | Default | Effect |
|----------|---------|--------|
| `show-all-if-ambiguous` | `off` | Show all matches on first TAB (instead of bell) |
| `show-all-if-unmodified` | `off` | Show all matches if no unique prefix remains |
| `completion-query-items` | `100` | Prompt before showing more than N completions |
| `completion-display-width` | `-1` | Screen columns for match display; `0` = one per line |
| `page-completions` | `on` | Use internal pager for long completion lists |
| `print-completions-horizontally` | `off` | Sort horizontally instead of vertically |
| `colored-stats` | `off` | Use LS_COLORS to indicate file types |
| `visible-stats` | `off` | Append file-type indicator character |
| `colored-completion-prefix` | `off` | Color the common prefix differently |
| `completion-prefix-display-length` | `0` | Ellipsize common prefix longer than N |
| `menu-complete-display-prefix` | `off` | Show common prefix before cycling |
| `skip-completed-text` | `off` | Don't duplicate text after cursor on mid-word completion |
| `mark-directories` | `on` | Append `/` to directory completions |
| `mark-symlinked-directories` | `off` | Append `/` to symlinks pointing to directories |
| `match-hidden-files` | `on` | Include dotfiles in completion |
| `completion-ignore-case` | `off` | Case-insensitive matching |
| `completion-map-case` | `off` | Treat `-` and `_` as equivalent (with `completion-ignore-case`) |
| `expand-tilde` | `off` | Perform tilde expansion during completion |
| `disable-completion` | `off` | Disable completion entirely |

See [references/readline.md](references/readline.md) for full Readline variable reference.

## Debugging Completion Functions

### Trace Execution

```bash
_mycomp() {
    set -x  # Enable tracing
    # ... completion logic ...
    set +x  # Disable tracing
}
```

### Log to File

```bash
_mycomp() {
    echo "COMP_WORDS: ${COMP_WORDS[@]}" >> ~/comp-debug.log
    echo "COMP_CWORD: $COMP_CWORD" >> ~/comp-debug.log
    echo "cur: $cur" >> ~/comp-debug.log
    # ...
}
```

### Test with compgen

```bash
# Source the completion
source /usr/share/bash-completion/completions/mytool

# Test compgen directly
compgen -W "start stop restart" -- "re"
# Output: restart

# Simulate completion function
_mytool_completion; echo "${COMPREPLY[@]}"
```

### Inspect Registered Completions

```bash
# Print all compspecs
complete -p

# Print compspec for specific command
complete -p mytool

# Print compopt settings
compopt mytool
```

## Installing Completions

### User-Specific

```bash
~/.local/share/bash-completion/completions/mytool
```

### System-Wide

```bash
/usr/share/bash-completion/completions/mytool
```

### Legacy

```bash
/etc/bash_completion.d/mytool
```

### In .bashrc

```bash
# For completions not in standard directories
source /path/to/mytool-completion.bash
```

## References

- [GNU Bash Manual: Programmable Completion](https://www.gnu.org/software/bash/manual/html_node/Programmable-Completion.html)
- [GNU Bash Manual: Programmable Completion Builtins](https://www.gnu.org/software/bash/manual/html_node/Programmable-Completion-Builtins.html)
- [GNU Bash Manual: A Programmable Completion Example](https://www.gnu.org/software/bash/manual/html_node/A-Programmable-Completion-Example.html)
- [bash-completion GitHub Repository](https://github.com/scop/bash-completion)
- [bash-completion Architecture (DeepWiki)](https://deepwiki.com/scop/bash-completion/3.1-completion-architecture)
- [bash-completion Helper Functions (DeepWiki)](https://deepwiki.com/scop/bash-completion/2.2-helper-functions)

## Related Skills

- **references/readline.md** — Readline variables and key bindings that control completion display
- **references/quoting-expansion.md** — How quoting and COMP_WORDBREAKS affect completion
- **references/execution.md** — How completion functions execute (subshell vs current shell)
- **references/startup.md** — Where to install completions and when they're loaded
