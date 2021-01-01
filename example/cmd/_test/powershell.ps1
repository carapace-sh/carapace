using namespace System.Management.Automation
using namespace System.Management.Automation.Language
Function _example_completer {
    [System.Diagnostics.CodeAnalysis.SuppressMessageAttribute("PSAvoidUsingInvokeExpression", "", Scope="Function", Target="*")]
    param($wordToComplete, $commandAst) #, $cursorPosition)
    $commandElements = $commandAst.CommandElements


    $completions = @(
      if (!$wordToComplete) {
        example _carapace powershell _ $($commandElements| ForEach-Object {$_.Extent}) '""' | Out-String | Invoke-Expression
      } else {
        example _carapace powershell _ $($commandElements| ForEach-Object {$_.Extent}) | Out-String | Invoke-Expression
      }
    )

    if ($completions.count -eq 0) {
      return "" # prevent default file completion
    }

    $completions.Where{ ($_.CompletionText -replace '`','') -like "$wordToComplete*" } |
        Sort-Object -Property ListItemText
}
Register-ArgumentCompleter -Native -CommandName 'example' -ScriptBlock (Get-Item "Function:_example_completer").ScriptBlock

