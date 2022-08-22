let external_completer = {|spans| 
  {
    $spans.0: { } # default
    example: { example _carapace nushell $spans | from json }
  } | get $spans.0 | each {|it| do $it}
}

let-env config = {
  external_completer: $external_completer
}
