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

