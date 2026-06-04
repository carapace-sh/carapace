# Xonsh Startup, Configuration, and Events

In-depth reference for xonsh's startup process, RC files, xontrib system, event system, and environment variable configuration.

## Startup Process

### Initialization Flow

1. **TTY Setup** — `_setup_controlling_terminal()` runs first to ensure xonsh is in the foreground process group
2. **Argument Parsing** — `premain()` parses command-line arguments, determines execution mode
3. **Environment Setup** — `start_services()` creates the Execer, loads the xonsh session (`XSH.load()`), installs import hooks
4. **RC File Loading** — `_load_rc_files()` sources run control files
5. **Xontrib Autoloading** — External xontribs are loaded from entry points
6. **Shell Creation** — Shell instance is created (PTK or readline)

### Execution Modes

| Mode | Trigger | Description |
|------|---------|-------------|
| `interactive` | Default (no args, TTY) | Interactive shell with prompt |
| `single_command` | `xonsh -c "cmd"` | Execute a single command |
| `script_from_file` | `xonsh script.xsh` | Execute a script file |
| `script_from_stdin` | Piped input | Execute from stdin |
| `source` | `source script.xsh` | Source a script in current context |

### Command-Line Options

| Option | Description |
|--------|-------------|
| `-c COMMAND` | Execute a single command |
| `-i, --interactive` | Force interactive mode |
| `-l, --login` | Run as login shell |
| `--no-rc` | Don't load any RC files |
| `--rc FILE...` | Load only specified RC files |
| `--no-env` | Don't inherit parent environment |
| `--save-origin-env` | Save origin environment to file |
| `--load-origin-env` | Load saved origin environment |
| `-D NAME=VAL` | Define environment variable |

## RC (Run Control) Files

### Default RC File Locations

| File | When Loaded | Description |
|------|-------------|-------------|
| `$XONSH_SYS_CONFIG_DIR/xonshrc` | Always | System-wide config (e.g., `/etc/xonshrc`) |
| `$XONSH_CONFIG_DIR/xonsh/rc.xsh` | Always | User XDG config (`~/.config/xonsh/rc.xsh`) |
| `~/.xonshrc` | Interactive only | Home directory RC file |

### RC Directories

| Directory | Description |
|-----------|-------------|
| `$XONSH_SYS_CONFIG_DIR/rc.d` | System-wide drop-in configs |
| `$XONSH_CONFIG_DIR/rc.d` | User drop-in configs (`~/.config/xonsh/rc.d/`) |

RC directories are scanned for `.xsh` and `.py` files, sorted lexicographically.

### RC File Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `$XONSHRC` | `['$XONSH_SYS_CONFIG_DIR/xonshrc', '$XONSH_CONFIG_DIR/xonsh/rc.xsh', '~/.xonshrc']` | List of RC file locations |
| `$XONSHRC_DIR` | `['$XONSH_SYS_CONFIG_DIR/rc.d', '$XONSH_CONFIG_DIR/rc.d']` | Directories containing RC files |

### RC File Loading Rules

- `--no-rc`: Skip all RC files
- `--rc FILE...`: Load only specified files (replaces defaults)
- In non-interactive mode, `~/.xonshrc` is excluded
- RC files are executed in the xonsh execution context (Python + shell syntax)

### Typical rc.xsh

```python
# Aliases
aliases['ll'] = 'ls -la'
aliases['gs'] = 'git status'

# Environment
$PATH.prepend('~/bin')
$EDITOR = 'vim'

# Completion
from xonsh.completers.completer import add_one_completer
add_one_completer('my_completer', my_completer_func, 'start')

# Xontribs
xontrib load autovox
xontrib load z

# Prompt customization
$PROMPT = '{env_name}{BOLD_GREEN}{user}@{hostname}{BOLD_BLUE} {cwd}{RESET} $ '

# Event handlers
@events.on_chdir
def _source_env_xsh(olddir, newdir, **_):
    env_file = pathlib.Path(newdir) / 'env.xsh'
    if env_file.exists():
        source @(env_file)
```

## Xontrib System

Xontribs (xonsh contributions) are plugins that extend xonsh's functionality.

### The `xontrib` Command

```python
class XontribAlias(ArgParserAlias):
    """Manage xonsh extensions"""

    # Subcommands:
    # xontrib load   — Load xontrib(s)
    # xontrib unload — Unload xontrib(s)
    # xontrib reload — Reload xontrib(s)
    # xontrib list   — List xontribs (default)
    # xontrib info   — Show xontrib info
```

