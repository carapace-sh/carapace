module github.com/carapace-sh/carapace/example-multi

go 1.24.0

require (
	github.com/carapace-sh/carapace v1.18.2-0.20261107171937-195323ce78d4
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/carapace-sh/carapace-shlex/v2 v2.1.4 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.1.0 // indirect
	go.yaml.in/yaml/v3 v3.0.6 // indirect
)

replace github.com/carapace-sh/carapace => ../

replace github.com/spf13/pflag => github.com/carapace-sh/carapace-pflag v1.4.0
