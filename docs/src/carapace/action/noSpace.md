# NoSpace

[`NoSpace`] disables space suffix for given character(s).

```go
carapace.ActionValues(
	"one,",
	"two/",
	"three",
).NoSpace(',', '/')
```

![](./nospace.cast)

[`NoSpace`]: https://pkg.go.dev/github.com/carapace-sh/carapace#Action.NoSpace
