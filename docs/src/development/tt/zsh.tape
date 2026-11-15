Set Shell "zsh"
Set Width 100
Set Height 16
Sleep 2000ms
Hide
Type@10ms "autoload -U compinit && compinit" Enter
Type@10ms "PROMPT='% '" Enter
Type@10ms "source <(example _carapace zsh)" Enter
Type@10ms "clear" Enter
WaitStable @250ms
Show
WaitStable @250ms
