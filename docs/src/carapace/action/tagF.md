# TagF

[`TagF`] sets the tag using a function.

```go
carapace.ActionValues(
	"one.png",
	"two.gif",
	"three.txt",
	"four.md",
).StyleF(style.ForPathExt).TagF(func(s string) string {
	switch filepath.Ext(s) {
	case ".png", ".gif":
		return "images"
	case ".txt", ".md":
		return "documents"
	default:
		return ""
	}
})
```

![](./tagF.cast)

[`TagF`]:https://pkg.go.dev/github.com/rsteube/carapace#Action.TagF
