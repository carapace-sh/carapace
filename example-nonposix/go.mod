module github.com/rsteube/carapace/example-nonposix

go 1.15

require (
	github.com/rsteube/carapace v0.31.1
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
)

replace github.com/rsteube/carapace => ../

replace github.com/spf13/pflag => github.com/rsteube/carapace-pflag v0.2.0
