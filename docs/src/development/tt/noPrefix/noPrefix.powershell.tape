Set Shell "pwsh"
Source "../powershell.tape"
Type "example modifier --noprefix one"
Wait+Screen /--noprefix one/
Sleep 1 Tab@1 2 Enter
Sleep 2
Wait+Screen /PS>\ example\ modifier\ \-\-noprefix\ one\-\-\-/
