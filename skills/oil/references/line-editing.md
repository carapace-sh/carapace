# Oil Shell Line Editing and Readline

In-depth reference for Oil shell's line editing, GNU Readline integration, completion display, and how the interactive editor works.

## How Readline Works in Oil

Oil is typically built with **GNU Readline**, which provides command line editing and completion capabilities. The architecture:

1. **GNU Readline** — handles basic line editing, key input, and the completion callback interface
2. **Oil Parser** — used as a library for history expansion and completion (unlike bash's ad-hoc parsers)
3. **InteractiveLineReader** — handles multi-line input and PS1/PS2 prompts
4. **comp_ui.py** — provides completion display classes (minimal or nice)

### ReadlineCallback Integration

Oil hooks into Readline via the `ReadlineCallback` class:

```python
class ReadlineCallback(object):
    def _GetNextCompletion(self, state):
        if state == 0:
            buf = self.readline.get_line_buffer()
            begin = self.readline.get_begidx()
            end = self.readline.get_endidx()

            comp = Api(line=buf, begin=begin, end=end)
            self.comp_iter = self.root_comp.Matches(comp)

        try:
            next_completion = self.comp_iter.next()
        except StopIteration:
            next_completion = None

        return next_completion
```

The `state == 0` call initializes the completion; subsequent calls iterate through candidates. This is the standard Readline completer protocol.

### Readline Init Commands

Oil's NiceDisplay sends initialization commands to Readline:

```python
def ReadlineInitCommands(self):
    return ['set horizontal-scroll-mode on']
```

This prevents line wrapping from clobbering completion output.

## Editing Modes

| Mode | How to Enable | Description |
|------|---------------|-------------|
| Emacs | Default | Standard Readline emacs key bindings |
| Vi | `set -o vi` | Vi-style key bindings via Readline |

## The .inputrc File

Oil uses GNU Readline's standard configuration. Readline reads:

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
| `$if term=xterm` | Test terminal type |
| `$if application=Bash` | Test application name (Oil may not match "Bash") |
| `$else` | Execute if test fails |
| `$endif` | Close conditional block |
| `$include /path/file` | Include another init file |

### Key Readline Variables for Completion

| Variable | Default | Description |
|----------|---------|-------------|
| `completion-display-width` | `-1` | Width of completion display (auto) |
| `completion-prefix-display-length` | `0` | Length of common prefix shown differently |
| `show-all-if-ambiguous` | `off` | List matches on first TAB instead of ringing bell |
| `show-all-if-unmodified` | `off` | List matches if no unique prefix |
| `colored-stats` | `off` | Color completions by type |
| `visible-stats` | `off` | Add character indicator to completions (`*` executable, `/` dir) |
| `colored-completion-prefix` | `off` | Color the common prefix differently |
| `completion-ignore-case` | `off` | Case-insensitive completion |
| `menu-complete-display-prefix` | `off` | Show prefix before cycling in menu-complete |
| `skip-completed-text` | `off` | Skip text already on the line during completion |
| `page-completions` | `on` | Page long completion lists |
| `completion-query-items` | `100` | Ask before listing this many completions |
| `bell-style` | `audible` | How to ring the bell (`audible`, `visible`, `none`) |
| `horizontal-scroll-mode` | `off` | Scroll horizontally instead of wrapping (set by NiceDisplay) |

## Completion Display: Minimal vs Nice

### MinimalDisplay

A display with minimal dependencies:

- **No color output**
- **No terminal width dependency**
- Could be useful for browser builds
- Uses `_RedrawPrompt()` to reprint prompt and command line
- Uses `_PrintPacked()` for simple packed column output

```python
def _RedrawPrompt(self):
    self.f.write(self.prompt_state.last_prompt_str)
    self.f.write(self.comp_state.line_until_tab)
```

### NiceDisplay

Full-featured completion display:

| Feature | Description |
|---------|-------------|
| Line erasure | Remembers lines drawn, erases before new content |
| Common prefix stripping | Strips common prefixes (Oil's rules, not Readline's) |
| Description display | Shows flag/builtin descriptions in yellow |
| Horizontal scrolling | Prevents line wrapping |
| Multi-press detection | Tracks repeated TAB presses via hash of matches |
| Right-side info | Shows additional info with reverse video |
| Cursor restoration | Uses ANSI escape codes to return cursor to prompt |

### ANSI Escape Codes Used by NiceDisplay

| Code | Meaning |
|------|---------|
| `\x1b[NdA` | Move cursor UP N lines |
| `\x1b[NdC` | Move cursor RIGHT N columns |
| `\x1b[2K` | Clear entire line |
| `\x1b[1B` | Move cursor down one line |
| `\x01`...`\x02` | Mark non-printing characters in prompts (Readline convention) |

### Display Functions

**`_PrintPacked(matches, max_match_len, term_width, max_lines, f)`** — Prints matches in columns:
- Width of each candidate = `max_match_len + 2`
- Number per line = `(term_width - 2) // w`
- Shows "... and N more" when exceeding `max_lines`

**`_PrintLong(matches, max_match_len, term_width, max_lines, descriptions, f)`** — Prints with descriptions:
- Format: `%-max_match_len s %s` with descriptions in YELLOW
- Truncates long descriptions with " ..."

## Prompt Handling

### PS1 and PS2

| Variable | Purpose |
|----------|---------|
| `PS1` | Primary prompt (first line of input) |
| `PS2` | Continuation prompt (multi-line input) |

### InteractiveLineReader

Oil handles the PS1/PS2 decision **outside** the lexer and parser, in the `InteractiveLineReader`:

- Manages the PS1/PS2 prompt decision
- Confines prompt knowledge to itself rather than spreading it throughout the parser
- Reads lines interactively and decides continuation

This contrasts with bash/dash, which "litter" their parsers with references to prompts through global variables like `doprompt` and `checkkwd`.

### Prompt Escape Sequences

Oil supports standard bash prompt escapes in PS1:

| Escape | Description |
|--------|-------------|
| `\u` | Username |
| `\h` | Hostname (short) |
| `\H` | Hostname (full) |
| `\w` | Working directory |
| `\W` | Basename of working directory |
| `\s` | Shell name |
| `\$` | `$` for regular users, `#` for root |
| `\t` | Time (24-hour) |
| `\A` | Time (24-hour, no seconds) |
| `\d` | Date |
| `\n` | Newline |
| `\[...\]` | Non-printing characters (for Readline) |

### YSH Prompt

YSH uses a function-based prompt instead of escape sequences:

```bash
func renderPrompt(io) {
  return (io.promptVal('s') ++ ' ')
}
```

## History

### History File

Oil stores history in `~/.local/share/oils/osh_history`.

### History Expansion

Oil uses its **full parser** for history expansion, unlike bash which uses a partial/incorrect parser:

```bash
# Bash: incorrectly splits the word
bash$ echo ${x:-a b c}
a b c
bash$ echo !$
echo c}                 # Bug: splits incorrectly

# Oil: correctly handles the last argument
osh$ echo ${x:-a b c}
a b c
osh$ echo !$
echo ${x:-a b c}       # Correct: full last argument
```

### Parser Architecture for History

Oil's parser has two outputs:

1. **The lossless syntax tree** — when parsing succeeds
2. **A "trail" of words and tokens** — when parsing fails on incomplete input

The trail enables both history expansion and autocompletion to use the same parser. This is possible because Oil uses **top-down parsing**, unlike bash's bottom-up LR parsing (yacc/Bison), which makes generating trails for incomplete input difficult.

## Alias Handling and the Parser

Oil re-invokes the parser as a **library** for alias expansion, rather than scattering `CHKALIAS` checks throughout the parser:

| Approach | Shell | Mechanism |
|----------|-------|-----------|
| Global flags | bash, dash | `checkkwd = CHKNL \| CHKKWD \| CHKALIAS` |
| Parser re-invocation | Oil | Parser called as library with alias context |

This keeps global variables out of the grammar and makes behavior "documentable."

### Shell Disagreements on Alias Behavior

Spec tests reveal shells disagree on:

- First and second word being the same alias with no trailing space
- Loop split across iterative and recursive aliases
- Alias respected inside `eval`
- Here doc inside alias

## Key Bindings

### Default Emacs Mode Bindings (Completion-Related)

| Key | Function |
|-----|----------|
| `TAB` | `complete` — invoke completion |
| `M-?` | `possible-completions` — list completions |
| `M-*` | `insert-completions` — insert all completions |
| `M-<TAB>` | `tab-insert` — insert literal TAB |
| `C-x /` | `complete-filename` — filename completion |
| `M-~` | `complete-username` — username completion |
| `M-$` | `complete-variable` — variable completion |
| `M-@` | `complete-hostname` — hostname completion |
| `M-!` | `complete-command` — command completion |
| `C-M-i` | `dynamic-complete-history` — history completion |

### Vi Mode Bindings (Completion-Related)

| Key (command mode) | Function |
|---------------------|----------|
| `\t` | `complete` — invoke completion |
| `*` | `vi-complete` — complete and list |
| `=` | `possible-completions` — list completions |

## Differences from Bash: Summary

| Feature | Bash | Oil |
|---------|------|-----|
| Completion source | Ad-hoc parsers | Parser as library |
| History expansion | Buggy on complex expressions | Accurate (full parser) |
| Prompt handling | Littered in parser | Isolated in InteractiveLineReader |
| Alias expansion | Global flags in parser | Parser re-invocation |
| Completion UI | Readline only | `minimal` or `nice` (OILS_COMP_UI) |
| Description display | Not supported | NiceDisplay shows descriptions |
| Line erasure | Not supported | NiceDisplay erases and redraws |

## References

- [Oil Shell Front End Reference](https://oils.pub/release/latest/doc/ref/chap-front-end.html)
- [Oil Shell Blog: History and Completion](https://www.oilshell.org/blog/2020/01/history-and-completion.html)
- [Oil Shell Blog: Alias and Prompt](https://www.oilshell.org/blog/2020/01/alias-and-prompt.html)
- [GNU Readline Manual](https://tiswww.case.edu/php/chet/readline/rltop.html)

## Related Skills

- **bash skill → references/readline.md** — bash Readline integration (shares GNU Readline library)
- **bash skill → references/completion.md** — bash programmable completion
