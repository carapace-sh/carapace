# Carapace Library: Mock System

Reference for [carapace](https://github.com/carapace-sh/carapace)'s mock system in `internal/mock/` used by sandbox tests.

## Mock Struct

```go
type Mock struct {
    Dir     string
    Replies map[string]string
}
```

The mock holds:
- **Dir**: the sandbox temp directory path
- **Replies**: a map of JSON-serialized command args → output strings

## Directory Structure

```
<tempdir>/carapace-sandbox_<testname>_xxxxx/
├── cache/   ← mock.CacheDir()
└── work/    ← mock.WorkDir()
```

## Mock Methods

```go
func (m Mock) CacheDir() string {
    return m.Dir + "/cache/"
}

func (m Mock) WorkDir() string {
    return m.Dir + "/work/"
}
```

## NewMock — Create a Mock

```go
func NewMock(t t) *Mock {
    tempDir, err := os.MkdirTemp(os.TempDir(), fmt.Sprintf("carapace-sandbox_%v_", t.Name()))
    if err != nil {
        t.Fatal("failed to create sandbox dir: " + err.Error())
    }

    m := &Mock{
        Dir:     tempDir,
        Replies: make(map[string]string),
    }
    if err := os.Mkdir(m.CacheDir(), os.ModePerm); err != nil {
        t.Fatal("failed to create sandbox cache dir: " + err.Error())
    }
    if err := os.Mkdir(m.WorkDir(), os.ModePerm); err != nil {
        t.Fatal("failed to create sandbox work dir: " + err.Error())
    }
    return m
}
```

Creates the temp directory and `cache/` and `work/` subdirectories. Uses `t.Fatal` for error reporting (idiomatic Go test pattern).

## Sandbox Integration

The mock is JSON-marshaled and injected as `CARAPACE_SANDBOX`:

```go
func (s *Sandbox) Run(args ...string) run {
    m, _ := json.Marshal(s.mock)
    r := run{
        t:       s.t,
        id:      string(m),
        dir:     s.mock.WorkDir(),
        context: s.NewContext(args...),
    }

    r.actual = carapace.ActionCallback(func(c carapace.Context) carapace.Action {
        b, err := json.Marshal(s.mock)
        if err != nil {
            return carapace.ActionMessage(err.Error())
        }
        c.Setenv("CARAPACE_SANDBOX", string(b))
        return carapace.ActionExecute(s.cmdF()).Invoke(c).ToA()
    }).Invoke(r.context).ToA()

    return r
}
```

`CARAPACE_SANDBOX` is read by `env.Sandbox()` to get the mock directory and replies.

## Reply Mechanism

```go
func (s *Sandbox) Reply(args ...string) reply {
    m, _ := json.Marshal(args)
    return reply{s, string(m)}
}

func (r reply) With(s string) {
    r.mock.Replies[r.call] = s
}
```

Command arguments are JSON-serialized as the map key (e.g., `["git", "remote"]` → `"[\"git\",\"remote\"]"`). The output string is stored as the value.

At invocation time, `Context.Command()` checks `env.Sandbox()` for mocked replies.

## The t Interface

```go
type t interface {
    Name() string
    Fatal(...any)
}
```

Minimal interface compatible with `*testing.T`. Allows the mock to work with any testing framework that provides `Name()` and `Fatal()`.

## How Context.Command() Uses Mock

`Context.Command()` (in `context.go`) checks for `CARAPACE_SANDBOX` and uses mocked replies when available. The actual implementation in `Context.Command()` reads `env.Sandbox()` and looks up the command key in `m.Replies`.

## Gotchas

- **JSON-serialized keys**: `s.Reply("git", "remote")` stores under `"[\"git\",\"remote\"]"`, not the plain string. Exact match required.
- **Temp directory cleanup**: `s.remove()` in `sandbox.go` only removes directories under `os.TempDir()`. If `s.Keep()` was called, the directory is retained.
- **No mocking for exec.LookPath**: `condition.Executable()` uses `exec.LookPath` directly, bypassing the sandbox mock. Only `Context.Command()` is mocked.
- **No mocking for os.Stat**: File existence checks (`condition.File()`) use real filesystem, not sandbox files. Use `s.Files()` to create real files in the sandbox's `work/` directory instead.
- **Cache dir in mock**: Each mock has its own `cache/` subdirectory. `s.ClearCache()` removes it.

## Related Skills

- **references/sandbox.md** — how mock integrates into the sandbox
- **internal/env/** — env.Sandbox() reads CARAPACE_SANDBOX
- **pkg/sandbox/** — sandbox uses mock.NewMock