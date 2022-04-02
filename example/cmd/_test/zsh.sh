#compdef example
function _example_completion {
  local IFS=$'\n'
  
  # shellcheck disable=SC2207,SC2086,SC2154
  if echo ${words}"''" | xargs echo 2>/dev/null > /dev/null; then
    # shellcheck disable=SC2207,SC2086
    local c=($(echo ${words}"''" | xargs example _carapace zsh ))
  elif echo ${words} | sed "s/\$/'/" | xargs echo 2>/dev/null > /dev/null; then
    # shellcheck disable=SC2207,SC2086
    local c=($(echo ${words} | sed "s/\$/'/" | xargs example _carapace zsh))
  else
    # shellcheck disable=SC2207,SC2086
    local c=($(echo ${words} | sed 's/$/"/'  | xargs example _carapace zsh))
  fi

  # shellcheck disable=SC2034,2206
  local vals=(${c%%$'\t'*})
  # shellcheck disable=SC2034,2206
  local displays=(${c##*$'\t'})

  local suffix=' '
  [[ ${vals[1]} == *$'\001' ]] && suffix=''
  # shellcheck disable=SC2034,2206
  vals=(${vals%%$'\001'*})

  if [[ ${displays[*]} == *$'\002'* ]]; then # with descriptions
    compadd -l -S "${suffix}" -d displays -a -- vals
  else
    compadd -S "${suffix}" -d displays -a -- vals
  fi
}
compquote '' 2>/dev/null && _example_completion
compdef _example_completion example

