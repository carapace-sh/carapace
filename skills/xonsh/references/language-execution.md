# Xonsh Language and Execution Model

In-depth reference for xonsh's Python/shell hybrid language — subprocess execution, capture modes, quoting, string types, environment access, aliases, and how Python and shell syntax interleave.

## The Hybrid Model

Xonsh is a Python-powered shell that blends Python and shell syntax. The key insight is that **every line is Python**, but xonsh extends Python syntax with shell-like constructs for subprocess execution:

| Syntax | Meaning | Returns |
|--------|---------|---------|
| `cmd arg1 arg2` | Subprocess (bare command) | `None` (output to terminal) |
| `$(cmd arg1 arg2)` | Captured subprocess | `str` (stdout) |
| `!(cmd arg1 arg2)` | Captured subprocess (object) | `CommandPipeline` |
| `$[cmd arg1 arg2]` | Uncaptured subprocess | `None` (output to terminal) |
| `![cmd arg1 arg2]` | Hidden captured subprocess | `HiddenCommandPipeline` |

### Python vs Shell Mode

Xonsh determines whether a line is Python or shell based on parsing:

- **Lines starting with a known command** → subprocess mode
- **Lines with Python syntax** → Python mode
- **Ambiguous cases** → Python mode takes precedence

The parser is a modified Python parser that adds subprocess syntax as additional grammar rules.

## Subprocess Execution

### SubprocSpec

The core class for subprocess execution:

```python
class SubprocSpec:
    """A container for specifying how a subprocess command should be executed."""

    def __init__(
        self,
        cmd,
        cls=subprocess.Popen,
        stdin=None,
        stdout=None,
        stderr=None,
        universal_newlines=False,
        close_fds=False,
        captured=False,
        env=None,
    ):
```

**Key attributes:**

| Attribute | Type | Description |
|-----------|------|-------------|
| `cmd` | `list[str]` | Command to be run |
| `args` | `list[str]` | Arguments as originally supplied |
| `alias` | `list/callable/None` | Resolved alias for command |
| `binary_loc` | `str/None` | Path to binary to execute |
| `is_proxy` | `bool` | Whether subprocess is run as proxy |
| `background` | `bool` | Whether subprocess runs in background |
| `threadable` | `bool` | Whether subprocess can run in background thread |
| `pipeline_index` | `int/None` | Index into the pipeline |
| `last_in_pipeline` | `bool` | Whether this is the last in the pipeline |
| `captured_stdout` | `file-like` | Handle to captured stdout |
| `captured_stderr` | `file-like` | Handle to captured stderr |
| `stack` | `list/None` | Stack frame for callable alias |

### Build Process

`SubprocSpec.build()` creates a spec and applies resolution steps in order:

1. `resolve_decorators()` — Apply `DecoratorAlias` modifications
2. `resolve_args_list()` — Flatten argument lists
3. `resolve_redirects()` — Handle `>`, `>>`, `<`, `2>&1`, etc.
4. `resolve_alias()` — Look up command in alias table
5. `resolve_binary_loc()` — Find binary in PATH
6. `resolve_auto_cd()` — Auto-cd to directory if command is a directory
7. `resolve_executable_commands()` — Check if command is executable
8. `resolve_alias_cls()` — Determine alias class (callable vs binary)
9. `resolve_stack()` — Set up stack frame for callable aliases

### Capture Modes

| Mode | Syntax | `captured` Value | Returns | Output |
|------|--------|-----------------|---------|--------|
| Uncaptured | `cmd` or `$[cmd]` | `False` | `None` | Goes to terminal |
| Stdout capture | `$(cmd)` | `"stdout"` | `str` | Captured, hidden |
| Object capture | `!(cmd)` | `"object"` | `CommandPipeline` | Captured, visible |
| Hidden object | `![cmd]` | `"hiddenobject"` | `HiddenCommandPipeline` | Captured, hidden |

**Key differences:**

