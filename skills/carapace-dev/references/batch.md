# Carapace Library: Batch — Parallel Invocation

Reference for [carapace](https://github.com/carapace-sh/carapace)'s `Batch` API and parallel action invocation model.

## Batch Type

```go
type batch []Action
type invokedBatch []InvokedAction
```

`Batch` wraps a slice of `Action` values for parallel invocation.

## Invocation Model

```go
func (b batch) Invoke(c Context) invokedBatch {
    invokedActions := make([]InvokedAction, len(b))
    functions := make([]func(), len(b))

    for index, action := range b {
        localIndex := index
        localAction := action
        functions[index] = func() {
            invokedActions[localIndex] = localAction.Invoke(c)
        }
    }
    parallelize(functions...)
    return invokedActions
}
```

Each action is captured by index into a separate goroutine closure. `parallelize()` launches all goroutines and waits for completion before returning.

## Parallelize

```go
func parallelize(functions ...func()) {
    var waitGroup sync.WaitGroup
    waitGroup.Add(len(functions))

    defer waitGroup.Wait()

    for _, function := range functions {
        go func(copy func()) {
            defer waitGroup.Done()
            copy()
        }(function)
    }
}
```

Uses `sync.WaitGroup` to coordinate goroutines. Each goroutine captures its function via closure copy (not a reference loop variable).

## Merge

```go
func (b invokedBatch) Merge() InvokedAction {
    switch len(b) {
    case 0:
        return ActionValues().Invoke(Context{})
    case 1:
        return b[0]
    default:
        return b[0].Merge(b[1:]...)
    }
}
```

Deduplicates by `(Value, Uid)` key, concatenates raw values, merges metadata. Single-item batch returns as-is (no extra allocation).

## ToA — Lazy Batch

```go
func (b batch) ToA() Action {
    return ActionCallback(func(c Context) Action {
        return b.Invoke(c).Merge().ToA()
    })
}
```

Wraps the batch as a lazy action. The goroutines only launch when `Invoke()` is called.

## Usage Pattern

```go
carapace.Batch(
    carapace.ActionDirectories(),
    carapace.ActionExecutables(),
).ToA()
```

Equivalent to:
```go
carapace.ActionCallback(func(c carapace.Context) carapace.Action {
    return carapace.Batch(
        carapace.ActionDirectories(),
        carapace.ActionExecutables(),
    ).Invoke(c).Merge().ToA()
})
```

## InvokedAction Merge

`InvokedAction.Merge()` combines multiple invoked actions:

```go
func (ia InvokedAction) Merge(others ...InvokedAction) InvokedAction {
    merged := InvokedAction{
        action: ia.action,
    }
    for _, o := range others {
        merged.action.rawValues = append(merged.action.rawValues, o.action.rawValues...)
        merged.action.meta.Merge(o.action.meta)
    }
    merged.action.rawValues.Unique() // dedup by (Value, Uid)
    return merged
}
```

Called by `invokedBatch.Merge()` to combine parallel results.

## Combined with Traverse

`traverse()` uses `Batch` when a positional argument has subcommands but no positionals consumed:

```go
Batch(getPositional(cmd, 0)) + ActionCommands(cmd)
```

This is the standard case for CLI tools where the first positional is a subcommand name.

## Gotchas

- **No early cutoff**: All goroutines launch regardless of whether earlier ones already produced enough results. For very slow actions, consider `Action.Timeout()` as a guard.
- **Context is shared**: All parallel invocations receive the same `Context`. Mutations to `c.Dir` or `c.Env` by one goroutine are visible to others. Use `Action.Chdir()` on individual actions, not on the batch.
- **Ordering preserved**: `invokedActions[i]` corresponds to `batch[i]`, so merge order is deterministic.
- **PreInvoke in batch**: Actions in a batch have already gone through PreInvoke before being captured. Parallelism applies to the callback execution phase only.

## Related Skills

- **references/action.md** — InvokedAction type and its modifiers
- **references/traverse.md** — where Batch is used in the completion decision