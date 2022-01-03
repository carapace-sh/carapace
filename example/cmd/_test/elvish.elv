set edit:completion:arg-completer[example] = {|@arg|
    example _carapace elvish (all $arg) | from-json | all (one) | each {|c| 
        if (eq $c[Description] "") {
            edit:complex-candidate $c[Value] &display=$c[Display] &code-suffix=$c[CodeSuffix]
        } else {
            edit:complex-candidate $c[Value] &display=$c[Display]" ("(styled $c[Description] magenta)")" &code-suffix=$c[CodeSuffix]
        }
    }
}

