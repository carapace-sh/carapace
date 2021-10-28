# Chdir

[`Chdir`] changes the current working directory to the named directory during invocation.

```go
// completes files for path relative to current wd
carapace.ActionFiles()

// complete files for path relative to `/tmp`
carapace.ActionFiles().Chdir("/tmp")
```

[`Chdir`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.Chdir
