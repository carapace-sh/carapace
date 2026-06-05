# Nushell Completion System

In-depth reference for nushell's completion system — the NuCompleter, external completers, custom completions, spans, matching algorithms, and how external tools like carapace hook into the completion menu.

## The Completion Flow

When the user presses TAB (or another completion key), nushell:

1. **Adds a placeholder** — appends `'a'` to the input to ensure valid AST generation when the cursor is at an incomplete position
2. **Parses the input** — constructs the AST from the modified input
3. **Identifies the expression** — `find_pipeline_element_by_position()` traverses the AST to locate the innermost expression containing the cursor
4. **Dispatches to a completer** — based on the AST node type and context
5. **Filters and sorts** — `NuMatcher` applies the configured algorithm (prefix/substring/fuzzy) and case sensitivity
6. **Strips the placeholder** — adjusts spans to remove the placeholder character
7. **Returns suggestions** — `Vec<SemanticSuggestion>` to Reedline for display in the completion menu

### Expression Identification

`find_pipeline_element_by_position()` handles:

- **Nested blocks** — subexpressions, closures, row conditions
- **Pipeline elements** — call arguments, external call arguments
- **Cell paths** — variable access with member paths (e.g., `$record.key`)
- **Binary operators** — for operator completion
- **Attribute blocks** — for `@` syntax
- **Redirection targets** — checked via `check_redirection_in_block()`

### Dispatch Order

Nushell tries completers in this priority order:

1. **CustomCompletion** — if the parameter has a `@completer` annotation
2. **CommandWideCompletion** — if a command-wide completer is defined
3. **FlagCompletion** — if completing flags (after `-` or `--`)
4. **Built-in completers** — CommandCompletion, FileCompletion, DirectoryCompletion, VariableCompletion, etc.
5. **External completer** — the closure in `$env.config.completions.external.completer` (only if no other completer produced results)

## The External Completer Closure

External completers are the primary mechanism for tools like carapace to hook into nushell's completion system.

### Configuration

```nu
$env.config.completions.external = {
    enable: true
    max_results: 100
    completer: $my_completer
}
```

### Closure Signature

```nu
let my_completer = {|spans|
    # $spans is a list<string> of command-line parts
    # spans.0 = command name
    # spans.1.. = arguments
    # trailing empty string = cursor position
}
```

### What the Closure Receives

The `spans` parameter is a `list<string>` of the command-line arguments as parsed by nushell:

| User types | Spans received |
|-----------|----------------|
| `git <TAB>` | `[git ""]` |
| `git checkout <TAB>` | `[git checkout ""]` |
| `git checkout --<TAB>` | `[git checkout "--"]` |
| `git checkout -b <TAB>` | `[git checkout -b ""]` |
| `cargo +<TAB>` | `[cargo "+"]` |

Key properties:

- **Quoting is preserved** — nushell passes spans with quoting intact (e.g., `"act` becomes `'"act'` or `'"act'`). External completers must handle quote stripping.
- **Empty string for cursor** — a trailing empty string or space represents the cursor position
- **No COMP_* variables** — unlike bash, nushell does not set special environment variables for the completer
- **No COMP_WORDBREAKS** — nushell's parser doesn't split on `:` or `=`, so no wordbreak prefix handling is needed

### What the Closure Returns

The closure can return:

| Return type | Effect |
|------------|--------|
| `list<string>` | Simple completion values |
| `list<record>` | Records with `value`, optional `description`, optional `style` |
| `null` | Fall back to nushell's built-in file completion |
| `[]` (empty list) | Suppress completions entirely |

### Completion Record Fields

```nu
[
    {
        value: "main",           # Text inserted into the command line
        display: "main",         # Text shown in the menu (overrides value if different)
        description: "Default branch",  # Help text shown alongside
        style: { fg: green, bg: "#66078c", attr: ub }  # Styling
    }
]
```

