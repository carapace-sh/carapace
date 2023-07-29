# MultiPartsP

[`MultiPartsP`] is like [MultiParts] but with placeholder completion.

```go
carapace.ActionStyledValuesDescribed(
	"keys/<key>/<value>", "key/value example", style.Default,
	"styles/<style>", "details", style.Default,
	"styles/custom", "custom style", style.Of(style.Blue, style.Blink),
	"styles", "list", style.Yellow,
).MultiPartsP("/", "<.*>", func(segment string, matches map[string]string) carapace.Action {
	switch segment {
	case "<style>":
		return carapace.ActionStyles()
	case "<key>":
		return carapace.ActionValues("key1", "key2")
	case "<value>":
		switch matches["<key>"] {
		case "key1":
			return carapace.ActionValues("val1", "val2")
		case "key2":
			return carapace.ActionValues("val3", "val4")
		default:
			return carapace.ActionValues()
		}
	default:
		return carapace.ActionValues()
	}
})
```

![](./multiPartsP.cast)

[MultiParts]:./multiParts.md
[`MultiPartsP`]:https://pkg.go.dev/github.com/rsteube/carapace#Action.MultiPartsP
