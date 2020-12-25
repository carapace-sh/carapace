#compdef example
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
    "1: :{_example_callback '_example__multiparts#1'}"
}
compquote '' 2>/dev/null && _example
compdef _example example

