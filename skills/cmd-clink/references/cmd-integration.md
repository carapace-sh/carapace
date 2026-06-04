# Clink cmd.exe Integration

In-depth reference for how clink integrates with Windows cmd.exe — the injection mechanism, autorun, profile directory, settings system, prompt customization, and cmd.exe-specific features.

## How Clink Injects into cmd.exe

Clink works by being **injected** into a `cmd.exe` process, where it intercepts a handful of Windows API functions to replace the prompt and input line editing with Readline-powered enhancements.

### The Injection Process

1. `clink.bat` detects the current CPU mode (x64/x86/ARM64) and invokes the correct `clink_*.exe`
2. `clink_*.exe` checks whether the parent process is supported (must be `cmd.exe` with `/k` or without `/c`/`/r`)
3. `clink_*.exe` injects a remote thread into the parent CMD process
4. `clink_x64_dll.dll` is loaded into the CMD process
5. The DLL hooks OS functions

### Hooked OS Functions

| Function | Reason |
|----------|--------|
| `GetEnvironmentVariableW()` | Let the DLL finish its initialization (always called before CMD displays prompt) |
| `SetEnvironmentVariableW()` | Intercept setting the `PROMPT` variable, add a special tag |
| `SetEnvironmentStringsW()` | Intercept setting the `PROMPT` variable, add a special tag |
| `WriteConsoleW()` | Capture the current prompt and defer printing it |
| `ReadConsoleW()` | Replace command line editing with Readline-powered editing |
| `SetConsoleTitleW()` | Enable replacing the "Administrator:" prefix in the title bar |

### Order of Execution During AutoRun

1. CMD executes commands in its AutoRun regkey
2. `clink.bat` checks whether to skip injecting (e.g., `%CLINK_NOAUTORUN%`)
3. `clink.bat` executes `clink_x64.exe` with `inject` command
4. `clink_x64.exe` injects a remote thread into parent CMD process
5. `clink_x64_dll.dll` is loaded and hooks OS functions
6. CMD executes any commands from command line arguments (`/c` or `/k`)
7. If a startup script is present, Clink intercepts the prompt and runs it
8. Clink initializes its Lua engine and loads Lua scripts

## Starting Clink

### Method 1: Autorun (Recommended)

```cmd
clink autorun install
```

Installs clink's entry in CMD.EXE's AutoRun registry key. Clink starts automatically with every new cmd.exe session.

```cmd
clink autorun uninstall   -- Remove autorun
clink autorun show        -- Show current autorun configuration
```

**Registry key:** `HKCU\Software\Microsoft\Command Processor\AutoRun`

### Method 2: Manual Start

Run the Clink shortcut from Start menu or `clink.bat`.

### Method 3: Inject into Existing cmd.exe

```cmd
clink inject
```

Injects clink into the current CMD process. Useful for adding clink to an already-running session.

### Preventing Injection

| Method | Effect |
|--------|--------|
| `cmd /d` | Disables CMD's AutoRun regkey processing entirely |
| `%CLINK_NOAUTORUN%` | Clink autorun script exits quickly |
| `clink autorun uninstall` | Removes autorun entry permanently |

## Profile Directory

### Default Locations

| Windows Version | Path |
|----------------|------|
| Windows XP | `c:\Documents and Settings\username\Local Settings\Application Data\clink` |
| Windows Vista+ | `c:\Users\username\AppData\Local\clink` |

### Override Methods (priority order)

1. `--profile path` command line option when using `clink inject`
2. `%CLINK_PROFILE%` environment variable

### Files in Profile Directory

| File | Purpose |
|------|---------|
| `.inputrc` | Readline configuration (key bindings, settings) |
| `clink_settings` | Clink settings |
| `clink_history` | Command history |
| `clink.log` | Diagnostic log file |
| `default_settings` | Optional default values for settings |
| `default_inputrc` | Optional default Readline init file |

## Settings System

### Commands

```cmd
clink set                        -- List all settings and values
clink set --describe             -- List all settings with descriptions
clink set setting_name           -- Describe setting and show current value
clink set setting_name value     -- Set a value
clink set setting_name clear     -- Reset to default
```

