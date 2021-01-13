edit:completion:arg-completer[example] = [@arg]{
    if (eq $arg[-1] "") {
        arg[-1] = "''"
    }
    example _carapace elvish _ (all $arg) | from-json | all (one) | each [c]{ edit:complex-candidate $c[Value] &display=$c[Display] }
}

