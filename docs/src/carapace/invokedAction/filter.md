# Filter

[`Filter`](https://pkg.go.dev/github.com/rsteube/carapace#InvokedAction.Filter) filters values within an [InvokedAction](../invokedAction.md).
E.g. completing a unique list of values in an [ActionMultiParts](../defaultActions/actionMultiParts.md):

```go
carapace.ActionMultiParts(",", func(c carapace.Context) carapace.Action {
 	return carapace.ActionValues("one", "two", "three").Invoke(c).Filter(c.Parts...).ToA()
}
```
