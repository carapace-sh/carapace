from shlex import split
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

