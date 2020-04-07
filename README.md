# carapace

[![CircleCI](https://circleci.com/gh/rsteube/carapace.svg?style=svg)](https://circleci.com/gh/rsteube/carapace)

Completion script generator for [cobra] with support for:

- [Bash](https://www.gnu.org/software/bash/manual/html_node/A-Programmable-Completion-Example.html)
- [Fish](https://fishshell.com/docs/current/#writing-your-own-completions)
- [Powershell](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/register-argumentcompleter) _(in progress)_
- [Zsh](https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org)


## Status

**WIP**: works, but expect some api changes and small hiccups like a special character not yet escaped

## Usage

Calling `carapace.Gen` on any command is enough for adding completion script generation using the [hidden command](#hidden-command).

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

./example _carapace fish '_example__condition#1' example condition --required invalid
#echo -e ERR\tflag --required must be set to valid: invalid\n_\t\n\n

./example _carapace powershell '_example__condition#1' example condition --required invalid
#[CompletionResult]::new('ERR', 'ERR', [CompletionResultType]::ParameterValue, ' ')
#[CompletionResult]::new('flag --required must be set to valid: invalid', 'flag --required must be set to valid: invalid', [CompletionResultType]::ParameterValue, ' ')

./example _carapace zsh '_example__condition#1' example condition --required invalid
# _message -r 'flag --required must be set to valid: invalid'
```

### Custom Action

For [actions](#action) that aren't implemented or missing required options, a custom action can be defined.

```go
carapace.Action{Zsh: "_most_recent_file 2"}

// #./example action --custom <TAB>
```

Additional information can be found at:
- Bash: [bash-programmable-completion-tutorial](https://iridakos.com/programming/2018/03/01/bash-programmable-completion-tutorial) and [Programmable-Completion-Builtins](https://www.gnu.org/software/bash/manual/html_node/Programmable-Completion-Builtins.html#Programmable-Completion-Builtins)
- Fish: [fish-shell/share/functions](https://github.com/fish-shell/fish-shell/tree/master/share/functions) and [writing your own completions](https://fishshell.com/docs/current/#writing-your-own-completions)
- Powershell: [Dynamic Tab Completion](https://adamtheautomator.com/powershell-parameters-argumentcompleter/) and [Register-ArgumentCompleter](https://docs.microsoft.com/en-us/powershell/module/microsoft.powershell.core/register-argumentcompleter)
- Zsh: [zsh-completions-howto](https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org#functions-for-performing-complex-completions-of-single-words) and [Completion-System](http://zsh.sourceforge.net/Doc/Release/Completion-System.html#Completion-System).


## Example

An example implementation can be found in the [example](./example/) folder.

```sh
cd example
go build .

# bash
PATH=$PATH:$(pwd)
source <(example _carapace bash)

# fish
set PATH $PATH (pwd) 
example _carapace fish | source

# powershell
Set-PSReadlineKeyHandler -Key Tab -Function MenuComplete
$env:PATH += ":$pwd"
example _carapace powershell | out-string | Invoke-Expression

# zsh
PATH=$PATH:$(pwd)
source <(example _carapace zsh)

example <TAB>
```

or use the preconfigured docker containers (on linux):
```sh
cd example
go build .
docker-compose run --rm [bash|fish|powershell|zsh]

example <TAB>
```

[cobra]:https://github.com/spf13/cobra
