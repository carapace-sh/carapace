# Fish Editor and Interactive Features

In-depth reference for fish's interactive editing system — key bindings, editing modes, autosuggestions, the pager, syntax highlighting, abbreviations, kill ring, multiline editing, prompts, and the input processing pipeline.

## The Reader/Editor System

The `Reader` manages the interactive command line as the high-level coordinator:

1. **Input reading** via `InputEventQueue`
2. **Command line editing** with undo support
3. **History search** integration
4. **Screen rendering** coordination
5. **Completion triggering** and pager display

### Input Processing Pipeline

Raw terminal bytes are transformed into structured events through a multi-stage pipeline:

```
Terminal FD → select() → Raw Bytes → Escape Sequence Parser → KeyEvent → CharEvent → InputQueue → Reader
```

### Core Data Structures

| Type | Purpose |
|------|---------|
| `InputEventTrigger` | Result of `select()`: `Byte(u8)`, `Eof`, `Interrupted`, `UvarNotified`, `IOPortNotified`, `TimeoutElapsed` |
| `KeyEvent` | Keyboard input with modifiers: `key`, `shifted_codepoint`, `base_layout_codepoint` |
| `Key` | Basic key representation: `modifiers`, `codepoint` |
| `CharEvent` | Union of all input events: `Key`, `Readline`, `Command`, `Implicit`, `QueryResult` |
| `InputData` | Input processing state: `queue`, `paste_buffer`, `blocking_query` |

### Escape Sequence Parsing

Fish uses a **timed read mechanism** to distinguish between a standalone Escape key and the start of an escape sequence (e.g., arrow keys):

- `WAIT_ON_ESCAPE_MS`: 30ms default (configurable via `fish_escape_delay_ms`)
- `WAIT_ON_SEQUENCE_KEY_MS`: Infinite for multi-key sequences
- Handles "torn" escape sequences (interrupted/incomplete) via the `InputData` queue
- Supports modern protocols like **Kitty keyboard protocol**

### Key Canonicalization

- Bytes 0-31 → Ctrl+letter (e.g., `\x01` → Ctrl-a)
- Special chars get named keys: `\x0d` → Enter, `\x09` → Tab, `\x1b` → Escape
- Private Unicode area U+F500-U+F5FF for named keys like `Backspace`, `Up`, `F1`

## Key Binding Modes

### Emacs Mode (Default)

Standard Emacs-style key bindings. Switch with `fish_default_key_bindings`.

**Essential commands:**

| Key | Action |
|-----|--------|
| `Home` / `Ctrl-a` | Move to beginning of line |
| `End` / `Ctrl-e` | Move to end of line; at end accepts autosuggestion |
| `Ctrl-f` / `→` | Move forward one character; at end accepts autosuggestion |
| `Ctrl-b` / `←` | Move backward one character |
| `Alt-f` / `Alt-→` | Move forward one word; empty line = directory history; accepts first word of autosuggestion |
| `Alt-b` / `Alt-←` | Move backward one word; empty line = directory history |
| `Ctrl-l` | Clear screen, reprint line |
| `Ctrl-k` | Kill to end of line |
| `Alt-d` | Kill next word |
| `Ctrl-w` | Kill previous path component |
| `Ctrl-y` | Yank last killed text |
| `Alt-y` | Rotate kill-ring and yank |
| `Ctrl-_` / `Ctrl-z` | Undo most recent edit |
| `Alt-/` / `Ctrl-Shift-z` | Redo most recent undo |
| `Ctrl-t` | Transpose last two characters |
| `Alt-t` | Transpose last two words |
| `Alt-c` | Capitalize current word |
| `Alt-u` | Uppercase current word |
| `Ctrl-r` | Open history in pager with search |
| `Escape` / `Ctrl-g` | Cancel operation or undo unambiguous completion |

### Vi Mode

Enable with `fish_vi_key_bindings` or `set -g fish_key_bindings fish_vi_key_bindings`.

**Command mode** (normal mode, entered with `Escape`):

