# ActionMultiParts

[`ActionMultiParts`] completes parts of an argument separately (e.g. `user:group` from chown).
For this the `Context.Value` is split with given divider and then updated to only contain the currently completed part.
`Context.Parts` contains the preceding parts and can be used in a `switch` statement to return the corresponding [Action](../action.md).

> An empty divider splits per character, but be aware that fish will add space suffix for anything other than `/=@:.,`.

```go
carapace.ActionMultiParts(":", func(c carapace.Context) carapace.Action {
	switch len(c.Parts) {
	case 0:
		return carapace.ActionValues("userA", "UserB").Invoke(c).Suffix(":").ToA()
	case 1:
		return carapace.ActionValues("groupA", "groupB")
	default:
		return carapace.ActionValues()
	}
})
```

- Values **must not** contain the separator as a simple `strings.Split()` is used to separate the parts.
- It is however **allowed as suffix** to enable fluent tab completion (like `/` for a directory).
- The divider is implicitly added to [`NoSpace`]
- If no suffix is added [`NoSpace`] can be used in the preceding parts to prevent a space suffix.

![](./actionMultiParts.cast)

## Nesting

[`ActionMultiParts`] can be nested as well, e.g. completing multiple `KEY=VALUE` pairs separated by `,`.

```go
carapace.ActionMultiParts(",", func(cEntries carapace.Context) carapace.Action {
	return carapace.ActionMultiParts("=", func(c carapace.Context) carapace.Action {
		switch len(c.Parts) {
		case 0:
			keys := make([]string, len(cEntries.Parts))
			for index, entry := range cEntries.Parts {
				keys[index] = strings.Split(entry, "=")[0]
			}
			return carapace.ActionValues("FILE", "DIRECTORY", "VALUE").Filter(keys...).Suffix("=")
		case 1:
			switch c.Parts[0] {
			case "FILE":
				return carapace.ActionFiles("").NoSpace()
			case "DIRECTORY":
				return carapace.ActionDirectories().NoSpace()
			case "VALUE":
				return carapace.ActionValues("one", "two", "three").NoSpace()
			default:
				return carapace.ActionValues()

			}
		default:
			return carapace.ActionValues()
		}
	})
})
```

![](./actionMultiParts-nested.cast)

[`ActionMultiParts`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionMultiParts
[`NoSpace`]:../action/noSpace.md