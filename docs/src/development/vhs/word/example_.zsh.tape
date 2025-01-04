Set Theme "Snazzy"
Set FontSize 30
Set Width 1400
Set Height 300
Set Padding 0

Hide
Set Shell "zsh"
Type "autoload -U compinit && compinit" Enter
Type "zstyle ':completion:*' menu select" Enter
Type "zstyle ':completion:*' format $'\e[2;37mCompleting %d\e[m'" Enter
Type "source <(example _carapace)" Enter
Type "clear" Enter
Show

Output     out/example_.zsh.gif
Output     out/example_.zsh.ascii

Type "example "
Sleep 1 Tab@1 2 Enter

Sleep 2
