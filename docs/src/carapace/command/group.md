# Group

[Command Groups] are implicitly used as `tag` for commands.

```go
groupCmd.AddGroup(
	&cobra.Group{ID: "main", Title: "Main Commands"},
	&cobra.Group{ID: "setup", Title: "Setup Commands"},
)

run := func(cmd *cobra.Command, args []string) {}
groupCmd.AddCommand(
	&cobra.Command{Use: "sub1", GroupID: "main", Run: run},
	&cobra.Command{Use: "sub2", GroupID: "main", Run: run},
	&cobra.Command{Use: "sub3", GroupID: "setup", Run: run},
	&cobra.Command{Use: "sub4", GroupID: "setup", Run: run},
	&cobra.Command{Use: "sub5", Run: run},
)
```

![](./group.cast)

[Command Groups]:https://github.com/spf13/cobra/blob/main/site/content/user_guide.md#grouping-commands-in-help