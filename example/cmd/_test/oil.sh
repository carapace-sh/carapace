#!/bin/osh
_example_completion() {
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local IFS=$'\n'
  mapfile -t COMPREPLY < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs example _carapace oil "_")
  [[ ""${COMPREPLY[@]}"" == "" ]] && setvar COMPREPLY = %() # fix for mapfile creating a non-empty array from empty command output
  [[ ${COMPREPLY[@]} == *[/=@:.,] ]] && compopt -o nospace
}

complete -F _example_completion example