### Override Settings File Location

`%CLINK_SETTINGS%` environment variable (not recommended — breaks `clink set`).

### Key Settings Categories

| Category | Key Settings |
|----------|-------------|
| `autosuggest.*` | Auto-suggestion configuration |
| `clink.*` | Core clink behavior (default_bindings, etc.) |
| `color.*` | Color settings for input, completions, prompts |
| `doskey.*` | Doskey expansion behavior |
| `exec.*` | Execution settings |
| `history.*` | History configuration |
| `match.*` | Matching and completion configuration |
| `terminal.*` | Terminal integration |
| `textform.*` | Text formatting |

## Environment Variables

| Variable | Purpose |
|----------|---------|
| `%CLINK_PROFILE%` | Profile directory location |
| `%CLINK_SETTINGS%` | Settings file location |
| `%CLINK_NOAUTORUN%` | Skip automatic injection |
| `%CLINK_ANSI_HOST%` | Override terminal detection |
| `%CLINK_INPUTRC%` | Custom .inputrc location |
| `%CLINK_PATH%` | Additional Lua script directories |
| `%CLINK_COMPLETIONS_DIR%` | Additional completion script directories |
| `%CLINK_THEMES_DIR%` | Prompt theme directories |
| `%CLINK_CUSTOMPROMPT%` | Override custom prompt file |
| `%CLINK_COLORTHEME%` | Override color theme file |
| `%CLINK_RPROMPT%` | Right-side prompt string |
| `%CLINK_PROMPT_PREFIX%` | Escape codes before prompt |
| `%CLINK_PROMPT_SUFFIX%` | Escape codes after prompt |
| `%CLINK_TRANSIENT_PROMPT%` | Transient prompt string |
| `%CLINK_TRANSIENT_RPROMPT%` | Transient right-side prompt |
| `%CLINK_MATCH_COLORS%` | Match coloring rules |
| `%CLINK_HISTORY_LABEL%` | History file label (up to 32 alphanumeric chars) |
| `%CLINK_TERM_VE%` | Normal cursor style escape codes |
| `%CLINK_TERM_VS%` | Insert mode cursor style escape codes |
| `%CLINK_MOUSE_MODIFIER%` | Mouse input modifier keys |
| `%LS_COLORS%` | File coloring (Readline compatibility) |

## cmd.exe Features

### Command Separators

| Separator | Meaning |
|-----------|---------|
| `&` | Sequential execution (run left, then right) |
| `|` | Pipe (redirect output from left to right) |
| `&&` | Conditional AND (run right only if left succeeds) |
| `||` | Conditional OR (run right only if left fails) |

Clink colors separators with `color.cmdsep`.

### Redirection

| Symbol | Meaning |
|--------|---------|
| `>` | Redirect stdout to file |
| `>>` | Append stdout to file |
| `2>` | Redirect stderr to file |
| `2>&1` | Redirect stderr to stdout |

Clink colors redirection symbols with `color.cmdredir`.

### Doskey Aliases

Doskey macros use syntax: `name=command $*` or `name=command $1 $2`

**Macro tokens:**

| Token | Meaning |
|-------|---------|
| `$*` | All arguments |
| `$1`–`$9` | Positional arguments |
| `$T` | Command separator in macros |

**Enhanced Doskey Expansion** (enabled by default via `doskey.enhanced`):
- Expands doskey macros that follow `|` and `&` command separators
- Respects quotes around words when expanding `$1`–`$9` tags

**Suppressing expansion:**

| Prefix | Doskey Expansion | History |
|--------|-----------------|--------|
| (none) | Yes | Yes |
| Space | No | No |
| Semicolon | No | Yes |

**API functions:** `os.setalias()`, `os.resolvealias()`, `clink.getalias()`

### Built-in CMD Commands

Clink recognizes and colors built-in CMD commands using `color.cmd`. Examples: `dir`, `cd`, `set`, `call`, `if`, `for`, `echo`, `type`, `del`, `copy`, `move`, `mkdir`, `rmdir`, `start`, etc.

