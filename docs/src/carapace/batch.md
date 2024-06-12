# Batch

[`Batch`](https://pkg.go.dev/github.com/carapace-sh/carapace#Batch) bundles [callback actions](./defaultActions/actionCallback.md) so they can be [invoked] concurrently using goroutines.

```go
carapace.ActionCallback(func(c carapace.Context) carapace.Action {
	return carapace.Batch(
		carapace.ActionValues("A", "B"),
		carapace.ActionValues("C", "D"),
		carapace.ActionValues("E", "F"),
	).Invoke(c).Merge().ToA()
})
```

[invoked]:https://pkg.go.dev/github.com/carapace-sh/carapace#Action.Invoke
