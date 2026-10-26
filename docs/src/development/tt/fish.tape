Set Shell "fish"
Set Width 100
Set Height 16
Sleep 3000ms
Hide
Type@10ms "function fish_prompt; printf '$ '; end" Enter
Type@10ms "example _carapace fish | source" Enter
Type@10ms "clear" Enter
WaitStable @250ms
Show
WaitStable @250ms
