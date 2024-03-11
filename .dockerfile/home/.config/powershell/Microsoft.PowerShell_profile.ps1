$env:SHELL = 'powershell'
$env:STARSHIP_SHELL = 'powershell'
$env:LS_COLORS = (&vivid generate dracula)

Invoke-Expression (&starship init powershell)


Set-PSReadLineOption -Colors @{ "Selection" = "`e[7m" }
Set-PSReadlineKeyHandler -Key Tab -Function MenuComplete

if ($Env:TARGET){
	& $Env:TARGET _carapace | out-string | Invoke-Expression
}
