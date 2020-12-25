function _example_quote_suffix
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

