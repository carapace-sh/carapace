# Clink Auto-Suggestions

In-depth reference for clink's auto-suggestion system — inline suggestions, suggestion lists, the suggester API, and how to create custom suggesters.

## Overview

Clink can suggest command lines as you type, based on command history and completions. It offers two display modes:

1. **Suggestion List** — An interactive list of several suggestions (press `F2`)
2. **Inline Suggestions** — A single suggestion in a muted color at the end of the input line

## Settings

| Setting | Default | Description |
|---------|---------|-------------|
| `autosuggest.enable` | `true` | Enables/disables automatic suggestions |
| `autosuggest.inline` | `false` | Shows suggestions in muted color after cursor |
| `autosuggest.hint` | `true` | Shows right-aligned usage hint (e.g., `F2=List Suggestions`) |
| `autosuggest.strategy` | `match_prev_cmd history completion` | Determines how suggestions are generated |
| `autosuggest.async` | `true` | Generates matches asynchronously for responsiveness |
| `autosuggest.original_case` | `true` | Preserves original capitalization when inserting |

## Suggestion List

Press `F2` to toggle the suggestion list. When active:

| Key | Action |
|-----|--------|
| `F2` | Toggle suggestion list on/off |
| `Esc` | Revert to original input and clear list |
| `Up`/`Down` | Navigate selection and replace input line |
| `F7` or `Ctrl-Alt-Up` | Show popup list of command history |

Related settings:

| Setting | Default | Description |
|---------|---------|-------------|
| `suggestionlist.default` | `false` | Whether sessions start with the list on |
| `suggestionlist.autooff` | `true` | Automatically turns off after use |

## Inline Suggestions

Disabled by default. Enable with `clink set autosuggest.inline true`.

**Key bindings for accepting suggestions:**

| Key | Action |
|-----|--------|
| `Right` or `End` | Insert whole suggestion |
| `Ctrl-Right` | Insert next word |
| `Shift-Right` | Insert next full word (up to space) |
| `F2` | Show suggestion list instead |

Inline suggestions appear in `color.suggestion` (muted color by default).

## Suggester API

Create custom suggestion generators using `clink.suggester()`:

### Basic Suggester

```lua
local suggestor = clink.suggester("name_goes_here")

function suggestor:suggest(line_state, matches)
    -- line_state: contains the input line
    -- matches: contains possible completions
    -- Return nil to pass to next suggester
    -- Return a string (or empty string) to use as suggestion
    -- Return "suggestion", offset for custom insertion point
end
```

### Example: Longest Common Prefix Suggester

```lua
local prefix_suggestor = clink.suggester("completion_prefix")

function prefix_suggestor:suggest(line_state, matches)
    if not line_state:getline():match("[^ ]") then
        return
    end
    local prefix = matches:getprefix()
    if prefix == "" then
        return
    end
    local info = line_state:getwordinfo(line_state:getwordcount())
    return prefix, info.offset
end
```

### Example: Doskey Macro Suggester

```lua
local doskeyarg = clink.suggester("doskeyarg")
function doskeyarg:suggest(line, matches)
   if line:getword(1) == "doskey" and
           line:getline():match("[ \t][^ \t/][^ \t]+=") and
           not line:getline():match("%$%*") then
       if line:getline():sub(#line:getline()) == " " then
           return "$*"
       else
           return " $*"
       end
   end
end
```

### Multiple Suggestions (v1.8.0+)

```lua
local suggestor = clink.suggester("multi_suggest")

function suggestor:suggest(line_state, matches, limit)
    if not limit then
        -- Single suggestion mode
        return "single suggestion", 5
    else
        -- Multiple suggestions mode (for suggestion list)
        local suggestions = {}
        table.insert(suggestions, {
            "suggestion text",      -- suggestion string
            5,                      -- offset where suggestion begins
            highlight = { 5, 8 },   -- optional: highlight offsets
            tooltip = "description" -- optional: description shown below
        })
        return suggestions
    end
end
```

**Suggestion table fields:**

| Field | Type | Description |
|-------|------|-------------|
| `[1]` | string | Suggestion text |
| `[2]` | integer | Offset where suggestion begins |
| `highlight` | {start, end} | Optional: highlight range offsets |
| `tooltip` | string | Optional: description shown below |

## Built-in Suggestion Strategies

The `autosuggest.strategy` setting determines suggestion order. Strategies are space-separated:

| Strategy | Description |
|----------|-------------|
| `match_prev_cmd` | Matches the previous command context |
| `history` | Most recent matching command from history |
| `completion` | First matching completion |

Custom suggesters can be added to this list by their registered name:

```cmd
clink set autosuggest.strategy "match_prev_cmd history completion my_custom_suggester"
```

## rl_buffer API for Suggestions

| Method | Description |
|--------|-------------|
| `rl_buffer:hassuggestion()` | Check if suggestion is available (v1.6.4+) |
| `rl_buffer:insertsuggestion([amount])` | Insert suggestion (v1.6.4+) |

**`amount` parameter:**

| Value | Meaning |
|-------|---------|
| `"all"` (default) | Insert entire suggestion |
| `"word"` | Insert next word |
| `"fullword"` | Insert next full word (up to space) |

## Usage Hints

When `autosuggest.hint` is enabled, clink shows right-aligned hints:

| Hint | Meaning |
|------|---------|
| `F2=List Suggestions` | Press F2 for suggestion list |
| `Right=Insert F2=List` | Press Right to accept, F2 for list |

Turn off with `clink set autosuggest.hint false`.

## How Suggestions Relate to Completions

Suggestions and completions are distinct but related:

| Aspect | Completions | Suggestions |
|--------|------------|-------------|
| Trigger | TAB or Ctrl-Space | Automatic (as you type) |
| Scope | Current word | Whole command line |
| Source | Argmatchers, generators | History, completions, custom suggesters |
| Display | Completion list/popup | Inline text or suggestion list |
| Key | TAB | F2 (list), Right (accept inline) |

Custom suggesters can use the `matches` parameter to access completion results and derive suggestions from them.
