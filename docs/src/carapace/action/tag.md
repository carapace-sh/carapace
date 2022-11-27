# Tag

[`Tag`] sets the `tag` for all values. This enables additional grouping and filtering.

```go
ActionValues(
  "192.168.1.1",
  "127.0.0.1",
).Tag("interfaces")
```

## Command Groups

[Command Groups] are implicitly used as `tag` for commands.

```go
func init() {
	carapace.Gen(rootCmd).Standalone()

	rootCmd.AddGroup(
		&cobra.Group{ID: "main", Title: "Main Commands"},
		&cobra.Group{ID: "other", Title: "Other Commands"},
	)
	rootCmd.AddCommand(
		&cobra.Command{Use: "sub1", GroupID: "main", Run: func(cmd *cobra.Command, args []string) {}},
		&cobra.Command{Use: "sub2", GroupID: "main", Run: func(cmd *cobra.Command, args []string) {}},
		&cobra.Command{Use: "sub3", GroupID: "other", Run: func(cmd *cobra.Command, args []string) {}},
		&cobra.Command{Use: "sub4", GroupID: "other", Run: func(cmd *cobra.Command, args []string) {}},
		&cobra.Command{Use: "sub5", Run: func(cmd *cobra.Command, args []string) {}},
	)
}
```

```json
{
  "Version": "v0.26.5",
  "Nospace": "",
  "RawValues": [
    {
      "Value": "sub1",
      "Display": "sub1",
      "Tag": "main commands"
    },
    {
      "Value": "sub2",
      "Display": "sub2",
      "Tag": "main commands"
    },
    {
      "Value": "sub3",
      "Display": "sub3",
      "Tag": "other commands"
    },
    {
      "Value": "sub4",
      "Display": "sub4",
      "Tag": "other commands"
    },
    {
      "Value": "sub5",
      "Display": "sub5",
      "Tag": "additional commands"
    }
  ]
}
```

[Command Groups]:https://github.com/spf13/cobra/blob/main/user_guide.md#grouping-commands-in-help
[`Tag`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.Tag
