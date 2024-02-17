#!/bin/bash
_example_completion() {
  export COMP_LINE
  export COMP_POINT
  export COMP_TYPE
  export COMP_WORDBREAKS

  local nospace data compline="${COMP_LINE:0:${COMP_POINT}}"

  if echo ${compline}"''" | xargs echo 2>/dev/null > /dev/null; then
  	data=$(echo ${compline}"''" | xargs example _carapace bash)
  elif echo ${compline} | sed "s/\$/'/" | xargs echo 2>/dev/null > /dev/null; then
  	data=$(echo ${compline} | sed "s/\$/'/" | xargs example _carapace bash)
  else
  	data=$(echo ${compline} | sed 's/$/"/' | xargs example _carapace bash)
  fi

  IFS=$'\001' read -r -d '' nospace data <<<"${data}"
  mapfile -t COMPREPLY < <(echo "${data}")
  unset COMPREPLY[-1]

  [ "${nospace}" = true ] && compopt -o nospace
  local IFS=$'\n'
  [[ "${COMPREPLY[*]}" == "" ]] && COMPREPLY=() # fix for mapfile creating a non-empty array from empty command output
}

complete -o noquote -F _example_completion example

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

