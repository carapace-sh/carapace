let example__multi_completer = {|spans|
    load-env { CARAPACE_SHELL: 'nushell' }
    example-multi $spans.0 _carapace nushell ...$spans | from json
}

mut current = (($env | default {} config).config | default {} completions)
$current.completions = ($current.completions | default {} external)
$current.completions.external = ($current.completions.external
    | default true enable
    |# backwards compatible workaround for default, see nushell #15654
    | upsert completer { if $in == null { $example__multi_completer } else { $in } })

$env.config = $current

