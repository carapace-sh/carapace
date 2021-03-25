#compdef example
function _example_completion {
  local IFS=$'\n'
  
  # shellcheck disable=SC2207,SC2086,SC2154
  if echo ${words}"''" | xargs echo 2>/dev/null > /dev/null; then
    # shellcheck disable=SC2207,SC2086
    local c=($(echo ${words}"''" | xargs example _carapace zsh _ ))
  elif echo ${words} | sed "s/\$/'/" | xargs echo 2>/dev/null > /dev/null; then
    # shellcheck disable=SC2207,SC2086
    local c=($(echo ${words} | sed "s/\$/'/" | xargs example _carapace zsh _ ))
  else
    # shellcheck disable=SC2207,SC2086
    local c=($(echo ${words} | sed 's/$/"/'  | xargs example _carapace zsh _ ))
  fi

  # shellcheck disable=SC2034,2206
  local vals=(${c%%$'\t'*})
  # shellcheck disable=SC2034,2206
  local descriptions=(${c##*$'\t'})
  compadd -S '' -d descriptions -a -- vals
}
compquote '' 2>/dev/null && _example_completion
compdef _example_completion example

