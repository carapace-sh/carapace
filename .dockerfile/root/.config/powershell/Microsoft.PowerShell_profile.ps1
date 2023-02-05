$env:SHELL = 'powershell'
$env:STARSHIP_SHELL = 'powershell'
$env:LS_COLORS = (&vivid generate dracula)
Invoke-Expression (&starship init powershell)
Set-PSReadlineKeyHandler -Key Tab -Function MenuComplete
& $Env:TARGET _carapace | out-string | Invoke-Expression