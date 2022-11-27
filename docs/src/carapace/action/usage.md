# Usage

[`Usage`] sets the `usage` message.

> `usage` is implicitly set by default with [`Flag.Usage`] for flag and [`Command.Use`] for positional arguments.

```go
carapace.ActionValues().Usage("explicit usage")
````

```json
{
  "Version": "unknown",
  "Usage": "explicit usage",
  "Nospace": "",
  "RawValues": []
}
```

[`Usage`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.Usage
[`Command.Use`]:https://pkg.go.dev/github.com/spf13/cobra#Command
[`Flag.Usage`]:https://pkg.go.dev/github.com/spf13/pflag#Flag