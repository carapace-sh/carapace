# PreRun

[`PreRun`] is called before arguments are parsed for the current command and allows modification of its structure.

```go
carapace.Gen(rootCmd).PreRun(func(cmd *cobra.Command, args []string) {
	pluginCmd := &cobra.Command{
		Use:     "plugin",
		Short:   "dynamic plugin command",
		GroupID: "plugin",
		Run:     func(cmd *cobra.Command, args []string) {},
	}

	carapace.Gen(pluginCmd).PositionalCompletion(
		carapace.ActionValues("pl1", "pluginArg1"),
	)

	cmd.AddCommand(pluginCmd)
})
```

![](./preRun.cast)

[`PreRun`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.PreRun