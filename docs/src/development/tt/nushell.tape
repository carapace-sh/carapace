Set Shell "nu"
Set Width 100
Set Height 16
Sleep 2500ms
Hide
Type@10ms "mkdir ~/.cache" Enter
Type@10ms "example _carapace nushell | save --force ~/.cache/example.nu" Enter
Type@10ms "source ~/.cache/example.nu" Enter
Type@10ms "$env.config.completions.external.completer = $example_completer" Enter
Type@10ms "$env.PROMPT_COMMAND = { '$ ' }" Enter
Type@10ms "$env.PROMPT_COMMAND_RIGHT = { '' }" Enter
Type@10ms "$env.PROMPT_INDICATOR = ''" Enter
Type@10ms "$env.PROMPT_INDICATOR_VI_INSERT = ''" Enter
Type@10ms "clear" Enter
WaitStable @250ms
Show
WaitStable @250ms
