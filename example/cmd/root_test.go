package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/internal/assert"
)

func TestBash(t *testing.T) {
	expected := `#!/bin/bash
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
  curprefix="$(echo "$cur" | sed -r 's_^(.*[:/=])?.*_\1_')"
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
            COMPREPLY=($(compgen -W $'--array (multiflag)\n-a (multiflag)\n--persistentFlag (Help message for persistentFlag)\n-p (Help message for persistentFlag)\n--toggle (Help message for toggle)\n-t (Help message for toggle)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;
        esac
      else
        case $prev in
          -a | --array)
            COMPREPLY=($())
            ;;

          *)
            COMPREPLY=($(compgen -W $'action (action example)\nalias (action example)\ncallback (callback example)\ncondition (condition example)\nhelp (Help about any command)\ninjection (just trying to break things)\nmultiparts (multiparts example)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
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
            COMPREPLY=($(compgen -W $'blue\nred\ngreen\nyellow' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;

          *)
            COMPREPLY=($(compgen -W $'--directories (files flag)\n--files (files flag)\n-f (files flag)\n--groups (groups flag)\n-g (groups flag)\n--hosts (hosts flag)\n--kill (kill signals)\n-k (kill signals)\n--message (message flag)\n-m (message flag)\n--net_interfaces (net_interfaces flag)\n-n (net_interfaces flag)\n--optarg (optional arg with default value blue)\n-o (optional arg with default value blue)\n--usergroup (user:group flag)\n--users (users flag)\n-u (users flag)\n--values (values flag)\n-v (values flag)\n--values_described (values with description flag)\n-d (values with description flag)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
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

          --hosts)
            COMPREPLY=($(eval $(_example_callback '_example__action##hosts')))
            ;;

          -k | --kill)
            COMPREPLY=($(compgen -W $'ABRT (Abnormal termination)\nALRM (Virtual alarm clock)\nBUS (BUS error)\nCHLD (Child status has changed)\nCONT (Continue stopped process)\nFPE (Floating-point exception)\nHUP (Hangup detected on controlling terminal)\nILL (Illegal instruction)\nINT (Interrupt from keyboard)\nKILL (Kill, unblockable)\nPIPE (Broken pipe)\nPOLL (Pollable event occurred)\nPROF (Profiling alarm clock timer expired)\nPWR (Power failure restart)\nQUIT (Quit from keyboard)\nSEGV (Segmentation violation)\nSTKFLT (Stack fault on coprocessor)\nSTOP (Stop process, unblockable)\nSYS (Bad system call)\nTERM (Termination request)\nTRAP (Trace/breakpoint trap)\nTSTP (Stop typed at keyboard)\nTTIN (Background read from tty)\nTTOU (Background write to tty)\nURG (Urgent condition on socket)\nUSR1 (User-defined signal 1)\nUSR2 (User-defined signal 2)\nVTALRM (Virtual alarm clock)\nWINCH (Window size change)\nXCPU (CPU time limit exceeded)\nXFSZ (File size limit exceeded)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;

          -m | --message)
            COMPREPLY=($(compgen -W $'_\nERR (message example)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
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
            COMPREPLY=($(compgen -W $'values\nexample' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;

          -d | --values_described)
            COMPREPLY=($(compgen -W $'values (valueDescription)\nexample (exampleDescription)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
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
            COMPREPLY=($(compgen -W $'--callback (Help message for callback)\n-c (Help message for callback)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
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
            COMPREPLY=($(compgen -W $'--required (required flag)\n-r (required flag)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
            ;;
        esac
      else
        case $prev in
          -r | --required)
            COMPREPLY=($(compgen -W $'valid\ninvalid' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
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
            COMPREPLY=($(compgen -W $'--at (multiparts with @ as divider)\n--colon (multiparts with : as divider )\n--comma (multiparts with , as divider)\n--dot (multiparts with . as divider)\n--dotdotdot (multiparts with ... as divider)\n--equals (multiparts with = as divider)\n--none (multiparts without divider)\n--slash (multiparts with / as divider)' -- "${cur//\\ / }" | sed "s"$'\001'"^${curprefix//\\ / }"$'\001'$'\001'))
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
  [[ ${#COMPREPLY[*]} -eq 1 ]] && COMPREPLY=( ${COMPREPLY[0]% (*} )  # https://stackoverflow.com/a/10130007
  [[ ${COMPREPLY[0]} == *[/=@:.,] ]] && compopt -o nospace
}

complete -F _example_completions example
`
	rootCmd.InitDefaultHelpCmd()
	s, _ := carapace.Gen(rootCmd).Snippet("bash", false)
	assert.Equal(t, expected, s)
}

