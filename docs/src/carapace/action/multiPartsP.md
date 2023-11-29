# MultiPartsP

[`MultiPartsP`] is like [MultiParts] but with placeholders.

```go
carapace.ActionStyledValuesDescribed(
	"keys/<key>", "key example", style.Default,
	"keys/<key>/<value>", "key/value example", style.Default,
	"styles/custom", "custom style", style.Of(style.Blue, style.Blink),
	"styles", "list", style.Yellow,
	"styles/<style>", "details", style.Default,
).MultiPartsP("/", "<.*>", func(placeholder string, matches map[string]string) carapace.Action {
	switch placeholder {
	case "<key>":
		return carapace.ActionValues("key1", "key2")
	case "<style>":
		return carapace.ActionStyles()
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
