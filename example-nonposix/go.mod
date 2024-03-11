module github.com/carapace-sh/carapace/example-nonposix

go 1.15

require (
	github.com/carapace-sh/carapace v0.50.3-0.20240311124258-a5adf91d8b8f
	github.com/spf13/cobra v1.8.0
	github.com/spf13/pflag v1.0.5
)

replace github.com/carapace-sh/carapace => ../

replace github.com/spf13/pflag => github.com/carapace-sh/carapace-pflag v1.0.0
