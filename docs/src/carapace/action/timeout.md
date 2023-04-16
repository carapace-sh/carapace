# Timeout

[`Timeout`] sets the maximum duration an [Action] may take to [invoke].

```go
carapace.ActionCallback(func(c carapace.Context) carapace.Action {
	time.Sleep(3*time.Second)
	return carapace.ActionValues("within timeout")
}).Timeout(2*time.Second, carapace.ActionMessage("timeout exceeded"))
```

![](./timeout.cast)

[Action]:../action.md
[invoke]:./invoke.md
[`Timeout`]:https://pkg.go.dev/github.com/rsteube/carapace#Action.Timeout
