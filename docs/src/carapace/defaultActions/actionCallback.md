# ActionCallback

[`ActionCallback`] completes values with given function.
It is invoked after the arguments are parsed which enables contextual completion.

> All [DefaultActions] are implicitly wrapped in an [`ActionCallback`] for performance.

```go
carapace.ActionCallback(func(c carapace.Context) carapace.Action {
	if flag := actionCmd.Flag("values"); flag.Changed {
		return carapace.ActionMessage("values flag is set to: '%v'", flag.Value.String())
	}
	return carapace.ActionMessage("values flag is not set")
})
```

- `c.Value` provides access to the current (partial) value of the flag or positional argument being completed
- return [ActionValues](./actionValues.md) without arguments to silently skip completion
- return [ActionMessage](./actionMessage.md) to provide an error message (e.g. failure during invocation of an external command)
- `c.Args` provides access to the positional arguments of the current subcommand (excluding the one currently being completed)

![](./actionCallback.cast)


[`ActionCallback`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionCallback
[DefaultActions]:../defaultActions.md