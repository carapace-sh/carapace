Set Shell "pwsh"
Source "../powershell.tape"
Type `example special p`
Wait+Screen /special p/
Sleep 1
Tab@1
Enter
Sleep 2
Wait+Screen />>/
