Require example

Set Theme "Snazzy"
Set FontFamily "JetBrains Mono"
Set FontSize 30
Set Width 1400
Set Height 300
Set Padding 0

Set Shell "zsh"

Hide
Type@10ms "autoload -U compinit && compinit" Enter
Type@10ms "zstyle ':completion:*' menu select" Enter
Type@10ms "zstyle ':completion:*' format $'\e[2;37mCompleting %d\e[m'" Enter
Type@10ms "source <(example _carapace zsh)" Enter
Type@10ms "clear" Enter
Show
