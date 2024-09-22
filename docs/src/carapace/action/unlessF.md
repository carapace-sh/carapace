# UnlessF

[`UnlessF`] skips invocation if given [condition] returns `true`.

```go
carapace.ActionValues(
	"./local",
	"~/home",
	"/abs",
	"one",
	"two",
	"three",
).UnlessF(condition.CompletingPath)
```

![](./unlessF.cast)

[`UnlessF`]:https://pkg.go.dev/github.com/carapace-sh/carapace#Action.UnlessF
[condition]:https://pkg.go.dev/github.com/carapace-sh/carapace/pkg/condition
