Set Shell "elvish"
Set Width 100
Set Height 16
Sleep 1500ms
Hide
Type@10ms "eval (example _carapace elvish | slurp)" Enter
Type@10ms "set edit:prompt = { put '$ ' }" Enter
Type@10ms "set edit:rprompt = { put '' }" Enter
Type@10ms "edit:clear" Enter
WaitStable @250ms
Show
WaitStable @250ms
