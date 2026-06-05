# BLE Syntax Highlighting

In-depth reference for ble.sh's syntax highlighting system — the incremental parser, face configuration, gspec color specification, and how highlighting integrates with the completion system.

## How the Syntax Parser Works

ble.sh performs **full syntactic analysis** of Bash command lines — not simple regex-based highlighting. This enables correct handling of complex structures like nested command substitutions, multiple here documents, and case statements.

### Core Data Structures

| Variable | Type | Purpose |
|----------|------|---------|
| `_ble_syntax_text` | string | The command line text being parsed |
| `_ble_syntax_stat[i]` | array | Parse state before character `i` |
| `_ble_syntax_nest[i]` | array | Nesting information at index `i` |
| `_ble_syntax_tree[i]` | array | Completed parse nodes ending at position `i+1` |
| `_ble_syntax_attr[i]` | array | Display attributes for character `i` |

### Parse State Format

Each element of `_ble_syntax_stat[i]` represents the parse state **immediately before** character `i`:

```
"ctx wlen wtype nlen tclen tplen nparam lookahead"
```

| Field | Type | Description |
|-------|------|-------------|
| `ctx` | int | Current context (command, quote, parameter expansion) |
| `wlen` | int | Length of current word (negative if no active word) |
| `wtype` | string | Type of current word |
| `nlen` | int | Length of current nesting level (negative if none) |
| `tclen` | int | Offset to last child tree node |
| `tplen` | int | Offset to previous sibling tree node |
| `nparam` | string | Nest-specific parameters (or "none") |
| `lookahead` | int | Characters examined to determine state |

The `lookahead` field enables the parser to handle **multi-character constructs** like `<<EOF` (heredoc) or `$((expr))` (arithmetic expansion).

### Nesting Structure

The `_ble_syntax_nest[inest]` array tracks contexts that span multiple characters:

```
"ctx wlen wtype inest tclen tplen nparam ntype"
```

- First 6 fields: Restoration state when the nest is closed
- `nparam`: Parameters to restore when returning
- `ntype`: Type identifier (e.g., `"$("` for command substitution)

### Parse Tree Node Format

```
"(wtype wlen tclen tplen wattr)*"
```

| Field | Description |
|-------|-------------|
| `wtype` | Word/nest type: integer context for words, string identifier for nests |
| `wlen` | Length of the range (start = `i - wlen`) |
| `tclen` | Offset to child node ending position (negative if no child) |
| `tplen` | Offset to previous sibling ending position (negative if first child) |
| `wattr` | Attributes: `-` (uncomputed), `g` (gflags), `m...` (multi-part), `d` (deleted) |

### Incremental Parsing

The parser operates **incrementally** — it only re-parses the affected range after an edit:

1. Identify affected range using dirty range tracking
2. Invalidate parse state in affected region
3. Re-parse from first invalidated position
4. Stop when reaching unchanged state

This approach efficiently handles partial edits without re-parsing the entire line.

### Context System

The parser tracks different syntactic contexts through integer context values:

- **Command context** — Expecting commands, assignments, redirections
- **Quote contexts** — Single quotes, double quotes, `$'...'`, `$"..."`
- **Parameter expansion** — `${...}`, `$((...))`
- **Command substitution** — `$(...)`, `` `...` ``
- **Arithmetic** — `$((...))`, `((...))`
- **Heredoc** — Within `<<EOF` content
- **Conditional** — `[[ ... ]]`, `[ ... ]`
- **Case patterns** — Within `case ... in` patterns

### Known Limitations

The parser deliberately simplifies certain complex Bash behaviors for performance:

- **Brace expansion + tilde** — `a=~:{a,b}` — Bash disables tilde expansion with brace expansion, but ble.sh highlights the first `~` as active
- **Extended globs in braces** — `[[ @({a,b}) ]]` — ble.sh highlights as if brace expansion is active
- **Array subscript evaluation** — Bash 4.3 and earlier: subscripts may be evaluated in false branches
- **Pathname expansion edge cases** — `@()` parsing discrepancy between Bash's initial parse and pathname expansion parse