## Prompt Customization

### Prompt Filters

```lua
local p = clink.promptfilter(priority)
function p:filter(prompt)
    return "new prefix "..prompt.." new suffix"
end
```

- `priority` determines call order (lower values first)
- Return `nil` → no effect
- Return string → new filtered prompt
- Return string + `false` → stops further filtering

### Right-Side Prompt

```lua
function p:rightfilter(rprompt)
    return os.date().."  "..rprompt
end
```

Displayed on the right side of the first input line. Hidden when input reaches it.

### Asynchronous Prompt Filtering

```lua
local function collect_git_info()
    -- Runs in background coroutine
    return { branch = git.getbranch() }
end

local git_prompt = clink.promptfilter(55)
function git_prompt:filter(prompt)
    local info = clink.promptcoroutine(collect_git_info)
    if info == nil then return end
    return prompt.." "..info.branch
end
```

- `clink.promptcoroutine(func)` runs `func()` in a coroutine
- Returns `nil` until complete, then returns the function's result
- Automatically refreshes the prompt when complete
- Each prompt filter can have at most one prompt coroutine
- Optional `cookie` argument (v1.7.0+) allows multiple coroutines per filter

### Transient Prompts

Replace past prompts with a simplified version:

```cmd
clink set prompt.transient always     -- Always use transient
clink set prompt.transient same_dir   -- Only when directory unchanged
```

```lua
function pf:transientfilter(prompt)
    return "> "  -- Simplified past prompt
end

function pf:transientrightfilter(rprompt)
    return ""  -- Hide right-side for past prompts
end
```

### .clinkprompt Files

Custom prompts packaged for easy sharing:

```cmd
clink config prompt use prompt_name    -- Activate
clink config prompt list               -- List available
clink config prompt show prompt_name    -- Preview
```

**Search locations:**
1. Directories in `%CLINK_THEMES_DIR%`
2. `themes\` subdirectory under each scripts directory
3. Full path to a file

**Exports table (optional):**

```lua
local exports = {
    onactivate = function_name,
    ondeactivate = function_name,
    demo = function_name,
    dependson = "dependency_name",
}
return exports
```

### Prompt Escape Codes

After all filters finish, prompts are surrounded with:

| Variable | Purpose |
|----------|---------|
| `%CLINK_PROMPT_PREFIX%` | Escape codes before prompt |
| `%CLINK_PROMPT_SUFFIX%` | Escape codes after prompt |
| `%CLINK_RPROMPT_PREFIX%` | Escape codes before right prompt |
| `%CLINK_RPROMPT_SUFFIX%` | Escape codes after right prompt |

Individual filters can add escape codes via `:surround()`:

```lua
function p:surround()
    return prompt_prefix, prompt_suffix, rprompt_prefix, rprompt_suffix
end
```

### Included Custom Prompts

| Prompt | Style |
|--------|-------|
| agnoster | Minimalistic with git branch |
| antares | Box-style with git info |
| bureau | Full info with user/host/path/time |
| darkblood | Dark theme variant |
| headline | Horizontal line style |
| jonathan | Double-line box style |
| oh-my-posh | Requires oh-my-posh program |
| pure | Minimal with git status symbols |
| starship | Requires starship program |

## Startup Script

When Clink is injected, it looks for `clink_start.cmd` in the binaries directory and profile directory. Clink automatically runs the script(s) when the first CMD prompt is shown.

Override with `clink.autostart` setting. Set to `"nul"` to run no command.

## Terminal Support

Clink's keyboard driver produces VT220 style key sequences with Xterm extensions.

Override terminal detection with `%CLINK_ANSI_HOST%`:

| Value | Terminal |
|-------|----------|
| `ansicon` | ANSICON |
| `clink` | Clink's built-in |
| `conemu` | ConEmu |
| `wezterm` | WezTerm |
| `winconsole` | Windows Console (legacy) |
| `winconsolev2` | Windows Console (v2) |
| `winterminal` | Windows Terminal |

Mouse input modes: `off`, `on`, `auto` (via `terminal.mouse_input` setting).
