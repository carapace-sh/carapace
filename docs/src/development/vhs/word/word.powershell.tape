Set Theme "Snazzy"
Set FontSize 30
Set Width 1400
Set Height 300
Set Padding 0

Hide
Set Shell "pwsh"
Type 'Set-PSReadLineOption -Colors @{ "Selection" = "`e[7m" }' Enter
Type "Set-PSReadlineKeyHandler -Key Tab -Function MenuComplete" Enter
Type "example _carapace powershell | Out-String | Invoke-Expression" Enter
Type "clear" Enter
Show

Output out/word.powershell.gif
Output out/word.powershell.ascii

Type "example " Sleep 1 Tab@1
Type "ac" Sleep 1 Tab@1

Sleep 2
