Set Shell "fish"
Source "../fish.tape"
Type `example action embeddedP1 "embeddedP2 with `
Wait+Screen /embeddedP2/
Sleep 1
Tab@1
Sleep 2
Wait+Screen /\$\ example\ action\ embeddedP1\ "embeddedP2\ with\ space"/
