# Split

[`Split`] splits `Context.Value` [lexicographically] and replaces `Context.Args` with the tokens.

```go
carapace.ActionCallback(func(c carapace.Context) carapace.Action {
	cmd := &cobra.Command{}
	carapace.Gen(cmd).Standalone()
	cmd.Flags().BoolP("bool", "b", false, "bool flag")
	cmd.Flags().StringP("string", "s", "", "string flag")

	carapace.Gen(cmd).FlagCompletion(carapace.ActionMap{
		"string": carapace.ActionValues("one", "two", "three"),
	})

	carapace.Gen(cmd).PositionalCompletion(
		carapace.ActionValues("pos1", "positional1"),
		carapace.ActionFiles(),
	)

	return carapace.ActionExecute(cmd)
}).Split()
```

![](./split.cast)

[lexicographically]:https://github.com/rsteube/carapace-shlex
[`Split`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.Split
