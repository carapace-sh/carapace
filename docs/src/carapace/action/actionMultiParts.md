# ActionMultiParts

[ActionMultiParts](https://pkg.go.dev/github.com/rsteube/carapace#ActionMultiParts) is a [callback action](./actionCallback.md) where parts of an argument can be completed separately (e.g. user:group from chown). Divider can be empty as well, but note that bash and fish will add the space suffix for anything other than `/=@:.,` (it still works, but after each selection backspace is needed to continue the completion).

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
