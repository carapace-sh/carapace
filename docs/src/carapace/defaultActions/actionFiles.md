# ActionFiles

[`ActionFiles`] completes files with optional suffix filtering.

```go
// all files
carapace.ActionFiles()

// files ending with `.md`, `go.mod` or `go.sum`
carapace.ActionFiles(".md", "go.mod", "go.sum"),
```

![](./actionFiles.cast)


[`ActionFiles`]:https://pkg.go.dev/github.com/rsteube/carapace#ActionFiles
