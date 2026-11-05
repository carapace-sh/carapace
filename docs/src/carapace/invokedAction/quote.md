# Quote

[`Quote`](https://pkg.go.dev/github.com/carapace-sh/carapace#InvokedAction.Quote) escapes values so they can be embedded in a double-quoted string literal. Only the inserted values are modified; display values are left untouched.

```go
carapace.ActionValues(`say "hi"`).Invoke(c).Quote().ToA()
```
