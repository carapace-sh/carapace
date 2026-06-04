# Fish Completion System

In-depth reference for fish's completion system — the `complete` builtin, `commandline` builtin, completion functions, autoloading, conditionals, wrapping, output format, internal data structures, and how external tools hook into the completion pipeline.

## The Completion Flow

When the user presses TAB (or another completion key), fish:

1. Identifies the command being completed and the current token
2. Searches for registered completions (compspecs) for that command
3. Evaluates condition functions (`-n`) to filter active completions
4. Invokes completion functions / command substitutions to generate candidates
5. Applies fuzzy matching and ranking via `sort_and_prioritize()`
6. If one match: inserts it directly; if multiple: opens the pager

### Compspec Search Order

Fish searches for completions in this order:

1. **Direct completions** — registered via `complete -c command_name`
2. **Path completions** — registered via `complete -p /path/to/command`
3. **Wraps chain** — transitively inherited completions from wrapped commands
4. **Default file completions** — if no `--no-files` was specified

### Match Generation Order

Within a compspec, matches are generated from multiple sources:

1. **Option arguments** (`-a` with `-s`/`-l`/`-o`) — completions for specific option arguments
2. **Non-option arguments** (`-a` without option flags) — positional argument completions
3. **Command substitutions** in `-a` — dynamic completions evaluated at completion time
4. **File completions** — added unless `--no-files` was specified
5. **Wildcard expansions** — from `complete_param_expand()`

After generation:

6. **Condition filtering** — `-n` conditions are evaluated; completions with failing conditions are removed
7. **Fuzzy matching and ranking** — `sort_and_prioritize()` ranks by match quality
8. **Deduplication** — duplicate completions are removed while retaining order

## The `complete` Builtin

```fish
complete ((-c | --command) | (-p | --path)) COMMAND [OPTIONS] [--color WHEN]
complete (-C | --do-complete) [--escape] STRING
```

### Command Specification

| Flag | Long Form | Description |
|------|-----------|-------------|
| `-c` | `--command` COMMAND | Register for command by name |
| `-p` | `--path` COMMAND | Register for command by absolute path (supports wildcards) |

### Option Definition

| Flag | Long Form | Argument | Description |
|------|----------|----------|-------------|
| `-s` | `--short-option` | SHORT | Adds a short option (single character, e.g. `-a`) |
| `-l` | `--long-option` | LONG | Adds a GNU-style long option (e.g. `--color`) |
| `-o` | `--old-option` | OPTION | Adds an old-style option (multi-char, single hyphen, e.g. `-Wall`) |

### Argument Handling

| Flag | Long Form | Argument | Description |
|------|----------|----------|-------------|
| `-a` | `--arguments` | ARGUMENTS | Adds completion candidates. Can be space-separated strings or command substitution. |
| `-k` | `--keep-order` | — | Keeps argument order instead of sorting alphabetically. Multiple `-k` calls display later arguments first. |

### File Handling

| Flag | Long Form | Description |
|------|-----------|-------------|
| `-f` | `--no-files` | This completion may NOT be followed by a filename. |
| `-F` | `--force-files` | This completion MAY be followed by a filename, even if another completion specified `--no-files`. |

### Requiring Parameters

| Flag | Long Form | Description |
|------|-----------|-------------|
| `-r` | `--require-parameter` | Option must have an argument. Next argument is the option argument. Without this, option argument must be attached (e.g. `-xFoo` or `--color=auto`). |
| `-x` | `--exclusive` | Short for `-r` and `-f` combined. |

### Description

| Flag | Long Form | Argument | Description |
|------|----------|----------|-------------|
| `-d` | `--description` | DESC | Description shown in the pager. Options with the same non-empty description are grouped as one candidate. |

### Wrapping

| Flag | Long Form | Argument | Description |
|------|----------|----------|-------------|
| `-w` | `--wraps` | WRAPPED | Inherit completions from WRAPPED command. Transitive: A→B→C means A inherits C's completions. Only works with `-c`, ignored with `-p`. |

### Conditionals

| Flag | Long Form | Argument | Description |
|------|----------|----------|-------------|
| `-n` | `--condition` | CONDITION | Shell command that must return 0 for the completion to be active. Multiple conditions are tried in order until one fails. |

### Other Options

| Flag | Long Form | Description |
|------|-----------|-------------|
| `-e` | `--erase` | Delete specified completion (or all completions for the command if no other options given) |
| `-C` | `--do-complete` STRING | Manually trigger completion engine for STRING. If no STRING, uses current commandline. |
| `--escape` | When used with `-C`, escape special characters in completions |
| `--color` WHEN | Controls syntax highlighting: `auto` (default), `always`, `never` |

### Option Styles (GNU getopt Compatibility)

Fish recognizes three option styles:

**Short options** (`-s`):
- Single character, single hyphen (`-a`)
- Can be grouped: `-la` = `-l -a`
- Option arguments: append value (`-w32`) or follow in next param (`-w 32`) if `--require-parameter` given

**Old-style options** (`-o`):
- Single or multi-character, single hyphen (`-Wall`, `-name`)
- Cannot be grouped
- Option arguments: space (`-foo null`) or `=` (`-foo=null`)

**GNU-style long options** (`-l`):
- Multi-character, double hyphen (`--colors`)
- Cannot be grouped
- Option arguments: `=` (`--quoting-style=shell`) or next param if `--require-parameter` given

### How `--arguments` Works

- **With option flags** (`-s`/`-l`/`-o`): Arguments only complete as arguments FOR those specific options
- **Without option flags**: Arguments complete for non-option arguments to the command (except when completing option argument specified with `-r`)

**Command substitutions in `-a`** should return newline-separated candidates. Each line may optionally contain a tab followed by a description:

```
argument1\tDescription for argument1
argument2\tDescription for argument2
```

These tab-separated descriptions override `-d` descriptions.

### Erasing Completions

Two ways to erase:

1. `complete -c COMMAND -e` — erases ALL completions for COMMAND
2. `complete -c COMMAND -e [options]` — erases specific completion option

### Displaying Completions

When `complete` is called without defining/erasing options, it shows matching completions:

- `complete` — shows ALL loaded completions
- `complete -c foo` — shows all completions for `foo`

Since completions are autoloaded, you must trigger them first.

### Windows/Cygwin/MSYS2 Handling

- Binary executables have `.exe` extension but extension not required when calling
- `complete -c myprog` works for both `myprog` and `myprog.exe`
- `complete -c myprog.exe` ONLY works for `myprog.exe`

## The `commandline` Builtin

```fish
commandline [OPTIONS] [CMD]
```

### Buffer Selection

| Option | Long Form | Description |
|--------|-----------|-------------|
| `-b` | `--current-buffer` | Select entire command line (default) |
| `-j` | `--current-job` | Select current job (pipeline) — stops at `;`, `&`, newlines |
| `-p` | `--current-process` | Select current process — stops at logical operators, terminators, and pipes |
| `-t` | `--current-token` | Select current token |
| `-s` | `--current-selection` | Select current selection |
| `--search-field` | | Use pager search field instead; returns false if not shown |

### Output/Cut Options

| Option | Long Form | Description |
|--------|-----------|-------------|
| `-c` | `--cut-at-cursor` | Print selection only up to cursor position. With `--tokens-expanded`, prints up to last completed token (excluding in-progress token). |
| `-x` | `--tokens-expanded` | Perform argument expansion, print one argument per line. Command substitutions are forwarded as-is. |
| `-o` | `tokenize` / `--tokens-raw` | Deprecated — do not use |
| `--input=INPUT` | | Operate on this string instead of command line (enables options like `--tokens-expanded`) |

### Cursor/Position

| Option | Long Form | Description |
|--------|-----------|-------------|
| `-C` | `--cursor` | Get/set current cursor position. If used with `-j`, `-p`, or `-t`, position is relative to that substring. |
| `-B` | `--selection-start` | Get current position of selection start |
| `-E` | `--selection-end` | Get current position of selection end |
| `-L` | `--line` | Print line number cursor is on (1-based), or set cursor to given line |
| `--column` | | Print 1-based Unicode code point offset from line start, or set cursor to given offset |

### Modification

| Option | Long Form | Description |
|--------|-----------|-------------|
| `-a` | `--append` | Append string at end (don't remove current content) |
| `-i` | `--insert` | Insert at current cursor position (don't remove content) |
| | `--insert-smart` | Insert with DWIM mode — strips `$` prefix from first command on each line |
| `-r` | `--replace` | Remove current and replace with specified string (default) |
| `-f` | `--function` | Queue arguments as input functions; cannot be combined with other options |

### Mode Queries

| Option | Long Form | Description |
|--------|-----------|-------------|
| `-S` | `--search-mode` | Returns true if command line is performing history search |
| `-P` | `--paging-mode` | Returns true if pager is showing (e.g., tab completions) |
| | `--paging-full-mode` | Returns true if pager is showing and all lines are visible (no "N more rows") |
| | `--is-valid` | Returns 0 if syntactically valid and complete; 2 if incomplete; 1 if erroneous |
| | `--showing-suggestion` | Returns true if shell is showing autosuggestion |

### Common Completion Patterns

```fish
# Get current process, tokenized, excluding in-progress token
set -l tokens (commandline -xpc)

# Get the in-progress token (what the cursor is on)
set -l current (commandline -ct)

# Combined: parse command line for completion logic
set -l tokens (commandline -xpc)
set -l current (commandline -ct)
```

### Example: Command Line Context

For the command line `echo $flounder >&2 | less; and echo $catfish` (cursor on the "o" of "flounder"):

| Command | Output |
|---------|--------|
| `commandline -t` | `$flounder` (current token) |
| `commandline -ct` | `$fl` (current token partial) |
| `commandline -b` | `echo $flounder >&2 \| less; and echo $catfish` |
| `commandline -p` | `echo $flounder >&2` (current process) |
| `commandline -j` | `echo $flounder >&2 \| less` (current job) |

### Using with `complete -C`

When `commandline` is called during a completion call via `complete -C STRING`, it considers that string to be the current command line contents. This enables testing completions:

```fish
complete -C 'git checkout m'
```

## Completion Functions (Helper Functions)

Functions beginning with `__fish_` are internal helpers. Two naming conventions:

- `__fish_print_*` — print newline-separated lists (no descriptions)
- `__fish_complete_*` — print completions with tab-separated descriptions

### Path Completions

| Function | Description |
|----------|-------------|
| `__fish_complete_directories STRING DESC` | Path completion allowing only directories |
| `__fish_complete_path STRING DESC` | Path completion with description |
| `__fish_complete_suffix SUFFIX` | File completion sorting files with given suffix first |

### System Completions

| Function | Description |
|----------|-------------|
| `__fish_complete_groups` | Lists all user groups with members as description |
| `__fish_complete_pids` | Lists all process IDs with command names as description |
| `__fish_complete_users` | Lists all users with full names as description |

### Informational Printers

| Function | Description |
|----------|-------------|
| `__fish_print_filesystems` | Lists all known file systems |
| `__fish_print_hostnames` | Lists hosts from fstab, ssh known_hosts, /etc/hosts |
| `__fish_print_interfaces` | Lists all known network interfaces |
| `__fish_print_xdg_mimetypes` | Lists XDG MIME types |
| `__fish_print_xdg_applications` | Lists XDG applications |

### Condition Functions

| Function | Description |
|----------|-------------|
| `__fish_seen_subcommand_from X ...` | Returns true if any of X was used on command line as a subcommand |
| `__fish_contains_opt -s SHORT LONG ...` | Checks if option was specified on command line |
| `__fish_is_first_arg` | Returns true if completing the first argument |
| `__fish_is_nth_arg N` | Returns true if completing the Nth argument |
| `__fish_use_subcommand` | Returns true if no subcommand has been given yet |

### Common Condition Patterns

**Subcommand-based conditions:**

```fish
set -l commands status set-time set-timezone list-timezones

# Offer subcommands when no subcommand given
complete -c timedatectl -n "not __fish_seen_subcommand_from $commands" \
    -a "status set-time set-timezone list-timezones"

# Offer timezones only for set-timezone subcommand
complete -c timedatectl -n "__fish_seen_subcommand_from set-timezone" \
    -a "(timedatectl list-timezones)"
```

**Option-dependent conditions:**

```fish
# Only show --nodeps when --erase is given
complete -c rpm -n "__fish_contains_opt -s e erase" \
    -l nodeps -d "Don't check dependencies"
```

## Completion Autoloading

Completions are autoloaded from directories in `$fish_complete_path`:

1. `~/.config/fish/completions` (user's own completions)
2. `/etc/fish/completions` (system-wide)
3. `~/.local/share/fish/vendor_completions.d` (third-party)
4. `/usr/share/fish/vendor_completions.d` (vendor-shipped)
5. `/usr/share/fish/completions` (fish shipped)
6. `~/.cache/fish/generated_completions` (auto-generated from man pages)

### File Naming

Completion files must be named `<command-name>.fish`. When a command is first invoked, fish loads the corresponding completion file.

### Registration in Completion Files

```fish
# Disable file completions for the entire command
complete -c myprog -f

# Offer subcommands
complete -c myprog -n "not __fish_seen_subcommand_from $commands" \
    -a "start stop status"

# Options with descriptions
complete -c myprog -s h -l help -d 'Print help text'
complete -c myprog -l output -r -d 'Output directory'

# Dynamic completions via command substitution
complete -c myprog -n "__fish_seen_subcommand_from set-timezone" \
    -a "(timedatectl list-timezones)"

# Wrap another command
complete -c hub -w git
```

## Completion Output Format

### Tab-Separated Format

Each line is one completion candidate. A tab character separates the **value** from the **description**:

```
value1\tdescription1
value2\tdescription2
value3\tdescription3
```

- If there is no tab, the entire line is the value with no description
- Descriptions are shown in the pager next to the value
- Options with the same non-empty description are grouped as one candidate

### How Fish Splits Completion Output

Fish splits command substitution output on **newlines only** (not on spaces). This is fundamentally different from bash's `IFS`-based splitting. This means:

- `origin/main`, `--option=value`, and `user@host` stay as single tokens
- No `COMP_WORDBREAKS`-style word-breaking problems
- Each line becomes one completion candidate

## Internal Completion Data Structures

### Completion Struct

```rust
pub struct Completion {
    pub completion: WString,       // Text to insert
    pub description: WString,      // Description shown in pager
    pub r#match: StringFuzzyMatch, // Fuzzy match metadata (rank, case info)
    pub flags: CompleteFlags,      // Behavior flags
}
```

### CompleteFlags

| Flag | Value | Description |
|------|-------|-------------|
| `NO_SPACE` | `1 << 0` | Do not append a space after completion |
| `REPLACES_TOKEN` | `1 << 1` | Completion replaces entire token, not just appends |
| `AUTO_SPACE` | `1 << 2` | Auto-detect whether to add space based on last character |
| `DONT_ESCAPE` | `1 << 3` | Insert without escaping |
| `DONT_SORT` | `1 << 4` | Keep original order |
| `DUPLICATES_ARGUMENT` | `1 << 5` | Marks completions that duplicate existing arguments |
| `REPLACES_LINE` | `1 << 6` | Replaces entire command line |
| `KEEP_VARIABLE_OVERRIDE_PREFIX` | `1 << 7` | Keep `foo=` prefix when replacing |
| `VARIABLE_NAME` | `1 << 8` | Indicates it's a variable name |
| `SUPPRESS_PAGER_PREFIX` | `1 << 9` | Suppress showing shared prefix in pager |

### CompletionReceiver

A wrapper around `Vec<Completion>` with size limits and convenience methods for adding/completing completions. Prevents unbounded memory growth during completion generation.

## NO_SPACE and AUTO_SPACE

### AUTO_SPACE Resolution

The `AUTO_SPACE` flag is resolved to `NO_SPACE` by checking the last character of the completion:

```rust
fn resolve_auto_space(comp: &wstr, mut flags: CompleteFlags) -> CompleteFlags {
    if flags.contains(CompleteFlags::AUTO_SPACE) {
        flags -= CompleteFlags::AUTO_SPACE;
        if let Some('/' | '=' | '@' | ':' | '.' | ',' | '-') = comp.as_char_slice().last() {
            flags |= CompleteFlags::NO_SPACE;
        }
    }
    flags
}
```

Characters that trigger `NO_SPACE` via `AUTO_SPACE`:

| Character | Use Case |
|-----------|----------|
| `/` | Directory paths |
| `=` | Option assignments (`--color=auto`) |
| `@` | User@host |
| `:` | Port specifications (`host:port`) |
| `.` | File extensions |
| `,` | Comma-separated lists |
| `-` | Chained flags (`-la`, `-vnc`) |

### NO_SPACE Behavior

- If only one completion exists and it's marked `NO_SPACE`, no space is added
- If multiple completions exist, space behavior depends on common prefix uniqueness
- `NO_SPACE` is set internally by fish for completions ending with the characters above
- External completion tools (like carapace) cannot set `NO_SPACE` through the tab-separated format — this is a known limitation

### REPLACES_TOKEN

When `REPLACES_TOKEN` is set, the completion replaces the entire current token rather than appending to it. This is used for:

- Variable name completions (replace `$pa` with `$PATH`)
- Completions that need to change the token prefix

## Wrapping and Inheritance

### How Wrapping Works

```fish
complete -c hub -w git  # hub inherits all git completions
```

- Wrapping is **transitive**: A→B→C means A inherits C's completions
- A command can wrap **multiple** commands
- Wrapping only works with `-c`/`--command`; ignored with `-p`/`--path`
- Wrapping is removed with `-e`/`--erase`

### Internal Implementation

- Mapping stored in `WRAPPER_MAP` (global `HashMap<WString, Vec<WString>>`)
- Recursion depth limited to 24 levels to prevent infinite loops
- Any wrapped command can disable file completions via the wrap chain

## The Pager

### Overview

The pager displays completions in a scrollable table when multiple matches exist.

### Navigation

| Key | Action |
|-----|--------|
| Arrow keys, Page Up/Down | Navigate completions |
| Tab / Shift+Tab | Move selection forward/backward |
| Ctrl+S or `/` (vi mode) | Open search menu to filter the list |

### Layout

The pager calculates a grid layout based on terminal dimensions:

- Maximum 6 columns (`PAGER_MAX_COLS`)
- Minimum 16 chars width, 4 rows height to display
- Initially shows 4 rows or half terminal height (undisclosed mode)
- Press Tab again to fully disclose all completions

### Prefix Suppression

If all completions share a common prefix, it's suppressed to save space using `CompleteFlags::SUPPRESS_PAGER_PREFIX`. For example, if all completions start with `/usr/share/`, that prefix is hidden and only the differing suffixes are shown.

### Fuzzy Matching in Pager

The pager supports fuzzy filtering when the search field is active:

- Matches against both description and completion strings
- Uses `string_fuzzy_match_string()` for matching
- Lower rank = better match
- Case sensitivity: case-insensitive unless uppercase character in search string

### Pager Color Variables

| Variable | Meaning |
|----------|---------|
| `fish_pager_color_progress` | Progress bar at bottom left |
| `fish_pager_color_background` | Background of each row |
| `fish_pager_color_prefix` | The prefix string to complete |
| `fish_pager_color_completion` | The proposed completion |
| `fish_pager_color_description` | Completion description |
| `fish_pager_color_selected_background` | Selected completion background |
| `fish_pager_color_secondary_background` | Alternating row background |

## Completion Engine Internals

### High-Level Flow

```
User presses Tab
  → ReadlineCmd::Complete dispatched
  → reader calls completion engine
  → perform_for_commandline_impl()
    → get_process_extent() — tokenize command line
    → Extract variable assignments before cursor
    → Determine if completing command name, arguments, or special cases
  → complete_cmd() — command name completions
  → complete_abbr() — abbreviation completions
  → complete_param_expand() — parameter/wildcard expansions
  → wildcard_complete() — file glob completions
  → Custom script completions (wrap chain support)
  → sort_and_prioritize()
    → Find best rank (lower is better from StringFuzzyMatch)
    → Retain only completions matching best rank
    → Deduplicate while retaining order
    → Apply natural sorting via wcsfilecmp()
  → Pager renders completions
```

### sort_and_prioritize()

This function is the core ranking and filtering step:

1. **Find best rank** — the lowest (best) `StringFuzzyMatch` rank among all completions
2. **Filter by rank** — retain only completions matching the best rank
3. **Deduplicate** — remove duplicate completions while retaining order
4. **Natural sort** — apply `wcsfilecmp()` for human-friendly ordering
5. **For autosuggestions** — prefer case matches, penalize duplicates and tilde-suffixed files

## How External Tools Hook Into Fish Completions

### The Standard Pattern

External completion tools (like carapace, zsh's compinit, etc.) register a single completion function that delegates to the external tool:

```fish
function _mytool_completion
    set --local data (commandline -cp | xargs mytool _complete fish 2>/dev/null)
    echo $data
end

complete -c mytool -f -a '(_mytool_completion)' -r
```

### Key Considerations for External Tools

1. **Output format** — must be tab-separated `value\tdescription` lines
2. **NO_SPACE limitation** — cannot communicate nospace information through the tab-separated format
3. **Open-quote handling** — fish's tokenizer uses `TOK_ACCEPT_UNFINISHED` which handles partial syntax, but command substitution via `xargs` may need retry logic for open quotes
4. **File completion** — use `-f`/`--no-files` to disable fish's default file completion, or `-F`/`--force-files` to re-enable for specific options
5. **Sorting** — fish sorts completions alphabetically by default; use `-k`/`--keep-order` to preserve custom ordering
6. **Condition evaluation** — `-n` conditions are evaluated at completion time, enabling context-sensitive completions

### Comparison with Bash's Completion Model

| Aspect | Fish | Bash |
|--------|------|------|
| Registration | `complete -c cmd -a '(...)'` | `complete -F func cmd` |
| Output format | Tab-separated `value\tdesc` | `COMPREPLY` array |
| Word breaking | No word breaks (splits on newlines only) | `COMP_WORDBREAKS` splits on `:`, `=`, `@`, etc. |
| Context access | `commandline -xpc`, `commandline -ct` | `COMP_WORDS`, `COMP_CWORD`, `COMP_LINE`, `COMP_POINT` |
| Conditionals | `-n "shell_command"` | Function logic |
| File completion | Default on, disable with `-f` | Default on, disable with `compopt +o default` |
| Nospace | `NO_SPACE` flag (internal only) | `compopt -o nospace` |
| Wrapping | `-w` (transitive) | No equivalent |
| Dynamic loading | Autoload from `fish_complete_path` | `bash-completion` lazy loading via exit status 124 |

## Example: Comprehensive Completion Script

```fish
# Disable file completions for the entire command
complete -c myapp -f

# Define subcommands
set -l commands status set-time set-timezone list-timezones

# Offer subcommands when no subcommand given
complete -c myapp -n "not __fish_seen_subcommand_from $commands" \
    -a "status set-time set-timezone list-timezones"

# Describe subcommands
complete -c myapp -n "not __fish_seen_subcommand_from $commands" \
    -a "status" -d "Show current time settings"
complete -c myapp -n "not __fish_seen_subcommand_from $commands" \
    -a "set-timezone" -d "Set the system timezone"

# Offer timezones for set-timezone subcommand
complete -c myapp -n "__fish_seen_subcommand_from set-timezone" \
    -a "(timedatectl list-timezones)"

# Options available only for specific subcommands
complete -c myapp -n "__fish_seen_subcommand_from set-local-rtc" \
    -l adjust-system-clock -d 'Synchronize system clock from RTC'

# Global options
complete -c myapp -s h -l help -d 'Print help text'
complete -c myapp -l version -d 'Print version string'
complete -c myapp -s v -l verbose -d 'Verbose output'

# Short option with required argument
complete -c myapp -s o -l output -r -d 'Output directory'

# Exclusive option (requires arg, no file completion)
complete -c myapp -s f -l format -x -a "json yaml xml" -d 'Output format'

# Wrap another command
complete -c myapp-wrapper -w myapp
```

## References

- [Writing your own completions — fish-shell docs](https://fishshell.com/docs/current/completions.html)
- [complete builtin — fish-shell docs](https://fishshell.com/docs/current/cmds/complete.html)
- [commandline builtin — fish-shell docs](https://fishshell.com/docs/current/cmds/commandline.html)
- [Tab Completion System — DeepWiki](https://deepwiki.com/fish-shell/fish-shell/5.4-tab-completion-system)
- [Fish shell source: src/complete.rs](https://github.com/fish-shell/fish-shell/blob/main/src/complete.rs)
- [Fish shell source: src/pager.rs](https://github.com/fish-shell/fish-shell/blob/main/src/pager.rs)

## Related Skills

- **carapace-dev** → `references/shell-fish.md` — carapace-specific fish integration (snippet, open-quote retry, value formatting, nospace limitation)
- **bash** → `references/completion.md` — bash programmable completion for comparison
