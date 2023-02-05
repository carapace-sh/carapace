set SHELL 'fish'
set STARSHIP_SHELL 'fish'
set LS_COLORS (vivid generate dracula)
starship init fish | source 
mkdir -p ~/.config/fish/completions
$TARGET _carapace fish | source