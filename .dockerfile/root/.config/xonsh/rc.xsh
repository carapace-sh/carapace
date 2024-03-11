$SHELL="xonsh"
$STARSHIP_SHELL="xonsh"
$LS_COLORS=$(vivid generate dracula)
$PROMPT=lambda: $(starship prompt)
$COMPLETIONS_CONFIRM=True
'TARGET' in ${...} && exec($($TARGET _carapace xonsh))
