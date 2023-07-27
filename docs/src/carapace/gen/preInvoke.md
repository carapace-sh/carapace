# PreInvoke

[`PreInvoke`] is called after arguments are parsed and allows generic modification of an [Action] before it is invoked.

```go
carapace.Gen(rootCmd).PreInvoke(func(cmd *cobra.Command, flag *pflag.Flag, action carapace.Action) carapace.Action {
	return action.Chdir(rootCmd.Flag("chdir").Value.String())
})
```

![](./preInvoke.cast)

[Action]:../action.md
[`PreInvoke`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.PreInvoke