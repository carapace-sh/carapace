let example_completer = {|spans|
    # if the current command is an alias, get it's expansion
    let expanded_alias = (scope aliases | where name == $spans.0 | $in.0?.expansion?)

    # overwrite
    let spans = (if $expanded_alias != null  {
      # put the first word of the expanded alias first in the span
      $spans | skip 1 | prepend ($expanded_alias | split row " " | take 1)
    } else {
      $spans | skip 1 | prepend ($spans.0)
    })

    example _carapace nushell ...$spans | from json
}
