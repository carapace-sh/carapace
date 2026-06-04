# Carapace-pflag: Flag Parsing Extensions

Reference for [carapace-pflag](https://github.com/carapace-sh/carapace-pflag) — the carapace fork of [spf13/pflag](https://github.com/spf13/pflag) with non-POSIX flag support and structured errors. Module path remains `github.com/spf13/pflag` for drop-in compatibility.

## Flag Modes

The `Mode` field on `Flag` controls how a flag is accessed on the command line:

| Mode | Value | Behavior | Example |
|------|-------|----------|---------|
| `Default` | `0` | Standard POSIX: `--name` and `-shorthand` | `--verbose`, `-v` |
| `ShorthandOnly` | `1` | Only `-shorthand` works; `--name` is treated as unknown | `-STOP` |
| `NameAsShorthand` | `2` | `-name` works alongside `--name` (single-dash longhand) | `-help`, `--help` |

### Effects of Mode

- **`ShorthandOnly`**: `--name` is silently skipped during parsing (like an unknown flag). Help text shows only `-shorthand`.
- **`NameAsShorthand`**: In `AddFlag`, the flag name is registered in `shorthands` map, enabling `-name` syntax. Automatically makes `FlagSet.IsPosix()` return `false`, disabling POSIX shorthand chaining (`-abc` no longer means `-a -b -c`).

## N/S/NF/SF Suffix Methods

Every flag type has registration methods with suffixes controlling the mode:

| Suffix | Mode set | Returns `*Flag`? | Example |
|--------|----------|-------------------|---------|
| (none) | `Default` | No | `Bool(name, value, usage)` |
| `P` | `Default` (with shorthand) | No | `BoolP(name, shorthand, value, usage)` |
| `N` | `NameAsShorthand` | No | `BoolN(name, shorthand, value, usage)` |
| `S` | `ShorthandOnly` | No | `BoolS(name, shorthand, value, usage)` |
| `NF` | `NameAsShorthand` | Yes | `VarNF(value, name, shorthand, usage)` |
| `SF` | `ShorthandOnly` | Yes | `VarSF(value, name, shorthand, usage)` |

Core methods (`VarN`, `VarS`, `VarNF`, `VarSF`) set the `Mode` field and register shorthands. All typed wrappers (`BoolN`, `StringS`, `IntSliceN`, `DurationS`, etc.) delegate to these.

### When to use each suffix

| Use case | Suffix | Reason |
|----------|--------|--------|
| Standard POSIX flag | (none) or `P` | Normal `--name`/`-s` behavior |
| Tool uses single-dash long flags (e.g., `tar -xzf`, `go -mod`) | `N` | `-name` works like `--name` |
| Flag should only be short form (e.g., `-STOP` signal flags) | `S` | No `--` form available |
| Need `*Flag` reference (e.g., for annotations) | `NF`/`SF` | Returns pointer for further modification |

## Nargs

The `Nargs` field on `Flag` controls how many positional arguments a flag consumes:

| Nargs value | Behavior |
|-------------|----------|
| `0` or `1` (default) | Consumes exactly 1 argument from the next position |
| `> 1` | Consumes exactly N arguments; joins them as CSV for `Set()` |
| `< 0` (e.g., `-1`) | Consumes arguments greedily until one starts with `-` (assumed to be next flag); joins them as CSV |

Designed for `Slice`/`Array` flag types. Example: `--env KEY=VAL KEY2=VAL2` with `Nargs: -1`.

## OptargDelimiter

The `OptargDelimiter` field on `Flag` controls the character that separates a flag from its attached argument:

```go
flag.OptargDelimiter = ':'  // enables -agentlib:jdwp instead of -agentlib=jdwp
```

Default is `'='`. Set to a different rune for Java-style (`:`), path-style (`/`), or other conventions. Set to `-1` to disable delimiter-based argument attachment entirely.

Affects both longhand (`--flag:val`) and shorthand (`-f:val`) parsing. Each flag can have its own delimiter — they coexist in the same `FlagSet`.

## ArgumentStyle

The `ArgumentStyle` bitmask on `Flag` controls which binding forms a flag accepts:

| Bit | Constant | Form | Example |
|-----|----------|------|---------|
| `1 << 0` | `AcceptNext` | `-f arg` (next positional) | `--output file.txt` |
| `1 << 1` | `AcceptDelimited` | `-f=arg` (delimiter-attached) | `--output=file.txt` |
| `1 << 2` | `AcceptAttached` | `-farg` (POSIX attached) | `-ooutput` |

Zero value (`0`) accepts all three forms (backward compatible). Combine with OR: `AcceptDelimited | AcceptNext` allows only `--flag=val` and `--flag val`.

When a form is not accepted, parsing returns `ValueRequiredError`.

## IsPosix()

`FlagSet.IsPosix()` returns `false` if any shorthand is multi-character (which happens when `NameAsShorthand` registers flag names as shorthands, or `ShorthandOnly` uses long shorthands like `"STOP"`).

When non-POSIX:
- Shorthand chaining is disabled (`-abc` is a single shorthand, not `-a -b -c`)
- `parseSingleShortArg` treats the entire string after `-` as one shorthand name
- Different error messages for unknown shorthands

## ParseErrorsAllowlist

```go
type ParseErrorsAllowlist struct {
    UnknownFlags bool // silently skip unknown flags during parsing
}
```

When `UnknownFlags` is `true`, unknown flags are silently ignored instead of causing an error. Useful for lenient parsing in completion contexts.

## Typed Errors

Carapace-pflag replaces upstream's `fmt.Errorf` strings with structured error types:

| Error Type | Replaces | Key Methods |
|------------|----------|-------------|
| `NotExistError` | `"flag does not exist"`, `"unknown flag/shorthand"` | `GetSpecifiedName()`, `GetSpecifiedShortnames()` |
| `ValueRequiredError` | `"flag needs an argument"` | `GetFlag()`, `GetSpecifiedName()` |
| `InvalidValueError` | `"invalid argument"` | `GetFlag()`, `GetValue()`, `Unwrap()` |
| `InvalidSyntaxError` | `"bad flag syntax"` | `GetSpecifiedFlag()` |

Error messages adapt to flag `Mode` — e.g., `InvalidValueError` shows `-s, --name` for `Default` but just `-s` for `ShorthandOnly`.

`NotExistError` has a `notExistErrorMessageType` enum with 6 variants including non-POSIX-specific message formats.

## Core Flag Type Extension Summary

| Field | Type | Default | Purpose |
|-------|------|---------|---------|
| `Mode` | `mode (int)` | `Default` | Flag access mode |
| `Nargs` | `int` | `0` | Number of args consumed |
| `OptargDelimiter` | `rune` | `'='` | Delimiter for attached args |
| `ArgumentStyle` | `ArgumentStyle (uint)` | `0` (accept all) | Which binding forms accepted |
| `NoOptDefVal` | `string` | `""` | Optional-arg sentinel (from upstream) |

## Integration with carapace

The carapace library's `internal/pflagfork` package reads these unexported fields via reflection:

- `Flag.Mode()` — reads `Mode` field
- `Flag.Nargs()` — reads `Nargs` field
- `Flag.OptargDelimiter()` — reads `OptargDelimiter` field
- `FlagSet.IsPosix()` — calls unexported `IsPosix()` method
- `FlagSet.IsInterspersed()` — reads unexported `interspersed` field

The `carapace-spec` library's `internal/pflagfork` has a separate, slim read-only version for code generation (no `FlagSet` wrapper, no parsing logic).

## Related Skills

- **references/traverse.md** — how pflagfork is used during argument traversal
- **carapace-integrate** — user-facing non-POSIX flag setup (from the integration perspective)
