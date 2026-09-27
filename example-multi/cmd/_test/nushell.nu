let example__multi_completer = {|place|
    # backwards compatible workaround for positional completer input, see nushell #18791
    let spans = (if ($place | describe) =~ "record" { $place.command } else { $place })
    example-multi $spans.0 _carapace nushell ...$spans | from json
}

mut current = (($env | default {} config).config | default {} completions)
$current.completions = ($current.completions | default {} external)
$current.completions.external = ($current.completions.external
    | default true enable
    |# backwards compatible workaround for default, see nushell #15654
    | upsert completer { if $in == null { $example__multi_completer } else { $in } })

$env.config = $current

