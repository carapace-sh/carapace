# MultiParts

[`MultiParts`] completes values splitted by given delimiter(s) separately.

```go
carapace.ActionValues(
	"dir/subdir1/fileA.txt",
	"dir/subdir1/fileB.txt",
	"dir/subdir2/fileC.txt",
).MultiParts("/")
```

![](./multiparts.cast)

[`MultiParts`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.Multiparts
