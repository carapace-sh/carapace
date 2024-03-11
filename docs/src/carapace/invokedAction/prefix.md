# Prefix

[`Prefix`](https://pkg.go.dev/github.com/carapace-sh/carapace#InvokedAction.Prefix) adds a prefix to all values within an [InvokedAction](../invokedAction.md).

```go
carapace.ActionValues("melon", "drop", "fall").Invoke(c).Prefix("water").ToA()
```
