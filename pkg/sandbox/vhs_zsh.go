package sandbox

func (v vhs) zsh() string {
	return `Hide
Set Shell "zsh"
Type "autoload -U compinit && compinit" Enter
Type "zstyle ':completion:*' menu select" Enter
Type "zstyle ':completion:*' format $'\e[2;37mCompleting %d\e[m'" Enter
Type "source <(example _carapace)" Enter
Type "clear" Enter
Show
`
}

// Type "example a" Tab Sleep 300ms

func (v vhs) output() string {
	return `Screenshot out/example_a.png
Output     out/example_a.gif
Output     out/example_a.ascii
`
}
