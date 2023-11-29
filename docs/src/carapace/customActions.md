# CustomActions

Custom Actions can be created by using a function that returns `carapace.Action`. A range of these can be found at [carapace-bin](https://pkg.go.dev/github.com/rsteube/carapace-bin/pkg/actions).

```go
type ExampleOpts struct {
	Static bool
}

//  ActionExample(ExampleOpts{Static: true})
func ActionExample(opts ExampleOpts) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if opts.Static {
			return carapace.ActionValues("a", "b")
		}
		if strings.HasPrefix(c.Value, "file://") {
			return carapace.ActionFiles().Invoke(c).Prefix("file://").ToA()
		}
		return carapace.ActionValues()
	})
}
```

> Unless static values are returned the code should be wrapped in a [callback](defaultActions/actionCallback.md) or the code would be executed at program start (and slow it down considerably).
> It is also mandatory when accessing the commands flag values as the callback function is invoked after these are parsed.
