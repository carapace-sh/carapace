# ActionCobra

[`ActionCobra`] bridges given cobra completion function.

```go
carapace.ActionCobra(func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"one", "two"}, cobra.ShellCompDirectiveNoSpace
})
```

![](./actionCobra.cast)

[`ActionCobra`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionCobra
