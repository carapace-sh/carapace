#compdef example
function _example_completion {
  # shellcheck disable=SC2086
  eval "$(example _carapace zsh _ ${^words//\\ / }'')"

}
compquote '' 2>/dev/null && _example_completion
compdef _example_completion example

