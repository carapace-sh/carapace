# Xonsh Completion System

In-depth reference for xonsh's completion system — the completer pipeline, RichCompletion, CommandContext, decorators, registration, filtering, sorting, and the built-in completers.

## Architecture Overview

Xonsh's completion system is implemented in Python using [prompt-toolkit](https://python-prompt-toolkit.readthedocs.io/). Unlike traditional shells, completers are **Python callables** registered in an `OrderedDict`, completions are **Python objects** (`RichCompletion`), and the completion protocol is **contextual** — completers receive a parsed `CommandContext` object.

### Key Files

| File | Purpose |
|------|---------|
| `xonsh/completer.py` | Main `Completer` class — orchestrates the pipeline |
| `xonsh/completers/tools.py` | `RichCompletion`, decorators, helpers |
| `xonsh/completers/completer.py` | `add_one_completer`, `list_completers`, `remove_completer` |
| `xonsh/parsers/completion_context.py` | `CompletionContextParser`, `CommandContext`, `CommandArg` |
| `xonsh/completers/base.py` | Base completer (commands + Python + paths) |
| `xonsh/completers/commands.py` | Command/alias completion, `CommandCompleter` (xompletions) |
| `xonsh/completers/path.py` | File/directory path completion |
| `xonsh/completers/python.py` | Python expression completion |
| `xonsh/completers/bash.py` | Bash completion bridge (high level) |
| `xonsh/completers/bash_completion.py` | Bash completion bridge (low level) |

## The Completion Pipeline

When the user presses TAB in xonsh:

1. **Parse** — `Completer.parse()` → `CompletionContextParser.parse()` → `CompletionContext`
2. **Generate** — `Completer.complete_from_context()` → `Completer.generate_completions()`
3. **Iterate completers** — For each completer in `XSH.completers` (OrderedDict):
   - Fire `on_completer_filter` event (handlers can veto)
   - Call completer function (contextual or legacy style)
   - Filter results using `get_filter_function()`
   - Format via `_format_completion()`
   - **If exclusive completer returns items → break** (skip remaining)
4. **Enforce query limit** — `COMPLETION_QUERY_LIMIT` truncates with warning
5. **Deduplicate** — Remove trailing-space variants (keep spaced version)
6. **Sort** — By match quality tier, then underscore prefix, position, name
7. **Return** — `(tuple_of_sorted_completions, lprefix)`

### Pipeline Flow Diagram

```
User presses TAB
  ↓
CompletionContextParser.parse(text, cursor_index)
  → CompletionContext(command=CommandContext | None, python=PythonContext | None)
  ↓
Completer.complete_from_context(completion_context)
  ↓
Completer.generate_completions(completion_context, old_args, trace)
  ├─ for name, func in XSH.completers.items():
  │   ├─ on_completer_filter event → veto check
  │   ├─ call func(completion_context) or func(*old_args)
  │   ├─ filter results with get_filter_function()
  │   ├─ format with _format_completion()
  │   └─ if exclusive and has results → BREAK
  ↓
Enforce COMPLETION_QUERY_LIMIT
  ↓
Deduplicate (trailing space variants)
  ↓
Sort by match quality tier
  ↓
Return (completions, lprefix)
```

## CompletionContext and Parsing

### CompletionContext

```python
class CompletionContext(NamedTuple):
    command: CommandContext | None = None
    python: PythonContext | None = None
```

Contains the parsed completion context. Either `command` or `python` may be `None` depending on what the cursor is in.

### CommandContext

```python
class CommandContext(NamedTuple):
    args: tuple[CommandArg, ...]    # Parsed arguments
    arg_index: int                   # Current argument's index (-1 if not in command)
    prefix: str = ""                 # Current arg's prefix (before cursor, without quotes)
    suffix: str = ""                 # Current arg's suffix (after cursor)
    opening_quote: str = ""          # Opening quote of current arg (e.g. "'", '"', 'r"')
    closing_quote: str = ""          # Closing quote of current arg
    is_after_closing_quote: bool = False  # Cursor is after a closed string literal
    subcmd_opening: str = ""         # If inside a subproc expression (e.g. '$(', '!')
```

**Key properties and methods:**

| Property/Method | Returns | Description |
|-----------------|---------|-------------|
| `command` | `str \| None` | First argument's `raw_value` (command name including quotes) |
| `raw_prefix` | `str` | Prefix including quotes: `f"{opening_quote}{prefix}{closing_quote}"` (or without closing quote if cursor is inside) |
| `completing_command(cmd)` | `bool` | `True` when `arg_index > 0` and first arg is `cmd` |
| `words_before_cursor` | `str` | Space-joined raw values of args before current |
| `text_before_cursor` | `str` | Full text before cursor including prefix |
| `begidx` | `int` | Cursor's position (length of text_before_cursor) |

