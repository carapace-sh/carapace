# Carapace Library: Storage & Hooks Internals

Reference for [carapace](https://github.com/carapace-sh/carapace)'s per-command completion storage — the global map, mutex protection, registration API, and PreRun/PreInvoke chains.

## The entry Struct

Every cobra command registered with `carapace.Gen()` gets an `entry` stored in a package-level `map[*cobra.Command]*entry`:

```go
type entry struct {
    flag          ActionMap
    flagMutex     sync.RWMutex
    positional    []Action
    positionalAny *Action
    dash          []Action
    dashAny       *Action
    preinvoke     func(cmd *cobra.Command, flag *pflag.Flag, action Action) Action
    prerun        func(cmd *cobra.Command, args []string)
    bridged       bool
    initialized   bool
}
```

The storage map is private (`_storage`) with access routed through methods.

## Global Storage Map

```go
var storageMutex sync.RWMutex

func (s _storage) get(cmd *cobra.Command) *entry {
    storageMutex.RLock()
    e, ok := s[cmd]
    storageMutex.RUnlock()

    if !ok {
        storageMutex.Lock()
        defer storageMutex.Unlock()
        if e, ok = s[cmd]; !ok {
            e = &entry{}
            s[cmd] = e
        }
    }
    return e
}
```

Double-checked locking: read-lock first, then write-lock only if the entry doesn't exist. This avoids locking on every access after warm-up.

## Registration API

`carapace.Gen(cmd)` returns a `Carapace` wrapper that provides the registration chain:

```go
type Carapace struct{ *cobra.Command }

func (c Carapace) FlagCompletion(actions ActionMap) Carapace
func (c Carapace) PositionalCompletion(actions ...Action) Carapace
func (c Carapace) PositionalAnyCompletion(action Action) Carapace
func (c Carapace) DashCompletion(actions ...Action) Carapace
func (c Carapace) DashAnyCompletion(action Action) Carapace
func (c Carapace) PreRun(fn func(cmd *cobra.Command, args []string)) Carapace
func (c Carapace) PreInvoke(fn func(cmd *cobra.Command, flag *pflag.Flag, action Action) Action) Carapace
```

Each method calls `storage.get(cmd)` and mutates the corresponding entry field, protected by the per-entry `flagMutex` for flag actions.

## PreRun — Before Traversal

Executed in `storage.preRun()` before `traverse()` walks the command tree. Use for dynamically adding/removing subcommands or flags:

```go
carapace.Gen(rootCmd).PreRun(func(cmd *cobra.Command, args []string) {
    if pluginEnabled {
        cmd.AddCommand(pluginCommand)
    }
})
```

Called from `complete()` before any argument parsing or traversal.

## PreInvoke — Before Action Invocation

Executed for every action resolution (flag or positional) in `storage.preinvoke()`. Chains from child command to parent, so a PreInvoke on the root command affects all subcommands:

```go
carapace.Gen(rootCmd).PreInvoke(func(cmd *cobra.Command, flag *pflag.Flag, action carapace.Action) carapace.Action {
    if flag != nil && flag.Name == "chdir" && cmd.Flag("chdir").Changed {
        return action.Chdir(cmd.Flag("chdir").Value.String())
    }
    return action
})
```

The `flag` parameter is `nil` for positional arguments. The chain runs after `traverse()` has classified the argument but before the action callback is executed.

## Flag Lookup

`storage.getFlag(cmd, name)` looks up a flag action, checking parent commands recursively:

```go
func (s _storage) getFlag(cmd *cobra.Command, name string) Action {
    if flag := cmd.LocalFlags().Lookup(name); flag == nil && cmd.HasParent() {
        return s.getFlag(cmd.Parent(), name)
    } else {
        entry := s.get(cmd)
        entry.flagMutex.RLock()
        defer entry.flagMutex.RUnlock()

        flagAction, ok := entry.flag[name]
        if !ok {
            if f, ok := cmd.GetFlagCompletionFunc(name); ok {
                flagAction = ActionCobra(f)
            }
        }
        a := s.preinvoke(cmd, flag, flagAction)
        return ActionCallback(func(c Context) Action {
            invoked := a.Invoke(c)
            if invoked.action.meta.Usage == "" {
                invoked.action.meta.Usage = flag.Usage
            }
            return invoked.ToA()
        })
    }
}
```

Falls back to cobra's native `GetFlagCompletionFunc` if no carapace action is registered.

## Positional Lookup

`storage.getPositional(cmd, index)` and `storage.hasPositional(cmd, index)` handle both normal and `--` (dash) positionals. Dash positionals use a separate array (`entry.dash`/`entry.dashAny`) activated when `common.IsDash(cmd)` returns true.

## Validation

`storage.check()` iterates all stored flag names and verifies each one exists on the command's `LocalFlags()`. Called by `carapace.Test()`. Reports unknown flags with `uid.Command(cmd)` as the context.

## Bridge Integration

`storage.bridge(cmd)` wires cobra's native completion functions (`ValidArgsFunction`, `RegisterFlagCompletionFunc`) via `cobra.OnInitialize()`. Called lazily on first completion request. Uses a separate `bridgeMutex` to prevent concurrent initialization.

## Gotchas

- **Global map**: `_storage` is a package-level global. Tests that register completions in parallel (`t.Parallel()`) will conflict. Sandbox tests run in subprocesses to avoid this.
- **Per-entry mutex**: Flag actions use `entry.flagMutex.RLock()`, but positional/dash actions have no mutex. This is safe because positional registration is single-threaded after init.
- **Cache keys from caller location**: `Action.Cache()` uses `runtime.Caller(1)` as the cache key base. Moving a `.Cache()` call to a different file:line changes the cache key and invalidates existing cache entries.
- **PreInvoke chain order**: PreInvoke chains from child → parent. A PreInvoke on the root command runs first, then each intermediate command's PreInvoke, then the target command's PreInvoke.
- **PreRun runs before Traversal**: `preRun()` is called from `complete()` before `traverse()`, so any command structure changes (added subcommands, flags) are visible to traversal.

## Related Skills

- **references/traverse.md** — how storage feeds into the traversal decision tree
- **references/action.md** — the Action API that storage returns