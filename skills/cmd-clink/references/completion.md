# Clink Completion System

In-depth reference for clink's completion pipeline — how matches are generated, filtered, displayed, and how external tools hook into the process.

## The Completion Flow

When the user presses TAB (or another completion key), clink:

1. Identifies the command word and searches for a matching argmatcher
2. If an argmatcher is found, it parses the input line to determine the current argument position
3. Calls match generators in increasing priority order
4. Applies match filters (`clink.onfiltermatches`)
5. Applies display filters (`clink.ondisplaymatches`)
6. Passes matches to Readline for display/insertion

### Argmatcher Resolution

When a command word is typed, clink searches for an argmatcher:

1. **Exact command name** — e.g., `git`
2. **Fully qualified path** (v1.3.38+) — e.g., `C:\Program Files\Git\git.exe`
3. **Lazy-loaded completion scripts** — searches `completions\` directories for matching `.lua` file
4. **No argmatcher** — falls back to match generators and file completion

### Match Generation Order

Within the completion pipeline:

1. **Argmatcher parsing** — if an argmatcher exists, it determines the current position and generates positional/flag matches
2. **Match generators** — called in increasing priority order (lower numbers first)
3. **File completion** — if no argmatcher or generator claims the line, file matches are generated

## Match Generators

Match generators are Lua functions called during Readline's completion process. They enable custom parsing for the input line or provide completions for the first word.

### Creating a Generator

```lua
local my_generator = clink.generator(priority)
function my_generator:generate(line_state, match_builder)
    -- Examine line_state and add matches via match_builder
    return true  -- Stop other generators
    -- return false  -- Let other generators continue
end
```

**Priority values:**
- Lower numbers execute before higher numbers
- Built-in generators use specific ranges (e.g., envvar_generator uses priority 10)
- Typical custom generators use priorities 10–90

### The `:getwordbreakinfo()` Function (Optional)

A generator can influence word breaking for the end word:

```lua
function my_generator:getwordbreakinfo(line_state)
    -- Returns word break length and optional end word length
    -- nil or 0: End word is truncated to 0 length (normal behavior)
    -- Two numbers: The end word is split at the word break length
    return 3, 1  -- First 3 chars in first word, remaining 1 char in end word
end
```

**Example from envvar_generator:**
- For `abc%USER` → returns `3,1` so words become `"abc"` and `"%"`
- For `abc%FOO%def` → returns `8,0` so words become `"abc%FOO%"` and `""`

> **Note:** `:getwordbreakinfo()` is called very often and must be very fast.

## match_builder API

The `match_builder` object is passed to generator functions and is used to submit completion matches.

### Methods

| Method | Description |
|--------|-------------|
| `addmatch(match, [type])` | Add a single match |
| `addmatches(matches, [type])` | Add multiple matches (returns count, success) |
| `setsuppressappend([state])` | Suppress appending anything after the match except a possible closing quote |
| `setsuppressquoting([state])` | Suppress quoting. 0 = normal, 1 = suppress quoting, 2 = suppress end quotes |
| `setappendcharacter([char])` | Set character to append after matches (e.g., `"="` for `set` command) |
| `setforcequoting()` | Force quoting rules even for non-filenames |
| `setfullyqualify([bool])` | Force completions to be inserted as fully qualified paths |
| `setnosort()` | Turn off sorting the matches |
| `setvolatile()` | Force matches to be used only once (regenerated each time) |
| `isempty()` | Return whether any matches have been added |

### Match Types

| Type | Description |
|------|-------------|
| `"word"` | Shows whole word including slashes |
| `"arg"` | Avoids appending space after colon/equal sign |
| `"cmd"` | Uses `color.cmd` |
| `"alias"` | Uses `color.doskey` |
| `"file"` | Last path component with file coloring |
| `"dir"` | Last path component with directory coloring |
| `"none"` | Uses `color.filtered` (display-only, not inserted) |

**Type modifiers** (appended with comma): `"hidden"`, `"readonly"`, `"link"`, `"orphaned"`

### Extended Match Table Format

```lua
{
    match           = "...",    -- The match text (inserted on selection)
    display         = "...",    -- Alternative display text (shown instead of match)
    arginfo         = "...",    -- Argument info string
    description     = "...",    -- Description for the match
    type            = "...",    -- Match type (see above)
    appendchar      = "...",    -- Character to append after match
    suppressappend  = true/false  -- Suppress appending after match
}
```

When `display` differs from `match`, the `color.filtered` setting is used for display.

## line_state API

The `line_state` object provides information about the current input line.

| Method | Returns | Description |
|--------|---------|-------------|
| `getline()` | string | The entire current line |
| `getword(index)` | string | Word at specified index (strips quotes) |
| `getendword()` | string | The last word (the word being completed) |
| `getwordcount()` | integer | Number of words in the current line |
| `getcursor()` | integer | Position of the cursor |
| `getwordinfo(index)` | table | Info about Nth word |
| `getcommandoffset()` | integer | Offset to start of the delimited command |
| `getcommandwordindex()` | integer | Index of the command word (usually 1) |
| `getendwordoffset()` | integer | Offset of the last word |
| `getrangeoffset()` | integer | Offset to start of the range |
| `getrangelength()` | integer | Length of the range |

### `getwordinfo()` Return Table

```lua
{
    offset   -- Offset where word starts in line string
    length   -- Length of word (includes embedded quotes)
    quoted   -- Boolean: word is quoted
    delim    -- Delimiter character or empty string
    cmd      -- true if built-in CMD command
    alias    -- true if doskey alias
    redir    -- true if redirection arg
}
```

## Match Filtering

### `clink.onfiltermatches(func)` (v1.1.41+)

Registers a function called after matches are generated but before they are displayed or inserted. This is reset every time match generation is invoked.

```lua
local function filter_matches(matches, completion_type, filename_completion_desired)
    -- matches: table of match objects
    -- completion_type: "?", "*", "\t", "!", "@", "%"
    -- filename_completion_desired: boolean
    local ret = {}
    for _, m in ipairs(matches) do
        if not m.match:find("%.exe$") then
            table.insert(ret, m)
        end
    end
    return ret
