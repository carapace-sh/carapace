# Usage

[`Usage`] sets the `usage` message.

```go
carapace.ActionMultiParts(":", func(c carapace.Context) carapace.Action {
	switch len(c.Parts) {
	case 0:
		return carapace.ActionValues("explicit", "implicit").Suffix(":")
	case 1:
		if c.Parts[0] == "explicit" {
			return carapace.ActionValues().Usage("explicit usage")
		}
		return carapace.ActionValues()

	default:
		return carapace.ActionValues()
	}
})
````
![](./usage.cast)

> It is implicitly set by default to [`Flag.Usage`] for flag and [`Command.Use`] for positional arguments.

[`Usage`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.Usage
[`Command.Use`]:https://pkg.go.dev/github.com/spf13/cobra#Command
[`Flag.Usage`]:https://pkg.go.dev/github.com/spf13/pflag#Flag