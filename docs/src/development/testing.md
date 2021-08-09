# Testing

Since callbacks are simply invocations of the program they can be tested directly.
```sh
example _carapace bash _ example condition --required ''
valid
invalid

example _carapace elvish _ example condition --required ''
[{"Value":"valid","Display":"valid"},{"Value":"invalid","Display":"invalid"}]

example _carapace fish _ example condition --required ''
valid
invalid

example _carapace powershell _ example condition --required ''
[{"CompletionText":"valid","ListItemText":"valid","ToolTip":" "},{"CompletionText":"invalid","ListItemText":"invalid","ToolTip":" "}]

example _carapace xonsh _ example condition --required ''
[{"Value":"valid","Display":"valid","Description":""},{"Value":"invalid","Display":"invalid","Description":""}]

example _carapace zsh _ example condition --required ''
valid   valid
invalid invalid
```
