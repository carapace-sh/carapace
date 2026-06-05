# Carapace Library: Xonsh Shell Integration Deep Dive

Reference for [carapace](https://github.com/carapace-sh/carapace)'s xonsh completion integration — how the snippet works, how completion output is formatted, and how carapace handles xonsh-specific edge cases including quoting, RichCompletion construction, and style conversion. For cross-shell comparisons, see the **references/shell.md**.

## Source Files

| File | Purpose |
|------|---------|
| `internal/shell/xonsh/snippet.go` | Xonsh completion script generation (Python `@contextual_command_completer`) |
| `internal/shell/xonsh/action.go` | Value formatting, JSON serialization, quoting/nospace logic |
| `internal/shell/xonsh/style.go` | Style conversion from carapace format to xonsh `bg: fg:` format |
| `internal/shell/shell.go` | Shared dispatch — message integration, nospace propagation, filtering |
| `complete.go` | Entry point — no xonsh-specific patching (unlike bash/nushell/cmd-clink) |

## Xonsh Completion System Background

### Xonsh: A Python-Powered Shell

Xonsh is a Python-powered shell that blends Python and shell syntax. Its completion system is implemented in Python using [prompt-toolkit](https://python-prompt-toolkit.readthedocs.io/), which provides the completion UI. Xonsh's completer infrastructure is fundamentally different from traditional shells:

- Completers are **Python callables** registered in an `OrderedDict` (`__xonsh__.completers`)
- Completions are **Python objects** (`RichCompletion`), not raw text
- The completion protocol is **contextual** — completers receive a parsed `CommandContext` object
- There is **no shell script** in the completion callback — everything is Python

### The `@contextual_command_completer` Decorator

```python
from xonsh.completers.tools import contextual_command_completer

@contextual_command_completer
def _my_command_completer(context):
    """carapace completer for my_command"""
    if context.completing_command('my_command'):
        # ... return completions
```

`@contextual_command_completer` is a decorator that:

1. Marks the function as a **contextual completer** (sets `func.contextual = True`)
2. Wraps the function so it receives a `CommandContext` object directly (instead of the more general `CompletionContext`)
3. Returns `None` when not completing a command (i.e., when `context.command is None`), preventing the completer from running for Python expressions

The wrapped function only runs during command completion. For a command-specific variant, use `@contextual_command_completer_for('cmd')` which additionally checks `context.command.completing_command('cmd')`.

### The `CommandContext` Object

Xonsh parses the command line using its own Python-based parser and provides the completer with a structured `CommandContext`:

```python
class CommandContext(NamedTuple):
    args: tuple[CommandArg, ...]           # Parsed arguments
    arg_index: int                          # Current argument's index (-1 if not in command)
    prefix: str = ""                       # Current string arg's prefix (before cursor)
    suffix: str = ""                        # Current string arg's suffix (after cursor)
    opening_quote: str = ""                # Opening quote of current arg (e.g. "'", '"', 'r"')
    closing_quote: str = ""                # Closing quote of current arg
    is_after_closing_quote: bool = False   # Cursor is after a closed string literal
    subcmd_opening: str = ""               # If inside a subproc expression (e.g. '$(', '!')
```

Each `CommandArg` has:

```python
class CommandArg(NamedTuple):
    value: str                     # The argument's value (without quotes)
    opening_quote: str = ""       # Opening quote
    closing_quote: str = ""        # Closing quote
    is_io_redir: bool = False     # Whether the arg is IO redirection
```

Key methods and properties:

- `context.completing_command('cmd')` — returns `True` when `arg_index > 0` and the first arg is `cmd`
- `context.command` — returns `self.args[0].raw_value` (the command name including quotes)
- `context.raw_prefix` — prefix including quotes: `f"{self.opening_quote}{self.prefix}{self.closing_quote}"` (or without closing quote if cursor is inside)

### The `RichCompletion` Class

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
| `prefix_len` | `int \| None` | `None` | Length of prefix to replace; `None` = use default prefix len |
| `display` | `str \| None` | `None` | Display text; if set, common prefix stripping is disabled |
| `description` | `str` | `""` | Description shown when completion is selected |
| `style` | `str` | `""` | Style string for prompt-toolkit's `Completion` object |
| `append_closing_quote` | `bool` | `True` | Whether to append closing quote if cursor is after one |
| `append_space` | `bool` | `False` | Whether to append space after completion (after closing quote) |
| `provider` | `str \| None` | `None` | Debug-only tag for `$XONSH_COMPLETER_TRACE` |

#### `append_closing_quote`

When `True` (the default), xonsh appends the closing quote to the completion value when the cursor is positioned right after a closing quote. For example, with `ls "/usr/"⇥`, if the closing quote is `"`, the completion value gets `"` appended.

Carapace sets `append_closing_quote=False` on all its `RichCompletion` objects because it handles quoting in its own `ActionRawValues()` formatter — it wraps values containing special characters in single quotes or raw strings. Letting xonsh also add closing quotes would double-quote the completion.

#### `append_space`

When `True`, xonsh appends a trailing space to the completion after insertion. The space is placed **after** the closing quote (if applicable). This is used by xonsh's built-in completers for commands and single-match candidates.

Carapace does **not** use `append_space`. Instead, it encodes trailing space directly in the `Value` field of its JSON output — adding `" "` to the value when nospace does not match. This is simpler than setting `append_space=True` and avoids interaction with `append_closing_quote`.

#### `prefix_len`

Controls how many characters of the current word are replaced by the completion. When `None` (default), xonsh calculates the default prefix length from the cursor position. When explicitly set, the completer controls the replacement length.

Carapace sets `prefix_len=len(context.raw_prefix)` to ensure the full typed prefix — including any opening quotes — is replaced by the completion value.

### The `add_one_completer` Function

```python
add_one_completer(name, func, pos='start')
```

Registers a completer in xonsh's `OrderedDict` of completers. The `pos` parameter controls insertion order:

| Position | Behavior |
|----------|----------|
| `'start'` | Insert at the beginning (highest priority) |
| `'end'` | Insert at the end (lowest priority) |
| Other string | Insert **after** the completer with that name |

Carapace uses `'start'` to ensure its completer runs before xonsh's built-in completers.

### The `sub_proc_get_output` Function

```python
def sub_proc_get_output(*args, **env_vars: str) -> tuple[bytes, bool]:
```

Runs a subprocess and returns `(stdout_bytes, not_found_bool)`. It:

1. Copies the current xonsh environment (`XSH.env.detype()`)
2. Merges any `env_vars` keyword arguments (overriding existing env vars)
3. Executes the command via `subprocess.run`
4. Returns the raw stdout bytes and a boolean indicating command-not-found

Carapace uses `sub_proc_get_output` to invoke its own `_carapace` subcommand, passing command-line arguments and receiving JSON-formatted completions.

### Xonsh's Completion Pipeline

When the user presses TAB in xonsh:

1. xonsh's parser tokenizes the command line and builds a `CompletionContext`
2. The `Completer` iterates through registered completers in priority order
3. An **exclusive** completer (default) returning non-empty results stops iteration
4. A **non-exclusive** completer's results are merged with others
5. `RichCompletion` objects are processed by prompt-toolkit's completion UI
6. The completion with the matching prefix is inserted into the command line

Key environment variables:
- `$XONSH_COMPLETER_MODE` — `"substring_tier"` (default, case-insensitive substring) or `"prefix"` (prefix-only)
- `$COMPLETION_QUERY_LIMIT` — max completions before truncation warning
- `$XONSH_COMPLETER_TRACE` — debug trace showing which completer produced each result

## The Xonsh Snippet

Generated by `xonsh.Snippet(cmd)` in `internal/shell/xonsh/snippet.go`:

```python
from xonsh.completers.completer import add_one_completer
from xonsh.completers.tools import contextual_command_completer

@contextual_command_completer
def _example_completer(context):
    """carapace completer for example"""
    if context.completing_command('example'):
        from json import loads
        from xonsh.completers.tools import sub_proc_get_output, RichCompletion

        def fix_prefix(s):
            """quick fix for partially quoted prefix completion ('prefix',<TAB>)"""
            return s.translate(str.maketrans('', '', '\'"'))

        output, _ = sub_proc_get_output(
            'example', '_carapace', 'xonsh', *[a.value for a in context.args], fix_prefix(context.prefix)
        )

        try:
            result = {RichCompletion(c["Value"], display=c["Display"], description=c["Description"], prefix_len=len(context.raw_prefix), append_closing_quote=False, style=c["Style"]) for c in loads(output)}
        except:
            result = {}
        if len(result) == 0:
            result = {RichCompletion(context.prefix, display=context.prefix, description='', prefix_len=len(context.raw_prefix), append_closing_quote=False)}
        return result

add_one_completer('example', _example_completer, 'start')
```

### Snippet Walkthrough

**1. Define the completer with `@contextual_command_completer`**

```python
@contextual_command_completer
def _example_completer(context):
```

The decorator ensures this function:
- Only runs during command completion (not Python expression completion)
- Receives a `CommandContext` object directly

**2. Guard with `context.completing_command()`**

```python
if context.completing_command('example'):
```

Checks that the command being completed is `example`. Without this guard, the completer would attempt completion for any command, potentially interfering with other completers.

**3. The `fix_prefix` helper**

```python
def fix_prefix(s):
    """quick fix for partially quoted prefix completion ('prefix',<TAB>)"""
    return s.translate(str.maketrans('', '', '\'"'))
```

**Why this is needed**: When the user types a partially quoted argument like `'prefix` and presses TAB, `context.prefix` may include the opening quote character. Carapace's traversal engine doesn't expect quote characters in the completion prefix — it needs the raw prefix without quotes to match against completion values. `fix_prefix` strips both single and double quote characters from the prefix.

**4. Invoke carapace as a subprocess**

```python
output, _ = sub_proc_get_output(
    'example', '_carapace', 'xonsh', *[a.value for a in context.args], fix_prefix(context.prefix)
)
```

This constructs the carapace command line:
- `example` — the command executable
- `_carapace` — carapace's hidden completion subcommand
- `xonsh` — the target shell format
- `*[a.value for a in context.args]` — the parsed argument values (unquoted, since `CommandArg.value` strips quotes)
- `fix_prefix(context.prefix)` — the current prefix being completed, with quotes stripped

Key point: `context.args` provides `CommandArg` objects where `a.value` is the unquoted argument value. Xonsh's parser has already handled tokenization and quote stripping — no `xargs` is needed.

**5. Parse JSON output into `RichCompletion` objects**

```python
try:
    result = {RichCompletion(c["Value"], display=c["Display"], description=c["Description"], prefix_len=len(context.raw_prefix), append_closing_quote=False, style=c["Style"]) for c in loads(output)}
except:
    result = {}
```

Each completion candidate from carapace's JSON output is converted to a `RichCompletion`:

| JSON field | RichCompletion field | Notes |
|------------|---------------------|-------|
| `Value` | `value` (positional) | The completion value, possibly quoted by carapace |
| `Display` | `display` | Human-readable display text |
| `Description` | `description` | Description shown in the completion menu |
| `Style` | `style` | Xonsh style string (`bg:default fg:ansired bold`) |
| (computed) | `prefix_len` | `len(context.raw_prefix)` — replaces the full typed prefix |
| (hardcoded) | `append_closing_quote` | Always `False` — carapace handles its own quoting |

The `try/except` swallows any JSON parsing errors and falls back to an empty result set.

**6. Fallback on empty results**

```python
if len(result) == 0:
    result = {RichCompletion(context.prefix, display=context.prefix, description='', prefix_len=len(context.raw_prefix), append_closing_quote=False)}
```

When no completions are found, carapace returns a single candidate that is just the current prefix. This prevents xonsh from falling through to other completers that might produce irrelevant results.

**7. Register the completer**

```python
add_one_completer('example', _example_completer, 'start')
```

Registers at the **start** of the completer chain, giving carapace highest priority over built-in completers.

## Value Formatting: `ActionRawValues()`

`xonsh.ActionRawValues(currentWord, meta, values)` in `internal/shell/xonsh/action.go` formats completion values as a JSON array of `richCompletion` objects.

### Implementation

```go
var sanitizer = strings.NewReplacer(
    "\n", ``,
    "\t", ``,
    `'`, `\'`,
)

type richCompletion struct {
    Value       string
    Display     string
    Description string
    Style       string
}

func ActionRawValues(currentWord string, meta common.Meta, values common.RawValues) string {
    vals := make([]richCompletion, len(values))
    for index, val := range values {
        val.Value = sanitizer.Replace(val.Value)

        appendSpace := !meta.Nospace.Matches(val.Value)

        if strings.ContainsAny(val.Value, ` ()[]{}*$?\"|<>&;#`+"`") {
            if strings.Contains(val.Value, `\`) {
                val.Value = fmt.Sprintf("r'%v'", val.Value) // backslash needs raw string
            } else {
                val.Value = fmt.Sprintf("'%v'", val.Value) // regular single-quoted string
            }
        }

        if appendSpace {
            val.Value = val.Value + " "
        }

        vals[index] = richCompletion{
            Value:       val.Value,
            Display:     val.Display,
            Description: val.TrimmedDescription(),
            Style:       convertStyle("bg-default fg-default " + val.Style),
        }
    }
    m, _ := json.Marshal(vals)
    return string(m)
}
```

### Step 1: Sanitize

The sanitizer strips characters that would break the JSON output or xonsh's completion display:

- `\n` → empty — newlines would break the JSON structure
- `\t` → empty — tabs would corrupt display formatting
- `'` → `\'` — escapes single quotes for Python string literals (since values containing special chars are wrapped in single quotes)

### Step 2: Quote Special Characters

Values containing xonsh/shell metacharacters are wrapped in Python string literals to prevent xonsh from interpreting them. The check `strings.ContainsAny(val.Value, ` ()[]{}*$?\"|<>&;#`+"`")` covers:

| Character(s) | Why they're special |
|-------------|---------------------|
| Space ` ` | Word separator — would split the completion into multiple tokens |
| `()[]{}` | Shell syntax — parentheses for subprocess, brackets for lists/dicts |
| `*$?` | Glob/wildcard — `*` is glob, `$` is variable expansion, `?` is single-char glob |
| `\` | Escape character in Python strings — triggers raw string mode |
| `"` | Double quote — string delimiter |
| `|<>&;#` | Pipe, redirect, logical operators, background, comment |

Two quoting modes:

1. **`'value'`** — Single-quoted string: Used when the value contains special characters but **no backslashes**. The sanitizer has already escaped any single quotes to `\'`. Example: `'my file.txt'`

2. **`r'value'`** — Raw string: Used when the value contains **backslashes** (e.g., Windows paths like `C:\Users`). The `r` prefix prevents Python from interpreting `\U`, `\n`, etc. as escape sequences. Example: `r'C:\Users\file'`

### Step 3: Append Space (Nospace Handling)

```go
appendSpace := !meta.Nospace.Matches(val.Value)
if appendSpace {
    val.Value = val.Value + " "
}
```

This is xonsh's mechanism for "nospace" — whether a trailing space should be inserted after the completion. The logic is inverted from the carapace convention:

- `meta.Nospace.Matches(val.Value)` returns `true` when the completion should **NOT** have a trailing space
- So `appendSpace = !meta.Nospace.Matches(val.Value)` means: add a space when nospace does NOT match

The trailing space is added directly to the `Value` field of the JSON output, not via `RichCompletion.append_space`. This is because:

1. Carapace's nospace can be per-candidate (e.g., `ActionValues("a", "b=").NoSpace()` matches `=`)
2. The space is applied **before** quoting — so `file.txt` becomes `'file.txt '` not `file.txt` with a separate space
3. The snippet sets `append_closing_quote=False`, so xonsh won't add closing quotes after the space

**Important**: This means the `Value` field in the JSON already contains the trailing space (if any). The snippet's `RichCompletion` does NOT set `append_space=True` — the space is baked into the value itself.

### Step 4: Build JSON and Serialize

```go
vals[index] = richCompletion{
    Value:       val.Value,
    Display:     val.Display,
    Description: val.TrimmedDescription(),
    Style:       convertStyle("bg-default fg-default " + val.Style),
}
m, _ := json.Marshal(vals)
```

The output is a JSON array:

```json
[
    {"Value": "'file.txt' ", "Display": "file.txt", "Description": "a file", "Style": "bg:default fg:ansigreen"},
    {"Value": "--flag ", "Display": "--flag", "Description": "a flag", "Style": "bg:default fg:ansiyellow"}
]
```

Each field maps to a `RichCompletion` constructor argument in the snippet:

| JSON field | Maps to | Notes |
|------------|---------|-------|
| `Value` | First positional arg | Already contains trailing space if needed |
| `Display` | `display=` | Human-readable text (no quoting, no space) |
| `Description` | `description=` | Shown in completion menu |
| `Style` | `style=` | Xonsh-style color string |

## Style Conversion: `convertStyle()`

`convertStyle(s)` in `internal/shell/xonsh/style.go` converts carapace's style format to xonsh's prompt-toolkit style format. It is called as `convertStyle("bg-default fg-default " + val.Style)` — prepending defaults ensures both foreground and background are always set.

### Output Format

Xonsh style strings use the format:

```
bg:default fg:ansired bold italic underline blink reverse
```

Components are space-separated:
- **Background**: `bg:<color>` or `bg:default`
- **Foreground**: `fg:<color>` or `fg:default`
- **Attributes**: `bold`, `italic`, `underline`, `blink`, `reverse`

### Color Mapping

Named colors are converted to xonsh's `ansi` prefixed form:

| Carapace | Xonsh |
|----------|-------|
| `black` | `ansiblack` |
| `red` | `ansired` |
| `green` | `ansigreen` |
| `yellow` | `ansiyellow` |
| `blue` | `ansiblue` |
| `magenta` | `ansimagenta` |
| `cyan` | `ansicyan` |
| `white` | `ansigray` |
| `bright-black` | `ansibrightblack` |
| `bright-red` | `ansibrightred` |
| ... (bright variants) | `ansibright<name>` |

Hex colors (e.g., `#ff5500`) pass through unchanged — xonsh/prompt-toolkit supports hex natively.

XTerm256 colors (`color0` through `color255`) are mapped to their hex equivalents (e.g., `color0` → `#000000`, `color15` → `#ffffff`, `color196` → `#ff0000`). This is necessary because xonsh/prompt-toolkit does not have built-in `colorN` support.

### Attribute Mapping

| Carapace Attribute | Xonsh Attribute | Notes |
|--------------------|-----------------|-------|
| `Bold` | `bold` | Direct mapping |
| `Dim` | *(workaround)* | Not natively supported; maps foreground to `#808080` if white/default |
| `Italic` | `italic` | Direct mapping |
| `Underlined` | `underline` | Direct mapping |
| `Blink` | `blink` | Direct mapping |
| `Inverse` | `reverse` | Direct mapping |

**Dim workaround**: Xonsh/prompt-toolkit does not have a `dim` attribute. The current workaround (`TODO dim not supported`) sets the foreground to `#808080` (gray) when the foreground is white or default. This is a best-effort approximation — it does not preserve the original color at reduced brightness.

## Edge Cases and How Carapace Handles Them

### 1. Special Characters in Completion Values

**Problem**: Xonsh interprets shell metacharacters (` `, `(`, `)`, `[`, `]`, `{`, `}`, `*`, `$`, `?`, `\`, `"`, `|`, `<`, `>`, `&`, `;`, `` ` ``) in completion values. A value like `my file.txt` would be split into two arguments, or `$(var)` would trigger variable expansion.

**Solution**: Values containing any of these characters are wrapped in Python string literals:
- **No backslashes**: `'value'` — single-quoted Python string. The sanitizer escapes internal single quotes as `\'`.
- **With backslashes**: `r'value'` — raw string. Prevents Python from interpreting `\n`, `\t`, `\U`, etc. as escape sequences. This is critical for Windows paths like `C:\Users\new_folder`.

The snippet sets `append_closing_quote=False` to prevent xonsh from adding its own closing quote on top of carapace's quoting.

### 2. Partially Quoted Prefix Completion

**Problem**: When the user types `'partial` and presses TAB, `context.prefix` contains `'partial` — including the opening quote character. Carapace's traversal engine matches completion values against the raw prefix without quotes, so the quote character causes no matches.

**Solution**: The `fix_prefix` helper in the snippet strips quote characters from `context.prefix`:

```python
def fix_prefix(s):
    """quick fix for partially quoted prefix completion ('prefix',<TAB>)"""
    return s.translate(str.maketrans('', '', '\'"'))
```

The snippet passes `fix_prefix(context.prefix)` as the last argument to carapace.

### 3. Nospace (No Trailing Space)

**Problem**: Some completions should not add a trailing space — for example, `git checkout origin/` should not add a space after `/`, because the user likely wants to continue typing the branch name.

**Solution**: Carapace encodes trailing space directly in the JSON `Value` field:
- When `meta.Nospace.Matches(val.Value)` is `true` (nospace), the value is left as-is (no trailing space)
- When `nospace` does NOT match, `" "` is appended to the value


The snippet does NOT use `RichCompletion.append_space=True`. This is because:
1. `append_space` adds space **after** the closing quote, which doesn't work when carapace has already quoted the value
2. Nospace can be per-candidate (e.g., `ActionValues("a", "b=").NoSpace()` only suppresses space for `b=`)
3. Baking the space into the value ensures it works consistently regardless of quoting mode

### 4. Empty Completion Results

**Problem**: When carapace returns no completions, xonsh falls through to other completers (path completer, alias completer, etc.). This can produce irrelevant results.

**Solution**: The snippet provides a fallback:

```python
if len(result) == 0:
    result = {RichCompletion(context.prefix, display=context.prefix, description='',
                              prefix_len=len(context.raw_prefix), append_closing_quote=False)}
```

This single candidate just re-inserts the current prefix, effectively a no-op that prevents other completers from running (since the `@contextual_command_completer` decorator makes this an exclusive completer).

### 5. Hyphenated Command Names

**Problem**: The snippet function name is derived from the command name by replacing hyphens with double underscores: `strings.ReplaceAll(cmd.Name(), "-", "__")`. This avoids Python identifier issues, since hyphens are not valid in Python identifiers.

For example, `my-command` becomes `_my__command_completer`.

### 6. JSON Parsing Errors

**Problem**: If carapace crashes or produces malformed output, `json.loads()` would raise an exception, breaking the completion.

**Solution**: The snippet wraps JSON parsing in `try/except`:

```python
try:
    result = {RichCompletion(...) for c in loads(output)}
except:
    result = {}
```

This silently swallows any error and falls through to the empty-result fallback, ensuring completion never crashes.

### 7. Style/Color Support

**Problem**: Xonsh uses prompt-toolkit's style format, which differs from carapace's internal style format.

**Solution**: `convertStyle()` maps carapace styles to xonsh styles:
- Named colors get the `ansi` prefix (e.g., `red` → `ansired`)
- Hex colors pass through unchanged
- XTerm256 colors are mapped to hex equivalents
- Background and foreground are always set (defaulting to `bg:default`/`fg:default`)

The style string is always prefixed with `"bg-default fg-default "` to ensure both fields have a value. Missing fields are then overridden by the actual style, producing a complete xonsh style string like `bg:default fg:ansired bold`.

### 8. Messages (Error/Usage)

**Problem**: Xonsh doesn't have a native message display API like zsh's `_message` or elvish's `edit:notify`.

**Solution**: Carapace uses the shared `shell.Value()` dispatch in `internal/shell/shell.go`, which integrates messages as synthetic `RawValue` entries for shells without native message support. For xonsh, messages appear as regular completion candidates styled with `style.Carapace.Error` (red underline) or `style.Carapace.Usage` (italic), making them visible but not insertable.

## Related Skills

- **references/shell.md** — cross-shell feature comparison and shared dispatch
- **references/traverse.md** — the completion engine that produces Actions before formatting
- **references/style.md** — how styles are resolved before shell rendering