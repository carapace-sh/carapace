# carapace

[![CircleCI](https://circleci.com/gh/rsteube/carapace.svg?style=svg)](https://circleci.com/gh/rsteube/carapace)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/rsteube/carapace)](https://pkg.go.dev/github.com/rsteube/carapace)
[![documentation](https://img.shields.io/badge/&zwnj;-documentation-blue?logo=gitbook)](https://rsteube.github.io/carapace/)
[![GoReportCard](https://goreportcard.com/badge/github.com/rsteube/carapace)](https://goreportcard.com/report/github.com/rsteube/carapace)
[![Coverage Status](https://coveralls.io/repos/github/rsteube/carapace/badge.svg?branch=master)](https://coveralls.io/github/rsteube/carapace?branch=master)

Command argument completion generator for [cobra]. You can read more about it here: _[A pragmatic approach to shell completion](https://dev.to/rsteube/a-pragmatic-approach-to-shell-completion-4gp0)_.


Supported shells:
- [Bash](https://www.gnu.org/software/bash/)
- [Elvish](https://elv.sh/)
- [Fish](https://fishshell.com/)
- [Ion](https://doc.redox-os.org/ion-manual/) ([experimental](https://github.com/rsteube/carapace/issues/88))
- [Nushell](https://www.nushell.sh/)
- [Oil](http://www.oilshell.org/)
- [Powershell](https://microsoft.com/powershell)
- [Tcsh](https://www.tcsh.org/) ([experimental](https://github.com/rsteube/carapace/issues/331))
- [Xonsh](https://xon.sh/)
- [Zsh](https://www.zsh.org/)

## Usage

Calling `carapace.Gen` on the root command is sufficient to enable completion script generation using the [hidden command](https://rsteube.github.io/carapace/carapace/gen/hiddenSubcommand.html).

```go
import (
    "github.com/rsteube/carapace"
)

carapace.Gen(rootCmd)
```

## Standalone Mode

Carapace can also be used to provide completion for arbitrary commands as well (similar to [aws_completer](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-completion.html)).
See [rsteube/carapace-bin](https://github.com/rsteube/carapace-bin) for examples. There is also a binary to parse flags from gnu help pages at [caraparse](https://github.com/rsteube/carapace-bin/tree/master/cmd/caraparse).

## Example

An example implementation can be found in the [example](./example/) folder.

```sh
cd example
go build .

# bash
PATH=$PATH:$(pwd)
source <(example _carapace bash)

# elvish
paths=[$@paths (pwd)]
eval (example _carapace elvish | slurp)

# fish
set PATH $PATH (pwd) 
example _carapace fish | source

# nushell
example _carapace nushell

# oil
PATH=$PATH:$(pwd)
source <(example _carapace oil)

# powershell
Set-PSReadLineOption -Colors @{ "Selection" = "`e[7m" }
Set-PSReadlineKeyHandler -Key Tab -Function MenuComplete
$env:PATH += ":$pwd"
example _carapace powershell | out-string | Invoke-Expression

# tcsh
set autolist
eval `example _carapace tcsh`

# xonsh
$COMPLETION_QUERY_LIMIT = 500 # increase limit
$PATH.append($(pwd))
exec($(example _carapace xonsh))

# zsh
PATH=$PATH:$(pwd)
source <(example _carapace zsh)

example <TAB>
```

or use [docker-compose](https://docs.docker.com/compose/):
```sh
docker-compose pull
docker-compose run --rm build
docker-compose run --rm [bash|elvish|fish|ion|nushell|oil|powershell|tcsh|xonsh|zsh]

example <TAB>
```

## Projects

- [carapace-bin](https://github.com/rsteube/carapace-bin) multi-shell multi-command argument completer
- [carapace-spec](https://github.com/rsteube/carapace-spec) define simple completions using a spec file
- [freckles](https://github.com/rsteube/freckles) simple dotfiles manager
- [go-jira-cli](https://github.com/rsteube/go-jira-cli) simple jira command line client
- [knoxite](https://github.com/knoxite/knoxite) A data storage & backup system
- [lab](https://github.com/zaquestion/lab) cli client for GitLab

[cobra]:https://github.com/spf13/cobra
