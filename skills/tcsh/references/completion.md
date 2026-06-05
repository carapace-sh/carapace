# Tcsh Programmable Completion

In-depth reference for tcsh's programmable completion system — the `complete` builtin, word specifications, list types, external command integration, and the internal completion engine.

## The Completion Flow

When the user presses TAB (or another completion key), tcsh:

1. Parses the input buffer to determine the type of word being completed
2. Searches for a matching completion specification registered via the `complete` builtin
3. If found, evaluates word specifications left-to-right to determine the completion type
4. Generates candidates from the matching list type
5. Applies the select pattern filter (if any)
6. Inserts the unique match or lists possibilities

### Word Type Detection (Built-in Logic)

Before consulting user-defined completions, tcsh determines the default completion type:

| Context | Default Type |
|---------|-------------|
| First word in buffer | Command |
| Word after `;`, `\|`, `\|&`, `&&`, `\|\|` | Command |
| Word beginning with `$` | Variable |
| Everything else | Filename |
| Empty line | Filename |

This is implemented in `tenematch()` in `tw.parse.c`:

```c
looking = starting_a_command(qline.s + word - 1, qline.s) ?
    TW_COMMAND : TW_ZERO;
```

### Completion Specification Search

When `tenematch()` identifies the word context, it calls `tw_complete()` to look up user-defined completions:

1. **Exact command match** — searches the `completions` varent tree for the command name
2. **Glob-pattern match** — the `complete` builtin allows command names to be glob patterns
3. **Ambiguous prefix match** — completions starting with `-` are used only when the command is ambiguous

If no user-defined completion is found, the built-in type (command, variable, or filename) is used.

### Match Generation Order

Within a completion specification, word specifications are evaluated **left to right**. The first matching specification wins. This is critical for ordering:

```tcsh
# CORRECT: positional first, then fallback
complete dbx 'p/2/(core)/' 'p/*/c/'

# WRONG: n-spec always matches, p/2 never reached
complete dbx 'n/*/\c/' 'p/2/(core)/'
```

## The `complete` Builtin

```
complete [command [word/pattern/list[:select][/suffix] ...]]
```

### Forms

| Invocation | Behavior |
|-----------|----------|
| `complete` | List all registered completions |
| `complete cmd` | List completions for `cmd` |
| `complete cmd 'spec' ...` | Define completions for `cmd` |

### Command Argument

The first argument is the command name to complete. It can be:

- **Exact name** — `complete cd ...` matches only `cd`
- **Glob pattern** — `complete 'co*' ...` matches any command starting with `co`
- **Ambiguous prefix** — `complete -co* ...` (leading `-`) used only when the command is ambiguous (e.g., `co` could be `compress` or `cp`)

### Removing Completions

```
uncomplete pattern
```

Removes all completions whose names match the glob pattern. `uncomplete *` removes all.

## Word Specifications

Each word specification has the form: `type/pattern/list[:select][/suffix]`

The delimiter between fields is `/`. The `@` character can also be used as a delimiter (common in practice to avoid conflicts with `/` in paths).

### Word Types

| Type | Name | Pattern Meaning | When It Matches |
|------|------|----------------|-----------------|
| `p` | Positional | Numeric range (e.g., `1`, `1-3`, `*`) | When current word position is within the range |
| `c` | Current-word | Glob pattern matching beginning of current word | Pattern is **ignored** when completing the current word |
| `C` | Current-word (inclusive) | Glob pattern matching beginning of current word | Pattern is **included** when completing the current word |
| `n` | Next-word | Glob pattern matching beginning of **previous** word | When previous word starts with pattern |
| `N` | Next-next-word | Glob pattern matching beginning of **word two before** current | When the word two positions back starts with pattern |

### Positional Pattern Syntax

The `p` type uses numeric ranges identical to shell variable indexing:

