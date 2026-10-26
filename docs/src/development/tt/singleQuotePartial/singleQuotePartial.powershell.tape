Set Shell "pwsh"
Source "../powershell.tape"
Type `example action embeddedP1 'embeddedP2 w'i`
Wait+Screen /embeddedP2/
Sleep 1
Tab@1
Sleep 2
Wait+Screen /PS>\ example\ action\ embeddedP1\ 'embeddedP2\ w'i/
