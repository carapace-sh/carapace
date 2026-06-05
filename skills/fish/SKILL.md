---
name: fish
description: >
  Use when working with fish shell internals — completion system, complete builtin,
  commandline builtin, pager, autosuggestions, syntax highlighting, abbreviations,
  key bindings, editing modes, language fundamentals, variables, functions, job control,
  startup configuration, or fish helper functions. Triggers on: "fish", "fish completion",
  "fish complete builtin", "fish commandline", "fish pager", "fish autosuggestion",
  "fish abbreviation", "fish key binding", "fish config", "fish_function_path",
  "fish_complete_path", "__fish_seen_subcommand_from", "__fish_contains_opt",
  "fish_mode_prompt", "fish_config", "conf.d", "argparse", "string builtin",
  "set builtin", "universal variable", "wrap completions", "NO_SPACE", "REPLACES_TOKEN".
user-invocable: true
---

# Fish Shell In-Depth Reference

Comprehensive reference for [fish](https://fishshell.com/) shell internals, with emphasis on the completion system and how external tools hook into it.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| complete, complete builtin, commandline, completion function, autoloading, conditionals, wrapping, completion output format, NO_SPACE, REPLACES_TOKEN, AUTO_SPACE, DONT_ESCAPE, DONT_SORT, SUPPRESS_PAGER_PREFIX, Completion struct, CompletionReceiver, sort_and_prioritize, fuzzy match, __fish_seen_subcommand_from, __fish_contains_opt, __fish_complete_directories, __fish_complete_path, __fish_print_hostnames, fish_complete_path, -C/--do-complete, --escape, --keep-order, --no-files, --force-files, --require-parameter, --exclusive, --wraps, --condition, --arguments, --old-option, --short-option, --long-option, option styles, completion insertion, token replacement, prefix suppression, dedup, completion engine, perform_for_commandline_impl, complete_cmd, complete_param_expand, wildcard_complete, wrap chain, WRAPPER_MAP | [references/completion.md](references/completion.md) |
| key binding, bind, editing mode, emacs, vi, autosuggestion, pager, syntax highlighting, abbreviation, abbr, kill ring, multiline editing, fish_key_reader, fish_cursor_*, fish_color_*, fish_pager_color_*, fish_autosuggestion_enabled, bracketed paste, ReadlineCmd, reader, screen rendering, dual-buffer, input event, escape sequence, CharEvent, KeyEvent, fish_escape_delay_ms, Ctrl-S search, transient prompt, fish_mode_prompt, fish_right_prompt, fish_title, fish_greeting, directory history, prevd, nextd, cdh, dirh | [references/editor.md](references/editor.md) |
| quoting, single quotes, double quotes, escape character, variable expansion, command substitution, brace expansion, wildcard expansion, globbing, lists, arrays, set, variable scope, universal variable, global, local, function scope, PATH variable, slices, index ranges, functions, autoload, $argv, control flow, if, switch, while, for, begin, and, or, combiners, test, string, math, argparse, read, count, contains, status, printf, echo, process substitution, psub, redirection, pipe, special variables, $status, $pipestatus, $fish_pid, $CMD_DURATION, $history | [references/language.md](references/language.md) |
| job control, pipeline, process group, foreground, background, fg, bg, jobs, wait, disown, signals, SIGINT, SIGHUP, SIGTERM, fork, exec, posix_spawn, internal process, External, Builtin, Function, BlockNode, Exec, process type, job lifecycle, TtyHandoff, job-control, topic monitor, iothread, thread pool, fd monitor, universal variable sync, EnvStack, EnvUniversal, execution context, exec_job, populate_job, job group, pipeline node, cancel signal | [references/execution.md](references/execution.md) |
| config.fish, conf.d, fish_function_path, fish_complete_path, XDG, __fish_config_dir, __fish_sysconf_dir, __fish_user_data_dir, startup order, login shell, interactive shell, fish_config, theme, prompt, fish_add_path, fish_user_paths, universal variables, fish_variables, vendor_completions.d, vendor_functions.d, vendor_conf.d, generated_completions, man page completions, environment variables, fish_pid, fish_version, EUID, fish_private_mode, fish_history, event handlers, fish_prompt, fish_preexec, fish_postexec, fish_exit, fish_cancel, fish_focus_in, fish_focus_out, --private, --no-config, --init-command | [references/startup-config.md](references/startup-config.md) |

## Quick Guide

- **How do I write a completion function?** → [references/completion.md](references/completion.md)
- **How does the `complete` builtin work?** → [references/completion.md](references/completion.md)
- **How does `commandline` provide context?** → [references/completion.md](references/completion.md)
- **How do conditionals (`-n`) work in completions?** → [references/completion.md](references/completion.md)
- **How does wrapping (`-w`) inherit completions?** → [references/completion.md](references/completion.md)
- **How does the pager display completions?** → [references/completion.md](references/completion.md) and [references/editor.md](references/editor.md)
- **How does NO_SPACE / AUTO_SPACE work?** → [references/completion.md](references/completion.md)
- **How do autosuggestions work?** → [references/editor.md](references/editor.md)
- **How do I configure key bindings?** → [references/editor.md](references/editor.md)
- **How do abbreviations work?** → [references/editor.md](references/editor.md)
- **How does syntax highlighting work?** → [references/editor.md](references/editor.md)
- **How do fish variables and scoping work?** → [references/language.md](references/language.md)
- **How do fish lists differ from bash arrays?** → [references/language.md](references/language.md)
- **How does command substitution split output?** → [references/language.md](references/language.md)
- **How does job control work in fish?** → [references/execution.md](references/execution.md)
- **Which config files does fish read at startup?** → [references/startup-config.md](references/startup-config.md)
- **Where do I install completion files?** → [references/startup-config.md](references/startup-config.md)
- **How do universal variables persist?** → [references/startup-config.md](references/startup-config.md)

## Cross-Project References

For carapace-specific fish integration (snippet, open-quote retry, value formatting, nospace limitation), see the **carapace-dev** skill → `references/shell-fish.md`.
