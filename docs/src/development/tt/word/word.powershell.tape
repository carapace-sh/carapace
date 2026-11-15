Set Shell "pwsh"
Source "../powershell.tape"
Type "example "
Wait+Screen /(?m)example/
Sleep 1
Tab@1
Type "ac"
Sleep 1
Tab@1
Sleep 2
Wait+Screen /PS>\ example\ action\ ac/