| Pattern | Meaning |
|---------|---------|
| `1` | First argument only |
| `2` | Second argument only |
| `1-3` | Arguments 1 through 3 |
| `*` | Any argument position |
| `0` | The command word itself |

### Current-Word Pattern Examples

```tcsh
# Complete options (words starting with -) for find
complete find 'c/-/(name newer user group type)/'

# Complete -I option prefix with directories
complete cc 'c/-I/d/'

# Complete @ prefix with hostnames
complete finger 'c/*@/$hostnames/'
```

The `c` type is commonly used for option completion because the pattern matches the prefix of the current word (e.g., `-` for options), but is ignored when the word doesn't start with that prefix.

### Next-Word Pattern Examples

```tcsh
# When previous word is -name, complete with filenames
complete find 'n/-name/f/'

# When previous word is -user, complete with usernames
complete find 'n/-user/u/'

# When previous word is -group, complete with group names
complete find 'n/-group/g/'
```

## List Types

The list type determines what candidates are generated:

### Built-in List Types

| Type | Completion Source | Internal Constant |
|------|-------------------|-------------------|
| `a` | Aliases | `TW_ALIAS` |
| `b` | Editor bindings (editor commands) | `TW_BINDING` |
| `c` | Commands (builtin or external) | `TW_COMMAND` |
| `C` | External commands with path prefix | `TW_PATH \| TW_COMMAND` |
| `d` | Directories | `TW_DIRECTORY` |
| `D` | Directories with path prefix | `TW_PATH \| TW_DIRECTORY` |
| `e` | Environment variables | `TW_ENVVAR` |
| `f` | Filenames (all files) | `TW_FILE` |
| `F` | Filenames with path prefix | `TW_PATH \| TW_FILE` |
| `g` | Group names | `TW_GRPNAME` |
| `j` | Jobs | `TW_JOB` |
| `l` | Resource limits | `TW_LIMIT` |
| `n` | Nothing (no completions) | `TW_NONE` |
| `s` | Shell variables | `TW_SHELLVAR` |
| `S` | Signal names | `TW_SIGNAL` |
| `t` | Plain text files | `TW_TEXT` |
| `T` | Plain text files with path prefix | `TW_PATH \| TW_TEXT` |
| `u` | Usernames | `TW_USER` |
| `v` | Any variables (shell + environment) | `TW_VARIABLE` |
| `x` | Like `n`, but prints select message when listing | `TW_EXPLAIN` |
| `X` | Completions (other completion specs) | `TW_COMPLETION` |

### Dynamic List Types

| Type | Syntax | Description |
|------|--------|-------------|
| Variable | `$varname` | Words from the shell variable `varname` |
| Word list | `(word1 word2 ...)` | Explicit list of words |
| Command | `` `command` `` | Words from command output |

### Variable List Type

```tcsh
# Complete from a shell variable
set hostnames = (localhost example.com remote.host)
complete ftp 'p/1/$hostnames/'
```

The variable is expanded at completion time, so changes to the variable take effect immediately.

### Word List Type

```tcsh
# Complete with explicit words
complete dbx 'p/2/(core)/'

# Complete find -type with specific types
complete find 'n/-type/(b c d f l p s)/'

# Complete find -fstype
complete find 'n/-fstype/(nfs 4.2)/'
```

### Command List Type (Backtick)

```tcsh
# Complete kill with PIDs from ps
complete kill 'p/*/`ps | awk \{print \$1\}`/'

# Complete with custom script
complete mycmd 'p/*/`my_completion_script`/'
```

When the backtick list type is used, tcsh:

1. Sets the `COMMAND_LINE` environment variable to the current command line
2. Executes the command via command substitution
3. Splits the output into words at blanks, tabs, and newlines
4. Discards null words
5. Uses the resulting words as completion candidates
6. Unsets `COMMAND_LINE`

## COMMAND_LINE Environment Variable

The `COMMAND_LINE` variable is the primary mechanism for external completion generators to hook into tcsh. It is set only during backtick command execution.

### What It Contains

