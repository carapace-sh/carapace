# Zsh Completion System

In-depth reference for zsh's completion system (compsys) — the initialization, registration, builtins, utility functions, zstyle configuration, tag system, matcher specifications, caching, and how external tools hook into the completion pipeline.

## The Completion Flow

When the user presses TAB (or another completion key), zsh:

1. The ZLE widget (e.g., `complete-word`, `expand-or-complete`) triggers the completion system
2. `_main_complete` is called as the shell-side entry point
3. The completer loop iterates through functions defined in the `completer` style (defaults to `_complete _ignored`)
4. The command-specific function is looked up via the `_comps` associative array and executed
5. The function calls `compadd` or higher-level utilities (`_describe`, `_arguments`, etc.) to register candidates
6. Zsh's completion code builds the "unambiguous" string and handles insertion
7. Matches are displayed in the completion UI

### Completer Functions (Control Flow)

The `completer` zstyle determines which strategies are tried, in order:

```zsh
zstyle ':completion:*' completer _complete _ignored _approximate
```

| Completer | Purpose |
|----------|---------|
| `_complete` | Normal contextual completion |
| `_approximate` | Allow corrections (controlled by `max-errors` style) |
| `_correct` | Generate corrections only (spell-check) |
| `_expand` | Expand word (variables, braces, etc.) |
| `_expand_alias` | Expand aliases |
| `_history` | Complete from command history |
| `_ignored` | Restore matches suppressed by `ignored-patterns` |
| `_list` | Delay insertion until second completion attempt |
| `_match` | Pattern matching completion |
| `_menu` | Enable menu completion |
| `_oldlist` | Use existing completion list with new completions |
| `_prefix` | Complete ignoring suffix |
| `_user_expand` | User-defined expansions |
| `_all_matches` | Add string with all matches |

Each completer is tried in sequence; if one produces matches, later completers may still be tried depending on configuration.

## compinit — Initialization

`compinit` initializes the completion system. It must be called before any completions work.

```zsh
autoload -U compinit && compinit
```

### What compinit Does

1. Loads completion functions from directories in `$fpath`
2. Scans for `#compdef` directives in files to build the `_comps` mapping
3. Creates the `_main_complete` dispatcher
4. Sets up keybindings (TAB → `complete-word`)
5. Initializes the completion cache (`.zcompdump`)

### compinit Options

