package xonsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

func Snippet(cmd *cobra.Command) string {
	functionName := strings.Replace(cmd.Name(), "-", "__", -1)
	return fmt.Sprintf(`def _%v_completer(prefix, line, begidx, endidx, ctx):
    """carapace completer for %v"""
    if not line.startswith('%v '):
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
    
    for index, item in enumerate(wordsNonPosix):
        if item.startswith("r'"):
            words[index] = words[index][1:] # fix rawstring not handled by shlex

    current=words[-1]
    previous=words[-2]

    result = {}

    def _%v_callback():
        from json import loads
        from subprocess import Popen, PIPE
        from xonsh.completers.tools import RichCompletion
        cb, _ = Popen(['%v', '_carapace', 'xonsh', '_', *words], stdout=PIPE, stderr=PIPE).communicate()
        try:
          return {RichCompletion(c["Value"], display=c["Display"], description=c["Description"], prefix_len=0) for c in loads(cb)}
        except:
          return {}
  
    result = _%v_callback()

    if len(result) == 0:
        result = {RichCompletion(current, display=current, description='', prefix_len=0)}

    result = set(map(lambda x: RichCompletion(x, display=x.display, description=x.description, prefix_len=endidx-begidx), result))
    return result

from xonsh.completers._aliases import _add_one_completer
_add_one_completer('%v', _%v_completer, 'start')
`, functionName, cmd.Name(), cmd.Name(), functionName, uid.Executable(), functionName, cmd.Name(), functionName)
}