### CommandArg

```python
class CommandArg(NamedTuple):
    value: str                     # The argument's value (without quotes)
    opening_quote: str = ""        # Opening quote
    closing_quote: str = ""        # Closing quote
    is_io_redir: bool = False      # Whether the arg is IO redirection
```

| Property | Returns | Description |
|----------|---------|-------------|
| `raw_value` | `str` | `f"{opening_quote}{value}{closing_quote}"` — complete argument including quotes |

### PythonContext

```python
class PythonContext(NamedTuple):
    multiline_code: str            # The multi-line Python code
    cursor_index: int              # Cursor's index in the code
    is_sub_expression: bool = False  # Whether this is @(...)
    ctx: dict[str, Any] | None = None  # Objects in current execution context
```

| Property | Returns | Description |
|----------|---------|-------------|
| `prefix` | `str` | Code from start to cursor (`multiline_code[:cursor_index]`) |

### CompletionContextParser

The parser tokenizes the command line using PLY (Python Lex-Yacc) with an LALR grammar. Key behaviors:

- **Token types used**: `STRING`, `ANY` (catch-all for unknown tokens)
- **Parenthesis matching**: Tracks `DOLLAR_LPAREN`/`RPAREN`, `BANG_LPAREN`/`RPAREN`, `AT_LPAREN`/`RPAREN`, etc.
- **Multi-token handling**: `AND`/`OR` tokens (`&&`/`||`) become `ANY` if cursor is inside
- **Ignored tokens**: `INDENT`, `DEDENT`, `WS`
- **Cursor detection**: The parser determines which token the cursor is in and builds the appropriate context

## RichCompletion

`RichCompletion` is a subclass of `str` that carries rich metadata for each completion candidate:

```python
class RichCompletion(str):
    def __init__(
        self,
        value: str,
        prefix_len: int | None = None,
        display: str | None = None,
        description: str = "",
        style: str = "",
        append_closing_quote: bool = True,
        append_space: bool = False,
        provider: str | None = None,
    ):
```

| Field | Type | Default | Purpose |
|-------|------|---------|---------|
| `value` | `str` | (required) | The completion value (also the string value via `str.__new__`) |
| `prefix_len` | `int \| None` | `None` | Length of prefix to replace; `None` = use default |
| `display` | `str \| None` | `None` | Display text; if set, common prefix stripping is disabled |
| `description` | `str` | `""` | Description shown when completion is selected |
| `style` | `str` | `""` | Style string for prompt-toolkit's `Completion` object |
| `append_closing_quote` | `bool` | `True` | Whether to append closing quote if cursor is after one |
| `append_space` | `bool` | `False` | Whether to append space after completion (after closing quote) |
| `provider` | `str \| None` | `None` | Debug-only tag for `$XONSH_COMPLETER_TRACE` |

### Key Methods

```python
@property
def value(self):
    return str(self)

def replace(self, **kwargs):
    """Create a new RichCompletion with replaced attributes"""
```

### `append_closing_quote`

When `True` (the default), xonsh appends the closing quote to the completion value when the cursor is positioned right after a closing quote. For example, with `ls "/usr/"⇥`, if the closing quote is `"`, the completion value gets `"` appended.

External completers (like carapace) typically set `append_closing_quote=False` because they handle quoting in their own formatters — letting xonsh also add closing quotes would double-quote the completion.

### `append_space`

When `True`, xonsh appends a trailing space to the completion after insertion. The space is placed **after** the closing quote (if applicable). Used by xonsh's built-in completers for commands and single-match candidates.

External completers typically do **not** use `append_space`. Instead, they encode trailing space directly in the `Value` field — adding `" "` to the value when nospace does not match. This is simpler and avoids interaction with `append_closing_quote`.

### `prefix_len`

Controls how many characters of the current word are replaced by the completion. When `None` (default), xonsh calculates the default prefix length from the cursor position. When explicitly set, the completer controls the replacement length.

External completers typically set `prefix_len=len(context.raw_prefix)` to ensure the full typed prefix — including any opening quotes — is replaced by the completion value.

### `display`

When set, this text is shown in the completion menu instead of the value. Setting `display` also disables prompt-toolkit's common prefix stripping for that completion. This is useful when the completion value differs from what should be displayed (e.g., tilde expansion on Windows).

