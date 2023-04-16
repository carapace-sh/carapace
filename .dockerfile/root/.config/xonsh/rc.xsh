$SHELL="xonsh"
$STARSHIP_SHELL="xonsh"
$LS_COLORS=$(vivid generate dracula)
$PROMPT=lambda: $(starship prompt)
$COMPLETIONS_CONFIRM=True
exec($($TARGET _carapace xonsh))