# ActionExecCommand

[`ActionExecCommand`] executes given command and transforms its output using the provided function. If an error occurs during execution [ActionMessage](./actionMessage.md) is returned instead with the first line of `stderr` or an exit code.

```go
carapace.ActionExecCommand("git", "remote")(func(output []byte) carapace.Action {
	lines := strings.Split(string(output), "\n")
	return carapace.ActionValues(lines[:len(lines)-1]...)
})
```

![](./actionExecCommand.cast)

[`ActionExecCommand`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionExecCommand
