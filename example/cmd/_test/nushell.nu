let example_completer = {|place| 
    # backwards compatible workaround for positional completer input, see nushell #18791
    let spans = (if ($place | describe) =~ "record" { $place.command } else { $place })
    example _carapace nushell ...$spans | from json
}
