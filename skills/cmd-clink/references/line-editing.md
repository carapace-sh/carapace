# Clink Readline Integration and Line Editing

In-depth reference for clink's GNU Readline integration — key bindings, .inputrc configuration, editing modes, completion commands, and clink-specific extensions.

## How Clink Uses Readline

Clink embeds the GNU Readline library (the same library used by bash) and replaces cmd.exe's native line editing with Readline-powered editing. This happens via DLL injection — the hooked `ReadConsoleW()` function is replaced with Readline's input handling.

### Key Differences from Bash's Readline

| Aspect | Bash | Clink |
|--------|------|-------|
| Platform | Unix/POSIX | Windows |
| Shell | bash | cmd.exe |
| Word breaks | `COMP_WORDBREAKS` | `nowordbreakchars` in argmatchers |
| Completion registration | `complete -F func cmd` | `clink.argmatcher("cmd")` or Lua generators |
| Key binding extension | `bind -x` | `luafunc:` prefix in .inputrc |
| Meta key | Alt key | Alt key (always enabled) |
| Eight-bit input | Configurable | Always enabled |
| Bracketed paste | Supported | Use `clink-paste` command instead |

## The .inputrc File

### File Locations (searched in order)

1. `%CLINK_INPUTRC%` — if set, used instead of default
2. Clink profile directory
3. `%USERPROFILE%`
4. `%LOCALAPPDATA%`
5. `%APPDATA%`
6. `%HOME%` or `%HOMEDRIVE%%HOMEPATH%`

### Syntax

```inputrc
# Comment lines
set variable value           # Set a Readline variable
keyname: function-name       # Bind key to Readline function
"keyseq": function-name      # Bind key sequence to function
"keyseq": "macro text"       # Bind key sequence to macro
"keyseq": "luafunc:name"     # Bind key to Lua function (clink-specific)
$include /path/to/file       # Include another init file
$if condition                 # Conditional construct
$else
$endif
```

### Conditional Constructs

| Construct | Purpose |
|-----------|---------|
| `$if mode=emacs` | Test editing mode (emacs/vi) |
| `$if term=xterm` | Test terminal type |
| `$if clink` | Clink-specific settings (not processed by bash) |
| `$if clink_version >= 1.6.1` | Test clink version (supports `=`, `==`, `!=`, `<=`, `>=`, `<`, `>`) |
| `$if variable == value` | Test Readline variable equality |
| `$else` | Execute if test fails |
| `$endif` | Close conditional block |

### Key Binding Syntax

**Keyname format (unquoted):**
```inputrc
Control-u: universal-argument
Meta-Rubout: backward-kill-word
Space: self-insert
```

**String format (quoted, for key sequences):**
```inputrc
"\C-x\C-r": re-read-init-file
"\e[11~": "Function Key 1"
"\t": old-menu-complete
```

**Escape sequences:**

| Sequence | Meaning |
|----------|---------|
| `\C-` | Control prefix |
| `\M-` | Meta prefix |
| `\e` | Escape character |
| `\\` | Backslash |
| `\"` | Double quote |
| `\'` | Single quote |
| `\a` | Bell |
| `\b` | Backspace |
| `\d` | Delete |
| `\f` | Form feed |
| `\n` | Newline |
| `\r` | Carriage return |
| `\t` | Tab |
| `\v` | Vertical tab |
| `\nnn` | Octal character value (1-3 digits) |
| `\xHH` | Hexadecimal character value |

**Lua function binding (clink-specific):**
```inputrc
"\e[27;8;72~": "luafunc:my_lua_function"
```

## Editing Modes

### Emacs Mode (Default)

Standard Emacs-style key bindings. Control keys operate on characters, Meta keys on words.

### Vi Mode

Enable with `set editing-mode vi` in `.inputrc`. Toggle between modes with `Alt-Ctrl-j`.

- Starts in insertion mode
- Press `Esc` for command mode
- Standard vi movement: `h`/`l` (char), `b`/`w` (word), `0`/`$` (line)

### Keymaps

| Keymap | When Active |
|--------|-------------|
| `emacs` | Emacs mode default |
| `emacs-standard` | Emacs mode (same as `emacs`) |
| `emacs-meta` | Meta key prefix in emacs mode |
| `emacs-ctlx` | `C-x` prefix in emacs mode |
| `vi` | Vi mode default |
| `vi-command` | Vi command mode |
| `vi-move` | Vi movement mode (same as `vi-command`) |
| `vi-insert` | Vi insertion mode |

