let example_completer = {|spans|
    example _carapace nushell ...$spans | from json
}

let example_previous_completer = ($env.config.completions.external.completer? | default {|spans| null })
$env.config = ($env.config | upsert completions.external {
    enable: true
    completer: {|spans|
        if $spans.0 == "example" {
            do $example_completer $spans
        } else {
            do $example_previous_completer $spans
        }
    }
})
