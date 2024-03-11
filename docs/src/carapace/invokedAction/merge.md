# Merge

[`Merge`](https://pkg.go.dev/github.com/carapace-sh/carapace#InvokedAction.Merge) combines values of multiple [InvokedActions](../invokedAction.md).

```go
carapace.ActionValues("one", "two").Invoke(c).Merge(carapace.ActionValues("three", "four").Invoke(c)).ToA()
```
