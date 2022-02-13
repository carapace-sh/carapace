set edit:completion:arg-completer[example] = {|@arg|
    example _carapace elvish (all $arg) | from-json | all (one) | each {|c| edit:complex-candidate $c[Value] &display=$c[Display] &code-suffix=$c[CodeSuffix] }
}

