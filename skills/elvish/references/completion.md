# Elvish Completion System

In-depth reference for elvish's completion system — the arg-completer mechanism, complex-candidate API, matcher system, complete-getopt helper, filter DSL, and the internal completion pipeline.

## The Completion Flow

When the user presses TAB (or another completion key), elvish:

1. Parses the current input line into an AST
2. Identifies the completion context (command, argument, variable, index, redirection)
3. Extracts the **seed** — the text being completed (the partial word before the cursor)
4. Calls the appropriate completer to generate raw candidates
5. Applies the matcher to filter candidates against the seed
6. Sorts and deduplicates candidates
7. Applies the quoting pass (adds appropriate quoting for insertion)
8. Displays candidates in the completion UI (ComboBox widget)

### Completer Types (Tried in Order)

The completion engine tries these completers in sequence; the first one that produces results wins:

| Order | Completer | Context | Example |
|-------|-----------|---------|---------|
| 1 | `completeCommand` | Command position (start of line, after `;`, after pipe `|`) | `gi⇥` → `git` |
| 2 | `completeIndex` | Index position (`$var[⇥`) | `$edit:completion:[⇥` |
| 3 | `completeRedir` | Redirection target (`> ⇥`) | `> out⇥` |
| 4 | `completeVariable` | Variable name (`$⇥`) | `$edi⇥` → `$edit` |
| 5 | `completeArg` | Command argument | `git che⇥` → `checkout` |

Each completer returns a **context** struct containing:
- `name` — completion type string (`"command"`, `"argument"`, `"variable"`, `"index"`, `"redir"`)
- `seed` — the text being completed
- `quote` — quoting mode (determines how the candidate is quoted on insertion)
- `interval` — the byte range in the input to replace

## The `arg-completer` Mechanism

### Registration

```elvish
set edit:completion:arg-completer[command-name] = {|@args|
    # $args is [command-name arg1 arg2 ...]
    # Output completions via put, echo, or edit:complex-candidate
}
```

The `$edit:completion:arg-completer` variable is a **map** from command names to completer functions. When elvish completes an argument for command `$x`, it looks up `$edit:completion:arg-completer[$x]` and calls it.

### How the Completer Is Called

When the user types `man 1⇥`:

```elvish
$edit:completion:arg-completer[man] man 1
```

When the user types `man 1 ⇥` (cursor after space, starting a new argument):

```elvish
$edit:completion:arg-completer[man] man 1 ""
```

The trailing empty string signals that a new argument is being started. The completer receives **all existing arguments** including the command name as the first element.

### Fallback Behavior

If no arg-completer is registered for a command, elvish falls back to **filename completion** (`GenerateFileNames`). This is the default for any command not in `$edit:completion:arg-completer`.

If the arg-completer is explicitly set to `$nil`, elvish also falls back to filename completion.

### Output Methods for Candidates

A completer can output candidates in three ways:

**1. Byte output** — each line becomes a candidate:

```elvish
set edit:completion:arg-completer[mycmd] = {|@args|
    echo candidate1
    echo candidate2
}
```

Limitation: candidates cannot contain newlines.

**2. Value output** — each value becomes a candidate:

```elvish
set edit:completion:arg-completer[mycmd] = {|@args|
    put candidate1 candidate2
}
```

**3. Complex candidates** — full control over display, suffix, and styling:

```elvish
set edit:completion:arg-completer[mycmd] = {|@args|
    edit:complex-candidate checkout &display=(styled checkout blue) &code-suffix=' '
    edit:complex-candidate --verbose &display=(styled --verbose yellow) &code-suffix=''
}
```

### Example: Simple Command Completer

```elvish
set edit:completion:arg-completer[apt] = {|@args|
    var n = (count $args)
    if (== $n 2) {
        put install uninstall update upgrade search show list
    } elif (== $n 3) {
        put $@all-packages
    }
}
```

### Example: Completer with Subcommand Dispatch

