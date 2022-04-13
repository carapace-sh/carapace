module carapace_example {
  def "nu-complete example" [line: string, pos: int] {
    let words = ($line | str substring [0 $pos] | split row " ")
    if ($line | str substring [0 $pos] | str ends-with " ") {
      example _carapace nushell ($words | append "") | from json
    } else {
      example _carapace nushell $words | from json
    }
  }
  
  export extern "example" [
    ...args: string@"nu-complete example"
  ]
}
use carapace_example *