## Face System

The `ble-face` function controls colors and styling for syntax highlighting and UI elements.

### Syntax Highlighting Faces

| Face | Default | Purpose |
|------|---------|---------|
| `syntax_default` | `none` | Default text |
| `syntax_command` | `fg=brown` | Commands |
| `syntax_quoted` | `fg=green` | Quoted content |
| `syntax_quotation` | `fg=green,bold` | Quotation marks |
| `syntax_escape` | `fg=magenta` | Escape sequences `\?` |
| `syntax_expr` | `fg=33` | Arithmetic expressions |
| `syntax_error` | `bg=203,fg=231` | Syntax errors |
| `syntax_varname` | `fg=202` | Variable names |
| `syntax_delimiter` | `bold` | Delimiters `;`, `&`, pipes, redirections |
| `syntax_param_expansion` | `fg=133` | Parameter expansions |
| `syntax_history_expansion` | `bg=94,fg=231` | History expansions |
| `syntax_function_name` | `fg=99,bold` | Function definitions |
| `syntax_comment` | `fg=242` | Comments |
| `syntax_glob` | `fg=198,bold` | Glob pattern operators |
| `syntax_brace` | `fg=37,bold` | Brace expansions |
| `syntax_tilde` | `fg=63,bold` | Tilde expansions |
| `syntax_document` | `fg=100` | Here document contents |
| `syntax_document_begin` | `fg=100,bold` | Here document delimiters |

### Command Highlighting Faces

| Face | Default | Purpose |
|------|---------|---------|
| `command_builtin_dot` | `fg=red,bold` | Builtin `.` |
| `command_builtin` | `fg=red` | Other builtins |
| `command_alias` | `fg=teal` | Aliases |
| `command_function` | `fg=99` | Function names |
| `command_file` | `fg=green` | File commands |
| `command_keyword` | `fg=blue` | Keywords |
| `command_jobs` | `fg=red` | Job specs |
| `command_directory` | `fg=33,underline` | Directories |
| `command_suffix` | (v0.4+) | Suffix sabbrevs |
| `command_suffix_new` | (v0.4+) | New suffix sabbrevs |

### Filename Highlighting Faces

| Face | Purpose |
|------|---------|
| `filename_directory` | Directories |
| `filename_directory_sticky` | Sticky-bit directories |
| `filename_link` | Symbolic links |
| `filename_orphan` | Broken symlinks |
| `filename_setuid` | Setuid files |
| `filename_setgid` | Setgid files |
| `filename_executable` | Executable files |
| `filename_other` | Other files |
| `filename_socket` | Sockets |
| `filename_pipe` | Named pipes |
| `filename_character` | Character devices |
| `filename_block` | Block devices |
| `filename_url` | URL-like filenames |
| `filename_warning` | Warning filenames |
| `filename_ls_colors` | Additional attributes from `$LS_COLORS` |

### Variable Highlighting Faces (v0.4+)

| Face | Purpose |
|------|---------|
| `varname_unset` | Unset variables |
| `varname_new` | New (not yet set) variables |
| `varname_empty` | Empty variables |
| `varname_export` | Exported variables |
| `varname_readonly` | Read-only variables |
| `varname_array` | Array variables |
| `varname_hash` | Associative array variables |
| `varname_number` | Numeric variables |
| `varname_transform` | Transform variables |
| `varname_expr` | Expression variables |

### Completion Menu Faces

| Face | Default | Purpose |
|------|---------|---------|
| `auto_complete` | `fg=238,bg=254` | Auto-completion suggestion |
| `menu_complete_selected` | `reverse` | Selected menu item |
| `menu_complete_match` | `bold` | Matching prefix in items |
| `menu_filter_fixed` | `bold` | Fixed portion of filter text |
| `menu_filter_input` | `fg=16,bg=229` | Actively typed filter text |

### Editing Faces

| Face | Default | Purpose |
|------|---------|---------|
| `region` | `bg=60,fg=231` | Selected region |
| `region_target` | `bg=153,fg=black` | Target region (for swap) |
| `region_match` | `bg=55,fg=231` | Matching bracket |
| `disabled` | `fg=242` | Disabled text |
| `overwrite_mode` | `fg=black,bg=51` | Overwrite mode indicator |