### `provider`

A debug-only tag identifying the sub-source of the completion. Used by `$XONSH_COMPLETER_TRACE` to show which completer produced each result. Common values: `"alias"`, `"command"`, `"python"`, `"path"`, `"bash"`.

## Type Aliases

```python
Completion = Union[RichCompletion, str]
CompleterResult = Union[set[Completion], tuple[set[Completion], int], None]
ContextualCompleter = Callable[[CompletionContext], CompleterResult]
```

A completer can return:
- `None` — no completions (not applicable)
- `set[Completion]` — a set of completions
- `(set[Completion], lprefix)` — completions with explicit prefix length override

## Completer Decorators

### `@contextual_completer`

```python
def contextual_completer(func: ContextualCompleter):
    """Decorator for a contextual completer that uses the parsed completion context."""
    func.contextual = True
    return func
```

Marks the function as a contextual completer. The function receives a `CompletionContext` object directly. Checked by `is_contextual_completer(func)`.

### `@contextual_command_completer`

```python
def contextual_command_completer(func: Callable[[CommandContext], CompleterResult]):
    """Like contextual_completer, but only runs when completing a command."""
    @contextual_completer
    @wraps(func)
    def _completer(context: CompletionContext) -> CompleterResult:
        if context.command is not None:
            return func(context.command)
        return None
    return _completer
```

The wrapped function:
1. Only runs during command completion (not Python expressions)
2. Receives a `CommandContext` object directly (not the full `CompletionContext`)
3. Returns `None` when `context.command is None`

### `@contextual_command_completer_for(cmd)`

```python
def contextual_command_completer_for(cmd: str):
    """Like contextual_command_completer, but only for a specific command."""
    def decor(func):
        @contextual_completer
        @wraps(func)
        def _completer(context: CompletionContext) -> CompleterResult:
            if context.command is not None and context.command.completing_command(cmd):
                return func(context.command)
            return None
        return _completer
    return decor
```

Convenience decorator for command-specific completers. Only runs when completing arguments for the specified command.

### `@non_exclusive_completer`

```python
def non_exclusive_completer(func):
    """Decorator for a non-exclusive completer."""
    func.non_exclusive = True
    return func
```

Non-exclusive completers' results are **merged** with other completers' results. The pipeline continues after them. Checked by `is_exclusive_completer(func)` (returns `not getattr(func, 'non_exclusive', False)`).

## Completer Registration

### `add_one_completer(name, func, loc='end')`

```python
def add_one_completer(name, func, loc="end"):
```

Registers a completer in `XSH.completers` (an `OrderedDict`). The `loc` parameter controls insertion order:

| Position | Behavior |
|----------|----------|
| `'start'` | Insert before the first exclusive completer (highest priority for non-exclusive) |
| `'end'` | Insert at the end (lowest priority) |
| `'<name'` | Insert **before** the named completer |
| `'>name'` | Insert **after** the named completer |

**`'start'` behavior detail**: The completer is inserted before the first exclusive completer in the OrderedDict. This means non-exclusive completers at `'start'` run before any exclusive completer, but existing non-exclusive completers at the start are shifted after the new one.

### `list_completers()`

Lists all active completers with their descriptions.

### `remove_completer(name)`

Removes a completer from xonsh.

## Exclusive vs Non-Exclusive Completers

### Exclusive Completers (default)

- If they return **any** items, the pipeline **breaks immediately**
- Other completers (both exclusive and non-exclusive) are skipped
- Used when a completer has authoritative completions (e.g., bash completions, command-specific completers)
- Checked by: `is_exclusive_completer(func)` returns `True`

### Non-Exclusive Completers

- Decorated with `@non_exclusive_completer`
- All non-exclusive completers are called and their results are **accumulated**
- Pipeline continues until all completers are exhausted or query limit reached
- Used for supplementary completions (e.g., path completions, end-of-pipe tokens)
- Checked by: `is_exclusive_completer(func)` returns `False`

### Interaction Rules

1. Non-exclusive completers run in order and accumulate results
2. When an exclusive completer returns results, the pipeline stops
3. If an exclusive completer returns `None` or empty, the pipeline continues
4. Non-exclusive completers after a successful exclusive completer are skipped

## Filtering

### `get_filter_function()`

Returns the completion filter function based on `$XONSH_COMPLETER_MODE`:

