# Invoke

[`Invoke`] explicitly executes the [callback] of an [Action].

```go
carapace.ActionCallback(func(c carapace.Context) carapace.Action {
	if !strings.HasPrefix(c.Value, "file://") {
		return carapace.ActionValues("file://").NoSpace()
	}

	c.Value = strings.TrimPrefix(c.Value, "file://")
	return carapace.ActionFiles().Invoke(c).Prefix("file://").ToA()
})
```

![](./invoke.cast)

[callback]:../defaultActions/actionCallback.md
[`Invoke`]:https://pkg.go.dev/github.com/rsteube/carapace#Action.Invoke
[Action]:../action.md
