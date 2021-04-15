# Gen

Calling [`Root()`](https://pkg.go.dev/github.com/rsteube/carapace#Carapace.Root) on the root command is sufficient to enable completion script generation using the [hidden subcommand](./gen/hiddenSubcommand.md).

```go
import (
    "github.com/rsteube/carapace"
)

carapace.Gen(myCmd).Root()
```

Additionally invoke [`carapace.Test`](https://pkg.go.dev/github.com/rsteube/carapace#Test) in a [test](https://golang.org/doc/tutorial/add-a-test) to verify configuration during build time.
```go
func TestCarapace(t *testing.T) {
    carapace.Test(t)
}
```

