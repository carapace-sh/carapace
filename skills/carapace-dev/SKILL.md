---
name: carapace-dev
description: >
  Use when developing or debugging the carapace library — Action API, traverse engine,
  shell formatters, style system, pflag extensions, or internal implementation details.
  Triggers on: "carapace" (development context), "carapace Action", "carapace traverse",
  "carapace shell", "carapace style", "pflagfork", "carapace-bridge".
user-invocable: true
---

# Carapace Development Reference

Reference for developing and debugging the [carapace](https://github.com/carapace-sh/carapace) library.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| Action type, modifier, InvokedAction, Context, callback, Batch, rawValues, suffix naming | [references/action.md](references/action.md) |
| traverse, completion engine, argument classification, pflagfork, shell patching, complete(), dispatch | [references/traverse.md](references/traverse.md) |
| pflag, flag mode, NameAsShorthand, ShorthandOnly, Nargs, OptargDelimiter, ArgumentStyle, non-POSIX flags | [references/pflag.md](references/pflag.md) |
| Shell overview, cross-shell comparison, dispatch pipeline, supported shells, nospace/message/quoting comparison | [references/shell.md](references/shell.md) |
| Bash, COMP_TYPE, COMP_WORDBREAKS, redirect patching, partial completion, bash quoting | [references/shell-bash.md](references/shell-bash.md) |
| Bash BLE, ble.sh, tab-delimited format, per-candidate suffix, ble/complete/cand/yield | [references/shell-bash-ble.md](references/shell-bash-ble.md) |
| Oil, OSH, inline nospace indicator, simpler snippet | [references/shell-oil.md](references/shell-oil.md) |
| Zsh, _describe, zstyle, named directory, quoting state machine, CARAPACE_COMPLINE | [references/shell-zsh.md](references/shell-zsh.md) |
| Fish, commandline, tab-separated format, open-quote retry, fish pager | [references/shell-fish.md](references/shell-fish.md) |
| Elvish, complexCandidate, CodeSuffix, styled, edit:notify, from-json, ParseStyling | [references/shell-elvish.md](references/shell-elvish.md) |
| Nushell, from json, JSON record, open-quote patching, style conversion, nushell.Patch | [references/shell-nushell.md](references/shell-nushell.md) |
| PowerShell, CompletionResult, Register-ArgumentCompleter, commandAst, SGR color, tooltip | [references/shell-powershell.md](references/shell-powershell.md) |
| Xonsh, RichCompletion, prefix_len, append_closing_quote, contextual_command_completer | [references/shell-xonsh.md](references/shell-xonsh.md) |
| Style system, XTerm256, TrueColor, ForKeyword, ForPath, ForPathExt, LS_COLORS, style config | [references/style.md](references/style.md) |
| Bridge, bridgeActions, actionCommand, output parsing, NoSpace detection, shell config, Detect() | [references/bridge.md](references/bridge.md) |

## Quick Guide

- **How do I add/modify an Action?** → [references/action.md](references/action.md)
- **How does argument traversal work?** → [references/traverse.md](references/traverse.md)
- **How do non-POSIX flags work?** → [references/pflag.md](references/pflag.md)
- **Which shells support which features?** → [references/shell.md](references/shell.md)
- **How does a specific shell format completions?** → [references/shell-{name}.md](references/shell-bash.md)
- **How do styles/colors work?** → [references/style.md](references/style.md)
- **How do bridges work internally?** → [references/bridge.md](references/bridge.md)

## Cross-Project References

For user-facing topics (integrating carapace into a CLI, writing specs, macros, setup, choices), use the **carapace** skill (in the carapace-bin repo).
