Require example

Set Theme "Snazzy"
Set FontFamily "JetBrains Mono"
Set FontSize 30
Set Width 1400
Set Height 300
Set Padding 0

Set Shell "pwsh"

Hide
Type@10ms 'Set-PSReadLineOption -Colors @{ "Selection" = "`e[7m" }' Enter
Type@10ms "Set-PSReadlineKeyHandler -Key Tab -Function MenuComplete" Enter
Type@10ms "example _carapace powershell | Out-String | Invoke-Expression" Enter
Type@10ms "clear" Enter
Show
