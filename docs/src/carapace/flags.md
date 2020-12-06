# Flags

| type | example |
| --- | --- |
| shorthand | `ls -l` |
| longhand | `ls --all` |
| optional argument | `tail --follow=descriptor` |

> - `shorthand-only` flags are supported using the [cornfeedhobo/pflag](https://github.com/cornfeedhobo/pflag) fork


## Completion

Completion for flags can be configured with [`FlagCompletion`](https://pkg.go.dev/github.com/rsteube/carapace#Carapace.FlagCompletion) and a map consisting of name and [Action](./action.md).

```go
carapace.Gen(myCmd).FlagCompletion(carapace.ActionMap{
    "flagName": carapace.ActionValues("a", "b", "c"),
    // ...
})
```

To mark a flag argument as optional the `NoOptDefVal` needs to be set to anything other than empty string.

```go
rootCmd.Flag("optarg").NoOptDefVal = " "
```


