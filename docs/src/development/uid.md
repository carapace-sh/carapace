# Uid

> deprecated

Uids are generated to identify corresponding completions:

- positional arguments
```handlebars
_{{rootCmd}}__{{subCommand1}}__{{subCommand2}}#{{position}}
```

- flags
```handlebars
_{{rootCmd}}__{{subCommand1}}__{{subCommand2}}##{{flagName}}
```

- state
```handlebars
_{{rootCmd}}__{{subCommand1}}__{{subCommand2}}
```

