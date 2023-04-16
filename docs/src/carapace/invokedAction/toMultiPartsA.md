# ToMultiPartsA

[`ToMultiPartsA`] creates an [ActionMultiParts](../defaultActions/actionMultiParts.md) from values containing a specific separator.
E.g. completing the contents of a zip file (`dir/subdir/file`) by each path segment separately like [ActionFiles](../defaultActions/actionFiles.md):

```go
func ActionZipFileContents(file string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if reader, err := zip.OpenReader(file); err != nil {
			return carapace.ActionMessage(err.Error())
		} else {
			defer reader.Close()
			vals := make([]string, len(reader.File))
			for index, f := range reader.File {
				vals[index] = f.Name
			}
			return carapace.ActionValues(vals...).Invoke(c).ToMultiPartsA("/")
		}
	})
}
```

[`ToMultiPartsA`]:https://pkg.go.dev/github.com/rsteube/carapace#InvokedAction.ToMultiPartsA
