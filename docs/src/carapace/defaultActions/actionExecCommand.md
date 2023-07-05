# ActionExecCommand

[`ActionExecCommand`] executes an external command.

```go
carapace.ActionExecCommand("git", "remote")(func(output []byte) carapace.Action {
	lines := strings.Split(string(output), "\n")
	return carapace.ActionValues(lines[:len(lines)-1]...)
})
```

![](./actionExecCommand.cast)

[`ActionExecCommand`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionExecCommand
