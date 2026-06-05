# Zsh Line Editor (ZLE)

In-depth reference for zsh's line editor — widgets, keymaps, key bindings, completion widgets, and how the completion system interacts with ZLE.

## How ZLE Works

ZLE (Zsh Line Editor) is the command-line editor built into zsh. When the `ZLE` option is set (default in interactive shells) and input is attached to a terminal, ZLE handles all keyboard input and display.

Two display modes:
- **Multiline mode** (default) — works with valid TERM setting
- **Single line mode** — used if TERM is invalid or `SINGLE_LINE_ZLE` is set

## Keymaps

A keymap contains bindings between key sequences and ZLE commands (widgets). The initially available keymaps:

| Keymap | Description |
|--------|-------------|
| `emacs` | EMACS-style editing |
| `viins` | vi insert mode |
| `vicmd` | vi command mode |
| `isearch` | Incremental search mode |
| `command` | Read a command name |
| `.safe` | Fallback keymap (never altered, always available) |
| `main` | Alias linked to either `emacs` or `viins` based on `VISUAL`/`EDITOR` env vars |

### Keymap Selection

```zsh
bindkey -e              # Select emacs and link to main
bindkey -v              # Select viins and link to main
bindkey -A mymap main   # Link mymap to main
```

### Keymap Management

```zsh
bindkey -l              # List all keymap names
bindkey -N newmap       # Create new empty keymap
bindkey -N mymap emacs  # Create keymap from emacs
bindkey -d              # Delete all keymaps, reset to default
bindkey -r key          # Unbind a key sequence
```

## The bindkey Command

`bindkey` manipulates keymaps and key bindings.

### Syntax

```zsh
bindkey [ options ] key sequence widget
```

### Key Notation

| Notation | Meaning |
|----------|---------|
| `^A` or `Ctrl-A` | Control character |
| `^[` or `Esc` | Escape prefix |
| `^Xh` | Ctrl-X followed by h |
| `\C-x` | Control-x (in key string) |
| `\M-x` | Meta-x (Alt-x) |
| `\e` | Escape character |

### Common Bindings for Completion

```zsh
bindkey '^I' complete-word              # TAB (default)
bindkey '^D' delete-char-or-list        # Ctrl-D
bindkey '^Xh' _complete_help            # Ctrl-X h (debug)
bindkey '^X?' _complete_debug           # Ctrl-X ? (trace)
bindkey '^X^I' menu-select              # Ctrl-X TAB (menu selection)
bindkey '\e/' _history_complete_word     # Alt-/ (history completion)
bindkey '^Xc' _correct_word             # Ctrl-X c (correction)
bindkey '^Xa' _expand_alias             # Ctrl-X a (alias expansion)
bindkey '^Xe' _expand_word              # Ctrl-X e (expansion)
```

## ZLE Widgets

All actions in the editor are performed by widgets. Widgets can be:
- **Built-in**: Defined by ZLE or other modules
- **User-defined**: Created with `zle -N`, implemented as shell functions

### Built-in Completion Widgets

| Widget | Default Binding | Description |
|--------|----------------|-------------|
| `complete-word` | unbound | Attempt completion on current word |
| `expand-or-complete` | TAB (emacs) | Shell expansion, then completion |
| `expand-or-complete-prefix` | unbound | Expand up to cursor, then complete |
| `menu-complete` | TAB (viins) | Like complete-word with menu cycling |
| `menu-select` | unbound | Menu selection with cursor navigation |
| `reverse-menu-complete` | unbound | Menu completion going backward |
| `delete-char-or-list` | ^D (emacs) | Delete char; at EOL, list completions |
| `list-choices` | ESC-^D | List possible completions |
| `accept-and-menu-complete` | unbound | Insert match, advance to next |

### Menu-Select Widget

`menu-select` provides interactive selection with cursor navigation through completions. Enable with:

```zsh
zstyle ':completion:*' menu select
bindkey '^X\t' menu-select
```

When `menu select` is set in zstyle, the completion system automatically uses menu selection when there are multiple matches.

### Other Useful Widgets

