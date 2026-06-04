# Carapace Library: Styling System Internals

Reference for [carapace](https://github.com/carapace-sh/carapace)'s styling system — how styles are represented, resolved, configured, and rendered across shells.

## String-Based Composable Styles

Styles are space-separated token strings, not structs. This makes them serializable, composable, and easy to use in YAML specs:

```go
style.Of(style.Bold, style.Red)           // "bold red"
style.Of(style.BgBlue, style.BrightYellow) // "bg-blue bright-yellow"
```

### Predefined Style Constants

| Category | Constants |
|----------|-----------|
| Foreground | `Black`, `Red`, `Green`, `Yellow`, `Blue`, `Magenta`, `Cyan`, `White` |
| Bright foreground | `BrightBlack`, `BrightRed`, `BrightGreen`, `BrightYellow`, `BrightBlue`, `BrightMagenta`, `BrightCyan`, `BrightWhite` |
| Background | `BgBlack`, `BgRed`, `BgGreen`, `BgYellow`, `BgBlue`, `BgMagenta`, `BgCyan`, `BgWhite` |
| Bright background | `BgBrightBlack`..`BgBrightWhite` |
| Attributes | `Bold`, `Dim`, `Italic`, `Underlined`, `Blink`, `Inverse` |
| Default | `Default` = `""` (no style) |

### Color Functions

| Function | Signature | Returns |
|----------|-----------|---------|
| `XTerm256Color` | `(i uint8) string` | 256-color palette: `"color<N>"` |
| `TrueColor` | `(r, g, b uint8) string` | 24-bit true color: `"#rrggbb"` |

### Core Functions

| Function | Purpose |
|----------|---------|
| `Of(s ...string)` | Combine style tokens (space-joined, trimmed) |
| `SGR(s string)` | Convert style string to ANSI SGR escape sequence |
| `Parse(s string) ui.Style` | Parse style string into elvish `ui.Style` |
| `Register(name, i)` | Register a style config struct with the config system |
| `Set(key, value) error` | Set a style value by dotted key (e.g., `"carapace.Value"`) |

The `Parse` function uses the elvish `ui` library's `ParseStyling` for each space-separated token, then composes them via `ui.ApplyStyling`.

## Semantic Style Resolution

### ForKeyword

```go
style.ForKeyword(value, context) string
```

Maps semantic keywords to colors with normalization (strip `-`/`_`, lowercase):

| Keyword pattern | Style |
|----------------|-------|
| `yes`, `true`, `on`, `enable`, `allow`, `installed`, `active`, `success`, `ok`, `healthy`, `running` | `style.Carapace.KeywordPositive` (default: `green`) |
| `no`, `false`, `off`, `disable`, `deny`, `missing`, `inactive`, `failed`, `error`, `unhealthy`, `stopped`, `exited` | `style.Carapace.KeywordNegative` (default: `red`) |
| `auto`, `default`, `if`, `interactive`, `ask` | `style.Carapace.KeywordAmbiguous` (default: `yellow`) |
| Everything else | `style.Carapace.KeywordUnknown` (default: `Default`) |

### ForLogLevel

```go
style.ForLogLevel(value, context) string
```

Maps log level strings to colors:

| Level | Style |
|-------|-------|
| trace | `style.Carapace.LogLevelTrace` |
| debug | `style.Carapace.LogLevelDebug` |
| info | `style.Carapace.LogLevelInfo` |
| warn/warning | `style.Carapace.LogLevelWarn` |
| error/err | `style.Carapace.LogLevelError` |
| fatal/critical | `style.Carapace.LogLevelFatal` |

### ForPath / ForPathExt / ForExtension

```go
style.ForPath(path, context) string    // style based on full path (file must exist)
style.ForPathExt(path, context) string  // style based on file extension only (file need not exist)
style.ForExtension(ext, context) string // style for extension string (e.g., "json")
```

These use the elvish `lscolors` library to honor the user's `LS_COLORS` environment variable. Resolution:

- `ForPath`: checks if path exists, applies `LS_COLORS` based on type (dir, symlink, executable) and extension
- `ForPathExt`: applies `LS_COLORS` based on extension pattern only
- `ForExtension`: wraps extension in `*.<ext>` pattern for `LS_COLORS` matching

The `Context` interface provides `Abs()`, `Getenv()`, `LookupEnv()` — `carapace.Context` implements this.

## The Carapace Config Struct

`style.Carapace` is the central style configuration struct. All fields are `string` type with `description` and `tag` struct tags for runtime introspection:

| Field | Default | Description |
|-------|---------|-------------|
| `Value` | `Default` | Style for completion values |
| `Description` | `Dim` | Style for descriptions |
| `Error` | `Red` + `Underlined` | Style for error messages |
| `Usage` | `Italic` | Style for usage hints |
| `KeywordPositive` | `Green` | Style for positive keywords |
| `KeywordNegative` | `Red` | Style for negative keywords |
| `KeywordAmbiguous` | `Yellow` | Style for ambiguous keywords |
| `KeywordUnknown` | `Default` | Style for unknown keywords |
| `LogLevelTrace` | `Dim` | Log level: trace |
| `LogLevelDebug` | `Blue` | Log level: debug |
| `LogLevelInfo` | `Default` | Log level: info |
| `LogLevelWarn` | `Yellow` | Log level: warn |
| `LogLevelError` | `Red` | Log level: error |
| `LogLevelFatal` | `Red` + `Bold` | Log level: fatal |
| `Highlight1`..`Highlight12` | Cyclic ANSI colors | Cyclic highlight styles for groups |
| `FlagArg` | Style for flags that take an argument | |
| `FlagMultiArg` | Style for flags that take multiple args | |
| `FlagNoArg` | Style for flags that take no argument | |
| `FlagOptArg` | Style for flags with optional argument | |

`Highlight(n int)` returns the highlight style for level 0..11, cycling through `Highlight1`–`Highlight12`.

## Runtime Style Configuration

### Register

```go
style.Register("carapace", &style.Carapace)
```

Registers a style config struct with the config system. Struct fields must be `string` type with `description`/`tag` struct tags. This is called automatically during initialization.

### Set

```go
style.Set("carapace.Value", "bold magenta")  // override value style
style.Set("carapace.Value", "")              // delete override (revert to default)
```

Persists style overrides to `$XDG_CONFIG_HOME/carapace/styles.json`. Key format is `"GroupName.FieldName"`.

### Load

`config.Load()` is called during completion in `complete()` — loads `styles.json` and applies all overrides to registered structs.

### Introspection

| Function | Returns |
|----------|---------|
| `GetStyleConfigs()` | Names of all registered style config groups |
| `GetStyleFields(name)` | `[]Field` descriptors for a config group |

`Field` struct: `{ Name, Description, Style, Tag string }`.

## SGR & Shell Rendering

`style.SGR(s)` converts a style string to an ANSI SGR escape sequence. This is used by shells that support inline SGR codes (bash, oil, tcsh). Other shells use different rendering:

| Shell | Style rendering |
|-------|---------------|
| zsh | `zstyle` color rules with `(#b)` pattern matching |
| fish | Native `--description` styling |
| elvish | `ui.ParseStyling` + `styled` builtin |
| nushell | Custom `fg`/`bg`/`attr` format |
| xonsh | `bg:ansi<N> fg:ansi<N> bold italic` format |
| powershell | SGR escape codes in `ListItemText` |

## Related Skills

- **references/action.md** — how `.Style()`/`.StyleF()`/`.StyleR()` apply styles to Actions
- **references/shell.md** — per-shell style rendering details
- **carapace-action** — using styles in carapace-bin shared actions
