module github.com/carapace-sh/carapace/example-nonposix

go 1.24

require (
	github.com/carapace-sh/carapace v0.50.3-0.20240311124258-a5adf91d8b8f
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.9
)

require (
	github.com/carapace-sh/carapace-shlex v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/carapace-sh/carapace => ../

replace github.com/spf13/pflag => github.com/carapace-sh/carapace-pflag v1.1.0
