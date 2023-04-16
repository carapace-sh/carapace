# Batch

[`Batch`](https://pkg.go.dev/github.com/rsteube/carapace#Batch) bundles [callback actions](./defaultActions/actionCallback.md) so they can be [invoked](https://pkg.go.dev/github.com/rsteube/carapace#Action.Invoke) in parallel using goroutines.

```go
carapace.ActionCallback(func(c carapace.Context) carapace.Action {
	return carapace.Batch(
		carapace.ActionValues("A", "B"),
		carapace.ActionValues("C", "D"),
		carapace.ActionValues("E", "F"),
	).Invoke(c).Merge().ToA()
})
```
