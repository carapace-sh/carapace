set edit:completion:arg-completer[example] = {|@arg|
    example _carapace elvish (all $arg) | from-json | all (one) | each {|c| 
        if (eq $c[Description] "") {
            edit:complex-candidate $c[Value] &display=(styled $c[Display] $c[Style]) &code-suffix=$c[CodeSuffix]
        } else {
            edit:complex-candidate $c[Value] &display=(styled $c[Display] $c[Style])" ("(styled $c[Description] $c[DescriptionStyle])")" &code-suffix=$c[CodeSuffix]
        }
    }
}