## Readline Configuration Variables

### Standard Variables Supported by Clink

| Variable | Description |
|----------|-------------|
| `bell-style` | `none`, `visible` (default), or `audible` |
| `blink-matching-paren` | Flash cursor on matching parenthesis |
| `colored-completion-prefix` | Color common prefix in completions |
| `colored-stats` | Color completions by file type |
| `comment-begin` | Comment string (for CMD, usually `#`) |
| `completion-ignore-case` | Case-insensitive completion |
| `editing-mode` | `emacs` or `vi` |
| `enable-active-region` | Enable text selection highlighting |
| `expand-tilde` | Expand tilde in word completion |
| `history-point-at-end-of-anchored-search` | Cursor at end when searching history |
| `history-search-ignore-case` | Case-insensitive history search |
| `horizontal-scroll-mode` | Scroll horizontally instead of wrapping |
| `isearch-terminators` | Keys that terminate incremental search |
| `keymap` | Set keymap |
| `mark-directories` | Append `/` to directory completions |
| `mark-symlinked-directories` | Mark symlinked directories with `@` |
| `match-hidden-files` | Match hidden files in completion |
| `output-meta` | Display eight-bit characters (always on in clink) |
| `print-completions-horizontally` | List completions horizontally |
| `show-all-if-ambiguous` | Show all matches on ambiguous completion |
| `show-all-if-unmodified` | Show all matches if no unique match |
| `show-mode-in-prompt` | Show editing mode in prompt |
| `skip-completed-text` | Skip common prefix when completing |
| `visible-stats` | Add indicator after completions |

### Clink-Specific Readline Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `completion-auto-query-items` | on | Prompt before showing many completions |
| `menu-complete-wraparound` | on | Wrap around in menu-complete |
| `search-ignore-case` | on | Case-insensitive history search |

### Variables with Different Defaults than Bash

| Variable | Bash Default | Clink Default | Notes |
|----------|-------------|---------------|-------|
| `bell-style` | `audible` | `visible` | Clink prefers visual bell |
| `search-ignore-case` | off | on | Clink enables case-insensitive search by default |
| `completion-auto-query-items` | off | on | Clink prompts before listing many completions |

### Unsupported/Deprecated Variables

| Variable | Status |
|----------|--------|
| `bind-tty-special-chars` | Not used by clink |
| `completion-map-case` | Use `match.ignore_case` clink setting instead |
| `convert-meta` | Always on |
| `enable-bracketed-paste` | Not supported — use `clink-paste` command |
| `enable-meta-key` | Always on |
| `history-size` | Use `history.max_lines` clink setting instead |
| `input-meta` | Always on |
| `keyseq-timeout` | Not supported |
| `output-meta` | Always on |

## Completion Commands

| Command | Default Key | Description |
|---------|-------------|-------------|
| `complete` | Tab | Standard completion |
| `old-menu-complete` | Tab (enhanced bindings) | CMD-style cycling completion |
| `menu-complete` | (unbind) | Cycle through completions |
| `menu-complete-backward` | (unbind) | Cycle backward |
| `possible-completions` | `Alt-=` | List all completions |
| `insert-completions` | `Alt-*` | Insert all completions |
| `clink-select-complete` | `Ctrl-Space` | Interactive completion selection |

### clink-select-complete

Interactive completion list with rich navigation:

- Arrow keys to navigate
- Typing filters the list
- `F4` toggles "find" vs "filter" search mode
- `F1` toggles description display (inline vs bottom)
- `Ctrl-Home`/`Ctrl-End` select first/last match
- `match.max_rows` limits displayed rows
- `match.preview_rows` shows preview
- `Shift-Enter`/`Ctrl-Enter` insert without executing

## Clink-Specific Commands