```elvish
use re

fn all-git-branches {
    git branch -a --format="%(refname:strip=2)" | re:awk {|0 1 @rest| put $1 }
}

var common-git-commands = [
  add branch checkout clone commit diff init log merge
  pull push rebase reset revert show stash status
]

set edit:completion:arg-completer[git] = {|@args|
    var n = (count $args)
    if (== $n 2) {
        put $@common-git-commands
    } elif (>= $n 3) {
        all-git-branches
    }
}
```

### How External Tools Hook In

External completion frameworks (like carapace) register themselves as arg-completers that invoke an external command and parse its output:

```elvish
set edit:completion:arg-completer[myapp] = {|@arg|
    myapp _carapace elvish (all $arg) | from-json | each {|completion|
        # Process JSON output and emit edit:complex-candidate calls
    }
}
```

Key integration points:
- `(all $arg)` — spreads the arg list into individual arguments
- `| from-json` — parses JSON output into elvish values
- `edit:complex-candidate` — provides display/suffix/styling control
- `edit:notify` — displays error messages above the editor

## `edit:complex-candidate`

```elvish
edit:complex-candidate $stem &display='' &code-suffix=''
```

| Parameter | Type | Purpose |
|-----------|------|---------|
| `$stem` | positional string | The actual value to insert into the command line |
| `&display` | string or styled text | How the candidate appears in the completion UI. If empty, `$stem` is used. Accepts output of `styled` builtin. |
| `&code-suffix` | string | Suffix appended to the quoted stem on insertion. **Not quoted** — added verbatim. Empty = no suffix. Space `' '` = trailing space after insertion. |

### How Insertion Works

1. When a candidate is accepted, elvish inserts a **quoted version** of `$stem`
2. If `&code-suffix` is non-empty, it is **appended verbatim** (not quoted)
3. This gives the completer precise control over what appears after the inserted value

### CodeSuffix and Nospace

The `&code-suffix` parameter is elvish's equivalent of the "nospace" concept in other shells:

| `&code-suffix` | Effect | Use Case |
|-----------------|--------|----------|
| `' '` (space) | Normal insertion with trailing space | Regular arguments, subcommands |
| `''` (empty) | No trailing space after insertion | Directories, `--flag=`, partial completions |

### Display with Styled Text

The `&display` parameter accepts either a plain string or styled text (output of the `styled` builtin):

```elvish
# Plain display
edit:complex-candidate checkout &display='checkout (Switch branches)'

# Styled display
edit:complex-candidate checkout &display=(styled checkout blue)(styled ' (Switch branches)' dim)
```

When `&display` differs from `$stem`, the completion UI shows the display text but inserts the stem text.

### Examples from Elvish's Own File Completion

```elvish
# Directory — no trailing space (the / is part of stem)
edit:complex-candidate d/ &code-suffix='' &display=(styled-segment d/ &fg-color=blue &bold)

# Regular file — trailing space for next argument
edit:complex-candidate bar &code-suffix=' ' &display=[^styled bar]
```

### Deprecated Options

