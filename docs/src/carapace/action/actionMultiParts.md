# ActionMultiParts

[`ActionMultiParts`] is a [callback action](./actionCallback.md) where parts of an argument can be completed separately (e.g. `user:group` from chown). Divider can be empty as well, but note that bash and fish will add the space suffix for anything other than `/=@:.,` (it still works, but after each selection backspace is needed to continue the completion).

```go
carapace.ActionMultiParts(":", func(args []string, parts []string) carapace.Action {
	switch len(parts) {
	case 0:
		return ActionUsers().Invoke(args).Suffix(":").ToA()
	case 1:
		return ActionGroups()
	default:
		return carapace.ActionValues()
	}
})
```

- [`carapace.CallbackValue`] is updated during execution to only contain the currently completed part
- values **must not** contain the separator as a simple `strings.Split()` is used to separate the parts
- it is however **allowed as suffix** to enable fluent tab completion (like `/` for a directory)

> There are still some [issues](https://github.com/rsteube/carapace/issues?q=is%3Aissue+is%3Aopen+ActionMultiParts+) with this so a couple of edge cases might not work

## Nesting

[`ActionMultiParts`] can be nested as well, e.g. completing multiple `KEY=VALUE` pairs separated by `,`.

```go
carapace.ActionMultiParts(",", func(args, entries []string) carapace.Action {
	return carapace.ActionMultiParts("=", func(args, parts []string) carapace.Action {
		switch len(parts) {
		case 0:
			keys := make([]string, len(entries))
			for index, entry := range entries {
				keys[index] = strings.Split(entry, "=")[0]
			}
			return carapace.ActionValues("FILE", "DIRECTORY", "VALUE").Invoke(args).Filter(keys).Suffix("=").ToA()
		case 1:
			switch parts[0] {
			case "FILE":
				return carapace.ActionFiles("")
			case "DIRECTORY":
				return carapace.ActionDirectories()
			case "VALUE":
				return carapace.ActionValues("one", "two", "three")
			default:
				return carapace.ActionValues()

			}
		default:
			return carapace.ActionValues()
		}
	})
}),
```

[`carapace.CallbackValue`]:https://pkg.go.dev/github.com/rsteube/carapace#pkg-variables
[`ActionMultiParts`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionMultiParts