# Quote

[`Quote`](https://pkg.go.dev/github.com/carapace-sh/carapace#InvokedAction.Quote) quotes values as double-quoted string literals. Only the inserted values are modified; display values are left untouched.

```go
carapace.ActionValues(`say "hi"`).Invoke(c).Quote().ToA()
```
