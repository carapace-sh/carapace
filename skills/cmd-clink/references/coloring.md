# Clink Input Line Coloring and Classification

In-depth reference for clink's input line coloring system — how words are classified, the classifier API, color settings, and completion menu styling.

## Overview

When the `clink.colorize_input` setting is enabled (default), argmatchers automatically apply colors to the input text as they parse it. When disabled, the entire input line uses `color.input`.

## How Words are Colored

### Command Word Coloring (Priority Order)

The command word is colored based on type, in priority order:

1. **Commands with argmatcher** → `color.argmatcher`
2. **Built-in CMD commands** → `color.cmd`
3. **Doskey aliases** → `color.doskey`
4. **Recognized executable files** → `color.executable` (if set)
5. **Unrecognized command words** → `color.unrecognized` (if set)
6. **None of above** → `color.input`

### Other Input Text Coloring

| Element | Color Setting |
|---------|--------------|
| Command separators (`&`, `|`) | `color.cmdsep` |
| Redirection symbols (`<`, `>`, `>&`) | `color.cmdredir` |
| Redirected file names | `color.input` |
| Flags defined by argmatcher | `color.flag` |
| Arguments defined by argmatcher | `color.arg` |
| Unexpected text (past what argmatcher expects) | `color.unexpected` |
| Auto-suggestion text | `color.suggestion` |
| No argmatcher defined | `color.input` |

## Classification Codes

Used with `word_classifications:classifyword()`:

| Code | Classification | Color Setting |
|------|---------------|---------------|
| `"a"` | Argument (preset matches) | `color.arg` or `color.input` |
| `"c"` | Shell command (CMD) | `color.cmd` |
| `"d"` | Doskey alias | `color.doskey` |
| `"f"` | Flag (preset matches) | `color.flag` |
| `"x"` | Executable (exists, not command/alias) | `color.executable` |
| `"u"` | Unrecognized (not command, alias, or executable) | `color.unrecognized` |
| `"o"` | Other (filenames, etc.) | `color.input` |
| `"n"` | None (not recognized as part of expected syntax) | `color.unexpected` |
| `"m"` | Prefix for argmatcher (e.g., `"mc"` or `"md"`) | `color.argmatcher` or other code's color |

The `"m"` prefix combines with other codes: `"mc"` = argmatcher + CMD command, `"md"` = argmatcher + doskey.

## Classifier API

### Argmatcher Classifier

Set via `_argmatcher:setclassifier(func)`:

```lua
local function classify_handler(arg_index, word, word_index, line_state, classifications, user_data)
    if arg_index == 0 then
        classifications:classifyword(word_index, "f")  -- Flag
    elseif arg_index == 2 and word:sub(-1) == "\\" then
        classifications:classifyword(word_index, "n")  -- Unexpected
    end
end

clink.argmatcher("myapp")
    :addflags("--help", "--version")
    :addarg({ "init", "build" })
    :addarg({ clink.filematches })
    :setclassifier(classify_handler)
```

**Parameters:**
- `arg_index` — Argument index in argmatcher (0 = flag)
- `word` — Partial string for word under cursor
- `word_index` — Word index in line_state
- `line_state` — line_state object
- `classifications` — word_classifications object
- `user_data` — (v1.5.17+) table for parsing assistance

### Global Classifier

Set via `clink.classifier(priority)`:

```lua
local my_classifier = clink.classifier(50)

function my_classifier:classify(commands)
    -- commands: table of {line_state, classifications}
    -- Return true to stop further classifiers, false/nil to continue
    for _, cmd in ipairs(commands) do
        local line_state = cmd.line_state
        local classifications = cmd.classifications
        local line = line_state:getline()
        -- Apply custom coloring...
    end
end
```

Priority is a number; lower values are called before higher values.

## word_classifications API

### `:classifyword(word_index, word_class, [overwrite])`

Classifies a word for coloring:

```lua
classifications:classifyword(word_index, "a")  -- Argument
classifications:classifyword(word_index, "f")  -- Flag
classifications:classifyword(word_index, "n")  -- Unexpected
```

- `overwrite` defaults to `true` (overwrites existing classification)
- Set `overwrite=false` to only apply if not already classified

### `:applycolor(start, length, color, [overwrite])`

Applies an ANSI SGR escape code to characters in the input line:

```lua
classifications:applycolor(pos, length, "95")  -- Bright magenta
```

- `start`: position to begin (1-based)
- `length`: number of characters
- `color`: SGR parameters (e.g., `"7"` for reverse video)
- An input line can have up to 100 unique color strings
- `overwrite` defaults to `true`

## Color Settings

### Input Line Colors

