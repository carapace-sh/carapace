package cmd

import (
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/assert"
)

func TestBash(t *testing.T) {
	expected := `#!/bin/bash
_example_callback() {
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local last="${COMP_WORDS[${COMP_CWORD}]}"
  if [[ $last =~ ^[\"\'] ]] && ! echo "$last" | xargs echo 2>/dev/null >/dev/null ; then
      compline="${compline}${last:0:1}"
      last="${last// /\\\\ }" 
  fi

  echo "$compline" | sed -e 's/ $/ _/' -e 's/"/\"/g' | xargs example _carapace bash "$1"
}

_example_completions() {
  local compline="${COMP_LINE:0:${COMP_POINT}}"
  local last="${COMP_WORDS[${COMP_CWORD}]}"
  
  if [[ $last =~ ^[\"\'] ]] && ! echo "$last" | xargs echo 2>/dev/null >/dev/null ; then
      compline="${compline}${last:0:1}"
      last="${last// /\\\\ }" 
  else
      last="${last// /\\\ }" 
  fi

  local state
  state="$(echo "$compline" | sed -e "s/ \$/ _/" -e 's/"/\"/g' | xargs example _carapace bash state)"
  local previous="${COMP_WORDS[$((COMP_CWORD-1))]}"
  local IFS=$'\n'

  case $state in

    '_example' )
      if [[ $last == -* ]]; then
        COMPREPLY=($(compgen -W $'--array\n-a\n--persistentFlag\n-p\n--toggle\n-t' -- "$last"))
      else
        case $previous in
          -a | --array)
            COMPREPLY=($())
            ;;



          *)
            COMPREPLY=($(compgen -W $'action\nalias\ncallback\ncondition\nhelp\ninjection' -- "$last"))
            ;;
        esac
      fi
      ;;


    '_example__action' )
      if [[ $last == -* ]]; then
        COMPREPLY=($(compgen -W $'--custom\n-c\n--directories\n--files\n-f\n--groups\n-g\n--hosts\n--message\n-m\n--multi_parts\n--net_interfaces\n-n\n--users\n-u\n--values\n-v\n--values_described\n-d' -- "$last"))
      else
        case $previous in
          -c | --custom)
            COMPREPLY=($())
            ;;

          --directories)
            COMPREPLY=($(compgen -S / -d -- "$last"))
            ;;

          -f | --files)
            COMPREPLY=($(compgen -S / -d -- "$last"; compgen -f -X '!*.go' -- "$last"))
            ;;

          -g | --groups)
            COMPREPLY=($(compgen -g -- "${last//[\"\|\']/}"))
            ;;

          --hosts)
            COMPREPLY=($(compgen -W "$(cut -d ' ' -f1 < ~/.ssh/known_hosts | cut -d ',' -f1)" -- "$last"))
            ;;

          -m | --message)
            COMPREPLY=($(compgen -W $'ERR\nmessage\\\ example' -- "$last"))
            ;;

          --multi_parts)
            COMPREPLY=($(compgen -W $'multi/parts\nmulti/parts/example\nmulti/parts/test\nexample/parts' -- "$last"))
            ;;

          -n | --net_interfaces)
            COMPREPLY=($(compgen -W "$(ifconfig -a | grep -o '^[^ :]\+')" -- "$last"))
            ;;

          -u | --users)
            COMPREPLY=($(compgen -u -- "${last//[\"\|\']/}"))
            ;;

          -v | --values)
            COMPREPLY=($(compgen -W $'values\nexample' -- "$last"))
            ;;

          -d | --values_described)
            COMPREPLY=($(compgen -W $'values\nexample\n\n' -- "$last"))
            ;;

          *)
            COMPREPLY=($(eval $(_example_callback '_')))
            ;;
        esac
      fi
      ;;


    '_example__callback' )
      if [[ $last == -* ]]; then
        COMPREPLY=($(compgen -W $'--callback\n-c' -- "$last"))
      else
        case $previous in
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
      if [[ $last == -* ]]; then
        COMPREPLY=($(compgen -W $'--required\n-r' -- "$last"))
      else
        case $previous in
          -r | --required)
            COMPREPLY=($(compgen -W $'valid\ninvalid' -- "$last"))
            ;;

          *)
            COMPREPLY=($(eval $(_example_callback '_')))
            ;;
        esac
      fi
      ;;


    '_example__help' )
      if [[ $last == -* ]]; then
        COMPREPLY=($())
      else
        case $previous in

          *)
            COMPREPLY=($(eval $(_example_callback '_')))
            ;;
        esac
      fi
      ;;


    '_example__injection' )
      if [[ $last == -* ]]; then
        COMPREPLY=($())
      else
        case $previous in

          *)
            COMPREPLY=($(eval $(_example_callback '_')))
            ;;
        esac
      fi
      ;;

  esac

  [[ $last =~ ^[\"\'] ]] && COMPREPLY=("${COMPREPLY[@]//\\ /\ }")
  [[ ${COMPREPLY[0]} == */ ]] && compopt -o nospace
}

complete -F _example_completions example
`
	rootCmd.InitDefaultHelpCmd()
	assert.Equal(t, expected, carapace.Gen(rootCmd).Bash())
}

func TestElvish(t *testing.T) {
	expected := `edit:completion:arg-completer[example] = [@arg]{
  fn _example_callback [uid]{
    # TODO there is no 'eval' in elvish and '-source' needs a file so use a tempary one for callback 
    tmpfile=(mktemp -t carapace_example_callback-XXXXX.elv)
    echo (joins ' ' $arg) | xargs example _carapace elvish $uid > $tmpfile
    -source $tmpfile
    rm $tmpfile
  }

  fn subindex [subcommand]{
    # TODO 'edit:complete-getopt' needs the arguments shortened for subcommmands - pretty optimistic here
    index=1
    for x $arg { if (eq $x $subcommand) { break } else { index = (+ $index 1) } } 
    echo $index
  }
  
  state=(echo (joins ' ' $arg) | xargs example _carapace elvish state)
  if (eq 1 0) {
  }  elif (eq $state '_example') {
    opt-specs = [
        [&long=array &desc='multiflag' &short=a &arg-required=$true &completer=[_]{  }]
        [&long=persistentFlag &desc='Help message for persistentFlag' &short=p]
        [&long=toggle &desc='Help message for toggle' &short=t]
    ]
    arg-handlers = [
        [_]{ edit:complex-candidate action &display-suffix=' (action example)'
edit:complex-candidate alias &display-suffix=' (action example)'
edit:complex-candidate callback &display-suffix=' (callback example)'
edit:complex-candidate condition &display-suffix=' (condition example)'
edit:complex-candidate help &display-suffix=' (Help about any command)'
edit:complex-candidate injection &display-suffix=' (just trying to break things)'





 }
    ]
    subargs = $arg[(subindex example):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }  elif (eq $state '_example__action') {
    opt-specs = [
        [&long=custom &desc='custom flag' &short=c &arg-required=$true &completer=[_]{  }]
        [&long=directories &desc='files flag' &arg-required=$true &completer=[_]{ edit:complete-filename $arg[-1] }]
        [&long=files &desc='files flag' &short=f &arg-required=$true &completer=[_]{ edit:complete-filename $arg[-1] }]
        [&long=groups &desc='groups flag' &short=g &arg-required=$true &completer=[_]{  }]
        [&long=hosts &desc='hosts flag' &arg-required=$true &completer=[_]{  }]
        [&long=message &desc='message flag' &short=m &arg-required=$true &completer=[_]{ edit:complex-candidate ERR &display-suffix=' (message example)'
edit:complex-candidate _ &display-suffix=' ()'

 }]
        [&long=multi_parts &desc='multi_parts flag' &arg-required=$true &completer=[_]{ put multi/parts multi/parts/example multi/parts/test example/parts }]
        [&long=net_interfaces &desc='net_interfaces flag' &short=n &arg-required=$true &completer=[_]{  }]
        [&long=users &desc='users flag' &short=u &arg-required=$true &completer=[_]{  }]
        [&long=values &desc='values flag' &short=v &arg-required=$true &completer=[_]{ put values example }]
        [&long=values_described &desc='values with description flag' &short=d &arg-required=$true &completer=[_]{ edit:complex-candidate values &display-suffix=' (valueDescription)'
edit:complex-candidate example &display-suffix=' (exampleDescription)'

 }]
    ]
    arg-handlers = [
      [_]{ put positional1 p1 }
      [_]{ put positional2 p2 }
    ]
    subargs = $arg[(subindex action):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }  elif (eq $state '_example__callback') {
    opt-specs = [
        [&long=callback &desc='Help message for callback' &short=c &arg-required=$true &completer=[_]{ _example_callback '_example__callback##callback' }]
    ]
    arg-handlers = [
      [_]{ _example_callback '_example__callback#1' }
    ]
    subargs = $arg[(subindex callback):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }  elif (eq $state '_example__condition') {
    opt-specs = [
        [&long=required &desc='required flag' &short=r &arg-required=$true &completer=[_]{ put valid invalid }]
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
      [_]{ put echo fail }
      [_]{ put echo fail }
      [_]{ put echo fail }
      [_]{ put  echo fail  }
      [_]{ put  echo fail  }
      [_]{ put  echo fail  }
      [_]{ put echo fail }
      [_]{ edit:complex-candidate ERR &display-suffix=' (no values to complete)'
edit:complex-candidate _ &display-suffix=' ()'

 }
      [_]{ put LAST POSITIONAL VALUE }
    ]
    subargs = $arg[(subindex injection):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }
}
`
	rootCmd.InitDefaultHelpCmd()
	assert.Equal(t, expected, carapace.Gen(rootCmd).Elvish())
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
  set -lx CALLBACK (commandline -cp | sed "s/\$/"(_example_quote_suffix)"/" | sed "s/ \$/ _/" | xargs example _carapace fish $argv )
  eval "$CALLBACK"
end

complete -c example -f

complete -c example -f -n '_example_state _example' -l array -s a -d 'multiflag' -r
complete -c example -f -n '_example_state _example' -l persistentFlag -s p -d 'Help message for persistentFlag'
complete -c example -f -n '_example_state _example' -l toggle -s t -d 'Help message for toggle' -a '(echo -e "true\nfalse")' -r
complete -c example -f -n '_example_state _example ' -a 'action alias' -d 'action example'
complete -c example -f -n '_example_state _example ' -a 'callback ' -d 'callback example'
complete -c example -f -n '_example_state _example ' -a 'condition ' -d 'condition example'
complete -c example -f -n '_example_state _example ' -a 'help ' -d 'Help about any command'
complete -c example -f -n '_example_state _example ' -a 'injection ' -d 'just trying to break things'


complete -c example -f -n '_example_state _example__action' -l custom -s c -d 'custom flag' -a '()' -r
complete -c example -f -n '_example_state _example__action' -l directories -d 'files flag' -a '(__fish_complete_directories)' -r
complete -c example -f -n '_example_state _example__action' -l files -s f -d 'files flag' -a '(__fish_complete_suffix ".go")' -r
complete -c example -f -n '_example_state _example__action' -l groups -s g -d 'groups flag' -a '(__fish_complete_groups)' -r
complete -c example -f -n '_example_state _example__action' -l hosts -d 'hosts flag' -a '(__fish_print_hostnames)' -r
complete -c example -f -n '_example_state _example__action' -l message -s m -d 'message flag' -a '(echo -e "ERR\tmessage example\n_")' -r
complete -c example -f -n '_example_state _example__action' -l multi_parts -d 'multi_parts flag' -a '(echo -e "multi/parts\nmulti/parts/example\nmulti/parts/test\nexample/parts")' -r
complete -c example -f -n '_example_state _example__action' -l net_interfaces -s n -d 'net_interfaces flag' -a '(__fish_print_interfaces)' -r
complete -c example -f -n '_example_state _example__action' -l users -s u -d 'users flag' -a '(__fish_complete_users)' -r
complete -c example -f -n '_example_state _example__action' -l values -s v -d 'values flag' -a '(echo -e "values\nexample")' -r
complete -c example -f -n '_example_state _example__action' -l values_described -s d -d 'values with description flag' -a '(echo -e "values\tvalueDescription\nexample\texampleDescription\n\n")' -r
complete -c example -f -n '_example_state _example__action' -a '(_example_callback _)'


complete -c example -f -n '_example_state _example__callback' -l callback -s c -d 'Help message for callback' -a '(_example_callback _example__callback##callback)' -r
complete -c example -f -n '_example_state _example__callback' -a '(_example_callback _)'


complete -c example -f -n '_example_state _example__condition' -l required -s r -d 'required flag' -a '(echo -e "valid\ninvalid")' -r
complete -c example -f -n '_example_state _example__condition' -a '(_example_callback _)'


complete -c example -f -n '_example_state _example__help' -a '(_example_callback _)'


complete -c example -f -n '_example_state _example__injection' -a '(_example_callback _)'
`
	rootCmd.InitDefaultHelpCmd()
	assert.Equal(t, expected, carapace.Gen(rootCmd).Fish())
}

func TestPowershell(t *testing.T) {
	expected := `using namespace System.Management.Automation
using namespace System.Management.Automation.Language
Register-ArgumentCompleter -Native -CommandName 'example' -ScriptBlock {
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
        example _carapace powershell "$uid" $($commandElements| Foreach {$_.Extent}) _ | Out-String | Invoke-Expression
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
            }
            break
        }
                }
            }

        '_example__action' {
            switch -regex ($previous) {
                '^(-c|--custom)$' {
                         
                        break
                      }
                '^(--directories)$' {
                        [CompletionResult]::new('', '', [CompletionResultType]::ParameterValue, '') 
                        break
                      }
                '^(-f|--files)$' {
                        [CompletionResult]::new('', '', [CompletionResultType]::ParameterValue, '') 
                        break
                      }
                '^(-g|--groups)$' {
                        $(Get-Localgroup).Name 
                        break
                      }
                '^(--hosts)$' {
                         
                        break
                      }
                '^(-m|--message)$' {
                        [CompletionResult]::new('_ ', '_', [CompletionResultType]::ParameterValue, 'message example')
                        [CompletionResult]::new('ERR ', 'ERR', [CompletionResultType]::ParameterValue, 'message example')
                        
                         
                        break
                      }
                '^(--multi_parts)$' {
                        [CompletionResult]::new('multi/parts ', 'multi/parts', [CompletionResultType]::ParameterValue, ' ')
                        [CompletionResult]::new('multi/parts/example ', 'multi/parts/example', [CompletionResultType]::ParameterValue, ' ')
                        [CompletionResult]::new('multi/parts/test ', 'multi/parts/test', [CompletionResultType]::ParameterValue, ' ')
                        [CompletionResult]::new('example/parts ', 'example/parts', [CompletionResultType]::ParameterValue, ' ') 
                        break
                      }
                '^(-n|--net_interfaces)$' {
                        $(Get-NetAdapter).Name 
                        break
                      }
                '^(-u|--users)$' {
                        $(Get-LocalUser).Name 
                        break
                      }
                '^(-v|--values)$' {
                        [CompletionResult]::new('values ', 'values', [CompletionResultType]::ParameterValue, ' ')
                        [CompletionResult]::new('example ', 'example', [CompletionResultType]::ParameterValue, ' ') 
                        break
                      }
                '^(-d|--values_described)$' {
                        [CompletionResult]::new('values ', 'values', [CompletionResultType]::ParameterValue, 'valueDescription')
                        [CompletionResult]::new('example ', 'example', [CompletionResultType]::ParameterValue, 'exampleDescription')
                        
                         
                        break
                      }
                default {

            if ($wordToComplete -like "-*") {
                [CompletionResult]::new('-c ', '-c', [CompletionResultType]::ParameterName, 'custom flag')
                [CompletionResult]::new('--custom ', '--custom', [CompletionResultType]::ParameterName, 'custom flag')
                [CompletionResult]::new('--directories ', '--directories', [CompletionResultType]::ParameterName, 'files flag')
                [CompletionResult]::new('-f ', '-f', [CompletionResultType]::ParameterName, 'files flag')
                [CompletionResult]::new('--files ', '--files', [CompletionResultType]::ParameterName, 'files flag')
                [CompletionResult]::new('-g ', '-g', [CompletionResultType]::ParameterName, 'groups flag')
                [CompletionResult]::new('--groups ', '--groups', [CompletionResultType]::ParameterName, 'groups flag')
                [CompletionResult]::new('--hosts ', '--hosts', [CompletionResultType]::ParameterName, 'hosts flag')
                [CompletionResult]::new('-m ', '-m', [CompletionResultType]::ParameterName, 'message flag')
                [CompletionResult]::new('--message ', '--message', [CompletionResultType]::ParameterName, 'message flag')
                [CompletionResult]::new('--multi_parts ', '--multi_parts', [CompletionResultType]::ParameterName, 'multi_parts flag')
                [CompletionResult]::new('-n ', '-n', [CompletionResultType]::ParameterName, 'net_interfaces flag')
                [CompletionResult]::new('--net_interfaces ', '--net_interfaces', [CompletionResultType]::ParameterName, 'net_interfaces flag')
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

        '_example__callback' {
            switch -regex ($previous) {
                '^(-c|--callback)$' {
                        _example_callback '_example__callback##callback' 
                        break
                      }
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

        '_example__condition' {
            switch -regex ($previous) {
                '^(-r|--required)$' {
                        [CompletionResult]::new('valid ', 'valid', [CompletionResultType]::ParameterValue, ' ')
                        [CompletionResult]::new('invalid ', 'invalid', [CompletionResultType]::ParameterValue, ' ') 
                        break
                      }
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

        '_example__help' {
            switch -regex ($previous) {

                default {

            if ($wordToComplete -like "-*") {
            } else {
                _example_callback '_'
            }
            break
        }
                }
            }

        '_example__injection' {
            switch -regex ($previous) {

                default {

            if ($wordToComplete -like "-*") {
            } else {
                _example_callback '_'
            }
            break
        }
                }
            }

    })

    if ($completions.count -eq 0) {
      return "" # prevent default file completion
    }

    $completions.Where{ $_.CompletionText -like "$wordToComplete*" } |
        Sort-Object -Property ListItemText
}`
	rootCmd.InitDefaultHelpCmd()
	assert.Equal(t, expected, carapace.Gen(rootCmd).Powershell())
}

func TestZsh(t *testing.T) {
	expected := `#compdef example
function _example {
  local -a commands
  # shellcheck disable=SC2206
  local -a -x os_args=(${words})

  _arguments -C \
    "(*-a *--array)"{\*-a,\*--array}"[multiflag]: :" \
    "(-p --persistentFlag)"{-p,--persistentFlag}"[Help message for persistentFlag]" \
    "(-t --toggle)"{-t,--toggle}"[Help message for toggle]: :_values '' true false" \
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
  esac
}

function _example__action {
    _arguments -C \
    "(-c --custom)"{-c,--custom}"[custom flag]: :_most_recent_file 2" \
    "--directories[files flag]: :_files -/" \
    "(-f --files)"{-f,--files}"[files flag]: :_files -g '*.go'" \
    "(-g --groups)"{-g,--groups}"[groups flag]: :_groups" \
    "--hosts[hosts flag]: :_hosts" \
    "(-m --message)"{-m,--message}"[message flag]: : _message -r 'message example'" \
    "--multi_parts[multi_parts flag]: :_multi_parts / '(multi/parts multi/parts/example multi/parts/test example/parts)'" \
    "(-n --net_interfaces)"{-n,--net_interfaces}"[net_interfaces flag]: :_net_interfaces" \
    "(-u --users)"{-u,--users}"[users flag]: :_users" \
    "(-v --values)"{-v,--values}"[values flag]: :_values '' values example" \
    "(-d --values_described)"{-d,--values_described}"[values with description flag]: :_values '' 'values[valueDescription]' 'example[exampleDescription]'  " \
    "1:: :_values '' positional1 p1" \
    "2:: :_values '' positional2 p2"
}

function _example__callback {
    _arguments -C \
    "(-c --callback)"{-c,--callback}"[Help message for callback]: : eval \$(example _carapace zsh '_example__callback##callback' ${${os_args:1:gs/\"/\\\"}:gs/\'/\\\"})" \
    "1:: : eval \$(example _carapace zsh '_example__callback#1' ${${os_args:1:gs/\"/\\\"}:gs/\'/\\\"})"
}

function _example__condition {
    _arguments -C \
    "(-r --required)"{-r,--required}"[required flag]: :_values '' valid invalid" \
    "1:: : eval \$(example _carapace zsh '_example__condition#1' ${${os_args:1:gs/\"/\\\"}:gs/\'/\\\"})"
}

function _example__help {
    _arguments -C \
    "*::arg:->args"
}

function _example__injection {
    _arguments -C \
    "1:: :_values '' echo\ fail" \
    "2:: :_values '' echo\ fail" \
    "3:: :_values '' echo\ fail" \
    "4:: :_values '' \ echo\ fail\ " \
    "5:: :_values '' \ echo\ fail\ " \
    "6:: :_values '' \ echo\ fail\ " \
    "7:: :_values '' echo\ fail" \
    "8:: : _message -r 'no values to complete'" \
    "9:: :_values '' LAST\ POSITIONAL\ VALUE"
}
if compquote '' 2>/dev/null; then _example; else compdef _example example; fi
`
	rootCmd.InitDefaultHelpCmd()
	assert.Equal(t, expected, carapace.Gen(rootCmd).Zsh())
}
