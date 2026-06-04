# Carapace Library: Cache System

Reference for [carapace](https://github.com/carapace-sh/carapace)'s file-based cache for action callback results.

## Action.Cache — Lazy Caching Modifier

```go
func (a Action) Cache(timeout time.Duration, keys ...key.Key) Action {
    if a.callback != nil {
        cachedCallback := a.callback
        _, file, line, _ := runtime.Caller(1)
        a.callback = func(c Context) Action {
            cacheFile, err := cache.File(file, line, keys...)
            if err != nil {
                return cachedCallback(c)
            }

            if cached, err := cache.LoadE(cacheFile, timeout); err == nil {
                return Action{meta: cached.Meta, rawValues: cached.Values}
            }

            invokedAction := (Action{callback: cachedCallback}).Invoke(c)
            if invokedAction.action.meta.Messages.IsEmpty() {
                if cacheFile, err := cache.File(file, line, keys...); err == nil {
                    _ = cache.WriteE(cacheFile, invokedAction.export())
                }
            }
            return invokedAction.ToA()
        }
    }
    return a
}
```

Only wraps actions with a `callback` (pre-computed `rawValues` skip caching). Uses `runtime.Caller(1)` to get the caller's file:line as the cache key base.

## Cache Key Structure

Cache files are stored at:

```
$XDG_CACHE_HOME/carapace/<executable>/<filehash>/<keyhash>
```

Where:
- `<executable>` = `uid.Executable()` — name of the compiled binary
- `<filehash>` = SHA1 of `file\001line` (the caller's `runtime.Caller(1)` location)
- `<keyhash>` = SHA1 of all `key.Key` values joined by `\001`

Example path:

```
~/.cache/carapace/git/3a2f1c9b/1d4e5f7a2b3c...
```

## Cache Key Types

`pkg/cache/key/cache.go` provides `Key` func types:

| Key | Purpose |
|-----|---------|
| `key.String(s ...string)` | Static strings as key components |
| `key.FileChecksum(file)` | SHA1 of file content |
| `key.FileStats(file)` | `path + size + mtime` for file changes |
| `key.FolderStats(folder)` | Aggregate stats of all files in folder |

```go
// Cache until go.mod changes
ActionExecCommand("go", "list", "-m")(...).Cache(10*time.Minute, key.FileChecksum("go.mod"))

// Cache until any file in vendor/ changes
ActionExecCommand("ls", "vendor/")(...).Cache(5*time.Minute, key.FolderStats("vendor"))
```

## TTL / Timeout

```go
cache.Load(file, timeout time.Duration)
```

`timeout < 0` means infinite (never expire). `timeout = 0` means always revalidate. `timeout > 0` expires entries older than the duration.

## Cache File Format

Cache stores a JSON-encoded `export.Export`:

```go
type Export struct {
    Version string       `json:"version"`
    Meta    common.Meta  `json:"meta"`
    Values  common.RawValues `json:"values"`
}
```

Same format as the `export` shell target. Written via `cache.WriteE()` which JSON-marshals and calls `os.WriteFile` with `0600` permissions.

## Sandbox Override

When `CARAPACE_SANDBOX` is set, `cache.CacheDir()` redirects to the sandbox's cache directory:

```go
func CacheDir(name string) (dir string, err error) {
    userCacheDir, err = xdg.UserCacheDir()
    if m, sandboxErr := env.Sandbox(); sandboxErr == nil {
        userCacheDir = m.CacheDir() // sandbox temp dir
    }
    dir = fmt.Sprintf("%v/carapace/%v/%v", userCacheDir, uid.Executable(), name)
    os.MkdirAll(dir, 0700)
    return
}
```

Sandbox cache dir: `<sandbox_tmpdir>/cache/carapace/<executable>/...`

## Coverage Integration

`CARAPACE_COVERDIR` is a custom env variable (distinct from `GOCOVERDIR`) that tells the sandbox where to write coverage data:

```go
func CoverDir() string {
    return os.Getenv("CARAPACE_COVERDIR")
}
```

Used by integration tests to collect coverage from sandbox subprocesses alongside unit test coverage.

## Write Path

`cache.WriteE()` uses JSON marshal, then `Write()` which calls `os.WriteFile`. No intermediate temp file — directly writes to the target path.

## Gotchas

- **Cache key = file:line**: Moving a `.Cache()` call changes the cache key and invalidates all cached entries for that action. Cache is per-call-site, not per-action-definition.
- **Messages suppress cache write**: If `Messages` is not empty, the cache is not written. This prevents caching error states or informational messages.
- **Cache keys regenerated after invoke**: The second call to `cache.File()` inside the callback regenerates the key in case the keys depend on `Context` values that changed between invocation start and now.
- **0600 permissions**: Cache files are owner-read/write only to prevent exposure of potentially sensitive completion data.
- **Sandbox isolation**: Each sandbox test gets its own cache directory via `mock.CacheDir()`. `s.ClearCache()` removes it between subtests.
- **Timeout = 0**: Means always revalidate (stat the file and check mtime). The `Load` function returns error if the file is older than `timeout`.

## Related Skills

- **references/action.md** — the `Cache()` modifier on Action
- **references/sandbox.md** — sandbox cache directory and ClearCache
- **internal/export/** — the Export struct stored in cache files