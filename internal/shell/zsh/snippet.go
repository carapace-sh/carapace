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
  
  local go_start
  go_start=$(zsh_timer)
  
  local lines
  lines="$(echo "${completion_input}" | CARAPACE_COMPLINE="${words}" CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs %v _carapace zsh)"
  
  local zstyle message data
  IFS=$'\001' read -r -d '' zstyle message data <<<"${lines}"
  # shellcheck disable=SC2154
  zstyle ":completion:${curcontext}:*" list-colors "${zstyle}"
  zstyle ":completion:${curcontext}:*" group-name ''
  [ -z "$message" ] || _message -r "${message}"
  
  local block tag suffix displays values
  local -A tags
  while IFS=$'\002' read -r -d $'\002' block; do
    IFS=$'\003' read -r -d '' tag suffix displays values <<<"${block}"
    
    local -a displaysArr=("${(f@)displays}")
    local -a valuesArr=("${(f@)values}")

    local -a describe_args
    if [[ -z ${tags[$tag]} ]]; then
      describe_args=(-t "${tag}" "${tag}")
      tags[$tag]=1
    else
      describe_args=(-t "${tag}" "")
    fi
    
    local separators=" /,.':@="
    if [[ "$suffix" == "" ]]; then
      _describe "${describe_args[@]}" displaysArr valuesArr -Q -S ' ' -r "${separators}0-9a-zA-Z"
    elif [[ "$separators" == *"$suffix"* ]]; then
      _describe "${describe_args[@]}" displaysArr valuesArr -Q -S "$suffix" -r ' '
    else
      _describe "${describe_args[@]}" displaysArr valuesArr -Q -S "$suffix" -r "0-9a-zA-Z"
    fi
  done <<<"${data}"
}
compquote '' 2>/dev/null && _%v_completion
compdef _%v_completion %v
`, cmd.Name(), cmd.Name(), uid.Executable(), cmd.Name(), cmd.Name(), cmd.Name())
}
