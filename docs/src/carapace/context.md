# Context

[`Context`] provides information during completion.

```go
type Context struct {
	Value string
	Args []string
	Parts []string
	Env []string
	Dir string
}
```

| Key            | Description                                  |
|----------------|----------------------------------------------|
| Value          | current value being completed                | 
| Args           | positional arguments of current (sub)command |
| Parts          | splitted Value during an [ActionMultiParts]  |
| Dir            | working directory                            |


## Examples

Default with flag parsing enabled.
```sh
command pos1 --flag1 pos2 --f<TAB>
# Value: --f
# Args: [pos1, pos2]
```

After encountering `--` (dash) further flag parsing is disabled and `Context.Args` is reset to only contain dash arguments.
```sh
command pos1 --flag1 pos2 -- dash1 <TAB>
# Value:
# Args: [dash1]
```

With [`Command.DisableFlagParsing`] to `true` all arguments are handled as positional.
```sh
command pos1 --flag1 pos2 -- dash1 d<TAB>
# Value: d
# Args: [pos1, --flag1, pos2, --, dash1]
```

With [`SetInterspersed`] to `false` flag parsing is disabled after encountering the first positional argument.
```sh
command --flag1 flagArg1 pos1 -- dash1 --flag2 d<TAB>
# Value: d
# Args: [pos1, --, dash1, --flag2]
```

[ActionMultiParts] is a special case where `Context.Parts` is filled with the splitted `Context.Value`.
```go
ActionValues("part1", "part2", "part3").UniqueList(",")
````

```sh
command pos1 part1,part2,p<TAB>
# Value: p
# Args: [pos1]
# Parts: [part1, part2]
```
 

[ActionMultiParts]:./defaultActions/actionMultiParts.md
[`Command.DisableFlagParsing`]:https://pkg.go.dev/github.com/spf13/cobra#Command
[`Context`]:https://pkg.go.dev/github.com/rsteube/carapace#Context
[`SetInterspersed`]:https://pkg.go.dev/github.com/spf13/pflag#SetInterspersed