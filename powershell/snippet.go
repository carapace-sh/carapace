package powershell

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/rsteube/carapace/uid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Snippet(cmd *cobra.Command, actions map[string]string) string {
	buf := new(bytes.Buffer)

	var subCommandCases bytes.Buffer
	generatePowerShellSubcommandCases(&subCommandCases, cmd, actions)
	fmt.Fprintf(buf, powerShellCompletionTemplate, cmd.Name(), uid.Executable(), cmd.Name(), uid.Executable(), uid.Executable(), subCommandCases.String())

	return buf.String()
}

var powerShellCompletionTemplate = `using namespace System.Management.Automation
using namespace System.Management.Automation.Language
Register-ArgumentCompleter -Native -CommandName '%s' -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $commandElements = $commandAst.CommandElements
    $previous = $commandElements[-1].Extent
    if ($wordToComplete) {
        $previous = $commandElements[-2].Extent
    }

    $state = %v _carapace powershell state $($commandElements| Foreach {$_.Extent})
    
    Function _%v_callback {
      param($uid)
      if (!$wordToComplete) {
        %v _carapace powershell "$uid" $($commandElements| Foreach {$_.Extent}) "''" | Out-String | Invoke-Expression
      } else {
        %v _carapace powershell "$uid" $($commandElements| Foreach {$_.Extent}) | Out-String | Invoke-Expression
      }
    }
    
    $completions = @(switch ($state) {%s
    })

    if ($completions.count -eq 0) {
      return "" # prevent default file completion
    }

    $completions.Where{ $_.CompletionText -like "$wordToComplete*" } |
        Sort-Object -Property ListItemText
}`

func generatePowerShellSubcommandCases(out io.Writer, cmd *cobra.Command, actions map[string]string) {
	var cmdName = fmt.Sprintf("%v", uid.Command(cmd))

	fmt.Fprintf(out, "\n        '%s' {", cmdName)
	fmt.Fprintf(out, `
            switch -regex ($previous) {
%v
                default {
%v
                }
            }
`, snippetFlagActions(cmd, actions), snippetTODO(cmd))

	for _, subCmd := range cmd.Commands() {
		if !subCmd.Hidden {
			generatePowerShellSubcommandCases(out, subCmd, actions)
		}
	}
}

func snippetFlagActions(cmd *cobra.Command, actions map[string]string) string {
	flagActions := make([]string, 0)
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if flag.NoOptDefVal != "" {
			return
		}

		match := fmt.Sprintf(`^(--%v)$`, flag.Name)
		if flag.Shorthand != "" {
			match = fmt.Sprintf(`^(-%v|--%v)$`, flag.Shorthand, flag.Name)
		} else if flag.ShorthandOnly {
			match = fmt.Sprintf(`^(-%v)$`, flag.Shorthand)
		}
		var action = ""
		if a, ok := actions[uid.Flag(cmd, flag)]; ok { // TODO cleanup
			action = a
		}
		flagActions = append(flagActions, fmt.Sprintf(`                '%v' {
                        %v 
                        break
                      }`, match, strings.Replace(action, "\n", "\n                        ", -1)))
	})
	return strings.Join(flagActions, "\n")
}

func snippetTODO(cmd *cobra.Command) string {
	result := ""
	result += fmt.Sprint("\n            if ($wordToComplete -like \"-*\") {")

	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		if !flag.Hidden {
			usage := escapeStringForPowerShell(flag.Usage)
			if len(flag.Shorthand) > 0 {
				result += fmt.Sprintf("\n                [CompletionResult]::new('-%s ', '-%s', [CompletionResultType]::ParameterName, '%s')", flag.Shorthand, flag.Shorthand, usage)
			}
			if !flag.ShorthandOnly {
				result += fmt.Sprintf("\n                [CompletionResult]::new('--%s ', '--%s', [CompletionResultType]::ParameterName, '%s')", flag.Name, flag.Name, usage)
			}
		}
	})

	result += fmt.Sprint("\n            } else {")
	for _, subCmd := range cmd.Commands() {
		if !subCmd.Hidden {
			usage := escapeStringForPowerShell(subCmd.Short)
			result += fmt.Sprintf("\n                [CompletionResult]::new('%s ', '%s', [CompletionResultType]::Command, '%s')", subCmd.Name(), subCmd.Name(), usage)
		}
	}

	if !cmd.HasAvailableSubCommands() {
		result += fmt.Sprintf("\n                _%v_callback '_'", cmd.Root().Name())
	}
	result += fmt.Sprint("\n            }")
	result += fmt.Sprint("\n            break\n        }")
	return result
}

func escapeStringForPowerShell(s string) string {
	if s == "" {
		return " " // completion fails if empty (fallback to file completion)
	}
	return strings.Replace(s, "'", "''", -1)
}
