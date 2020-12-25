use str
edit:completion:arg-completer[example] = [@arg]{
  fn _example_callback [uid]{
    if (eq $arg[-1] "") {
        arg[-1] = "''"
    }
    eval (echo (str:join "\001" $arg) | xargs --delimiter="\001" example _carapace elvish $uid | slurp) &ns=(ns [&arg=$arg])
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
    ]
    subargs = $arg[(subindex multiparts):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }
}

