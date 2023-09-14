# UniqueListF

[`UniqueListF`] is like [UniqueList] but uses a function to transform values before filtering.

```go
carapace.ActionMultiPartsN(":", 2, func(c carapace.Context) carapace.Action {
	switch len(c.Parts) {
	case 0:
		return carapace.ActionValues("one", "two", "three")
	default:
		return carapace.ActionValues("1", "2", "3")
	}
}).UniqueListF(",", func(s string) string {
	return strings.SplitN(s, ":", 2)[0]
})
```

![](./uniquelistF.cast)

[UniqueList]:./uniqueList.md
[`UniqueListF`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.UniqueListF
