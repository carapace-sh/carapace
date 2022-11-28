#!/bin/bash
_example_completion() {
  export COMP_WORDBREAKS

  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local IFS=$'\n'
  mapfile -t COMPREPLY < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs example _carapace bash)
  [[ "${COMPREPLY[*]}" == "" ]] && COMPREPLY=() # fix for mapfile creating a non-empty array from empty command output

  compopt -o nospace
}



_example_completion_ble() {
  if [[ ${BLE_ATTACHED-} ]]; then
    [[ :$comp_type: == *:auto:* ]] && return

    compopt -o ble/no-default
    bleopt complete_menu_style=desc

    local compline="${COMP_LINE:0:${COMP_POINT}}"
    local IFS=$'\n'
    local c
    mapfile -t c < <(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs example _carapace bash-ble)
    [[ "${c[*]}" == "" ]] && c=() # fix for mapfile creating a non-empty array from empty command output

    local cand
    for cand in "${c[@]}"; do
      [ ! -z "$cand" ] && ble/complete/cand/yield mandb "${cand%$'\t'*}" "${cand##*$'\t'}"
    done
  else
    complete -F _example_completion example
  fi
}

complete -F _example_completion_ble example