func TestElvish(t *testing.T) {
	expected := `use str
edit:completion:arg-completer[example] = [@arg]{
  fn _example_callback [uid]{
    if (eq $arg[-1] "") {
        arg[-1] = "''"
    }
    eval (echo (str:join ' ' $arg) | xargs example _carapace elvish $uid | slurp) &ns=(ns [&arg=$arg])
  }

  fn subindex [subcommand]{
    # TODO 'edit:complete-getopt' needs the arguments shortened for subcommmands - pretty optimistic here
    index=1
    for x $arg { if (eq $x $subcommand) { break } else { index = (+ $index 1) } } 
    echo $index
  }
  
  state=(echo (str:join ' ' $arg) | xargs example _carapace elvish state)
  if (eq 1 0) {
  }  elif (eq $state '_example') {
    opt-specs = [
        [&long='array' &short='a' &desc='multiflag' &arg-required=$true &completer=[_]{  }]
        [&long='persistentFlag' &short='p' &desc='Help message for persistentFlag' &arg-optional=$true &completer=[_]{  }]
        [&long='toggle' &short='t' &desc='Help message for toggle' &arg-optional=$true &completer=[_]{ edit:complex-candidate 'true' &display='true'
edit:complex-candidate 'false' &display='false' }]
    ]
    arg-handlers = [
        [_]{ edit:complex-candidate 'action' &display='action (action example)'
edit:complex-candidate 'alias' &display='alias (action example)'
edit:complex-candidate 'callback' &display='callback (callback example)'
edit:complex-candidate 'condition' &display='condition (condition example)'
edit:complex-candidate 'help' &display='help (Help about any command)'
edit:complex-candidate 'injection' &display='injection (just trying to break things)'
edit:complex-candidate 'multiparts' &display='multiparts (multiparts example)' }
    ]
    subargs = $arg[(subindex example):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }  elif (eq $state '_example__action') {
    opt-specs = [
        [&long='directories' &desc='files flag' &arg-required=$true &completer=[_]{ edit:complete-filename $arg[-1] }]
        [&long='files' &short='f' &desc='files flag' &arg-required=$true &completer=[_]{ edit:complete-filename $arg[-1] }]
        [&long='groups' &short='g' &desc='groups flag' &arg-required=$true &completer=[_]{ _example_callback '_example__action##groups' }]
        [&long='hosts' &desc='hosts flag' &arg-required=$true &completer=[_]{ _example_callback '_example__action##hosts' }]
        [&long='kill' &short='k' &desc='kill signals' &arg-required=$true &completer=[_]{ edit:complex-candidate 'ABRT' &display='ABRT (Abnormal termination)'
edit:complex-candidate 'ALRM' &display='ALRM (Virtual alarm clock)'
edit:complex-candidate 'BUS' &display='BUS (BUS error)'
edit:complex-candidate 'CHLD' &display='CHLD (Child status has changed)'
edit:complex-candidate 'CONT' &display='CONT (Continue stopped process)'
edit:complex-candidate 'FPE' &display='FPE (Floating-point exception)'
edit:complex-candidate 'HUP' &display='HUP (Hangup detected on controlling terminal)'
edit:complex-candidate 'ILL' &display='ILL (Illegal instruction)'
edit:complex-candidate 'INT' &display='INT (Interrupt from keyboard)'
edit:complex-candidate 'KILL' &display='KILL (Kill, unblockable)'
edit:complex-candidate 'PIPE' &display='PIPE (Broken pipe)'
edit:complex-candidate 'POLL' &display='POLL (Pollable event occurred)'
edit:complex-candidate 'PROF' &display='PROF (Profiling alarm clock timer expired)'
edit:complex-candidate 'PWR' &display='PWR (Power failure restart)'
edit:complex-candidate 'QUIT' &display='QUIT (Quit from keyboard)'
edit:complex-candidate 'SEGV' &display='SEGV (Segmentation violation)'
edit:complex-candidate 'STKFLT' &display='STKFLT (Stack fault on coprocessor)'
edit:complex-candidate 'STOP' &display='STOP (Stop process, unblockable)'
edit:complex-candidate 'SYS' &display='SYS (Bad system call)'
edit:complex-candidate 'TERM' &display='TERM (Termination request)'
edit:complex-candidate 'TRAP' &display='TRAP (Trace/breakpoint trap)'
edit:complex-candidate 'TSTP' &display='TSTP (Stop typed at keyboard)'
edit:complex-candidate 'TTIN' &display='TTIN (Background read from tty)'
edit:complex-candidate 'TTOU' &display='TTOU (Background write to tty)'
edit:complex-candidate 'URG' &display='URG (Urgent condition on socket)'
edit:complex-candidate 'USR1' &display='USR1 (User-defined signal 1)'
edit:complex-candidate 'USR2' &display='USR2 (User-defined signal 2)'
edit:complex-candidate 'VTALRM' &display='VTALRM (Virtual alarm clock)'
edit:complex-candidate 'WINCH' &display='WINCH (Window size change)'
edit:complex-candidate 'XCPU' &display='XCPU (CPU time limit exceeded)'
edit:complex-candidate 'XFSZ' &display='XFSZ (File size limit exceeded)' }]
        [&long='message' &short='m' &desc='message flag' &arg-required=$true &completer=[_]{ edit:complex-candidate '_' &display='_'
edit:complex-candidate 'ERR' &display='ERR (message example)' }]
        [&long='net_interfaces' &short='n' &desc='net_interfaces flag' &arg-required=$true &completer=[_]{ _example_callback '_example__action##net_interfaces' }]
        [&long='optarg' &short='o' &desc='optional arg with default value blue' &arg-optional=$true &completer=[_]{ edit:complex-candidate 'blue' &display='blue'
edit:complex-candidate 'red' &display='red'
edit:complex-candidate 'green' &display='green'
edit:complex-candidate 'yellow' &display='yellow' }]
        [&long='usergroup' &desc='user\:group flag' &arg-required=$true &completer=[_]{ _example_callback '_example__action##usergroup' }]
        [&long='users' &short='u' &desc='users flag' &arg-required=$true &completer=[_]{ _example_callback '_example__action##users' }]
        [&long='values' &short='v' &desc='values flag' &arg-required=$true &completer=[_]{ edit:complex-candidate 'values' &display='values'
edit:complex-candidate 'example' &display='example' }]
        [&long='values_described' &short='d' &desc='values with description flag' &arg-required=$true &completer=[_]{ edit:complex-candidate 'values' &display='values (valueDescription)'
edit:complex-candidate 'example' &display='example (exampleDescription)' }]
    ]
    arg-handlers = [
      [_]{ edit:complex-candidate 'positional1' &display='positional1'
edit:complex-candidate 'p1' &display='p1' }
      [_]{ edit:complex-candidate 'positional2' &display='positional2'
edit:complex-candidate 'p2' &display='p2' }
    ]
    subargs = $arg[(subindex action):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }  elif (eq $state '_example__callback') {
    opt-specs = [
        [&long='callback' &short='c' &desc='Help message for callback' &arg-required=$true &completer=[_]{ _example_callback '_example__callback##callback' }]
    ]
    arg-handlers = [
      [_]{ _example_callback '_example__callback#1' }
      [_]{ _example_callback '_example__callback#2' }
      [_]{ _example_callback '_example__callback#0' }
      ...
    ]
    subargs = $arg[(subindex callback):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }  elif (eq $state '_example__condition') {
    opt-specs = [
        [&long='required' &short='r' &desc='required flag' &arg-required=$true &completer=[_]{ edit:complex-candidate 'valid' &display='valid'
edit:complex-candidate 'invalid' &display='invalid' }]
    ]
    arg-handlers = [
      [_]{ _example_callback '_example__condition#1' }
    ]
    subargs = $arg[(subindex condition):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }  elif (eq $state '_example__help') {
    opt-specs = [

    ]
    arg-handlers = [

    ]
    subargs = $arg[(subindex help):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }  elif (eq $state '_example__injection') {
    opt-specs = [

    ]
    arg-handlers = [
      [_]{ edit:complex-candidate 'echo fail' &display='echo fail' }
      [_]{ edit:complex-candidate 'echo fail' &display='echo fail' }
      [_]{ edit:complex-candidate 'echo fail' &display='echo fail' }
      [_]{ edit:complex-candidate ' echo fail ' &display=' echo fail ' }
      [_]{ edit:complex-candidate ' echo fail ' &display=' echo fail ' }
      [_]{ edit:complex-candidate ' echo fail ' &display=' echo fail ' }
      [_]{ edit:complex-candidate 'echo fail' &display='echo fail' }
      [_]{ edit:complex-candidate '' &display='' }
      [_]{ edit:complex-candidate 'LAST POSITIONAL VALUE' &display='LAST POSITIONAL VALUE' }
    ]
    subargs = $arg[(subindex injection):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }  elif (eq $state '_example__multiparts') {
    opt-specs = [
        [&long='at' &desc='multiparts with @ as divider' &arg-required=$true &completer=[_]{ _example_callback '_example__multiparts##at' }]
        [&long='colon' &desc='multiparts with \: as divider ' &arg-required=$true &completer=[_]{ _example_callback '_example__multiparts##colon' }]
        [&long='comma' &desc='multiparts with , as divider' &arg-required=$true &completer=[_]{ _example_callback '_example__multiparts##comma' }]
        [&long='dot' &desc='multiparts with . as divider' &arg-required=$true &completer=[_]{ _example_callback '_example__multiparts##dot' }]
        [&long='dotdotdot' &desc='multiparts with ... as divider' &arg-required=$true &completer=[_]{ _example_callback '_example__multiparts##dotdotdot' }]
        [&long='equals' &desc='multiparts with = as divider' &arg-required=$true &completer=[_]{ _example_callback '_example__multiparts##equals' }]
        [&long='none' &desc='multiparts without divider' &arg-required=$true &completer=[_]{ _example_callback '_example__multiparts##none' }]
        [&long='slash' &desc='multiparts with / as divider' &arg-required=$true &completer=[_]{ _example_callback '_example__multiparts##slash' }]
    ]
    arg-handlers = [
      [_]{ _example_callback '_example__multiparts#1' }
      [_]{ _example_callback '_example__multiparts#0' }
      ...
    ]
    subargs = $arg[(subindex multiparts):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }
}
`
	rootCmd.InitDefaultHelpCmd()
	s, _ := carapace.Gen(rootCmd).Snippet("elvish", false)
	assert.Equal(t, expected, s)
}