`COMMAND_LINE` contains the **entire current command line** as typed by the user, including the command name and all arguments up to (and including) the word being completed.

### What It Does NOT Contain

Unlike bash's `COMP_LINE` + `COMP_POINT`, tcsh provides **only** the raw command line string. There is:

- **No cursor position** — no equivalent to bash's `COMP_POINT`
- **No word array** — no equivalent to bash's `COMP_WORDS` / `COMP_CWORD`
- **No completion type** — no equivalent to bash's `COMP_TYPE`
- **No word-break characters** — no equivalent to bash's `COMP_WORDBREAKS`

External completion generators must parse `COMMAND_LINE` themselves to determine:

- Which argument position is being completed
- What the partial word is
- The cursor position (if needed)

### Example: Custom Completion Script

```tcsh
complete myapp 'p/*/`myapp-completer`/'
```

The `myapp-completer` script receives `COMMAND_LINE` in its environment:

```sh
#!/bin/sh
# myapp-completer - completion generator for myapp
# COMMAND_LINE is set by tcsh to the full command line

# Parse the command line to find current word
set -- $COMMAND_LINE
cmd="$1"
shift

# Last argument is the word being completed
current_word="$_"

# Generate completions based on context
case "$#" in
  0) echo "help version config" ;;
  1) echo "init build test deploy" ;;
  *) echo "--verbose --quiet --help" ;;
esac
```

### Comparison with Bash's COMP_* Variables

| Feature | Bash | Tcsh |
|---------|------|------|
| Full command line | `COMP_LINE` | `COMMAND_LINE` |
| Cursor position | `COMP_POINT` | **Not provided** |
| Word array | `COMP_WORDS` | **Not provided** |
| Current word index | `COMP_CWORD` | **Not provided** |
| Completion type | `COMP_TYPE` | **Not provided** |
| Word break chars | `COMP_WORDBREAKS` | **Not provided** (see `wordchars`) |
| Key that triggered | `COMP_KEY` | **Not provided** |

## Select Pattern

The optional select pattern filters which candidates from the list are considered. It is a glob pattern appended after `:`:

```tcsh
# Only .c, .a, .o files
complete cc 'p/*/f:*.[cao]/'

# Exclude .c, .h, .tex files
complete rm 'p/*/f:^*.{c,h,tex}/'

# Only files matching *.py
complete python 'p/*/f:*.py/'
```

When a select pattern is provided, the `fignore` shell variable is **ignored** for that completion.

### Negation in Select Patterns

The `^` prefix negates the glob pattern, excluding matches rather than including them:

```tcsh
# Complete with all files EXCEPT .o and .a
complete vim 'p/*/f:^*.{o,a}/'
```

## Suffix Field

The optional suffix field controls what character is appended after a successful completion:

| Suffix Value | Behavior |
|-------------|----------|
| Omitted | `/` appended to directories, space to other words |
| Empty (null) | No character appended |
| Single character | That character is appended |

### Examples

```tcsh
# Append '@' after username completion
complete finger 'c/*@/$hostnames/' 'p/1/u/@'

# No suffix after completion
complete mycmd 'p/*/c/' ''

# Append ':' after completion
complete ssh 'p/1/$hostnames/:'
```

### Relationship to `addsuffix`

The `addsuffix` shell variable controls the **default** suffix behavior when the suffix field is omitted:

- **Set** (default): `/` for directories, space for other words
- **Unset**: No suffix is added regardless of the suffix field

When a suffix is explicitly specified in the completion spec, it overrides the `addsuffix` default.

## Shell Variables Affecting Completion

