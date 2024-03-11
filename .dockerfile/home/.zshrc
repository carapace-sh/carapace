export SHELL=zsh
export STARSHIP_SHELL=zsh
export LS_COLORS="$(vivid generate dracula)"
eval "$(starship init zsh)"

zstyle ':completion:*' menu select 
zstyle ':completion:*' matcher-list 'm:{a-zA-Z}={A-Za-z}' 'r:|=*' 'l:|=* r:|=*' 

autoload -U compinit && compinit 
source <($TARGET _carapace zsh)