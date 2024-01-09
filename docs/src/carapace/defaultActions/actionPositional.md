# ActionPositional

[`ActionPositional`] completes positional arguments for given command ignoring `--` (dash).

```go
carapace.Gen(cmd).DashAnyCompletion(
	carapace.ActionPositional(cmd),
)
```

> It resets `Context.Args` to contain the full arguments and is meant as a means to continue positional completion on dash positions.

![](./actionPositional.cast)

[`ActionPositional`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionPositional
