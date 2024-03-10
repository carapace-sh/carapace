set SHELL 'fish'
set STARSHIP_SHELL 'fish'
set LS_COLORS (vivid generate dracula)
starship init fish | source 
mkdir -p ~/.config/fish/completions
set -q TARGET &&  $TARGET _carapace fish | source