| Option | Description |
|--------|-------------|
| `-D` | Disable dump file creation |
| `-d dumpfile` | Specify dump file location (default: `$ZDOTDIR/.zshcompdump` or `~/.zcompdump`) |
| `-C` | Skip security checks, use cached dump if it exists (fastest) |
| `-u` | Use insecure files without asking |
| `-i` | Silently ignore insecure files |
| `-n` | Skip function check (don't check if functions exist) |

### Security Model

`compinit` checks that completion files are owned by root or the current user, and that parent directories are not world- or group-writable. Insecure files are reported and skipped unless `-u` or `-i` is given.

Audit insecure files:

```zsh
compaudit | xargs chmod g-w
```

### The Dump File (zcompdump)

`compinit` caches its mapping in a dump file for faster subsequent invocations. The dump stores the `_comps`, `_services`, `_patcomps`, and `_postpatcomps` associative arrays.

Optimized initialization pattern:

```zsh
autoload -Uz compinit
if [[ -n ~/.zcompdump(#qN.mh+24) ]]; then
  compinit
else
  compinit -C
fi
```

This only does a full scan if the dump file is older than 24 hours.

## compdef — Registration

`compdef` associates completion functions with commands or contexts.

### Syntax Forms

```zsh
compdef function command...           # Register function for commands
compdef -d command...                 # Delete completion for commands
compdef -k function style keyseqs...  # Create widget with key bindings
compdef -K function name style keyseqs...  # Create multiple widgets
compdef -p pattern function           # Pattern-based (tried before exact)
compdef -P pattern function           # Pattern-based (tried after exact)
compdef -n                            # Don't overwrite existing
compdef newcmd=oldcmd                 # Copy completions from oldcmd
```

### The `#compdef` Directive

Completion script files can declare their target commands via a first-line directive:

```zsh
#compdef git
#compdef git stg            # Multiple commands
#compdef -p git*             # Pattern-based
#compdef -k _complete complete-word ^I  # Key binding
```

`compinit` scans for these directives when building the `_comps` mapping.

### The `_comps` Associative Array

Internally, `compdef` populates the `_comps` associative array:

```zsh
echo $_comps[git]   # prints _git
```

This is the primary lookup table used during completion dispatch.

## Special Parameters During Completion

When a completion function is invoked, zsh sets these parameters:

| Parameter | Type | Description |
|-----------|------|-------------|
| `words` | Array | All words on the command line (1-based indexing) |
| `CURRENT` | Integer | 1-based index of the word being completed in `words` |
| `PREFIX` | String | Text before cursor in the current word |
| `SUFFIX` | String | Text after cursor in the current word |
| `IPREFIX` | String | Initial prefix (ignored for matching but inserted) |
| `ISUFFIX` | String | Initial suffix (ignored for matching but inserted) |
| `QIPREFIX` | String | Quoted version of IPREFIX |
| `QISUFFIX` | String | Quoted version of ISUFFIX |
| `curcontext` | String | Current context string (e.g., `:complete:git:`) |
| `NORMARG` | Integer | Position of first normal (non-option) argument |
| `opt_args` | Assoc Array | Parsed option arguments (key: option name, value: argument) |
| `line` | Array | Command and arguments array (set by `_arguments`) |
| `state` | Array | State strings from `->state` actions in `_arguments` |
| `state_descr` | Assoc Array | Descriptions for each state |
| `context` | Array | Context strings for state processing |
| `service` | String | Current service being completed |
| `_compskip` | String | Set to skip further completion (`-`, `*`, `patterns`) |
| `expl` | Array | Explanation options array (by convention, set by `_description`) |

### Key Differences from Bash

Unlike bash's `COMP_WORDS` / `COMP_CWORD`:
- Zsh's `words` is 1-based (not 0-based)
- Zsh preserves quoting in `words` — no `COMP_WORDBREAKS`-style word splitting
- `PREFIX` / `SUFFIX` give the partial word around the cursor, not the full word
- `IPREFIX` / `ISUFFIX` allow ignoring parts of the word for matching while still inserting them

## compadd — The Core Completion Builtin

`compadd` is the low-level builtin that adds completion candidates. All higher-level utilities (`_describe`, `_arguments`, etc.) ultimately call `compadd`.

### Syntax

```zsh
compadd [ -akqQfenUl12C ] [ -F array ]
        [ -P prefix ] [ -S suffix ]
        [ -p hidden-prefix ] [ -s hidden-suffix ]
        [ -i ignored-prefix ] [ -I ignored-suffix ]
        [ -W file-prefix ] [ -d array ]
        [ -J name ] [ -V name ] [ -X explanation ] [ -x message ]
        [ -r remove-chars ] [ -R remove-func ]
        [ -D array ] [ -O array ] [ -A array ]
        [ -E number ]
        [ -M match-spec ] [ -- ] [ words ... ]
```

### Display and Insertion Control

| Option | Description |
|--------|-------------|
| `-P prefix` | String inserted before each match; not part of match, metacharacters not quoted |
| `-S suffix` | String inserted after each match (e.g., `-S '='` for `--flag=`) |
| `-p hidden-prefix` | Inserted before match but not shown in listings; must be matched unless `-U` given |
| `-s hidden-suffix` | Inserted after match but hidden from listings |
| `-i ignored-prefix` | Inserted before `-P` string; part of match but not shown |
| `-I ignored-suffix` | Like `-i`, but for suffix |
| `-d array` | Per-match display strings; shows array element instead of word |
| `-l` | With `-d`: display strings listed one per line, not in columns |
| `-o` | With `-d`: order output by match strings (not display strings) |
| `-q` | Suffix from `-S` auto-removed if next char is blank, inserts nothing, or same char |
| `-r remove-chars` | Suffix removed if next char inserts one of these characters (character class syntax) |
| `-R remove-func` | Function called after suffix accepted; passed suffix length as argument |
| `-Q` | Don't quote metacharacters when inserting (completer handles its own quoting) |
| `-n` | Words are matches but don't appear in listings (hidden matches) |

### Grouping and Explanation

| Option | Description |
|--------|-------------|
| `-J name` | Name of sorted group for matches |
| `-V name` | Like `-J` but for unsorted group (preserves insertion order) |
| `-1` | With `-V`: only remove consecutive duplicates |
| `-2` | With `-J` or `-V`: keep all duplicates |
| `-X explanation` | String printed above matches listing |
| `-x message` | Like `-X` but printed even if no matches |

### Match Generation

| Option | Description |
|--------|-------------|
| `-a` | Words are array names; matches are their values |
| `-k` | Words are associative array names; matches are their keys |
| `-f` | Matches marked as filenames; adds slashes for directories if `LIST_TYPES` set |
| `-e` | Matches are parameter names for parameter expansion |
| `-F array` | Patterns; words matching are ignored |
| `-M match-spec` | Local match specification (overrides `matcher` style) |
| `-U` | All words accepted; no matching performed against PREFIX |
| `-W file-prefix` | Prepended to matches with `-p` to form complete filename |

### Array Manipulation

| Option | Description |
|--------|-------------|
| `-O array` | Words NOT added to completions; matching words stored in array |
| `-A array` | Like `-O` but stores strings generated by completion code (after match specs) |
| `-D array` | Words NOT added; non-matching words removed from array |
| `-C` | Adds special match expanding to all other matches when inserted |
| `-E number` | Adds `number` empty matches; implies `-V` and `-2` |

### Suffix Removal Examples

The `-r` option uses character class syntax:

```zsh
compadd -S '=' -r '= \t\n\-' -- options   # Remove '=' suffix on space, tab, newline, or nothing
compadd -S '/' -r '/ \t' -- dirs           # Remove '/' suffix on space or tab
compadd -S '=' -q -- options               # -q shorthand for common -r pattern
```

`\-` in the character class stands for characters that insert nothing (e.g., when another completion is triggered immediately).

## compset — Modifying Completion Parameters

`compset` simplifies modification of the special parameters (`PREFIX`, `SUFFIX`, `IPREFIX`, `ISUFFIX`, etc.). It returns 0 if the modification was performed.

### Syntax

```zsh
compset -p number
compset -P [ number ] pattern
compset -s number
compset -S [ number ] pattern
compset -n begin [ end ]
compset -N beg-pat [ end-pat ]
compset -q
```

| Option | Description |
|--------|-------------|
| `-p number` | If `PREFIX` > `number` chars, first `number` chars moved to `IPREFIX` |
| `-P [number] pattern` | Match `PREFIX` with pattern; matched portion moved to `IPREFIX`. Without `number`, longest match taken. With negative `number`, `number`'th longest match moved |
| `-s number` | Transfer last `number` chars from `SUFFIX` to `ISUFFIX` |
| `-S [number] pattern` | Match last portion of `SUFFIX`; transfer to `ISUFFIX` |
| `-n begin [end]` | If `CURRENT` >= `begin`, remove words up to `begin`'th; decrement `CURRENT`. If `end` given, also remove words from position `end` onwards |
| `-N beg-pat [end-pat]` | Remove words up to matching `beg-pat` if before `CURRENT`. With `end-pat`, also remove from that word onwards if after cursor |
| `-q` | Split current word on spaces respecting quoting; store in `words` array; update `CURRENT`, `PREFIX`, `SUFFIX`, `QIPREFIX`, `QISUFFIX` |

### Common Pattern: Stripping Option Prefix

```zsh
if compset -P '*\='; then
    # After --option=, complete the value
    _values "option value" ...
fi
```

This moves everything up to and including `=` into `IPREFIX`, so `PREFIX` now contains only the value part being completed.

## compquote and comptilde

### compquote

```zsh
compquote string...
```

`compquote` only succeeds when called from within the completion system. It is used as a guard to conditionally execute completion functions:

```zsh
compquote '' 2>/dev/null && _example_completion
```

This prevents the completion function from running when sourced outside a completion context (e.g., in an interactive shell).

### comptilde

```zsh
comptilde -q name
```

Handles tilde expansion within completion. Expands `~user` prefixes in completion candidates. Called internally by the completion system; rarely used directly in completion functions.

## _describe — Simple Completions with Descriptions

`_describe` is the primary utility for completions where candidates have fixed descriptions. It aligns candidates in columns and groups those with the same description together.

### Syntax

```zsh
_describe [-12JVx] [ -oO | -t tag ] descr name1 [ name2 ] [ opt ... ] [ -- name1 opt ... ]
```

### Options

| Option | Description |
|--------|-------------|
| `-t tag` | Specify tag name (default: `values` or `options` with `-o`) |
| `-o` | Complete as command options (handles `prefix-hidden`, `prefix-needed`, `verbose` styles) |
| `-O` | Like `-o` but doesn't handle `prefix-needed` style |
| `-V name` | Unsorted group (passed to `compadd -V`) |
| `-J name` | Named sorted group (passed to `compadd -J`) |
| `-1` | Only remove consecutive duplicates (with `-V`) |
| `-2` | Keep all duplicates |
| `-x` | Description shown even if no matches |

### Display Array Format

The array elements use colons to separate display text from description:

```zsh
local -a options
options=(
  'checkout:Switch branches or restore working tree files'
  'commit:Record changes to the repository'
  'push:Update remote refs along with associated objects'
)
_describe 'command' options
```

**Literal colons** in display text must be escaped with `\:`:

```zsh
options=('http\://example.com:URL with colon in display')
```

### Multiple Groups

Separate groups with `--`:

```zsh
local -a subcmds topics
subcmds=('checkout:Switch branches' 'commit:Record changes')
topics=('rebase:Reapply commits' 'merge:Join histories')
_describe 'command' subcmds -- topics
```

Each group gets its own tag and can be styled independently via zstyle.

### How _describe Works Internally

`_describe` uses C-level utilities (`cd_group` and `cd_calc` in `computil.c`) to format lists and align them into columns. It:

1. Collects candidates with the same description together on the same row
2. Ensures descriptions are neatly aligned in columns
3. Calls `compadd` with appropriate `-d`, `-J`/`-V`, `-X` options
4. Supports separate display and value arrays (display shown to user, value inserted)

### Separate Display and Value Arrays

```zsh
local -a displays values
displays=('v1.0 (stable)' 'v2.0 (latest)')
values=('v1.0' 'v2.0')
_describe 'version' displays values
```

When two arrays are given, the first is displayed and the second is inserted. This is how carapace provides human-readable `display:description` pairs while inserting clean values.

## _arguments — Option and Argument Completion

`_arguments` is the workhorse utility for Unix-style commands with options and positional arguments. It parses the command line and dispatches to the appropriate completion action.

### Syntax

```zsh
_arguments [ -nswWCRS ] [ -A pat ] [ -O name ] [ -M matchspec ] spec ...
```

### Global Options

| Option | Description |
|--------|-------------|
| `-n` | Set `NORMARG` to position of first normal argument |
| `-s` | Enable option stacking for single-letter options (e.g., `-xy` for `-x -y`) |
| `-w` | Allow option stacking even when options take arguments |
| `-W` | Allow option stacking after an argument in the same word |
| `-S` | Do not complete options after `--` on the line |
| `-A pat` | Stop completing options after first non-option argument matching `pat` |
| `-C` | Modify `curcontext` for `->state` actions |
| `-R` | Return status 300 for `->state` actions |
| `-O name` | Pass array elements as arguments to action functions |
| `-M matchspec` | Match specification for option names and values |

### Option Spec Formats

**Simple flags:**

```zsh
'-h[Show help]' '--verbose[Enable verbose output]'
```

**Options with arguments:**

```zsh
'-f[input file]:filename:_files'
'--port[Port number]:port:(8000 8080 3000)'
'-o[output file]:output file:_files -g "*.txt"'
```

**Option argument follows in same word (no space):**

```zsh
'-f+[follow symlinks]:target:_files'   # -ftarget (no space)
'-f=[use equals sign]:target:_files'   # -f=target
'-f-[must use equals]:target:_files'   # -f=target only
```

**Mutually exclusive options:**

```zsh
'(-q --quiet)'{-v,--verbose}'[Enable verbose output]'
'(-v --verbose)'{-q,--quiet}'[Suppress output]'
```

The exclusion list in parentheses specifies which other options/arguments are incompatible.

**Combined short+long options with brace expansion:**

```zsh
{-v,--verbose}'[Enable verbose output]'
{-q,--quiet}'[Suppress output]'
```

### Positional Argument Specs

| Format | Meaning |
|--------|---------|
| `1:message:action` | 1st argument (required) |
| `2::message:action` | 2nd argument (optional) |
| `*:message:action` | All remaining arguments |
| `*::message:action` | All remaining (words modified to normal args) |
| `*:::message:action` | Covered words only |

### Action Types

| Action | Description |
|--------|-------------|
| `(item1 item2)` | Literal list of matches |
| `((item1\:'desc1' item2\:'desc2'))` | List with descriptions (colons escaped) |
| `->state` | Set `$state` for case-statement dispatch |
| `function_name` | Call completion function |
| `{eval-string}` | Evaluate shell code to generate matches |
| `= action` | Insert dummy word for parsing; adds to `$words` |
| `( )` | Required but no matches (message only) |
| ` ` (single space) | No completions, display message only |

### State Machine Pattern

The `->state` action is the most powerful feature of `_arguments`. It sets the `$state` array and `$context`, `$line`, `$opt_args` parameters, then returns control to the calling function:

```zsh
_arguments -C \
  '1:subcommand:->subcommand' \
  '*:argument:->argument' \
  && return 0

case "$state" in
  subcommand)
    local -a subcmds
    subcmds=('checkout:Switch branches' 'commit:Record changes')
    _describe 'subcommand' subcmds
    ;;
  argument)
    case "$words[2]" in
      checkout)
        _git_branches
        ;;
      commit)
        _files
        ;;
    esac
    ;;
esac
```

### _arguments Sets These Variables

| Variable | Description |
|----------|-------------|
| `$state` | Array of state strings from `->state` actions |
| `$state_descr` | Associative array of descriptions for each state |
| `$line` | Command and arguments array |
| `$opt_args` | Associative array of parsed option arguments (key: option, value: argument) |
| `$context` | Context array for state processing |

## Other Utility Functions

### _alternative — Mixed Completion Types

`_alternative` combines different completion types for a single argument position. Each spec is `'tag:description:action'`:

```zsh
_alternative \
  'args:custom arg:(a b c)' \
  'files:filename:_files' \
  'dirs:directory:_files -/'
```

Unlike `_describe`, `_alternative` can call other completion functions as actions. Actions starting with a space (e.g., `' _values ...'`) are called without standard compadd options.

### _values — Keyword-Value Pairs and Comma-Separated Lists

```zsh
_values [-s sep] [-S sep] [-wC] description spec...
```

For completing arbitrary keywords/values, optionally with arguments:

```zsh
_values -s , 'flags' 'a[flag A]' 'b[flag B]' 'c[flag C]'
```

The `-s sep` option sets the separator character for multi-value words. Sets `$val_args` associative array for parsed values.

### _multi_parts — Path-Like Completion with Fixed Separator

```zsh
_multi_parts separator array
```

Completes parts of words separately, where each part is separated by a fixed character. Ideal for partial path completion:

```zsh
_multi_parts / '(usr/local/bin usr/local/sbin usr/bin usr/sbin)'
```

The user can complete one path component at a time.

### _sep_parts — Different Separators at Different Parts

```zsh
_sep_parts '(foo bar)' @ '(news ftp)' : '(woo laa)'
```

Like `_multi_parts` but allows different separators at each position. Completes `foo@news:woo`, `bar@ftp:laa`, etc.

### _regex_arguments — Regex-Based Completion

```zsh
_regex_arguments name spec...
```

For complex command lines with multiple possible argument sequences. Specs use `/pattern/:tag:descr:action` format:

```zsh
_regex_arguments _cmd \
  /$'[^\0]##\0/ \
  \( /$'word1(a|b|c)\0/' ':word:first word:(word1a word1b word1c)' \| \
     /$'word2(a|b|c)\0/' ':word:second word:(word2a word2b word2c)' \)
```

The `_regex_words` helper simplifies specification:

```zsh
local -a firstword
_regex_words subcmds 'Subcommands' 'cmd1:first command' 'cmd2:second command'
firstword="$reply[@]"
_regex_arguments _cmd /$'[^\0]##\0/ "$firstword[@]"
_cmd "$@"
```

### _combination — Combined Value Completion

```zsh
_combination [-s pattern] tag style spec... field opts...
```

Completes combinations of values (e.g., `user@host` pairs). Uses zstyle `style` value for lookup.

### _sequence — Delimited List from Another Function

```zsh
_sequence -s sep function
```

Wraps another completion function to handle comma-separated (or other delimiter) lists.

### _path_files and _files — File Completion

```zsh
_path_files [-f] [-/] [-g pattern] [-W paths] [-F ignored]
_files  # wrapper that respects file-patterns style
```

`_path_files` is the low-level file completion function. `_files` is a wrapper that respects the `file-patterns` and `list-dirs-first` styles.

### _gnu_generic — Automatic Option Completion from --help

```zsh
compdef _gnu_generic mycommand
```

Parses `--help` output to extract long options. Works for commands following GNU coding standards.

### _message — Display Informational Text

```zsh
_message [-r12] [-VJ group] descr
```

Displays text in the completion area when no completions are available. The `-r` flag treats the argument as a raw string (not a format string).

```zsh
_message -r "No matches found"
```

## The zstyle System

`zstyle` is the configuration mechanism for the completion system. It uses pattern-matched context strings to apply settings.

### Context Format

```
:completion:function:completer:command:argument:tag
```

| Component | Description |
|-----------|-------------|
| `:completion:` | Fixed literal string identifying the completion system |
| `function` | Calling widget name (often empty for normal completion) |
| `completer` | Active completer name (without underscore, e.g., `complete`, `approximate`) |
| `command` | Command name or special context (e.g., `-default-`, `-command-`) |
| `argument` | Argument position (`argument-N` or `option-opt-N`) |
| `tag` | Classification tag (`files`, `options`, `hosts`, etc.) |

### Special Contexts

| Context | When Used |
|---------|-----------|
| `-command-` | Command position (first word) |
| `-default-` | Commands with no specific completion |
| `-value-` | Parameter value assignment |
| `-array-value-` | Array parameter value |
| `-assign-parameter-` | Parameter name in assignment |
| `-brace-parameter-` | Parameter name in `${...}` |
| `-condition-` | Inside `[[...]]` conditional |
| `-math-` | Inside `((...))` arithmetic |
| `-parameter-` | Parameter name after `$` |
| `-redirect-` | After redirection operator |
| `-subscript-` | Inside parameter subscript |
| `-tilde-` | After `~` |

### Lookup Rules

Patterns are matched by **specificity**:
1. More colon-separated components = more specific
2. Literal strings beat patterns at the same depth
3. Among patterns, more complex patterns beat `*`
4. First defined pattern wins when equal specificity

### zstyle Command Forms

```zsh
zstyle context style value...     # Set style
zstyle -t context style           # Test boolean (returns 0 if true/yes/on/1)
zstyle -m context style pattern   # Test pattern match
zstyle -s context style var       # Get as string into var
zstyle -a context style arr      # Get as array into arr
zstyle -b context style var      # Get as boolean into var
zstyle -d context style          # Delete style
zstyle -L                         # List definitions in reusable form
```

### Key Completion Styles

**Display and formatting:**

| Style | Description |
|-------|-------------|
| `list-colors` | Color specifications for completion list (uses `(#b)` extended glob patterns) |
| `list-packed` | Compact listing (boolean) |
| `list-rows-first` | Row-first listing (boolean) |
| `list-dirs-first` | List directories before files (boolean) |
| `list-prompt` | Prompt for scrolling long completion lists |
| `list-separator` | Separator between match and description |
| `max-matches-width` | Trade-off between match/description column width |
| `format` | Format string for group headers and descriptions |
| `group-name` | Name for completion group (empty string = merge all groups) |
| `group-order` | Display order for groups |
| `extra-verbose` | More verbose completion listing |
| `verbose` | Verbose completion listing (boolean) |
| `description` | Show descriptions (boolean) |
| `auto-description` | Auto-generate option descriptions |

**Matching and filtering:**

| Style | Description |
|-------|-------------|
| `matcher` | Match specification per tag |
| `matcher-list` | Global match specifications (tried in sequence) |
| `ignored-patterns` | Patterns to exclude from completion |
| `ignore-line` | Ignore words already on command line |
| `sort` | Sort matches (boolean or `yes`/`no`/`menu`) |
| `file-sort` | Sort order for files (`size`, `links`, `time`, `modification`, `access`, `change`, `name`) |
| `file-patterns` | Custom file pattern tags |
| `special-dirs` | Include `.` and `..` directories (boolean) |
| `squeeze-slashes` | Treat `//` as single `/` (boolean) |

**Insertion behavior:**

| Style | Description |
|-------|-------------|
| `add-space` | Add space after expansion (boolean or `true`/`false`/`auto`) |
| `insert` | Insert all matches unconditionally (boolean) |
| `insert-unambiguous` | Start menu only on unambiguous prefix (boolean) |
| `insert-tab` | Insert TAB when no char left of cursor (boolean) |
| `menu` | Menu completion control (boolean, `select`, `select=N`) |
| `accept-exact` | Accept exact match immediately (boolean) |
| `accept-exact-dirs` | Accept directory without checking contents (boolean) |
| `suffix` | Expand tilde/parameter without suffix (boolean) |
| `keep-prefix` | Keep tilde/parameter prefix during expansion (boolean) |

**Completer control:**

| Style | Description |
|-------|-------------|
| `completer` | List of completer functions to try |
| `max-errors` | Max errors for approximate completion |
| `condition` | Delay insertion unconditionally |
| `old-list` | Use existing completion list |
| `old-menu` | Continue with old list on standard key |
| `original` | Add original string to corrections |
| `force-list` | Always show completion list (boolean) |
| `stop` | Stop at history boundaries |

**Caching:**

| Style | Description |
|-------|-------------|
| `use-cache` | Enable completion caching (boolean) |
| `cache-path` | Cache directory (default: `~/.zcompcache`) |
| `cache-policy` | Function to determine cache validity |

**Tag control:**

| Style | Description |
|-------|-------------|
| `tag-order` | Order for trying tags |
| `tag-label` | Label for tag (for display) |
| `hidden` | Hide matches from listing (boolean) |
| `host-names` | Hostname list |
| `users` | Username list |
| `groups` | Group name list |

### The list-colors Style

`list-colors` controls per-candidate colors using `(#b)` extended glob patterns with capture groups:

```zsh
zstyle ':completion:*' list-colors ${(s.:.)LS_COLORS}
```

The `(#b)` flag enables backreferences. Pattern format:

```
=(#b)(PATTERN)(GROUP_PATTERN)=0=VALUE_COLOR=DESC_COLOR
```

When zsh displays a candidate, it matches the display text against these patterns. Capture groups allow coloring different parts of the display differently.

Custom per-command coloring:

```zsh
zstyle ':completion:*:*:git:*' list-colors '=*=32'  # green for all git completions
```

### The format Style

The `format` style controls group header formatting with prompt escapes:

```zsh
zstyle ':completion:*:descriptions' format '%F{green}-- %d --%f'
zstyle ':completion:*:options' format '%F{yellow}-- %d --%f'
```

Prompt escape sequences:

| Sequence | Meaning |
|----------|---------|
| `%F{color}` / `%f` | Foreground color on/off |
| `%K{color}` / `%k` | Background color on/off |
| `%B` / `%b` | Bold on/off |
| `%U` / `%u` | Underline on/off |
| `%S` / `%s` | Standout on/off |
| `%d` | Description text |
| `%n` | Number of matches |

## The Tag System

Tags categorize completion matches. They serve a dual role: **classification** (grouping matches by type) and **customization** (providing hooks for zstyle-based control).

### Standard Tags

`accounts`, `all-expansions`, `all-files`, `arguments`, `arrays`, `association-keys`, `bookmarks`, `builtins`, `characters`, `colormapids`, `colors`, `commands`, `contexts`, `corrections`, `cursors`, `default`, `descriptions`, `devices`, `directories`, `directory-stack`, `displays`, `domains`, `expansions`, `extensions`, `file-descriptors`, `files`, `fonts`, `fstypes`, `functions`, `globbed-files`, `groups`, `history-words`, `hosts`, `indexes`, `jobs`, `interfaces`, `keymaps`, `keysyms`, `libraries`, `limits`, `local-directories`, `manuals`, `mailboxes`, `maps`, `messages`, `modifiers`, `modules`, `my-accounts`, `named-directories`, `names`, `newsgroups`, `nicknames`, `options`, `original`, `other-accounts`, `other-files`, `packages`, `parameters`, `path-directories`, `paths`, `pods`, `ports`, `prefixes`, `printers`, `processes`, `processes-names`, `sequences`, `sessions`, `signals`, `strings`, `styles`, `suffixes`, `tags`, `targets`, `time-zones`, `types`, `urls`, `users`, `values`, `variant`, `visuals`, `warnings`, `widgets`, `windows`, `zsh-options`

### Tag Management Functions

**_tags** — Register and iterate over valid tags:

```zsh
_tags friends users hosts
while _tags; do
  _requested friends expl friend compadd alice bob && ret=0
  _requested users && _users && ret=0
  _requested hosts && _hosts && ret=0
  (( ret )) || break
done
```

**_requested** — Test if a tag is in the current set:

```zsh
_requested tag [name descr [command args...]]
```

Returns 0 if the tag is requested. Optionally calls `_description` and/or `_all_labels`.

**_wanted** — Convenience wrapper combining `_tags` + `_requested`:

```zsh
_wanted names expl 'name' compadd - alice bob
```

**_next_label** — Iterate over tag labels for `tag-order` style:

```zsh
while _next_label names expl 'name' "$@"; do
  compadd "$expl[@]" - alice bob && ret=0
done
```

**_all_labels** — Convenient interface to `_next_label`:

```zsh
_all_labels names expl 'name' compadd - alice bob
```

### The tag-order Style

Controls which tags are tried and in what order:

```zsh
zstyle ':completion:*:cd:*' tag-order 'directories' 'files'
zstyle ':completion:*' tag-order 'options' 'arguments' 'files'
```

Tags not listed are still tried after the listed ones.

## Matcher Specification

Matcher specifications control how typed text matches completion candidates. They are set via the `matcher` or `matcher-list` styles.

### Syntax

```zsh
zstyle ':completion:*' matcher-list 'spec1' 'spec2' ...
```

Each spec is tried in sequence; the first that produces matches wins.

### Specification Codes

| Code | Type | Description |
|------|------|-------------|
| `m:` | Match | Case-insensitive or character mapping |
| `l:` | Left | Match at beginning of word |
| `r:` | Right | Match at end of word |
| `b:` | Boundary | Left-anchored partial word |
| `e:` | End | Right-anchored partial word |

Uppercase variants (`M:`, `L:`, `R:`) apply to the entire line rather than individual words.

### Common Matcher Examples

```zsh
# Case-insensitive completion
zstyle ':completion:*' matcher-list 'm:{a-zA-Z}={A-Za-z}'

# Partial word matching (e.g., "1704" matches "_DSC1704.JPG")
zstyle ':completion:*' matcher-list 'r:|[._-]=* r:|=*'

# Both partial word and case-insensitive (tried in sequence)
zstyle ':completion:*' matcher-list \
  '' \
  'm:{a-zA-Z}={A-Za-z}' \
  'r:|[._-]=* r:|=*' \
  'l:|=* r:|=*'
```

The empty string `''` as the first spec means "try exact matching first."

### Matcher Specification Details

- `m:{from}={to}` — Maps characters in `from` to corresponding characters in `to`
- `l:|left=right` — If the word starts with `left`, match against `right` part
- `r:left=|right` — If the word ends with `left`, match against `right` part
- `b:left=|right` — Left-anchored partial word matching
- `e:left=|right` — Right-anchored partial word matching

## Completion Caching

For expensive completions (e.g., querying a remote service), zsh provides a caching mechanism.

### Cache Functions

```zsh
# Check if cache is valid
if _cache_invalid mycompletions; then
    # Generate completions
    local -a mycompletions
    mycompletions=( ... )
    # Store in cache
    _store_cache mycompletions mycompletions
fi

# Retrieve from cache
_retrieve_cache mycompletions
```

### Cache Configuration

```zsh
zstyle ':completion:*' use-cache on
zstyle ':completion:*' cache-path ~/.zcompcache
zstyle ':completion:*:mycommand:*' cache-policy _mycommand_cache_policy
```

The cache policy function determines when the cache is stale:

```zsh
_mycommand_cache_policy() {
    # Rebuild if older than 1 hour
    local oldp=( "$1"(mh+1) )
    (( $#oldp )) && return 0
    return 1
}
```

## How External Tools Hook Into Zsh Completion

External completion frameworks (like carapace) register themselves as completion functions that invoke an external command and parse its output.

### Pattern 1: Direct compdef Registration

```zsh
compdef _myapp_completion myapp

_myapp_completion() {
    local compline="${words[@]:0:$CURRENT}"
    local output
    output=($(myapp _carapace zsh "${words[@]:1:$CURRENT-1}" 2>/dev/null))
    # Parse output and call compadd/_describe
}
```

### Pattern 2: Wrapper Function with JSON Parsing

```zsh
_myapp_completion() {
    local output
    output=("${(@f)$(myapp _carapace zsh "${words[@]:1:$CURRENT-1}" 2>/dev/null)}")
    # Parse structured output (e.g., zstyle/message/data sections)
    local zstyle message data
    IFS=$'\001' read -r -d '' zstyle message data <<<"${output}"
    zstyle ":completion:${curcontext}:*" list-colors "${zstyle}"
    zstyle ":completion:${curcontext}:*" group-name ''
    [ -z "$message" ] || _message -r "${message}"
    # Parse tag groups and call _describe
}
```

### Key Integration Points

| Mechanism | Purpose |
|-----------|---------|
| `compdef` / `#compdef` | Register completion function for a command |
| `words` / `CURRENT` | Access command-line words and cursor position |
| `compadd` | Add raw completion candidates |
| `_describe` | Add candidates with display/value separation and descriptions |
| `zstyle list-colors` | Per-candidate coloring via `(#b)` patterns |
| `zstyle group-name ''` | Merge all tag groups into single display |
| `_message -r` | Display error/usage messages in completion area |
| `compquote` guard | Prevent execution outside completion context |
| `CARAPACE_COMPLINE` | Environment variable for raw command line (carapace-specific) |
| `CARAPACE_ZSH_HASH_DIRS` | Named directory mappings (carapace-specific) |

### The compquote Guard Pattern

```zsh
compquote '' 2>/dev/null && _myapp_completion
compdef _myapp_completion myapp
```

`compquote` only succeeds within the completion system context. This prevents the completion function from executing when the script is sourced outside of completion (e.g., in an interactive shell).

### The group-name '' Pattern

```zsh
zstyle ":completion:${curcontext}:*" group-name ''
```

Setting `group-name` to an empty string merges all tag groups into a single display group. External tools like carapace handle their own grouping via `_describe -t tag` tags, so the additional zsh group-name headers would be redundant.

## Writing Completion Functions — Best Practices

### Choosing the Right Utility

| Scenario | Utility |
|----------|---------|
| Simple static list with descriptions | `_describe` |
| Unix-style options and positional arguments | `_arguments` |
| Mixed completion types at one position | `_alternative` |
| Comma-separated or keyword values | `_values` |
| Path-like completion with separator | `_multi_parts` |
| Different separators at different parts | `_sep_parts` |
| Complex keyword-based command structure | `_regex_arguments` |
| File completion | `_files` / `_path_files` |
| Commands with `--help` | `_gnu_generic` |

### Debugging Completions

| Key Binding | Action |
|-------------|--------|
| `Ctrl+X H` | Show context names, tags, and functions (`_complete_help`) |
| `Ctrl+X ?` | Capture trace in temporary file (`_complete_debug`) |
| `Alt+2 Ctrl+X` | Even more detailed information |

Reload a completion function after editing:

```zsh
unfunction _mycmd && autoload -U _mycmd
```

### Common Patterns

**Subcommand dispatch:**

```zsh
_arguments -C \
  '1:subcommand:->subcommand' \
  '*::arg:->args' && return 0

case $state in
  subcommand)
    _describe 'subcommand' subcmds
    ;;
  args)
    case $words[1] in
      build)  _arguments ... ;;
      deploy) _arguments ... ;;
    esac
    ;;
esac
```

**Option with `=` argument:**

```zsh
_arguments \
  '--color=[color]:color:(always auto never)' \
  '--output=-[output format]:format:(json yaml)'
```

**Comma-separated values:**

```zsh
_arguments \
  '--features[features]:feature:_values -s , "feature" a b c'
```

## References

### Documentation

- [Zsh Manual: Completion System](https://zsh.sourceforge.io/Doc/Release/Completion-System.html) — official reference for `compdef`, `_describe`, `compadd`, `zstyle`, `_message`
- [Zsh Manual: Completion Widgets](https://zsh.sourceforge.io/Doc/Release/Zsh-Modules.html#Completion-Widgets) — `compadd`, `compstate`, `compquote` builtins
- [Zsh Manual: Zsh Modules — zsh/computil](https://zsh.sourceforge.io/Doc/Release/Zsh-Modules.html#The-zsh_002fcomputil-Module) — `compquote`, `compargv`
- [zshcompsys(1) man page](https://linux.die.net/man/1/zshcompsys) — completion system reference
- [zshcompwid(1) man page](https://linux.die.net/man/1/zshcompwid) — completion widgets reference

### Tutorials & Blog Posts

- [A Guide to ZSH Completion With Examples](https://thevaluable.dev/zsh-completion-guide-examples/) — comprehensive walkthrough of writing zsh completions
- [Zsh Completions HOWTO](https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org) — community guide for writing completions
- [Writing ZSH Completion Functions](https://github.com/zsh-users/zsh-completions/issues/2) — practical patterns and tips
- [Mastering ZSH: Completion System Demystified](https://blog.blackwell-systems.com/posts/zsh-completion-system-explained/) — architecture and context system explained
- [Zsh Completion Architecture (DeepWiki)](https://deepwiki.com/zsh-users/zsh/4.1-completion-architecture) — internal architecture analysis

### Source Code

- [zsh `_describe` source](https://github.com/zsh-users/zsh/blob/master/Completion/Base/Utility/_describe) — the utility function carapace's output targets
- [zsh `compadd` source](https://github.com/zsh-users/zsh/blob/master/Src/Zle/complete.c) — the core completion builtin
- [zsh `computil.c`](https://github.com/zsh-users/zsh/blob/master/Src/Zle/computil.c) — `compquote` and utility implementation
- [zsh Completion Style Guide](https://github.com/zsh-users/zsh/blob/master/Etc/completion-style-guide) — conventions for writing completions
