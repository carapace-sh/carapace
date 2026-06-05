# Tcsh Command-Line Editor

In-depth reference for tcsh's command-line editor — key bindings, editing modes, completion editor commands, prompt configuration, and terminal handling.

## Overview

tcsh includes a built-in command-line editor that supports emacs-style and vi-style key bindings. The editor is active only when the `edit` shell variable is set (default in interactive shells).

The editor is a tcsh extension — the original csh has no command-line editing. The "t" in tcsh comes from TENEX, the operating system that inspired the command-completion feature.

## Editing Modes

### Emacs Mode (Default)

```tcsh
bindkey -e
```

Binds all keys to emacs-style bindings and unsets the `vimode` variable. This is the default for interactive shells.

### Vi Mode

```tcsh
bindkey -v
```

Binds all keys to vi-style bindings and sets the `vimode` variable. In vi mode, the editor operates in command, delete, insert, or replace states.

### Mode Differences

| Aspect | Emacs Mode | Vi Mode |
|--------|-----------|---------|
| Default | Yes | No |
| `vimode` variable | Unset | Set |
| `wordchars` default | `*?_-.[]~=` | `_` |
| Word boundary | `wordchars` vs other chars | `wordchars` vs whitespace vs other chars |
| Insert mode | Always in insert mode | Must enter insert mode explicitly |
| Command mode | All keys are editing commands | Escape enters command mode |

## The `bindkey` Builtin

```
bindkey [-l | -d | -e | -v | -u]
bindkey [-a] [-b] [-k] [-r] [--] key
bindkey [-a] [-b] [-k] [-c | -s] [--] key command
```

### Forms

| Form | Behavior |
|------|----------|
| No arguments | List all bound keys and their editor commands |
| `-l` | List all editor commands with short descriptions |
| `-d` | Bind all keys to standard default bindings |
| `-e` | Bind all keys to emacs-style; unset `vimode` |
| `-v` | Bind all keys to vi-style; set `vimode` |
| `-u` | Print usage message |
| `key` | Show the editor command bound to `key` |
| `key command` | Bind `command` to `key` |

### Options

| Option | Description |
|--------|-------------|
| `-a` | List or change bindings in the alternative key map (vi command mode) |
| `-b` | Interpret `key` as a control character (`^X`), meta character (`M-X`), function key (`F-string`), or extended prefix key (`X-`) |
| `-c` | Interpret `command` as a builtin or external command (not an editor command) |
| `-k` | Interpret `key` as a symbolic arrow key name (down, up, left, right) |
| `-r` | Remove `key`'s binding completely (does not bind to self-insert) |
| `-s` | Interpret `command` as a literal string treated as terminal input |
| `--` | Force end of option processing |

### Key Format

| Format | Example | Meaning |
|--------|---------|---------|
| Literal character | `a` | The character `a` |
| Control | `^A` | Ctrl-A |
| Meta | `M-a` | Meta-a (Alt-a) |
| Function key | `F1` | F1 key |
| Extended prefix | `X-A` | Extended key sequence |
| Backslash escapes | `\t`, `\n`, `\e` | Tab, newline, escape |

### Binding a Shell Command

The `-c` flag binds a key to execute a shell command:

```tcsh
# Bind F1 to run ls
bindkey -c F1 ls

# Bind Ctrl-X Ctrl-L to run a pipeline
bindkey -c '^X^L' 'ls -la | less'
```

### Binding a Literal String

The `-s` flag binds a key to insert a literal string:

```tcsh
# Bind Ctrl-X h to insert "hello world"
bindkey -s '^Xh' 'hello world'
```

## Editor Commands

### Movement

| Command | Default Binding | Description |
|---------|-----------------|-------------|
| `backward-char` | `^B`, left | Move back one character |
| `forward-char` | `^F`, right | Move forward one character |
| `backward-word` | `M-b`, `M-B` | Move to beginning of current word |
| `forward-word` | `M-f`, `M-F` | Move forward to end of current word |
| `beginning-of-line` | `^A`, home | Move to beginning of line |
| `end-of-line` | `^E`, end | Move to end of line |