| Field | Type | Required | Purpose |
|-------|------|----------|---------|
| `value` | string | Yes | Text inserted into the command line |
| `display` | string | No | Text shown in the completion menu (defaults to `value`) |
| `description` | string | No | Help text shown alongside the suggestion |
| `style` | string or record | No | Styling — color name, hex code, or `{fg, bg, attr}` record |
| `span` | record | No | `{start: int, end: int}` byte range to replace (nushell computes if omitted) |
| `append_whitespace` | bool | No | Whether to auto-append trailing space after insertion (defaults to `false`) |

### Carapace Integration

```nu
let carapace_completer = {|spans|
    carapace $spans.0 nushell ...$spans | from json
}
```

How it works:

1. `carapace $spans.0` — invokes carapace with the command name
2. `nushell` — tells carapace to use the nushell JSON formatter
3. `...$spans` — spread operator passes each span as a separate argument
4. `| from json` — nushell pipeline: parses JSON output into nushell records

Carapace encodes nospace as a trailing space in the `value` field (since `append_whitespace` defaults to `false`).

### Fish Integration

```nu
let fish_completer = {|spans|
    fish --command $"complete '--do-complete=($spans | str replace --all "'" "\\'" | str join ' ')'"
    | from tsv --flexible --noheaders --no-infer
    | rename value description
    | update value {|row|
      let value = $row.value
      let need_quote = ['\' ',' '[' ']' '(' ')' ' ' '\t' "'" '"' "`"] | any {$in in $value}
      if ($need_quote and ($value | path exists)) {
        let expanded_path = if ($value starts-with ~) {$value | path expand --no-symlink} else {$value}
        $'"($expanded_path | str replace --all "\"" "\\\"")"'
      } else {$value}
    }
}
```

Key notes on fish integration:

- Fish returns TSV output (value tab description)
- `--flexible` handles missing descriptions
- `--noheaders` output has no header row
- `--no-infer` keeps all values as strings (preserves git hashes)
- Path escaping: fish uses POSIX escapes; nushell doesn't parse them, so paths need explicit double-quoting

### Multiple Completers

Combine multiple completers using `match` on `$spans.0`:

```nu
let multiple_completers = {|spans|
    match $spans.0 {
        nu => $fish_completer
        git => $fish_completer
        _ => $carapace_completer
    } | do $in $spans
}
```

### Alias Workaround

Nushell has a [bug (#8483)](https://github.com/nushell/nushell/issues/8483) where completions don't work for aliases. Workaround:

```nu
let expanded_alias = (scope aliases | where name == $spans.0 | get -i 0 | get -i expansion)
let spans = (if $expanded_alias != null {
    $spans | skip 1 | prepend ($expanded_alias | split row " " | take 1)
} else { $spans })
```

This replaces `$spans.0` with the first word of the alias expansion so the external completer sees the actual command.

### CARAPACE_LENIENT

Carapace returns `ERR unknown shorthand flag` for non-supported flags. Set this environment variable to prevent errors:

```nu
$env.CARAPACE_LENIENT = 1
```

Or inline in the completer:

```nu
let carapace_completer = {|spans|
    CARAPACE_LENIENT=1 carapace $spans.0 nushell ...$spans | from json
}
```

### Complete Example

```nu
let fish_completer = {|spans|
    fish --command $"complete '--do-complete=($spans | str replace --all "'" "\\'" | str join ' ')'"
    | from tsv --flexible --noheaders --no-infer
    | rename value description
}

let carapace_completer = {|spans: list<string>|
    CARAPACE_LENIENT=1 carapace $spans.0 nushell ...$spans | from json
}

let external_completer = {|spans|
    let expanded_alias = scope aliases
    | where name == $spans.0
    | get -o 0.expansion

    let spans = if $expanded_alias != null {
        $spans
        | skip 1
        | prepend ($expanded_alias | split row ' ' | take 1)
    } else {
        $spans
    }

    match $spans.0 {
        nu => $fish_completer
        git => $fish_completer
        _ => $carapace_completer
    } | do $in $spans
}