- **`&display-suffix`** — deprecated in elvish 0.14.0. Use `&display` instead.
- **`&style`** — removed from `edit:complex-candidate`. Styling is now done via `styled` text in `&display`. Fixed in elvish 0.18.0 ([issue #1011](https://github.com/elves/elvish/issues/1011)).

## The Matcher System

After the completer outputs candidates, elvish matches them against the user's typed text (the **seed**). The matcher table `$edit:completion:matcher` maps completion types to matcher functions.

### Matcher Lookup

```elvish
set edit:completion:matcher[argument] = {|seed| edit:match-prefix $seed &ignore-case=$true }
```

Elvish indexes `$edit:completion:matcher` with the completion type name. Completion types:

| Type | Context |
|------|---------|
| `variable` | Variable name completion (`$⇥`) |
| `index` | Index completion (`$var[⇥`) |
| `command` | Command name completion |
| `redir` | Redirection target completion |
| `argument` | Command argument completion |

If the matcher map lacks a key for the current type, `$edit:completion:matcher['']` (empty string key) is used as the default.

### Matcher Protocol

A matcher is a function that:
1. Receives the **seed** as its first argument
2. Reads candidate strings from its **value input**
3. Outputs an identical number of **booleans** — one per candidate — indicating whether to keep it

```elvish
# Custom matcher: keep candidates that contain the seed as substring
set edit:completion:matcher[argument] = {|seed|
    each {|cand| has-prefix $cand $seed }
}
```

### Built-in Matchers

| Matcher | Description |
|---------|-------------|
| `edit:match-prefix` | Matches candidates that start with the seed |
| `edit:match-substr` | Matches candidates that contain the seed as a substring |
| `edit:match-subseq` | Matches candidates that contain the seed characters as a subsequence |

All three accept these options:

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `&ignore-case` | boolean | `$false` | Case-insensitive matching |
| `&smart-case` | boolean | `$false` | Case-insensitive unless seed contains uppercase |

### Matcher Configuration Examples

```elvish
# Case-insensitive prefix matching for arguments
set edit:completion:matcher[argument] = {|seed| edit:match-prefix $seed &ignore-case=$true }

# Substring matching for all types
set edit:completion:matcher[''] = {|seed| edit:match-substr $seed }

# Smart-case subsequence matching for commands
set edit:completion:matcher[command] = {|seed| edit:match-subseq $seed &smart-case=$true }
```

Default value: `[&''=$edit:match-prefix~]` (prefix matching for all completion types).

## `edit:complete-getopt`

A helper function for completing commands that follow GNU-style option conventions.

```elvish
edit:complete-getopt $args $opt-specs $arg-handlers
```

### Parameters

**1. `$args`** — array of current command-line arguments (without the command name itself).

**2. `$opt-specs`** — array of maps, each defining one command-line option:

| Key | Type | Description |
|-----|------|-------------|
| `short` | string | One-letter short option (without the dash), e.g., `"a"` for `-a` |
| `long` | string | Long option name (without `--`), e.g., `"all"` for `--all` |
| `arg-optional` | boolean | `$true` if option receives an optional argument |
| `arg-required` | boolean | `$true` if option receives a mandatory argument |
| `desc` | string | Human-readable description for the completion menu |
| `completer` | function | Function to generate completions for the option's argument |

Only one of `arg-optional` or `arg-required` can be `$true`.

**3. `$arg-handlers`** — array of functions for positional argument completion. Each function receives the last element of `$args`. If the last element is `...`, the last handler is reused for all following positional arguments.

### Usage Example

```elvish
fn complete {|@args|
    var opt-specs = [
        [&short=a &long=all &desc="Show all"]
        [&short=n &desc="Set name" &arg-required=$true
         &completer= {|_| put name1 name2 }]
    ]
    var arg-handlers = [
        {|_| put first1 first2 }
        {|_| put second1 second2 }
        ...
    ]
    edit:complete-getopt $args $opt-specs $arg-handlers
}
```

### Behavior Examples

```
~> complete ''
▶ first1
▶ first2

~> complete '-'
▶ (edit:complex-candidate -a &display='-a (Show all)')
▶ (edit:complex-candidate --all &display='--all (Show all)')
▶ (edit:complex-candidate -n &display='-n (Set name)')

~> complete -n ''
▶ name1
▶ name2

~> complete -a ''
▶ first1
▶ first2

~> complete arg1 ''
▶ second1
▶ second2

~> complete arg1 arg2 ''
▶ second1
▶ second2
```

## `edit:complete-filename` and `edit:complete-dirname`

### `edit:complete-filename`

```elvish
edit:complete-filename $args...
```

Produces filename completions for the last argument:
- Determines the directory to complete from the path
- If basename starts with `.`, completes hidden files; otherwise non-hidden files
- Returns `edit:complex-candidate` objects with styles from `$E:LSCOLORS`
- Directories have trailing `/`; non-directory files have space as code suffix

This is the **default handler** for commands without explicit arg-completers.

### `edit:complete-dirname`

```elvish
edit:complete-dirname $args...
```

Like `edit:complete-filename` but only produces directory completions.

## The Filter DSL

Completion mode (and other listing modes) support filtering via a mini DSL:

| Expression | Meaning |
|------------|---------|
| `literal` | Matches items containing the literal string (case-insensitive if all lower-case) |
| `[re $pattern]` | Matches items matching the regular expression |
| `[and $expr...]` | Matches items matching all expressions |
| `[or $expr...]` | Matches items matching any expression |

Multiple filter expressions are ANDed together.

Examples:
- `git` — matches items containing "git" (case-insensitive)
- `[re ^git]` — matches items starting with "git"
- `[and git [re ^git]]` — matches items containing "git" AND starting with "git"

## The Internal Completion Pipeline

### Source Code Architecture

The completion system spans multiple packages:

| Package | Files | Purpose |
|---------|-------|---------|
| `pkg/edit/completion.go` | Editor-level completion API | `completionStart`, `complexCandidate`, `initCompletion`, arg-completer adapter |
| `pkg/edit/complete/` | Core engine | `Complete()`, `Config`, `Result`, `RawItem`, `ComplexItem`, `PlainItem`, completers, generators, filterers |
| `pkg/cli/modes/completion.go` | UI widget | `CompletionSpec`, `CompletionItem`, ComboBox-based display |

### The `Complete()` Function

```go
func Complete(code CodeBuffer, ev *eval.Evaler, cfg Config) (*Result, error)
```

Pipeline:
1. Parse code into AST
2. Walk the AST to find the path to the cursor position
3. Try each completer in order (command → index → redir → variable → argument)
4. First successful completer produces raw items
5. Filter using `FilterPrefix` (default filterer)
6. Sort alphabetically
7. Cook items (add quoting for insertion)
8. Deduplicate

### RawItem Types

```go
type RawItem interface {
    String() string
    Cook(parse.PrimaryType) modes.CompletionItem
}

type PlainItem string  // Simple string candidate

type ComplexItem struct {  // Rich candidate
    Stem       string   // Used in code and menu
    CodeSuffix string   // Appended to code (e.g., space for executables)
    Display    ui.Text  // Display text (if empty, defaults to Stem)
}
```

### Config Struct

```go
type Config struct {
    Filterer     Filterer      // Filters raw candidates. Default: FilterPrefix
    ArgGenerator ArgGenerator  // Generates candidates for command arguments. Default: GenerateFileNames
}
```

### Arg Generator Protocol

```go
type ArgGenerator func(args []string) ([]RawItem, error)
```

The `args` slice contains the command name as `args[0]` followed by all existing arguments. The generator returns a slice of `RawItem` values.

### How arg-completers Are Called Internally

The `adaptArgGeneratorMap` function in `pkg/edit/completion.go` bridges the elvish-level arg-completer map to the Go-level `ArgGenerator` interface:

1. Looks up the completer function from `$edit:completion:arg-completer[command-name]`
2. If not found or nil, falls back to `GenerateFileNames`
3. Calls the completer function via the `Evaler` with all arguments
4. Captures value output through a pipe
5. Converts output values to `RawItem`:
   - `complexItem` → `complete.ComplexItem`
   - `string` → `complete.PlainItem`
6. Returns the collected items

### The Quoting Pass (Cook)

After raw items are generated and filtered, each item is "cooked" by calling `Cook(quoteType)`:

- `PlainItem` — the string is quoted according to the detected quote type (bareword, single-quoted, double-quoted)
- `ComplexItem` — the `Stem` is quoted, and `CodeSuffix` is appended verbatim after the quoted stem

The `Display` field of `ComplexItem` is used as-is for the UI display (it's already `ui.Text`).

### Completion UI Widget

The completion mode uses a `ComboBox` widget (from `pkg/cli/tk`):

- **CodeArea** — shows the filter text with a "COMPLETING" prompt
- **ListBox** — displays candidates horizontally with styling
- **OnSelect** — when a candidate is selected, stores the `ToInsert` text as pending code
- **OnAccept** — when a candidate is accepted (Enter), applies the pending code and closes the mode
- **OnFilter** — when the filter text changes, re-filters candidates using the filter DSL predicate

The `CompletionItem` struct:

```go
type CompletionItem struct {
    ToShow   ui.Text  // Used in the UI and for filtering
    ToInsert string   // Used when inserting a candidate
}
```

### Smart Completion

`edit:completion:smart-start` (the default Tab binding) adds autofix handling before starting completion:

1. Applies any pending autofixes (e.g., missing `use` statements)
2. Calls `complete.Complete()` to get candidates
3. If all candidates share a common prefix longer than the seed, inserts the common prefix
4. Otherwise, opens the completion UI

`edit:smart-enter` (the default Enter binding) also applies autofixes before accepting the line.

## Completion Mode Keybindings

The completion mode has its own binding table: `$edit:completion:binding`.

Default bindings (from insert mode):

| Key | Action |
|-----|--------|
| `Tab` | `edit:completion:smart-start` (with autofix) or `edit:completion:start` |
| `Up` / `Down` | Navigate candidates |
| `Enter` | Accept selected candidate |
| `Ctrl-[` | Close completion mode |

Custom binding example:

```elvish
set edit:completion:binding[Ctrl-F] = { edit:close-mode }
```

## Comparison with Other Shells

| Feature | Elvish | Bash | Zsh | Fish |
|---------|--------|------|-----|------|
| Registration | `$edit:completion:arg-completer[cmd]` | `complete -F func cmd` | `compdef func cmd` | `complete -c cmd` |
| Candidate output | `put`, `echo`, `edit:complex-candidate` | `COMPREPLY` array | `compadd`, `_describe` | `complete` command |
| Display vs insert | `&display` vs `$stem` | `COMPREPLY` only | `-d` display array | `--description` |
| Nospace control | `&code-suffix=''` | `compopt +o nospace` | `compadd -S ''` | `--no-space` |
| Per-candidate styling | `styled` in `&display` | No (Readline colors) | `zstyle` colors | `--color` |
| Matcher customization | `$edit:completion:matcher` | No (prefix only) | `zstyle matcher-list` | No (prefix only) |
| Error messages | `edit:notify` | `printf >&2` | `_message` | No standard way |
| Word breaking | No word breaks (clean args) | `COMP_WORDBREAKS` | `COMPWORDS` | Tokenized args |
| Open-quote problem | No (parser provides clean args) | Yes (manual handling) | Yes (COMPQUOTED) | Yes (commandline) |

## Edge Cases and Known Issues

### Single-Element List Unwrapping

Elvish's `put` command unwraps single-element lists. When iterating over a one-element array, use `all (one)` to re-wrap it:

```elvish
put $completion[Candidates] | all (one) | each {|c| ... }
```

Without `all (one)`, a single-element array would be unwrapped and `each` would receive the element directly instead of the list.

### Newlines in Candidates

Byte-output candidates (via `echo`) cannot contain newlines — each line becomes a separate candidate. Use `put` with string values or `edit:complex-candidate` for values containing newlines (though newlines in completion values are generally problematic).

### edit:notify Persistence

`edit:notify` messages persist in the notification area until the next editor action. Carapace suppresses usage hints when candidates exist to avoid notification spam (see `meta.Usage = ""` when `len(values) > 0`).

### Matcher Input Protocol

Custom matchers must read the same number of candidates from input as they output booleans. A mismatch causes a runtime error.

### Parallel Candidate Emission

Carapace uses `peach` (parallel `each`) in its snippet for performance. This means candidates may be emitted out of order. Elvish's completion UI sorts candidates alphabetically after receiving them, so order is not preserved from the completer.

## References

- [Elvish Editor API Reference](https://elv.sh/ref/edit.html) — official documentation for the `edit:` module
- [Elvish Language Reference](https://elv.sh/ref/language.html) — language specification
- [Elvish Source: pkg/edit/completion.go](https://github.com/elves/elvish/blob/master/pkg/edit/completion.go) — editor-level completion API
- [Elvish Source: pkg/edit/complete/](https://github.com/elves/elvish/tree/master/pkg/edit/complete) — core completion engine
- [Elvish Source: pkg/cli/modes/completion.go](https://github.com/elves/elvish/blob/master/pkg/cli/modes/completion.go) — completion UI widget

## Related Skills

- For carapace-specific elvish integration (snippet, value formatting, JSON output, ParseStyling validation), see the **carapace-dev** skill → `references/shell-elvish.md`.
- For elvish editor API details (keybindings, prompts, hooks, modes), see [references/editor.md](editor.md).
- For elvish styling system details (styled, styled-segment, colors), see [references/styling.md](styling.md).
