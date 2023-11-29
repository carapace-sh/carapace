# ActionExecCommandE

[`ActionExecCommandE`] is like [ActionExecCommand] but with custom error handling.

```go
carapace.ActionExecCommandE("false")(func(output []byte, err error) carapace.Action {
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return carapace.ActionMessage("failed with %v", exitErr.ExitCode())
		}
		return carapace.ActionMessage(err.Error())
	}
	return carapace.ActionValues()
})
```

![](./actionExecCommandE.cast)

[ActionExecCommand]:./actionExecCommand.md
[`ActionExecCommandE`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionExecCommandE
