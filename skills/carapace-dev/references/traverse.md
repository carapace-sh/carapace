# Carapace Library: Completion Engine & Traversal

Reference for [carapace](https://github.com/carapace-sh/carapace)'s core completion engine — how CLI arguments are classified, traversed, and resolved to completion actions.

## Dispatch Flow: `complete()`

```
complete(cmd, args) entry point
  │
  ├─ len(args) == 0 → Snippet for auto-detected shell
  ├─ len(args) == 1 → Snippet for named shell
  └─ len(args) >= 2 → Perform actual completion
       │
       ├─ Shell-specific patching
       │   ├─ bash: bash.Patch(args) — re-lex COMP_LINE, handle redirects
       │   ├─ nushell: nushell.Patch(args) — strip quote delimiters
       │   └─ cmd-clink: cmd_clink.Patch(args) — re-lex CARAPACE_COMPLINE
       │
       ├─ traverse(cmd, args[2:]) → (Action, Context)
       │
       ├─ config.Load() — load styles.json
       │
       └─ action.Invoke(context).value(shell, currentValue)
            └─ shell.Value() → shell-specific output formatting
```

## Shell Patching

Before traversal, shell-specific patching corrects args that shells pass in ways incompatible with carapace's parser:

| Shell | Patch function | What it does |
|-------|---------------|-------------|
| bash | `bash.Patch(args)` | Re-lexes `COMP_LINE`/`COMP_POINT` via `carapace-shlex`; filters out redirects (`>`, `2>`); on `RedirectError`, returns `ActionFiles()`; extracts `wordbreakPrefix` and `compType` as side effects |
| nushell | `nushell.Patch(args)` | Strips quote delimiters (`"`, `'`, `` ` ``) from args since nushell passes raw quoted strings |
| cmd-clink | `cmd_clink.Patch(args)` | Re-lexes `CARAPACE_COMPLINE` via `carapace-shlex`; filters redirects like bash |

After patching, all `COMP_*` env vars are unset for bash.

## The `traverse()` State Machine

`traverse(cmd *cobra.Command, args []string) (Action, Context)` walks the cobra command tree classifying each argument:

```
State variables:
  inArgs[]         — args consumed by current command
  inPositionals[]  — positionals consumed by current command
  inFlag           — last flag that still expects arguments (*pflagfork.Flag)
  dash             — whether "--" was encountered
```

### Argument Classification (priority order)

Each arg in the main loop is classified into one of 5 categories:

| Priority | Category | Condition |
|----------|----------|-----------|
| 1 | **Flag argument** | `inFlag != nil && inFlag.Consumes(arg)` |
| 2 | **Dash** | `arg == "--"` |
| 3 | **Flag** | Starts with `-`, flag parsing enabled, interspersed or no positionals yet |
| 4 | **Subcommand** | Matches a subcommand name → **recurse** |
| 5 | **Positional** | Default |

### Subcommand Recursion

When a subcommand is found:
1. Parse accumulated `inArgs` as flags via `cmd.ParseFlags(inArgs)`
2. Update `context.Args = cmd.Flags().Args()`
3. Recurse: `traverse(subcommand, args[i+1:])`

### Post-Loop Edge Cases

After the main loop, before deciding what to complete:

| Case | Handling |
|------|----------|
| `dash` | Skip flag parsing for remaining args |
| Flag missing its argument | Remove the flag from `toParse` so cobra doesn't choke |
| Shorthand series (e.g., `-abcX`) | Look up partial flag; either remove suffix or add the value to `toParse` |

### Completion Decision

After parsing flags, `traverse()` decides what to complete:

| Case | Returns |
|------|---------|
| Dash argument (`--` encountered) | `storage.getPositional(cmd, len(dashPositionals))` |
| Flag argument (`inFlag.Consumes(context.Value)`) | `storage.getFlag(cmd, inFlag.Name)` with `context.Parts = inFlag.Args` |
| Optional flag arg (`NoOptDefVal` or attached `=`) | `storage.getFlag(cmd, f.Name).Prefix(f.Prefix)` or `ActionValues("true","false")` for bool |
| Attached flag arg (POSIX shorthand) | `storage.getFlag(cmd, f.Name).Prefix(f.Prefix)` |
| Flag name (starts with `-`) | `actionFlags(cmd)` — all available flags |
| Positional + subcommands | `Batch(getPositional) + ActionCommands(cmd)` if subcommands exist and no positionals consumed |

## pflagfork: Flag & FlagSet Wrappers

`internal/pflagfork` extends pflag with carapace-specific metadata needed for traversal. All unexported fields are read via reflection.

### Flag Wrapper

```go
type Flag struct {
    *pflag.Flag
    Prefix string   // e.g., "--verbose=" when completing --verbose=foo
    Args   []string // already-consumed flag arguments
}
```

| Method | What it reads | Purpose |
|--------|-------------|---------|
| `Nargs()` | Unexported `Nargs` field | Multi-arg flag consumption |
| `Mode()` | Unexported `Mode` field | Default/ShorthandOnly/NameAsShorthand |
| `OptargDelimiter()` | Unexported `OptargDelimiter` field | Delimiter for `--flag=val` vs `--flag:val` |
| `IsRepeatable()` | Value type string | True for Slice/Array/count |
| `TakesValue()` | Value type | False for bool/boolSlice/count |
| `IsOptarg()` | `NoOptDefVal != ""` | Flag arg is optional |
| `Consumes(arg)` | Nargs + Args + TakesValue | Does this flag still need more args? |
| `Style()` | Type + Optarg | Style constant based on flag kind |
| `Required()` | Cobra annotation | `BashCompOneRequiredFlag` |
| `Definition()` | All metadata | Human-readable: `-v, --verbose!*?` |

### Consumes() Logic

`Consumes(arg)` determines if a flag still needs more arguments:

```
Consumes(arg) == true when:
  - Flag takes value (not bool/count)
  - Flag is not optarg (optarg never consumes the next positional)
  - No args consumed yet, OR
  - Nargs > 1 and fewer args consumed than Nargs, OR
  - Nargs < 0 and arg doesn't look like a flag (doesn't start with "-")
```

### FlagSet Wrapper

```go
type FlagSet struct {
    *pflag.FlagSet
}
```

| Method | Purpose |
|--------|---------|
| `IsInterspersed()` | Reads unexported `interspersed` field |
| `IsPosix()` | Calls unexported `IsPosix()` method — controls shorthand chaining |
| `IsShorthandSeries(arg)` | Regex `^-(?P<shorthand>[^-].*)` + IsPosix |
| `IsMutuallyExclusive(flag)` | Checks `cobra_annotation_mutually_exclusive` annotations |
| `LookupArg(arg)` | **Critical dispatcher** — routes to correct lookup method |
| `ShorthandLookup(name)` | Wraps pflag's ShorthandLookup with `Args: []string{}` |
| `VisitAll(fn)` | Iterates flags, initializing `Args: []string{}` |

### LookupArg Dispatch

`LookupArg(arg)` routes based on prefix and POSIX mode:

```
arg starts with "--" → lookupPosixLonghandArg
  - Splits on flag's OptargDelimiter
  - Sets Prefix and Args on the returned Flag

arg starts with "-" and IsPosix → lookupPosixShorthandArg
  - Iterates each character of the shorthand series
  - Handles attached values, optarg delimiters
  - Sets Prefix and Args

arg starts with "-" and !IsPosix → two-phase lookup:
  1. LookupNonPosixLonghandArg (single-dash longhand: -name)
  2. Fallback: lookupNonPosixShorthandArg (single-dash shorthand: -n)
```

## Lenient Mode

When `CARAPACE_LENIENT` is set, `traverse()` enables `cmd.FParseErrWhitelist.UnknownFlags = true` before traversal, allowing unknown flags to be silently ignored.

## Post-Processing in `shell.Value()`

After traversal and invocation, `shell.Value()` applies post-processing before formatting:

1. **Color disable**: if `env.ColorDisabled()`, strip styles and set fallback styles
2. **Prefix filtering**: `values.FilterPrefix(value)` unless `CARAPACE_UNFILTERED`
3. **Flag merging**: merge "shorthand flags"/"longhand flags" tags into "flags" — implicit for zsh, explicit via `CARAPACE_MERGEFLAGS`
4. **Message integration**: For shells without native message support (not elvish/export/zsh), messages are integrated as synthetic completion values
5. **Nospace propagation**: If messages exist or `CARAPACE_NOSPACE` is set, add `*` to nospace set
6. **Sort + dedup**: `sort.Sort(ByDisplay(values))` → clear UIDs → `values.Unique()`

## Related Skills

- **references/pflag.md** — the flag extensions (Mode, Nargs, OptargDelimiter) that pflagfork reads
- **references/action.md** — the Action API that traverse() returns
- **references/shell.md** — the per-shell formatting that `shell.Value()` dispatches to
