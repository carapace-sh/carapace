// Package zsh provides zsh completion
package zsh

import (
	"fmt"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

// Snippet creates the zsh completion script.
func Snippet(cmd *cobra.Command) string {
	return fmt.Sprintf(`#compdef %v

function _%v_completion {
  local IFS=$'\n'
  
  # shellcheck disable=SC2207,SC2086,SC2154
  if echo ${words}"''" | xargs echo 2>/dev/null > /dev/null; then
    # shellcheck disable=SC2207,SC2086
    local lines=$(echo ${words}"''" | CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs %v _carapace zsh)
  elif echo ${words} | sed "s/\$/'/" | xargs echo 2>/dev/null > /dev/null; then
    # shellcheck disable=SC2207,SC2086
    local lines=$(echo ${words} | sed "s/\$/'/" | CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs %v _carapace zsh)
  else
    # shellcheck disable=SC2207,SC2086
    local lines=$(echo ${words} | sed 's/$/"/' | CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs %v _carapace zsh)
  fi

  # Header ("message":"suffix specs")
  read -r header <<< "${lines}"
  header=${header//$(printf '\t')/:}
  IFS=$':' read message suffix rm_suffix <<< "${header}"
  [[ -n ${message} ]] && _message -r "${message}"
  lines=${lines#*$'\n'}

  # Completion options
  local compOpts=( -qS "${suffix}" -r "${rm_suffix}" -Q )
  
  # Styles
  # shellcheck disable=SC2034,2206
  read -r styles <<< "${lines}"
  export ZLS_COLOURS="${styles}"
  zstyle ":completion:${curcontext}:*" list-colors "${styles}"
  lines=${lines#*$'\n'}

  # Scan/build/register completions groups
  local i=1 s=1
  local tag group vals=() displays=()

  while read -r line ; do
    # A new line: done with current group, register its completions and reset.
    if [[ $line == "" ]]; then
        _describe -t "$tag-${group// /-}" "$group" displays vals "${compOpts[@]}" 
        ((i++)) && s=1 && vals=() && displays=() && tag= && group= && continue
    fi

    # Add completion candidate to current group
    if [[ $s == 0 ]]; then
        vals+=( "${line%%%%$'\t'*}" )
        displays+=( "${line##*$'\t'}" ) 
        continue
    fi

    # Or read group header
    IFS=$':' read tag group <<< "${line}" && s=0 
  done < <(printf "%%s\n\n" "${lines}") 
}
compquote '' 2>/dev/null && _%v_completion
compdef _%v_completion %v
`, cmd.Name(), cmd.Name(), uid.Executable(), uid.Executable(), uid.Executable(), cmd.Name(), cmd.Name(), cmd.Name())
}
