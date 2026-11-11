Set Shell "bash"
Source "../bash.tape"
Type `example 'ac'`
Wait+Screen /example 'ac'/
Sleep 1
Tab@1
Sleep 2
Wait+Screen /\$\ example\ action/