### Xontrib Discovery Order

When loading a xontrib, xonsh searches in this order:

1. **Autoloaded cache** — `XSH.builtins.autoloaded_xontribs` (common interactive path)
2. **Entry-point lookup** — For `--no-rc` or `$XONTRIBS_AUTOLOAD_DISABLED`
3. **Legacy namespace** — `xontrib.<name>` namespace package
4. **Top-level module** — `import <name>` as fallback

### Xontrib Package Layout

```
xontrib/
  my_xontrib.py        # Python file
  my_xontrib.xsh       # Xonsh script file
  my_xontrib/
    __init__.py         # Package style
```

### Writing a Xontrib

#### Method 1: Using `__all__`

```python
# xontrib/my_xontrib.py
__all__ = ['my_alias', 'my_env_var']

def my_alias(args, stdout=None):
    print("Hello from my alias!", file=stdout)

my_env_var = "some_value"
```

#### Method 2: Using `_load_xontrib_()`

```python
# xontrib/my_xontrib.py
def _load_xontrib_(xsh):
    """Called when xontrib is loaded. Receives XSH."""
    return {
        'my_alias': my_alias_func,
        'my_env_var': 'value',
    }
```

#### Method 3: With Cleanup

```python
# xontrib/my_xontrib.py
def _load_xontrib_(xsh):
    # Register event handler
    @xsh.builtins.events.on_precommand
    def _my_handler(cmd, **_):
        pass
    return {}

def _unload_xontrib_(xsh):
    """Called when xontrib is unloaded. Clean up handlers, aliases, etc."""
    # Remove event handlers, env vars, aliases, completers
    pass
```

#### Package as Entry Point

Add to `pyproject.toml`:

```toml
[project.entry-points."xonsh.xontribs"]
my_xontrib = "my_xontrib_module"
```

### Xontrib Meta Class

```python
class Xontrib(NamedTuple):
    module: str                                    # Path to the xontrib module
    distribution: "Distribution | None" = None    # Distribution metadata
```

| Property | Returns | Description |
|----------|---------|-------------|
| `get_description()` | `str` | Distribution summary or module docstring |
| `url` | `str` | Home-page from distribution metadata |
| `license` | `str` | License from distribution metadata |
| `is_loaded` | `bool` | `True` if module is in `sys.modules` |
| `is_auto_loaded` | `bool` | `True` if in `autoloaded_xontribs` |

## Event System

Xonsh's event system allows loose coupling between components. Events are created dynamically on first access.

### Event Species

| Species | Class | Behavior |
|---------|-------|---------|
| Default | `Event` | Handlers called immediately; returns list of return values; supports scatter-gather |
| Load | `LoadEvent` | Each handler called exactly once (after event fires OR handler registers, whichever is later); no return values |

### EventManager

```python
class EventManager:
    """Container singleton for all events."""
```

| Method | Description |
|--------|-------------|
| `register(func)` | Decorator that extracts name and doc from function |
| `doc(name, docstring)` | Apply documentation to an event |
| `transmogrify(name, species)` | Convert event from one species to another |
| `exists(name)` | Check if event exists without creating it |
| `handlers(name=None)` | Return dict of events and their handlers |

### Handler Registration

```python
# Decorator style
@events.on_precommand
def my_handler(cmd, **_):
    pass

# Direct call
events.on_precommand(my_handler)

# With validator
@events.on_completer_filter
def _only_bash_for_some(completer, command, context, **_):
    return command in {'kubectl', 'docker'} if completer == 'bash' else True

@events.on_completer_filter.validator
def _filter_validator(completer, command, context, **_):
    return True
```

**Handler requirements:**
- **MUST** have `**kwargs` parameter for future-proofing (enforced in debug mode)
- Named module-level functions use `(__module__, __qualname__)` as uniqueness key (enables replacement on reload)
- Closures/lambdas use `id()` for uniqueness

### Built-in Events

#### Directory Events

**`on_chdir(olddir: str, newdir: str) -> None`** — Fired when the current directory changes.

#### Command Events

