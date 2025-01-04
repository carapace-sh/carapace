Set Theme "Snazzy"
Set FontSize 32
Set Width 1600
Set Height 400
Set Padding 0
Set TypingSpeed 0

Hide
Set Shell "zsh"
Type "autoload -U compinit && compinit" Enter
Type "zstyle ':completion:*' menu select" Enter
Type "zstyle ':completion:*' format $'\e[2;37mCompleting %d\e[m'" Enter
Type "source <(example _carapace)" Enter
Type "clear" Enter
Show

Type "example a" Tab Sleep 300ms

Screenshot out/example_a.png
Output     out/example_a.gif
Output     out/example_a.ascii