- `$(cmd)` returns the stdout text as a string (like bash's `$()`)
- `!(cmd)` returns a `CommandPipeline` object with `.output`, `.returncode`, `.proc`, etc.
- `![cmd]` is like `!(cmd)` but doesn't block (useful for background processing)
- `$[cmd]` is like bare `cmd` but explicitly marks it as subprocess syntax

### Pipeline Execution

```python
def run_subproc(cmds, captured=False, envs=None, in_boolop=False):
    """Runs a subprocess, in its many forms."""
```

1. `cmds_to_specs()` — Convert command lists to `SubprocSpec` objects
2. `_run_specs()` — Execute specs via `_run_command_pipeline()`
3. Return based on capture mode

### Pipe and Redirect Handling

#### Redirect Syntax

| Syntax | Meaning |
|--------|---------|
| `cmd > file` | Redirect stdout to file |
| `cmd >> file` | Append stdout to file |
| `cmd < file` | Redirect stdin from file |
| `cmd 2> file` | Redirect stderr to file |
| `cmd 2>&1` | Redirect stderr to stdout |
| `cmd >&2` | Redirect stdout to stderr |
| `cmd e>p` | Redirect stderr into pipe |
| `cmd a>p` | Redirect all (stdout+stderr) into pipe |

#### Redirect Regex

```python
_REDIR_REGEX = r"(o(?:ut)?|e(?:rr)?|a(?:ll)?|&?\d?)(>?>|<)(o(?:ut)?|e(?:rr)?|a(?:ll)?|&?\d?)$"
```

#### Pipe Channel Setup

For `|` redirects:
1. Creates `PipeChannel.from_pipe()`
2. Sets upstream `stdout` to pipe's write fd
3. Sets downstream `stdin` to pipe's read fd
4. Handles `e>p` and `a>p` sentinels

### Subprocess Expressions

| Syntax | Name | Description |
|--------|------|-------------|
| `$(...)` | Command substitution | Captures stdout as string |
| `!(...)` | Object capture | Returns `CommandPipeline` |
| `@[...]` | Python sub-expression | Evaluates Python inside subprocess context |
| `$[...]` | Uncaptured subprocess | Output goes to terminal |
| `![...]` | Hidden object capture | Returns `HiddenCommandPipeline` |
| `@$(...)` | String interpolation | Embeds subprocess output in Python string |
| `@(...)` | Python evaluation | Evaluates Python expression |

## Quoting and String Types

Xonsh supports all Python string types plus some extensions:

### String Literals

| Type | Syntax | Description |
|------|--------|-------------|
| Single-quoted | `'...'` | Standard Python string |
| Double-quoted | `"..."` | Standard Python string |
| Triple-quoted | `'''...'''` or `"""..."""` | Multi-line string |
| Raw string | `r'...'` or `r"..."` | Backslashes are literal |
| F-string | `f'...'` or `f"..."` | Formatted string literal |
| P-string | `p'...'` or `p"..."` | Path string (xonsh extension) |
| PR-string | `pr'...'` or `pr"..."` | Raw path string (xonsh extension) |

### Path Strings (xonsh extension)

Path strings (`p"..."`, `pr"..."`) are xonsh-specific string types for filesystem paths:

- `p"..."` — Path string, behaves like a regular string but signals path intent
- `pr"..."` — Raw path string, backslashes are literal (useful for Windows paths)
- The `p` prefix is stripped during evaluation
- `pr` is converted to `r` during evaluation

### String Quoting in Completions

When completing paths, xonsh's path completer applies per-path quoting:

- **Plain names** (no special chars) → unquoted
- **Names with spaces, `$`, `\`** → quoted with `'` or `"` (whichever doesn't appear in the name)
- **Names with `$` or `\`** → promoted to raw string (`r'...'`)
- **Names with control chars** (`\n`, `\t`) → non-raw string with escape sequences
- **Raw strings ending with `\`** → doubled backslash to avoid syntax error

### Quote Handling in CommandContext

The `CommandContext` provides quote information for the current argument:

```python
opening_quote: str = ""      # e.g. "'", '"', 'r"'
closing_quote: str = ""      # e.g. "'", '"'
is_after_closing_quote: bool = False  # Cursor is after a closed string
```

- `raw_prefix` includes the opening quote (and closing quote if `is_after_closing_quote`)
- `prefix` is the text without quotes
- Completers must handle both quoted and unquoted arguments

## Environment Access

### Variable Syntax

| Syntax | Meaning | Example |
|--------|---------|---------|
| `$VAR` | Environment variable | `$HOME` → `/home/user` |
| `${VAR}` | Environment variable (braced) | `${HOME}` → `/home/user` |
| `${{VAR}}` | Python dict lookup | `${{mydict['key']}}` |
| `@VAR` | Xonsh built-in reference | `@aliases` |

### Environment Variable Types

Xonsh's `Env` class supports typed environment variables with validation, conversion, and detyping:

```python
class Env(cabc.MutableMapping):
    """A xonsh environment, whose variables have limited typing."""
```

**Type system:**

| Type Name | Validator | Converter | Detyper |
|-----------|-----------|-----------|---------|
| `bool` | `is_bool` | `to_bool` | `bool_to_str` |
| `str` | `is_string` | `ensure_string` | `ensure_string` |
| `path` | `is_path` | `str_to_path` | `path_to_str` |
| `env_path` | `is_env_path` | `str_to_env_path` | `env_path_to_str` |
| `abs_path` | `is_path` | `str_to_abs_path` | `abs_path_to_str` |
| `float` | `is_float` | `float` | `str` |
| `int` | `is_int` | `int` | `str` |

### VarPattern (Dynamic Type Matching)

```python
class VarPattern:
    """Pattern rule for dynamic env var typing."""
    def __init__(self, pattern, var_type, exclude=None):
```

Example:
```python
XONSH_ENV_PATTERN_PATH = VarPattern(r"\w*PATH$", "env_path")
XONSH_ENV_PATTERN_DIRS = VarPattern(r"\w*DIRS$", "env_path", exclude=["JUPYTER_PLATFORM_DIRS"])
```

Any environment variable whose name matches the pattern gets the specified type handling.

### EnvPath

```python
class EnvPath(cabc.MutableSequence):
    """A list of strings representing paths with expansion support."""
