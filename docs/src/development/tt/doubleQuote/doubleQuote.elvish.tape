Set Shell "elvish"
Source "../elvish.tape"
Type `example "ac"`
Wait+Screen /example "ac"/
Sleep 1
Tab@1
Wait+Screen /example "action"/
Sleep 2
Wait+Screen /\$\ example\ "action"/
