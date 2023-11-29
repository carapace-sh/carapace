# Invoke

[`Invoke`] explicitly executes the [callback] of an [Action].

```go
carapace.ActionCallback(func(c carapace.Context) carapace.Action {
	switch {
	case strings.HasPrefix(c.Value, "file://"):
		c.Value = strings.TrimPrefix(c.Value, "file://")
	case strings.HasPrefix("file://", c.Value):
		c.Value = ""
	default:
		return carapace.ActionValues()
	}
	return carapace.ActionFiles().Invoke(c).Prefix("file://").ToA()
})
```

![](./invoke.cast)

[callback]:../defaultActions/actionCallback.md
[`Invoke`]:https://pkg.go.dev/github.com/rsteube/carapace#Action.Invoke
[Action]:../action.md
