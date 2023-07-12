# Shift

[`Shift`] shifts positional arguments left `n` times.

```go
carapace.ActionCallback(func(c carapace.Context) carapace.Action {
	return carapace.ActionMessage("%#v", c.Args)
}).Shift(1)
```

![](./shift.cast)

[`Shift`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.Shift
