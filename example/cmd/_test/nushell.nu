let example_completer = {|spans| 
    example _carapace nushell ...$spans | from json
}