end
```

**Completion types:**

| Type | Meaning |
|------|---------|
| `?` | List alternatives |
| `*` | Insert all completions |
| `\t` | Standard completion (TAB) |
| `!` | List if ambiguous |
| `@` | List if no partial match |
| `%` | Menu complete |

**Important:** Filter functions can remove matches but cannot add new ones.

### `clink.ondisplaymatches(func)` (v1.1.12+)

Registers a function called before matches are displayed. Can modify `display` and `description` fields. Affects display but not insertion.

```lua
local function my_display_filter(matches, popup)
    -- matches: table of match objects (v1.3.1+ includes all fields)
    -- popup: boolean indicating popup display
    local new_matches = {}
    for _, m in ipairs(matches) do
        if m.type:find("^dir") then
            m.display = "\x1b[35m*" .. m.match  -- Magenta asterisk for dirs
        end
        table.insert(new_matches, m)
    end
    return new_matches
end
```

## Match Volatility

By default, clink caches and reuses match lists for performance. Use `match_builder:setvolatile()` to force regeneration each time. This is essential for arguments with special syntax (e.g., email addresses, URLs) where the match set depends on runtime state.

## How Carapace Hooks Into Clink

Carapace generates a Lua snippet that creates an argmatcher for the target command. The argmatcher's completion function:

1. Sets `CARAPACE_COMPLINE` environment variable to the current line state
2. Invokes the carapace binary with `_carapace cmd-clink <command>` arguments
3. Parses the tab-delimited output (value, display, description, appendchar)
4. Adds each match via `match_builder:addmatch()`

The snippet uses `match_builder:setnosort()` and `match_builder:setvolatile()` to let carapace handle sorting and to ensure fresh completions on each invocation.

See the carapace-dev skill → `references/shell.md` for the cmd-clink row in the cross-shell comparison table.

## Completion Display

### Standard Completion (TAB)

- If one match: inserts immediately
- If multiple matches: inserts common prefix, then lists on second TAB
- `completion-auto-query-items` controls whether to prompt before showing many matches

### clink-select-complete (Ctrl-Space)

Interactive completion list with:
- Arrow key navigation
- Typing to filter the list
- `F4` toggles between "find" and "filter" search modes
- `F1` toggles description display (inline vs bottom)
- `Ctrl-Home`/`Ctrl-End` select first/last match
- `match.max_rows` limits displayed rows
- `match.preview_rows` shows preview

### old-menu-complete

CMD-style completion that cycles through matches. Default TAB binding with enhanced key bindings. Not limited by `completion-query-items`.

### Popup Windows

- `clink-popup-history` (F7) — popup list of command history
- `clink-popup-directories` (Ctrl-Alt-PgUp) — popup directory list
- `clink-popup-show-help` (F1) — popup key binding help

Popup navigation: `Esc`/`Ctrl-G` cancel, `Enter` insert, `Shift-Enter`/`Ctrl-Enter` insert without executing, `Del` delete entry, `F3`/`Shift-F3` next/prev match.

## Sorting

- Matches are normally sorted in Unicode-aware locale order
- Disable with `nosort=true` in arg tables, `:addargunsorted()`, `:addflagsunsorted()`, or `match_builder:setnosort()`
- `match.sort_dirs` setting controls directory position: `first`, `mixed`, or `last`
