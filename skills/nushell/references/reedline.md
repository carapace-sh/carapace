# Nushell Reedline Integration

In-depth reference for Reedline — nushell's line editor — the Completer trait, Suggestion struct, menu system, keybindings, and how completions flow from NuCompleter through Reedline to the user.

## How Reedline Works

Reedline is nushell's cross-platform line editor. It sits between the user and the terminal:

1. **Character input** — reads characters from the terminal, interprets key sequences, maps them to editing commands
2. **Buffer management** — text is stored in an internal `LineBuffer`; characters are inserted at the cursor position
3. **Completion** — when the user triggers completion, Reedline calls the configured `Completer`
4. **Menu display** — suggestions are shown in a menu; user navigates and selects
5. **Buffer update** — selected suggestion replaces the appropriate span in the buffer

## The Completer Trait

```rust
pub trait Completer {
    fn complete(&self, line: &str, pos: usize) -> Vec<Suggestion>;
    fn complete_with_base_ranges(&self, line: &str, pos: usize) -> (Vec<Suggestion>, Vec<Range<usize>>);
    fn partial_complete(&self, line: &str, pos: usize, start: usize, offset: usize) -> Vec<Suggestion>;
    fn total_completions(&self, line: &str, pos: usize) -> usize;
}
```

| Method | Required | Returns | Purpose |
|--------|----------|---------|---------|
| `complete` | Yes | `Vec<Suggestion>` | Primary method — generate suggestions from buffer state |
| `complete_with_base_ranges` | No | `(Vec<Suggestion>, Vec<Range<usize>>)` | Suggestions + deduplicated target spans |
| `partial_complete` | No | `Vec<Suggestion>` | Subset of completions for lazy loading/pagination |
| `total_completions` | No | `usize` | Count of available completions |

`NuCompleter` implements this trait, bridging nushell's engine state with Reedline's completion API.

## The Suggestion Struct

```rust
pub struct Suggestion {
    pub value: String,
    pub display_override: Option<String>,
    pub description: Option<String>,
    pub style: Option<Style>,
    pub span: Span,
    pub append_whitespace: bool,
    pub match_indices: Option<Vec<usize>>,
}
```

| Field | Type | Purpose |
|-------|------|---------|
| `value` | `String` | Text inserted into the buffer |
| `display_override` | `Option<String>` | If present, shown in menu instead of `value` |
| `description` | `Option<String>` | Help text for the suggestion |
| `style` | `Option<Style>` | `nu_ansi_term::Style` for specialized rendering |
| `span` | `Span` | Byte range `{start, end}` to replace in the buffer |
| `append_whitespace` | `bool` | If true, space is added after the suggestion is accepted |
| `match_indices` | `Option<Vec<usize>>` | Indices of matching characters (for fuzzy highlighting) |

### Span

```rust
pub struct Span {
    pub start: usize,
    pub end: usize,
}
```

Positions are in **bytes** (not characters) to align with Rust's `String` indexing. The span indicates which portion of the line buffer should be replaced by `value`.

### How `append_whitespace` Works

- **`true`** — Reedline adds a space after inserting the suggestion
- **`false`** (default) — no space is added; the external completer must encode trailing space in `value` if desired

Carapace uses `false` (default) and encodes nospace/space behavior directly in the `value` field.

## The Menu System

Reedline provides four standard menu implementations:

| Menu Type | Layout | Description |
|-----------|--------|-------------|
| `ColumnarMenu` | Grid | Multiple columns; supports 2D navigation; auto-adjusts based on text size |
| `ListMenu` | List | Vertical list; supports `!N` quick selection |
| `DescriptionMenu` | Description | Prioritizes extended descriptions and extra information |
| `IdeMenu` | Floating | IDE-style with borders and metadata |

### Menu Trait

```rust
pub trait Menu: Send {
    fn menu_string(&self, available_lines: usize, use_ansi_coloring: bool) -> String;
    fn menu_required_lines(&self, terminal_columns: usize) -> usize;
    fn min_rows(&self) -> u16;
    fn get_values(&self) -> Vec<Suggestion>;
    // ... update and event handling methods
}
```

### Menu Configuration in nushell

