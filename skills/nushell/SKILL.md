---
name: nushell
description: >
  Use when working with nushell shell internals — completion system, external completers,
  custom completions, Reedline, completion menus, spans, quoting, SyntaxShape, externs,
  $env.config, modules, or nushell's type system. Triggers on: "nushell", "nu shell",
  "nushell completion", "nushell completer", "nushell spans", "nushell extern",
  "nushell custom completion", "Reedline", "nushell menu", "nushell quoting",
  "nushell SyntaxShape", "nushell config", "nushell module", "nushell type",
  "nushell style", "nushell startup", "nushell pipeline", "NuCompleter",
  "nushell external completer", "nushell alias completion".
user-invocable: true
---

# Nushell Shell In-Depth Reference

Comprehensive reference for nushell shell internals, with emphasis on the completion system and how external tools hook into it.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| NuCompleter, completer closure, spans, external completer, custom completion, @ syntax, completion algorithm, fuzzy, prefix, substring, NuMatcher, SemanticSuggestion, Suggestion struct, append_whitespace, completion_options, CommandCompletion, FileCompletion, FlagCompletion, CustomCompletion, CommandWideCompletion, context-aware completion, completion fallback, null fallback, carapace completer, fish completer, alias workaround, CARAPACE_LENIENT, max_results | [references/completion.md](references/completion.md) |
| Reedline, Completer trait, Suggestion, Span, menu system, ColumnarMenu, ListMenu, DescriptionMenu, IdeMenu, menu layout, menu marker, menu style, menu keybindings, menu event, menu source, menu record, partial_complete, replace_in_buffer, MenuEvent, traversal direction, page_size, selection_rows, description_rows, only_buffer_difference | [references/reedline.md](references/reedline.md) |
| quoting, single quotes, double quotes, raw string, backtick, bare word, string interpolation, $'...', $\"...\", escape sequence, metacharacters, tilde expansion, path quoting, shlex, open quote, Patch | [references/quoting.md](references/quoting.md) |
| SyntaxShape, Type, Value, type system, CompleterWrapper, shape_filepath, shape_string, shape_flag, shape_internalcall, shape_external, type annotation, input/output types, subtyping, cell-path, record type, list type, table type, @ syntax, type-driven completion | [references/types.md](references/types.md) |
| extern, export extern, extern signature, flag completion, argument completion, subcommand completion, module-based extern, completion function, def completer, rest parameter, optional parameter, quoted command string, caret sigil, externs and modules | [references/externs.md](references/externs.md) |
| $env.config, completions, external completer, completion algorithm, sort, case_sensitive, max_results, env.nu, config.nu, login.nu, autoload, NU_LIB_DIRS, startup order, XDG, default_config.nu, default_env.nu, vendor-autoload, user-autoload, config nu, config env | [references/configuration.md](references/configuration.md) |
| style, color, attr, fg, bg, color_config, shape_*, hex color, named color, abbreviation, light variant, background color, LS_COLORS, nu_ansi_term, style record, completion style, menu style, selected_text, description_text | [references/style.md](references/style.md) |
| pipeline, PipelineData, execution model, two-stage, parse, evaluate, external command, caret sigil, ExternalStream, ListStream, $in, input/output types, structured data, command search, module system, use, hide, export, submodule | [references/execution.md](references/execution.md) |

## Quick Guide

- **How do I write a custom completer?** → [references/completion.md](references/completion.md) and [references/externs.md](references/externs.md)
- **How does the external completer closure work?** → [references/completion.md](references/completion.md)
- **What are spans and how are they passed?** → [references/completion.md](references/completion.md)
- **How do I configure the completion menu?** → [references/reedline.md](references/reedline.md)
- **How does Reedline handle completion display?** → [references/reedline.md](references/reedline.md)
- **How does quoting affect completions?** → [references/quoting.md](references/quoting.md)
- **What is SyntaxShape and how does it drive completions?** → [references/types.md](references/types.md)
- **How do I define extern commands with completions?** → [references/externs.md](references/externs.md)
- **Which startup files does nushell read?** → [references/configuration.md](references/configuration.md)
- **How do I set up carapace as external completer?** → [references/completion.md](references/completion.md)
- **How do completion styles and colors work?** → [references/style.md](references/style.md)
- **How does nushell execute commands and pipelines?** → [references/execution.md](references/execution.md)
- **How do I handle alias completion?** → [references/completion.md](references/completion.md)
- **What is the @ syntax for attaching completers?** → [references/types.md](references/types.md) and [references/externs.md](references/externs.md)
- **How do I configure completion matching (fuzzy/prefix)?** → [references/completion.md](references/completion.md)
- **How do I return styled completions?** → [references/completion.md](references/completion.md) and [references/style.md](references/style.md)

## Cross-Project References

For carapace-specific nushell integration (snippet, patch phase, value formatting, style conversion), see the **carapace-dev** skill → `references/shell-nushell.md`.
