// Package zsh provides zsh completion
package zsh

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/carapace-sh/carapace/pkg/uid"
)

// Snippet creates the zsh completion script
func Snippet(cmd *cobra.Command) string {
	return fmt.Sprintf(`#compdef %v
function _%v_completion {
  local words=${words[@]:0:$CURRENT}
  local IFS=$'\n'

  # shellcheck disable=SC2086,SC2154,SC2155
  local completion_input
  if echo ${words}"''" | xargs echo 2>/dev/null > /dev/null; then
    completion_input="${words}''"
  elif echo "${words[1,-2]} ${words[-1]}'" | xargs echo 2>/dev/null > /dev/null; then
    completion_input="${words[1,-2]} ${words[-1]}'"
  else
    completion_input="${words[1,-2]} ${words[-1]}\\"
  fi

  local lines
  lines="$(echo "${completion_input}" | CARAPACE_COMPLINE="${words}" CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs %v _carapace zsh)"

  local zstyle message data
  IFS=$'\001' read -r -d '' zstyle message data <<<"${lines}"
  # shellcheck disable=SC2154
  zstyle ":completion:${curcontext}:*" list-colors "${zstyle}"
  zstyle ":completion:${curcontext}:*" group-name ''
  [ -z "$message" ] || _message -r "${message}"
  
  local block tag displays values suffixes
  while IFS=$'\002' read -r -d $'\002' block; do
    IFS=$'\003' read -r -d '' tag displays values suffixes <<<"${block}"

    local -a displaysArr=("${(f@)displays}")
    local -a valuesArr=("${(f@)values}")
    local -a suffixesArr=("${(f@)suffixes}")

    typeset -A grouped_values=()
    typeset -A grouped_displays=()

    for i in {1..${#valuesArr[@]}}; do
      local suffix_key="${suffixesArr[$i]:-__NOSUFFIX__}"
      grouped_values[$suffix_key]+="${valuesArr[$i]}"$'\n'
      grouped_displays[$suffix_key]+="${displaysArr[$i]}"$'\n'
    done

    local first_call=1
    for suffix_key in "${(@k)grouped_values}"; do
      local -a s_values=("${(f@)${grouped_values[$suffix_key]%%$'\n'}}")
      local -a s_displays=("${(f@)${grouped_displays[$suffix_key]%%$'\n'}}")

      if [[ ${#s_values[@]} -eq 0 ]]; then
        continue
      fi

      local -a describe_args
      if (( first_call )); then
        describe_args=(-t "${tag}" "${tag}")
        first_call=0
      else
		describe_args=(-t "${tag}" "")
      fi

      local separators=" /,.':@"
      if [[ "$suffix_key" == "__NOSUFFIX__" ]]; then
        _describe "${describe_args[@]}" s_displays s_values -Q
      elif [[ "$separators" == *"$suffix_key"* ]]; then
        _describe "${describe_args[@]}" s_displays s_values -Q -S "$suffix_key" -r ' '
      else
        _describe "${describe_args[@]}" s_displays s_values -Q -S "$suffix_key" -r "0-9a-zA-Z"
      fi
    done
  done <<<"${data}"
}
compquote '' 2>/dev/null && _%v_completion
compdef _%v_completion %v
`, cmd.Name(), cmd.Name(), uid.Executable(), cmd.Name(), cmd.Name(), cmd.Name())
}
