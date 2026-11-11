Set Shell "pwsh"
Source "../powershell.tape"
Type `example "ac"`
Wait+Screen /example "ac"/
Sleep 1
Tab@1
Sleep 2
Wait+Screen /PS>\ example\ "ac"/
