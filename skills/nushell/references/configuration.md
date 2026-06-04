# Nushell Configuration and Startup

In-depth reference for nushell's configuration system — `$env.config`, startup files, autoload directories, and where to set up completions.

## Startup File Loading Order

Nushell loads configuration files in this specific order:

| Step | Stage | Action |
|------|-------|--------|
| 0 | misc | Sets internal defaults via Rust implementation |
| 1 | main | Inherits initial environment from calling process |
| 2 | main | Gets configuration directory |
| 3 | main | Creates initial `$env.NU_LIB_DIRS` (empty by default) |
| 4 | main | Creates `$NU_LIB_DIRS` const (includes scripts dir and completions dir) |
| 5 | main | Creates initial `$env.NU_PLUGIN_DIRS` (empty by default) |
| 6 | main | Creates `$NU_PLUGIN_DIRS` const |
| 7 | main | Initializes SQLite database for `stor` commands |
| 8 | main | Processes commandline arguments |
| 9 | main | Gets paths to `env.nu` and `config.nu` |
| 10 | main | Applies `--include-path` (-I) flag if used |
| 11 | main | Loads `$env.config` values from internal defaults |
| 12 | main | Converts PATH from string to list |
| 13 | stdlib | Loads Standard Library into virtual filesystem |
| 14 | stdlib | Parses and evaluates `std/prelude` |
| 15 | main | Generates `$nu` record constant |
| 16 | main | Loads plugins from `--plugin` flag |
| 17 | repl | Sets default prompt-related environment variables |
| 18 | config files | Processes plugin signatures from `plugin.msgpackz` |
| 19 | config files | Creates config directory on first launch |
| 20 | config files | Creates empty `env.nu` and `config.nu` on first launch |
| 21 | config files | Loads `default_env.nu` |
| 22 | config files | Converts PATH to list in `env.nu` |
| 23 | config files | Loads user's `env.nu` |
| 24 | config files | Loads `$env.config` from `default_config.nu` |
| 25 | config files | Loads user's `config.nu` |
| 26 | config files | Loads `login.nu` (login shell only) |
| 27 | config files | Loads vendor autoload directories (alphabetically) |
| 28 | config files | Loads user autoload directories (alphabetically) |
| 29 | repl | Shows banner if configured, enters REPL |

## Configuration Files

| File | When Loaded | Typical Contents |
|------|------------|------------------|
| `env.nu` | Always (before `config.nu`) | Environment variables, PATH modifications |
| `config.nu` | Always (after `env.nu`) | `$env.config` settings, aliases, custom commands, completions |
| `login.nu` | Login shell only | Login-specific commands |
| `default_env.nu` | Before user `env.nu` | Default environment values |
| `default_config.nu` | Before user `config.nu` | Default `$env.config` values |

### Config Directory Locations

| Platform | Default Path |
|----------|-------------|
| macOS | `~/Library/Application Support/nushell` |
| Linux | `~/.config/nushell` |
| Windows | `C:\Users\<user>\AppData\Roaming\nushell` |

Override with `XDG_CONFIG_HOME` environment variable.

### Quick Commands

```nu
config nu          # Open config.nu in editor
config env         # Open env.nu in editor
config nu --doc    # View available settings with documentation
```

## `$env.config`

The primary mechanism for changing nushell's behavior. Key properties:

- **Not inherited** from parent process — populated by nushell with defaults
- **Not exported** to child processes
- **Must be modified incrementally** — overwriting the entire record resets all other settings

### Correct Modification

```nu
$env.config.show_banner = false
$env.config.completions.algorithm = "fuzzy"
```

### Wrong Modification

```nu
$env.config = { show_banner: false }  # Resets ALL other settings!
```

For record-type keys, overwrite the entire record with all values:

```nu
$env.config.completions = {
    algorithm: "fuzzy"
    sort: "smart"
    case_sensitive: false
    external: {
        enable: true
        max_results: 100
        completer: $my_completer
    }
}
```

## Completions Configuration

```nu
$env.config.completions = {
    algorithm: "prefix"    # or "substring" or "fuzzy"
    sort: "alphabetical"  # or "smart"
    case_sensitive: false
    external: {
        enable: true
        max_results: 100
        completer: $completer
    }
}
```

### Completion Algorithm

| Algorithm | Description |
|-----------|-------------|
| `prefix` | Only matches if haystack starts with needle |
| `substring` | Matches if needle appears anywhere in haystack |
| `fuzzy` | Matches if all needle characters appear in order (uses `nucleo_matcher` for scoring) |

### Sort

| Option | Description |
|--------|-------------|
| `alphabetical` | Sort alphabetically |
| `smart` | Sort based on match algorithm relevance |

### External Completer

```nu
$env.config.completions.external = {
    enable: true
    max_results: 100
    completer: $my_completer
}
```

| Setting | Default | Description |
|---------|---------|-------------|
| `enable` | `true` | Whether external completers are used |
| `max_results` | `100` | Maximum number of suggestions from external completer |
| `completer` | `null` | Closure to invoke for external completion |

## Menu Configuration

```nu
$env.config.menus ++= [{
    name: completion_menu
    only_buffer_difference: false
    marker: "| "
    type: {
        layout: columnar
        columns: 4
        col_width: 20
        col_padding: 2
    }
    style: {
        text: green
        selected_text: green_reverse
        description_text: yellow
    }
}]
```