```nu
$env.config.menus ++= [{
    name: completion_menu
    only_buffer_difference: false
    marker: "| "
    type: {
        layout: columnar
        columns: 4
        col_width: 20
        col_padding: 2
    }
    style: {
        text: green
        selected_text: green_reverse
        description_text: yellow
    }
}]
```

### Menu Configuration Fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Unique menu name (referenced by keybindings) |
| `only_buffer_difference` | bool | `true` = search on text after activation; `false` = search on all text |
| `marker` | string | Indicator shown when menu is active (e.g., `"\| "`, `"? "`) |
| `type.layout` | string | Menu layout: `columnar`, `list`, or `description` |
| `type.columns` | int | Number of columns (for columnar/description) |
| `type.col_width` | int or null | Fixed column width (null for auto-calculation) |
| `type.col_padding` | int | Spaces between columns |
| `type.page_size` | int | Entries shown when activating menu (for list) |
| `type.selection_rows` | int | Rows for found options (for description) |
| `type.description_rows` | int | Rows for command description (for description) |
| `style.text` | style | Unselected suggestion style |
| `style.selected_text` | style | Selected suggestion style |
| `style.description_text` | style | Description text style |
| `source` | closure | Custom menu source (for user-defined menus) |

### ColumnarMenu

Default completion menu. Displays suggestions in a grid:

- **Columns** — number of columns (default: 4)
- **Traversal direction** — `horizontal` (left-to-right, top-to-bottom) or `vertical` (top-to-bottom, left-to-right)
- **Auto-column** — when suggestions have descriptions, switches to single-column mode
- **Scrolling** — supported when rows exceed available lines
- **Partial completion** — supports completing common prefixes

### ListMenu

Simple vertical list:

- **Page size** — number of entries per page
- **Quick selection** — type `!N` + Enter to select the Nth entry directly
- **Filtering** — type keywords to filter entries

### DescriptionMenu

Layout that prioritizes descriptions:

- **Selection rows** — rows for the selection area
- **Description rows** — rows for the description area
- **Extra field** — supports `extra: [string]` for additional information

### User-Defined Menus

Custom menus use the `source` field:

```nu
$env.config.menus ++= [{
    name: vars_menu
    only_buffer_difference: true
    marker: "# "
    type: { layout: list, page_size: 10 }
    style: {
        text: green
        selected_text: green_reverse
        description_text: yellow
    }
    source: { |buffer, position|
        scope variables
        | where name =~ $buffer
        | sort-by name
        | each { |row| {value: $row.name description: $row.type} }
    }
}]
```

The `source` closure receives:
- `$buffer` — captured text (based on `only_buffer_difference`)
- `$position` — cursor position

Must return records with:

```nu
{
    value: "text",           # Required — inserted into buffer
    description: "help",     # Optional — displayed with value
    span: { start: 0, end: 5 },  # Optional — section to replace
    extra: ["more info"]     # Optional — for description menu
}
```

## Menu Keybindings

### Default Keybindings (when menu is active)

| Key | Action |
|-----|--------|
| Tab | Select next item |
| Shift+Tab | Select previous item |
| Enter | Accept selection |
| ↑/↓ | Move menu up/down |
| ←/→ | Move menu left/right |
| Ctrl+P/N | Move menu up/down |
| Ctrl+B/F | Move menu left/right |
| Ctrl+Z | Previous page |
| Ctrl+X | Next page |

### Binding Completion to a Key

```nu
$env.config.keybindings ++= [{
    name: completion_menu
    modifier: control
    keycode: char_t
    mode: emacs
    event: { send: menu name: completion_menu }
}]
```

### Tab Cycling (menu-complete behavior)

Use `until` to make Tab cycle through completions:

```nu
$env.config.keybindings ++= [{
    name: completion_menu
    modifier: none
    keycode: tab
    mode: emacs
    event: {
        until: [
            { send: menu name: completion_menu }
            { send: menunext }
        ]
    }
}]
```

First press opens the menu; subsequent presses navigate to the next element.

## MenuEvent

Reedline's `MenuEvent` enum handles user interactions:

