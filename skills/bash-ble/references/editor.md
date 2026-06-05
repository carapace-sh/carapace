# BLE Editor Architecture

In-depth reference for ble.sh's editor architecture — how it replaces Readline, the widget system, keymaps, key binding, line editing, and terminal rendering.

## How ble.sh Replaces Readline

ble.sh **replaces** GNU Readline entirely rather than augmenting it. Key mechanisms:

1. **Input capture via `bind -x`** — Routes keyboard input to `_ble_decode_hook` instead of Readline
2. **Complete widget implementation** — All editing operations implemented as Bash functions (not delegated to Readline)
3. **Custom rendering** — Terminal output goes directly via the canvas system, bypassing Readline's display management
4. **History integration** — Wraps Bash's builtin history with enhancements for asynchronous loading, multiline handling
5. **Hook system** — `blehook` allows customization at lifecycle points without Readline hooks
6. **Readline compatibility layer** — Maps Readline function names to internal widgets for backward compatibility with `~/.inputrc` configurations

## Core Editor State

| Variable | Purpose |
|----------|---------|
| `_ble_edit_str` | Current command line content |
| `_ble_edit_ind` | Cursor position index |
| `_ble_edit_mark` | Selection mark position |
| `_ble_edit_overwrite_mode` | Insert/overwrite mode flag |
| `_ble_edit_arg` | Numeric prefix argument |
| `_ble_edit_kill_ring` | Circular buffer of killed/copied text |
| `_ble_edit_kill_ring_index` | Current position for yank |
| `_ble_edit_undo` | Stack of undo records |
| `_ble_edit_undo_index` | Current position in undo history |

### Buffer Manipulation

Buffer modifications go through controlled operations:

- `ble-edit/content/replace` — Core replacement with undo recording
- `ble/widget/.insert-string` — Insert with overwrite mode handling
- `ble/widget/.delete-range` — Delete range, optionally to kill ring

## Widget System

Every user input is handled by executing a **widget**. Each widget is a shell function `ble/widget/WidgetName`.

### Widget Categories

| Category | Examples |
|----------|---------|
| Character | `self-insert`, `delete-char`, `delete-char-backward` |
| Word movement | `forward-cword`, `backward-cword`, `kill-word` |
| Line movement | `beginning-of-line`, `end-of-line` |
| Kill/Yank | `kill-region`, `yank`, `yank-pop`, `kill-line` |
| History | `history-search-backward`, `history-search-forward` |
| Completion | `complete`, `menu-complete`, `dabbrev-expand` |
| Sabbrev | `sabbrev-expand`, `magic-space` |
| Undo | `undo`, `redo` |
| Mark | `set-mark`, `exchange-point-and-mark` |

### Widget Dispatch Flow

1. Key pressed → decode to key code
2. Check keymaps in stack for binding
3. Execute bound widget
4. Before/after hooks wrap execution

### Custom Widgets

```bash
function ble/widget/my/example1 {
  ble/widget/beginning-of-logical-line
  ble/widget/insert-string 'echo $('
  ble/widget/end-of-logical-line
  ble/widget/insert-string ')'
}
ble-bind -f C-t my/example1
```

### Special Binding Names

| Name | Purpose |
|------|---------|
| `__default__` | Fallback handler for unbound keys |
| `__defchar__` | Printable character handler |
| `__batch_char__` | Batch character input handler |
| `__before_widget__` | Hook before widget execution |
| `__after_widget__` | Hook after widget execution |
| `__attach__` | Keymap activation hook |
| `__detach__` | Keymap deactivation hook |

## Keymaps

### Available Keymaps

| Keymap | Description |
|--------|-------------|
| `emacs` | Emacs editing mode (default) |
| `vi_imap` | Vim Insert mode |
| `vi_nmap` | Vim Normal mode |
| `vi_omap` | Vim Operator-pending mode |
| `vi_xmap` | Vim Visual mode |
| `vi_smap` | Vim Select mode |
| `vi_cmap` | Vim Command-line mode |
| `vi_digraph` | Vim digraph input |
| `safe` | Emergency base keymap |
| `isearch` | Incremental search mode |
| `nsearch` | Non-incremental search mode |
| `read` | For `read -e` builtins |
| `auto_complete` | Auto-completion mode |
| `menu_complete` | Menu completion mode |

### Keymap Stack

ble.sh maintains a stack (`_ble_decode_keymap_stack`) allowing multiple keymaps to be active simultaneously. When a key is pressed, ble.sh searches from top of stack downward.

## Key Binding

### ble-bind Command

```bash
ble-bind [-m keymap] -f key widget     # Bind key to widget
ble-bind [-m keymap] -c key 'cmd'      # Bind to shell command
ble-bind [-m keymap] -s key 'macro'     # Bind to macro (key sequence)
ble-bind [-m keymap] -p                 # Print bindings
ble-bind [-m keymap] -r key             # Remove binding
ble-bind [-m keymap] -u key             # Unbind (undefined)
ble-bind [-m keymap] -x key 'code'     # Bind to editing function
ble-bind -T timeout                     # Set timeout for sequences
ble-bind -k byte_seq                    # Register byte sequences
ble-bind --csi seq                      # Register CSI sequences
```

### Key Notation

| Format | Meaning |
|--------|---------|
| `C-x` | Control + X |
| `M-x` | Meta/Alt + X |
| `S-x` | Shift + X |
| `A-x` | AltGr + X |
| `s-x` | Super + X |
| `H-x` | Hyper + X |
| `f1`, `home`, `up` | Special keys |
| `\xHH` | Hex byte |
| `C-x C-c` | Key sequence (chord) |

