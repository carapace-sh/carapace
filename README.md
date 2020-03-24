# carapace

[![CircleCI](https://circleci.com/gh/rsteube/carapace.svg?style=svg)](https://circleci.com/gh/rsteube/carapace)

Completion script generator for [cobra] with support for:

- Bash
- [Fish](https://fishshell.com/docs/current/#writing-your-own-completions)
- [ZSH](https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org)


## Status

**WIP**: works, but expect some api changes and small hiccups like a special character not yet escaped

## Usage

Calling `carapace.Gen` on any command is enough for adding completion script generation using the [hidden command](#hidden-command).

> Invocations to `carapace.Gen` must be **after** the command was to the parent command so that the [uids](#uid) are correct.

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

When `carapace.Gen(myCmd)` is invoked a hidden command (`_carapace`) is added to the root command unless it already exists. This handles completion script generation (when invoked with arg `zsh` or `fish`) and [callbacks](#actioncallback).


### Uid

Uids are generated to identify corresponding completions:
- `_{rootCmd}__{subCommand1}__{subCommand2}#{position}` for positional arguments
- `_{rootCmd}__{subCommand1}__{subCommand2}##{flagName}` for flags


## Action
An [action](#action) indicates how to complete a flag or a positional argument. See [action.go](./action.go) and the examples below for current implementations.

Additional information can be found at [zsh-completions-howto](https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org#functions-for-performing-complex-completions-of-single-words) and the [official documentation](http://zsh.sourceforge.net/Doc/Release/Completion-System.html#Completion-System).

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

./example _carapace zsh '_example__condition#1' example condition --required invalid
# _message -r 'flag --required must be set to valid: invalid'
```

### Custom Action

For [actions](#action) that aren't implemented or missing required options, a custom action can be defined.

```go
carapace.Action{Value: "_most_recent_file 2"}

// #./example action --custom <TAB>
```

## Example

An example implementation can be found in the [example](./example/) folder.

```sh
cd example
go build .
source <(./example _carapace zsh)
./example <TAB>
```

[cobra]:https://github.com/spf13/cobra
