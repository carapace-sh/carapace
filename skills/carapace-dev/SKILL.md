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
| Storage, global map, mutex, entry, FlagCompletion, PositionalCompletion, PreRun, PreInvoke, storage.get, gotchas | [references/storage.md](references/storage.md) |
| Sandbox, Command, Package, Action, Files, Reply, Run, Expect, ExpectNot, ClearCache, SANDBOX env, subprocess isolation | [references/sandbox.md](references/sandbox.md) |
| Batch, Invoke, parallelize, Merge, ToA, invokedBatch, goroutine model, shared Context, parallel invocation | [references/batch.md](references/batch.md) |
| Cache, Action.Cache, File, Load, Write, cache key, file:line, key.Key types, FolderStats, TTL, sandbox override | [references/cache.md](references/cache.md) |
| UID, url.URL, cmd:// scheme, Command, Flag, mLocalFlags mutex, PathEscape, Executable, UidF, uid.Map | [references/uid.md](references/uid.md) |
| Config, styles.json, RegisterStyle, Load, CARAPACE_* env vars, LOG, ColorDisabled, Hidden, Sandbox, isGoRun | [references/config-env.md](references/config-env.md) |
| Spec, spec.Spec, yaml.Marshal, Command struct, Flags, PersistentFlags, _carapace skip, pflagfork.Flag.Definition | [references/spec.md](references/spec.md) |
| Export, Export struct, MarshalJSON, version, RawValue, Meta, cache format, sandbox invoke, ActionImport | [references/export.md](references/export.md) |
| Conditionals, Arch, Os, Executable, File, CompletingPath, UnlessF, Context predicates, pkg/condition | [references/conditionals.md](references/conditionals.md) |
| _carapace CLI, command.go, addCompletionCommand, spec, style set, export re-invocation, PositionalAnyCompletion | [references/command-cli.md](references/command-cli.md) |
| Mock, Mock struct, NewMock, Dir, Replies, CacheDir, WorkDir, t interface, Reply mechanism, sandbox integration | [references/mock.md](references/mock.md) |

## Quick Guide

- **How do I add/modify an Action?** → [references/action.md](references/action.md)
- **How does argument traversal work?** → [references/traverse.md](references/traverse.md)
- **How do non-POSIX flags work?** → [references/pflag.md](references/pflag.md)
- **Which shells support which features?** → [references/shell.md](references/shell.md)
- **How does a specific shell format completions?** → [references/shell-{name}.md](references/shell-bash.md)
- **How do styles/colors work?** → [references/style.md](references/style.md)
- **How do bridges work internally?** → [references/bridge.md](references/bridge.md)
- **How does the global storage map work?** → [references/storage.md](references/storage.md)
- **How do I write/run sandbox tests?** → [references/sandbox.md](references/sandbox.md)
- **How does Batch parallel invocation work?** → [references/batch.md](references/batch.md)
- **How does Action.Cache work?** → [references/cache.md](references/cache.md)
- **How do UID identifiers work?** → [references/uid.md](references/uid.md)
- **How is styles.json loaded?** → [references/config-env.md](references/config-env.md)
- **How does spec YAML generation work?** → [references/spec.md](references/spec.md)
- **What is the JSON wire format?** → [references/export.md](references/export.md)
- **What conditional helpers exist?** → [references/conditionals.md](references/conditionals.md)
- **How does the _carapace CLI work?** → [references/command-cli.md](references/command-cli.md)
- **How does the mock system work?** → [references/mock.md](references/mock.md)

## Cross-Project References

- For user-facing topics (integrating carapace into a CLI, writing specs, macros, setup, choices), use the **carapace** skill (in the carapace-bin repo).
- For in-depth bash shell knowledge (programmable completion, Readline, quoting/expansion, execution model, startup files), use the **bash** skill (in this repo).
- For in-depth bash BLE/ble.sh knowledge (cand/yield, completion sources and actions, menu styles, menu-filter, progcomp integration, sabbrev, dabbrev, auto-complete, syntax highlighting, faces, gspec colors, widgets, keymaps, bleopt, ble-bind, blehook, blerc, ble-import), use the **bash-ble** skill (in this repo).
- For in-depth clink/cmd.exe knowledge (argmatcher API, match generators, Readline integration, DLL injection, prompt filtering, auto-suggestions, input line coloring), use the **cmd-clink** skill (in this repo).
- For in-depth elvish shell knowledge (arg-completer, complex-candidate, matcher, editor API, styled text, language fundamentals, modules, rc.elv), use the **elvish** skill (in this repo).
- For in-depth xonsh shell knowledge (RichCompletion, CommandContext, contextual_command_completer, completer pipeline, prompt-toolkit integration, Python/shell hybrid model, xontribs, event system, rc.xsh), use the **xonsh** skill (in this repo).
- For in-depth fish shell knowledge (complete builtin, commandline builtin, pager, autosuggestions, syntax highlighting, abbreviations, key bindings, language fundamentals, job control, startup configuration), use the **fish** skill (in this repo).
- For in-depth PowerShell knowledge (Register-ArgumentCompleter, PSReadLine, CompletionResult, CommandAst, tab completion, menu completion, prediction plugins, $PSStyle ANSI rendering, quoting rules, argument passing to native commands, profiles and configuration), use the **powershell** skill (in this repo).
- For in-depth tcsh shell knowledge (complete builtin, COMMAND_LINE, programmable completion, bindkey, command-line editor, quoting/expansion, execution model, startup files, addsuffix, autolist, wordchars, tenematch), use the **tcsh** skill (in this repo).
- For in-depth Oil/OSH/YSH shell knowledge (programmable completion, COMP_ARGV, compadjust, compexport, OILS_COMP_UI, NiceDisplay, Readline, quoting model, simple word evaluation, lastpipe, strict:all, ysh:all, headless mode, parser-as-library, startup configuration), use the **oil** skill (in this repo).