func TestFish(t *testing.T) {
	expected := `function _example_quote_suffix
  if not commandline -cp | xargs echo 2>/dev/null >/dev/null
    if commandline -cp | sed 's/$/"/'| xargs echo 2>/dev/null >/dev/null
      echo '"'
    else if commandline -cp | sed "s/\$/'/"| xargs echo 2>/dev/null >/dev/null
      echo "'"
    end
  else 
    echo ""
  end
end

function _example_state
  set -lx CURRENT (commandline -cp)
  if [ "$LINE" != "$CURRENT" ]
    set -gx LINE (commandline -cp)
    set -gx STATE (commandline -cp | sed "s/\$/"(_example_quote_suffix)"/" | xargs example _carapace fish state)
  end

  [ "$STATE" = "$argv" ]
end

function _example_callback
  set -lx CALLBACK (commandline -cp | sed "s/\$/"(_example_quote_suffix)"/" | sed "s/ \$/ ''/" | xargs example _carapace fish $argv )
  eval "$CALLBACK"
end

complete -c example -f

complete -c 'example' -f -n '_example_state _example' -l 'array' -s 'a' -d 'multiflag' -r
complete -c 'example' -f -n '_example_state _example' -l 'persistentFlag' -s 'p' -d 'Help message for persistentFlag'
complete -c 'example' -f -n '_example_state _example' -l 'toggle' -s 't' -d 'Help message for toggle' -a '(echo -e "true	\nfalse	")'
complete -c 'example' -f -n '_example_state _example ' -a 'action alias' -d 'action example'
complete -c 'example' -f -n '_example_state _example ' -a 'callback ' -d 'callback example'
complete -c 'example' -f -n '_example_state _example ' -a 'condition ' -d 'condition example'
complete -c 'example' -f -n '_example_state _example ' -a 'help ' -d 'Help about any command'
complete -c 'example' -f -n '_example_state _example ' -a 'injection ' -d 'just trying to break things'
complete -c 'example' -f -n '_example_state _example ' -a 'multiparts ' -d 'multiparts example'


complete -c 'example' -f -n '_example_state _example__action' -l 'directories' -d 'files flag' -a '(__fish_complete_directories)' -r
complete -c 'example' -f -n '_example_state _example__action' -l 'files' -s 'f' -d 'files flag' -a '(__fish_complete_suffix ".go")' -r
complete -c 'example' -f -n '_example_state _example__action' -l 'groups' -s 'g' -d 'groups flag' -a '(_example_callback _example__action##groups)' -r
complete -c 'example' -f -n '_example_state _example__action' -l 'hosts' -d 'hosts flag' -a '(_example_callback _example__action##hosts)' -r
complete -c 'example' -f -n '_example_state _example__action' -l 'kill' -s 'k' -d 'kill signals' -a '(echo -e "ABRT	Abnormal termination\nALRM	Virtual alarm clock\nBUS	BUS error\nCHLD	Child status has changed\nCONT	Continue stopped process\nFPE	Floating-point exception\nHUP	Hangup detected on controlling terminal\nILL	Illegal instruction\nINT	Interrupt from keyboard\nKILL	Kill, unblockable\nPIPE	Broken pipe\nPOLL	Pollable event occurred\nPROF	Profiling alarm clock timer expired\nPWR	Power failure restart\nQUIT	Quit from keyboard\nSEGV	Segmentation violation\nSTKFLT	Stack fault on coprocessor\nSTOP	Stop process, unblockable\nSYS	Bad system call\nTERM	Termination request\nTRAP	Trace/breakpoint trap\nTSTP	Stop typed at keyboard\nTTIN	Background read from tty\nTTOU	Background write to tty\nURG	Urgent condition on socket\nUSR1	User-defined signal 1\nUSR2	User-defined signal 2\nVTALRM	Virtual alarm clock\nWINCH	Window size change\nXCPU	CPU time limit exceeded\nXFSZ	File size limit exceeded")' -r
complete -c 'example' -f -n '_example_state _example__action' -l 'message' -s 'm' -d 'message flag' -a '(echo -e "_	\nERR	message example")' -r
complete -c 'example' -f -n '_example_state _example__action' -l 'net_interfaces' -s 'n' -d 'net_interfaces flag' -a '(_example_callback _example__action##net_interfaces)' -r
complete -c 'example' -f -n '_example_state _example__action' -l 'optarg' -s 'o' -d 'optional arg with default value blue' -a '(echo -e "blue	\nred	\ngreen	\nyellow	")'
complete -c 'example' -f -n '_example_state _example__action' -l 'usergroup' -d 'user\:group flag' -a '(_example_callback _example__action##usergroup)' -r
complete -c 'example' -f -n '_example_state _example__action' -l 'users' -s 'u' -d 'users flag' -a '(_example_callback _example__action##users)' -r
complete -c 'example' -f -n '_example_state _example__action' -l 'values' -s 'v' -d 'values flag' -a '(echo -e "values	\nexample	")' -r
complete -c 'example' -f -n '_example_state _example__action' -l 'values_described' -s 'd' -d 'values with description flag' -a '(echo -e "values	valueDescription\nexample	exampleDescription")' -r
complete -c 'example' -f -n '_example_state _example__action' -a '(_example_callback _)'


complete -c 'example' -f -n '_example_state _example__callback' -l 'callback' -s 'c' -d 'Help message for callback' -a '(_example_callback _example__callback##callback)' -r
complete -c 'example' -f -n '_example_state _example__callback' -a '(_example_callback _)'


complete -c 'example' -f -n '_example_state _example__condition' -l 'required' -s 'r' -d 'required flag' -a '(echo -e "valid	\ninvalid	")' -r
complete -c 'example' -f -n '_example_state _example__condition' -a '(_example_callback _)'


complete -c 'example' -f -n '_example_state _example__help' -a '(_example_callback _)'


complete -c 'example' -f -n '_example_state _example__injection' -a '(_example_callback _)'


complete -c 'example' -f -n '_example_state _example__multiparts' -l 'at' -d 'multiparts with @ as divider' -a '(_example_callback _example__multiparts##at)' -r
complete -c 'example' -f -n '_example_state _example__multiparts' -l 'colon' -d 'multiparts with \: as divider ' -a '(_example_callback _example__multiparts##colon)' -r
complete -c 'example' -f -n '_example_state _example__multiparts' -l 'comma' -d 'multiparts with , as divider' -a '(_example_callback _example__multiparts##comma)' -r
complete -c 'example' -f -n '_example_state _example__multiparts' -l 'dot' -d 'multiparts with . as divider' -a '(_example_callback _example__multiparts##dot)' -r
complete -c 'example' -f -n '_example_state _example__multiparts' -l 'dotdotdot' -d 'multiparts with ... as divider' -a '(_example_callback _example__multiparts##dotdotdot)' -r
complete -c 'example' -f -n '_example_state _example__multiparts' -l 'equals' -d 'multiparts with = as divider' -a '(_example_callback _example__multiparts##equals)' -r
complete -c 'example' -f -n '_example_state _example__multiparts' -l 'none' -d 'multiparts without divider' -a '(_example_callback _example__multiparts##none)' -r
complete -c 'example' -f -n '_example_state _example__multiparts' -l 'slash' -d 'multiparts with / as divider' -a '(_example_callback _example__multiparts##slash)' -r
complete -c 'example' -f -n '_example_state _example__multiparts' -a '(_example_callback _)'
`
	rootCmd.InitDefaultHelpCmd()
	s, _ := carapace.Gen(rootCmd).Snippet("fish", false)
	assert.Equal(t, expected, s)
}

