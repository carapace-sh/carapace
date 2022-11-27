# ActionStyledValues

[`ActionStyledValues`] is like [ActionValues](./actionValues.md) but accepts an additional [style](https://pkg.go.dev/github.com/rsteube/carapace/pkg/style).

```go
carapace.ActionStyledValues(
	"first", style.Default,
	"second", style.Blue,
	"third", style.Of(style.BgBrightBlack, style.Magenta, style.Bold),
)
```

![](./actionStyledValues.cast)

[`ActionStyledValues`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionStyledValues
