#compdef example
function _example_completion {
  local IFS=$'\n'
  
  # shellcheck disable=SC2207,SC2086,SC2154
  if echo ${words}"''" | xargs echo 2>/dev/null > /dev/null; then
    # shellcheck disable=SC2207,SC2086
    local lines=($(echo ${words}"''" | CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs example _carapace zsh ))
  elif echo ${words} | sed "s/\$/'/" | xargs echo 2>/dev/null > /dev/null; then
    # shellcheck disable=SC2207,SC2086
    local lines=($(echo ${words} | sed "s/\$/'/" | CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs example _carapace zsh))
  else
    # shellcheck disable=SC2207,SC2086
    local lines=($(echo ${words} | sed 's/$/"/' | CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs example _carapace zsh))
  fi

  export ZLS_COLOURS="${lines[1]}"
  local line_break=$'\n'
  [[ ! "${lines[2]}" == "NONE" ]] && _message -r "${lines[2]//$'\t'/${line_break}}"
  
  # shellcheck disable=SC2034,2206
  lines=(${lines[@]:2})
  # shellcheck disable=SC2034,2206
  local vals=(${lines%$'\t'*})
  # shellcheck disable=SC2034,2206
  local displays=(${lines##*$'\t'})

  local compOpts=( -qS "" -Q )
  _describe -t "sometag" "sometag" displays vals "${compOpts[@]}"
}
compquote '' 2>/dev/null && _example_completion
compdef _example_completion example

