package powershell

import (
	"fmt"

	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

func Snippet(cmd *cobra.Command, actions map[string]string) string {
	return fmt.Sprintf(`using namespace System.Management.Automation
using namespace System.Management.Automation.Language
Function _%v_completer {
    [System.Diagnostics.CodeAnalysis.SuppressMessageAttribute("PSAvoidUsingInvokeExpression", "", Scope="Function", Target="*")]
    param($wordToComplete, $commandAst) #, $cursorPosition)
    $commandElements = $commandAst.CommandElements


    $completions = @(
      if (!$wordToComplete) {
        %v _carapace powershell _ $($commandElements| ForEach-Object {$_.Extent}) '""' | ConvertFrom-Json | ForEach-Object { [CompletionResult]::new($_.CompletionText, $_.ListItemText, [CompletionResultType]::ParameterValue, $_.ToolTip) }
      } else {
        %v _carapace powershell _ $($commandElements| ForEach-Object {$_.Extent}) | ConvertFrom-Json | ForEach-Object { [CompletionResult]::new($_.CompletionText, $_.ListItemText, [CompletionResultType]::ParameterValue, $_.ToolTip) }
      }
    )

    if ($completions.count -eq 0) {
      return "" # prevent default file completion
    }

    $completions.Where{ ($_.CompletionText -replace '`+"`"+`','') -like "$wordToComplete*" } |
        Sort-Object -Property ListItemText
}
Register-ArgumentCompleter -Native -CommandName '%v' -ScriptBlock (Get-Item "Function:_%v_completer").ScriptBlock
`, cmd.Name(), uid.Executable(), uid.Executable(), cmd.Name(), cmd.Name())
}
