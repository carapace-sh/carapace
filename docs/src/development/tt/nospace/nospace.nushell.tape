Set Shell "nu"
Source "../nushell.tape"
Type "example modifier --nospace o"
Wait+Screen /--nospace o/
Sleep 1
Tab@1
Sleep 2
Wait+Screen /\$\ example\ modifier\ \-\-nospace\ one,/
