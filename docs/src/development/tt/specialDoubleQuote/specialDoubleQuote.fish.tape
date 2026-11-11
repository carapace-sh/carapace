Set Shell "fish"
Source "../fish.tape"
Type `example special "p`
Wait+Screen /special "p/
Sleep 1
Tab@1
Enter
Sleep 2
Wait+Screen /\$/
