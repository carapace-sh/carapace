# Carapace Library: Export Format

Reference for [carapace](https://github.com/carapace-sh/carapace)'s JSON wire format in `internal/export/`.

## Export Struct

```go
type Export struct {
    Version string       `json:"version"`
    common.Meta
    Values common.RawValues `json:"values"`
}
```

The universal JSON representation of a completed action. Contains:
- **Version**: carapace library version from `debug.BuildInfo`
- **Meta**: messages, nospace set, usage hint, queries
- **Values**: sorted list of `RawValue` entries

## MarshalJSON

```go
func (e Export) MarshalJSON() ([]byte, error) {
    sort.Sort(common.ByValue(e.Values))
    return json.Marshal(&struct {
        Version string `json:"version"`
        common.Meta
        Values common.RawValues `json:"values"`
    }{
        Version: version(),
        Meta:    e.Meta,
        Values:  e.Values,
    })
}
```

Always re-sorts values by `Value` before serializing. Values must be sorted for deterministic output.

## Version Lookup

```go
func version() string {
    if info, ok := debug.ReadBuildInfo(); ok {
        for _, dep := range info.Deps {
            if dep.Path == "github.com/carapace-sh/carapace" {
                return dep.Version
            }
        }
    }
    return "unknown"
}
```

Reads `debug.BuildInfo` to find the carapace library version. Returns `"unknown"` if not available (e.g., in tests without module info).

## Usage in Cache

Cache files store `export.Export` JSON. The `common.Meta` includes messages and nospace state, so cached completions restore the full completion context:

```go
// cache.LoadE returns *export.Export
if cached, err := cache.LoadE(cacheFile, timeout); err == nil {
    return Action{meta: cached.Meta, rawValues: cached.Values}
}
```

## Usage in Sandbox

Sandbox's `run.invoke()` converts an `InvokedAction` to JSON using the export format:

```go
func (r run) invoke(a carapace.Action) string {
    meta, rawValues := common.FromInvokedAction(a.Invoke(r.context))
    rawValues = rawValues.FilterPrefix(r.context.Value)
    sort.Sort(common.ByValue(rawValues))

    m, _ := json.MarshalIndent(export.Export{
        Meta:   meta,
        Values: rawValues,
    }, "", "  ")
    return string(m)
}
```

The sandbox JSON is compared with `assert.Equal` for test assertions.

## Usage in ActionImport

`ActionImport` parses JSON back into an action:

```go
func ActionImport(output string) Action {
    var e export.Export
    if err := json.Unmarshal([]byte(output), &e); err != nil {
        return ActionMessage(err.Error())
    }
    return Action{meta: e.Meta, rawValues: e.Values}
}
```

Used by the `_carapace export` subcommand to re-import completions.

## The export Shell

The `export` shell target in `complete()` re-invokes the `_carapace` command with `export` as the shell argument:

```go
if shell == "export" {
    // Re-invokes as: myapp _carapace export <root-cmd> <args...>
    // Outputs JSON to stdout
}
```

This enables subcommand completion to work across process boundaries — the parent process re-invokes itself with the `export` shell to resolve nested subcommand completions, then parses the JSON.

## RawValue Structure

`common.RawValue` is the unit of completion data:

```go
type RawValue struct {
    Value       string
    Display     string
    Kind        rawValueKind
    Style       string
    Tag         string
    Uid         string
    Description string
}
```

- **Value**: inserted text (what the shell actually receives)
- **Display**: shown text (may differ from Value for placeholders/prefixes)
- **Kind**: type indicator (plain, path, keyword, etc.)
- **Style**: ANSI style string
- **Tag**: grouping tag for filtering
- **Uid**: unique identifier
- **Description**: shown description text

## Meta Structure

```go
type Meta struct {
    Messages Messages
    Nospace  SuffixMatcher
    Usage    string
    Queries  []string
}
```

- **Messages**: info/warning messages to display
- **Nospace**: characters after which the shell should NOT add a trailing space
- **Usage**: usage hint string
- **Queries**: query metadata for UID resolution

## Related Skills

- **references/cache.md** — export.Export is the cache file format
- **references/sandbox.md** — sandbox uses export for test assertions
- **references/action.md** — Action → InvokedAction → export cycle