| Setting | Purpose |
|---------|---------|
| `color.input` | Default for unrecognized text |
| `color.arg` | Arguments from argmatcher |
| `color.flag` | Flags from argmatcher |
| `color.cmd` | CMD commands |
| `color.doskey` | Doskey aliases |
| `color.argmatcher` | Commands with argmatcher (supersedes cmd/doskey/input) |
| `color.executable` | Executable files |
| `color.unrecognized` | Unrecognized command words |
| `color.unexpected` | Text not in expected syntax |
| `color.cmdsep` | Command separators (`&`, `|`) |
| `color.cmdredir` | Redirection symbols |
| `color.suggestion` | Auto-suggestion text |

### Completion and Display Colors

| Setting | Purpose |
|---------|---------|
| `color.readonly` | Readonly files/directories |
| `color.hidden` | Hidden files/directories |
| `color.filtered` | Filtered matches (display differs from match) |
| `color.common_match_prefix` | Common prefix in completions |
| `color.description` | Default match description |
| `color.arginfo` | Arguments for flags/arguments in completions |
| `color.selected_completion` | Selected completion |

### Popup Colors

| Setting | Purpose |
|---------|---------|
| `color.popup` | Popup window colors |
| `color.popup_desc` | Popup description column |
| `color.popup_border` | Popup border |
| `color.popup_header` | Popup header |
| `color.popup_footer` | Popup footer |
| `color.popup_select` | Selected popup item |
| `color.popup_selectdesc` | Selected popup description |

### Other Colors

| Setting | Purpose |
|---------|---------|
| `color.message` | Readline message area |
| `color.modmark` | Modified mark (for changed lines) |
| `color.histexpand` | History expansion highlighting |
| `color.comment_row` | Comment row in clink-select-complete |
| `color.prompt` | Prompt color (backward compatibility) |

## Color Syntax

```
[attributes] [foreground_color] [on [background_color]]
```

**Attributes:** `bold`, `underline`, `italic`, `reverse`

**Named colors:** `default`, `black`, `red`, `green`, `yellow`, `blue`, `cyan`, `magenta`, `white`, `bright` (prefix for high intensity)

**RGB:** `#XXXXXX` (24-bit) or `#XXX` (12-bit)

**SGR:** `sgr N` (raw SGR escape code number)

**Examples:**

```
bold yellow on blue
#ff8030
sgr 7
bri yel
bold underline green on #222
default on blue
```

## Completion Match Coloring

### Using `match.coloring_rules` (Recommended)

```cmd
clink set match.coloring_rules di=93:ro ex=1;32:ex=1:ro=32:*.tmp=90
```

**Type codes:**

| Code | Meaning | Default Color |
|------|---------|---------------|
| `di` | Directory | bright blue (`01;34`) |
| `ex` | Executable | bright green (`01;32`) |
| `fi` | Normal file | normal color |
| `ro` | Readonly | from `color.readonly` |
| `hi` | Hidden | from `color.hidden` |
| `mi` | Missing | — |
| `ln` | Symlink | `ln=target` for target color |
| `or` | Orphaned symlink | — |
| `no` | Normal color | — |
| `any` | Matches any type | — |

**Glob patterns** can be appended: `*.tmp=90` colors `.tmp` files gray.

### Using `%LS_COLORS%` (Backward Compatibility)

```cmd
set LS_COLORS=so=90:fi=97:di=93:ex=92:*.pdf=30;105
```

**Type codes:** `di`, `ex`, `fi`, `ln`, `mi`, `no`, `or`, `so`

## Color Themes

### .clinktheme Files

Save and share color configurations:

```cmd
clink config theme save theme_name     -- Save current colors as theme
clink config theme use theme_name      -- Apply a theme
clink config theme list                 -- List available themes
clink config theme show theme_name     -- Preview a theme
```

**Search locations:**
1. Directories in `%CLINK_THEMES_DIR%`
2. `themes\` subdirectory under each scripts directory

### Override

`%CLINK_COLORTHEME%` environment variable overrides the active theme.

## Completion Menu Styling

### Match Display Types

| Type | Behavior |
|------|----------|
| `"word"` | Shows whole word including slashes |
| `"arg"` | Avoids appending space after colon/equal sign |
| `"cmd"` | Uses `color.cmd` |
| `"alias"` | Uses `color.doskey` |
| `"file"` | Last path component with file coloring |
| `"dir"` | Last path component with directory coloring |

**Type modifiers** (appended with comma): `"hidden"`, `"readonly"`, `"link"`, `"orphaned"`

### Description Display

- `F1` toggles between descriptions at bottom vs inline with matches
- When more than 9 matches: descriptions shown in single line at bottom
- Inline descriptions shown when they don't require more than 9 rows
- `color.description` sets default description color
- `color.arginfo` sets argument info color

### Common Prefix Highlighting

- `colored-completion-prefix` setting controls coloring the typed prefix
- `color.common_match_prefix` sets the color
- The `"so"` color from `%LS_COLORS%` is used (default bright magenta)

### Filtered Match Display

- Matches with `type="none"` or where `display` differs from `match` use `color.filtered`
- `display` field is shown instead of `match` when inserting
- `description` field shows aligned descriptions next to matches
