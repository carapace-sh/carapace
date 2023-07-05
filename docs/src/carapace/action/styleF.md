# StyleF

[`StyleF`] sets the [style](https://pkg.go.dev/github.com/rsteube/carapace/pkg/style) for all values using a function.

```go
carapace.ActionValues(
	"one",
	"two",
	"three",
).StyleF(func(s string, sc style.Context) string {
	switch s {
	case "one":
		return style.Green
	case "two":
		return style.Red
	default:
		return style.Default
	}
})
```

![](./styleF.cast)

[`StyleF`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.StyleF