| Mode | Filter | Description |
|------|--------|-------------|
| `"substring_tier"` (default) | `_filter_substring` | Case-insensitive substring match |
| `"prefix"` | `_filter_prefix` | Only completions starting with prefix (case-insensitive) |

### Filter Implementation

```python
def _filter_substring(text, prefix):
    """Case-insensitive substring match"""
    func = lambda txt, pre: pre.lower() in txt.lower()
    return _filter_with_func(text, prefix, func)

def _filter_prefix(text, prefix):
    """Case-insensitive prefix match"""
    func = lambda txt, pre: txt.lower().startswith(pre.lower())
    return _filter_with_func(text, prefix, func)
```

For `RichCompletion` objects with `display` set, the filter matches against comma-separated parts of the display text (after stripping whitespace).

## Sorting

Completions are sorted by match quality tier:

| Tier | Match Type | Description |
|------|-----------|-------------|
| 0 | Case-sensitive prefix | `text.startswith(prefix)` |
| 1 | Case-insensitive prefix | `text.lower().startswith(prefix.lower())` |
| 2 | Case-sensitive substring | `prefix in text` |
| 3 | Case-insensitive substring | `prefix.lower() in text.lower()` |
| 4 | No match | Sorted last |

**Within each tier**, the sort key is: `(tier, has_leading_underscore, match_position, lowercased_text)`

- Names whose last component starts with `_` are sorted last (handles both `_codecs` and `json._default_decoder`)
- Earlier match position first
- Alphabetically by lowercased text

**Fallback sort** (when no prefix/context): `s.lstrip('\'"').lower()`

## Deduplication

After all completers have been consulted, xonsh deduplicates completions that differ only by a trailing space:

```python
spaced = {str(c) for c in completions if str(c).endswith(" ")}
if spaced:
    completions = {
        c: None
        for c in completions
        if str(c).endswith(" ") or (str(c) + " ") not in spaced
    }
```

This handles cases like `_cd` (from Python name completions) and `_cd ` (from command completions with `append_space=True`) — the spaced variant is kept.

## The `Completer` Class

```python
class Completer:
    """This provides a list of optional completions for the xonsh shell."""

    def __init__(self):
        self.context_parser = CompletionContextParser()
```

### Key Methods

| Method | Signature | Purpose |
|--------|-----------|---------|
| `parse` | `(text, cursor_index=None, ctx=None) → CompletionContext` | Parse text into completion context |
| `complete_line` | `(text) → (completions, lprefix)` | Complete when cursor is at end |
| `complete` | `(prefix, line, begidx, endidx, ctx=None, ...) → (completions, lprefix)` | Main entry point (legacy API) |
| `complete_from_context` | `(completion_context, old_completer_args=None) → (completions, lprefix)` | Main orchestration method |
| `generate_completions` | `(completion_context, old_completer_args, trace) → Iterator[(Completion, int)]` | Generator yielding completions from all completers |
| `_format_completion` | `(completion, context, completing_contextual, lprefix, custom_lprefix) → (Completion, int)` | Format completion with closing quotes and spaces |

### `_format_completion` Details

Handles the special case when the cursor is appending to a closed string literal (e.g., `ls "/usr/"⇥`):

1. If `completing_contextual_command` and context is after closing quote:
   - If not `custom_lprefix`, extends lprefix to include closing quote
   - If completion has `append_closing_quote=True`, appends closing quote
   - Applies lprefix via `apply_lprefix()`
2. If completion has `append_space=True`, appends space after closing quote
3. Returns `(completion, lprefix)`

## Built-in Completers

### Base Completer (`complete_base`)

```python
@contextual_completer
def complete_base(context: CompletionContext):
    """If the line is empty, complete based on valid commands, python names, and paths."""
```

Only runs when completing the first argument (`context.command.arg_index == 0`). Chains:
1. `complete_python(context)` → tagged `"python"`
2. `complete_command(context.command)` → tagged `"alias"` or `"command"`
3. `contextual_complete_path(context.command, cdpath=False)` → tagged `"path"` (only if prefix is empty)

### Command Completer (`complete_command`)

```python
@contextual_command_completer
def complete_command(command: CommandContext):
```

Iterates through `XSH.commands_cache` to find matching commands and aliases. For aliases, generates descriptions from docstrings or alias contents. For regular commands, shows path as description if `CMD_COMPLETIONS_SHOW_DESC` is enabled.

Yields `RichCompletion` with `append_space=True` and `provider="alias"` or `provider="command"`.

### Skipper Completer (`complete_skipper`)

```python
@contextual_command_completer
def complete_skipper(command_context: CommandContext):
```

