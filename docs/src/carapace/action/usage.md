# Usage

[`Usage`] sets the usage message.

```go
carapace.ActionValues().Usage("explicit usage")
````

![](./usage.cast)

> It is implicitly set by default to [`Flag.Usage`] for flag and [`Command.Use`] for positional arguments.

[`Usage`]: https://pkg.go.dev/github.com/rsteube/carapace#Action.Usage
[`Command.Use`]:https://pkg.go.dev/github.com/spf13/cobra#Command
[`Flag.Usage`]:https://pkg.go.dev/github.com/spf13/pflag#Flag