# Carapace Library: Config & Environment

Reference for [carapace](https://github.com/carapace-sh/carapace)'s runtime configuration — `styles.json` loading and CARAPACE_* environment variables.

## styles.json — User Style Configuration

`config.Load()` reads `$XDG_CONFIG_HOME/carapace/styles.json` (fallback: `~/.config/carapace/styles.json`):

```json
{
  "style": {
    "keyword": "blue bold",
    "path":    "underline"
  }
}
```

The file maps style category names to style strings. Each registered style config struct (`config.RegisterStyle`) is a field on the `config.Styles` map. Field values are injected via reflection:

```go
config.RegisterStyle("style", &struct {
    Keyword string `description:"keyword style"`
    Path    string `description:"path style"`
}{})
```

`config.Load()` unmarshals the JSON and sets struct fields by name match.

## Config Loading Flow

```
config.Load()
  └─ load("styles", config.Styles)
       └─ xdg.UserConfigDir() → ~/.config (or $XDG_CONFIG_HOME)
            └─ ReadFile ~/.config/carapace/styles.json
                 └─ Unmarshal into map[name]map[field]value
                      └─ reflect.ValueOf(s).Elem().FieldByName(k).SetString(v)
```

## CARAPACE_* Environment Variables

All defined in `internal/env/env.go`:

| Constant | Env Var | Purpose |
|----------|---------|---------|
| `CARAPACE_COLOR` | `CARAPACE_COLOR` | Enable color explicitly (0=disable) |
| `CARAPACE_COMPLINE` | `CARAPACE_COMPLINE` | Raw completion line (cmd-clink) |
| `CARAPACE_COVERDIR` | `CARAPACE_COVERDIR` | Coverage directory for sandbox tests |
| `CARAPACE_DESCRIPTION_LENGTH` | `CARAPACE_DESCRIPTION_LENGTH` | Max description length (default 80) |
| `CARAPACE_EXPERIMENTAL` | `CARAPACE_EXPERIMENTAL` | Enable experimental features |
| `CARAPACE_HIDDEN` | `CARAPACE_HIDDEN` | Show hidden cmds/flags (1=exclude carapace, 2=include) |
| `CARAPACE_LENIENT` | `CARAPACE_LENIENT` | Allow unknown flags |
| `CARAPACE_LOG` | `CARAPACE_LOG` | Enable debug logging |
| `CARAPACE_MATCH` | `CARAPACE_MATCH` | Match mode (e.g., `CASE_INSENSITIVE`) |
| `CARAPACE_MERGEFLAGS` | `CARAPACE_MERGEFLAGS` | Merge shorthand/longhand flags into single tag |
| `CARAPACE_NOSPACE` | `CARAPACE_NOSPACE` | Additional nospace suffixes |
| `CARAPACE_SANDBOX` | `CARAPACE_SANDBOX` | JSON mock context for sandbox tests |
| `CARAPACE_TOOLTIP` | `CARAPACE_TOOLTIP` | Enable tooltip style |
| `CARAPACE_UNFILTERED` | `CARAPACE_UNFILTERED` | Skip prefix filtering |
| `CARAPACE_ZSH_HASH_DIRS` | `CARAPACE_ZSH_HASH_DIRS` | Zsh hash directories |
| `CLICOLOR` | `CLICOLOR` | Disable color (standard CLI convention) |
| `NO_COLOR` | `NO_COLOR` | Disable color (standard NO_COLOR convention) |

## Boolean Env Parsing

```go
func getBool(s string) bool {
    switch os.Getenv(s) {
    case "true", "1":
        return true
    default:
        return false
    }
}
```

All boolean env helpers use this. "true" and "1" are truthy; everything else (including empty string) is falsy.

## Hidden Mode

```go
const (
    HIDDEN_NONE hidden = iota
    HIDDEN_EXCLUDE_CARAPACE
    HIDDEN_INCLUDE_CARAPACE
)

func Hidden() hidden {
    switch parsed, _ := strconv.Atoi(os.Getenv(CARAPACE_HIDDEN)); parsed {
    case 1:  return HIDDEN_EXCLUDE_CARAPACE
    case 2:  return HIDDEN_INCLUDE_CARAPACE
    default: return HIDDEN_NONE
    }
}
```

Integer env var: `0`/default = show normal, `1` = exclude `_carapace`, `2` = include all including hidden.

## Sandbox Detection

```go
func Sandbox() (m *mock.Mock, err error) {
    sandbox := os.Getenv(CARAPACE_SANDBOX)
    if sandbox == "" || !isGoRun() {
        return nil, errors.New("no sandbox")
    }
    err = json.Unmarshal([]byte(sandbox), &m)
    return
}

func isGoRun() bool {
    return strings.Contains(os.Args[0], "/go-build")
}
```

Only activates under `go run` (not compiled binaries). Returns the unmarshaled `*mock.Mock`.

## Logging

`log.LOG` is a `*log.Logger` that writes to a file when `CARAPACE_LOG` is set:

```go
func init() {
    if !env.Log() {
        LOG = log.New(io.Discard, "", log.Flags())
        return
    }
    tmpdir := fmt.Sprintf("%v/carapace", os.TempDir())
    os.MkdirAll(tmpdir, os.ModePerm)
    file := fmt.Sprintf("%v/%v.log", tmpdir, uid.Executable())
    logfileWriter, _ := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    LOG = log.New(logfileWriter, ps.DetermineShell()+" ", log.Flags()|log.Lmsgprefix|log.Lmicroseconds)
}
```

Log file: `/tmp/carapace/<executable>.log`. Prefix includes the detected shell name and microsecond timestamps.

## LOG Usage

```go
LOG.Printf("executing PreRun for %#v with args %#v", cmd.Name(), args)
```

Conditional — no-op when `CARAPACE_LOG` is not set.

## Color Disabled Logic

```go
func ColorDisabled() bool {
    if v, ok := os.LookupEnv(CARAPACE_COLOR); ok {
        return v == "0"
    }
    return getBool(NO_COLOR) || os.Getenv(CLICOLOR) == "0"
}
```

Priority: explicit `CARAPACE_COLOR=0` > `NO_COLOR` > `CLICOLOR=0`. Explicit `CARAPACE_COLOR=1` enables even when `NO_COLOR` is set.

## Description Length

```go
func DescriptionLength() int {
    if parsed, err := strconv.Atoi(os.Getenv(CARAPACE_DESCRIPTION_LENGTH)); err == nil && parsed > 0 {
        return parsed
    }
    return 80
}
```

Integer env var with default 80.

## Related Skills

- **references/style.md** — style system that uses config values
- **references/sandbox.md** — CARAPACE_SANDBOX injection
- **references/cache.md** — CARAPACE_COVERDIR for coverage