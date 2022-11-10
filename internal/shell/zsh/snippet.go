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
    [[ ${#vals[@]} -gt 0 ]] && _describe -t "$tag-${group// /-}" "$group" displays vals ${suffix} -Q
  done
}
compquote '' 2>/dev/null && _%v_completion
compdef _%v_completion %v
`, cmd.Name(), cmd.Name(), uid.Executable(), uid.Executable(), uid.Executable(), cmd.Name(), cmd.Name(), cmd.Name())
}

// Notes: Alternative system based _description and compadd
//
// A few problems:
// If the completions share the same tag, they are not grouped under their group description.
// Completions are not grouped together when they share the same COMPLETION description (short/long flags)
// Does not seem to go significantly faster than _describe calls.
//
// Advantages
// We can pass our options to compadd more easily, which (-l -Q -S)

// Code:
//     local expl=(-S "${suffix}")
//     _description "${tag}" expl ${group}
//     compadd -Q -S${suffix} -l "${expl[@]}" -d displays -a -- vals
