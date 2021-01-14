#compdef example
function _example_completion {
  local IFS=$'\n'
  # shellcheck disable=SC2207,SC2086
  local c=($(example _carapace zsh _ ${^words//\\ / }''))
  # shellcheck disable=SC2034,2206
  local vals=(${c%%$'\t'*})
  # shellcheck disable=SC2034,2206
  local descriptions=(${c##*$'\t'})
  compadd -Q -S '' -d descriptions -a -- vals
}
compquote '' 2>/dev/null && _example_completion
compdef _example_completion example

