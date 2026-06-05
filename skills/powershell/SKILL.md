---
name: powershell
description: >
  Use when working with PowerShell internals — tab completion, Register-ArgumentCompleter,
  PSReadLine, CompletionResult, CommandAst, argument passing to native commands, quoting rules,
  $PSStyle ANSI rendering, profiles and configuration, or the prediction plugin system.
  Triggers on: "powershell", "pwsh", "Register-ArgumentCompleter", "PSReadLine",
  "CompletionResult", "CommandAst", "TabExpansion2", "MenuComplete", "ArgumentCompleter",
  "PSStyle", "powershell profile", "powershell completion", "powershell quoting",
  "powershell argument passing", "IPredictor", "predictive intellisense".
user-invocable: true
---

# PowerShell Shell In-Depth Reference

Comprehensive reference for PowerShell internals, with emphasis on the completion system, PSReadLine, and how external tools hook into the completion menu.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| Register-ArgumentCompleter, ArgumentCompleter, ArgumentCompletions, TabExpansion2, CommandCompletion, CompletionResult, CompletionResultType, CommandAst, CompletionAnalysis, CompletionCompleters, CompleteInput, native mode, parameter mode, NativeFallback, completion registration, completion routing, pseudo parameter binding, type inference, ValidateSet, ending keys, completion display, nospace, CompletionContext | [references/completion.md](references/completion.md) |
| PSReadLine, Set-PSReadLineKeyHandler, Set-PSReadLineOption, edit mode, emacs, vi, MenuComplete, TabCompleteNext, PossibleCompletions, completion menu, menu rendering, column layout, tooltip, selection color, prediction, IPredictor, Predictive IntelliSense, InlineView, ListView, CommandPrediction, prediction plugin, PSConsoleReadLine, key binding, custom handler, rendering pipeline, token coloring, differential rendering | [references/editor.md](references/editor.md) |
| $PSStyle, ANSI escape, VT100, SGR, Write-Host, ForegroundColor, BackColor, virtual terminal, escape sequence, `` `e ``, CSI, true color, FromRgb, OutputRendering, StringDecorated, terminal rendering, console color, formatting color, FileInfo coloring, progress bar | [references/styling.md](references/styling.md) |
| quoting, single quotes, double quotes, here-strings, backtick, escape character, variable expansion, subexpression, argument passing, native command, stop-parsing, --%, PSNativeCommandArgumentPassing, Legacy, Standard, Windows mode, comma, pipeline, parameter binding, expression mode, argument mode | [references/language.md](references/language.md) |
| profile, $PROFILE, AllUsers, CurrentUser, startup, powershell.config.json, execution policy, $PSHOME, PSModulePath, experimental features, configuration, Group Policy, logging, transcription, script block logging, module logging | [references/startup-config.md](references/startup-config.md) |

## Quick Guide

- **How do I register a custom completer?** → [references/completion.md](references/completion.md)
- **How does PSReadLine render the completion menu?** → [references/editor.md](references/editor.md)
- **How do I create a prediction plugin?** → [references/editor.md](references/editor.md)
- **How do CompletionResultType values affect behavior?** → [references/completion.md](references/completion.md)
- **How does PowerShell handle ANSI colors?** → [references/styling.md](references/styling.md)
- **How do I use $PSStyle for styled output?** → [references/styling.md](references/styling.md)
- **How does quoting affect argument passing to native commands?** → [references/language.md](references/language.md)
- **What is the stop-parsing token (--%)?** → [references/language.md](references/language.md)
- **How does PSNativeCommandArgumentPassing work?** → [references/language.md](references/language.md)
- **Which profile files does PowerShell load?** → [references/startup-config.md](references/startup-config.md)
- **How do I configure powershell.config.json?** → [references/startup-config.md](references/startup-config.md)
- **How does TabExpansion2 route to completers?** → [references/completion.md](references/completion.md)
- **How do I customize PSReadLine key bindings?** → [references/editor.md](references/editor.md)
- **How does the AST provide context to completers?** → [references/completion.md](references/completion.md)

## Cross-Project References

For carapace-specific PowerShell integration (snippet, value formatting, SGR styling, nospace, single-quote stripping), see the **carapace-dev** skill → `references/shell-powershell.md`.
