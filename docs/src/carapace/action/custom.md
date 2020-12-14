# Custom

Custom Actions can be created by using a function that returns `carapace.Action`. A range of these can be found at [carapace-bin](https://pkg.go.dev/github.com/rsteube/carapace-bin/pkg/actions).

```go
func ActionTheme() carapace.Action {
	return carapace.ActionCallback(func(args []string) carapace.Action {
		if output, err := exec.Command("bat", "--list-themes").Output(); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			return carapace.ActionValues(strings.Split(string(output), "\n")...)
		}
	})
}
```

> Unless static values are returned the code should be wrapped in a  [callback](./actionCallback.md) or the code would be executed at program start (and slow it down considerably).
