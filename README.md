# carapace

[![CircleCI](https://circleci.com/gh/rsteube/carapace.svg?style=svg)](https://circleci.com/gh/rsteube/carapace)

Completion script generator for [cobra] with support for:

- [Bash](https://www.gnu.org/software/bash/manual/html_node/A-Programmable-Completion-Example.html)
- [Elvish](https://elv.sh/ref/edit.html#editcompletionarg-completer)
- [Fish](https://fishshell.com/docs/current/#writing-your-own-completions)
- [Oil](http://www.oilshell.org/blog/2018/10/10.html) *broken* ([#86](https://github.com/rsteube/carapace/issues/86))
- [Powershell](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/register-argumentcompleter)
- [Xonsh](https://xon.sh/tutorial_completers.html#writing-a-new-completer) *experimental*
- [Zsh](https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org)


## Status

**WIP**: works, but expect some api changes and small hiccups like a special character not yet escaped

## Usage

Calling `carapace.Gen` on any command is enough to enable completion script generation using the [hidden command](#hidden-command).

> Invocations to `carapace.Gen` must be **after** the command was added to the parent command so that the [uids](#uid) are correct.

```go
import (
    "github.com/rsteube/carapace"
)

carapace.Gen(myCmd)
```

### FlagCompletions

Completion for flags can be configured with `FlagCompletion` and a map consisting of name and [action](#action).

```go
carapace.Gen(myCmd).FlagCompletion(carapace.ActionMap{
    "flagName": carapace.ActionValues("a", "b", "c"),
    ...
})
```

### PositionalCompletions

Completion for positional arguments can be configured with `PositionalCompletion` and a list of [actions](#action).

```go
carapace.Gen(callbackCmd).PositionalCompletion(
    carapace.ActionValues("a", "b", "c"),
    ...
)
```

## Hidden Command

When `carapace.Gen(myCmd)` is invoked a hidden command (`_carapace`) is added to the root command unless it already exists. This handles completion script generation and [callbacks](#actioncallback).


### Uid

Uids are generated to identify corresponding completions:
- `_{rootCmd}__{subCommand1}__{subCommand2}#{position}` for positional arguments
- `_{rootCmd}__{subCommand1}__{subCommand2}##{flagName}` for flags


## Action
An [action](#action) indicates how to complete a flag or a positional argument. See [action.go](./action.go) and the examples below for current implementations.

### ActionMessage

```go
carapace.ActionMessage("message example")

// #./example action --message <TAB>
// message example
```

### ActionValuesDescribed

```go
carapace.ActionValuesDescribed(
  "values", "valueDescription",
  "example", "exampleDescription"),

// #./example action --values_described <TAB>
// example  -- exampleDescription
// values   -- valueDescription
```

### ActionCallback

ActionCallback is a special [action](#action) where the program itself provides the completion dynamically. For this the [hidden command](#hidden-command) is called with a [uid](#uid) and the current command line content which then lets [cobra] parse existing flags and invokes the callback function after that.

```go
carapace.ActionCallback(func(args []string) carapace.Action {
  if conditionCmd.Flag("required").Value.String() == "valid" {
    return carapace.ActionValues("condition fulfilled")
  } else {
    return carapace.ActionMessage("flag --required must be set to valid: " + conditionCmd.Flag("required").Value.String())
  }
})

// #./example condition --required invalid <TAB>
// flag --required must be set to valid: invalid
```

Since callbacks are simply invocations of the program they can be tested directly.

```sh
./example _carapace bash '_example__condition#1' example condition --required invalid
#compgen -W "ERR flag_--required_must_be_set_to_valid:_invalid" -- $last

./example _carapace elvish '_example__condition#1' example condition --required invalid
#edit:complex-candidate ERR &display-suffix=' (flag --required must be set to valid: invalid)'
#edit:complex-candidate _ &display-suffix=' ()'

./example _carapace fish '_example__condition#1' example condition --required invalid
#echo -e ERR\tflag --required must be set to valid: invalid\n_\t\n\n

./example _carapace powershell '_example__condition#1' example condition --required invalid
#[CompletionResult]::new('ERR', 'ERR', [CompletionResultType]::ParameterValue, ' ')
#[CompletionResult]::new('flag --required must be set to valid: invalid', 'flag --required must be set to valid: invalid', [CompletionResultType]::ParameterValue, ' ')

./example _carapace xonsh '_example__condition#1' example condition --required invalid
#{
#  RichCompletion('_', display='_', description='flag --required must be set to valid: invalid', prefix_len=0),
#  RichCompletion('ERR', display='ERR', description='flag --required must be set to valid: invalid', prefix_len=0),
#}

./example _carapace zsh '_example__condition#1' example condition --required invalid
# _message -r 'flag --required must be set to valid: invalid'
```

### ActionMultiParts

> This is an initial version which still got some quirks, expect some changes here (in the long term this shall return Action as well)

ActionMultiParts is a [callback action](#actioncallback) where parts of an argument can be completed separately (e.g. user:group from [chown](https://github.com/rsteube/carapace-completers/blob/master/completers/chown_completer/cmd/root.go)). Divider can be empty as well, but note that `bash` and `fish` will add the space suffix for anything other than `/=@:.,` (it still works, but after each selection backspace is needed to continue the completion).

```go
carapace.ActionMultiParts(":", func(args []string, parts []string) []string {
	switch len(parts) {
	case 0:
		return []{"user1:", "user2:", "user3:"}
	case 1:
		return []{"groupA", "groupB", "groupC"}
	default:
		return []string{}
	}
})
```

### Custom Action

For [actions](#action) that aren't implemented or missing required options, a custom action can be defined.

```go
carapace.Action{Zsh: "_most_recent_file 2"}

// #./example action --custom <TAB>
```

Additional information can be found at:
- Bash: [bash-programmable-completion-tutorial](https://iridakos.com/programming/2018/03/01/bash-programmable-completion-tutorial) and [Programmable-Completion-Builtins](https://www.gnu.org/software/bash/manual/html_node/Programmable-Completion-Builtins.html#Programmable-Completion-Builtins)
- Elvish: [using-and-writing-completions-in-elvish](https://zzamboni.org/post/using-and-writing-completions-in-elvish/) and [argument-completer](https://elv.sh/ref/edit.html#argument-completer)
- Fish: [fish-shell/share/functions](https://github.com/fish-shell/fish-shell/tree/master/share/functions) and [writing your own completions](https://fishshell.com/docs/current/#writing-your-own-completions)
- Powershell: [Dynamic Tab Completion](https://adamtheautomator.com/powershell-parameters-argumentcompleter/) and [Register-ArgumentCompleter](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/register-argumentcompleter)
- Xonsh: [Programmable Tab-Completion](https://xon.sh/tutorial_completers.html) and [RichCompletion(str)](https://github.com/xonsh/xonsh/blob/master/xonsh/completers/tools.py)
- Zsh: [zsh-completions-howto](https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org#functions-for-performing-complex-completions-of-single-words) and [Completion-System](http://zsh.sourceforge.net/Doc/Release/Completion-System.html#Completion-System).

## Standalone Mode

Carapace can also be used to provide completion for arbitrary commands as well (similar to [aws_completer](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-completion.html)).
See [rsteube/carapace-bin](https://github.com/rsteube/carapace-bin) for examples. There is also a binary to parse flags from gnu help pages at [caraparse](https://github.com/rsteube/carapace-bin/tree/master/caraparse).

## Example

An example implementation can be found in the [example](./example/) folder.

```sh
cd example
go build .

# bash
PATH=$PATH:$(pwd)
source <(example _carapace bash)

# elvish (-source will be replaced in next version: `eval (example _carapace elvish | slurp`)
)
paths=[$@paths (pwd)]
example _carapace elvish > example.elv
-source example.elv

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
docker-compose run --rm build
docker-compose run --rm [bash|elvish|fish|oil|powershell|xonsh|zsh]

example <TAB>
```

## Projects

- [carapace-bin](https://github.com/rsteube/carapace-bin) multi-shell multi-command argument completer
- [lab](https://github.com/zaquestion/lab) cli client for GitLab

[cobra]:https://github.com/spf13/cobra
