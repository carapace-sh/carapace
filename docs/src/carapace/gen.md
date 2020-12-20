# Gen

Calling [`carapace.Gen`](https://pkg.go.dev/github.com/rsteube/carapace#Gen) on any command is sufficient to enable completion script generation using the [hidden subcommand](./gen/hiddenSubcommand.md).

```go
import (
    "github.com/rsteube/carapace"
)

carapace.Gen(myCmd)
```

> Invocations to `carapace.Gen` must be **after** the command was added to the parent command so that the [uids](./gen/uid.md) are correct.