Skips over prefix commands like `time`, `timeit`, `which`, `showcmd`, `man` and delegates to the inner command's completer. Creates a modified context with skipped tokens removed.

### End-of-Pipe Token Completers

```python
@non_exclusive_completer
@contextual_command_completer
def complete_end_proc_tokens(command_context: CommandContext):
    """Insert space after |, ;, && when there's no space following."""

@non_exclusive_completer
@contextual_command_completer
def complete_end_proc_keywords(command_context: CommandContext):
    """Insert space after 'and' or 'or' keywords."""
```

Both return `{RichCompletion(prefix, append_space=True)}` when the prefix matches.

### Path Completer (`complete_path`)

```python
@contextual_completer
def complete_path(context: CompletionContext):
```

Handles file/directory completion with:
- Quote-aware path parsing (raw strings, p-strings, regular strings)
- Tilde expansion (skipped for raw strings)
- Glob matching (case-insensitive on non-Windows)
- Subsequence matching (fish/zsh-style, e.g., `~/u/ro` → `~/lou/carcohl`)
- Fuzzy matching (Levenshtein distance with `SUGGEST_THRESHOLD`)
- CDPATH completion (when `cd` is in the command)
- Per-path quoting (each path quoted individually)
- Raw string promotion for paths with `\` or `$`
- Control character handling (filenames with `\n`, `\t` use non-raw strings)
- Dot completion (`.`/`..`) controlled by `COMPLETE_DOTS`

### Python Completer (`complete_python`)

```python
@contextual_completer
def complete_python(context: CompletionContext):
```

Handles Python expression completion:
- Attribute access completion (`obj.attr`)
- Function signature completion (parameter names with `=` suffix)
- Import completion (`__xonsh__.imp.<module>`)
- Context dictionary completion
- Builtins from `dir(builtins)`
- Xonsh-specific tokens (`XONSH_EXPR_TOKENS`, `XONSH_STMT_TOKENS`)
- Bracket behavior controlled by `COMPLETIONS_BRACKETS`:
  - Callables with useful attributes → plain name
  - Simple functions/methods → adds `(` for calling
  - Sequences/mappings → adds `[` for indexing
- Empty prefix filtering: excludes private/dunder attributes and all tokens

### Bash Completion Bridge (`complete_from_bash`)

```python
@contextual_command_completer
def complete_from_bash(context: CommandContext):
    """Completes based on results from BASH completion."""
```

Wraps bash's programmable completion system by:
1. Transforming `CommandContext` into bash's `COMP_*` variable format
2. Executing bash with a script that sources completion files and invokes the completion function
3. Parsing `COMPREPLY` output
4. Converting results to `RichCompletion` with `append_closing_quote=False`
5. Adjusting `lprefix` for quote handling

The bash script template:
- Sources bash completion framework from `$BASH_COMPLETIONS` paths
- Overrides `quote_readline` and `_quote_readline_by_ref` to prevent double-quoting
- Locates the completion function via `complete -p`
- Calls `_completion_loader` for lazy loading
- Sets `COMP_WORDS`, `COMP_LINE`, `COMP_POINT`, `COMP_CWORD`
- Invokes the completion function and prints `COMPREPLY`

### CommandCompleter (xompletions)

```python
class CommandCompleter:
    """Lazily complete commands from xompletions package."""
```

Delegates to external completer modules in the `xompletions` package or `$XONSH_COMPLETER_DIRS`. Looks up modules by command name (exact or regex pattern). Calls the module's `xonsh_complete` function if it exists.

```python
complete_xompletions = CommandCompleter()  # Singleton instance
```

## Helper Functions

### `apply_lprefix(comps, lprefix)`

Apply `lprefix` to completions, converting plain strings to `RichCompletion` if needed. If a `RichCompletion` already has `prefix_len` set, it is preserved.

### `tag_provider(result, provider)`

Tag completer output with `provider` for `$XONSH_COMPLETER_TRACE`. Accepts any of the three standard return shapes and preserves the shape.

### `complete_from_sub_proc(*args, sep=None, filter_prefix=None, **env_vars)`

Helper to complete from subprocess output. Runs the command, parses output (by lines or custom separator), and yields `RichCompletion` objects. If only one completion candidate, sets `append_space=True` (unless it ends with a path separator).

### `comp_based_completer(ctx, start_index=0, **env)`

Helper for commands that use bash's `COMP_*` variable interface (e.g., pip, django-admin). Sets `COMP_WORDS` and `COMP_CWORD` and delegates to `complete_from_sub_proc`.

### `completion_from_cmd_output(line, append_space=False)`

Parses a single line of completion output. Tab-separated format: `cmd\tdescription`. Paths ending with a separator get `append_space=False`.

## Event: `on_completer_filter`

```python
events.doc("on_completer_filter", """
on_completer_filter(completer: str, command: str, context: CompletionContext) -> bool | None

Fires once for every registered completer just before it is invoked,
allowing handlers to veto the call.
""")
```

- Return `False` to skip the completer
- Return `True` or `None` to let it run
- Multiple handlers: any handler returning `False` vetoes the completer

Example — only allow bash completer for specific commands:

```python
@events.on_completer_filter
def _only_bash_for_some(completer, command, context, **_):
    return command in {'kubectl', 'docker'} if completer == 'bash' else True
```

## Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| `XONSH_COMPLETER_MODE` | `"substring_tier"` | Filter mode: substring or prefix |
| `XONSH_COMPLETER_TRACE` | `False` | Show completer trace for debugging |
| `COMPLETION_QUERY_LIMIT` | `1000` | Max completions before truncation |
| `XONSH_COMPLETER_DIRS` | (platform-specific) | Paths to search for command completers |
| `XONSH_COMPLETER_EMOJI_PREFIX` | `None` | Trigger prefix for emoji completion |
| `XONSH_COMPLETER_SYMBOLS_PREFIX` | `None` | Trigger prefix for unicode symbol completion |
| `BASH_COMPLETIONS` | (platform-specific) | Paths to bash completion scripts |
| `COMPLETE_DOTS` | `"matching"` | `.`/`..` completion: `"always"`, `"never"`, `"matching"` |
| `SUBSEQUENCE_PATH_COMPLETION` | `True` | Enable fish/zsh-style subsequence path matching |
| `FUZZY_PATH_COMPLETION` | `True` | Enable fuzzy path matching (Levenshtein) |
| `SUGGEST_THRESHOLD` | `3` | Levenshtein distance threshold for fuzzy matching |
| `CMD_COMPLETIONS_SHOW_DESC` | `False` | Show path/alias in command completion descriptions |
| `COMPLETIONS_BRACKETS` | `True` | Include `(`/`[` in Python attribute completions |
| `ALIAS_COMPLETIONS_OPTIONS_BY_DEFAULT` | `False` | Show options without `-` prefix for argparse completions |
| `ALIAS_COMPLETIONS_OPTIONS_LONGEST` | `False` | Show only longest option variant |

## Writing a Custom Completer

### Basic Command Completer

```python
from xonsh.completers.tools import contextual_command_completer, RichCompletion

@contextual_command_completer
def _my_command_completer(context):
    """Completer for my_command"""
    if context.completing_command('my_command'):
        return {
            RichCompletion("start", description="Start the service"),
            RichCompletion("stop", description="Stop the service"),
            RichCompletion("status", description="Show status"),
        }
```

### Completer with Subprocess

```python
from xonsh.completers.tools import contextual_command_completer, RichCompletion, sub_proc_get_output
from json import loads

@contextual_command_completer
def _my_command_completer(context):
    """Completer for my_command using subprocess"""
    if context.completing_command('my_command'):
        output, not_found = sub_proc_get_output(
            'my_command', 'complete', *[a.value for a in context.args], context.prefix
        )
        if not_found:
            return None
        try:
            return {RichCompletion(c["value"], description=c.get("desc", ""))
                    for c in loads(output)}
        except:
            return None
```

### Non-Exclusive Completer

```python
from xonsh.completers.tools import non_exclusive_completer, contextual_command_completer, RichCompletion

@non_exclusive_completer
@contextual_command_completer
def _my_supplementary_completer(context):
    """Adds extra completions alongside other completers"""
    return {RichCompletion("--verbose", description="Enable verbose mode")}
```

### Registering in rc.xsh

```python
from xonsh.completers.completer import add_one_completer

add_one_completer('my_completer', _my_command_completer, 'start')
```

## Debugging Completions

### `$XONSH_COMPLETER_TRACE`

Set to `True` to see which completers are invoked and their results. Each completion is printed with its source completer and non-default `RichCompletion` attributes.

Format:
```
"value": src=<completer>, pvd=<sub-source>, type=<exclusive|non-exclusive>, <non-default attrs>
```

### `list_completers()`

Lists all registered completers with their descriptions and order.

### `on_completer_filter` Event

Use this event to selectively enable/disable specific completers for debugging or customization.
