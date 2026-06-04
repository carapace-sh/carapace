# Bash Readline

In-depth reference for the GNU Readline library as used by bash — key bindings, variables, editing modes, macros, and how Readline controls the completion display.

## How Readline Works

Readline is the line-editing library that bash uses for interactive input. It sits between the user and the terminal:

1. **Initialization** — reads `~/.inputrc` (or `$INPUTRC`) for key bindings and variable settings
2. **Character input** — reads characters from the terminal, interprets key sequences, maps them to editing commands
3. **Buffer management** — text is stored in an internal buffer; characters are inserted at the cursor position
4. **Re-reading config** — `C-x C-r` re-reads the init file to incorporate changes

## The .inputrc File

### File Locations

| Location | Priority |
|----------|----------|
| `$INPUTRC` | If set, used instead of default |
| `~/.inputrc` | Default user config |
| `/etc/inputrc` | Fallback if user file doesn't exist |

### Syntax

```
# Comment lines
set variable value           # Set a Readline variable
keyname: function-name       # Bind key to Readline function
"keyseq": function-name      # Bind key sequence to function
"keyseq": "macro text"       # Bind key sequence to macro
$include /path/to/file       # Include another init file
$if condition                 # Conditional construct
$else
$endif
```

### Conditional Constructs

| Construct | Purpose |
|-----------|---------|
| `$if mode=emacs` | Test editing mode (emacs/vi) |
| `$if term=xterm` | Test terminal type (matches full name or prefix) |
| `$if application=Bash` | Test application name |
| `$if version >= 8.0` | Test Readline version (supports `=`, `==`, `!=`, `<=`, `>=`, `<`, `>`) |
| `$if variable == value` | Test Readline variable equality |
| `$else` | Execute if test fails |
| `$endif` | Close conditional block |
| `$include /path/file` | Include another init file |

### Key Binding Syntax

**Keyname format:**
```
Control-u: universal-argument
Meta-Rubout: backward-kill-word
```