| Widget | Default Binding | Description |
|--------|----------------|-------------|
| `_complete_help` | ^Xh | Show context names, tags, and functions |
| `_complete_debug` | ^X? | Debug completion with trace file |
| `_correct_word` | ^Xc | Correct current argument |
| `_expand_alias` | ^Xa | Expand alias at cursor |
| `_expand_word` | ^Xe | Expand current word |
| `_history_complete_word` | \e/ | Complete from history |
| `_most_recent_file` | ^Xm | Complete most recently modified file |
| `_next_tags` | ^Xn | Cycle through next tag set |

## Creating User-Defined Widgets

### Standard Widgets

```zsh
my-widget() {
    # Access special parameters
    local buf=$BUFFER
    local cur=$CURSOR
    # ... do something ...
    zle -U "text"            # Push characters onto input stack
    zle reset-prompt          # Redraw prompt
}
zle -N my-widget my-widget
bindkey '^Xw' my-widget
```

### Completion Widgets

Completion widgets are defined with `zle -C`:

```zsh
zle -C widget-name base-widget function-name
```

| Argument | Description |
|----------|-------------|
| `widget-name` | Name of the new widget |
| `base-widget` | Built-in widget that handles match insertion (e.g., `complete-word`, `menu-complete`) |
| `function-name` | Shell function that generates matches |

Example:

```zsh
zle -C my-complete complete-word _my_completer
bindkey '^X\t' my-complete
```

When the widget is triggered:
1. The shell function `_my_completer` is called with stdin closed
2. The function uses `compadd`, `compset`, etc. to generate matches
3. After the function returns, the base widget (`complete-word`) handles insertion

## The zle Command

```zsh
zle -l              # List user-defined widgets
zle -l -a           # List all widgets
zle -N widget func  # Create user-defined widget
zle -C widget base func  # Create completion widget
zle -D widget       # Delete widget
zle -A old new      # Alias widgets
zle -R [str]        # Redisplay command line
zle widget [args]   # Invoke widget from function
```

## Special Parameters Inside Widgets

| Parameter | Type | Description |
|-----------|------|-------------|
| `BUFFER` | String | Entire edit buffer contents |
| `CURSOR` | Integer | Cursor offset (0 to `$#BUFFER`) |
| `LBUFFER` | String | Portion of buffer left of cursor |
| `RBUFFER` | String | Portion of buffer right of cursor |
| `KEYMAP` | String | Current keymap name (read-only) |
| `WIDGET` | String | Name of widget being executed (read-only) |
| `KEYS` | String | Keys typed to invoke widget (read-only) |
| `NUMERIC` | Integer | Numeric argument (from prefix digit) |
| `CONTEXT` | String | Context: `start`, `cont`, `select`, `vared` (read-only) |
| `WIDGETFUNC` | String | Function name for user-defined widgets (read-only) |
| `WIDGETSTYLE` | String | Base widget name for completion widgets (read-only) |

## The compstate Variable

`compstate` is an associative array that exchanges information between the completion code and the widget. It is the primary mechanism for controlling completion behavior from within a completion function.

### Context (compstate[context])

| Value | When |
|-------|------|
| `command` | Normal command completion |
| `array_value` | Inside array parameter assignment value |
| `brace_parameter` | Parameter name in `${...}` |
| `assign_parameter` | Parameter name in assignment |
| `condition` | Inside `[[...]]` conditional |
| `math` | Inside `((...))` construct |
| `parameter` | Parameter name after `$` |
| `redirect` | After redirection operator |
| `subscript` | Inside parameter subscript |
| `value` | Value of parameter assignment |

### Insertion Control (compstate[insert])

| Value | Effect |
|-------|--------|
| unset | No change to command line |
| `unambiguous` | Insert common prefix |
| `automenu-unambiguous` | Insert common prefix, may start menu completion |
| `menu` or `automenu` | Start menu completion |
| `tab` | Only insert TAB character |
| `all` | Insert all matches |
| Number | Insert that match (negative counts backward, -1 = last) |

Values ending in space insert without space. `menu:2` starts menu with second match.

### List Control (compstate[list])

| Value | Effect |
|-------|--------|
| unset/empty | Never list |
| `list` | Always list |
| `autolist` | List when `AUTO_LIST` would trigger |
| `ambiguous` | List when `LIST_AMBIGUOUS` would trigger |
| `force` | Show even with single match |
| `packed` | `LIST_PACKED` behavior for the group |
| `rows` | `LIST_ROWS_FIRST` behavior |
| `explanations` | Only explanation strings listed |
| `messages` | Only messages added with `-x` listed |

