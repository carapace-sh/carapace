# StyleR

[`StyleR`] sets the [style](https://pkg.go.dev/github.com/rsteube/carapace/pkg/style) for all values using a reference.

```go
carapace.ActionValues(
	"one",
	"two",
).StyleR(&style.Carapace.KeywordAmbiguous)
```

![](./styleR.cast)

> Using a reference avoids having to wrap the [Action] in an [ActionCallback] as style configurations are not yet loaded
> when registering the completion.

[Action]:../action.md
[ActionCallback]:../defaultActions/actionCallback.md
[`StyleR`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.StyleR
