# Example

```sh
go install .

# bash
source <(example _carapace bash)

# elvish
eval (example _carapace elvish | slurp)

# fish
example _carapace fish | source

# nushell
example _carapace nushell # update config.nu according to output

# oil
source <(example _carapace oil)

# powershell
Set-PSReadLineOption -Colors @{ "Selection" = "`e[7m" }
Set-PSReadlineKeyHandler -Key Tab -Function MenuComplete
example _carapace powershell | out-string | Invoke-Expression

# tcsh
set autolist
eval `example _carapace tcsh`

# xonsh
$COMPLETION_QUERY_LIMIT = 500 # increase limit
exec($(example _carapace xonsh))

# zsh
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

