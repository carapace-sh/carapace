---
name: tcsh
description: >
  Use when working with tcsh shell internals — programmable completion, complete builtin,
  COMMAND_LINE, command-line editor, bindkey, quoting/expansion, execution model, startup files,
  or tcsh-specific features. Triggers on: "tcsh", "tcsh completion", "tcsh complete",
  "tcsh bindkey", "tcsh quoting", "tcsh startup", "COMMAND_LINE", "complete builtin",
  "addsuffix", "autolist", "wordchars", "tenematch", "csh", "TENEX".
user-invocable: true
---

# Tcsh Shell In-Depth Reference

Comprehensive reference for tcsh shell internals, with emphasis on the completion system and how external tools hook into it.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| complete, uncomplete, COMMAND_LINE, p/c/n/N word spec, list types, select pattern, suffix, addsuffix, autolist, fignore, complete enhance, recexact, listmax, listmaxrows, tenematch, tw_complete, tw_result, complete-word, list-choices, complete-word-raw, complete-word-fwd, complete-word-back, external command completion, backtick completion, COMP_WORDBREAKS, wordchars, spelling correction, autocorrect, correct variable | [references/completion.md](references/completion.md) |
| bindkey, editor, emacs mode, vi mode, vimode, key binding, editing commands, wordchars, complete-word, list-choices, history-search, i-search, dabbrev-expand, expand-history, expand-glob, normalize-command, which-command, run-help, prompt, rprompt, prompt2, prompt3, visiblebell, color, ls-F | [references/editor.md](references/editor.md) |
| quoting, single quotes, double quotes, backslash, backtick, backslash_quote, variable expansion, $name, ${name}, command substitution, history expansion, !, alias expansion, globbing, filename substitution, *, ?, [...], brace expansion, tilde expansion, globstar, nonomatch, noglob, fignore, csubstnonl | [references/quoting-expansion.md](references/quoting-expansion.md) |
| execution, pipeline, job control, background, foreground, suspend, stop, notify, nohup, hup, onintr, signal, trap, autologout, exit, exec, source, repeat, foreach, while, if, switch, goto, cdpath, hash, rehash, noclobber, redirect, pipefail | [references/execution.md](references/execution.md) |
| startup, .tcshrc, .cshrc, .login, .logout, /etc/csh.cshrc, /etc/csh.login, /etc/csh.logout, login shell, interactive shell, set, setenv, unset, unsetenv, shell variables, environment variables, path, cwd, home, history, savehist, histfile, prompt, term, tcsh version, -f flag, -l flag, -m flag | [references/startup-config.md](references/startup-config.md) |

## Quick Guide

- **How do I write a completion for a command?** → [references/completion.md](references/completion.md)
- **How does an external completion generator hook into tcsh?** → [references/completion.md](references/completion.md)
- **What is COMMAND_LINE and how is it set?** → [references/completion.md](references/completion.md)
- **How does the complete builtin word specification work?** → [references/completion.md](references/completion.md)
- **How does the command-line editor handle completion?** → [references/editor.md](references/editor.md)
- **How do I bind keys for completion?** → [references/editor.md](references/editor.md)
- **How does quoting affect completions?** → [references/quoting-expansion.md](references/quoting-expansion.md)
- **How does history expansion work?** → [references/quoting-expansion.md](references/quoting-expansion.md)
- **How does tcsh execute commands and handle signals?** → [references/execution.md](references/execution.md)
- **Which startup files does tcsh read?** → [references/startup-config.md](references/startup-config.md)
- **How do I set shell and environment variables?** → [references/startup-config.md](references/startup-config.md)
- **What is the addsuffix variable?** → [references/completion.md](references/completion.md)
- **How does the complete builtin differ from bash's complete?** → [references/completion.md](references/completion.md)
- **How does spelling correction work?** → [references/completion.md](references/completion.md)
- **How do I configure the prompt?** → [references/editor.md](references/editor.md)

## Cross-Project References

For carapace-specific tcsh integration (snippet, value formatting, quoting), see the **carapace-dev** skill → `references/shell.md` (tcsh is covered in the secondary shells comparison table).