| Event | Signature | When Fired |
|-------|-----------|------------|
| `on_precommand` | `(cmd: str) -> None` | Before command execution |
| `on_postcommand` | `(cmd: str, rtn: int, out: str \| None, ts: list) -> None` | After command execution |
| `on_command_not_found` | `(cmd: list[str]) -> list[str] \| tuple[str, ...] \| dict \| None` | When command is not found; return replacement command |
| `on_pre_spec_run` | `(spec: SubprocSpec) -> None` | Before any command runs |
| `on_pre_spec_run_<cmd>` | `(spec: SubprocSpec) -> None` | Before specific command runs |
| `on_post_spec_run` | `(spec: SubprocSpec) -> None` | After any command starts |
| `on_post_spec_run_<cmd>` | `(spec: SubprocSpec) -> None` | After specific command starts |
| `on_transform_command` | `(cmd: str) -> str \| None` | Transform command before execution |

#### Prompt Events

| Event | Signature | When Fired |
|-------|-----------|------------|
| `on_pre_prompt` | `() -> None` | Before showing prompt |
| `on_pre_prompt_format` | `() -> None` | Before prompt is formatted |
| `on_post_prompt` | `() -> None` | After prompt returns |

#### PTK Events

**`on_ptk_create(prompter, history, completer, bindings)`** — After PTK initialization. Use to customize key bindings, completer, or other PTK components.

#### Completion Events

**`on_completer_filter(completer: str, command: str, context: CompletionContext) -> bool | None`** — Before each completer is invoked. Return `False` to veto.

#### Environment Events

| Event | Signature | When Fired |
|-------|-----------|------------|
| `on_envvar_change` | `(name: str, oldvalue: Any, newvalue: Any) -> None` | After env var changes |
| `on_envvar_new` | `(name: str, value: Any) -> None` | After new env var created |
| `on_lscolors_change` | `(key: str, oldvalue: Any, newvalue: Any) -> None` | After LS_COLORS value changes |

#### Lifecycle Events

| Event | When Fired |
|-------|------------|
| `on_init` | After initialization finished |
| `on_pre_rc` | Before RC files loaded |
| `on_post_rc` | After RC files loaded |
| `on_pre_cmdloop` | Before command loop starts |
| `on_post_cmdloop` | After command loop finishes |
| `on_exit` | Before shell exits |
| `on_xontribs_loaded` | After external xontribs loaded |

#### Import Events

| Event | When Fired |
|-------|------------|
| `on_import_pre_find_spec` | Before finding a module spec |
| `on_import_post_find_spec` | After finding a module spec |
| `on_import_pre_create_module` | Before creating a module |
| `on_import_post_create_module` | After creating a module |
| `on_import_pre_exec_module` | Before executing a module |
| `on_import_post_exec_module` | After executing a module |

### Custom Events

```python
events.doc('myxontrib_on_spam', """
myxontrib_on_spam(can: Spam) -> bool?

Fired in case of spam. Return True if it's been eaten.
""")
```

### Debugging Events

```python
# Print all events and handlers
print(events)

# Get handlers programmatically
events.handlers()
events.handlers('on_precommand')
```

## Environment Variables

### Configuration Directories

| Variable | Default | Description |
|----------|---------|-------------|
| `XONSH_CONFIG_DIR` | `$XDG_CONFIG_HOME/xonsh` | User configuration directory |
| `XONSH_SYS_CONFIG_DIR` | (platform-specific) | System-wide config directory |
| `XONSH_DATA_DIR` | `$XDG_DATA_HOME/xonsh` | Data files (history, completers) |
| `XONSH_CACHE_DIR` | (platform-specific) | Cache directory |

### Startup/Mode Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `XONSH_INTERACTIVE` | `True` | Running interactively |
| `XONSH_LOGIN` | `False` | Running as login shell |
| `XONSH_MODE` | `'interactive'` | Current mode |
| `XONSH_SOURCE` | `''` | Path to currently executing script |
| `XONSH_DEBUG` | `0` | Debug level (0-2+) |
| `XONSH_SHOW_TRACEBACK` | `False` | Show tracebacks on errors |
| `XONSH_TRACEBACK_LOGFILE` | `None` | Logfile for tracebacks |

### Shell/Prompt Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SHELL_TYPE` | `'prompt_toolkit'` | Shell type (`prompt_toolkit` or `readline`) |
| `PROMPT` | (default prompt) | Left prompt template |
| `RIGHT_PROMPT` | `""` | Right-side prompt |
| `BOTTOM_TOOLBAR` | `None` | Bottom toolbar |
| `TITLE` | `""` | Terminal title |
| `MULTILINE_PROMPT` | `""` | Continuation prompt |
| `VI_MODE` | `False` | Enable vi editing mode |
| `XONSH_AUTOPAIR` | `False` | Auto-close brackets and quotes |
| `XONSH_COLOR_STYLE` | `'default'` | Color theme |
| `XONSH_STYLE_OVERRIDES` | `{}` | Pygments style overrides |

