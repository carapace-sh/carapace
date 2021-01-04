def _example_completer(prefix, line, begidx, endidx, ctx):
    """carapace completer for example"""
    if not line.startswith('example '):
        return # not the expected command to complete

    from shlex import split
    from xonsh.completers.tools import RichCompletion
  
    words=""
    wordsNonPosix=""
    try:
        words=split(line[0:endidx] + "_") # ensure last word is empty when ends with space
        wordsNonPosix=split(line[0:endidx], posix=False)
    except:
        try:
            words=split(line[0:endidx] + '"' + "_") # ensure last word is empty when ends with space
            wordsNonPosix=split(line[0:endidx] + '"', posix=False)
            wordsNonPosix[-1]= wordsNonPosix[-1][:-1]
        except:
            words=split(line[0:endidx] + "'" + "_") # ensure last word is empty when ends with space
            wordsNonPosix=split(line[0:endidx] + "'", posix=False)
            wordsNonPosix[-1]= wordsNonPosix[-1][:-1]
    
    words[-1]=words[-1][0:-1]
    if len(words[-1]) != 0:
        begidx = endidx
        for word in reversed(wordsNonPosix):
            begidx = begidx - len(word)
            if line[begidx-1] == " ":
                break

    current=words[-1]
    previous=words[-2]

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

    def _example_quote(s):
        if " " in s:
            return '"' + s + '"'
        return s
        
    result = set(map(lambda x: RichCompletion(_example_quote(x), display=x.display, description=x.description, prefix_len=endidx-begidx), result))
    return result

from xonsh.completers._aliases import _add_one_completer
_add_one_completer('example', _example_completer, 'start')

