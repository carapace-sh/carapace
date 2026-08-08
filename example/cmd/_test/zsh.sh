#compdef example
function _example_completion {
  local compline=${words[@]:0:$CURRENT}
  local IFS=$'\n'
  local lines

  # shellcheck disable=SC2086,SC2154,SC2155
  lines="$(echo "${compline}''" | CARAPACE_COMPLINE="${compline}" CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs example _carapace zsh 2>/dev/null)"
  if [ $? -eq 1 ]; then
    lines="$(echo "${compline}'" | CARAPACE_COMPLINE="${compline}" CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs example _carapace zsh 2>/dev/null)"
    if [ $? -eq 1 ]; then
      lines="$(echo "${compline}\"" | CARAPACE_COMPLINE="${compline}" CARAPACE_ZSH_HASH_DIRS="$(hash -d)" xargs example _carapace zsh 2>/dev/null)"
    fi
  fi

  local zstyle message noprefix data
  IFS=$'\001' read -r -d '' zstyle message noprefix data <<<"${lines}"
  # shellcheck disable=SC2154
  zstyle ":completion:${curcontext}:*" list-colors "${zstyle}"
  zstyle ":completion:${curcontext}:*" group-name ''
  [ -z "$message" ] || _message -r "${message}"
  [[ "${noprefix}" = "true" ]] && compstate[insert]=menu
  
  local block tag suffix displays values
  local -A tags
  while IFS=$'\002' read -r -d $'\002' block; do
    IFS=$'\003' read -r -d '' tag suffix displays values <<<"${block}"
    # shellcheck disable=SC2034
    local -a displaysArr=("${(@f)displays}")
    # shellcheck disable=SC2034
    local -a valuesArr=("${(@f)values}")

    # Skip empty blocks. Upstream renders nothing for a single-value default
    # (empty-suffix) group, so preserve that behaviour exactly; a suffix group
    # must still render a lone value since zsh manages the separator itself.
    if [[ -z "$suffix" ]]; then
      [[ ${#valuesArr[@]} -gt 1 ]] || continue
    else
      [[ ${#valuesArr[@]} -gt 0 ]] || continue
    fi

    # Emit the group header only once per tag: a tag may be split across
    # multiple blocks (one per removable suffix) and reusing the same header
    # would render duplicate group titles. Done after the guard so a skipped
    # group does not consume the tag's one-time header.
    local -a describe_args
    if [[ -z ${tags[$tag]} ]]; then
      describe_args=(-t "${tag}" "${tag}")
      tags[$tag]=1
    else
      describe_args=(-t "${tag}" "")
    fi

    if [[ -z "$suffix" ]]; then
      _describe "${describe_args[@]}" displaysArr valuesArr -Q -S ''
    else
      # Let zsh insert the removable suffix and strip it again when the user
      # continues with a space (native directory/optarg-style handling).
      _describe "${describe_args[@]}" displaysArr valuesArr -Q -S "$suffix" -r ' '
    fi
  done <<<"${data}"
}
compquote '' 2>/dev/null && _example_completion
compdef _example_completion example
