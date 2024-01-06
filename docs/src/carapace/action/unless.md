# Unless

[`Unless`] skips invokation if given [condition] succeeds.

```go
carapace.ActionValues(
	"./local",
	"~/home",
	"/abs",
	"one",
	"two",
	"three",
).Unless(condition.CompletingPath)
```

![](./unless.cast)

[`Unless`]:https://pkg.go.dev/github.com/rsteube/carapace#Action.Unless
[condition]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/condition
