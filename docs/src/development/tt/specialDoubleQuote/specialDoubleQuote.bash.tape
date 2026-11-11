Set Shell "bash"
Source "../bash.tape"
Type `example special "p`
Wait+Screen /special "p/
Sleep 1
Tab@1
Enter
Sleep 2
Wait+Screen /\$/
