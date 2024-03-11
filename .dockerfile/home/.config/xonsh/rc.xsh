$SHELL="xonsh"
$STARSHIP_SHELL="xonsh"
$LS_COLORS=$(vivid generate dracula)
$PROMPT=lambda: $(starship prompt)

$COMPLETIONS_CONFIRM=True
$COMPLETION_QUERY_LIMIT = 500

'TARGET' in ${...} && exec($($TARGET _carapace xonsh))