func TestPowershell(t *testing.T) {
	expected := `using namespace System.Management.Automation
using namespace System.Management.Automation.Language
$_example_completer = {
    param($wordToComplete, $commandAst, $cursorPosition)
    $commandElements = $commandAst.CommandElements
    $previous = $commandElements[-1].Extent
    if ($wordToComplete) {
        $previous = $commandElements[-2].Extent
    }

    $state = example _carapace powershell state $($commandElements| Foreach {$_.Extent})
    
    Function _example_callback {
      param($uid)
      if (!$wordToComplete) {
        example _carapace powershell "$uid" $($commandElements| Foreach {$_.Extent}) '""' | Out-String | Invoke-Expression
      } else {
        example _carapace powershell "$uid" $($commandElements| Foreach {$_.Extent}) | Out-String | Invoke-Expression
      }
    }
    
    $completions = @(switch ($state) {
        '_example' {
            switch -regex ($previous) {
                '^(-a|--array)$' {
                         
                        break
                      }
                default {
                    switch -regex ($wordToComplete) {
                '^(-p=*|--persistentFlag=*)$' {
                        @(
                        
                        ) | ForEach-Object{ [CompletionResult]::new($wordToComplete.split("=")[0] + "=" + $_.CompletionText, $_.ListItemText, $_.ResultType, $_.ToolTip) }
                        break
                      }

                        default {

            if ($wordToComplete -like "-*") {
                [CompletionResult]::new('-a ', '-a', [CompletionResultType]::ParameterName, 'multiflag')
                [CompletionResult]::new('--array ', '--array', [CompletionResultType]::ParameterName, 'multiflag')
                [CompletionResult]::new('-p ', '-p', [CompletionResultType]::ParameterName, 'Help message for persistentFlag')
                [CompletionResult]::new('--persistentFlag ', '--persistentFlag', [CompletionResultType]::ParameterName, 'Help message for persistentFlag')
                [CompletionResult]::new('-t ', '-t', [CompletionResultType]::ParameterName, 'Help message for toggle')
                [CompletionResult]::new('--toggle ', '--toggle', [CompletionResultType]::ParameterName, 'Help message for toggle')
            } else {
                [CompletionResult]::new('action ', 'action', [CompletionResultType]::Command, 'action example')
                [CompletionResult]::new('callback ', 'callback', [CompletionResultType]::Command, 'callback example')
                [CompletionResult]::new('condition ', 'condition', [CompletionResultType]::Command, 'condition example')
                [CompletionResult]::new('help ', 'help', [CompletionResultType]::Command, 'Help about any command')
                [CompletionResult]::new('injection ', 'injection', [CompletionResultType]::Command, 'just trying to break things')
                [CompletionResult]::new('multiparts ', 'multiparts', [CompletionResultType]::Command, 'multiparts example')
            }
            break
        }
                        }
                    }
                }
            }

        '_example__action' {
            switch -regex ($previous) {
                '^(--directories)$' {
                        [CompletionResult]::new('', '', [CompletionResultType]::ParameterValue, '') 
                        break
                      }
                '^(-f|--files)$' {
                        [CompletionResult]::new('', '', [CompletionResultType]::ParameterValue, '') 
                        break
                      }
                '^(-g|--groups)$' {
                        _example_callback '_example__action##groups' 
                        break
                      }
                '^(--hosts)$' {
                        _example_callback '_example__action##hosts' 
                        break
                      }
                '^(-k|--kill)$' {
                        [CompletionResult]::new('ABRT', 'ABRT ', [CompletionResultType]::ParameterValue, 'Abnormal termination ')
                        [CompletionResult]::new('ALRM', 'ALRM ', [CompletionResultType]::ParameterValue, 'Virtual alarm clock ')
                        [CompletionResult]::new('BUS', 'BUS ', [CompletionResultType]::ParameterValue, 'BUS error ')
                        [CompletionResult]::new('CHLD', 'CHLD ', [CompletionResultType]::ParameterValue, 'Child status has changed ')
                        [CompletionResult]::new('CONT', 'CONT ', [CompletionResultType]::ParameterValue, 'Continue stopped process ')
                        [CompletionResult]::new('FPE', 'FPE ', [CompletionResultType]::ParameterValue, 'Floating-point exception ')
                        [CompletionResult]::new('HUP', 'HUP ', [CompletionResultType]::ParameterValue, 'Hangup detected on controlling terminal ')
                        [CompletionResult]::new('ILL', 'ILL ', [CompletionResultType]::ParameterValue, 'Illegal instruction ')
                        [CompletionResult]::new('INT', 'INT ', [CompletionResultType]::ParameterValue, 'Interrupt from keyboard ')
                        [CompletionResult]::new('KILL', 'KILL ', [CompletionResultType]::ParameterValue, 'Kill, unblockable ')
                        [CompletionResult]::new('PIPE', 'PIPE ', [CompletionResultType]::ParameterValue, 'Broken pipe ')
                        [CompletionResult]::new('POLL', 'POLL ', [CompletionResultType]::ParameterValue, 'Pollable event occurred ')
                        [CompletionResult]::new('PROF', 'PROF ', [CompletionResultType]::ParameterValue, 'Profiling alarm clock timer expired ')
                        [CompletionResult]::new('PWR', 'PWR ', [CompletionResultType]::ParameterValue, 'Power failure restart ')
                        [CompletionResult]::new('QUIT', 'QUIT ', [CompletionResultType]::ParameterValue, 'Quit from keyboard ')
                        [CompletionResult]::new('SEGV', 'SEGV ', [CompletionResultType]::ParameterValue, 'Segmentation violation ')
                        [CompletionResult]::new('STKFLT', 'STKFLT ', [CompletionResultType]::ParameterValue, 'Stack fault on coprocessor ')
                        [CompletionResult]::new('STOP', 'STOP ', [CompletionResultType]::ParameterValue, 'Stop process, unblockable ')
                        [CompletionResult]::new('SYS', 'SYS ', [CompletionResultType]::ParameterValue, 'Bad system call ')
                        [CompletionResult]::new('TERM', 'TERM ', [CompletionResultType]::ParameterValue, 'Termination request ')
                        [CompletionResult]::new('TRAP', 'TRAP ', [CompletionResultType]::ParameterValue, 'Trace/breakpoint trap ')
                        [CompletionResult]::new('TSTP', 'TSTP ', [CompletionResultType]::ParameterValue, 'Stop typed at keyboard ')
                        [CompletionResult]::new('TTIN', 'TTIN ', [CompletionResultType]::ParameterValue, 'Background read from tty ')
                        [CompletionResult]::new('TTOU', 'TTOU ', [CompletionResultType]::ParameterValue, 'Background write to tty ')
                        [CompletionResult]::new('URG', 'URG ', [CompletionResultType]::ParameterValue, 'Urgent condition on socket ')
                        [CompletionResult]::new('USR1', 'USR1 ', [CompletionResultType]::ParameterValue, 'User-defined signal 1 ')
                        [CompletionResult]::new('USR2', 'USR2 ', [CompletionResultType]::ParameterValue, 'User-defined signal 2 ')
                        [CompletionResult]::new('VTALRM', 'VTALRM ', [CompletionResultType]::ParameterValue, 'Virtual alarm clock ')
                        [CompletionResult]::new('WINCH', 'WINCH ', [CompletionResultType]::ParameterValue, 'Window size change ')
                        [CompletionResult]::new('XCPU', 'XCPU ', [CompletionResultType]::ParameterValue, 'CPU time limit exceeded ')
                        [CompletionResult]::new('XFSZ', 'XFSZ ', [CompletionResultType]::ParameterValue, 'File size limit exceeded ') 
                        break
                      }
                '^(-m|--message)$' {
                        [CompletionResult]::new('_', '_ ', [CompletionResultType]::ParameterValue, ' ')
                        [CompletionResult]::new('ERR', 'ERR ', [CompletionResultType]::ParameterValue, 'message example ') 
                        break
                      }
                '^(-n|--net_interfaces)$' {
                        _example_callback '_example__action##net_interfaces' 
                        break
                      }
                '^(--usergroup)$' {
                        _example_callback '_example__action##usergroup' 
                        break
                      }
                '^(-u|--users)$' {
                        _example_callback '_example__action##users' 
                        break
                      }
                '^(-v|--values)$' {
                        [CompletionResult]::new('values', 'values ', [CompletionResultType]::ParameterValue, ' ')
                        [CompletionResult]::new('example', 'example ', [CompletionResultType]::ParameterValue, ' ') 
                        break
                      }
                '^(-d|--values_described)$' {
                        [CompletionResult]::new('values', 'values ', [CompletionResultType]::ParameterValue, 'valueDescription ')
                        [CompletionResult]::new('example', 'example ', [CompletionResultType]::ParameterValue, 'exampleDescription ') 
                        break
                      }
                default {
                    switch -regex ($wordToComplete) {
                '^(-o=*|--optarg=*)$' {
                        @(
                        [CompletionResult]::new('blue', 'blue ', [CompletionResultType]::ParameterValue, ' ')
                        [CompletionResult]::new('red', 'red ', [CompletionResultType]::ParameterValue, ' ')
                        [CompletionResult]::new('green', 'green ', [CompletionResultType]::ParameterValue, ' ')
                        [CompletionResult]::new('yellow', 'yellow ', [CompletionResultType]::ParameterValue, ' ')
                        ) | ForEach-Object{ [CompletionResult]::new($wordToComplete.split("=")[0] + "=" + $_.CompletionText, $_.ListItemText, $_.ResultType, $_.ToolTip) }
                        break
                      }

                        default {

            if ($wordToComplete -like "-*") {
                [CompletionResult]::new('--directories ', '--directories', [CompletionResultType]::ParameterName, 'files flag')
                [CompletionResult]::new('-f ', '-f', [CompletionResultType]::ParameterName, 'files flag')
                [CompletionResult]::new('--files ', '--files', [CompletionResultType]::ParameterName, 'files flag')
                [CompletionResult]::new('-g ', '-g', [CompletionResultType]::ParameterName, 'groups flag')
                [CompletionResult]::new('--groups ', '--groups', [CompletionResultType]::ParameterName, 'groups flag')
                [CompletionResult]::new('--hosts ', '--hosts', [CompletionResultType]::ParameterName, 'hosts flag')
                [CompletionResult]::new('-k ', '-k', [CompletionResultType]::ParameterName, 'kill signals')
                [CompletionResult]::new('--kill ', '--kill', [CompletionResultType]::ParameterName, 'kill signals')
                [CompletionResult]::new('-m ', '-m', [CompletionResultType]::ParameterName, 'message flag')
                [CompletionResult]::new('--message ', '--message', [CompletionResultType]::ParameterName, 'message flag')
                [CompletionResult]::new('-n ', '-n', [CompletionResultType]::ParameterName, 'net_interfaces flag')
                [CompletionResult]::new('--net_interfaces ', '--net_interfaces', [CompletionResultType]::ParameterName, 'net_interfaces flag')
                [CompletionResult]::new('-o ', '-o', [CompletionResultType]::ParameterName, 'optional arg with default value blue')
                [CompletionResult]::new('--optarg ', '--optarg', [CompletionResultType]::ParameterName, 'optional arg with default value blue')
                [CompletionResult]::new('--usergroup ', '--usergroup', [CompletionResultType]::ParameterName, 'user:group flag')
                [CompletionResult]::new('-u ', '-u', [CompletionResultType]::ParameterName, 'users flag')
                [CompletionResult]::new('--users ', '--users', [CompletionResultType]::ParameterName, 'users flag')
                [CompletionResult]::new('-v ', '-v', [CompletionResultType]::ParameterName, 'values flag')
                [CompletionResult]::new('--values ', '--values', [CompletionResultType]::ParameterName, 'values flag')
                [CompletionResult]::new('-d ', '-d', [CompletionResultType]::ParameterName, 'values with description flag')
                [CompletionResult]::new('--values_described ', '--values_described', [CompletionResultType]::ParameterName, 'values with description flag')
            } else {
                _example_callback '_'
            }
            break
        }
                        }
                    }
                }
            }

        '_example__callback' {
            switch -regex ($previous) {
                '^(-c|--callback)$' {
                        _example_callback '_example__callback##callback' 
                        break
                      }
                default {
                    switch -regex ($wordToComplete) {


                        default {

            if ($wordToComplete -like "-*") {
                [CompletionResult]::new('-c ', '-c', [CompletionResultType]::ParameterName, 'Help message for callback')
                [CompletionResult]::new('--callback ', '--callback', [CompletionResultType]::ParameterName, 'Help message for callback')
            } else {
                _example_callback '_'
            }
            break
        }
                        }
                    }
                }
            }

        '_example__condition' {
            switch -regex ($previous) {
                '^(-r|--required)$' {
                        [CompletionResult]::new('valid', 'valid ', [CompletionResultType]::ParameterValue, ' ')
                        [CompletionResult]::new('invalid', 'invalid ', [CompletionResultType]::ParameterValue, ' ') 
                        break
                      }
                default {
                    switch -regex ($wordToComplete) {


                        default {

            if ($wordToComplete -like "-*") {
                [CompletionResult]::new('-r ', '-r', [CompletionResultType]::ParameterName, 'required flag')
                [CompletionResult]::new('--required ', '--required', [CompletionResultType]::ParameterName, 'required flag')
            } else {
                _example_callback '_'
            }
            break
        }
                        }
                    }
                }
            }

        '_example__help' {
            switch -regex ($previous) {

                default {
                    switch -regex ($wordToComplete) {


                        default {

            if ($wordToComplete -like "-*") {
            } else {
                _example_callback '_'
            }
            break
        }
                        }
                    }
                }
            }

        '_example__injection' {
            switch -regex ($previous) {

                default {
                    switch -regex ($wordToComplete) {


                        default {

            if ($wordToComplete -like "-*") {
            } else {
                _example_callback '_'
            }
            break
        }
                        }
                    }
                }
            }

        '_example__multiparts' {
            switch -regex ($previous) {
                '^(--at)$' {
                        _example_callback '_example__multiparts##at' 
                        break
                      }
                '^(--colon)$' {
                        _example_callback '_example__multiparts##colon' 
                        break
                      }
                '^(--comma)$' {
                        _example_callback '_example__multiparts##comma' 
                        break
                      }
                '^(--dot)$' {
                        _example_callback '_example__multiparts##dot' 
                        break
                      }
                '^(--dotdotdot)$' {
                        _example_callback '_example__multiparts##dotdotdot' 
                        break
                      }
                '^(--equals)$' {
                        _example_callback '_example__multiparts##equals' 
                        break
                      }
                '^(--none)$' {
                        _example_callback '_example__multiparts##none' 
                        break
                      }
                '^(--slash)$' {
                        _example_callback '_example__multiparts##slash' 
                        break
                      }
                default {
                    switch -regex ($wordToComplete) {


                        default {

            if ($wordToComplete -like "-*") {
                [CompletionResult]::new('--at ', '--at', [CompletionResultType]::ParameterName, 'multiparts with @ as divider')
                [CompletionResult]::new('--colon ', '--colon', [CompletionResultType]::ParameterName, 'multiparts with : as divider ')
                [CompletionResult]::new('--comma ', '--comma', [CompletionResultType]::ParameterName, 'multiparts with , as divider')
                [CompletionResult]::new('--dot ', '--dot', [CompletionResultType]::ParameterName, 'multiparts with . as divider')
                [CompletionResult]::new('--dotdotdot ', '--dotdotdot', [CompletionResultType]::ParameterName, 'multiparts with ... as divider')
                [CompletionResult]::new('--equals ', '--equals', [CompletionResultType]::ParameterName, 'multiparts with = as divider')
                [CompletionResult]::new('--none ', '--none', [CompletionResultType]::ParameterName, 'multiparts without divider')
                [CompletionResult]::new('--slash ', '--slash', [CompletionResultType]::ParameterName, 'multiparts with / as divider')
            } else {
                _example_callback '_'
            }
            break
        }
                        }
                    }
                }
            }

    })

    if ($completions.count -eq 0) {
      return "" # prevent default file completion
    }

    $completions.Where{ $_.CompletionText -like "$wordToComplete*" } |
        Sort-Object -Property ListItemText
}
Register-ArgumentCompleter -Native -CommandName 'example' -ScriptBlock $_example_completer
`
	rootCmd.InitDefaultHelpCmd()
	s, _ := carapace.Gen(rootCmd).Snippet("powershell", false)
	assert.Equal(t, expected, s)
}