| Key | Action |
|-----|--------|
| `h` / `l` | Move left/right |
| `k` / `j` | History search previous/next |
| `i` / `I` | Enter insert mode at cursor/beginning |
| `a` / `A` | Enter insert mode after cursor/end |
| `o` / `O` | Insert new line below/above and enter insert mode |
| `0` | Move to beginning of line |
| `$` | Move to end of line |
| `dd` / `D` | Delete line/text and put in kill ring |
| `p` | Paste from kill ring |
| `u` | Undo |
| `Ctrl-r` | Redo |
| `[` / `]` | Search history for token under cursor |
| `/` | Open history in pager with search |
| `g,g` / `G` | Go to beginning/end of command line |
| `~` | Toggle case and move to next character |
| `v` | Enter visual mode |
| `:q` | Exit fish |

**Insert mode:**

| Key | Action |
|-----|--------|
| `Escape` | Enter command mode |
| `Backspace` | Delete character left |
| `Ctrl-n` | Accept autosuggestion |

**Visual mode:**

| Key | Action |
|-----|--------|
| `←` / `→` / `h` / `l` | Extend selection by character |
| `k` / `j` | Move up/down |
| `b` / `w` | Extend selection by word |
| `d` / `x` | Move selection to kill ring, enter command mode |
| `c` / `s` | Remove selection, enter insert mode |
| `X` | Move entire line to kill ring |
| `y` | Copy selection to kill ring |
| `~` | Toggle case on selection |
| `g,u` / `g,U` | Lowercase/uppercase selection |
| `"` / `*` / `y` | Copy selection to clipboard |

### Vi Cursor Configuration

```fish
set fish_cursor_default block              # normal/visual mode
set fish_cursor_insert line                # insert mode
set fish_cursor_replace_one underscore     # replace mode
set fish_cursor_external line              # external cursor
set fish_cursor_visual block               # visual mode
# Add "blink" for blinking cursors
set fish_cursor_default block blink
```

### Shared Bindings (Both Modes)

| Key | Action |
|-----|--------|
| `Tab` | Complete current token |
| `Shift-Tab` | Complete and start pager search |
| `Enter` | Execute or insert newline if incomplete |
| `Alt-Enter` | Insert newline at cursor |
| `Ctrl-←` / `Ctrl-→` | Move by word; accepts one word of autosuggestion |
| `Shift-←` / `Shift-→` | Move by word without stopping on punctuation |
| `↑` / `↓` | History search for previous/next command containing string |
| `Alt-↑` / `Alt-↓` | Search history for token under cursor |
| `Ctrl-c` | Interrupt/kill (SIGINT) |
| `Ctrl-d` | Delete right; exit if line empty |
| `Ctrl-u` | Kill from beginning to cursor |
| `Ctrl-x` | Copy to system clipboard |
| `Ctrl-v` | Paste from system clipboard |
| `Alt-h` / `F1` | Show manual page for current command |
| `Alt-l` | List current directory |
| `Alt-o` | Open file in pager or editor |
| `Alt-p` | Append `&\| less;` to command |
| `Alt-w` | Print short description of command |
| `Alt-e` / `Alt-v` | Edit command line in external editor |
| `Alt-s` | Prepend `sudo` (or `doas`/`please`/`run0`) |
| `Ctrl-Space` | Insert space without expanding abbreviation |

### Custom Key Bindings

```fish
# In config.fish or fish_user_key_bindings function
function fish_user_key_bindings
    bind ctrl-c cancel-commandline
    bind alt-x 'commandline -r sudo; commandline -a " "(commandline -p)'
end

# Remove a custom binding
bind --erase ctrl-c

# Discover key names
fish_key_reader
# Press a key to see its name, e.g.: bind alt-right 'do something'
```

### Escape Delay

```fish
set -g fish_escape_delay_ms 100  # Higher value = more time to distinguish escape from alt
```

## Autosuggestions

Fish suggests commands as you type based on command history, completions, and valid file paths.

### Behavior

- Suggestion appears after cursor in muted gray color (configurable via `fish_color_autosuggestion`)
- Press `→` or `Ctrl-f` to accept the full suggestion
- Press `Alt-→` or `Alt-f` to accept just the first suggested word
- Ignoring the suggestion means it won't execute
- Suggestions are computed asynchronously — path validation is debounced

### Disabling

```fish
set -g fish_autosuggestion_enabled 0
```

### How Autosuggestions Work Internally

1. After each keystroke, the reader computes an autosuggestion
2. History is searched for matching commands
3. If a history match is found, it's offered as the suggestion
4. Path validation (for file/directory suggestions) is offloaded to a thread pool
5. The `sort_and_prioritize()` function ranks candidates — for autosuggestions, it prefers case matches and penalizes duplicates and tilde-suffixed files

## Syntax Highlighting

Fish interprets the command line as you type and provides visual feedback. **Errors are marked red** by default.

### Detected Errors

- Non-existing commands
- Reading from or appending to non-existing files
- Incorrect use of output redirects
- Mismatched parentheses

### Syntax Highlighting Variables

| Variable | Meaning |
|----------|---------|
| `fish_color_normal` | Default color |
| `fish_color_command` | Commands like `echo` |
| `fish_color_keyword` | Keywords like `if` |
| `fish_color_quote` | Quoted text like `"abc"` |
| `fish_color_redirection` | I/O redirections |
| `fish_color_end` | Process separators like `;` and `&` |
| `fish_color_error` | Syntax errors |
| `fish_color_param` | Ordinary parameters |
| `fish_color_valid_path` | Valid filenames |
| `fish_color_option` | Options starting with `-` |
| `fish_color_comment` | Comments like `# important` |
| `fish_color_operator` | Parameter expansion like `*` and `~` |
| `fish_color_escape` | Character escapes like `\n` |
| `fish_color_autosuggestion` | Autosuggestions |
| `fish_color_search_match` | History search matches |

### Theme Management

```fish
fish_config theme choose none       # disable nearly all coloring
fish_config theme choose default   # restore default theme
fish_config theme show             # see all themes in terminal
fish_config theme save             # save current theme
fish_config theme dump             # dump current theme as variables
```

### Theme Synchronization Across Sessions

```fish
function apply-my-theme --on-variable=my_theme
    fish_config theme choose $my_theme
end
set -U my_theme lava  # All fish sessions update
```

## Abbreviations

Abbreviations expand a short token into a longer command when `Space` or `Enter` is pressed.

### Basic Abbreviations

```fish
abbr -a gco git checkout
abbr -a gs git status
abbr -a ll 'ls -la'
```

After typing `gco` and pressing `Space` or `Enter`, it expands to `git checkout`. The expanded form is what gets stored in history.

- Use `Ctrl-Space` to insert a literal abbreviation without expanding
- Abbreviations are visible before expansion — you can see and edit the actual command

### Regex-Based Abbreviations

```fish
function multicd
    echo cd (string repeat -n (math (string length -- $argv[1]) - 1) ../)
end
abbr --add dotdot --regex '^\.\.+$' --function multicd
```

Now `..` → `cd ../`, `...` → `cd ../../`, `....` → `cd ../../../`, etc.

### Abbreviation Management

```fish
abbr -a name expansion     # Add abbreviation
abbr -e name               # Erase abbreviation
abbr -l                    # List abbreviations
abbr -s                    # Show abbreviations as commands
```

### Advantages Over Aliases

- Can see the actual command before using it
- Can add to or change it before execution
- Actual command is stored in history (not the abbreviation)
- Regex abbreviations enable dynamic expansion

## The Pager

### Overview

The pager displays completions in a scrollable table when multiple matches exist. It is managed by the `Pager` struct.

### Pager Structure

```rust
pub struct Pager {
    pub available_term_width: usize,
    pub available_term_height: usize,
    pub selected_completion_idx: Option<usize>,
    pub fully_disclosed: bool,           // Show all completions?
    pub search_field_shown: bool,
    pub search_field_line: EditableLine, // User's filter text
    completion_infos: Vec<PagerComp>,     // Filtered completions
    unfiltered_completion_infos: Vec<PagerComp>,
    prefix: Cow<'static, wstr>,          // Shared prefix (e.g., directory path)
}
```

