Require example

Set Theme "Snazzy"
Set FontFamily "JetBrains Mono"
Set FontSize 30
Set Width 1400
Set Height 300
Set Padding 0

Set Shell "zsh"

Hide
Type@0.01 "autoload -U compinit && compinit" Enter
Type@0.01 "zstyle ':completion:*' menu select" Enter
Type@0.01 "zstyle ':completion:*' format $'\e[2;37mCompleting %d\e[m'" Enter
Type@0.01 "source <(example _carapace zsh)" Enter
Type@0.01 "clear" Enter
Show
