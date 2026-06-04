# Clink Argmatcher API

In-depth reference for clink's argmatcher system — the primary mechanism for defining command completions in Lua.

## Overview

Argmatchers describe a command's arguments and flags. When clink detects the command being edited, it uses the argmatcher parser to generate matches and apply input line coloring.

```lua
clink.argmatcher([commands...]) : _argmatcher
```

- If one command is provided and an argmatcher already exists for it, returns the existing parser
- Passing more than one command registers the argmatcher with all specified commands
- Passing no arguments creates an unregistered argmatcher (useful for linked argmatchers)
- In v1.3.38+, fully qualified paths can be used for command-specific argmatchers

## Creating an Argmatcher

### Standard Form

```lua
clink.argmatcher("git")
    :addarg({ "add", "status", "commit", "checkout" })
    :addflags("-a", "-g", "-p", "--help")
```

### Shorthand Form

Positional tables without `:addarg()` become argument positions; strings starting with `-` or `/` become flags:

```lua
clink.argmatcher("myapp")
    { "one", "won" }                -- Arg #1
    { "two", "too" }                -- Arg #2
    { "-a", "-b", "/?", "/h" }      -- Flags
```

### Linking Parsers

Parsers can be linked to arguments or flags using concatenation:

```lua
local sub_parser = clink.argmatcher():addarg({ "foo", "bar" })
clink.argmatcher("main"):addarg({ "action" .. sub_parser })
```

A `:` or `=` at end of flag indicates it takes an argument:

```lua
"/f:"..clink.argmatcher():addarg(clink.filematches)
```

## _argmatcher Methods

### `:addarg(choices...) : self`

Adds a new argument position with the specified matches.

**Parameters:** Strings, string linked to parsers via concatenation, tables of arguments, or functions returning argument tables.

**Table entries for `addarg`:**

| Entry | Type | Description | Version |
|-------|------|-------------|---------|
| `delayinit` | function | Delayed initialization for argument position | v1.3.10+ |
| `fromhistory` | boolean | Generate matches from command history | v1.3.9+ |
| `hide` | table/string | Hide strings during completion | v1.8.3+ |
| `hint` | string/function | Show usage hint | v1.7.0+ |
| `loopchars` | string | Delimiter characters for looping | v1.3.37+ |
| `nosort` | boolean | Disable sorting matches | v1.3.3+ |
| `nowordbreakchars` | string | Override word break characters | v1.5.17+ |
| `onadvance` | function | Called before parsing a word | v1.5.14+ |
| `onalias` | function | Called before parsing to alias text | v1.6.18+ |
| `onarg` | function | Called when parsing a word | v1.3.13+ |
| `onlink` | function | Called after parsing a word | v1.5.14+ |

### `:addargunsorted(choices...) : self`

Same as `:addarg()` but disables sorting matches. (v1.3.3+)

### `:addflags(flags...) : self`

Adds flag matches. Flags are position-independent and can appear anywhere.

**Table entries for `addflags`:**

| Entry | Type | Description | Version |
|-------|------|-------------|---------|
| `delayinit` | function | Delayed initialization | v1.3.10+ |
| `fromhistory` | boolean | Generate from history | v1.3.9+ |
| `hide` | table/string | Hide strings | v1.8.3+ |
| `nosort` | boolean | Disable sorting | v1.3.3+ |
| `onarg` | function | Called when parsing | v1.3.13+ |

### `:addflagsunsorted(flags...) : self`

Same as `:addflags()` but disables sorting for flags. (v1.3.3+)

### `:adddescriptions([descriptions...]) : self`

Adds descriptions for arg matches and/or flag matches.

**Two table schemes:**

```lua
-- Scheme 1: Strings with description field
:adddescriptions({ "-h", "--help", description = "Show help" })

-- Scheme 2: Key/value pairs
:adddescriptions({
    ["--user"] = { " name", "Specify user name" },
    ["--port"] = { " port", "Specify port number" },
})
```

### `:chaincommand([modes], [hint]) : self`

Makes the rest of the line be parsed as a separate command.

**Modes (v1.6.2+):**

| Mode | Behavior |
|------|----------|
| `cmd` | Processed like CMD.exe (default) |
| `start` | Processed like START command |
| `run` | Processed like Windows Run dialog |
| `doskey` | Doskey alias expansion (modifier) |

### `:hideflags(flags...) : self`

Hides specified flags from completions but they are still recognized during parsing. (v1.3.3+)

### `:loop([index]) : self`

Makes parser loop back to argument position when out of positional args.

- If `index` omitted, loops back to position 1
- Used for commands accepting repeated arguments

```lua
clink.argmatcher("echo")
    :addarg({ "hello", "world" })
    :loop()  -- After arg #1, loop back to arg #1
```

### `:nofiles() : self`

Prevents invoking match generators and "dead ends" the parser. No file completions will be offered.

### `:reset() : self`

Resets argmatcher to empty state. Clears all flags, arguments, and settings. (v1.3.10+)

### `:setclassifier(func) : self`

Registers a function called for each word to classify for input text coloring. (v1.1.18+)

**Function signature:**

```lua
function(arg_index, word, word_index, line_state, classifications, user_data)
```

- `arg_index` — Argument index (0 = flag)
- `word` — Word being parsed
- `word_index` — Word index in line_state
- `line_state` — line_state object
- `classifications` — word_classifications object
- `user_data` — Table for parsing (v1.5.17+)

