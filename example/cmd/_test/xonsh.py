def _example_completer(prefix, line, begidx, endidx, ctx):
    """carapace completer for example"""
    if not line.startswith('example '):
        return # not the expected command to complete

    from shlex import split
    from xonsh.completers.tools import RichCompletion
    
    full_words=split(line + "_") # ensure last word is empty when ends with space
    full_words[-1]=full_words[-1][0:-1]
    words=split(line[0:endidx] + "_") # ensure last word is empty when ends with space
    words[-1]=words[-1][0:-1]
    current=words[-1]
    previous=words[-2]
    suffix=full_words[len(words)-1][len(current):]

    result = {}

    def _example_callback():
        from subprocess import Popen, PIPE
        from xonsh.completers.tools import RichCompletion
        cb, _ = Popen(['example', '_carapace', 'xonsh', '_', *words],
                                     stdout=PIPE,
                                     stderr=PIPE).communicate()
        cb = cb.decode('utf-8')
        if cb == "":
            return {}
        else:
            nonlocal prefix, line, begidx, endidx, ctx
            return eval(cb)
  
    result = _example_callback()

    if len(result) == 0:
        result = {RichCompletion(current, display=current, description='', prefix_len=0)}

    result = set(map(lambda x: RichCompletion(x[:len(x)-(len(suffix)+suffix.count(' '))], display=x.display, description=x.description, prefix_len=len(current)+current.count(' ')), result))
    return result

from xonsh.completers._aliases import _add_one_completer
_add_one_completer('example', _example_completer, 'start')

