Set Shell "zsh"
Source "../zsh.tape"
Type "example modifier --noprefix one"
Wait+Screen /--noprefix one/
Sleep 1 Tab@1 2 Enter
Sleep 2
Wait+Screen /example\ modifier\ \-\-noprefix\ one\-\-\-/
