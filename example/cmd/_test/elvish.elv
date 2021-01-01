edit:completion:arg-completer[example] = [@arg]{
    if (eq $arg[-1] "") {
        arg[-1] = "''"
    }
    eval (example _carapace elvish _ (all $arg) | slurp) &ns=(ns [&arg=$arg])
}

