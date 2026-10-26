Set Shell "pwsh"
Source "../powershell.tape"
Type "example >/dev/null "
Wait+Screen /dev.null/
Sleep 1
Tab@1
Type "ac"
Sleep 1
Tab@1
Sleep 2
Wait+Screen /PS>\ example\ >\/dev\/null\ action/
