# AGENTS.md

## Project Overview

Carapace is a Go library that provides command argument completion for [spf13/cobra](https://github.com/spf13/cobra)-based CLI applications. It generates shell completion scripts and handles runtime completion callbacks for 12 shells: bash, bash-ble, cmd-clink, elvish, export, fish, ion, nushell, oil, powershell, tcsh, xonsh, and zsh.

This is the **core library** (`github.com/carapace-sh/carapace`). Companion projects:
- `carapace-bin` — pre-built completions for 500+ commands
- `carapace-bridge` — bridges to other completion frameworks
- `carapace-pflag` — forked spf13/pflag with non-POSIX modes

## Build & Test Commands

```sh
# Build
go build ./...

# Run all tests (unit + example-nonposix via go.work)
go test ./...

# Run with coverage (matches CI)
mkdir -p .cover && CARAPACE_COVERDIR="$(pwd)/.cover" go test -v -coverpkg ./... -coverprofile=unit.cov ./... ./example-nonposix/...

# Benchmarks
go test -bench ./...

# Generate
go generate ./...

# Format check (CI uses this; must produce no output)
gofmt -d -s .

# Lint
golangci-lint run
```

## Architecture & Control Flow

### Two-Phase Completion Model

1. **Callback phase** — `Action` objects hold a `CompletionCallback` function. When `Invoke(ctx)` is called, the callback recursively resolves until it produces an `Action` with `rawValues` (no callback). Actions are lazy — callbacks execute only at completion time.

2. **Formatting phase** — `InvokedAction` is passed to the shell-specific formatter (`internal/shell/<shell>/action.go`) which produces the output string the shell consumes.

### Completion Dispatch Flow

```
cobra command invoked with _carapace subcommand
  → complete()                    # shell-specific arg patching
  → traverse()                   # classify args (flag, positional, subcommand, dash)
  → storage.get(cmd)            # resolve Action from per-command entry
  → action.Invoke(context)       # resolve callback chain
  → invokedAction.value(shell)   # format for target shell
```

### Core Types

| Type | File | Purpose |
|------|------|---------|
| `Carapace` | `carapace.go` | Wrapper around `*cobra.Command`; entry via `Gen()` |
| `Action` | `action.go` | Lazy completion definition (callback + meta + rawValues) |
| `InvokedAction` | `invokedAction.go` | Resolved action; supports Filter/Merge/Prefix/Suffix |
| `Context` | `context.go` | Completion context: Value, Args, Parts, Env, Dir |
| `Batch` | `batch.go` | Parallel action invocation via goroutines |
| `entry` | `storage.go` | Per-command storage: flag actions, positional actions, hooks |
| `RawValue` | `internal/common/value.go` | Single candidate: Value, Display, Description, Style, Tag, Uid |
| `Meta` | `internal/common/meta.go` | Metadata: Messages, Nospace, Usage, Queries |
| `SuffixMatcher` | `internal/common/suffix.go` | Controls when shell should NOT add trailing space |
| `pflagfork.FlagSet` | `internal/pflagfork/` | pflag wrapper with POSIX/non-POSIX flag lookup |
| `pflagfork.Flag` | `internal/pflagfork/` | Flag wrapper with Mode, Nargs, OptargDelimiter, Consumes |

## Directory Structure

```
/                           Root package — public API
  carapace.go               Gen(), PreRun(), PreInvoke(), Snippet(), IsCallback(), Test()
  action.go                 Action type, modifiers (Cache, Chdir, Filter, MultiParts, etc.)
  invokedAction.go          InvokedAction type, Merge, Prefix, Suffix, Filter, ToMultiPartsA
  complete.go               complete() — main dispatch with shell patching
  traverse.go               traverse() — argument classification state machine
  storage.go                Global per-command completion storage (mutex-protected)
  command.go                _carapace hidden subcommand + spec/style subcommands
  compat.go                 Cobra bridge: registerValidArgsFunction, registerFlagCompletion
  defaultActions.go         ActionValues, ActionFiles, ActionDirectories, ActionCommands, etc.
  internalActions.go        ActionCallback, ActionExecCommand, ActionImport, ActionMessage
  diff.go                   Diff() action
  log.go                    LOG (conditional logging)
  experimental.go           x.ClearStorage, x.Complete

/internal/
  common/                   RawValue, Meta, SuffixMatcher, Messages, Queries
  shell/                    Per-shell formatters (bash, zsh, fish, nushell, elvish, etc.)
    <shell>/action.go       ActionRawValues() — format completion output
    <shell>/snippet.go      Snippet() — generate shell setup script
    bash/patch.go           Redirect handling
    nushell/patch.go        Open-quote handling
    cmd_clink/patch.go      cmd-clink patching
    zsh/zstyle.go           zstyle integration, named directories
  pflagfork/                Non-POSIX flag mode handling
  spec/                     YAML spec generation from cobra commands
  config/                   Runtime config loading
  env/                      Environment variable accessors (CARAPACE_* constants)
  mock/                     Sandbox test mocking
  cache/                    File-based caching
  log/                      Logging
  man/                      Man page integration
  export/                   Export struct (JSON wire format)

/pkg/                        Public sub-packages (importable by consumers)
  sandbox/                  Test sandbox: Command(), Package(), Action()
  style/                    Style system (colors, semantic: ForPath, ForKeyword, ForLogLevel)
  uid/                      Unique identifiers (URL-based: cmd://host/path?flag=name)
  match/                    Prefix matching (CARAPACE_MATCH env)
  assert/                   Test assertions with myers diff
  cache/ /key/              Cache key types
  condition/                Conditional completion helpers
  ps/                       Shell detection from process tree
  execlog/                  exec.Command wrapper with logging
  util/                     HasVolumePrefix, etc.
  x/                        Cross-package hooks (ClearStorage, Complete)
  xdg/                      XDG directory resolution
  traverse/                 Context type for external traverse callbacks

/example/                    Integration test app (posix mode)
/example-nonposix/           Integration test app (non-posix mode, go.work module)
/third_party/                Vendored code (elvish UI, gotextdiff, envsubst, stripansi, etc.)
/skills/                     Agent skill definitions
  carapace-dev/                 Composite skill for carapace library development
    SKILL.md                      Entry point with routing table
    references/action.md          Action API internals
    references/traverse.md        Completion engine & traversal
    references/style.md           Styling system internals
    references/shell.md           Per-shell overview & cross-shell comparisons
    references/shell-bash.md     Bash deep dive
    references/shell-bash-ble.md Bash BLE deep dive
    references/shell-oil.md      Oil deep dive
    references/shell-elvish.md   Elvish deep dive
    references/shell-fish.md    Fish deep dive
    references/shell-nushell.md  Nushell deep dive
    references/shell-powershell.md PowerShell deep dive
    references/shell-xonsh.md    Xonsh deep dive
    references/shell-zsh.md      Zsh deep dive
    references/pflag.md          carapace-pflag extensions (from carapace-pflag repo)
    references/bridge.md         carapace-bridge internals (from carapace-bridge repo)
```

## Conventions

### Action Naming Suffixes

- **No suffix** — basic constructor (e.g., `ActionValues`)
- **`F`** — takes a function (e.g., `StyleF`, `ChdirF`, `UidF`)
- **`P`** — takes placeholders (e.g., `MultiPartsP`)
- **`E`** — error-returning variant (e.g., `ActionExecCommandE`)
- **`N`** — name variant for pflag (e.g., `BoolN`, `StringN` from carapace-pflag)

### Modifier Chaining

Modifiers return new `Action` values wrapping the original callback. They are lazy — the inner callback is only executed when `Invoke()` is called:

```go
carapace.ActionValues("a", "b", "c").
    Filter("b").
    StyleF(style.ForKeyword).
    Tag("keywords").
    NoSpace().
    Usage("some flag")
```

### Defining Completions

```go
carapace.Gen(rootCmd)
    .FlagCompletion(carapace.ActionMap{
        "flag-name": carapace.ActionValues("a", "b"),
    })
    .PositionalCompletion(
        carapace.ActionValues("first"),
        carapace.ActionValues("second"),
    )
    .PositionalAnyCompletion(carapace.ActionFiles())
    .DashCompletion(carapace.ActionValues("after-dash"))
    .DashAnyCompletion(carapace.ActionFiles())
```

### PreRun / PreInvoke

- **PreRun** — called before traversal; use for dynamically adding commands/flags. Multiple handlers chain sequentially.
- **PreInvoke** — transforms actions before invocation; receives `(cmd, flag, action)`. Chains from child to parent. Common use: `action.Chdir(flagValue)`.

## Testing

### Sandbox Pattern

```go
sandbox.Command(t, func() *cobra.Command {
    cmd := &cobra.Command{Use: "myapp"}
    carapace.Gen(cmd).PositionalCompletion(carapace.ActionValues("a", "b"))
    return cmd
})(func(s *sandbox.Sandbox) {
    s.Files("test.txt", "content")
    s.Reply("git", "remote").With("origin\nfork")
    s.Run("", "").Expect(carapace.ActionValues("a", "b"))
})
```

- `sandbox.Command(t, cmdF)` — test a cobra command directly
- `sandbox.Package(t, "pkg/path")` — test via `go run` (integration)
- `sandbox.Action(t, actionF)` — test a standalone action
- `s.Files()` — create files in sandbox directory (path, content pairs)
- `s.Reply()` — mock command output (works with `Context.Command()` only)
- `s.Run(args...)` — execute completion
- `run.Expect(action)` — assert output matches expected action
- `run.ExpectNot(action)` — assert output differs

### Assert Package

`assert.Equal(t, expected, actual)` serializes both to JSON and shows a myers diff on mismatch. Expected actions in tests must include all modifiers (Style, NoSpace, Tag, Chdir, Usage) for exact matching.

### Validation

```go
func TestCarapace(t *testing.T) {
    carapace.Test(t) // validates that flag names in ActionMap exist on the command
}
```

## Gotchas

### Storage is a global map

`_storage` is a package-level `map[*cobra.Command]*entry` with mutex protection. Tests that register completions may conflict. Sandbox tests avoid this by running in subprocesses.

### uid.Flag() uses a mutex

`uid.Flag()` acquires `mLocalFlags` because `cmd.LocalFlags()` triggers `mergePersistentFlags()` which is not thread-safe internally.

### Context.Value is the last arg

`NewContext(args...)` treats the last argument as the value being completed. An empty args slice becomes `Value: ""`.

### Shell patching before traversal

`complete()` applies shell-specific patches to args (bash redirects, nushell quotes, cmd-clink) before `traverse()`. The traversal sees modified args, not raw shell input.

### go.work includes example-nonposix

The `go.work` file adds `./example-nonposix` as a separate module. CI tests both: `go test ./... ./example-nonposix/...`.

### Cache keys come from caller location

`Action.Cache()` uses `runtime.Caller(1)` (file:line) as cache key. Moving a `.Cache()` call changes the key.

### Sandbox tests cannot run in parallel

Storage mutations make `t.Parallel()` unsafe. There are TODO comments about this in sandbox code.

### The `export` shell target

`internal/shell/export/` outputs completion results as JSON. The `_carapace` subcommand re-invokes itself with the `export` shell to resolve subcommand completions.

### third_party/ is excluded from linting

`.golangci.yml` excludes `third_party/`. Do not modify vendored code.

### Non-POSIX flag modes

`pflagfork` handles three modes via reflection on carapace-pflag extensions:
- **Default** — standard POSIX (`-s`, `--long`)
- **ShorthandOnly** — only shorthand valid
- **NameAsShorthand** — name with single dash (non-POSIX, e.g. `-bool-long`)

In non-POSIX mode, longhand lookup is tried before shorthand to correctly handle cases where a flag name overlaps with its shorthand.

### Environment variables

Key env vars (see `internal/env/env.go`):
- `CARAPACE_LOG` — enable debug logging
- `CARAPACE_SANDBOX` — JSON mock context (set by sandbox tests)
- `CARAPACE_COVERDIR` — coverage directory for integration tests
- `CARAPACE_LENIENT` — allow unknown flags
- `CARAPACE_MATCH` — set to `CASE_INSENSITIVE` for case-insensitive matching
- `CARAPACE_NOSPACE` — additional nospace suffixes
- `CARAPACE_UNFILTERED` — skip prefix filtering
- `CARAPACE_EXPERIMENTAL` — enable experimental features
- `CARAPACE_HIDDEN` — show hidden cmds/flags (1=exclude carapace, 2=include)
- `CARAPACE_TOOLTIP` — enable tooltip style
- `CARAPACE_DESCRIPTION_LENGTH` — max description length (default 80)
- `NO_COLOR` / `CLICOLOR=0` — disable colors

## Shell Skill Maintenance

Shell integration documentation lives in the `carapace-dev` composite skill (`skills/carapace-dev/`). The generic overview is `references/shell.md` and per-shell deep dives are `references/shell-{name}.md` (bash, bash-ble, oil, zsh, fish, elvish, nushell, xonsh, powershell). Minor shells (tcsh, ion, cmd-clink, export) are covered only in `references/shell.md`.

### Structure

- **`references/shell.md`** — Generic overview, shared dispatch pipeline, and cross-shell comparison tables. Does NOT contain shell-specific implementation details (quoting rules, snippet walkthroughs, etc.).
- **`references/shell-{shell}.md`** — Shell-specific deep dive. Contains ONLY information unique to that shell. Follows a consistent section order:
  1. Source Files
  2. Shell Background
  3. The Snippet (with walkthrough)
  4. Patch Phase (if applicable — only bash, nushell, cmd-clink have `Patch()`)
  5. Value Formatting (`ActionRawValues()`)
  6. Nospace Handling
  7. Message Handling
  8. Edge Cases and Known Issues
  9. Completion Dispatch Flow
  10. References
  11. Related Skills

### Bash Family

`references/shell-bash.md` covers regular bash. Bash BLE and Oil have their own deep dives:
- **`references/shell-bash.md`** — bash only (COMP_TYPE, COMP_WORDBREAKS, redirect patching, quoting, partial completion workaround)
- **`references/shell-bash-ble.md`** — bash BLE (tab-delimited format, per-candidate suffix, ble/complete/cand/yield)
- **`references/shell-oil.md`** — oil (simpler snippet, inline `\001` nospace indicator, no patching)

### When Adding a New Shell Skill

1. Add `references/shell-{name}.md` to `skills/carapace-dev/` following the section order above
2. Include shell-specific details only — link to `references/shell.md` for cross-shell comparisons
3. Add the shell to the "Supported Shells" table in `references/shell.md`
4. Add any new comparison rows if the shell introduces a new mechanism
5. Update the routing table in `skills/carapace-dev/SKILL.md`

### When Updating Shell Skills

- **`references/shell.md`**: Update when adding/removing shells, or when a cross-shell mechanism changes
- **`references/shell-{shell}.md`**: Update when the shell's formatter, snippet, or patch logic changes
- **Avoid duplication**: Link to `references/shell.md` rather than repeating cross-shell comparisons
