let example_completer = {|spans| 
    load-env { CARAPACE_SHELL: 'nushell' }
    example _carapace nushell ...$spans | from json
}