| Variable | Effect |
|----------|--------|
| `addsuffix` | Append `/` to directories and space to other words on exact match (default: set) |
| `autolist` | List choices when completion fails; `ambiguous` variant lists only when no new chars added |
| `autoexpand` | Run `expand-history` before each completion attempt |
| `autocorrect` | Run `spell-word` before each completion attempt |
| `complete` | Set to `enhance` for case-insensitive completion treating `-`, `_`, `.` as word separators; `Enhance` for case-insensitive only when typing lowercase/hyphen |
| `correct` | Set to `cmd` (correct commands), `complete` (complete commands), or `all` (correct entire line) |
| `fignore` | List of file suffixes to ignore during completion (e.g., `(o a class)`) |
| `listmax` | Maximum items `list-choices` will list without asking |
| `listmaxrows` | Maximum rows `list-choices` will list without asking |
| `matchbeep` | Control when completion beeps: `never`, `nomatch`, `ambiguous` (default), `notunique` |
| `nobeep` | Disable all beeping |
| `nostat` | List of directories to skip `stat(2)` calls during completion |
| `recexact` | Complete on shortest unique match even if longer matches exist |
| `recognize_only_executables` | List only executable files when listing commands |

### The `complete enhance` Variable

Setting `complete` to `enhance` enables smarter completion:

- **Case-insensitive matching** — `FOO` matches `foo`, `Foo`, etc.
- **Word separator handling** — `-`, `_`, and `.` are treated as word separators for partial matching
- This means `f-n` can match `file-name`, `file_name`, or `file.name`

Setting to `Enhance` (capital E) applies these rules only when the user types lowercase letters or hyphens.

## Spelling Correction

tcsh has built-in spelling correction that integrates with the completion system:

### The `correct` Variable

| Value | Behavior |
|-------|----------|
| `cmd` | Automatically correct misspelled command names |
| `complete` | Automatically complete command names |
| `all` | Correct the entire command line |

When correction is triggered, tcsh prompts:

```
> lz /usr/bin
CORRECT>ls /usr/bin (y|n|e|a)?
```

| Response | Action |
|----------|--------|
| `y` or space | Execute the corrected line |
| `n` or anything else | Execute the original line |
| `e` | Leave the uncorrected command in the input buffer for editing |
| `a` | Abort the command (like `^C`) |

### The `autocorrect` Variable

When set, the `spell-word` editor command is invoked automatically before each completion attempt. This can correct the current word before completion candidates are generated.

### Editor Commands for Spelling

| Command | Default Binding | Description |
|---------|-----------------|-------------|
| `spell-word` | `M-s`, `M-S` | Correct spelling of current word |
| `spell-line` | `M-$` | Correct spelling of each word in the buffer |

## Internal Completion Engine

### Source Files

| File | Purpose |
|------|---------|
| `tw.parse.c` | Core completion entry point (`tenematch()`), word parsing, quote handling |
| `tw.comp.c` | Completion engine: `docomplete()`, `tw_find()`, `tw_complete()`, `tw_result()` |
| `tw.h` | Completion type constants (`TW_COMMAND`, `TW_FILE`, etc.) |
| `complete.tcsh` | Example completions shipped with tcsh (~1277 lines) |

### `tenematch()` — The Entry Point

```c
int tenematch(Char *inputline, int num_read, COMMAND command)
```

Called when the user invokes a completion editor command. Returns:

| Return Value | Meaning |
|-------------|---------|
| `> 1` | Number of items found (multiple matches) |
| `= 1` | Exactly one match (or spelling corrected) |
| `= 0` | No match (or spelling was already correct) |
| `< 0` | Error |

Processing steps:

1. Track quote state (single, double, backtick) while scanning the line
2. Detect command boundaries (after `|`, `;`, `&&`, `||`)
3. Determine default completion type (`TW_COMMAND` or `TW_ZERO`)
4. Call `tw_complete()` to look up user-defined completions
5. If found, override the default type with the user-specified list type

### `tw_complete()` — User Completion Lookup

```c
int tw_complete(const Char *line, Char **word, Char **pat, int looking, eChar *suf)
```

Parses the completion specification string and determines the appropriate list type. Handles all word specification types (`p`, `c`, `C`, `n`, `N`) and evaluates them left-to-right.

