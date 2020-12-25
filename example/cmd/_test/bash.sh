#!/bin/bash
_example_callback() {
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  # TODO
  #if [[ $last =~ ^[\"\'] ]] && ! echo "$last" | xargs echo 2>/dev/null >/dev/null ; then
  #    compline="${compline}${last:0:1}"
  #    last="${last// /\\\\ }" 
  #fi

  echo "$compline" | sed -e "s/ $/ ''/" -e 's/"/\"/g' | xargs example _carapace bash "$1"
}

_example_completions() {
  local cur prev #words cword split
  _init_completion -n /=:.,
  local curprefix
  curprefix="$(echo "$cur" | sed -r 's_^(.*[:=])?.*_\1_')"
  local compline="${COMP_LINE:0:${COMP_POINT}}"
 
  # TODO
  #if [[ $last =~ ^[\"\'] ]] && ! echo "$last" | xargs echo 2>/dev/null >/dev/null ; then
  #    compline="${compline}${last:0:1}"
  #    last="${last// /\\\\ }" 
  #else
  #    last="${last// /\\\ }" 
  #fi

  local state
  state="$(echo "$compline" | sed -e "s/ \$/ ''/" -e 's/"/\"/g' | xargs example _carapace bash state)"

  local IFS=$'\n'

  case $state in

    '_example' )
      if [[ $cur == -* ]]; then
        case $cur in
          -p=* | --persistentFlag=*)
            cur=${cur#*=}
            curprefix=${curprefix#*=}
            COMPREPLY=($())
            ;;

          *)
            COMPREPLY=($(compgen -W $'--array\t--array (multiflag)\n-a\t-a (multiflag)\n--persistentFlag\t--persistentFlag (Help message for persistentFlag)\n-p\t-p (Help message for persistentFlag)\n--toggle\t--toggle (Help message for toggle)\n-t\t-t (Help message for toggle)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;
        esac
      else
        case $prev in
          -a | --array)
            COMPREPLY=($())
            ;;

          *)
            COMPREPLY=($(compgen -W $'action\taction (action example)\nalias\taction (action example)\ncallback\tcallback (callback example)\ncondition\tcondition (condition example)\nhelp\thelp (Help about any command)\ninjection\tinjection (just trying to break things)\nmultiparts\tmultiparts (multiparts example)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;
        esac
      fi
      ;;


    '_example__action' )
      if [[ $cur == -* ]]; then
        case $cur in
          -o=* | --optarg=*)
            cur=${cur#*=}
            curprefix=${curprefix#*=}
            COMPREPLY=($(compgen -W $'blue\tblue\nred\tred\ngreen\tgreen\nyellow\tyellow' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;

          *)
            COMPREPLY=($(compgen -W $'--directories\t--directories (files flag)\n--files\t--files (files flag)\n-f\t-f (files flag)\n--groups\t--groups (groups flag)\n-g\t-g (groups flag)\n--kill\t--kill (kill signals)\n-k\t-k (kill signals)\n--message\t--message (message flag)\n-m\t-m (message flag)\n--net_interfaces\t--net_interfaces (net_interfaces flag)\n-n\t-n (net_interfaces flag)\n--optarg\t--optarg (optional arg with default value blue)\n-o\t-o (optional arg with default value blue)\n--usergroup\t--usergroup (user:group flag)\n--users\t--users (users flag)\n-u\t-u (users flag)\n--values\t--values (values flag)\n-v\t-v (values flag)\n--values_described\t--values_described (values with description flag)\n-d\t-d (values with description flag)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;
        esac
      else
        case $prev in
          --directories)
            COMPREPLY=($(compgen -S / -d -- "$cur"))
            ;;

          -f | --files)
            COMPREPLY=($(compgen -S / -d -- "$cur"; compgen -f -X '!*.go' -- "$cur"))
            ;;

          -g | --groups)
            COMPREPLY=($(eval $(_example_callback '_example__action##groups')))
            ;;

          -k | --kill)
            COMPREPLY=($(compgen -W $'ABRT\tABRT (Abnormal termination)\nALRM\tALRM (Virtual alarm clock)\nBUS\tBUS (BUS error)\nCHLD\tCHLD (Child status has changed)\nCONT\tCONT (Continue stopped process)\nFPE\tFPE (Floating-point exception)\nHUP\tHUP (Hangup detected on controlling terminal)\nILL\tILL (Illegal instruction)\nINT\tINT (Interrupt from keyboard)\nKILL\tKILL (Kill, unblockable)\nPIPE\tPIPE (Broken pipe)\nPOLL\tPOLL (Pollable event occurred)\nPROF\tPROF (Profiling alarm clock timer expired)\nPWR\tPWR (Power failure restart)\nQUIT\tQUIT (Quit from keyboard)\nSEGV\tSEGV (Segmentation violation)\nSTKFLT\tSTKFLT (Stack fault on coprocessor)\nSTOP\tSTOP (Stop process, unblockable)\nSYS\tSYS (Bad system call)\nTERM\tTERM (Termination request)\nTRAP\tTRAP (Trace/breakpoint trap)\nTSTP\tTSTP (Stop typed at keyboard)\nTTIN\tTTIN (Background read from tty)\nTTOU\tTTOU (Background write to tty)\nURG\tURG (Urgent condition on socket)\nUSR1\tUSR1 (User-defined signal 1)\nUSR2\tUSR2 (User-defined signal 2)\nVTALRM\tVTALRM (Virtual alarm clock)\nWINCH\tWINCH (Window size change)\nXCPU\tXCPU (CPU time limit exceeded)\nXFSZ\tXFSZ (File size limit exceeded)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;

          -m | --message)
            COMPREPLY=($(compgen -W $'_\t_\nERR\tERR (message example)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;

          -n | --net_interfaces)
            COMPREPLY=($(eval $(_example_callback '_example__action##net_interfaces')))
            ;;

          --usergroup)
            COMPREPLY=($(eval $(_example_callback '_example__action##usergroup')))
            ;;

          -u | --users)
            COMPREPLY=($(eval $(_example_callback '_example__action##users')))
            ;;

          -v | --values)
            COMPREPLY=($(compgen -W $'values\tvalues\nexample\texample' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;

          -d | --values_described)
            COMPREPLY=($(compgen -W $'values\tvalues (valueDescription)\nexample\texample (exampleDescription)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;

          *)
            COMPREPLY=($(eval $(_example_callback '_')))
            ;;
        esac
      fi
      ;;


    '_example__callback' )
      if [[ $cur == -* ]]; then
        case $cur in

          *)
            COMPREPLY=($(compgen -W $'--callback\t--callback (Help message for callback)\n-c\t-c (Help message for callback)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;
        esac
      else
        case $prev in
          -c | --callback)
            COMPREPLY=($(eval $(_example_callback '_example__callback##callback')))
            ;;

          *)
            COMPREPLY=($(eval $(_example_callback '_')))
            ;;
        esac
      fi
      ;;


    '_example__condition' )
      if [[ $cur == -* ]]; then
        case $cur in

          *)
            COMPREPLY=($(compgen -W $'--required\t--required (required flag)\n-r\t-r (required flag)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;
        esac
      else
        case $prev in
          -r | --required)
            COMPREPLY=($(compgen -W $'valid\tvalid\ninvalid\tinvalid' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;

          *)
            COMPREPLY=($(eval $(_example_callback '_')))
            ;;
        esac
      fi
      ;;


    '_example__help' )
      if [[ $cur == -* ]]; then
        case $cur in

          *)
            COMPREPLY=($())
            ;;
        esac
      else
        case $prev in

          *)
            COMPREPLY=($(eval $(_example_callback '_')))
            ;;
        esac
      fi
      ;;


    '_example__injection' )
      if [[ $cur == -* ]]; then
        case $cur in

          *)
            COMPREPLY=($())
            ;;
        esac
      else
        case $prev in

          *)
            COMPREPLY=($(eval $(_example_callback '_')))
            ;;
        esac
      fi
      ;;


    '_example__multiparts' )
      if [[ $cur == -* ]]; then
        case $cur in

          *)
            COMPREPLY=($(compgen -W $'--at\t--at (multiparts with @ as divider)\n--colon\t--colon (multiparts with : as divider )\n--comma\t--comma (multiparts with , as divider)\n--dot\t--dot (multiparts with . as divider)\n--dotdotdot\t--dotdotdot (multiparts with ... as divider)\n--equals\t--equals (multiparts with = as divider)\n--none\t--none (multiparts without divider)\n--slash\t--slash (multiparts with / as divider)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;
        esac
      else
        case $prev in
          --at)
            COMPREPLY=($(eval $(_example_callback '_example__multiparts##at')))
            ;;

          --colon)
            COMPREPLY=($(eval $(_example_callback '_example__multiparts##colon')))
            ;;

          --comma)
            COMPREPLY=($(eval $(_example_callback '_example__multiparts##comma')))
            ;;

          --dot)
            COMPREPLY=($(eval $(_example_callback '_example__multiparts##dot')))
            ;;

          --dotdotdot)
            COMPREPLY=($(eval $(_example_callback '_example__multiparts##dotdotdot')))
            ;;

          --equals)
            COMPREPLY=($(eval $(_example_callback '_example__multiparts##equals')))
            ;;

          --none)
            COMPREPLY=($(eval $(_example_callback '_example__multiparts##none')))
            ;;

          --slash)
            COMPREPLY=($(eval $(_example_callback '_example__multiparts##slash')))
            ;;

          *)
            COMPREPLY=($(eval $(_example_callback '_')))
            ;;
        esac
      fi
      ;;

  esac

  [[ $cur =~ ^[\"\'] ]] && COMPREPLY=("${COMPREPLY[@]//\\ /\ }")
  
  [[ ${#COMPREPLY[@]} -gt 1 ]] && for entry in "${COMPREPLY[@]}"; do
    value="${entry%	*}"
    display="${entry#*	}"
    if [[ "${value::1}" != "${display::1}"  ]]; then # inserted value differs from display value
       [[ "$(printf  "%c\n" "${COMPREPLY[@]#*	}" | uniq | wc -l)" -eq 1 ]] && COMPREPLY=("${COMPREPLY[@]}" "") # prevent insertion if all have same first character (workaround for #164)
      break
    fi
  done

  [[ ${#COMPREPLY[@]} -gt 1 ]] && COMPREPLY=("${COMPREPLY[@]#*	}") # show visual part (all after tab)
  [[ ${#COMPREPLY[@]} -eq 1 ]] && COMPREPLY=( ${COMPREPLY[0]%	*} ) # show value to insert (all before tab) https://stackoverflow.com/a/10130007
  [[ ${COMPREPLY[0]} == *[/=@:.,] ]] && compopt -o nospace
}

complete -F _example_completions example

