# ActionFiles

[`ActionFiles`] completes files with optional suffix filtering.

```go
carapace.ActionFiles(".go")
```

> When used in a static way [`ActionFiles`] uses the shell's own functionality to complete files, but within a [callback action](./actionCallback.md) it does it itself. Thus there might be slight differences during completion.

[`ActionFiles`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionFiles
