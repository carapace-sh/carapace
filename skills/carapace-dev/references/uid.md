# Carapace Library: UID System

Reference for [carapace](https://github.com/carapace-sh/carapace)'s URL-based unique identifiers in `pkg/uid/`.

## Overview

UIDs are `url.URL` values with scheme `cmd` or `file`, providing stable, globally unique identifiers for commands and flags. They are used in cache key derivation, error messages, and debugging output.

## URL Structure

```go
// Command UID
cmd://rootcmd/subcmd?flag=flagname

// Example: git branch
cmd://git/branch

// Flag: git branch --force
cmd://git/branch?flag=force
```

- **Scheme**: `cmd` for commands, `file` for files (used in `Action.Uid()`)
- **Host**: root command name
- **Path**: subcommand path (slash-separated ancestors)
- **Query**: `flag=<name>` for flag UIDs

## uid.Command

```go
func Command(cmd *cobra.Command) *url.URL {
    path := []string{cmd.Name()}
    for parent := cmd.Parent(); parent != nil; parent = parent.Parent() {
        path = append(path, url.PathEscape(parent.Name()))
    }
    reverse(path)
    return &url.URL{
        Scheme: "cmd",
        Host:   path[0],
        Path:   strings.Join(path[1:], "/"),
    }
}
```

Builds the UID by walking up the cobra parent chain. The root command is always the host.

## uid.Flag

```go
var mLocalFlags sync.Mutex

func Flag(cmd *cobra.Command, flag *pflagfork.Flag) *url.URL {
    mLocalFlags.Lock()
    defer mLocalFlags.Unlock()
    return flagRecursive(cmd, flag)
}

func flagRecursive(cmd *cobra.Command, flag *pflagfork.Flag) *url.URL {
    _ = cmd.LocalFlags() // Force flag merge; not thread-safe internally

    if cmd.LocalFlags().Lookup(flag.Name) == nil && cmd.HasParent() {
        return flagRecursive(cmd.Parent(), flag)
    }
    uid := Command(cmd)
    values := uid.Query()
    values.Set("flag", flag.Name)
    uid.RawQuery = values.Encode()
    return uid
}
```

The mutex is necessary because `cmd.LocalFlags()` calls `mergePersistentFlags()` internally which is not thread-safe. The lock is held for the entire recursive traversal to prevent concurrent modification.

## uid.UidF — Experimental UID Function

```go
func UidF(scheme, host string, opts ...string) func(v string, uc Context) (*url.URL, error) {
    return func(v string, uc Context) (*url.URL, error) {
        uid := &url.URL{
            Scheme: scheme,
            Host:   url.PathEscape(host),
            Path:   PathEscape(v),
        }
        if len(opts) > 0 {
            values := uid.Query()
            for i := 0; i < len(opts); i += 2 {
                if opts[i+1] != "" {
                    values.Set(opts[i], opts[i+1])
                }
            }
            uid.RawQuery = values.Encode()
        }
        return uid, nil
    }
}
```

Used by `Action.Uid()` and `Action.UidF()`. Takes alternating key-value opts (e.g., `"flag", "name", "group", "build"`). Empty values are skipped.

## PathEscape

```go
func PathEscape(s string) string {
    s = strings.ReplaceAll(s, "/", "-")
    s = url.PathEscape(s)
    return s
}
```

Replaces `/` with `-` before escaping, since `/` is the path separator in the UID format. This allows values containing slashes (e.g., file paths) to be encoded as part of the path.

## uid.Map — Testing Helper

```go
func Map(uids ...string) func(s string) (*url.URL, error) {
    return func(s string) (*url.URL, error) {
        for i := 0; i < len(uids); i += 2 {
            if uids[i] == s {
                return url.Parse(uids[i+1])
            }
        }
        return &url.URL{}, nil
    }
}
```

Maps static values to UIDs for deterministic testing. Returns an empty UID for unmatched values.

## uid.Executable

```go
func Executable() string {
    executable, err := os.Executable()
    if err != nil {
        return "echo" // safe fallback
    }
    switch base := filepath.Base(executable); base {
    case "cmd.test":
        return "example" // for `go test -v ./...`
    case "ld-musl-x86_64.so.1":
        return filepath.Base(os.Args[0]) // alpine container workaround
    default:
        return base
    }
}
```

Used as the cache directory name component. Special-cases test binaries and Alpine's musl loader.

## Context Interface

```go
type Context interface {
    Abs(s string) (string, error)
    Getenv(key string) string
    LookupEnv(key string) (string, bool)
}
```

Passed to `UidF` functions at invocation time. Implemented by `carapace.Context`.

## Gotchas

- **mLocalFlags mutex**: `cmd.LocalFlags()` is not thread-safe internally. The mutex protects the entire recursive traversal. A recent fix (commit 324b32b) addressed a race condition by holding the mutex during recursion.
- **PathEscape converts `/` to `-`**: A value like `users/admin` becomes `users-admin` in the UID path. This is necessary because `/` is the path separator.
- **Empty UID for unmatched Map values**: `uid.Map()` returns `&url.URL{}` for values not in the map, not an error.
- **cmd.test → example**: During `go test`, `os.Executable()` returns the test binary name containing `.test`. The `cmd.test` case maps this to `example` for consistency with the integration test binary name.

## Related Skills

- **references/action.md** — Uid/UidF modifiers on Action
- **references/storage.md** — where uid.Command is used for error reporting
- **references/cache.md** — cache key derivation uses uid.Executable