# FlagCompletion

[`FlagCompletion`] defines completion for flags.

```go
carapace.Gen(myCmd).FlagCompletion(carapace.ActionMap{
    "flagName": carapace.ActionValues("a", "b", "c"),
})
```

## Optional argument

To mark a flag argument as optional (`--name=value`) the [`NoOptDefVal`] needs to be set to anything other than empty string.

```go
rootCmd.Flag("optarg").NoOptDefVal = " "
```

[`FlagCompletion`]:https://pkg.go.dev/github.com/rsteube/carapace#Carapace.FlagCompletion
[`NoOptDefVal`]:https://pkg.go.dev/github.com/spf13/pflag#Flag
