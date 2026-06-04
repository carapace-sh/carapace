# Carapace Library: Conditional Helpers

Reference for [carapace](https://github.com/carapace-sh/carapace)'s conditional completion helpers in `pkg/condition/`.

## Overview

These helpers produce `func(carapace.Context) bool` predicates used with `Action.UnlessF()` or `Action.Unless()` to conditionally suppress or include completions based on runtime context.

## Arch — GOARCH Filter

```go
func Arch(s ...string) func(c carapace.Context) bool {
    return func(c carapace.Context) bool {
        return slices.Contains(s, runtime.GOARCH)
    }
}
```

```go
ActionValues("linux", "darwin", "windows").
    UnlessF(condition.Arch("linux", "darwin"))
```

## Os — GOOS Filter

```go
func Os(s ...string) func(c carapace.Context) bool {
    return func(c carapace.Context) bool {
        return slices.Contains(s, runtime.GOOS)
    }
}
```

## Executable — PATH Lookup

```go
func Executable(s ...string) func(c carapace.Context) bool {
    return func(c carapace.Context) bool {
        for _, executable := range s {
            if _, err := exec.LookPath(executable); err == nil {
                return true
            }
        }
        return false
    }
}
```

```go
ActionValues("docker", "podman").
    UnlessF(condition.Executable("docker"))
```

Note: Uses `exec.LookPath` directly instead of `Context.Command()`, so it uses the real PATH, not the sandbox's mocked environment.

## File — File or Directory Existence

```go
func File(s string) func(c carapace.Context) bool {
    return func(c carapace.Context) bool {
        if _, err := os.Stat(s); err == nil {
            return true
        }
        return false
    }
}
```

```go
ActionValues("force", "interactive").
    UnlessF(condition.File("/etc/myapp.conf"))
```

## CompletingPath — Path Prefix Detection

```go
func CompletingPath(c carapace.Context) bool {
    return util.HasPathPrefix(c.Value)
}
```

Returns true when `Context.Value` starts with `/`, `./`, `../`, or `~`.

## CompletingPathS — Path with Separator

```go
func CompletingPathS(c carapace.Context) bool {
    return CompletingPath(c) || strings.Contains(c.Value, "/")
}
```

Like `CompletingPath` but also triggers when the value contains a `/` anywhere (not just as a prefix). Useful for completions where partial path segments are being completed.

## Usage with Unless / UnlessF

All conditionals are predicates for use with `Unless`:

```go
ActionCallback(func(c carapace.Context) carapace.Action {
    if condition.Executable("docker")(c) {
        return ActionValues("start", "stop", "rm")
    }
    return ActionMessage("docker not found")
})
```

`UnlessF` applies the inverse:

```go
ActionValues("linux-amd64", "darwin-arm64").
    UnlessF(condition.Arch("linux")).
    UnlessF(condition.Os("windows"))
```

## As Modifiers

Conditionals can also be composed inline:

```go
carapace.ActionValues("a", "b").
    UnlessF(func(c carapace.Context) bool {
        return runtime.GOOS == "windows"
    })
```

## Future / Additional Helpers

Additional conditionals may exist in the package (check `pkg/condition/*.go`). Common patterns not yet present but reasonable to add:

- `Env(key, value)` — check environment variable
- `FlagChanged(name)` — check if a flag was set
- `HasPrefix(s)` — check if value has a prefix
- `ToLower` — case conversion predicates

## Gotchas

- **Executable uses real PATH**: `exec.LookPath` reads the actual system PATH, not `Context.Env`. In sandbox tests, the real PATH is still used.
- **File uses absolute paths**: `os.Stat` requires an absolute path. Use `Context.Abs()` to resolve relative paths first.
- **No sandbox mocking for conditionals**: Unlike `Context.Command()`, conditionals directly call `os`, `exec`, `runtime` — no mocking via `CARAPACE_SANDBOX`.

## Related Skills

- **references/action.md** — Unless/UnlessF modifiers
- **pkg/util/** — HasPathPrefix utility used by CompletingPath