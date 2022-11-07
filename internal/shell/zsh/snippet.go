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
    local lines=($(echo ${words}"''" | CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs %v _carapace zsh ))
  elif echo ${words} | sed "s/\$/'/" | xargs echo 2>/dev/null > /dev/null; then
    # shellcheck disable=SC2207,SC2086
    local lines=($(echo ${words} | sed "s/\$/'/" | CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs %v _carapace zsh))
  else
    # shellcheck disable=SC2207,SC2086
    local lines=($(echo ${words} | sed 's/$/"/' | CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs %v _carapace zsh))
  fi

  # Return code and message (sanitized)
  header=${lines[1]//$(printf '\t')/:}
  IFS=$':' read retcode message <<< "${header}"
  [[ -n ${message} ]] && _message -r "${message}"

  # Styles
  export ZLS_COLOURS="${lines[2]}"
  zstyle ":completion:${curcontext}:*" list-colors "${lines[2]}"
  # zstyle ":completion:*:default*" list-colors "${lines[2]}"
  
  # shellcheck disable=SC2034,2206
  lines=(${lines[@]:2})

  # Completions (inserted and displayed)
  # shellcheck disable=SC2034,2206
  local vals=(${lines%%%%$'\t'*})
  # shellcheck disable=SC2034,2206
  local displays=(${lines##*$'\t'})

  ## Suffix
  local suffix=' '
  [[ ${vals[1]} == *$'\001' ]] && suffix=''
  # shellcheck disable=SC2034,2206
  vals=(${vals%%%%$'\001'*})

  # New completion generation
  #_message -e "test name" 'this is a message description'
  ISUFFIX="${suffix}"
  #_describe -t 'test name' "test comps" displays vals
  #_describe -t 'test name' "test compother" displays vals
  [[ ${#vals[@]} -gt 0 ]] && _describe "completions" displays vals
  #[[ ${#vals[@]} -gt 0 ]] && _describe -t 'test name' "test comps" displays vals
}
compquote '' 2>/dev/null && _%v_completion
compdef _%v_completion %v
`, cmd.Name(), cmd.Name(), uid.Executable(), uid.Executable(), uid.Executable(), cmd.Name(), cmd.Name(), cmd.Name())
}
