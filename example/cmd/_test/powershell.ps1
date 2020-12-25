using namespace System.Management.Automation
using namespace System.Management.Automation.Language
$_example_completer = {
    param($wordToComplete, $commandAst) #, $cursorPosition)
    $commandElements = $commandAst.CommandElements
    $previous = $commandElements[-1].Value
    if ($wordToComplete) {
        $previous = $commandElements[-2].Value
    }

    $state = example _carapace powershell state $($commandElements| ForEach-Object {$_.Value})

    Function _example_callback {
      [System.Diagnostics.CodeAnalysis.SuppressMessageAttribute("PSAvoidUsingInvokeExpression", "", Scope="Function", Target="*")]
      param($uid)
      if (!$wordToComplete) {
        example _carapace powershell "$uid" $($commandElements| ForEach-Object {$_.Value}) '""' | Out-String | Invoke-Expression
      } else {
        example _carapace powershell "$uid" $($commandElements| ForEach-Object {$_.Value}) | Out-String | Invoke-Expression
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
                        [CompletionResult]::new('KILL', 'KILL ', [CompletionResultType]::ParameterValue, 'Kill`, unblockable ')
                        [CompletionResult]::new('PIPE', 'PIPE ', [CompletionResultType]::ParameterValue, 'Broken pipe ')
                        [CompletionResult]::new('POLL', 'POLL ', [CompletionResultType]::ParameterValue, 'Pollable event occurred ')
                        [CompletionResult]::new('PROF', 'PROF ', [CompletionResultType]::ParameterValue, 'Profiling alarm clock timer expired ')
                        [CompletionResult]::new('PWR', 'PWR ', [CompletionResultType]::ParameterValue, 'Power failure restart ')
                        [CompletionResult]::new('QUIT', 'QUIT ', [CompletionResultType]::ParameterValue, 'Quit from keyboard ')
                        [CompletionResult]::new('SEGV', 'SEGV ', [CompletionResultType]::ParameterValue, 'Segmentation violation ')
                        [CompletionResult]::new('STKFLT', 'STKFLT ', [CompletionResultType]::ParameterValue, 'Stack fault on coprocessor ')
                        [CompletionResult]::new('STOP', 'STOP ', [CompletionResultType]::ParameterValue, 'Stop process`, unblockable ')
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
                [CompletionResult]::new('--comma ', '--comma', [CompletionResultType]::ParameterName, 'multiparts with `, as divider')
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

    $completions.Where{ ($_.CompletionText -replace '`','') -like "$wordToComplete*" } |
        Sort-Object -Property ListItemText
}
Register-ArgumentCompleter -Native -CommandName 'example' -ScriptBlock $_example_completer

