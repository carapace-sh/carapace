# ActionMultiPartsN

[`ActionMultiPartsN`] is like [ActionMultiParts] but limits the number of parts to `n`.


```go
carapace.ActionMultiPartsN("=", 2, func(c carapace.Context) carapace.Action {
	switch len(c.Parts) {
	case 0:
		return carapace.ActionValues("one", "two").Suffix("=")
	case 1:
		return carapace.ActionMultiParts("=", func(c carapace.Context) carapace.Action {
			switch len(c.Parts) {
			case 0:
				return carapace.ActionValues("three", "four").Suffix("=")
			case 1:
				return carapace.ActionValues("five", "six")
			default:
				return carapace.ActionValues()
			}
		})
	default:
		return carapace.ActionMessage("should never happen")
	}
})
```

![](./actionMultiPartsN.cast)

[ActionMultiParts]:./actionMultiParts.md
[`ActionMultiPartsN`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.MultipartsN

