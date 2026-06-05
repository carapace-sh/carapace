# Carapace Library: Action API Internals

Reference for the [carapace](https://github.com/carapace-sh/carapace) library's core completion API — the `Action` type, its modifiers, and the invocation model.

## The Action Type

`Action` is the fundamental completion type. It is **lazy** — it either holds pre-computed `rawValues` or a `callback` that produces them at invocation time.

```go
type Action struct {
    meta      common.Meta      // metadata (messages, nospace, usage, queries)
    rawValues common.RawValues  // pre-computed values (non-callback actions)
    callback  CompletionCallback // lazy function producing an Action
}

type CompletionCallback func(c Context) Action
type ActionMap map[string]Action
```

## Two-Phase Model: Action → InvokedAction

Actions are descriptions of completions, not the completions themselves. The two-phase model separates definition from execution:

1. **`Action`** (deferred) — constructed at init time, carries a callback chain
2. **`Invoke(c Context)`** → **`InvokedAction`** (materialized) — callback executed, values resolved

```go
// Phase 1: define (lazy, no work done yet)
action := carapace.ActionValues("a", "b", "c").Style(style.Blue).Tag("letters")

// Phase 2: invoke (work happens now)
invoked := action.Invoke(carapace.Context{Value: "a"})
```

`InvokedAction` has its own modifier set (`Filter`, `Merge`, `Prefix`, `Retain`, `Suffix`, `UidF`, `ToA`, `ToMultiPartsA`, `MarshalJSON`) that operates directly on already-materialized values.

### When to use each phase

| Phase | Use when |
|-------|----------|
| `Action` modifiers | At definition time — building the callback chain |
| `InvokedAction` modifiers | After invocation — direct value transforms, merging results from parallel batches |

## Context

`Context` provides all runtime state during completion:

```go
type Context struct {
    Value string     // current value being completed (or part of it during MultiParts)
    Args  []string   // positional arguments of current subcommand (excluding the one being completed)
    Parts []string   // split Value during ActionMultiParts (excluding current part)
    Env   []string   // environment variables
    Dir   string     // working directory
}
```

| Method | Purpose |
|--------|---------|
| `NewContext(args ...string)` | Last arg = `Value`, rest = `Args`; inherits `os.Environ()` and `os.Getwd()` |
| `LookupEnv(key)` | Lookup env var in context |
| `Getenv(key)` | Get env var (empty string if missing) |
| `Setenv(key, value)` | Add env var to context |
| `Envsubst(s)` | Replace `${var}` in string using context env |
| `Command(name, arg...)` | Create `exec.Cmd` with context's `Env` and `Dir` |
| `Abs(path)` | Resolve to absolute path (handles `~`, relative, volume prefixes) |

## Built-in Action Constructors

| Function | Returns | Purpose |
|----------|---------|---------|
| `ActionCallback(f)` | `Action` | Create from `CompletionCallback` |
| `ActionExecCommand(name, arg...)` | `func(func([]byte) Action) Action` | Execute external command; curried — call with command first, then callback |
| `ActionExecCommandE(name, arg...)` | `func(func([]byte, error) Action) Action` | Like above with custom error handling |
| `ActionImport(output)` | `Action` | Parse JSON export as Action |
| `ActionExecute(cmd)` | `Action` | Delegate to a cobra subcommand's completions |
| `ActionDirectories()` | `Action` | Complete directories |
| `ActionFiles(suffix...)` | `Action` | Complete files (optional suffix filter) |
| `ActionValues(vals...)` | `Action` | Static completion values |
| `ActionStyledValues(vals...)` | `Action` | Values with style (pairs: value, style) |
| `ActionValuesDescribed(vals...)` | `Action` | Values with description (pairs: value, description) |
| `ActionStyledValuesDescribed(vals...)` | `Action` | Values with description + style (triples: value, desc, style) |
| `ActionMessage(msg, args...)` | `Action` | Display help message (no completions) |
| `ActionMultiParts(divider, f)` | `Action` | Complete parts of arg separated by divider |
| `ActionMultiPartsN(divider, n, f)` | `Action` | Like above but limit number of parts |
| `ActionCommands(cmd)` | `Action` | Complete (sub)commands of cobra command |
| `ActionCobra(f)` | `Action` | Bridge cobra completion function |
| `ActionExecutables(dirs...)` | `Action` | Complete executables from PATH or dirs |
| `ActionPositional(cmd)` | `Action` | Complete positional args for cobra command |
| `ActionStyleConfig()` | `Action` | Complete style configuration strings |
| `ActionStyles()` | `Action` | Complete style names |

## Action Modifier Methods

All return `Action` (chainable). Most wrap in `ActionCallback` for lazy evaluation.

### Filtering & Value Selection

| Method | Signature | Purpose |
|--------|-----------|---------|
| `Filter` | `(values ...string)` | Remove given values |
| `FilterArgs` | `()` | Filter values that appear in `Context.Args` |
| `FilterParts` | `()` | Filter values that appear in `Context.Parts` |
| `Retain` | `(values ...string)` | Keep only given values |
| `Unique` | `()` | Deduplicate values by `(Value, Uid)` key |

### Prefix & Suffix

| Method | Signature | Purpose |
|--------|-----------|---------|
| `Prefix` | `(prefix string)` | Add prefix to inserted values (not display) |
| `Suffix` | `(suffix string)` | Add suffix to inserted values (not display) |
| `NoSpace` | `(suffixes ...rune)` | Disable trailing space for chars (`'*'` = all) |

### Multi-Part & Splitting

| Method | Signature | Purpose |
|--------|-----------|---------|
| `MultiParts` | `(dividers ...string)` | Split values by dividers, complete segments |
| `MultiPartsP` | `(delimiter, pattern, f)` | MultiParts with placeholder support |
| `List` | `(divider string)` | Wrap in `ActionMultiParts` with given divider |
| `UniqueList` | `(divider string)` | List with dedup via `ActionMultiParts` |
| `UniqueListF` | `(divider string, f)` | UniqueList with transform before filtering |
| `Split` | `()` | Lex split `Context.Value`, replace `Args` with tokens |
| `SplitP` | `()` | Like Split but supports pipelines (`\|`, `>`, `>>`) |

### Styling & Tagging

| Method | Signature | Purpose |
|--------|-----------|---------|
| `Style` | `(s string)` | Set static style |
| `StyleF` | `(f func(string, style.Context) string)` | Set style per value |
| `StyleR` | `(s *string)` | Set style by reference (runtime-configurable) |
| `Tag` | `(tag string)` | Set tag on all values |
| `TagF` | `(f func(string) string)` | Set tag per value |

### Directory & Context

| Method | Signature | Purpose |
|--------|-----------|---------|
| `Chdir` | `(dir string)` | Change working dir for invocation |
| `ChdirF` | `(f func(traverse.Context) (string, error))` | Chdir with function (deferred resolution) |
| `Shift` | `(n int)` | Shift `Context.Args` left by n |

### Caching & Timeout

| Method | Signature | Purpose |
|--------|-----------|---------|
| `Cache` | `(timeout time.Duration, keys ...key.Key)` | Cache callback results to disk |
| `Timeout` | `(d time.Duration, alternative Action)` | Set max invocation duration |

### Conditional & Error Handling

| Method | Signature | Purpose |
|--------|-----------|---------|
| `Unless` | `(condition bool)` | Skip action if condition true |
| `UnlessF` | `(condition func(Context) bool)` | Skip action if function returns true |
| `Suppress` | `(expr ...string)` | Suppress error messages matching regex |
| `Usage` | `(usage string, args ...any)` | Set usage string |
| `UsageF` | `(f func() string)` | Set usage via function |

### Experimental (Uid/Query)

| Method | Signature | Purpose |
|--------|-----------|---------|
| `Uid` | `(scheme, host string, opts ...string)` | Set UID |
| `UidF` | `(f func(string, uid.Context) (*url.URL, error))` | Set UID via function |
| `Query` | `(scheme, host, path string, opts ...string)` | Add query metadata |
| `QueryF` | `(f func(string, uid.Context) (*url.URL, error))` | Add query via function |

## Suffix Naming Convention

Carapace uses consistent suffixes to distinguish method variants:

| Suffix | Meaning | Examples |
|--------|---------|---------|
| `F` | Functional variant — takes a function instead of a static value | `Chdir`/`ChdirF`, `Style`/`StyleF`, `Unless`/`UnlessF`, `Uid`/`UidF`, `Usage`/`UsageF` |
| `P` | Pipeline-aware — supports shell pipes/redirects | `Split`/`SplitP`, `MultiParts`/`MultiPartsP` |
| `E` | Error-explicit — provides raw error to callback | `ActionExecCommand`/`ActionExecCommandE` |
| `N` | Limited/numbered variant | `ActionMultiParts`/`ActionMultiPartsN` |

## InvokedAction Methods

Operate on already-materialized values:

| Method | Purpose |
|--------|---------|
| `Filter(values...)` | Remove given values |
| `Merge(others...)` | Merge multiple invoked actions (dedup) |
| `Prefix(prefix)` | Add prefix to inserted values |
| `Retain(values...)` | Keep only given values |
| `Suffix(suffix)` | Add suffix to inserted values |
| `UidF(f)` | Set UID per value (experimental) |
| `ToA()` | Cast back to `Action` |
| `ToMultiPartsA(dividers...)` | Convert values to multipart action |
| `MarshalJSON()` | JSON serialization |

## Batch — Parallel Invocation

```go
batch := carapace.Batch(action1, action2, action3)

// Invoke all in parallel, merge results
result := batch.ToA() // shortcut for batch.Invoke(c).Merge().ToA()

// Manual control
invoked := batch.Invoke(c)
merged := invoked.Merge()
action := merged.ToA()
```

## Storage & Registration

`_storage` maps `*cobra.Command` → `entry` (internal, not user-facing). Registration API on `Carapace`:

| Method | Purpose |
|--------|---------|
| `Gen(cmd)` | Initialize carapace for a command |
| `Standalone()` | Prevent cobra defaults from interfering |
| `PreRun(f)` | Execute before argument parsing (modify command structure) |
| `PreInvoke(f)` | Modify actions before execution (e.g., Chdir) |
| `FlagCompletion(ActionMap)` | Map flags to completion actions |
| `PositionalCompletion(...Action)` | Map positional positions to actions |
| `PositionalAnyCompletion(Action)` | Action for all remaining positionals |
| `DashCompletion(...Action)` | Actions for positionals after `--` |
| `DashAnyCompletion(Action)` | Action for all remaining dash positionals |

## RawValue & Meta (internal/common)

`RawValue` is a single completion candidate:

```go
type RawValue struct {
    Value       string // completion text
    Display     string // shown text
    Description string // help text
    Style       string // ANSI style string
    Tag         string // group tag
    Uid         string // unique ID
}
```

`Meta` is attached to every `InvokedAction`:

```go
type Meta struct {
    Messages Messages       // error/warning messages
    Nospace  SuffixMatcher  // suffixes where no space is appended
    Usage    string         // contextual usage hint
    Queries  Queries        // deferred shell queries
}
```

## Related Skills

- **carapace-bin skill** (carapace-bin repo) — creating/modifying shared actions (Opts, Uid/QueryF, caching, macro exposure)
- **carapace-bin skill** (carapace-bin repo) — integrating carapace library into cobra CLIs (PreRun, PreInvoke, bridge, spec)
- **references/traverse.md** — how the completion engine resolves what to complete