### Other compstate Keys

| Key | Description |
|-----|-------------|
| `all_quotes` | Characters showing quoting levels broken by `compset -q` |
| `exact` | Set to `accept` if `REC_EXACT` would accept; otherwise unset |
| `exact_string` | String of exact match found, otherwise unset |
| `ignored` | Number of words ignored due to `-F` patterns |
| `insert_positions` | Colon-separated positions where chars differ in unambiguous string |
| `last_prompt` | Non-empty for every match = move cursor back to previous prompt |
| `list_lines` | Number of lines needed to display completions (read-only) |
| `list_max` | Initially `LISTMAX`; can be set to override |
| `nmatches` | Number of matches generated and accepted (read-only) |
| `old_insert` | Match number from old list inserted; unset if none |
| `old_list` | `yes` if valid list from previous completion; `shown` if also on screen |
| `parameter` | Parameter name when in subscript or value assignment |
| `pattern_insert` | `menu` if `GLOB_COMPLETE` set; controls how pattern matches insert |
| `pattern_match` | Controls `GLOB_COMPLETE` behavior locally; `**` adds wildcard at cursor |
| `quote` | Quotation character if inside quotes (single, double, or backtick) |
| `quoting` | `single`, `double`, or `backtick` depending on quote type |
| `redirect` | Redirection operator (`<`, `>`, etc.) |
| `restore` | `auto` before function entry = restore special params on exit |
| `to_end` | `single` = cursor moves on single unambiguous; `match` = always; `always`/`empty` = always/never |
| `unambiguous` | Common prefix for all matches added so far (read-only) |
| `unambiguous_cursor` | Cursor position relative to unambiguous string if inserted |
| `unambiguous_positions` | Positions where chars differ in unambiguous string |
| `vared` | Parameter name given to `vared` if active |

## ZLE Options and Variables

| Variable | Description |
|----------|-------------|
| `KEYTIMEOUT` | Timeout for key sequences (default 0.4 sec; 0.01 for instant) |
| `zle_highlight` | Character highlighting configuration array |
| `WORDCHARS` | Characters treated as part of a word (default: `*?_-.[]~=/&;!#$%^(){}<>`) |

### Character Highlighting (zle_highlight)

```zsh
zle_highlight=(region:standout special:standout isearch:underline)
zle_highlight=(none)  # Disable all highlighting
```

Highlighting regions:

| Region | Description |
|--------|-------------|
| `region` | Selected text region |
| `special` | Special characters (line endings, etc.) |
| `isearch` | Incremental search match |
| `suffix` | Auto-removable suffix |
| `paste` | Pasted text |

## The vared Command

`vared` invokes ZLE on a shell parameter value, allowing interactive editing:

```zsh
vared -c -p "prompt: " -r "right prompt" varname
```

| Option | Description |
|--------|-------------|
| `-c` | Create the parameter if it doesn't exist |
| `-p prompt` | Left prompt |
| `-r rprompt` | Right prompt |
| `-M keymap` | Use specified keymap |
| `-h` | Use history |

## References

### Documentation

- [Zsh Manual: Zsh Line Editor](https://zsh.sourceforge.io/Doc/Release/Zsh-Line-Editor.html) — official ZLE reference
- [Zsh Manual: Completion Widgets](https://zsh.sourceforge.io/Doc/Release/Zsh-Modules.html#Completion-Widgets) — `compadd`, `compstate`, `compquote` builtins
- [zshzle(1) man page](https://linux.die.net/man/1/zshzle) — ZLE reference
- [zshcompwid(1) man page](https://linux.die.net/man/1/zshcompwid) — completion widgets reference

### Source Code

- [zsh ZLE source](https://github.com/zsh-users/zsh/tree/master/Src/Zle) — C implementation of ZLE
- [zsh `complete.c`](https://github.com/zsh-users/zsh/blob/master/Src/Zle/complete.c) — core completion widget implementation
- [zsh `compcore.c`](https://github.com/zsh-users/zsh/blob/master/Src/Zle/compcore.c) — completion engine internals
- [zsh `compresult.c`](https://github.com/zsh-users/zsh/blob/master/Src/Zle/compresult.c) — match insertion and display
