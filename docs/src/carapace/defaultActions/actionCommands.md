# ActionCommands

[`ActionCommands`] completes (sub)commands of given command.

> `Context.Args` is used to traverse the command tree further down.
> Use [Shift](../action/shift.md) to avoid this.


```go
carapace.Gen(helpCmd).PositionalAnyCompletion(
	carapace.ActionCommands(rootCmd),
)
```

![](./actionCommands.cast)

[`ActionCommands`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionCommands
