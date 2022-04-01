# ActionStyledValuesDescribed

Same as [ActionValuesDescribed](./actionValuesDescribed.md) but accepts an additional [style](https://pkg.go.dev/github.com/rsteube/carapace/pkg/style).

```go
carapace.ActionStyledValuesDescribed(
  "default", "description of default", style.Default,
  "red", "description of red", style.Red,
  "green-underlined", "description of green-underlined", style.Of(style.Green, style.Underlined),
)
```