func TestXonsh(t *testing.T) {
	expected := `from shlex import split
import re
import pathlib
import subprocess
import xonsh
from xonsh.completers._aliases import _add_one_completer
from xonsh.completers.path import complete_dir, complete_path
from xonsh.completers.tools import RichCompletion

def _example_completer(prefix, line, begidx, endidx, ctx):
    """carapace completer for example"""
    if not line.startswith('example '):
        return # not the expected command to complete
    
    full_words=split(line + "_") # ensure last word is empty when ends with space
    full_words[-1]=full_words[-1][0:-1]
    words=split(line[0:endidx] + "_") # ensure last word is empty when ends with space
    words[-1]=words[-1][0:-1]
    current=words[-1]
    previous=words[-2]
    suffix=full_words[len(words)-1][len(current):]

    result = {}

    # TODO python retrieve state
    state, _ = subprocess.Popen(['example', '_carapace', 'xonsh', 'state', *words],
                                   stdout=subprocess.PIPE,
                                   stderr=subprocess.PIPE).communicate()
    state = state.decode('utf-8').split('\n')[0]
   
    # TODO python callback function
    def _example_callback(uid):
        cb, _ = subprocess.Popen(['example', '_carapace', 'xonsh', uid, *words],
                                     stdout=subprocess.PIPE,
                                     stderr=subprocess.PIPE).communicate()
        cb = cb.decode('utf-8')
        if cb == "":
            return {}
        else:
            nonlocal prefix, line, begidx, endidx, ctx
            return eval(cb)
   
    if False:
        pass
    elif state == '_example':
        if False: # switch previous
            pass
        elif re.search('^(-a|--array)$',previous):
            result = {}
                  
        else:
            if False:
                pass
            elif re.search('^(-p=.*|--persistentFlag=.*)$',current):
                result = {}
                result = set(map(lambda x: RichCompletion(current.split('=')[0]+'='+x, display=x.display, description=x.description, prefix_len=x.prefix_len), result))
    

            elif re.search("-.*",current):
                result = {
                    RichCompletion('-a', display='-a', description='multiflag', prefix_len=0),
                    RichCompletion('--array', display='--array', description='multiflag', prefix_len=0),
                    RichCompletion('-p', display='-p', description='Help message for persistentFlag', prefix_len=0),
                    RichCompletion('--persistentFlag', display='--persistentFlag', description='Help message for persistentFlag', prefix_len=0),
                    RichCompletion('-t', display='-t', description='Help message for toggle', prefix_len=0),
                    RichCompletion('--toggle', display='--toggle', description='Help message for toggle', prefix_len=0),
                }
            else:
                result = {
                RichCompletion('action', display='action', description='action example', prefix_len=0),
                RichCompletion('callback', display='callback', description='callback example', prefix_len=0),
                RichCompletion('condition', display='condition', description='condition example', prefix_len=0),
                RichCompletion('help', display='help', description='Help about any command', prefix_len=0),
                RichCompletion('injection', display='injection', description='just trying to break things', prefix_len=0),
                RichCompletion('multiparts', display='multiparts', description='multiparts example', prefix_len=0),
                }


    elif state == '_example__action':
        if False: # switch previous
            pass
        elif re.search('^(--directories)$',previous):
            result = { RichCompletion(f, display=pathlib.PurePath(f).name, description='', prefix_len=0) for f in complete_dir(prefix, line, begidx, endidx, ctx, True)[0]}
                  
        elif re.search('^(-f|--files)$',previous):
            result = { RichCompletion(f, display=pathlib.PurePath(f).name, description='', prefix_len=0) for f in complete_path(prefix, line, begidx, endidx, ctx)[0]}
                  
        elif re.search('^(-g|--groups)$',previous):
            result = _example_callback('_example__action##groups')
                  
        elif re.search('^(--hosts)$',previous):
            result = _example_callback('_example__action##hosts')
                  
        elif re.search('^(-k|--kill)$',previous):
            result = {
                          RichCompletion('ABRT', display='ABRT', description='Abnormal termination', prefix_len=0),
                          RichCompletion('ALRM', display='ALRM', description='Virtual alarm clock', prefix_len=0),
                          RichCompletion('BUS', display='BUS', description='BUS error', prefix_len=0),
                          RichCompletion('CHLD', display='CHLD', description='Child status has changed', prefix_len=0),
                          RichCompletion('CONT', display='CONT', description='Continue stopped process', prefix_len=0),
                          RichCompletion('FPE', display='FPE', description='Floating-point exception', prefix_len=0),
                          RichCompletion('HUP', display='HUP', description='Hangup detected on controlling terminal', prefix_len=0),
                          RichCompletion('ILL', display='ILL', description='Illegal instruction', prefix_len=0),
                          RichCompletion('INT', display='INT', description='Interrupt from keyboard', prefix_len=0),
                          RichCompletion('KILL', display='KILL', description='Kill, unblockable', prefix_len=0),
                          RichCompletion('PIPE', display='PIPE', description='Broken pipe', prefix_len=0),
                          RichCompletion('POLL', display='POLL', description='Pollable event occurred', prefix_len=0),
                          RichCompletion('PROF', display='PROF', description='Profiling alarm clock timer expired', prefix_len=0),
                          RichCompletion('PWR', display='PWR', description='Power failure restart', prefix_len=0),
                          RichCompletion('QUIT', display='QUIT', description='Quit from keyboard', prefix_len=0),
                          RichCompletion('SEGV', display='SEGV', description='Segmentation violation', prefix_len=0),
                          RichCompletion('STKFLT', display='STKFLT', description='Stack fault on coprocessor', prefix_len=0),
                          RichCompletion('STOP', display='STOP', description='Stop process, unblockable', prefix_len=0),
                          RichCompletion('SYS', display='SYS', description='Bad system call', prefix_len=0),
                          RichCompletion('TERM', display='TERM', description='Termination request', prefix_len=0),
                          RichCompletion('TRAP', display='TRAP', description='Trace/breakpoint trap', prefix_len=0),
                          RichCompletion('TSTP', display='TSTP', description='Stop typed at keyboard', prefix_len=0),
                          RichCompletion('TTIN', display='TTIN', description='Background read from tty', prefix_len=0),
                          RichCompletion('TTOU', display='TTOU', description='Background write to tty', prefix_len=0),
                          RichCompletion('URG', display='URG', description='Urgent condition on socket', prefix_len=0),
                          RichCompletion('USR1', display='USR1', description='User-defined signal 1', prefix_len=0),
                          RichCompletion('USR2', display='USR2', description='User-defined signal 2', prefix_len=0),
                          RichCompletion('VTALRM', display='VTALRM', description='Virtual alarm clock', prefix_len=0),
                          RichCompletion('WINCH', display='WINCH', description='Window size change', prefix_len=0),
                          RichCompletion('XCPU', display='XCPU', description='CPU time limit exceeded', prefix_len=0),
                          RichCompletion('XFSZ', display='XFSZ', description='File size limit exceeded', prefix_len=0),
                        }
                  
        elif re.search('^(-m|--message)$',previous):
            result = {
                          RichCompletion('_', display='_', description='', prefix_len=0),
                          RichCompletion('ERR', display='ERR', description='message example', prefix_len=0),
                        }
                  
        elif re.search('^(-n|--net_interfaces)$',previous):
            result = _example_callback('_example__action##net_interfaces')
                  
        elif re.search('^(--usergroup)$',previous):
            result = _example_callback('_example__action##usergroup')
                  
        elif re.search('^(-u|--users)$',previous):
            result = _example_callback('_example__action##users')
                  
        elif re.search('^(-v|--values)$',previous):
            result = {
                          RichCompletion('values', display='values', description='', prefix_len=0),
                          RichCompletion('example', display='example', description='', prefix_len=0),
                        }
                  
        elif re.search('^(-d|--values_described)$',previous):
            result = {
                          RichCompletion('values', display='values', description='valueDescription', prefix_len=0),
                          RichCompletion('example', display='example', description='exampleDescription', prefix_len=0),
                        }
                  
        else:
            if False:
                pass
            elif re.search('^(-o=.*|--optarg=.*)$',current):
                result = {
                              RichCompletion('blue', display='blue', description='', prefix_len=0),
                              RichCompletion('red', display='red', description='', prefix_len=0),
                              RichCompletion('green', display='green', description='', prefix_len=0),
                              RichCompletion('yellow', display='yellow', description='', prefix_len=0),
                            }
                result = set(map(lambda x: RichCompletion(current.split('=')[0]+'='+x, display=x.display, description=x.description, prefix_len=x.prefix_len), result))
    

            elif re.search("-.*",current):
                result = {
                    RichCompletion('--directories', display='--directories', description='files flag', prefix_len=0),
                    RichCompletion('-f', display='-f', description='files flag', prefix_len=0),
                    RichCompletion('--files', display='--files', description='files flag', prefix_len=0),
                    RichCompletion('-g', display='-g', description='groups flag', prefix_len=0),
                    RichCompletion('--groups', display='--groups', description='groups flag', prefix_len=0),
                    RichCompletion('--hosts', display='--hosts', description='hosts flag', prefix_len=0),
                    RichCompletion('-k', display='-k', description='kill signals', prefix_len=0),
                    RichCompletion('--kill', display='--kill', description='kill signals', prefix_len=0),
                    RichCompletion('-m', display='-m', description='message flag', prefix_len=0),
                    RichCompletion('--message', display='--message', description='message flag', prefix_len=0),
                    RichCompletion('-n', display='-n', description='net_interfaces flag', prefix_len=0),
                    RichCompletion('--net_interfaces', display='--net_interfaces', description='net_interfaces flag', prefix_len=0),
                    RichCompletion('-o', display='-o', description='optional arg with default value blue', prefix_len=0),
                    RichCompletion('--optarg', display='--optarg', description='optional arg with default value blue', prefix_len=0),
                    RichCompletion('--usergroup', display='--usergroup', description='user:group flag', prefix_len=0),
                    RichCompletion('-u', display='-u', description='users flag', prefix_len=0),
                    RichCompletion('--users', display='--users', description='users flag', prefix_len=0),
                    RichCompletion('-v', display='-v', description='values flag', prefix_len=0),
                    RichCompletion('--values', display='--values', description='values flag', prefix_len=0),
                    RichCompletion('-d', display='-d', description='values with description flag', prefix_len=0),
                    RichCompletion('--values_described', display='--values_described', description='values with description flag', prefix_len=0),
                }
            else:
                result = _example_callback('_')


    elif state == '_example__callback':
        if False: # switch previous
            pass
        elif re.search('^(-c|--callback)$',previous):
            result = _example_callback('_example__callback##callback')
                  
        else:
            if False:
                pass
    

            elif re.search("-.*",current):
                result = {
                    RichCompletion('-c', display='-c', description='Help message for callback', prefix_len=0),
                    RichCompletion('--callback', display='--callback', description='Help message for callback', prefix_len=0),
                }
            else:
                result = _example_callback('_')


    elif state == '_example__condition':
        if False: # switch previous
            pass
        elif re.search('^(-r|--required)$',previous):
            result = {
                          RichCompletion('valid', display='valid', description='', prefix_len=0),
                          RichCompletion('invalid', display='invalid', description='', prefix_len=0),
                        }
                  
        else:
            if False:
                pass
    

            elif re.search("-.*",current):
                result = {
                    RichCompletion('-r', display='-r', description='required flag', prefix_len=0),
                    RichCompletion('--required', display='--required', description='required flag', prefix_len=0),
                }
            else:
                result = _example_callback('_')


    elif state == '_example__help':
        if False: # switch previous
            pass

        else:
            if False:
                pass
    

            elif re.search("-.*",current):
                result = {
                }
            else:
                result = _example_callback('_')


    elif state == '_example__injection':
        if False: # switch previous
            pass

        else:
            if False:
                pass
    

            elif re.search("-.*",current):
                result = {
                }
            else:
                result = _example_callback('_')


    elif state == '_example__multiparts':
        if False: # switch previous
            pass
        elif re.search('^(--at)$',previous):
            result = _example_callback('_example__multiparts##at')
                  
        elif re.search('^(--colon)$',previous):
            result = _example_callback('_example__multiparts##colon')
                  
        elif re.search('^(--comma)$',previous):
            result = _example_callback('_example__multiparts##comma')
                  
        elif re.search('^(--dot)$',previous):
            result = _example_callback('_example__multiparts##dot')
                  
        elif re.search('^(--dotdotdot)$',previous):
            result = _example_callback('_example__multiparts##dotdotdot')
                  
        elif re.search('^(--equals)$',previous):
            result = _example_callback('_example__multiparts##equals')
                  
        elif re.search('^(--none)$',previous):
            result = _example_callback('_example__multiparts##none')
                  
        elif re.search('^(--slash)$',previous):
            result = _example_callback('_example__multiparts##slash')
                  
        else:
            if False:
                pass
    

            elif re.search("-.*",current):
                result = {
                    RichCompletion('--at', display='--at', description='multiparts with @ as divider', prefix_len=0),
                    RichCompletion('--colon', display='--colon', description='multiparts with : as divider ', prefix_len=0),
                    RichCompletion('--comma', display='--comma', description='multiparts with , as divider', prefix_len=0),
                    RichCompletion('--dot', display='--dot', description='multiparts with . as divider', prefix_len=0),
                    RichCompletion('--dotdotdot', display='--dotdotdot', description='multiparts with ... as divider', prefix_len=0),
                    RichCompletion('--equals', display='--equals', description='multiparts with = as divider', prefix_len=0),
                    RichCompletion('--none', display='--none', description='multiparts without divider', prefix_len=0),
                    RichCompletion('--slash', display='--slash', description='multiparts with / as divider', prefix_len=0),
                }
            else:
                result = _example_callback('_')



    result = set(filter(lambda x: x.startswith(current) and x.endswith(suffix), result))
    if len(result) == 0:
        result = {RichCompletion(current, display=current, description='', prefix_len=0)}

    result = set(map(lambda x: RichCompletion(x[:len(x)-len(suffix)], display=x.display, description=x.description, prefix_len=len(current)), result))
    return result
_add_one_completer('example', _example_completer, 'start')
`

	rootCmd.InitDefaultHelpCmd()
	s, _ := carapace.Gen(rootCmd).Snippet("xonsh", false)
	assert.Equal(t, expected, s)
}