See [references/reedline.md](references/reedline.md) for full menu configuration details.

## Keybindings Configuration

```nu
$env.config.keybindings ++= [{
    name: completion_menu
    modifier: control
    keycode: char_t
    mode: emacs
    event: { send: menu name: completion_menu }
}]
```

| Field | Description |
|-------|-------------|
| `name` | Unique keybinding name |
| `modifier` | `none`, `control`, `alt`, `shift`, `control_alt`, `control_shift`, etc. |
| `keycode` | The key (e.g., `char_t`, `enter`, `tab`, `backspace`) |
| `mode` | `emacs`, `vi_insert`, `vi_normal` (single or list) |
| `event` | Action: `{ send: ... }`, `{ edit: ... }`, `{ until: [...] }`, or `null` to disable |

### Event Types

| Type | Example | Description |
|------|---------|-------------|
| Send | `{ send: menu name: completion_menu }` | Send event to Reedline |
| Edit | `{ edit: insertstring, value: "text" }` | Edit command |
| Until | `{ until: [{ send: menu name: completion_menu }, { send: menunext }] }` | Chain events (stops on first success) |
| Null | `null` | Disable keybinding |

## Autoload Directories

### Vendor Autoload

`$nu.vendor-autoload-dirs` — intended for vendors' and package managers' startup files. Loaded alphabetically.

### User Autoload

`$nu.user-autoload-dirs` — good for modularizing configuration. Loaded alphabetically.

Files in these directories are sourced automatically. This is a good place for completion definitions:

```
~/.config/nushell/completions/
  git.nu
  ssh.nu
  docker.nu
```

### NU_LIB_DIRS

The `$NU_LIB_DIRS` constant includes the completions directory:

```nu
$nu.data-dir/completions  # Added to NU_LIB_DIRS by default
```

Files in these directories can be sourced using `source` or `use`:

```nu
use ~/.config/nushell/completions/git.nu
```

## Environment Variables

### Setting Environment Variables

| Method | Syntax | Use Case |
|--------|--------|----------|
| Direct | `$env.FOO = 'BAR'` | Simple, single variable |
| `load-env` | `load-env { BOB: FOO, JAY: BAR }` | Multiple variables at once |
| One-shot | `FOO=BAR $env.FOO` | Temporary, single command |
| `with-env` | `with-env { FOO: BAR } { ... }` | Explicit temporary |
| `def --env` | In custom commands | Environment from custom commands |

### ENV_CONVERSIONS

Converts environment variables between strings and structured values:

```nu
$env.ENV_CONVERSIONS = {
    FOO: {
        from_string: { |s| $s | split row '-' }
        to_string: { |v| $v | str join '-' }
    }
}
```

- `from_string` — converts string → value (after env.nu/config.nu loaded)
- `to_string` — converts value → string (every time an external command runs)

### Scoping

Environment variables are scoped — changes in inner blocks don't affect parent scopes:

```nu
$env.FOO = "BAR"
do {
    $env.FOO = "BAZ"  # only exists inside this block
}
# $env.FOO is still "BAR"
```

## Mode-Specific Behavior

| Mode | Command | Config Files | REPL | Autoload |
|------|---------|--------------|------|----------|
| Normal | `nu` | Yes | Yes | Yes |
| Login | `nu -l` | Yes | Yes | Yes |
| Command string | `nu -c` | No | No | No |
| Script file | `nu file.nu` | No | No | No |
| No config | `nu -n` | No | Yes | No |

## Setting Up Completions in config.nu

### Typical Pattern

```nu
# Define completers
let carapace_completer = {|spans: list<string>|
    CARAPACE_LENIENT=1 carapace $spans.0 nushell ...$spans | from json
}

let fish_completer = {|spans|
    fish --command $"complete '--do-complete=($spans | str replace --all "'" "\\'" | str join ' ')'"
    | from tsv --flexful --noheaders --no-infer
    | rename value description
}

# Handle aliases
let external_completer = {|spans|
    let expanded_alias = scope aliases
    | where name == $spans.0
    | get -o 0.expansion

    let spans = if $expanded_alias != null {
        $spans
        | skip 1
        | prepend ($expanded_alias | split row ' ' | take 1)
    } else {
        $spans
    }

    match $spans.0 {
        nu => $fish_completer
        git => $fish_completer
        _ => $carapace_completer
    } | do $in $spans
}

# Configure completions
$env.config.completions.external = {
    enable: true
    max_results: 100
    completer: $external_completer
}
```

### Using Autoload for Externs

Place extern definitions in autoload files:

```nu
# ~/.config/nushell/completions/git.nu
module "git-completions" {
    def "nu-complete git branches" [] {
        ^git branch 2>/dev/null | lines | each {|l| $l | str replace -r '^\*?\s+' '' }
    }

    export extern "git checkout" [
        branch?: string@"nu-complete git branches"
        -b
    ]
}
use "git-completions"
```

## References

- [Nushell Book: Configuration](https://www.nushell.sh/book/configuration.html)
- [Nushell Book: Environment](https://www.nushell.sh/book/environment.html)
- [Nushell Cookbook: External Completers](https://www.nushell.sh/cookbook/external_completers.html)

## Related Skills

- [references/completion.md](references/completion.md) — External completer closure details
- [references/reedline.md](references/reedline.md) — Menu and keybinding configuration
- [references/externs.md](references/externs.md) — Extern command definitions