### `tw_result()` — Action Resolution

Maps the list type character to an internal completion type constant and generates candidates. For the backtick list type:

1. Sets `COMMAND_LINE` via `tsetenv(STRCOMMAND_LINE, line)`
2. Executes the command via `globone()` for command substitution
3. Captures output as a word list
4. Returns `TW_WORDLIST`
5. Cleans up with `Unsetenv(STRCOMMAND_LINE)`

### Completion Storage

User-defined completions are stored in a global `varent` tree structure (`completions`). Each entry is stored as a string in the format:

```
"command pattern action suffix"
```

The `docomplete()` function (the `complete` builtin implementation) adds entries using `set1()`, and `tw_find()` searches the tree using `Gmatch()` for glob pattern matching.

## Completion Display

### How Candidates Are Shown

When multiple matches exist:

1. **Unique prefix insertion** — if all candidates share a common prefix, tcsh inserts it
2. **Autolist** — if `autolist` is set, remaining choices are listed automatically
3. **Manual listing** — `^D` (at end of line) or `M-^D` (anywhere) lists choices
4. **Columnar display** — `print_by_column()` sorts and displays matches in columns

### The `ls-F` Builtin and Completion Listing

Completion listing uses the same display logic as `ls-F`:

- `/` suffix for directories
- `@` suffix for symbolic links
- `*` suffix for executables
- When the `color` variable is set, `ls-F` passes `--color=auto` to `ls`

### Completion Cycling

| Command | Default Binding | Description |
|---------|-----------------|-------------|
| `complete-word-fwd` | Not bound | Cycle to next completion candidate |
| `complete-word-back` | Not bound | Cycle to previous completion candidate |

These replace the current word with the next/previous candidate from the match list, allowing manual cycling through possibilities.

### Raw Completion (Bypassing User Definitions)

| Command | Default Binding | Description |
|---------|-----------------|-------------|
| `complete-word-raw` | `^X-Tab` | Like `complete-word` but ignores user-defined completions |
| `list-choices-raw` | `^X-^D` | Like `list-choices` but ignores user-defined completions |

These fall back to the built-in completion logic (command, variable, or filename) without consulting the `complete` builtin registrations.

## wordchars vs COMP_WORDBREAKS

tcsh does **not** have a `COMP_WORDBREAKS` variable like bash. The equivalent concept is `wordchars`:

### `wordchars` Shell Variable

Defines non-alphanumeric characters considered part of a word by editor commands (`forward-word`, `backward-word`, etc.):

| Mode | Default |
|------|---------|
| Emacs (vimode unset) | `*?_-.[]~=` |
| Vi (vimode set) | `_` |

### Key Differences from Bash's COMP_WORDBREAKS

| Aspect | Bash COMP_WORDBREAKS | Tcsh wordchars |
|--------|---------------------|-----------------|
| Purpose | Characters that **break** words for completion | Characters that are **part of** words for editing |
| Scope | Affects what `$2` contains in completion functions | Affects editor word movement commands |
| Completion impact | Only last segment after wordbreak is passed to function | Does not directly affect completion candidate generation |
| Default value | `" \t\n\"'@><=;|&(:` | `*?_-.[]~=` (emacs) or `_` (vi) |

The `wordchars` variable affects **editor** word boundaries, not **shell** word boundaries. The shell itself recognizes only whitespace and metacharacters as word separators. This means completion candidates are generated based on shell word splitting, not `wordchars`.

## External Completion Generator Integration

### The Standard Pattern

External completion generators (like carapace, git-completion.tcsh, etc.) hook into tcsh using the backtick list type:

```tcsh
complete mycmd 'p/*/`my_external_completer`/'
```

### How It Works

1. User presses TAB on a `mycmd` command
2. tcsh finds the `complete` registration for `mycmd`
3. tcsh sets `COMMAND_LINE` to the full command line
4. tcsh executes `my_external_completer` via command substitution
5. The completer reads `COMMAND_LINE` and prints candidates (one per line)
6. tcsh parses the output into words and presents them as completions

