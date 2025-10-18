# Unique

[`Unique`] ensures the [Action] only contains unique values.

```go
carapace.ActionValues(
    "one",
    "two",
    "two",
    "three",
    "three",
    "three",
).Unique()
```

![](./unique.cast)

[Action]:../action.md
[`Unique`]: https://pkg.go.dev/github.com/carapace-sh/carapace#Action.Unique