### Deletion and Kill

| Command | Default Binding | Description |
|---------|-----------------|-------------|
| `delete-char` | Not bound | Delete character under cursor |
| `backward-delete-char` | `^H`, `^?` | Delete character before cursor |
| `delete-word` | `M-d`, `M-D` | Cut from cursor to end of current word |
| `backward-delete-word` | `M-^H`, `M-^?` | Cut from beginning of word to cursor |
| `kill-line` | `^K` | Cut from cursor to end of line |
| `backward-kill-line` | `^U` | Cut from beginning of line to cursor |
| `kill-region` | `^W` | Cut region between mark and cursor |
| `yank` | `^Y` | Paste last killed text |
| `yank-pop` | `M-y` | Rotate kill ring and paste |

### History Navigation

| Command | Default Binding | Description |
|---------|-----------------|-------------|
| `up-history` | up, `^P` | Copy previous history entry into buffer |
| `down-history` | down, `^N` | Step down through history |
| `history-search-backward` | `M-p`, `M-P` | Search history backward (prefix match) |
| `history-search-forward` | `M-n`, `M-N` | Search history forward (prefix match) |
| `i-search-back` | Not bound | Incremental search backward |
| `i-search-fwd` | Not bound | Incremental search forward |
| `vi-search-back` | `?` (vi mode) | Vi-style search backward |
| `vi-search-fwd` | `/` (vi mode) | Vi-style search forward |
| `insert-last-word` | `M-_` | Insert last word of previous input line |
| `dabbrev-expand` | `M-/` | Dynamic abbreviation — expand to most recent preceding match |

### History Expansion

| Command | Default Binding | Description |
|---------|-----------------|-------------|
| `expand-history` | `M-space` | Expand history substitutions in current word |
| `expand-line` | Not bound | Expand history substitutions in all words |
| `magic-space` | Not bound | Expand history and insert space |
| `toggle-literal-history` | `M-r`, `M-R` | Toggle between expanded and literal history form |

### Completion Editor Commands

| Command | Default Binding | Description |
|---------|-----------------|-------------|
| `complete-word` | Tab | Complete the current word |
| `complete-word-fwd` | Not bound | Cycle to next completion candidate |
| `complete-word-back` | Not bound | Cycle to previous completion candidate |
| `complete-word-raw` | `^X-Tab` | Complete word, ignoring user-defined completions |
| `list-choices` | `M-^D` | List completion possibilities |
| `list-choices-raw` | `^X-^D` | List choices, ignoring user-defined completions |
| `delete-char-or-list-or-eof` | `^D` | Delete char / list completions / EOF (context-dependent) |
| `list-or-eof` | Not bound | List choices or EOF on empty line |
| `delete-char-or-list` | Not bound | Delete char or list completions |

### Spelling Correction

| Command | Default Binding | Description |
|---------|-----------------|-------------|
| `spell-word` | `M-s`, `M-S` | Correct spelling of current word |
| `spell-line` | `M-$` | Correct spelling of each word in buffer |

### Expansion and Utility

| Command | Default Binding | Description |
|---------|-----------------|-------------|
| `expand-glob` | `^X-*` | Expand glob pattern to the left of cursor |
| `expand-variables` | `^X-$` | Expand variable to the left of cursor |
| `normalize-command` | `^X-?` | Replace command word with full path |
| `normalize-path` | `^X-n`, `^X-N` | Expand path through symlinks |
| `which-command` | `M-?` | Run `which` on the first word |
| `run-help` | `M-h`, `M-H` | Search for documentation on first word |
| `run-fg-editor` | `M-^Z` | Save line, resume stopped `+editor` job |

### Case Manipulation

| Command | Default Binding | Description |
|---------|-----------------|-------------|
| `capitalize-word` | `M-c`, `M-C` | Capitalize from cursor to end of word |
| `downcase-word` | `M-l`, `M-L` | Lowercase from cursor to end of word |
| `upcase-word` | `M-u`, `M-U` | Uppercase from cursor to end of word |

### Word and Region

| Command | Default Binding | Description |
|---------|-----------------|-------------|
| `copy-prev-word` | `M-^_` | Copy previous word into buffer |
| `exchange-point-and-mark` | `^X^X` | Swap cursor and mark positions |
| `set-mark-command` | `^@`, `^Space` | Set mark at cursor position |

## Arrow Key Bindings

Arrow keys are always bound (regardless of mode), using sequences from the `TERMCAP` environment variable:

| Arrow | Command |
|-------|---------|
| up | `up-history` |
| down | `down-history` |
| left | `backward-char` |
| right | `forward-char` |

ANSI/VT100 escape sequences for arrow keys are also always bound.

## The `wordchars` Variable

The `wordchars` shell variable defines non-alphanumeric characters that are considered part of a word by editor commands like `forward-word`, `backward-word`, `delete-word`, etc.

| Mode | Default `wordchars` |
|------|-------------------|
| Emacs (`vimode` unset) | `*?_-.[]~=` |
| Vi (`vimode` set) | `_` |

### How It Affects Editing

Characters in `wordchars` are treated as word constituents:

- `forward-word` moves past sequences of alphanumeric + `wordchars` characters
- `backward-word` moves back to the start of the current word
- `delete-word` cuts from cursor to end of the current word

Characters **not** in `wordchars` that are also non-alphanumeric act as word boundaries for the editor.

### wordchars vs Shell Word Boundaries

The editor's notion of a "word" differs from the shell's:

- **Editor**: Words are delimited by non-alphanumeric characters not in `wordchars`
- **Shell**: Words are delimited by whitespace and metacharacters (`&`, `|`, `;`, `<`, `>`, `(`, `)`)

This means the editor may treat `foo-bar` as two words (if `-` is not in `wordchars`), while the shell treats it as one word.

## Prompt Configuration

### `prompt` — Main Prompt

Default: `%# ` in interactive shells.

| Sequence | Expansion |
|----------|-----------|
| `%/` | Current working directory |
| `%~` | CWD with `~` substitution |
| `%c[[0]n]`, `%.[[0]n]` | Trailing n directory components |
| `%h`, `%!`, `!` | Current history event number |
| `%M` | Full hostname |
| `%m` | Short hostname (up to first `.`) |
| `%t`, `%@` | Time of day (12-hour AM/PM) |
| `%T` | Time of day (24-hour) |
| `%p` | Precise time with seconds (12-hour) |
| `%P` | Precise time with seconds (24-hour) |
| `%n` | Username |
| `%j` | Number of jobs |
| `%d` | Weekday name |
| `%D` | Weekday number (1=Monday) |
| `%w` | Month name |
| `%W` | Month number |
| `%y` | Year (2-digit) |
| `%Y` | Year (4-digit) |
| `%S` | Start standout mode (reverse video) |
| `%s` | Stop standout mode |
| `%B` | Start boldface |
| `%b` | Stop boldface |
| `%U` | Start underline |
| `%u` | Stop underline |
| `%L` | Clear to end of line |
| `%#` | `>` (normal user) or `#` (superuser) |
| `%?` | Return code of previous command |
| `%R` | Parser status (in prompt2), corrected string (in prompt3) |
| `%{...%}` | Literal escape sequence (terminal attributes) |
| `%%` | Literal `%` |

### `prompt2` — Continuation Prompt

Used for `while`/`foreach` loops and lines ending with `\`. Default: `%R? `

### `prompt3` — Correction Prompt

Used for spelling correction confirmation. Default: `CORRECT>%R (y|n|e|a)? `

### `rprompt` — Right-Side Prompt

Printed on the right side of the screen, after the command input. Recognizes the same format sequences as `prompt`. Automatically hides when command input is long enough to overlap. Only appears if everything fits on the first line.

```tcsh
set rprompt = "%~"
```

### Prompt Examples

```tcsh
# Simple prompt with hostname and directory
set prompt = "%m [%h] %B[%@]%b [%/] you rang? "

# Right-side prompt with directory
set rprompt = "%~"

# Git-style prompt with return code
set prompt = "%m:%/%# "
```

## Terminal Handling

### TERMCAP Environment Variable

tcsh uses the `TERMCAP` environment variable to determine terminal capabilities, including arrow key sequences. The shell always binds arrow keys based on this variable.

### Terminal Capability Commands

| Command | Description |
|---------|-------------|
| `telltc` | List all terminal capabilities |
| `settc cap value` | Set a terminal capability |
| `echotc cap [args]` | Output terminal capability string |

### Color Support

The `color` shell variable enables color display for the `ls-F` builtin:

```tcsh
set color
```

When set, `ls-F` passes `--color=auto` to `ls` for colored file type indicators.

### Visible Bell

The `visiblebell` variable causes the terminal to flash instead of beep:

```tcsh
set visiblebell
```

The `nobeep` variable disables all beeping entirely:

```tcsh
set nobeep
```

## The `ls-F` Builtin

`ls-F` lists files with type-indicating suffixes, using the same display logic as completion listing:

| Suffix | File Type |
|--------|-----------|
| `/` | Directory |
| `@` | Symbolic link |
| `*` | Executable file |
| `=` | Socket |
| `|` | Named pipe (FIFO) |
| `%` | Block device |
| `#` | Character device |

The `listflags` variable controls `ls-F` behavior:

| Value | Effect |
|-------|--------|
| Not set | Default listing |
| `x` | Multi-column listing (like `ls -x`) |
| Word beginning with `-` | Passed as flags to `ls` |

## History in the Editor

### History Search

- `history-search-backward` (`M-p`): Searches history for entries starting with the current word up to the cursor
- `history-search-forward` (`M-n`): Searches forward through history
- `i-search-back` / `i-search-fwd`: Incremental search (not bound by default)

### Dynamic Abbreviation

`dabbrev-expand` (`M-/`) expands the current word to the most recent preceding word that matches the prefix. This is similar to emacs' `dabbrev-expand`.

### History Expansion in the Editor

- `expand-history` (`M-space`): Expands history substitutions in the current word
- `expand-line`: Expands history substitutions in all words
- `toggle-literal-history` (`M-r`): Toggles between expanded and literal form

## Completion in the Editor

### How Completion Is Triggered

1. User presses Tab → `complete-word` editor command
2. The editor calls `tenematch()` in `tw.parse.c`
3. `tenematch()` parses the line, determines word context
4. If user completions exist, `tw_complete()` is called
5. Candidates are generated and either inserted or listed

### Completion Display Behavior

| Scenario | Behavior |
|----------|----------|
| Single unique match | Word is completed; suffix added (if `addsuffix` set) |
| Multiple matches, common prefix | Common prefix is inserted |
| Multiple matches, no common prefix | Terminal bell (if `matchbeep` allows); if `autolist` set, choices are listed |
| No matches | Terminal bell; nothing happens |

### Listing Completions

| Key | Command | Behavior |
|-----|---------|----------|
| `^D` (at end of line) | `delete-char-or-list-or-eof` | List completions |
| `M-^D` (anywhere) | `list-choices` | List completions |
| `^X-^D` | `list-choices-raw` | List completions (ignoring user definitions) |

### Cycling Through Completions

When multiple matches exist, `complete-word-fwd` and `complete-word-back` cycle through candidates, replacing the current word with each in turn. These are not bound by default but can be bound:

```tcsh
bindkey "^I" complete-word-fwd   # Tab cycles forward
bindkey "^X^I" complete-word-back # Ctrl-X Tab cycles backward
```

## References

- tcsh(1) man page — editor commands, bindkey, prompt format sequences
- `ed.defns.c` in the tcsh source — editor command definitions
- `ed.init.c` in the tcsh source — key binding initialization

## Related Skills

- [references/completion.md](references/completion.md) — the completion system that the editor invokes
- [references/quoting-expansion.md](references/quoting-expansion.md) — how quoting affects word boundaries in the editor
- [references/startup-config.md](references/startup-config.md) — where to set key bindings (~/.tcshrc)