```

Used for `PATH`, `CDPATH`, `XONSH_COMPLETER_DIRS`, etc. Supports:
- Path expansion (`~`, `$VAR`)
- `add(data, front=False, replace=False)`
- `prepend(value)`, `append(value)`, `remove(value)`
- Fires `on_envvar_change` events on mutation

### LsColors

```python
class LsColors(cabc.MutableMapping):
    """Helps convert to/from $LS_COLORS format, respecting xonsh color style."""
```

Maps file type keys to color tuples. Fires `on_lscolors_change` events on mutation.

### `detype()` and `detype_all()`

- `detype()` — Returns detyped variables for subprocess use (only explicitly set vars)
- `detype_all()` — Returns all detyped env vars including defaults

Detyping converts Python types back to string representations suitable for passing to subprocesses.

## Alias System

### Alias Types

| Type | Description | Example |
|------|-------------|---------|
| String alias | Simple text replacement | `aliases['ll'] = 'ls -la'` |
| List alias | Multi-word command | `aliases['gs'] = ['git', 'status']` |
| Callable alias | Python function | `aliases['hello'] = lambda args: print('hello')` |
| ExecAlias | Executable alias | `aliases['py'] = ExecAlias('python3')` |

### Callable Aliases

```python
def my_alias(args, stdin=None, stdout=None, stderr=None, spec=None, stack=None, **kwargs):
    """Custom callable alias."""
    print(f"Args: {args}", file=stdout)
```

Callable aliases receive:
- `args` — list of arguments
- `stdin`, `stdout`, `stderr` — file-like objects
- `spec` — `SubprocSpec` for the command
- `stack` — call stack frame

### DecoratorAlias

```python
class DecoratorAlias:
    """Decorator alias base class."""

    def decorate_spec(self, spec):
        """Modify spec immediately after modifier added."""

    def decorate_spec_pre_run(self, pipeline, spec, spec_num):
        """Modify spec before run."""
```

### SpecAttrDecoratorAlias

```python
class SpecAttrDecoratorAlias(DecoratorAlias):
    """Decorator Alias for spec attributes."""

    def __init__(self, set_attributes: dict, descr="", name=""):
```

Sets attributes on the `SubprocSpec` before execution.

## Word Splitting and Glob Expansion

### Word Splitting

Unlike bash, xonsh does **not** perform word splitting on unquoted variable expansions by default. `$VAR` expands to a single string, not multiple words.

To get bash-like word splitting, use `${{VAR}}` or explicit splitting:

```python
# No word splitting
cmd $VAR        # $VAR is a single argument

# Explicit word splitting
cmd ${{VAR.split()}}  # Split on whitespace
```

### Glob Expansion

Glob patterns (`*`, `?`, `[...]`) are expanded in subprocess context:

```python
ls *.py         # Glob expanded
echo @(*.py)    # Python glob via @()
```

### Tilde Expansion

| Syntax | Result |
|--------|--------|
| `~` | `$HOME` |
| `~/path` | `$HOME/path` |
| `~user` | Home directory of `user` |

Tilde expansion is **not** performed inside raw strings (`r'~'` is literal).

## Process Model

### Foreground/Background

| Syntax | Meaning |
|--------|---------|
| `cmd` | Foreground (waits for completion) |
| `cmd &` | Background (returns immediately) |
| `cmd1 && cmd2` | Run cmd2 only if cmd1 succeeds |
| `cmd1 \|\| cmd2` | Run cmd2 only if cmd1 fails |
| `cmd1 ; cmd2` | Run sequentially |
| `cmd1 \| cmd2` | Pipe stdout of cmd1 to stdin of cmd2 |

### Process Groups

On POSIX, xonsh sets up process groups for pipeline management:

```python
def prep_preexec_fn(self, kwargs, pipeline_group=None):
    """Prepares the 'preexec_fn' keyword argument"""
```

- Each pipeline gets its own process group
- The first process in the pipeline is the group leader
- Job control is handled via `SIGTTIN`, `SIGTTOU`, `SIGTSTP`

### Threaded Subprocesses

When `THREAD_SUBPROCS` is `True`, the last command in a pipeline runs in a background thread, allowing the prompt to remain responsive.

## XSH Built-in Object

The `XSH` object is xonsh's central built-in, providing access to:

| Attribute | Type | Description |
|-----------|------|-------------|
| `XSH.env` | `Env` | Environment variables |
| `XSH.completers` | `OrderedDict` | Registered completers |
| `XSH.commands_cache` | `CommandsCache` | Cache of available commands |
| `XSH.builtins` | `BuiltIns` | Shell builtins |
| `XSH.shell` | `BaseShell` | Current shell instance |
| `XSH.execer` | `Execer` | Code execution engine |
| `XSH.ctx` | `dict` | Current execution context |
| `XSH.last` | `CommandPipeline` | Last executed pipeline |
| `XSH.history` | `History` | Command history |

### Accessing XSH

```python
from xonsh.built_ins import XSH

# Get an environment variable
path = XSH.env.get("PATH")

# Access completers
completers = XSH.completers

# Get command cache
for cmd, (path, is_alias) in XSH.commands_cache.iter_commands():
    print(cmd, path, is_alias)
```
