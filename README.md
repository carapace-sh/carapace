# cobra-zsh-gen

[![Build Status](https://travis-ci.org/rsteube/cobra-zsh-gen.svg?branch=master)](https://travis-ci.org/rsteube/cobra-zsh-gen)
[![CircleCI](https://circleci.com/gh/rsteube/cobra-zsh-gen.svg?style=svg)](https://circleci.com/gh/rsteube/cobra-zsh-gen)

[ZSH completions](https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org) script generator for [cobra] (based on [spf13/cobra#646](https://github.com/spf13/cobra/pull/646)).


## Status

**WIP**: works, but expect some api changes and small hiccups like a special character not yet escaped

## Usage

Calling `zsh.Gen` on any command is enough for adding completion script generation using the [hidden command](#hidden-command).

> Invocations to `zsh.Gen` must be **after** the command was to the parent command so that the [uids](#uid) are correct.

```go
import (
    "github.com/rsteube/cobra-zsh-gen"
)

zsh.Gen(myCmd)
```

### FlagCompletions

Completion for flags can be configured with `FlagCompletion` and a map consisting of name and [action](#action).

```go
zsh.Gen(myCmd).FlagCompletion(zsh.ActionMap{
    "flagName": zsh.ActionValues("a", "b", "c"),
    ...
})
```

### PositionalCompletions

Completion for positional arguments can be configured with `PositionalCompletion` and a list of [actions](#action).

```go
zsh.Gen(callbackCmd).PositionalCompletion(
    zsh.ActionValues("a", "b", "c"),
    ...
)
```

## Hidden Command

When `zsh.Gen(myCmd)` is invoked a hidden command (`_zsh_completion`) is added to the root command unless it already exists. This handles zsh completion script generation (when invoked without args: `./executable _zsh_completion`) and [callbacks](#actioncallback).


### Uid

Uids are generated to identify corresponding completions:
- `_{rootCmd}__{subCommand1}__{subCommand2}#{position}` for positional arguments
- `_{rootCmd}__{subCommand1}__{subCommand2}##{flagName}` for flags


## Action
An [action](#action) indicates how to complete a flag or a positional argument. See [action.go](./action.go) and the examples below for current implementations.

Additional information can be found at [zsh-completions-howto](https://github.com/zsh-users/zsh-completions/blob/master/zsh-completions-howto.org#functions-for-performing-complex-completions-of-single-words) and the [official documentation](http://zsh.sourceforge.net/Doc/Release/Completion-System.html#Completion-System).

### ActionMessage

```go
zsh.ActionMessage("message example")

// #./example action --message <TAB>
// message example
```

### ActionValuesDescribed

```go
zsh.ActionValuesDescribed(
  "values", "valueDescription",
  "example", "exampleDescription"),

// #./example action --values_described <TAB>
// example  -- exampleDescription
// values   -- valueDescription
```

### ActionCallback

ActionCallback is a special [action](#action) where the program itself provides the completion dynamically. For this the [hidden command](#hidden-command) is called with a [uid](#uid) and the current command line content which then lets [cobra] parse existing flags and invokes the callback function after that.

```go
zsh.ActionCallback(func(args []string) zsh.Action {
  if conditionCmd.Flag("required").Value.String() == "valid" {
    return zsh.ActionValues("condition fulfilled")
  } else {
    return zsh.ActionMessage("flag --required must be set to valid: " + conditionCmd.Flag("required").Value.String())
  }
})

// #./example condition --required invalid <TAB>
// flag --required must be set to valid: invalid
```

Since callbacks are simply invocations of the program they can be tested directly.

```sh
./example _zsh_completion '_example__condition#1' condition --required invalid

# _message -r 'flag --required must be set to valid: invalid'
```

### Custom Action

For [actions](#action) that aren't implemented or missing required options, a custom action can be defined.

```go
zsh.Action{Value: "_most_recent_file 2"}

// #./example action --custom <TAB>
```

## Example

An example implementation can be found in the [example](./example/) folder.

```sh
cd example
go build .
source <(./example _zsh_completion)
./example <TAB>
```

[cobra]:https://github.com/spf13/cobra
