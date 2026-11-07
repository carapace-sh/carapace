module github.com/carapace-sh/carapace

go 1.24.0

require (
	github.com/carapace-sh/carapace-shlex/v2 v2.1.4
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.1.0
	go.yaml.in/yaml/v3 v3.0.6
)

require github.com/inconshreveable/mousetrap v1.1.0 // indirect

replace github.com/spf13/pflag => github.com/carapace-sh/carapace-pflag v1.4.0
