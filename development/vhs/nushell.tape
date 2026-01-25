Require example

Set Theme "Snazzy"
Set FontFamily "JetBrains Mono"
Set FontSize 30
Set Width 1400
Set Height 300
Set Padding 0

Set Shell "nu"

Hide
Type@10ms "example _carapace nushell | save --force /tmp/example.nu" Enter
Type@10ms "source /tmp/example.nu" Enter
Type@10ms "$env.config.completions.external.completer = $example_completer" Enter
Type@10ms "clear" Enter
Show