**String format (for key sequences):**
```
"\C-x\C-r": re-read-init-file
"\e[11~": "Function Key 1"
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

## The `bind` Builtin

```bash
bind [-m keymap] [-lsvSVX]
bind [-m keymap] -f filename
bind [-m keymap] -x keyseq:shell-command
bind [-m keymap] keyseq:function-name
bind [-m keymap] keyseq:readline-command
bind [-p|-P] [readline-command]
bind [-q function]
bind [-u function]
bind [-r keyseq]
```

### Options

| Option | Description |
|--------|-------------|
| `-m keymap` | Use specified keymap for subsequent bindings |
| `-l` | List names of all Readline functions |
| `-p` / `-P` | Display bindings in reusable format |
| `-s` / `-S` | Display macro bindings |
| `-v` / `-V` | Display variable names and values |
| `-f filename` | Read bindings from file |
| `-q function` | Show key sequences invoking function |
| `-u function` | Unbind all keys for function |
| `-r keyseq` | Remove binding for key sequence |
| `-x keyseq:shell-command` | Execute shell command when keyseq entered |
| `-X` | List key sequences bound to shell commands via `-x` |

### Variables Available to `bind -x` Shell Commands

| Variable | Description |
|----------|-------------|
| `READLINE_LINE` | Contents of Readline line buffer |
| `READLINE_POINT` | Current insertion point position |
| `READLINE_MARK` | Saved insertion point (mark) |
| `READLINE_ARGUMENT` | Numeric argument supplied by user |

These can be read and modified — changing `READLINE_LINE` or `READLINE_POINT` updates the line buffer.

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

## Editing Modes

### Emacs Mode (Default)

Standard Emacs-style key bindings. Control keys operate on characters, Meta keys on words.

**Essential commands:**

| Key | Command |
|-----|---------|
| `C-a` | Move to start of line |
| `C-e` | Move to end of line |
| `C-f` | Move forward one character |
| `C-b` | Move backward one character |
| `M-f` | Move forward one word |
| `M-b` | Move backward one word |
| `C-l` | Clear screen, reprint line |
| `C-k` | Kill to end of line |
| `M-d` | Kill word forward |
| `M-DEL` | Kill word backward |
| `C-w` | Kill backward to whitespace |
| `C-y` | Yank last killed text |
| `M-y` | Rotate kill-ring and yank |
| `C-_` | Undo |
| `C-x C-r` | Re-read init file |

### Vi Mode

Enable with `set -o vi` or `set editing-mode vi` in `.inputrc`.

- Starts in **insertion mode** (you can type immediately)
- Press `ESC` to enter **command mode**
- Standard vi movement: `h`/`l` (char), `b`/`w` (word), `0`/`$` (line)
- `i`/`a` to re-enter insertion mode
- `dd` to kill line, `p` to yank

## Completion Commands

| Key | Command | Description |
|-----|---------|-------------|
| `TAB` | `complete` | Attempt completion on text before point |
| `M-?` | `possible-completions` | List possible completions |
| `M-*` | `insert-completions` | Insert all completions separated by space |
| `M-/` | `complete-filename` | Filename completion |
| `M-~` | `complete-username` | Username completion |
| `M-$` | `complete-variable` | Shell variable completion |
| `M-@` | `complete-hostname` | Hostname completion |
| `M-!` | `complete-command` | Command name completion |
| `M-TAB` | `dynamic-complete-history` | Complete from history |
| `M-{` | `complete-into-braces` | Insert completions in brace format |
| (unbind) | `menu-complete` | Replace word with first match; cycle on repeat |
| (unbind) | `menu-complete-backward` | Cycle backward through matches |

### Menu Completion

To make TAB cycle through completions instead of showing all:

```bash
# In .inputrc
TAB: menu-complete
"\e[Z": menu-complete-backward

# Or via bind
bind 'TAB: menu-complete'
bind '"\e[Z": menu-complete-backward'
```

## Readline Variables

### Completion Variables

| Variable | Default | Effect |
|----------|---------|--------|
| `completion-ignore-case` | `off` | Case-insensitive filename matching |
| `completion-map-case` | `off` | Treat `-` and `_` as equivalent (with `completion-ignore-case`) |
| `colored-stats` | `off` | Use LS_COLORS to indicate file types in completion lists |
| `visible-stats` | `off` | Append file-type indicator character to completions |
| `show-all-if-ambiguous` | `off` | List all completions immediately on first TAB |
| `show-all-if-unmodified` | `off` | List completions immediately if no unique prefix |
| `completion-prefix-display-length` | `0` | Ellipsize common prefix longer than N; hidden files use `___` |
| `colored-completion-prefix` | `off` | Display common prefix in different color (uses LS_COLORS) |
| `mark-symlinked-directories` | `off` | Append `/` to symlinks pointing to directories |
| `mark-directories` | `on` | Append `/` to completed directory names |
| `match-hidden-files` | `on` | Match dotfiles during completion |
| `completion-query-items` | `100` | Prompt before showing more than N completions |
| `completion-display-width` | `-1` | Screen columns for match display; `0` = one per line |
| `page-completions` | `on` | Use internal pager for completion lists |
| `print-completions-horizontally` | `off` | Sort horizontally instead of vertically |
| `skip-completed-text` | `off` | Don't duplicate text after cursor on mid-word completion |
| `menu-complete-display-prefix` | `off` | Show common prefix before cycling menu completions |
| `disable-completion` | `off` | Disable completion entirely |
| `expand-tilde` | `off` | Perform tilde expansion during completion |

### Display Variables

| Variable | Default | Effect |
|----------|---------|--------|
| `bell-style` | `audible` | `none`, `visible`, or `audible` |
| `horizontal-scroll-mode` | `off` | Scroll horizontally vs. line wrapping |
| `enable-bracketed-paste` | `on` | Handle pasted text as single string |
| `blink-matching-paren` | `off` | Briefly move cursor to matching opening paren |

### Editing Mode Variables

| Variable | Default | Effect |
|----------|---------|--------|
| `editing-mode` | `emacs` | `emacs` or `vi` |
| `keymap` | `emacs` | Current keymap |
| `show-mode-in-prompt` | `off` | Show editing mode indicator in prompt |
| `emacs-mode-string` | `@` | String for emacs mode indicator |
| `vi-ins-mode-string` | `(ins)` | String for vi insertion mode |
| `vi-cmd-mode-string` | `(cmd)` | String for vi command mode |

### Input/Output Variables

| Variable | Default | Effect |
|----------|---------|--------|
| `input-meta` | `off` | Enable eight-bit input (synonym: `meta-flag`) |
| `output-meta` | `off` | Display eighth-bit characters directly |
| `convert-meta` | `on` | Convert eighth-bit chars to meta-prefixed sequences |
| `enable-meta-key` | `on` | Attempt to enable meta modifier key |
| `enable-keypad` | `off` | Enable application keypad mode |

### History Variables

| Variable | Default | Effect |
|----------|---------|--------|
| `history-preserve-point` | `off` | Keep cursor position when navigating history |
| `history-size` | `HISTSIZE` | Maximum history entries (0 = disable, <0 = unlimited) |
| `revert-all-at-newline` | `off` | Undo all history changes before accept-line |
| `mark-modified-lines` | `off` | Display `*` before modified history lines |
| `isearch-terminators` | `ESC C-j` | Characters that terminate incremental search |
| `search-ignore-case` | `off` | Case-insensitive history search |

### Other Variables

| Variable | Default | Effect |
|----------|---------|--------|
| `comment-begin` | `"#"` | String for `insert-comment` command |
| `keyseq-timeout` | `500` | Milliseconds to wait for ambiguous key sequence |
| `bind-tty-special-chars` | `on` | Bind kernel terminal control characters |
| `echo-control-characters` | `on` | Echo characters for keyboard-generated signals |
| `active-region-start-color` | standout | Terminal escape before active region text |
| `active-region-end-color` | restore | Terminal escape to restore after active region |
| `enable-active-region` | `on` | Enable region highlighting |

## Keyboard Macros

| Key | Command |
|-----|---------|
| `C-x (` | `start-kbd-macro` — begin saving |
| `C-x )` | `end-kbd-macro` — stop saving |
| `C-x e` | `call-last-kbd-macro` — replay last macro |
| (unbound) | `print-last-kbd-macro` — print in inputrc format |

### Macro Examples

```bash
# In .inputrc
"\C-xp": "PATH=${PATH}\e\C-e\C-a\ef\C-f"   # Insert PATH expansion
"\C-x\"": "\"\"\C-b"                          # Insert "" and move between
"\C-o": "> output"                            # Insert literal text
```

## Active Region

Readline supports highlighting a region between point (cursor) and mark:

- `C-@` or `C-x C-x` — set mark at current position
- `enable-active-region` — controls whether highlighting is active
- `active-region-start-color` / `active-region-end-color` — terminal escape sequences for highlighting

## Word Boundaries and Tokenization

Readline determines word boundaries for:

1. **Movement commands** (`M-f`, `M-b`) — words are composed of alphanumeric characters and underscores
2. **Completion tokenization** — when completing, Readline identifies the "word being completed" by finding boundaries around the cursor position
3. **COMP_WORDBREAKS** — characters that split words for completion purposes (default: `" \t\n\"'@><=;|&(:`)

### Completion Type Detection

Readline checks what the word begins with to determine completion type:

| Prefix | Completion Type |
|--------|----------------|
| `$` | Variable completion |
| `~` | Username/home directory completion |
| `@` | Hostname completion |
| (other) | Command, filename, or programmable completion |

## Sample .inputrc

```
# Editing mode
set editing-mode emacs

# Completion behavior
set completion-ignore-case on
set show-all-if-ambiguous on
set visible-stats on
set colored-stats on
set colored-completion-prefix on
set completion-prefix-display-length 2
set bell-style visible
set mark-symlinked-directories on

# Pager
set page-completions on
set completion-query-items 200

# Key bindings
$if mode=emacs
"\M-[D": backward-char
"\M-[C": forward-char
$endif

$if Bash
"\C-xp": "PATH=${PATH}\e\C-e\C-a\ef\C-f"
$endif
```

## References

- [GNU Bash Manual: Readline Interaction](https://www.gnu.org/software/bash/manual/html_node/Readline-Interaction.html)
- [GNU Bash Manual: Readline Init File](https://www.gnu.org/software/bash/manual/html_node/Readline-Init-File.html)
- [GNU Bash Manual: Readline Init File Syntax](https://www.gnu.org/software/bash/manual/html_node/Readline-Init-File-Syntax.html)
- [GNU Bash Manual: Bindable Readline Commands](https://www.gnu.org/software/bash/manual/html_node/Bindable-Readline-Commands.html)
- [GNU Bash Manual: Commands For Completion](https://www.gnu.org/software/bash/manual/html_node/Commands-For-Completion.html)
- [GNU Bash Manual: Readline vi Mode](https://www.gnu.org/software/bash/manual/html_node/Readline-vi-Mode.html)
- [GNU Bash Manual: Sample Init File](https://www.gnu.org/software/bash/manual/html_node/Sample-Init-File.html)

## Related Skills

- **references/completion.md** — how Readline variables interact with programmable completion
- **references/quoting-expansion.md** — how Readline tokenizes and quotes during completion
