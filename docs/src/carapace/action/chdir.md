# Chdir

[`Chdir`] changes the working directory.

```go
// completes files for path relative to current working directory
carapace.ActionFiles()

// complete files for path relative to `/tmp`
carapace.ActionFiles().Chdir("/tmp")
```

![](./chdir.cast)

[`Chdir`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.Chdir
