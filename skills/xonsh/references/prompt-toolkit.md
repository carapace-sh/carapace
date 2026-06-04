# Xonsh Prompt-Toolkit Integration

In-depth reference for how xonsh integrates with [prompt-toolkit](https://python-prompt-toolkit.readthedocs.io/) — the shell class, completer bridge, key bindings, completion display, and auto-suggestions.

## Architecture Overview

Xonsh uses prompt-toolkit (PTK) as its primary line-editing frontend. The integration is implemented in `xonsh/shells/ptk_shell/`:

| File | Purpose |
|------|---------|
| `xonsh/shells/ptk_shell/__init__.py` | `PromptToolkitShell` — main shell class |
| `xonsh/shells/ptk_shell/completer.py` | `PromptToolkitCompleter` — bridge to xonsh's completer |
| `xonsh/shells/ptk_shell/key_bindings.py` | Custom key bindings |
| `xonsh/shells/ptk_shell/formatter.py` | `PTKPromptFormatter` — prompt rendering |
| `xonsh/shells/ptk_shell/history.py` | `PromptToolkitHistory` — history backend |

## PromptToolkitShell

```python
class PromptToolkitShell(BaseShell):
    """The xonsh shell for prompt_toolkit v2 and later."""

    completion_displays_to_styles = {
        "multi": CompleteStyle.MULTI_COLUMN,
        "single": CompleteStyle.COLUMN,
        "readline": CompleteStyle.READLINE_LIKE,
        "none": None,
    }
```

### Initialization

```python
def __init__(self, **kwargs):
    self.prompter = PromptSession(**ptk_args)
    self.pt_completer = PromptToolkitCompleter(self.completer, self.ctx, self)
    self.key_bindings = load_xonsh_bindings(ptk_bindings)
```

### Shell Type Selection

Xonsh chooses between `prompt_toolkit` and `readline` shell types:

```python
@staticmethod
def choose_shell_type(init_shell_type=None, env=None):
    if shell_type == "prompt_toolkit":
        if not has_prompt_toolkit():
            # Falls back to readline
            shell_type = "readline"
        elif not ptk_above_min_supported():
            # Version too old, falls back to readline
            shell_type = "readline"
```

### The `on_ptk_create` Event

After PTK initialization, xonsh fires the `on_ptk_create` event:

```python
events.on_ptk_create.fire(
    prompter=self.prompter,
    history=self.history,
    completer=self.pt_completer,
    bindings=self.key_bindings,
)
```

This event allows xontribs and rc.xsh code to customize PTK behavior:

```python
@events.on_ptk_create
def _custom_keybindings(bindings, **kw):
    @bindings.add(prompt_toolkit.keys.Keys.ControlW)
    def say_hi(event):
        event.current_buffer.insert_text('hi')
```

## PromptToolkitCompleter

The bridge between xonsh's `Completer` and prompt-toolkit's `Completer` interface:

```python
class PromptToolkitCompleter(Completer):
    """Simple prompt_toolkit Completer object.
    It just redirects requests to normal Xonsh completer.
    """

    def __init__(self, completer, ctx, shell):
        self.completer = completer  # xonsh.completer.Completer instance
        self.ctx = ctx              # xonsh execution context
        self.shell = shell          # PromptToolkitShell instance
        self.hist_suggester = AutoSuggestFromHistory()
```

### `get_completions(document, complete_event)`

```python
def get_completions(self, document, complete_event):
    """Returns a generator for list of completions."""
```

**Flow:**

1. **Check if completions should be generated**:
   - `complete_event.completion_requested` (Tab key pressed), OR
   - `$UPDATE_COMPLETIONS_ON_KEYPRESS` is `True`

2. **Expand aliases** in the current line

3. **Calculate prefix and indices**:
   ```python
   line = document.current_line
   begidx = line[:endidx].rfind(" ") + 1
   prefix = line[begidx:endidx]
   ```

4. **Delegate to xonsh's Completer**:
   ```python
   completions, plen = self.completer.complete(
       prefix, line_ex, begidx, endidx, self.ctx,
       multiline_text=multiline_text, cursor_index=cursor_index,
   )
   ```

5. **Convert to PTK `Completion` objects** and yield

### Converting RichCompletion to PTK Completion

```python
if isinstance(comp, RichCompletion):
    desc = comp.description.replace(os.linesep, " ") if comp.description else None

    if comp.display:
        display = comp.display
    else:
        full_text = unquote(comp)
        disp = full_text[pre:]
        display = _highlight_match(disp, full_text, prefix, pre)

    yield Completion(
        comp,
        -comp.prefix_len if comp.prefix_len is not None else -plen,
        display=display,
        display_meta=desc,
        style=comp.style or "",
    )
```

| RichCompletion Field | PTK Completion Parameter | Notes |
|---------------------|--------------------------|-------|
| `value` (via `str(self)`) | First positional arg | The completion text |
| `prefix_len` | Second positional arg (negated) | Cursor movement offset |
| `display` | `display` | Shown in menu; computed if not set |
| `description` | `display_meta` | Newlines replaced with spaces |
| `style` | `style` | PTK style string |

### Display Text Highlighting

```python
def _highlight_match(display_text, full_text, prefix, pre):
    """Create FormattedText with underline on the matched substring."""
```

When the prefix matches as a substring (not at the start), the matched portion is underlined using PTK's `FormattedText`:

```python
def _underline_span(text, start, end):
    """Build a FormattedText that underlines text[start:end]."""
    parts = [("", text[:start]), ("underline", text[start:end])]
    if end < len(text):
        parts.append(("", text[end:]))
    return FormattedText(parts)
```

### `unquote(completion)`

Strips quotes from path completions for display text calculation:
- Handles raw strings (`r'...'`, `r"..."`)
- Handles triple-quoted strings
- Uses `ast.literal_eval` for escaped character resolution

### Auto-Suggestion

```python
def suggestion_completion(self, document, line):
    """Provides a completion based on the current auto-suggestion."""
```

When `AUTO_SUGGEST_IN_COMPLETIONS` is `True`, the auto-suggestion from history is included as the first completion option.

## Completion Display Configuration

### `$COMPLETIONS_DISPLAY`

Controls how completions are displayed in the PTK shell:

| Value | PTK Style | Description |
|-------|-----------|-------------|
| `"none"` or `False` | `None` | Don't display completions |
| `"single"` | `CompleteStyle.COLUMN` | Single column |
| `"multi"` or `True` | `CompleteStyle.MULTI_COLUMN` | Multiple columns (default) |
| `"readline"` | `CompleteStyle.READLINE_LIKE` | Emulate readline behavior |

### `$COMPLETIONS_MENU_ROWS`

Default: `5`. Number of rows to reserve for the completion menu. Only affects `single` and `multi` display modes.

### `$COMPLETION_MODE`

Default: `"default"`. Controls Tab completion behavior:

| Mode | Behavior |
|------|----------|
| `"default"` | First Tab selects common prefix, then cycles through completions |
| `"menu-complete"` | First Tab selects first whole completion, then cycles through remaining, then common prefix |

### `$COMPLETIONS_CONFIRM`

Default: `True`. When the completion menu is displayed, pressing Enter confirms the completion instead of running the command. Only affects the PTK shell.

### `$UPDATE_COMPLETIONS_ON_KEYPRESS`

Default: `False`. When `True`, completions are evaluated and presented on every keypress, not just on Tab. Only affects the PTK shell.

### `$COMPLETION_IN_THREAD`

Default: `False`. When `True`, completions are generated in a background thread to avoid blocking the UI. Useful when completers are slow.

### `$AUTO_SUGGEST_IN_COMPLETIONS`

Default: `False`. When `True`, the auto-suggestion result is placed as the first option in the completions list, enabling Tab completion of the auto-suggestion.

### `$XONSH_PROMPT_AUTO_SUGGEST`

Default: `True`. Enables automatic command suggestions based on history. Pressing the right arrow key inserts the displayed suggestion.

## Key Bindings

### Loading Custom Bindings

```python
def load_xonsh_bindings(ptk_bindings: KeyBindingsBase) -> KeyBindingsBase:
    """Load custom key bindings."""
    key_bindings = KeyBindings()
```

### Built-in Key Bindings

| Key | Handler | Description |
|-----|---------|-------------|
| `Tab` (with selection) | `indent_selection` | Indent selected lines |
| `Tab` (no selection, menu-complete) | `menu_complete_select` | Navigate completions |
| `Tab` (indent mode) | `insert_indent` | Insert indentation |
| `BackTab` | `dedent_current_line` | Dedent or navigate backward |
| `Ctrl+X, Ctrl+E` | `open_editor` | Open buffer in `$EDITOR` |
| `(`, `[`, `{` | `insert_left_paren` | Auto-pair brackets (if `$XONSH_AUTOPAIR`) |
| `'`, `"` | `insert_quote` | Auto-pair quotes (if `$XONSH_AUTOPAIR`) |
| `Ctrl+D` (empty buffer) | `call_exit_alias` | Exit xonsh |
| `Ctrl+M/J` (multiline) | `multiline_carriage_return` | Smart multiline handling |
| `Ctrl+Z` | `skip_control_z` | Prevent `^Z` writing |
| `Ctrl+Right` | `_fill` | Fill auto-suggestion |
| `Shift+Enter` | `shift_enter_newline` | Insert newline with indent |

### Condition Functions

```python
@Condition
def tab_insert_indent():
    """Check if <Tab> should insert indent instead of starting autocompletion."""
    before_cursor = get_app().current_buffer.document.current_line_before_cursor
    return bool(before_cursor.isspace())

@Condition
def tab_menu_complete():
    """Checks whether completion mode is `menu-complete`"""
    return XSH.env.get("COMPLETION_MODE") == "menu-complete"

@Condition
def should_confirm_completion():
    """Check if completion needs confirmation"""
    return XSH.env.get("COMPLETIONS_CONFIRM") and get_app().current_buffer.complete_state
```

### Adding Custom Key Bindings

Via `on_ptk_create` event:

```python
from prompt_toolkit.keys import Keys

@events.on_ptk_create
def _custom_bindings(bindings, **kw):
    @bindings.add(Keys.ControlG)
    def _grep_history(event):
        """Open grep-based history search"""
        event.current_buffer.insert_text('hello')
```

Or directly in rc.xsh:

```python
from xonsh.built_ins import XSH

@XSH.builtins.events.on_ptk_create
def _my_bindings(bindings, **kw):
    @bindings.add('c-t')
    def _transpose_chars(event):
        buf = event.current_buffer
        buf.swap_characters_before_cursor()
```

## Prompt Formatting

### PTKPromptFormatter

```python
class PTKPromptFormatter(PromptFormatter):
    """A subclass of PromptFormatter to support rendering prompt sections with/without threads."""

    def __call__(self, template=DEFAULT_PROMPT, fields=None, threaded=False,
                 prompt_name=None, **_) -> str:
```

Supports threaded prompt rendering for slow prompt segments (e.g., git status, kubernetes context).

### Prompt Configuration Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `PROMPT` | `"{env_name}{BOLD_GREEN}{user}@{hostname}{BOLD_BLUE} {cwd}{RESET} {BOLD_BLUE}${RESET} "` | Left prompt |
| `RIGHT_PROMPT` | `""` | Right-side prompt |
| `BOTTOM_TOOLBAR` | `None` | Bottom toolbar |
| `TITLE` | `""` | Terminal title |
| `MULTILINE_PROMPT` | `""` | Continuation prompt |
| `PROMPT_REFRESH_INTERVAL` | `None` | Auto-refresh interval in seconds |
| `UPDATE_PROMPT_ON_KEYPRESS` | `False` | Refresh prompt on each keypress |
| `PROMPT_FIELDS` | `{}` | Custom prompt fields |
| `VI_MODE` | `False` | Enable vi editing mode |
| `MOUSE_SUPPORT` | `False` | Enable mouse support |
| `PROMPT_TOOLKIT_COLOR_DEPTH` | `""` | Color depth override |
| `XONSH_PROMPT_CURSOR_SHAPE` | `None` | Cursor shape configuration |
| `XONSH_AUTOPAIR` | `False` | Auto-close brackets and quotes |
| `XONSH_COPY_ON_DELETE` | `False` | Copy deleted text to clipboard |
| `XONSH_USE_SYSTEM_CLIPBOARD` | `True` | Use system clipboard |

## History Integration

### PromptToolkitHistory

```python
class PromptToolkitHistory(prompt_toolkit.history.History):
    """History class that implements the prompt-toolkit history interface
    with the xonsh backend.
    """

    def store_string(self, entry):
        pass  # xonsh handles history independently

    def load_history_strings(self):
        """Loads synchronous history strings"""
```

### History Configuration

| Variable | Default | Purpose |
|----------|---------|---------|
| `XONSH_HISTORY_SIZE` | `8128` | Size of command history |
| `XONSH_HISTORY_FILE` | (platform-specific) | History file location |
| `XONSH_HISTORY_MATCH_ANYWHERE` | `False` | Match history anywhere in command |
| `XONSH_STORE_STDOUT` | `False` | Store stdout in history |
| `XONSH_STORE_STDIN` | `False` | Store stdin in history |

## Merged Key Bindings

The shell merges custom bindings with PTK's default bindings:

```python
self._key_bindings_merge = merge_key_bindings(
    [self.key_bindings, load_emacs_shift_selection_bindings()]
)
```

This allows xonsh's custom bindings to override PTK defaults while preserving unmodified defaults.