See [references/coloring.md](coloring.md) for classification codes.

### `:setdelayinit(func) : self`

Registers a function called the first time the argmatcher is used. (v1.3.10+)

**Function signature:**

```lua
function(argmatcher, command_word)
```

- `argmatcher` — The argmatcher to initialize
- `command_word` — Word that matched (v1.3.12+)

Useful for deferring expensive initialization (e.g., querying available subcommands) until the command is actually used.

### `:setendofflags([endofflags]) : self`

Sets a special flag signaling end of flags. (v1.3.12+)

- String: the end-of-flags flag
- `true` or `nil`: uses `--`
- `false`: clears end-of-flags

### `:setflagprefix([prefixes...]) : self`

Sets flag prefix characters. Usually not needed since `:addflags()` auto-detects prefixes (`-`, `/`).

### `:setflagsanywhere(anywhere) : self`

Controls where flags are recognized. (v1.3.12+)

- `true` (default): flags recognized anywhere
- `false`: flags only recognized until an argument is encountered

## Callback Functions

### Function Arguments (called every time matches are generated)

```lua
local function my_func(word, word_index, line_state, match_builder, user_data)
    -- word: partial string under cursor
    -- word_index: word index in line_state
    -- line_state: line_state object
    -- match_builder: builder object
    -- user_data: table for parsing
    return { "red", "green", "blue" }  -- table of matches
    -- or return true  to stop generating matches
    -- or return false/nil  to stop and use file completions
end
```

### Delayed Init Functions (called once per session)

```lua
local function delayinit_func(argmatcher, argindex)
    -- argmatcher: current argmatcher
    -- argindex: argument position (0=flags, 1+=args)
    return { "value1", "value2" }
end
```

### On Advance Function (v1.5.14+)

```lua
local function onadvance_func(arg_index, word, word_index, line_state, user_data)
    -- Returns:
    -- nil: advance to next position (default)
    -- 1: advance BEFORE parsing current word
    -- 0: repeat same position
    -- -1: chain command (start new command line)
    -- -1, modes: chain with mode specification (v1.6.2+)
end
```

### On Alias Function (v1.6.18+)

```lua
local function onalias_func(arg_index, word, word_index, line_state, user_data)
    -- Returns:
    -- nil: continue parsing normally
    -- string: parse this string instead
    -- string, true: parse as new command
end
```

### On Arg Function (v1.3.13+)

```lua
local function onarg_func(arg_index, word, word_index, line_state, user_data)
    -- Called when parsing a word
    -- Can use os.chdir() to change directory context
end
```

### On Link Function (v1.5.14+)

```lua
local function onlink_func(link, arg_index, word, word_index, line_state, user_data)
    -- link: linked argmatcher if function returns nil
    -- Returns:
    -- argmatcher: override which parser to use
    -- false or nil: continue normally
end
```

### Hint Function (v1.7.0+)

```lua
local function hint_func(arg_index, word, word_index, line_state, user_data)
    return "Hint text", position  -- position is optional
end
```

## Built-in Match Functions

| Function | Description |
|----------|-------------|
| `clink.dirmatches(word)` | Generate directory matches |
| `clink.dirmatchesexact(word)` | Directory matches without `*` append |
| `clink.filematches(word)` | Generate file matches |

## Loop Characters

For arguments accepting multiple values separated by delimiters:

```lua
:addarg({ loopchars=";", "red", "green", "blue" })
-- Accepts: "foo red;green;blue file"
```

## Responsive Argmatchers / User Data

The `user_data` table allows argument positions to share state:

- Set by "on advance", "on alias", "on arg", and "on link" functions
- Each linked argmatcher gets a separate `user_data` table
- `user_data.shared_user_data` (v1.6.10+) enables sharing between linked parsers

## Complete Example

```lua
local function git_remote_actions(word, word_index, line_state, match_builder)
    local handle = io.popen("git remote 2>nul")
    if not handle then return false end
    local remotes = {}
    for line in handle:lines() do
        table.insert(remotes, line)
    end
    handle:close()
    return remotes
end

clink.argmatcher("git")
    :addarg({ "add", "commit", "push", "remote" .. clink.argmatcher()
        :addarg({ "add", "remove", git_remote_actions })
        :addflags("--verbose", "--dry-run")
        :loop(1)
    })
    :addflags("--help", "--version", "-C" .. clink.argmatcher():addarg(clink.dirmatches))
    :setclassifier(function(arg_index, word, word_index, line_state, classifications)
        if arg_index == 0 then
            classifications:classifyword(word_index, "f")
        end
    end)
    :setdelayinit(function(argmatcher, command_word)
        -- Lazy-load subcommands on first use
        local handle = io.popen("git --list-cmds=main 2>nul")
        if handle then
            local cmds = {}
            for line in handle:lines() do
                table.insert(cmds, line)
            end
            handle:close()
            argmatcher:addarg(cmds)
        end
    end)
```

## Deprecated Methods

| Method | Replacement |
|--------|-------------|
| `:add_arguments()` | `:addarg()` |
| `:add_flags()` | `:addflags()` |
| `:set_arguments()` | `:addarg()` |
| `:set_flags()` | `:addflags()` |
| `:be_precise()` | (removed) |
| `:disable_file_matching()` | `:nofiles()` |
