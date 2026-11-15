Set Shell "xonsh"
Set Width 100
Set Height 16
Sleep 3000ms
Hide
Type@10ms "$PROMPT = '$ '" Enter
Type@10ms "exec($(example _carapace xonsh))" Enter
Type@10ms "clear" Enter
WaitStable @250ms
Show
WaitStable @250ms
