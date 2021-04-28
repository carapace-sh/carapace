# ActionCallback
[`ActionCallback`](https://pkg.go.dev/github.com/rsteube/carapace#ActionCallback) is a special action where the program itself provides the completion dynamically. For this the [hidden subcommand](../gen/hiddenSubcommand.md) is called with the current command line content which then lets cobra parse existing flags and invokes the callback function after that.

```go
carapace.ActionCallback(func(c carapace.Context) carapace.Action {
  if conditionCmd.Flag("required").Value.String() == "valid" {
    return carapace.ActionValues("condition fulfilled")
  } else {
    return carapace.ActionMessage("flag --required must be set to valid: " + conditionCmd.Flag("required").Value.String())
  }
})
```

- `c.CallbackValue` provides access to the current (partial) value of the flag or positional argument being completed
- return [ActionValues](./actionValues.md) without arguments to silently skip completion
- return [ActionMessage](./actionMessage.md) to provide an error message (e.g. failure during invocation of an external command)
- `c.Args` provides access to the positional arguments of the current subcommand (excluding the one currently being completed)
- [`IsCallback`](https://pkg.go.dev/github.com/rsteube/carapace#IsCallback) indicates if the current invocation of the program is a callback (useful to skip any lengthy init steps)
