# Suppress

[`Suppress`] suppresses specific error messages using regular expressions.

```go
carapace.Batch(
	carapace.ActionMessage("unexpected error"),
	carapace.ActionMessage("ignored error"),
).ToA().Suppress("ignored")
```

![](./suppress.cast)

[`Suppress`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.Suppress
