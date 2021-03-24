# carapace

[![CircleCI](https://circleci.com/gh/rsteube/carapace.svg?style=svg)](https://circleci.com/gh/rsteube/carapace)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/rsteube/carapace)](https://pkg.go.dev/github.com/rsteube/carapace)
[![documentation](https://img.shields.io/badge/documentation-grey)](https://rsteube.github.io/carapace/)
[![GoReportCard](https://goreportcard.com/badge/github.com/rsteube/carapace)](https://goreportcard.com/report/github.com/rsteube/carapace)
[![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/rsteube/carapace)](https://hub.docker.com/r/rsteube/carapace)

Completion script generator for [cobra] with support for:

- [Bash](https://www.gnu.org/software/bash/)
- [Elvish](https://elv.sh/)
- [Fish](https://fishshell.com/)
- [Ion](https://doc.redox-os.org/ion-manual/html/) ([experimental](https://github.com/rsteube/carapace/issues/88))
- [Nushell](https://www.nushell.sh/) ([experimental](https://github.com/rsteube/carapace/issues/89))
- [Oil](http://www.oilshell.org/)
- [Powershell](https://microsoft.com/powershell)
- [Xonsh](https://xon.sh/)
- [Zsh](https://www.zsh.org/)


## Status

**WIP**: works, but expect some api changes and small hiccups like a special character not yet escaped

## Usage

Calling `carapace.Gen` on any command is sufficient to enable completion script generation using the [hidden command](https://rsteube.github.io/carapace/carapace/gen/hiddenSubcommand.html).

```go
import (
    "github.com/rsteube/carapace"
)

carapace.Gen(myCmd)
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

# oil
PATH=$PATH:$(pwd)
source <(example _carapace oil)

# powershell
Set-PSReadlineKeyHandler -Key Tab -Function MenuComplete
$env:PATH += ":$pwd"
example _carapace powershell | out-string | Invoke-Expression

# xonsh
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
docker-compose run --rm [bash|elvish|fish|ion|nushell|oil|powershell|xonsh|zsh]

example <TAB>
```

## Projects

- [carapace-bin](https://github.com/rsteube/carapace-bin) multi-shell multi-command argument completer
- [gh](https://github.com/rsteube/gh) github cli with added completions (fork)
- [lab](https://github.com/zaquestion/lab) cli client for GitLab

[cobra]:https://github.com/spf13/cobra
