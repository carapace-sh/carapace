# Clink Lua Scripting System

In-depth reference for clink's Lua scripting system — how scripts are loaded, script directories, lazy loading, event handlers, and the Lua API.

## Overview

Clink uses **Lua 5.2** for extensibility. Lua scripts can create completions, custom prompts, color input, provide suggestions, and respond to events.

## Script Loading

### Script Directories (searched in order)

1. All directories listed in the `clink.path` setting (separated by semicolons)
2. If `clink.path` is not set: the DLL directory and the profile directory
3. All directories listed in `%CLINK_PATH%` environment variable (separated by semicolons)
4. All directories registered by `clink installscripts`

Run `clink info` to see the script paths for the current session.

### Loading Behavior

- Lua scripts are loaded **once** at startup
- Scripts are only reloaded if forced (e.g., `clink-reload` bound to `Ctrl-X Ctrl-R`)
- Loading a script executes it immediately
- Code not inside a function runs when the script loads (don't use for external programs or print messages)
- Use `clink.onbeginedit()` to defer initialization code

### Tilde Expansion

Clink performs tilde expansion on Lua script directory names: `~\` becomes `%HOME%` or `%HOMEDRIVE%%HOMEPATH%` or `%USERPROFILE%`.

## Completion Directories (Lazy Loading)

Lua completion scripts can be placed in a `completions\` directory to enable lazy loading:

- Prevents scripts from loading at clink startup
- Scripts load only when the corresponding command is typed
- Makes clink load faster with large quantities of argmatcher scripts
- When a command name is typed, clink searches for a matching `.lua` file in completions directories

### Completion Directory Search Order

1. Any directories in `%CLINK_COMPLETIONS_DIR%` environment variable
2. A `completions\` subdirectory under each scripts directory

**Important:** Only put scripts in completions directories if they specifically say they can go there. Scripts with additional functionality (like prompt filters) won't work correctly in completions directories until the corresponding command is typed.

### Installing Completion Scripts

```cmd
clink installscripts <path>
```

Registers a directory containing completion scripts. The path is stored and searched on every clink session.

## Event Handlers

### `clink.onbeginedit(func)` (v1.1.11+)

Called when clink's edit prompt is activated, before prompt filters run. Useful for deferred initialization:

```lua
local initialized = false
clink.onbeginedit(function()
    if not initialized then
        initialized = true
        -- One-time setup
    end
    -- Per-prompt setup
end)
```

### `clink.onendedit(func)` (v1.1.20+)

Called when clink's edit prompt ends. Receives input text as argument.

```lua
clink.onendedit(function(input_text)
    -- Cleanup after command is entered
end)
```

### `clink.onfilterinput(func)` (v1.2.16+)

Called after `onendedit`. Can replace the edit prompt text by returning a string. Can return multiple lines in a table (v1.3.13+).

```lua
clink.onfilterinput(function(text)
    -- Return nil to use text as-is
    -- Return string to replace text
    -- Return table of strings for multiple lines
end)
```

### `clink.oncommand(func)` (v1.3.12+)

Called when the command word changes in the edit line. Receives `line_state` and a table with `command`, `quoted`, `type`, and `file` fields.

```lua
clink.oncommand(function(line_state, command_info)
    if command_info.command == "git" then
        -- React to git command being typed
    end
end)
```

### `clink.onaftercommand(func)` (v1.2.50+)

Called after every editing command (key binding).

### `clink.oninputlinechanged(func)` (v1.4.18+)

Called after an editing command makes changes in the input line. Receives the new line contents.

### `clink.onhistory(func)` (v1.5.13+)

Called when input is about to be added to history. Can return `false` to cancel.

### `clink.oninject(func)` (v1.1.21+)

Called when clink is injected into CMD. Called only once per session.

### `clink.onprovideline(func)` (v1.3.18+)

Called after `onbeginedit` but before input editing starts. If it returns a string, it's executed as a command line without showing a prompt.

## Core clink Functions

| Function | Description |
|----------|-------------|
| `clink.argmatcher([commands...])` | Create/get an argument matcher |
| `clink.generator(priority)` | Create a match generator |
| `clink.promptfilter(priority)` | Create a prompt filter |
| `clink.classifier(priority)` | Create an input line classifier |
| `clink.hinter(priority)` | Create an input hint provider |
| `clink.suggester(name)` | Create a suggestion generator |
| `clink.popuplist(title, items)` | Display popup list |
| `clink.dirmatches(word)` | Get directory matches |
| `clink.dirmatchesexact(word)` | Directory matches without `*` append |
| `clink.filematches(word)` | Get file matches |
| `clink.getprompts()` | Get available custom prompts |
| `clink.getthemes()` | Get available color themes |
| `clink.applytheme(name)` | Apply a color theme |
| `clink.parseline(text)` | Parse command line into commands |
| `clink.recognizecommand(word)` | Get word classification |
| `clink.reload()` | Reload Lua scripts |
| `clink.setcoroutineinterval(ms)` | Set coroutine timing |
| `clink.addpackagepath(path)` | Add to Lua package search path |
| `clink.getclinkprompt()` | Get active .clinkprompt name |
| `clink.getalias(name)` | Get doskey alias definition |
| `clink.info()` | Get clink session information |

## Match Filtering Functions

### `clink.onfiltermatches(func)` (v1.1.41+)

Registers a function called after matches are generated but before they are displayed or inserted. Reset every time match generation is invoked.

```lua
clink.onfiltermatches(function(matches, completion_type, filename_completion_desired)
    local ret = {}
    for _, m in ipairs(matches) do
        if not m.match:find("%.exe$") then
            table.insert(ret, m)
        end
    end
    return ret
end)
```

Can remove matches but cannot add new ones.

### `clink.ondisplaymatches(func)` (v1.1.12+)

Registers a function called before matches are displayed. Can modify `display` and `description` fields.

```lua
clink.ondisplaymatches(function(matches, popup)
    for _, m in ipairs(matches) do
        if m.type:find("^dir") then
            m.display = "\x1b[35m*" .. m.match
        end
    end
    return matches
end)
```

## Additional Lua API Groups

| API | Description |
|-----|-------------|
| `console` | Console screen functions (scroll, write, title) |
| `git` | Git repository information (branch, status) |
| `http` | HTTP requests |
| `io` | Extended file I/O (popenyield, etc.) |
| `log` | Debug logging |
| `matches` | Match helper functions |
| `os` | Extended OS functions (setalias, resolvealias, getenv, setenv, chdir, cwd) |
| `path` | Path manipulation (getbasename, getdirectory, getextension, join, etc.) |
| `rl` | Readline functions (see line-editing.md) |
| `rl_buffer` | Readline buffer manipulation |
| `settings` | Settings management (get, set) |
| `string` | String utilities (explode, etc.) |
| `unicode` | Unicode handling |

## Module System

Scripts can use Lua's `require()` for sharing code:

```lua
-- In scripts/my_module.lua
local M = {}
M.hello = function() return "world" end
return M

-- In another script
local my_module = require("my_module")
print(my_module.hello())
```

Use `clink.addpackagepath(path)` to add directories to the Lua module search path.

## Debugging Lua Scripts

| Setting | Description |
|---------|-------------|
| `lua.debug` | Enable Lua debugger |
| `lua.break_on_error` | Auto-break on errors |
| `lua.verbose` | Enable verbose logging |

Add `pause()` line in scripts to break at a specific spot for debugging.

## Reloading Scripts

- `Ctrl-X Ctrl-R` (`clink-reload`) — reloads all Lua scripts and .inputrc
- `clink.reload()` — programmatic reload from Lua
- Scripts are reloaded from all script directories
