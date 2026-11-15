Set Shell "bash"
Set Width 100
Set Height 16
Sleep 1500ms
Hide
Type@10ms "PS1='$ '" Enter
Type@10ms "source <(example _carapace bash)" Enter
Type@10ms "clear" Enter
WaitStable @250ms
Show
WaitStable @250ms
