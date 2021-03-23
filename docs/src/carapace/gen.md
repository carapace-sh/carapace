# Gen

Calling [`carapace.Gen`](https://pkg.go.dev/github.com/rsteube/carapace#Gen) on any command is sufficient to enable completion script generation using the [hidden subcommand](./gen/hiddenSubcommand.md).

```go
import (
    "github.com/rsteube/carapace"
)

carapace.Gen(myCmd)
```

> Invocations to `carapace.Gen` must be **after** the command was added to the parent command so that the [uids](./gen/uid.md) are correct.

Additionally invoke [`carapace.Test`](https://pkg.go.dev/github.com/rsteube/carapace#Test) in a [test](https://golang.org/doc/tutorial/add-a-test) to verify configuration during build time.
```go
func TestCarapace(t *testing.T) {
    carapace.Test(t)
}
```