func TestZsh(t *testing.T) {
	expected := `#compdef example
function _example_callback {
  # shellcheck disable=SC2086
  eval "$(example _carapace zsh "$1" ${os_args})"
}
function _example {
  local -a commands
  # shellcheck disable=SC2206
  local -a -x os_args=(${words})

  _arguments -C \
    "(*-a *--array)"{\*-a,\*--array}"[multiflag]: :" \
    "(-p --persistentFlag)"{-p=-,--persistentFlag=-}"[Help message for persistentFlag]::" \
    "(-t --toggle)"{-t=-,--toggle=-}"[Help message for toggle]:: :{local _comp_desc=('true' 'false');compadd -S '' -d _comp_desc 'true' 'false'}" \
    "1: :->cmnds" \
    "*::arg:->args"

  # shellcheck disable=SC2154
  case $state in
    cmnds)
      # shellcheck disable=SC2034
      commands=(
        "action:action example"
        "alias:action example"
        "callback:callback example"
        "condition:condition example"
        "help:Help about any command"
        "injection:just trying to break things"
        "multiparts:multiparts example"
      )
      _describe "command" commands
      ;;
  esac
  
  case "${words[1]}" in
    action)
      _example__action
      ;;
    alias)
      _example__action
      ;;
    callback)
      _example__callback
      ;;
    condition)
      _example__condition
      ;;
    help)
      _example__help
      ;;
    injection)
      _example__injection
      ;;
    multiparts)
      _example__multiparts
      ;;
  esac
}

function _example__action {
    _arguments -C \
    "--directories[files flag]: :_files -/" \
    "(-f --files)"{-f,--files}"[files flag]: :_files -g '*.go'" \
    "(-g --groups)"{-g,--groups}"[groups flag]: :{_example_callback '_example__action##groups'}" \
    "--hosts[hosts flag]: :{_example_callback '_example__action##hosts'}" \
    "(-k --kill)"{-k,--kill}"[kill signals]: :{local _comp_desc=('ABRT (Abnormal termination)' 'ALRM (Virtual alarm clock)' 'BUS (BUS error)' 'CHLD (Child status has changed)' 'CONT (Continue stopped process)' 'FPE (Floating-point exception)' 'HUP (Hangup detected on controlling terminal)' 'ILL (Illegal instruction)' 'INT (Interrupt from keyboard)' 'KILL (Kill, unblockable)' 'PIPE (Broken pipe)' 'POLL (Pollable event occurred)' 'PROF (Profiling alarm clock timer expired)' 'PWR (Power failure restart)' 'QUIT (Quit from keyboard)' 'SEGV (Segmentation violation)' 'STKFLT (Stack fault on coprocessor)' 'STOP (Stop process, unblockable)' 'SYS (Bad system call)' 'TERM (Termination request)' 'TRAP (Trace/breakpoint trap)' 'TSTP (Stop typed at keyboard)' 'TTIN (Background read from tty)' 'TTOU (Background write to tty)' 'URG (Urgent condition on socket)' 'USR1 (User-defined signal 1)' 'USR2 (User-defined signal 2)' 'VTALRM (Virtual alarm clock)' 'WINCH (Window size change)' 'XCPU (CPU time limit exceeded)' 'XFSZ (File size limit exceeded)');compadd -S '' -d _comp_desc 'ABRT' 'ALRM' 'BUS' 'CHLD' 'CONT' 'FPE' 'HUP' 'ILL' 'INT' 'KILL' 'PIPE' 'POLL' 'PROF' 'PWR' 'QUIT' 'SEGV' 'STKFLT' 'STOP' 'SYS' 'TERM' 'TRAP' 'TSTP' 'TTIN' 'TTOU' 'URG' 'USR1' 'USR2' 'VTALRM' 'WINCH' 'XCPU' 'XFSZ'}" \
    "(-m --message)"{-m,--message}"[message flag]: :{local _comp_desc=('_' 'ERR (message example)');compadd -S '' -d _comp_desc '_' 'ERR'}" \
    "(-n --net_interfaces)"{-n,--net_interfaces}"[net_interfaces flag]: :{_example_callback '_example__action##net_interfaces'}" \
    "(-o --optarg)"{-o=-,--optarg=-}"[optional arg with default value blue]:: :{local _comp_desc=('blue' 'red' 'green' 'yellow');compadd -S '' -d _comp_desc 'blue' 'red' 'green' 'yellow'}" \
    "--usergroup[user\:group flag]: :{_example_callback '_example__action##usergroup'}" \
    "(-u --users)"{-u,--users}"[users flag]: :{_example_callback '_example__action##users'}" \
    "(-v --values)"{-v,--values}"[values flag]: :{local _comp_desc=('values' 'example');compadd -S '' -d _comp_desc 'values' 'example'}" \
    "(-d --values_described)"{-d,--values_described}"[values with description flag]: :{local _comp_desc=('values (valueDescription)' 'example (exampleDescription)');compadd -S '' -d _comp_desc 'values' 'example'}" \
    "1: :{local _comp_desc=('positional1' 'p1');compadd -S '' -d _comp_desc 'positional1' 'p1'}" \
    "2: :{local _comp_desc=('positional2' 'p2');compadd -S '' -d _comp_desc 'positional2' 'p2'}"
}

function _example__callback {
    _arguments -C \
    "(-c --callback)"{-c,--callback}"[Help message for callback]: :{_example_callback '_example__callback##callback'}" \
    "1: :{_example_callback '_example__callback#1'}" \
    "2: :{_example_callback '_example__callback#2'}" \
    "*: :{_example_callback '_example__callback#0'}"
}

function _example__condition {
    _arguments -C \
    "(-r --required)"{-r,--required}"[required flag]: :{local _comp_desc=('valid' 'invalid');compadd -S '' -d _comp_desc 'valid' 'invalid'}" \
    "1: :{_example_callback '_example__condition#1'}"
}

function _example__help {
    _arguments -C \
    "*::arg:->args"
}

function _example__injection {
    _arguments -C \
    "1: :{local _comp_desc=('echo fail');compadd -S '' -d _comp_desc 'echo fail'}" \
    "2: :{local _comp_desc=('echo fail');compadd -S '' -d _comp_desc 'echo fail'}" \
    "3: :{local _comp_desc=('echo fail');compadd -S '' -d _comp_desc 'echo fail'}" \
    "4: :{local _comp_desc=(' echo fail ');compadd -S '' -d _comp_desc ' echo fail '}" \
    "5: :{local _comp_desc=(' echo fail ');compadd -S '' -d _comp_desc ' echo fail '}" \
    "6: :{local _comp_desc=(' echo fail ');compadd -S '' -d _comp_desc ' echo fail '}" \
    "7: :{local _comp_desc=('echo fail');compadd -S '' -d _comp_desc 'echo fail'}" \
    "8: :{local _comp_desc=('');compadd -S '' -d _comp_desc ''}" \
    "9: :{local _comp_desc=('LAST POSITIONAL VALUE');compadd -S '' -d _comp_desc 'LAST POSITIONAL VALUE'}"
}

function _example__multiparts {
    _arguments -C \
    "--at[multiparts with @ as divider]: :{_example_callback '_example__multiparts##at'}" \
    "--colon[multiparts with \: as divider ]: :{_example_callback '_example__multiparts##colon'}" \
    "--comma[multiparts with , as divider]: :{_example_callback '_example__multiparts##comma'}" \
    "--dot[multiparts with . as divider]: :{_example_callback '_example__multiparts##dot'}" \
    "--dotdotdot[multiparts with ... as divider]: :{_example_callback '_example__multiparts##dotdotdot'}" \
    "--equals[multiparts with = as divider]: :{_example_callback '_example__multiparts##equals'}" \
    "--none[multiparts without divider]: :{_example_callback '_example__multiparts##none'}" \
    "--slash[multiparts with / as divider]: :{_example_callback '_example__multiparts##slash'}" \
    "1: :{_example_callback '_example__multiparts#1'}" \
    "*: :{_example_callback '_example__multiparts#0'}"
}
compquote '' 2>/dev/null && _example
compdef _example example
`
	rootCmd.InitDefaultHelpCmd()
	s, _ := carapace.Gen(rootCmd).Snippet("zsh", false)
	assert.Equal(t, expected, s)
}
