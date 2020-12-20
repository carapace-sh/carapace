# FlagCompletion

[`FlagCompletion`] defines completion for flags using a map consisting of name and [action](../action.md).

```go
carapace.Gen(myCmd).FlagCompletion(carapace.ActionMap{
    "flagName": carapace.ActionValues("a", "b", "c"),
    // ...
})
```

## Optional argument

To mark a flag argument as optional the [`NoOptDefVal`] needs to be set to anything other than empty string.

```go
rootCmd.Flag("optarg").NoOptDefVal = " "
```

| type | example |
| --- | --- |
| shorthand | `ls -l` |
| longhand | `ls --all` |
| optional argument | `tail --follow=descriptor` |


[`FlagCompletion`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.FlagCompletion
[`NoOptDefVal`]:https://pkg.go.dev/github.com/spf13/pflag#Flag
