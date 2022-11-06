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

  export ZLS_COLOURS="${lines[1]}"
  zstyle ":completion:${curcontext}:*" list-colors "${lines[1]}"
  zstyle ":completion:*:default*" list-colors "${lines[1]}"
  
  # shellcheck disable=SC2034,2206
  lines=(${lines[@]:1})

  # shellcheck disable=SC2034,2206
  local vals=(${lines%%%%$'\t'*})
  # shellcheck disable=SC2034,2206
  local displays=(${lines##*$'\t'})

  local suffix=' '
  [[ ${vals[1]} == *$'\001' ]] && suffix=''
  # shellcheck disable=SC2034,2206
  vals=(${vals%%%%$'\001'*})

  # -------- Quotes ---------- #
  # ------ OLD ------- #
  # compadd -S "${suffix}" -l -d displays -a -- vals
  # compadd -l -Q -S "${suffix}" -d displays -a -- vals
  # ------ OLD ------- #

  # ------- New ----------
   ISUFFIX="${suffix}"
  # compset -S "${suffix}"
  _describe "testing this" displays vals

  # ------- Alternate ----------
  # local expl

  # Display message description even if no matches, with -x
  # Problem is that the builtin message does not disappears
  # _description -x vals expl "testing that" 

  # _description vals expl "testing that" 

  # Add completions
  # compadd -l -Q -S "${suffix}" "$expl[@]" -d displays -a vals 
  # compadd -l -Q -S "${suffix}" -d displays -a vals 
}
compquote '' 2>/dev/null && _%v_completion
compdef _%v_completion %v
`, cmd.Name(), cmd.Name(), uid.Executable(), uid.Executable(), uid.Executable(), cmd.Name(), cmd.Name(), cmd.Name())
}