### Completion Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `XONSH_COMPLETER_MODE` | `"substring_tier"` | Filter mode |
| `XONSH_COMPLETER_TRACE` | `False` | Debug trace |
| `XONSH_COMPLETER_DIRS` | (platform-specific) | Completer search paths |
| `XONSH_COMPLETER_EMOJI_PREFIX` | `None` | Emoji completion trigger |
| `XONSH_COMPLETER_SYMBOLS_PREFIX` | `None` | Symbol completion trigger |
| `COMPLETION_QUERY_LIMIT` | `1000` | Max completions |
| `COMPLETIONS_DISPLAY` | `"multi"` | Display style |
| `COMPLETION_MODE` | `"default"` | Tab behavior |
| `COMPLETIONS_CONFIRM` | `True` | Enter confirms completion |
| `COMPLETIONS_MENU_ROWS` | `5` | Menu rows |
| `COMPLETION_IN_THREAD` | `False` | Background thread |
| `UPDATE_COMPLETIONS_ON_KEYPRESS` | `False` | Complete on keypress |
| `AUTO_SUGGEST_IN_COMPLETIONS` | `False` | Auto-suggest in completions |
| `XONSH_PROMPT_AUTO_SUGGEST` | `True` | Auto-suggest from history |
| `BASH_COMPLETIONS` | (platform-specific) | Bash completion script paths |
| `COMPLETE_DOTS` | `"matching"` | Dot completion behavior |
| `SUBSEQUENCE_PATH_COMPLETION` | `True` | Subsequence path matching |
| `FUZZY_PATH_COMPLETION` | `True` | Fuzzy path matching |
| `SUGGEST_THRESHOLD` | `3` | Levenshtein distance threshold |
| `SUGGEST_COMMANDS` | `True` | Suggest similar commands |
| `SUGGEST_MAX_NUM` | `5` | Max command suggestions |
| `CMD_COMPLETIONS_SHOW_DESC` | `False` | Show path in command descriptions |
| `COMPLETIONS_BRACKETS` | `True` | Include brackets in Python completions |
| `ALIAS_COMPLETIONS_OPTIONS_BY_DEFAULT` | `False` | Show options without `-` prefix |
| `ALIAS_COMPLETIONS_OPTIONS_LONGEST` | `False` | Show only longest option variant |

### History Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `XONSH_HISTORY_SIZE` | `8128` | History size |
| `XONSH_HISTORY_FILE` | (platform-specific) | History file location |
| `XONSH_HISTORY_MATCH_ANYWHERE` | `False` | Match history anywhere |
| `XONSH_STORE_STDOUT` | `False` | Store stdout in history |
| `XONSH_STORE_STDIN` | `False` | Store stdin in history |

### Subprocess Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `XONSH_CAPTURE_ALWAYS` | `False` | Always capture subprocess output |
| `THREAD_SUBPROCS` | (platform-specific) | Run subprocesses in threads |
| `XONSH_SUBPROC_TRACE` | `False` | Trace subprocess execution |
| `XONSH_SUBPROC_CMD_RAISE_ERROR` | `False` | Raise error on command failure |
| `XONSH_ENCODING` | `sys.getdefaultencoding()` | Subprocess encoding |
| `XONSH_ENCODING_ERRORS` | `'surrogateescape'` | Encoding error handling |

### Environment Pattern Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `XONSH_ENV_PATTERN_PATH` | `r'\w*PATH$'` | Pattern for `env_path` typed vars |
| `XONSH_ENV_PATTERN_DIRS` | `r'\w*DIRS$'` | Pattern for `env_path` typed vars (excludes `JUPYTER_PLATFORM_DIRS`) |

### Cache Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `XONSH_CACHE_SCRIPTS` | `True` | Cache scripts |
| `XONSH_CACHE_EVERYTHING` | `False` | Cache everything |

### Xontrib Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `XONTRIBS_AUTOLOAD_DISABLED` | `False` | Disable auto-loading xontribs |

## Configuration Wizards

- **`xonfig web`** — Web-based configuration UI at `http://localhost:8421`
- **`xonfig wizard`** — Interactive Q&A configuration wizard
- **`xonfig info`** — Show current configuration
