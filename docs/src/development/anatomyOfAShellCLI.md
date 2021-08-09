# Anatomy of a shell CLI

See [Anatomy_of_a_shell_CLI](https://en.wikipedia.org/wiki/Command-line_interface#Anatomy_of_a_shell_CLI) for a detailed description.

## Parameters
Generally a command accepts 0..n parameters:

```sh
command param1 param2 param3 ... paramN
```

Parameters are separated by space which depending on shell can either be escaped (e.g. `arg\ with\ space`) or quoted (e.g. `"arg with space"`).

Everything else is just a commonly used logical classification of the parameters itself as defined by the command.

## Positional Arguments

Any Parameter that is neither a flag, an argument for a flag or a subcommand is a positional argument.
While the order of flags and positional arguments can be mixed, the positional arguments must be after the (sub)command being executed.

```sh
command subCommand1 subCommand2 arg1 --flag1 arg2 --flag2 flagArg1 arg3
```

## Flags

Flags provide a predefined modification of a commands behaviour.
While positional arguments can only be provided once for the (sub)command being executed, flags can be set throughout the command structure.
Unless defined as [persistent](https://github.com/spf13/cobra#persistent-flags) flags need to be set for the (sub)command where they are defined.

```sh
command --flag1 subcommand1 --flag2 --persistentFlag1
```

Posix-style as supported by pflag:

| type                                 | example                    |
| -                                    | -                          |
| longhand flag                        | `command --flag1`          |
| shorthand flag                       | `command -f`               |
| longhand flag with argument          | `command --flag1 flagArg1` |
| shorthand flag with argument         | `command -f flagArg1`      |
| longhand flag with optional argument | `command --flag1=flagArg1` |
| shorthand flag chain                 | `command -abcdef flagArg1` |

Other:

| type                | example          |
| -                   | -                |
| long shorthand flag | `command -flag1` |
| dos/windows style   | `command /flag1` |

## Subcommand

```sh
command subcommand1 --flag1 arg1 --flag2 flagArg1 arg2
command subcommand2 arg1 arg2
```
