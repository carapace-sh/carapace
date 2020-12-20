# Merge

[`Merge`](https://pkg.go.dev/github.com/rsteube/carapace#InvokedAction.Merge) combines values of multiple [InvokedActions](../invokedAction.md).

```go
carapace.ActionValues("one", "two").Invoke(args).Merge(carapace.ActionValues("three", "four").Invoke(args)).ToA()
```