| Command | Default Key | Description |
|---------|-------------|-------------|
| `add-history` | `Alt-Ctrl-K` | Add line to history without executing |
| `clink-accept-suggested-line` | (unbind) | Accept auto-suggestion |
| `clink-backward-bigword` | (unbind) | Move back space-delimited word |
| `clink-complete-numbers` | (unbind) | Complete numbers from screen |
| `clink-copy-cwd` | `Alt-C` | Copy current directory |
| `clink-copy-line` | `Alt-Ctrl-C` | Copy input line |
| `clink-copy-word` | `Alt-Ctrl-W` | Copy word at cursor |
| `clink-ctrl-c` | `Ctrl-C` | Copy selection or cancel line |
| `clink-diagnostics` | `Ctrl-X Ctrl-Z` | Show diagnostics |
| `clink-dump-functions` | (unbind) | List functions/commands |
| `clink-dump-macros` | (unbind) | List macros |
| `clink-dump-variables` | (unbind) | List variables |
| `clink-echo` | (unbind) | Echo key sequences |
| `clink-expand-doskey-alias` | (unbind) | Expand doskey alias |
| `clink-expand-history` | `Alt-^` | Expand history |
| `clink-popup-directories` | `Ctrl-Alt-PgUp` | Popup directory list |
| `clink-popup-history` | `F7` | Popup history list |
| `clink-popup-show-help` | `F1` | Popup help |
| `clink-reload` | `Ctrl-X Ctrl-R` | Reload config and scripts |
| `clink-select-complete` | `Ctrl-Space` | Interactive completion selection |
| `clink-show-help` | `Alt-H` | Show all key bindings |
| `remove-history` | `Alt-Ctrl-D` | Remove history entry |

## Key Binding Styles

### Bash-style (Default)

```inputrc
$if clink
    set keymap emacs
    "\t":           old-menu-complete
    "\e[Z":         old-menu-complete-backward
    "\e[27;5;32~":  clink-select-complete
$endif
```

### Windows-style

Enable with `clink set clink.default_bindings windows`:

- `Ctrl-A`: Select All
- `Ctrl-F`: Find text
- `Ctrl-M`: Mark text
- `Tab`: Cycle through completions

## Sample .inputrc for Clink

```inputrc
$if clink
    set colored-completion-prefix    on
    set colored-stats                on
    set mark-symlinked-directories   on
    set visible-stats                off
    set completion-auto-query-items  on
    set menu-complete-wraparound     off
    set search-ignore-case           on

    set keymap emacs

    "\t":           old-menu-complete
    "\e[Z":         old-menu-complete-backward
    "\e[27;5;32~":  clink-select-complete

    "\e[A":         history-search-backward
    "\e[B":         history-search-forward
    "\e[5~":        clink-popup-history

    "\C-x\C-f":     clink-dump-functions
    "\C-x\C-m":     clink-dump-macros
    "\C-x\C-v":     clink-dump-variables
$endif
```

## The `bind` Builtin Equivalent

Clink does not have a `bind` builtin like bash. Key bindings are configured exclusively through:

1. `.inputrc` file
2. Lua scripts using `rl.setbinding()`
3. `clink set` for clink-specific settings

## rl Lua API

| Function | Description |
|----------|-------------|
| `rl.getvariable(name)` | Get Readline variable value |
| `rl.isvariabletrue(name)` | Check if Readline variable is true |
| `rl.setbinding(key, command)` | Bind key to command |
| `rl.getkeybindings()` | Get list of key bindings |
| `rl.gethistorycount()` | Number of history items |
| `rl.gethistoryitems(start, end)` | Table of history items |
| `rl.getpromptinfo()` | Prompt information table |
| `rl.getinputrcfilename()` | Path to .inputrc file |
| `rl.bracketpromptcodes(text)` | Surround escape codes with `\001`/`\002` |

## rl_buffer API

Available in Lua key binding functions:

| Method | Description |
|----------|-------------|
| `rl_buffer:getbuffer()` | Get current line text |
| `rl_buffer:setbuffer(text)` | Set current line text |
| `rl_buffer:getcursor()` | Get cursor position |
| `rl_buffer:setcursor(pos)` | Set cursor position |
| `rl_buffer:getmark()` | Get mark position |
| `rl_buffer:setmark(pos)` | Set mark position |
| `rl_buffer:insert(text)` | Insert text at cursor |
| `rl_buffer:delete(from, to)` | Delete text range |
| `rl_buffer:beginundogroup()` | Start undo group |
| `rl_buffer:endundogroup()` | End undo group |
| `rl_buffer:hassuggestion()` | Check if suggestion available (v1.6.4+) |
| `rl_buffer:insertsuggestion([amount])` | Insert suggestion (v1.6.4+) |
