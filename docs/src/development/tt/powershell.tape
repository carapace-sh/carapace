Set Shell "pwsh"
Set Width 100
Set Height 16
Wait+Screen /PowerShell/
Sleep 500ms
Hide
Type@10ms "example _carapace powershell | Out-String | Invoke-Expression" Enter
Sleep 1000ms
Type@10ms "Set-PSReadLineKeyHandler -Key Tab -Function MenuComplete" Enter
Type@10ms "function prompt { 'PS> ' }" Enter
Type@10ms "clear" Enter
WaitStable @250ms
Show
WaitStable @250ms
