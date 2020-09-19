package xonsh

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/rsteube/carapace/common"
	"github.com/rsteube/carapace/uid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Snippet(cmd *cobra.Command, actions map[string]string) string {
	buf := new(bytes.Buffer)

	var subCommandCases bytes.Buffer
	generateXonshSubcommandCases(&subCommandCases, cmd, actions)
	fmt.Fprintf(buf, xonshCompletionTemplate, cmd.Name(), cmd.Name(), uid.Executable(), cmd.Name(), uid.Executable(), subCommandCases.String(), cmd.Name(), cmd.Name())

	return buf.String()
}

var xonshCompletionTemplate = `from shlex import split
import re
import pathlib
import subprocess
import xonsh
from xonsh.completers._aliases import _add_one_completer
from xonsh.completers.path import complete_dir, complete_path
from xonsh.completers.tools import RichCompletion

def %v_completer(prefix, line, begidx, endidx, ctx):
    if not line.startswith('%v '):
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
    state, _ = subprocess.Popen(['%v', '_carapace', 'xonsh', 'state', *words],
                                   stdout=subprocess.PIPE,
                                   stderr=subprocess.PIPE).communicate()
    state = state.decode('utf-8').split('\n')[0]
   
    # TODO python callback function
    def _%v_callback(uid):
        cb, _ = subprocess.Popen(['%v', '_carapace', 'xonsh', uid, *words],
                                     stdout=subprocess.PIPE,
                                     stderr=subprocess.PIPE).communicate()
        cb = cb.decode('utf-8')
        return eval(cb)
   
    if False:
        pass%s

    result = set(filter(lambda x: x.startswith(current) and x.endswith(suffix), result))
    if len(result) == 0:
        result = {RichCompletion(current, display=current, description='', prefix_len=0)}

    result = set(map(lambda x: RichCompletion(x[:len(x)-len(suffix)], display=x.display, description=x.description, prefix_len=len(current)), result))
    return result
_add_one_completer('%v', %v_completer, 'start')
`

func generateXonshSubcommandCases(out io.Writer, cmd *cobra.Command, actions map[string]string) {
	var cmdName = fmt.Sprintf("%v", uid.Command(cmd))

	fmt.Fprintf(out, "\n    elif state == '%s':", cmdName)
	fmt.Fprintf(out, `
        if False: # switch previous
            pass
%v
        else:
            if False:
                pass
%v
%v
`, snippetFlagActions(cmd, actions, false), "    "+strings.Replace(snippetFlagActions(cmd, actions, true), "\n", "\n    ", -1), snippetTODO(cmd))

	for _, subCmd := range cmd.Commands() {
		if !subCmd.Hidden {
			generateXonshSubcommandCases(out, subCmd, actions)
		}
	}
}

func snippetFlagActions(cmd *cobra.Command, actions map[string]string, optArgFlag bool) string {
	flagActions := make([]string, 0)
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		// TODO cleanup this mess
		if flag.NoOptDefVal != "" && !optArgFlag {
			return
		}
		if flag.NoOptDefVal == "" && optArgFlag {
			return
		}

		optArgSuffix := ""
		if flag.NoOptDefVal != "" {
			optArgSuffix = "=.*"
		}

		match := fmt.Sprintf(`^(--%v)$`, flag.Name+optArgSuffix)
		if flag.Shorthand != "" {
			match = fmt.Sprintf(`^(-%v|--%v)$`, flag.Shorthand+optArgSuffix, flag.Name+optArgSuffix)
		} else if common.IsShorthandOnly(flag) {
			match = fmt.Sprintf(`^(-%v)$`, flag.Shorthand+optArgSuffix)
		}
		var action = "{}"
		if a, ok := actions[uid.Flag(cmd, flag)]; ok { // TODO cleanup
			action = a
		}
		if flag.NoOptDefVal != "" {
			flagActions = append(flagActions, fmt.Sprintf(`        elif re.search('%v',current):
            result = %v
            result = set(map(lambda x: RichCompletion(current.split('=')[0]+'='+x, display=x.display, description=x.description, prefix_len=x.prefix_len), result))
`, match, strings.Replace(action, "\n", "\n                        ", -1)))

		} else {
			flagActions = append(flagActions, fmt.Sprintf(`        elif re.search('%v',previous):
            result = %v
                  `, match, strings.Replace(action, "\n", "\n                        ", -1)))
		}
	})
	return strings.Join(flagActions, "\n")
}

func snippetTODO(cmd *cobra.Command) string {
	result := ""
	result += fmt.Sprint("\n            elif re.search(\"-.*\",current):")
	result += fmt.Sprint("\n                result = {")

	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden {
			if len(flag.Shorthand) > 0 {
				result += fmt.Sprintf("\n                    RichCompletion('-%v', display='-%v', description='%v', prefix_len=0),", flag.Shorthand, flag.Shorthand, sanitizer.Replace(flag.Usage))

			}
			if !common.IsShorthandOnly(flag) {
				result += fmt.Sprintf("\n                    RichCompletion('--%v', display='--%v', description='%v', prefix_len=0),", flag.Name, flag.Name, sanitizer.Replace(flag.Usage))
			}
		}
	})
	result += fmt.Sprint("\n                }")

	result += fmt.Sprint("\n            else:")
	// TODO wrap in `result = {}`
	if cmd.HasAvailableSubCommands() {
		result += fmt.Sprint("\n                result = {")
		for _, subCmd := range cmd.Commands() {
			if !subCmd.Hidden {
				result += fmt.Sprintf("\n                RichCompletion('%v', display='%v', description='%v', prefix_len=0),", subCmd.Name(), subCmd.Name(), sanitizer.Replace(cmd.Short))
			}
		}
		result += fmt.Sprint("\n                }")
	} else {
		if !cmd.HasAvailableSubCommands() {
			result += fmt.Sprintf("\n                result = _%v_callback('_')", cmd.Root().Name())
		}
	}
	result += fmt.Sprint("\n")
	return result
}
