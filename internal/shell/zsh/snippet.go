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
  
  # Grouped completions
  # shellcheck disable=SC2034,2206
  lines=(${lines[@]:2})

  # Ensure ordering of groups and tags
  local tag_order group_order

  # Process and generate completions by groups (one per line)
  for group in "${lines[@]}"; do
    candidates=($(xargs -n1 <<< ${group}))

    # Header (tag and group description)
    IFS=$':' read tag group <<< "${candidates[1]}"
    candidates=(${candidates[@]:1})

    # shellcheck disable=SC2034,2206
    local vals=(${candidates%%%%$'\t'*})
    # shellcheck disable=SC2034,2206
    local displays=(${candidates##*$'\t'})

    # Suffix
    local suffix=-S' '
    [[ ${vals[1]} == *$'\001' ]] && suffix=
    # shellcheck disable=SC2034,2206
    vals=(${vals%%%%$'\001'*})

    # Generate completions
    #ISUFFIX="${suffix}"
    [[ ${#vals[@]} -gt 0 ]] && _describe -t "$tag" "$group" displays vals ${suffix}

    # Append to tag/group ordering
    group_order+="$(printf %%q "$group") "
    zstyle ":completion:${curcontext}:*" tag-order "$tag:$(printf %%q "$group")"
  done

  zstyle ":completion:${curcontext}:*" group-order "$group_order"
}
compquote '' 2>/dev/null && _%v_completion
compdef _%v_completion %v
`, cmd.Name(), cmd.Name(), uid.Executable(), uid.Executable(), uid.Executable(), cmd.Name(), cmd.Name(), cmd.Name())
}
