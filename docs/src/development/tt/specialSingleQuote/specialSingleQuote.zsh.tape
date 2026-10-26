Set Shell "zsh"
Source "../zsh.tape"
Type `example special 'p`
Wait+Screen /special 'p/
Sleep 1
Tab@1
Enter
Sleep 2
Wait+Screen /ok/
