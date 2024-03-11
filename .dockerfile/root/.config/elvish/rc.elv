set-env SHELL elvish
set-env STARSHIP_SHELL elvish
set-env LS_COLORS (vivid generate dracula)
set edit:prompt = { starship prompt }
if (has-env TARGET) {
  eval ($E:TARGET _carapace|slurp)
}