$env.config.completions.external = {
    enable: true
    max_results: 100
    completer: $external_completer
}
```

## Custom Completions (the `@` Syntax)

Custom completions attach completer functions to command parameters using the `@` operator:

```nu
def animals [] { ["cat", "dog", "eel"] }
def my-command [animal: string@animals] { print $animal }
```

`string@animals` tells nushell two things:
1. The `SyntaxShape` (`string`) — used for type-checking
2. The completer (`animals`) — a command that returns suggestions

### Context-Aware Completions

Completers receive the command-line context:

```nu
def animal-names [context: string] {
    match ($context | split words | last) {
        cat => ["Missy", "Phoebe"]
        dog => ["Lulu", "Enzo"]
        eel => ["Eww", "Slippy"]
    }
}

def my-command [
    animal: string@animals
    name: string@animal-names
] { print $"The ($animal) is named ($name)." }
```

The completer can also receive cursor position:

```nu
def completer [context: string, position: int] {}
```

### Completion Options Override

Return a record with `options` to override global settings:

```nu
def animals [] {
    {
        options: {
            case_sensitive: false,
            completion_algorithm: substring,
            sort: false,
        },
        completions: [cat, rat, bat]
    }
}
```

### Styled Completions

Return records with `value`, `description`, and `style`:

```nu
def my_commits [] {
    [
        { value: "5c2464", description: "Add .gitignore", style: red },
        { value: "f3a377", description: "Initial commit", style: { fg: green, bg: "#66078c", attr: ub } }
    ]
}
```

### Suppressing Completions

- Return `null` to fall back to nushell's built-in file completion
- Return `[]` (empty list) to suppress completions entirely

## Built-in Completers

| Completer | Purpose | Source File |
|-----------|---------|-------------|
| `CommandCompletion` | Internal commands and external executables from PATH | `command_completions.rs` |
| `FileCompletion` | File and directory paths (handles `~`, n-dots `...`, quote handling, hidden files) | `file_completions.rs` |
| `DirectoryCompletion` | Directories only (used for `cd`) | `directory_completions.rs` |
| `VariableCompletion` | Variable names prefixed with `$` (special vars, working set, permanent state) | `variable_completions.rs` |
| `FlagCompletion` | Command flags (short and long forms) from command signature | `flag_completions.rs` |
| `CellPathCompletion` | Record keys or list columns with dot notation | `cell_path_completions.rs` |
| `DotNuCompletion` | Module names and `.nu` scripts for `use`/`source-env`/`overlay use` | `dotnu_completions.rs` |
| `ExportableCompletion` | Exportable items from modules | `exportable_completions.rs` |
| `AttributeCompletion` | Attribute names for `@` syntax | `attribute_completions.rs` |
| `OperatorCompletion` | Binary operators in expressions | `operator_completions.rs` |
| `CustomCompletion` | User-defined completion functions via `@` syntax | `custom_completions.rs` |
| `CommandWideCompletion` | Completion function for all arguments of a command | `custom_completions.rs` |

### FileCompletion Details

- **Path expansion** — `~` is expanded to home directory
- **N-dots support** — `...` expands to `../..`, `....` to `../../..`
- **Quote handling** — properly handles quoted paths
- **Hidden file handling** — non-hidden files appear first, then hidden files
- **Directory trailing slash** — directories get a trailing `/`

### CommandCompletion Details

- Searches internal command registry first
- Then searches PATH directories for external executables
- Uses hash table for caching (similar to bash's `hash` builtin)

## NuMatcher: Filtering and Sorting

`NuMatcher` performs filtering and sorting of completion candidates.

### Match Algorithms

| Algorithm | Description |
|-----------|-------------|
| `prefix` | Only matches if haystack starts with needle |
| `substring` | Matches if needle appears anywhere in haystack |
| `fuzzy` | Matches if all needle characters appear in order (uses `nucleo_matcher` for scoring) |

### Configuration

```nu
$env.config.completions = {
    algorithm: "prefix"    # or "substring" or "fuzzy"
    sort: "alphabetical"  # or "smart"
    case_sensitive: false
}
```

### Per-Completer Override

Custom completers can override global settings by returning a record with `options`:

```nu
{
    options: {
        case_sensitive: true,
        completion_algorithm: fuzzy,
        sort: false,
    },
    completions: [...]
}
```

## SemanticSuggestion and Suggestion

### SemanticSuggestion (nushell side)

Wraps `reedline::Suggestion` with additional nushell-specific data:

| Field | Type | Purpose |
|-------|------|---------|
| `value` | String | Completion text |
| `description` | Option<String> | Help text |
| `span` | Span | Byte range to replace |
| `style` | Option<Style> | Styling for the suggestion |
| `match_indices` | Option<Vec<usize>> | Character positions that matched the needle (for highlighting) |

### Suggestion (Reedline side)

| Field | Type | Purpose |
|-------|------|---------|
| `value` | String | Text to insert into the buffer |
| `display_override` | Option<String> | Text shown in menu instead of `value` |
| `description` | Option<String> | Help text |
| `style` | Option<Style> | `nu_ansi_term::Style` for rendering |
| `span` | Span | Byte range `{start, end}` to replace |
| `append_whitespace` | bool | Whether to add trailing space after insertion |
| `match_indices` | Option<Vec<usize>> | Indices of matching characters for fuzzy highlighting |

## Completion Dispatch Flow (End to End)

```
User presses TAB
  → Reedline calls NuCompleter::complete(line, pos)
  → NuCompleter adds placeholder 'a' and parses input
  → find_pipeline_element_by_position() identifies target expression
  → Dispatches to appropriate completer:
    ├── CustomCompletion (if @completer attached)
    ├── CommandWideCompletion (if defined)
    ├── FlagCompletion (if completing flags)
    ├── Built-in completers (Command, File, Directory, Variable, etc.)
    └── External completer closure (if no other results)
  → NuMatcher filters and sorts results
  → strip_placeholder_if_any() adjusts spans
  → Returns Vec<SemanticSuggestion> to Reedline
  → Reedline displays in completion menu
  → User selects → replace_in_buffer() applies Suggestion
