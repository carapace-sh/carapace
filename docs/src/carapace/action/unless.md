# Unless

[`Unless`] skips invocation if given condition is `true`.

```go
carapace.ActionMultiPartsN(":", 2, func(c carapace.Context) carapace.Action {
	switch len(c.Parts) {
	case 0:
		return carapace.ActionValues("true", "false").Suffix(":")
	default:
		return carapace.Batch(
			carapace.ActionValues(
				"yes",
				"positive",
			).Unless(c.Parts[0] != "true"),
			carapace.ActionValues(
				"no",
				"negative",
			).Unless(c.Parts[0] != "false"),
		).ToA()
	}
})
```

![](./unless.cast)

[`Unless`]:https://pkg.go.dev/github.com/carapace-sh/carapace#Action.Unless