### Output Format

The external command should print completion candidates to stdout, one per line:

```
candidate1
candidate2
candidate3
```

Output is split at blanks, tabs, and newlines. Null words are discarded.

### Description Support

tcsh's native completion system does **not** support descriptions alongside candidates. Unlike bash (which can show descriptions with `COMP_TYPE=63`), fish (tab-separated), or zsh (`_describe`), tcsh candidates are plain words only.

However, some implementations work around this by embedding descriptions in the candidate value using the `_(description)` suffix format, which tcsh displays but the user must delete. The carapace tcsh formatter uses this approach:

```
value_(description)
```

### Limitations for External Generators

| Limitation | Impact |
|-----------|--------|
| No cursor position | Cannot distinguish between completing at end of word vs. middle |
| No word array | Must parse COMMAND_LINE manually to find argument positions |
| No completion type | Cannot distinguish between TAB press, listing, or menu completion |
| No description support | Candidates are plain words; descriptions require workarounds |
| No style/color support | Cannot style individual candidates |
| No nospace per-candidate | `addsuffix` is global; per-candidate suffix requires the suffix field in `complete` |
| No dynamic reloading | Completions defined in `complete` are static; dynamic ones require backtick commands |

### Git's tcsh Completion (Example)

The git project ships `git-completion.tcsh` which demonstrates the external generator pattern:

```tcsh
complete git 'p/*/`__git_tcsh_complete`/'
```

The `__git_tcsh_complete` function:

1. Reads `COMMAND_LINE`
2. Invokes the bash completion script in a subshell
3. Captures the `COMPREPLY` output
4. Prints each candidate on a separate line

This bridges bash's richer completion API to tcsh's simpler model.

## Common Completion Patterns

### Directory-Only Completion

```tcsh
complete cd 'p/1/d/'
complete pushd 'p/1/d/'
complete rmdir 'p/*/d/'
```

### Option Completion with Argument Types

```tcsh
complete find \
  'n/-name/f/' \
  'n/-newer/f/' \
  'n/-user/u/' \
  'n/-group/g/' \
  'n/-fstype/(nfs 4.2)/' \
  'n/-type/(b c d f l p s)/' \
  'c/-/(name newer user group fstype type)/' \
  'p/*/d/'
```

### Variable-Based Completion

```tcsh
set hosts = (localhost server1 server2)
complete ssh 'p/1/$hosts/'
complete scp 'c/*@/$hosts/' 'p/*/f/'
```

### Command Output Completion

```tcsh
complete kill 'p/*/`ps -e -o pid=`/'
complete man 'p/*/c/'
```

### File Type Filtering

```tcsh
# Only C source files
complete gcc 'p/*/f:*.[cChSs]/'

# Exclude object and archive files
complete vim 'p/*/f:^*.{o,a,so}/'
```

## References

- tcsh(1) man page — the definitive reference for the `complete` builtin
- `tw.parse.c` in the tcsh source — `tenematch()` implementation
- `tw.comp.c` in the tcsh source — `tw_complete()`, `tw_result()`, `docomplete()`
- `complete.tcsh` in the tcsh source — example completions (~1277 lines)
- [Complete::Tcsh CPAN module](https://metacpan.org/pod/Complete::Tcsh) — Perl API for tcsh completions
- [git-completion.tcsh](https://github.com/git/git/blob/master/contrib/completion/git-completion.tcsh) — git's tcsh completion bridge

## Related Skills

- [references/editor.md](references/editor.md) — editor commands for completion (complete-word, list-choices, etc.)
- [references/quoting-expansion.md](references/quoting-expansion.md) — how quoting affects completion candidates
- [references/startup-config.md](references/startup-config.md) — where to install completions (~/.tcshrc)
- carapace-dev skill → `references/shell.md` — cross-shell comparison including tcsh