```

## Edge Cases and Known Issues

1. **Alias completion bug** — [nushell#8483](https://github.com/nushell/nushell/issues/8483): completions don't work for aliases. Use the alias workaround above.

2. **Open quoting** — nushell passes spans with quoting intact. External completers must strip quotes (carapace does this via `nushell.Patch()`).

3. **No COMP_WORDBREAKS** — nushell's parser doesn't split on `:` or `=`, so no wordbreak prefix handling is needed (unlike bash).

4. **No redirect handling** — unlike bash, nushell does not leak redirect operators into completion arguments.

5. **`append_whitespace` defaults to `false`** — external completers must encode trailing space in the `value` field if they want a space after insertion.

6. **`max_results` limits output** — the `max_results` setting in `$env.config.completions.external` caps the number of suggestions returned by the external completer.

7. **Fallback to file completion** — returning `null` from the external completer triggers nushell's built-in file completion. Returning `[]` suppresses all completions.

8. **`from json` is required** — carapace outputs JSON; nushell needs `| from json` to parse it into records. Without this, the raw JSON string is returned as a single completion.

## References

- [Nushell Book: Custom Completions](https://www.nushell.sh/book/custom_completions.html)
- [Nushell Cookbook: External Completers](https://www.nushell.sh/cookbook/external_completers.html)
- [Nushell Book: Externs](https://www.nushell.sh/book/externs.html)
- [DeepWiki: NuCompleter](https://deepwiki.com/nushell/nushell/6.3-completion-system)
- [Nushell GitHub: nu-cli/completions](https://github.com/nushell/nushell/tree/main/crates/nu-cli/src/completions)

## Related Skills

- [references/reedline.md](references/reedline.md) — Reedline menu system and how completions are displayed
- [references/externs.md](references/externs.md) — Defining extern commands with completions
- [references/types.md](references/types.md) — SyntaxShape and how types drive completions
- [references/quoting.md](references/quoting.md) — How quoting affects spans passed to completers
