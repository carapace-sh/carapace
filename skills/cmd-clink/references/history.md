# Clink History System

In-depth reference for clink's command history — how it's stored, configured, navigated, and the Lua API for programmatic access.

## How History is Stored

Clink stores command history in a **master history file** (`clink_history`) in the profile directory.

- When `history.save` is enabled, history persists across sessions
- Every time a new input line starts, clink reloads the master history list and prunes it to `history.max_lines`
- For performance, deleting a history line marks the line as deleted without rewriting the file
- When deleted lines exceed the max lines or 200 (whichever is larger), the history file is compacted
- Force compaction with `history compact`

### Multiple History Files

`%CLINK_HISTORY_LABEL%` environment variable (up to 32 alphanumeric characters) enables multiple master history files. Useful for project-specific history.

## History Settings

| Setting | Default | Description |
|---------|---------|-------------|
| `history.save` | `true` | Persist history between sessions |
| `history.max_lines` | `10000` (enhanced) / `2500` | Max lines to save (0 = unlimited) |
| `history.shared` | `false` | Share history across all clink instances |
| `history.dont_add_to_history_cmds` | `history exit` | Commands not added to history (space-separated) |
| `history.dupe_mode` | — | How duplicate entries are handled |
| `history.ignore_space` | — | Ignore lines prefixed with whitespace |
| `history.expand_mode` | — | Control history expansion in quotes |
| `history.auto_expand` | `true` | Auto-expand history on Enter |
| `history.show_preview` | — | Show preview of history expansions |
| `history.sticky_search` | `false` | Reusing history line doesn't add to end |
| `history.time_stamp` | `off` | Save/show timestamps (`off`, `save`, `show`) |
| `history.time_format` | `%F %T` | Timestamp format (strftime specifiers) |

### Shared History

When `history.shared` is `true`:
- All clink instances update the master history file
- They reload it every time a new input line starts
- Commands entered in one instance appear in other instances immediately

When `history.shared` is `false`:
- Each instance loads the master file but doesn't append its own history until exit
- Each instance's history is isolated during its session

### Controlling What Gets Added to History

| Prefix | Doskey Expansion | History |
|--------|-----------------|--------|
| (none) | Yes | Yes |
| Space | No | No |
| Semicolon | No | Yes |

## History Commands

| Command | Description |
|---------|-------------|
| `clink history` / `history` | List command history |
| `history --help` | Show usage information |
| `history compact` | Force history file compaction |
| `history compact <n>` | Prune to N items and compact |
| `history compact --unique` | Remove duplicate entries |

## Key Bindings for History Navigation

| Key | Command | Description |
|-----|---------|-------------|
| `Up` | `previous-history` | Previous history entry |
| `Down` | `next-history` | Next history entry |
| `PgUp` | `history-search-backward` | Search for matching prefix (backward) |
| `PgDn` | `history-search-forward` | Search for matching prefix (forward) |
| `Ctrl-R` | `reverse-search-history` | Incremental reverse search |
| `Ctrl-S` | `forward-search-history` | Incremental forward search |
| `F7` / `Ctrl-Alt-Up` | `clink-popup-history` | Popup list of history |
| `Alt-Ctrl-K` | `add-history` | Add current line without executing |
| `Alt-Ctrl-D` | `remove-history` | Delete selected history entry |
| `Ctrl-O` | `operate-and-get-next` | Execute and move to next entry |

## History Deduplication

- `history.dupe_mode` controls duplicate handling
- `erase_prev` removes previous entry when duplicate is added
- `history compact --unique` removes duplicate entries from the history file

## History Expansion

Clink uses Readline's History library for history expansion (using `!` character).

### Event Designators

| Designator | Meaning |
|------------|---------|
| `!n` | Command line n |
| `!-n` | Current command minus n |
| `!!` | Previous command |
| `!string` | Most recent command starting with string |
| `!?string?` | Most recent command containing string |
| `!#` | Current command line typed so far |

### Word Designators

| Designator | Meaning |
|------------|---------|
| `:$` | Last word |
| `:^` | First word |
| `:*` | All words |
| `:n` | Word n (0-based) |
| `:n-m` | Words n through m |

### Modifiers

| Modifier | Meaning |
|----------|---------|
| `:s/old/new/` | Substitute first occurrence |
| `:gs/old/new/` | Substitute all occurrences |
| `:h` | Remove trailing pathname component |
| `:t` | Remove leading pathname components |
| `:r` | Remove suffix |
| `:p` | Print without executing |

### Expansion Commands

| Key | Command | Description |
|-----|---------|-------------|
| `Alt-^` | `clink-expand-history` | Expand history |
| `Alt-Ctrl-E` | `clink-expand-line` | Expand the current line |

### Disabling Expansion

```cmd
clink set history.auto_expand off     -- Don't auto-expand on Enter
clink set history.expand_mode off     -- Turn off all expansion syntax
```

## History Lua API

### `rl.gethistorycount()` (v1.3.18+)

Returns the number of history items.

### `rl.gethistoryitems(start, end)` (v1.3.18+)

Returns a table of history items. Each item contains:

```lua
local h = rl.gethistoryitems(1, rl.gethistorycount())
-- h.line       [string] The item's command line string
-- h.time       [integer or nil] The item's time (os.time() compatible)
```

### `clink.onhistory(func)` (v1.5.13+)

Registers a callback when input line is accepted and about to be added to history.

```lua
clink.onhistory(function(input_text)
    -- Return false to cancel adding to history
    if input_text:match("^%s*$") then
        return false
    end
end)
```

## Popup History

Press `F7` or `Ctrl-Alt-Up` to show a popup list of command history.

**Popup navigation:**

| Key | Action |
|-----|--------|
| `Esc` / `Ctrl-G` | Cancel |
| `Enter` | Execute highlighted item |
| `Shift-Enter` / `Ctrl-Enter` | Insert without executing |
| `Del` | Delete selected entry |
| `F3` / `Shift-F3` | Next/previous match |
| `Ctrl-Home` / `End` | Scroll to extremes |
| `Left` / `Right` | Scroll horizontally |

**Settings:**

| Setting | Description |
|---------|-------------|
| `clink.popup_delete_direction` | Direction after deletion (`last` or `next`) |
| `clink.popup_search_mode` | Default search mode (`find` or `filter`) |
