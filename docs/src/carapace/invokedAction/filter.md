# Filter

[`Filter`](https://pkg.go.dev/github.com/rsteube/carapace#InvokedAction.Filter) filters values within an [InvokedAction](../invokedAction.md).
E.g. completing a unique list of values in an [ActionMultiParts](../action/actionMultiParts.md):

```go
carapace.ActionMultiParts(",", func(args, parts []string) carapace.Action {
 	return carapace.ActionValues("one", "two", "three").Invoke(args).Filter(parts).ToA()
}
```
