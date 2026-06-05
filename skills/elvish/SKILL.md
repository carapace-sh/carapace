---
name: elvish
description: >
  Use when working with elvish shell internals — completion system, arg-completer,
  complex-candidate, matcher, complete-getopt, editor API, keybindings, prompts,
  styled text, language fundamentals, pipelines, modules, or rc.elv configuration.
  Triggers on: "elvish", "elvish completion", "elvish editor", "elvish arg-completer",
  "elvish complex-candidate", "elvish styled", "elvish matcher", "elvish complete-getopt",
  "edit:completion", "edit:complex-candidate", "edit:notify", "edit:match-prefix",
  "from-json", "rc.elv", "elvish module", "elvish pipeline", "elvish namespace".
user-invocable: true
---

# Elvish Shell In-Depth Reference

Comprehensive reference for [elvish](https://elv.sh/) shell internals, with emphasis on the completion system and how external tools hook into it.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| arg-completer, complex-candidate, matcher, complete-getopt, complete-filename, complete-dirname, filter DSL, completion UI, completion pipeline, CodeSuffix, Display, seed, RawItem, ComplexItem, PlainItem, GenerateFileNames, FilterPrefix, completion mode, smart-start, completion:binding, completion:arg-completer, completion:matcher, adaptArgGeneratorMap, completionStart, quoting pass, cook, dedup | [references/completion.md](references/completion.md) |
| edit: module, editor API, modes, insert mode, command mode, navigation mode, history mode, histlist, location, lastcmd, instant mode, listing, keybindings, binding table, global-binding, prompts, rprompt, stale prompt, hooks, before-readline, after-readline, after-command, word types, autofix, smart-enter, notify, max-height, Editor struct, cli.App, CodeArea, ComboBox, ListBox | [references/editor.md](references/editor.md) |
| styled, styled-segment, style transformer, kebab-case, ANSI color, XTerm 256-color, TrueColor, foreground, background, bold, italic, underlined, dim, inverse, blink, ParseStyling, ui.Text, ui.Segment, Styledown, render-styledown, derender-styledown, NO_COLOR, LSCOLORS, color config | [references/styling.md](references/styling.md) |
| value types, string, number, list, map, pseudo-map, nil, boolean, exception, file, function, variable, scoping, closure, upvalues, expressions, compounding, indexing, output capture, exception capture, braced list, tilde expansion, wildcard expansion, pipeline, value pipeline, byte pipeline, redirection, special commands, var, set, del, fn, if, while, for, try, and, or, coalesce, pragma, namespace, module, use, bareword, quoting, metacharacters | [references/language.md](references/language.md) |
| rc.elv, config directory, ~/.config/elvish, module system, use, epm, package manager, runtime module, path module, os module, str module, re module, math module, platform module, unix module, md module, doc module, file module, environment variables, E: namespace, paths, startup order, XDG, db.elv, history store | [references/startup-config.md](references/startup-config.md) |

## Quick Guide

- **How do I write an arg-completer?** → [references/completion.md](references/completion.md)
- **How does the completion pipeline work internally?** → [references/completion.md](references/completion.md)
- **How do I use edit:complex-candidate?** → [references/completion.md](references/completion.md)
- **How do I configure matchers?** → [references/completion.md](references/completion.md)
- **How do I use edit:complete-getopt?** → [references/completion.md](references/completion.md)
- **How do I configure keybindings?** → [references/editor.md](references/editor.md)
- **How do I customize prompts?** → [references/editor.md](references/editor.md)
- **How do editor hooks work?** → [references/editor.md](references/editor.md)
- **How do I use styled text?** → [references/styling.md](references/styling.md)
- **What style transformers are available?** → [references/styling.md](references/styling.md)
- **How do elvish pipelines work?** → [references/language.md](references/language.md)
- **How do namespaces and modules work?** → [references/language.md](references/language.md) and [references/startup-config.md](references/startup-config.md)
- **How do I set up rc.elv?** → [references/startup-config.md](references/startup-config.md)
- **How do I write a custom module?** → [references/startup-config.md](references/startup-config.md)
- **How does the filter DSL work in completion mode?** → [references/completion.md](references/completion.md)

## Cross-Project References

- For carapace-specific elvish integration (snippet, value formatting, JSON output, ParseStyling validation), see the **carapace-dev** skill → `references/shell-elvish.md`.