### Binding Types

| Type | Flag | Description |
|------|------|-------------|
| Direct widget | `1:widget` | Immediately dispatches to widget |
| Prefix | `_:widget` | May continue with more keys |
| Timeout prefix | `_TIMEOUT:widget` | Prefix with timeout |

### Examples

```bash
# Bind key to widget
ble-bind -f 'C-x h' 'insert-string "Hello, world!"'

# Bind in specific keymap
ble-bind -m vi_imap -f 'C-a' 'beginning-of-line'

# Bind to shell command
ble-bind -c 'M-c' 'my-command'

# Bind to macro
ble-bind -s 'C-x e' 'C-a C-k echo Hello RET'

# Multiline RET
ble-bind -m emacs -f 'C-m' 'accept-line'
ble-bind -m vi_imap -f 'C-m' 'accept-line'
```

## Input Decoding System

### Processing Pipeline

```
Raw bytes → Byte Processing → Character Decoding → Key Processing → Widget Dispatch
```

1. **Byte Processing** — `_ble_decode_hook` receives bytes via `bind -x`
2. **Character Decoding** — Converts bytes to characters, handles CSI sequences
3. **Key Processing** — Applies modifiers, looks up in keymap
4. **Widget Dispatch** — Executes bound widget function

### Key Code Format

32-bit integers with bit layout:

| Bits | Content |
|------|---------|
| 0-20 | Character code |
| 21 | Unused |
| 22-28 | Modifier flags (Meta, Ctrl, Shift, Hyper, Super, AltGr) |
| 29 | Macro flag |
| 30 | Error flag |

## Editing Modes

### Emacs Mode (Default)

Standard Emacs-style key bindings. Control keys operate on characters, Meta keys on words.

**Essential commands:**

| Key | Widget |
|-----|--------|
| `C-a` | `beginning-of-line` |
| `C-e` | `end-of-line` |
| `C-f` | `forward-char` |
| `C-b` | `backward-char` |
| `M-f` | `forward-cword` |
| `M-b` | `backward-cword` |
| `C-l` | `clear-screen` |
| `C-k` | `kill-line` |
| `M-d` | `kill-word` |
| `M-DEL` | `backward-kill-word` |
| `C-w` | `backward-kill-word` |
| `C-y` | `yank` |
| `M-y` | `yank-pop` |
| `C-_` | `undo` |
| `C-x C-r` | Re-read init file |

### Vi Mode

Enable with `ble-bind -m vi_imap` or set in blerc:

```bash
bleopt edit_vim_mode=vi
```

- Starts in **insertion mode** (you can type immediately)
- Press `ESC` to enter **normal mode**
- Standard vi movement: `h`/`l` (char), `b`/`w` (word), `0`/`$` (line)
- `i`/`a` to re-enter insertion mode
- `dd` to kill line, `p` to yank
- Full support for operators, text objects, registers, keyboard macros, marks
- `vim-surround` available as an option

## Command Execution Flow

1. User presses Enter → `ble/widget/accept-line`
2. Command validated against syntax
3. `ble-edit/accept-line` adds to history
4. `ble/history/add` updates history
5. `blehook PREEXEC` fires
6. Command executes
7. `blehook POSTEXEC` fires

## Terminal Rendering

### Rendering Pipeline

1. `ble-edit/content/reset` marks buffer dirty
2. `ble/textarea#render` calculates layout
3. `ble/canvas/trace` generates ANSI sequences
4. `ble/canvas/flush.draw` outputs to terminal
5. Cursor positioned via `ble/canvas/put-move-x/y`

### Canvas Configuration

| Option | Default | Description |
|--------|---------|-------------|
| `tab_width` | 8 | Displayed width of tabs |
| `char_width_version` | auto | Unicode version for char width |
| `emoji_width` | 2 | Emoji display width |
| `emoji_version` | 14.0 | Unicode Emoji version |
| `grapheme_cluster` | extended | Grapheme cluster type (extended/legacy) |
| `canvas_winch_action` | `redraw-safe` | SIGWINCH behavior |

### SIGWINCH Actions

| Value | Description |
|-------|-------------|
| `redraw-safe` | Redraw if safe (default) |
| `redraw-prev` | Redraw previous content |
| `redraw-here` | Redraw at current position |
| `clear` | Clear and redraw |

### Terminal Requirements

- Terminal must support SGR sequences for colors
- For `tmux`, use `tmux-256color` or `xterm-256color` terminal type
- `term_index_colors` — Number of index colors (auto/256/88/0)
- `term_true_colors` — 24-bit color support (`semicolon`, `colon`, or empty)

## Readline Compatibility

ble.sh provides a Readline compatibility layer:

- Maps Readline function names to internal widgets
- Can load `~/.inputrc` settings via `ble-import config/readline`
- `bind -x` commands are intercepted and redirected to ble.sh's decode system
- `READLINE_LINE`, `READLINE_POINT`, `READLINE_MARK` are supported

```bash
# Import Readline configuration
ble-import config/readline
```

## References

- [ble.sh Key Binding Manual](https://github.com/akinomyoga/ble.sh/wiki/Manual-%C2%A73-Key-Binding)
- [Widget Design (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/2.4-keymap-and-widget-system)
- [Input Decoding (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/2.2-input-decoding-system)
- [Line Editor Core (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/2.3-line-editor-core)
- [Core Architecture (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/2-core-architecture)