## gspec — Graphics Specifier

The `gspec` describes styles as a comma-separated list of specifiers:

### Attribute Specifiers

| Specifier | Attribute |
|-----------|-----------|
| `bold` | Bold (high intensity in some terminals) |
| `underline` | Underline |
| `blink` | Blink |
| `invis` | Invisible characters |
| `reverse` | Reversed video |
| `strike` | Strike line |
| `italic` | Italic shape |
| `standout` | Emphasize (synonym for `bold,reverse`) |
| `fg=name` | Set foreground color |
| `bg=name` | Set background color |
| `none` | No attribute |

### Named Colors

Basic: `default`, `transparent`, `black`, `gray`, `brown`, `red`, `green`, `lime`, `olive`, `yellow`, `navy`, `blue`, `purple`, `magenta`, `teal`, `cyan`, `silver`, `white`, `orange`

### Color Specification Formats

| Format | Example | Description |
|--------|---------|-------------|
| `#RGB` | `#f00` | 1-digit hex RGB |
| `#RRGGBB` | `#ff0000` | 2-digit hex RGB |
| `rgb:R/G/B` | `rgb:255/0/0` | RGB with 0-255 or percentile |
| `cmy:C/M/Y` | `cmy:0/1/1` | CMY |
| `cmyk:C/M/Y/K` | `cmyk:0/1/1/0` | CMYK |
| `hsl:H/S%/L%` | `hsl:0/100/50` | HSL |
| `hsb:H/S%/B%` | `hsb:0/100/100` | HSB |
| `I` | `196` | 256 index color (0-255) |

### Type Prefixes for ble-face (v0.4+)

| Prefix | Description |
|--------|-------------|
| `gspec` | Graphics specifier (default) |
| `g` | Integer value |
| `ref` | Reference to another face |
| `copy` | Copy from another face |
| `sgrspec` | SGR parameters |
| `ansi` | ANSI escape sequences |

## ble-face Command

```bash
# Set face
ble-face FACEPAT=[TYPE:]SPEC
ble-face -s FACEPAT [TYPE:]SPEC

# Define new face (v0.4+)
ble-face FACE:=[TYPE:]SPEC
ble-face -d FACE [TYPE:]SPEC

# Show settings
ble-face [-u | --color[=WHEN]]... [FACEPAT...]

# Reset faces
ble-face -r [FACEPAT...]
```

The `@` wildcard can match multiple faces in patterns:

```bash
ble-face -s syntax_@ fg=white  # Set all syntax_* faces
```

### Referencing Other Faces

```bash
ble-face menu_desc_type=ref:syntax_delimiter  # Reference another face
```

## Highlighting Configuration

### Enable/Disable

```bash
bleopt highlight_syntax=1       # Enable syntax highlighting (default: 1)
bleopt highlight_filename=1     # Enable filename highlighting (default: 1)
bleopt highlight_variable=1     # Enable variable type highlighting (default: 1)
```

### Color Schemes

```bash
bleopt color_scheme=            # Preset color schemes
```

### Terminal Color Support

```bash
bleopt term_index_colors=auto  # Number of index colors (auto/256/88/0)
bleopt term_true_colors=       # 24-bit color support (semicolon/colon/empty)
```

## Integration with Completion

The syntax parser provides the foundation for the completion system:

1. **Context analysis** — The parse tree determines what kind of completion is needed (command, file, variable, argument)
2. **Quoting context** — `comps_flags` encodes the current quoting state for the quote-insert system
3. **Word boundaries** — Parse tree word boundaries determine `COMP1`, `COMP2`, `COMPS`, `COMPV`
4. **Vi text objects** — Word boundaries from parse tree support vim-style text objects

## References

- [ble.sh Graphics Manual](https://github.com/akinomyoga/ble.sh/wiki/Manual-%C2%A72-Graphics)
- [Syntax Parser (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/4.1-syntax-parser)
- [Core Architecture (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/2-core-architecture)
