# FilterParts

[`FilterParts`] filters `Context.Parts`.

```go
carapace.ActionMultiParts(",", func(c carapace.Context) carapace.Action {
	return carapace.ActionValues(
		"one",
		"two",
		"three",
	).FilterParts()
})
```

![](./filterParts.cast)

[`FilterParts`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.FilterParts
