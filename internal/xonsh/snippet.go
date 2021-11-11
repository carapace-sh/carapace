package xonsh

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

// Snippet creates the xonsh completion script
func Snippet(cmd *cobra.Command) string {
	functionName := strings.Replace(cmd.Name(), "-", "__", -1)
	return fmt.Sprintf(`from xonsh.completers.tools import *

@contextual_command_completer
def _%v_completer(context):
    """carapace completer for %v"""
    if context.completing_command('%v'):
        from json import loads
        from subprocess import Popen, PIPE
        from xonsh.completers.tools import RichCompletion
        output, _ = Popen(['%v', '_carapace', 'xonsh', '_', *[a.value for a in context.args], context.prefix], stdout=PIPE, stderr=PIPE).communicate()
        try:
            return {RichCompletion(c["Value"], display=c["Display"], description=c["Description"], prefix_len=len(context.raw_prefix)) for c in loads(output)}
        except:
            return


from xonsh.completers._aliases import _add_one_completer
_add_one_completer('%v', _%v_completer, 'start')
`, functionName, cmd.Name(), cmd.Name(), uid.Executable(), cmd.Name(), functionName)
}
