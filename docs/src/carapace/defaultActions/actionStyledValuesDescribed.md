# ActionStyledValuesDescribed

[`ActionStyledValuesDescribed`] is like [ActionValuesDescribed](./actionValuesDescribed.md) but accepts an additional [style](https://pkg.go.dev/github.com/carapace-sh/carapace/pkg/style).

```go
carapace.ActionStyledValuesDescribed(
	"first", "description of first", style.Blink,
	"second", "description of second", style.Of("color210", style.Underlined),
	"third", "description of third", style.Of("#112233", style.Italic),
)
```

![](./actionStyledValuesDescribed.cast)

[`ActionStyledValuesDescribed`]:https://pkg.go.dev/github.com/carapace-sh/carapace#ActionStyledValuesDescribed