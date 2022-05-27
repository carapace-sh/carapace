# Testing

Since callbacks are simply invocations of the program they can be tested directly.
```sh
example _carapace bash example condition --required ''
valid
invalid

example _carapace elvish example condition --required ''
[{"Value":"valid","Display":"valid"},{"Value":"invalid","Display":"invalid"}]

example _carapace fish example condition --required ''
valid
invalid

example _carapace powershell example condition --required ''
[{"CompletionText":"valid","ListItemText":"valid","ToolTip":" "},{"CompletionText":"invalid","ListItemText":"invalid","ToolTip":" "}]

example _carapace xonsh example condition --required ''
[{"Value":"valid","Display":"valid","Description":""},{"Value":"invalid","Display":"invalid","Description":""}]

example _carapace zsh example condition --required ''
valid   valid
invalid invalid
```