| Event | Description |
|-------|-------------|
| `Activate(bool)` | Activate menu (bool = edit mode) |
| `Deactivate` | Deactivate menu |
| `Edit(bool)` | Line buffer modification |
| `NextElement` | Select next element |
| `PreviousElement` | Select previous element |
| `MoveUp` | Move up in menu |
| `MoveDown` | Move down in menu |
| `MoveLeft` | Move left in menu |
| `MoveRight` | Move right in menu |
| `NextPage` | Go to next page |
| `PreviousPage` | Go to previous page |

## ReedlineMenu Wrapper

Reedline routes completion calls to the appropriate completer based on the menu variant:

| Variant | Completer Source |
|---------|------------------|
| `EngineCompleter` | Forwards the engine's main completer (NuCompleter) to the menu |
| `HistoryMenu` | Instantiates a `HistoryCompleter` using the engine's history |
| `WithCompleter` | Uses the specific completer provided during menu construction |

## Menu → Completer Integration Flow

```
User presses TAB (or bound key)
  → Reedline receives key event
  → MenuEvent::Activate sent to menu
  → Menu calls update_values()
    → Queries the completer (NuCompleter or custom)
    → Refreshes menu suggestions
  → update_working_details()
    → Recalculates display parameters
  → menu_string() generates ANSI-formatted output
  → Terminal displays menu
  → User navigates (MenuEvent::NextElement, etc.)
  → User selects (Enter)
  → replace_in_buffer()
    → Applies selected Suggestion to LineBuffer
    → Uses Suggestion.span to determine replacement range
    → Inserts Suggestion.value
    → If append_whitespace, adds trailing space
```

## Editing Modes

| Mode | Setting | Description |
|------|---------|-------------|
| Emacs | `$env.config.edit_mode = "emacs"` | Default; standard Emacs-style editing |
| Vi | `$env.config.edit_mode = "vi"` | Modal editor with Normal and Insert modes |

### Emacs Mode Keybindings

| Key | Action |
|-----|--------|
| Ctrl+A | Move to start of line |
| Ctrl+E | Move to end of line |
| Ctrl+F | Move forward one character / History-hint complete |
| Ctrl+B | Move backward one character |
| Alt+F | History-hint complete one word |
| Ctrl+L | Clear screen |
| Ctrl+R | Search history |
| Ctrl+C | Cancel current line |
| Ctrl+P | Move up |
| Ctrl+N | Move down |
| End | Move to end / Complete history hint |
| Home | Move to line start |

### Vi Mode

- Starts in **insertion mode**
- Press `Esc` to enter **command mode**
- Standard vi motions: `w`/`e`/`b` (word), `0`/`$` (line), `f`/`t` (find)
- `i`/`a` to re-enter insertion mode
- `d` (delete), `c` (change), `r` (replace), `x` (delete char), `u` (undo)

## Multi-line Editing

| Method | Description |
|--------|-------------|
| Enter with open bracket | Auto-continues when `{`, `(`, or `[` is open |
| Trailing pipe | Pressing Enter after `\|` continues the line |
| Alt+Enter / Shift+Enter | Manually insert a newline |
| Ctrl+O | Open in external editor; saving updates the line |

## DefaultCompleter

Reedline includes a trie-based `DefaultCompleter` for keyword-based completion:

- **Multi-line handling** — flattens multi-line input by replacing `\r\n` or `\n` with spaces
- **Word splitting** — splits the line by spaces and iterates backwards to find the word prefix under cursor
- **Minimum length** — ignores words shorter than `min_word_len` (default: 2)

## History Menu

- **Paging** — Ctrl+X appends more records; Ctrl+Z/Ctrl+X navigate pages
- **Filtering** — type keywords to filter history
- **Quick selection** — type `!N` + Enter to insert the Nth entry directly
- **Deduplication** — uses `HashSet` to filter identical command lines

## References

- [Nushell Book: Line Editor](https://www.nushell.sh/book/line_editor.html)
- [DeepWiki: Reedline Completion System](https://deepwiki.com/nushell/reedline/8-completion-system)
- [DeepWiki: Reedline Menu System](https://deepwiki.com/nushell/reedline/6.3-menu-system)
- [Reedline GitHub](https://github.com/nushell/reedline)

## Related Skills

- [references/completion.md](references/completion.md) — NuCompleter and the completion dispatch system
- [references/style.md](references/style.md) — How styles are applied in menus
