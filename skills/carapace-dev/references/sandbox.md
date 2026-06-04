# Carapace Library: Sandbox Testing

Reference for [carapace](https://github.com/carapace-sh/carapace)'s sandbox testing harness in `pkg/sandbox/`.

## Overview

The sandbox runs completion tests in an isolated subprocess with a temporary directory, mocked environment, and controlled command replies. It avoids conflicts with the global `storage` map by running each test as a separate process.

## Entry Points

### Command — test a cobra command directly

```go
sandbox.Command(t, func() *cobra.Command {
    cmd := &cobra.Command{Use: "myapp"}
    carapace.Gen(cmd).PositionalCompletion(carapace.ActionValues("a", "b"))
    return cmd
})(func(s *sandbox.Sandbox) {
    s.Files("test.txt", "content")
    s.Run("", "").Expect(carapace.ActionValues("a", "b"))
})
```

### Package — test via `go run` (integration)

```go
sandbox.Package(t, "example")(func(s *sandbox.Sandbox) {
    s.Run("myapp", "subcmd", "").Expect(carapace.ActionValues("x", "y"))
})
```

### Action — test a standalone action

```go
sandbox.Action(t, func() carapace.Action {
    return carapace.ActionValues("a", "b", "c").StyleF(style.ForKeyword)
})(func(s *sandbox.Sandbox) {
    s.Run("").Expect(carapace.ActionValues("a", "b", "c").StyleF(style.ForKeyword))
})
```

## Sandbox Setup

`newSandbox(t, cmdF)` creates:

- A `*mock.Mock` with a temp directory at `os.TempDir()/carapace-sandbox_<testname>_`
- `cache/` subdirectory for cache tests
- `work/` subdirectory for file tests
- A `*testing.T` reference for error reporting

The sandbox directory is auto-removed on test completion unless `s.Keep()` is called.

## Files — Create Test Files

```go
s.Files(
    "file1.txt", "content of file1.txt",
    "dir1/file2.md", "content of file2.md",
)
```

Files are created relative to the sandbox's `work/` directory. Filenames with `..` or leading `/` are rejected.

## Env — Set Environment Variables

```go
s.Env("CARAPACE_UNFILTERED", "true")
s.Env("HOME", "/tmp/testhome")
```

Set on the `Context` passed to action invocation.

## Reply — Mock Command Output

```go
s.Reply("git", "remote").With("origin\nfork")
s.Reply("docker", "ps", "-q").With("abc123\ndef456")
```

Mocks `Context.Command()` output. Calls are JSON-serialized as the map key (`["git", "remote"]`). Only works with `Context.Command()`, not raw `exec.Command`.

## Run — Execute Completion

```go
s.Run(args ...string) run
```

Creates a `Context` from the args (last arg = `Value`, rest = `Args`), injects `CARAPACE_SANDBOX` with the mock JSON, and executes `ActionExecute(cmdF())` via the sandbox.

Returns a `run` value with `Expect`, `ExpectNot`, and `Output` methods.

## Expect — Assert Output

```go
s.Run("", "").Expect(carapace.ActionValues("a", "b"))
s.Run("", "").ExpectNot(carapace.ActionValues("x"))
```

Serializes both the expected and actual actions to JSON via `export.Export` and compares with `assert.Equal`. The test is named after the args tuple (e.g., `run["",""]`).

## ClearCache — Clear Cache Between Tests

```go
s.ClearCache()
```

Removes the sandbox's `cache/` directory. Use between subtests to ensure cache-dependent tests don't interfere.

## Sandbox Env Vars

The sandbox sets `CARAPACE_SANDBOX` containing a JSON-encoded `*mock.Mock`:

```json
{
  "Dir": "/tmp/carapace-sandbox_testname_xxxxx",
  "Replies": { "[\"git\",\"remote\"]": "origin\nfork" }
}
```

`env.Sandbox()` reads and unmarshals this variable. The `isGoRun()` check ensures sandbox mode only activates under `go run` (not compiled binaries).

## Architecture

```
sandbox.Command/Package/Action
  └─ newSandbox
       └─ mock.NewMock  → creates temp dir, cache/, work/
  └─ defer s.remove()   → cleanup unless s.Keep()
  └─ f(&sandbox)       → test body
       ├─ s.Files()
       ├─ s.Reply().With()
       ├─ s.Env()
       └─ s.Run(args...)
            ├─ s.NewContext(args...)
            │    └─ Context{Dir: mock.WorkDir(), Env: s.env}
            └─ ActionExecute(cmdF())
                 └─ ActionCallback
                      └─ sets CARAPACE_SANDBOX=json(mock)
                           └─ env.Sandbox() in cache.File()
```

## Gotchas

- **No parallel tests**: `t.Parallel()` is explicitly TODO-commented — storage mutations make concurrent tests unsafe.
- **SANDBOX only under go run**: `env.Sandbox()` checks `isGoRun()` which matches `/go-build` in `os.Args[0]`. Compiled binaries skip sandbox mode.
- **Reply keys are JSON arrays**: `s.Reply("git", "remote")` serializes to `"[\"git\",\"remote\"]"` as the map key. Exact match required.
- **Cache keys include file:line**: `Action.Cache()` uses `runtime.Caller(1)` — moving a `.Cache()` call changes the key.
- **`s.Keep()` for debugging**: Call before the test body to retain the temp directory for inspection.

## Related Skills

- **references/action.md** — Action API tested via sandbox
- **references/cache.md** — cache integration in sandbox
- **pkg/assert/** — assertion library used by `Expect`