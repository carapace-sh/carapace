let example_completer = {|place| 
    example _carapace nushell ...$place.command | from json
}