### Navigation

| Key | Action |
|-----|--------|
| Arrow keys, Page Up/Down | Navigate completions |
| Tab / Shift+Tab | Move selection forward/backward |
| Ctrl+S or `/` (vi mode) | Open search menu to filter the list |
| Enter | Accept selected completion |
| Escape | Close pager |

### Layout

- Maximum 6 columns (`PAGER_MAX_COLS`)
- Minimum 16 chars width, 4 rows height to display
- Initially shows 4 rows or half terminal height (undisclosed mode)
- Press Tab again to fully disclose all completions
- Prefixes shared by all completions may be suppressed to save space

### Pager Color Variables

| Variable | Meaning |
|----------|---------|
| `fish_pager_color_progress` | Progress bar at bottom left |
| `fish_pager_color_background` | Background of each row |
| `fish_pager_color_prefix` | The prefix string to complete |
| `fish_pager_color_completion` | The proposed completion |
| `fish_pager_color_description` | Completion description |
| `fish_pager_color_selected_background` | Selected completion background |
| `fish_pager_color_selected_completion` | Selected completion text |
| `fish_pager_color_selected_description` | Selected completion description |
| `fish_pager_color_secondary_background` | Alternating row background |
| `fish_pager_color_secondary_completion` | Alternating row completion text |
| `fish_pager_color_secondary_description` | Alternating row description |

### Fuzzy Filtering

When the search field is active (Ctrl+S or `/`):

- Matches against both description and completion strings
- Uses `string_fuzzy_match_string()` for matching
- Lower rank = better match
- Case sensitivity: case-insensitive unless uppercase character in search string

## Screen Rendering

### Dual-Buffer Architecture

The `Screen` class maintains two `ScreenData` instances:

- **`desired`**: The intended screen state after rendering
- **`actual`**: The last-known real terminal state

This enables **differential updates** — only differences between desired and actual states are sent to the terminal.

### Rendering Pipeline

1. **Layout Computation** — determine text wrapping across screen width
2. **Desired State Construction** — build ideal state in `desired`
3. **Differential Update** — send only changes via `update()`
4. **State Synchronization** — update `actual` to match `desired`

### Optimization Techniques

- **Soft Wrap Exploitation** — write consecutive lines without explicit cursor repositioning
- **Line-by-Line Comparison** — check each line against actual state
- **Intra-Line Differential** — find first difference and write from there
- **Selective Clearing** — use `ClearToEndOfLine` only when necessary

### Color Transitions

The `Outputter` provides a unified interface for terminal writes. Color transitions are optimized by comparing new `TextFace` against the last written face — only necessary attribute changes (Bold, Dim, Italic, Underline) are written.

## Kill Ring

Fish uses an Emacs-style kill ring:

| Key | Action |
|-----|--------|
| `Ctrl-k` | Kill from cursor to end of line |
| `Ctrl-y` | Yank (paste) latest kill |
| `Alt-y` | Rotate kill-ring and yank previous entry |

### System Clipboard

| Key | Action |
|-----|--------|
| `Ctrl-x` | Copy to system clipboard |
| `Ctrl-v` | Paste from system clipboard |

### Bracketed Paste Mode

Fish enables bracketed paste mode for automatic detection of pasted content. Pasted text is inserted as a single unit rather than being interpreted as individual keystrokes.

### Selection Mode

```fish
set fish_cursor_selection_mode inclusive   # include character under cursor
set fish_cursor_selection_mode exclusive   # exclude character under cursor (default)
```

Kill ring entries are stored in the `fish_killring` variable.

## Multiline Editing

Three ways to span multiple lines:

1. **Unclosed blocks** — press `Enter` with unclosed `for`, `begin`, `if`, etc.
2. **Alt-Enter** — insert a newline at cursor position
3. **Backslash** — insert `\` before `Enter` for line continuation

## Searchable Command History

### Behavior

- Commands stored in history after execution
- Duplicate entries removed automatically
- `↑`/`↓` search forwards/backwards through history
- If command line not empty, only matching commands shown
- `Alt-↑`/`Alt-↓` searches for elements matching token under cursor
- `Ctrl-r` opens a pager for history searching

### Case Sensitivity

Search is case-insensitive unless an uppercase character appears in the search string.

### Private Mode

```fish
fish --private        # or fish -P
set fish_private_mode 1   # enable internally
```

Disables writing command history to disk. Useful for screencasts or handling sensitive information.

### History File Location

- Default: `~/.local/share/fish/fish_history`
- Or `$XDG_DATA_HOME/fish/fish_history`
- Session-specific files via `$fish_history` variable
- Prefix a command with space to prevent it from being stored in history

## Programmable Prompt

Fish runs `fish_prompt` and `fish_right_prompt` functions to display prompts.

### Prompt Functions

| Function | Purpose |
|----------|---------|
| `fish_prompt` | Left prompt |
| `fish_right_prompt` | Right prompt |
| `fish_mode_prompt` | Vi mode indicator (above left prompt) |
| `fish_title` | Terminal title (before and after commands) |
| `fish_greeting` | Greeting message on interactive startup |
| `fish_transient_prompt` | Redraw prompt before command execution |

### Managing Prompts

```fish
fish_config prompt show           # show all prompts in terminal
fish_config prompt choose disco   # temporarily select a prompt
fish_config prompt save           # save selection
funced fish_prompt                # edit the prompt function
funcsave fish_prompt              # save after editing
```

### Example: Custom Prompt

```fish
function fish_prompt
    set -l last_status $status
    set -l color (set_color $fish_color_cwd)
    set -l prompt_symbol '$'

    # Print working directory
    echo -n -s $color (prompt_pwd) (set_color normal) " $prompt_symbol "

    # Show exit status if non-zero
    if test $last_status -ne 0
        echo -n -s (set_color $fish_color_error) "[$last_status]" (set_color normal) " "
    end
end
```

### Transient Prompt

The `fish_transient_prompt` function allows redrawing the prompt before command execution. This is useful for showing a simplified prompt in scrollback:

```fish
function fish_transient_prompt
    echo -n '❯ '
end
```

## Directory Navigation

### Directory History

Fish keeps track of visited directories in `dirprev` and `dirnext` variables.

| Command | Description |
|---------|-------------|
| `dirh` | Print the directory history |
| `cdh` | Display prompt to navigate history |
| `prevd` | Move backward (bound to `Alt-←`) |
| `nextd` | Move forward (bound to `Alt-→`) |

### Directory Stack

| Command | Description |
|---------|-------------|
| `dirs` | Print the stack |
| `pushd <dir>` | Push directory onto stack and cd to it |
| `popd` | Pop directory from stack and cd to it |

## Configurable Greeting

Fish runs `fish_greeting` function on interactive startup.

```fish
set -U fish_greeting              # disable greeting
set -g fish_greeting 'Hey, stranger!'
function fish_greeting
    random choice "Hello!" "Hi" "G'day" "Howdy"
end
```

## References

- [Interactive use — fish-shell docs](https://fishshell.com/docs/current/interactive.html)
- [bind builtin — fish-shell docs](https://fishshell.com/docs/current/cmds/bind.html)
- [abbr builtin — fish-shell docs](https://fishshell.com/docs/current/cmds/abbr.html)
- [Input Processing and Terminal I/O — DeepWiki](https://deepwiki.com/fish-shell/fish-shell/3-input-processing-and-terminal-io)
- [Screen Rendering and Display — DeepWiki](https://deepwiki.com/fish-shell/fish-shell/3.3-screen-rendering-and-display)
- [Fish shell source: src/reader/reader.rs](https://github.com/fish-shell/fish-shell/blob/main/src/reader/reader.rs)
- [Fish shell source: src/pager.rs](https://github.com/fish-shell/fish-shell/blob/main/src/pager.rs)

## Related Skills

- **carapace-dev** → `references/shell-fish.md` — carapace-specific fish integration
- **fish** → `references/completion.md` — fish completion system details
