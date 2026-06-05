---
name: xonsh
description: >
  Use when working with xonsh shell internals — completion system, RichCompletion,
  CommandContext, contextual_command_completer, completer pipeline, prompt-toolkit
  integration, key bindings, Python/shell hybrid model, subprocess execution, xontribs,
  event system, rc.xsh configuration, or environment variables. Triggers on: "xonsh",
  "xonsh completion", "xonsh RichCompletion", "xonsh CommandContext",
  "xonsh completer", "contextual_command_completer", "xonsh prompt-toolkit",
  "xonsh xontrib", "xonsh event", "xonsh rc.xsh", "add_one_completer",
  "sub_proc_get_output", "XONSH_COMPLETER_MODE", "xonsh subprocess",
  "xonsh quoting", "xonsh environment".
user-invocable: true
---

# Xonsh Shell In-Depth Reference

Comprehensive reference for [xonsh](https://xon.sh/) shell internals, with emphasis on the completion system and how external tools hook into it.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| RichCompletion, CommandContext, CommandArg, CompletionContext, contextual_command_completer, contextual_command_completer_for, non_exclusive_completer, add_one_completer, Completer class, completer pipeline, exclusive vs non-exclusive, deduplication, sorting tiers, filter function, prefix_len, append_closing_quote, append_space, display, description, style, provider, tag_provider, apply_lprefix, complete_from_sub_proc, comp_based_completer, on_completer_filter, XONSH_COMPLETER_MODE, XONSH_COMPLETER_TRACE, COMPLETION_QUERY_LIMIT, XONSH_COMPLETER_DIRS, CommandCompleter, xompletions, bash completion bridge, path completer, command completer, python completer, base completer | [references/completion.md](references/completion.md) |
| PromptToolkitShell, PromptToolkitCompleter, PromptSession, Completion, CompleteStyle, key bindings, on_ptk_create, COMPLETIONS_DISPLAY, COMPLETION_MODE, COMPLETIONS_CONFIRM, COMPLETIONS_MENU_ROWS, UPDATE_COMPLETIONS_ON_KEYPRESS, COMPLETION_IN_THREAD, AUTO_SUGGEST_IN_COMPLETIONS, XONSH_PROMPT_AUTO_SUGGEST, auto-suggest, menu-complete, display_meta, _highlight_match, unquote, ptk shell, vi mode, emacs mode, XONSH_AUTOPAIR | [references/prompt-toolkit.md](references/prompt-toolkit.md) |
| Python-powered shell, subprocess execution, SubprocSpec, CommandPipeline, capture modes, $() vs !() vs $[] vs ![], pipe redirect, a>p, e>p, subprocess syntax, Python expressions, shell expressions, @(), $(), ![], subprocess wrapping, quoting, string literals, raw strings, p-strings, f-strings, path strings, word splitting, glob expansion, tilde expansion, environment access, $VAR, ${VAR}, ${{VAR}}, alias system, callable alias, ExecAlias | [references/language-execution.md](references/language-execution.md) |
| rc.xsh, .xonshrc, XONSHRC, XONSHRC_DIR, XONSH_CONFIG_DIR, xontrib, xontrib load, xontrib list, _load_xontrib_, _unload_xontrib_, entry points, event system, Event, LoadEvent, EventManager, on_chdir, on_precommand, on_postcommand, on_command_not_found, on_pre_prompt, on_post_prompt, on_envvar_change, on_transform_command, on_exit, on_ptk_create, on_completer_filter, XSH, XSH.env, XSH.completers, XSH.commands_cache, XSH.builtins, environment variables, Env class, Var, VarPattern, detype, LsColors, EnvPath, XONSH_* env vars, BASH_COMPLETIONS, COMPLETE_DOTS, SUBSEQUENCE_PATH_COMPLETION, FUZZY_PATH_COMPLETION, SUGGEST_THRESHOLD, CMD_COMPLETIONS_SHOW_DESC, ALIAS_COMPLETIONS_OPTIONS_BY_DEFAULT | [references/startup-config.md](references/startup-config.md) |

## Quick Guide

- **How do I write a contextual completer?** → [references/completion.md](references/completion.md)
- **How does the completer pipeline work?** → [references/completion.md](references/completion.md)
- **What is RichCompletion and its fields?** → [references/completion.md](references/completion.md)
- **How do I register a completer with add_one_completer?** → [references/completion.md](references/completion.md)
- **How does exclusive vs non-exclusive work?** → [references/completion.md](references/completion.md)
- **How does the bash completion bridge work?** → [references/completion.md](references/completion.md)
- **How does prompt-toolkit render completions?** → [references/prompt-toolkit.md](references/prompt-toolkit.md)
- **How do I configure completion display?** → [references/prompt-toolkit.md](references/prompt-toolkit.md)
- **How do I add custom key bindings?** → [references/prompt-toolkit.md](references/prompt-toolkit.md)
- **How does xonsh execute subprocess commands?** → [references/language-execution.md](references/language-execution.md)
- **What are the capture modes ($() vs !())?** → [references/language-execution.md](references/language-execution.md)
- **How does quoting work in xonsh?** → [references/language-execution.md](references/language-execution.md)
- **How do I set up rc.xsh?** → [references/startup-config.md](references/startup-config.md)
- **How do I write a xontrib?** → [references/startup-config.md](references/startup-config.md)
- **How does the event system work?** → [references/startup-config.md](references/startup-config.md)
- **What are the XONSH_* environment variables?** → [references/startup-config.md](references/startup-config.md)

## Cross-Project References

- For carapace-specific xonsh integration (snippet, patch phase, value formatting, JSON output, style conversion), see the **carapace-dev** skill → `references/shell-xonsh.md`.
