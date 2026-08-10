using namespace System.Management.Automation
using namespace System.Management.Automation.Language

$_example__multi_completer = {
    [System.Diagnostics.CodeAnalysis.SuppressMessageAttribute("PSAvoidUsingInvokeExpression", "", Scope="Function", Target="*")]
    param($wordToComplete, $commandAst, $cursorPosition)
    $env:CARAPACE_SHELL = "powershell"
    $commandElements = $commandAst.CommandElements

    $elems = @()
    foreach ($_ in $commandElements) {
      if ($_.Extent.StartOffset -gt $cursorPosition) {
          break
      }
      $t = $_.Extent.Text
      if ($_.Extent.EndOffset -gt $cursorPosition) {
          $t = $t.Substring(0, $_.Extent.Text.get_Length() - ($_.Extent.EndOffset - $cursorPosition))
      }

      if ($t.Substring(0,1) -eq "'"){
        $t = $t.Substring(1)
      }
      if ($t.get_Length() -gt 0 -and $t.Substring($t.get_Length()-1) -eq "'"){
        $t = $t.Substring(0,$t.get_Length()-1)
      }
      if ($t.get_Length() -eq 0){
        $t = '""'
      }
      $elems += $t.replace('`,', ',') # quick fix
    }

    $completions = @(
      if (!$wordToComplete) {
        example-multi ($elems[0] -replace ('\.exe$', '')) _carapace powershell $($elems| ForEach-Object {$_}) '' | ConvertFrom-Json | ForEach-Object { [CompletionResult]::new($_.CompletionText, $_.ListItemText.replace('`e[', "`e["), [CompletionResultType]::ParameterValue, $_.ToolTip.replace('`e[', "`e[")) }
      } else {
        example-multi ($elems[0] -replace ('\.exe$', '')) _carapace powershell $($elems| ForEach-Object {$_}) | ConvertFrom-Json | ForEach-Object { [CompletionResult]::new($_.CompletionText, $_.ListItemText.replace('`e[', "`e["), [CompletionResultType]::ParameterValue, $_.ToolTip.replace('`e[', "`e[")) }
      }
    )

    if ($completions.count -eq 0) {
      return "" # prevent default file completion
    }

    $completions
}

Register-ArgumentCompleter -Native -ScriptBlock $_example__multi_completer -CommandName 'example-multi' # 'example-multi.exe'
Register-ArgumentCompleter -Native -ScriptBlock $_example__multi_completer -CommandName 'identify' # 'identify.exe'
Register-ArgumentCompleter -Native -ScriptBlock $_example__multi_completer -CommandName 'convert' # 'convert.exe'

