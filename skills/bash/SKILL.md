---
name: bash
description: >
  Use when working with bash shell internals — programmable completion, Readline,
  quoting/expansion, execution model, startup files, prompt hooks, or bash-completion
  helpers. Triggers on: "bash", "bash completion", "bash readline", "bash quoting",
  "bash expansion", "bash startup", "bash prompt", "COMP_WORDBREAKS", "COMPREPLY",
  "complete builtin", "compgen", "compopt", "bash-completion", "inputrc", "bind builtin".
user-invocable: true
---

# Bash Shell In-Depth Reference

Comprehensive reference for bash shell internals, with emphasis on the completion system and how external tools hook into it.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| complete, compgen, compopt, COMPREPLY, COMP_WORDS, COMP_CWORD, COMP_LINE, COMP_POINT, COMP_TYPE, COMP_WORDBREAKS, compspec, completion function, _init_completion, _comp_initialize, _completion_loader, dynamic loading, exit status 124, bash-completion helpers, _filedir, _comp_compgen, _comp_split, _comp_quote, _comp_dequote, _comp_reassemble, _comp_xfunc, completion registration, completion options, nospace, filenames, noquote, nosort, plusdirs, bashdefault, default, dirnames, fullquote, menu-complete | [references/completion.md](references/completion.md) |
| Readline, .inputrc, bind, key binding, editing mode, emacs, vi, completion variables, mark-directories, show-all-if-ambiguous, colored-stats, visible-stats, completion-prefix-display-length, colored-completion-prefix, bell-style, completion-display-width, page-completions, completion-query-items, menu-complete-display-prefix, skip-completed-text, macro, kill-ring, yank, bracketed-paste, READLINE_LINE, READLINE_POINT, READLINE_MARK, READLINE_ARGUMENT | [references/readline.md](references/readline.md) |
| quoting, single quotes, double quotes, ANSI-C quoting, $'...', escape character, backslash, word splitting, IFS, pathname expansion, globbing, brace expansion, tilde expansion, parameter expansion, command substitution, arithmetic expansion, process substitution, COMP_WORDBREAKS quoting, dequoting, _comp_dequote, _comp_quote, FIGNORE, nullglob, dotglob | [references/quoting-expansion.md](references/quoting-expansion.md) |
| execution, subshell, fork, exec, process group, pipeline, job control, signal, trap, DEBUG trap, EXIT trap, ERR trap, RETURN trap, PROMPT_COMMAND, PS0, PS1, PS2, PS3, PS4, prompt escape, foreground, background, SIGINT, SIGHUP, SIGTTIN, SIGTTOU, disown, huponexit, lastpipe, pipefail, command_not_found_handle | [references/execution.md](references/execution.md) |
| startup, .bashrc, .bash_profile, .profile, /etc/profile, login shell, interactive shell, non-interactive, BASH_ENV, ENV, --login, --norc, --noprofile, sh invocation, POSIX mode, remote shell, SSH, elevated privileges, shell options, shopt, set | [references/startup.md](references/startup.md) |

## Quick Guide

- **How do I write a completion function?** → [references/completion.md](references/completion.md)
- **How does Readline handle completion display?** → [references/readline.md](references/readline.md)
- **How does quoting affect completions?** → [references/quoting-expansion.md](references/quoting-expansion.md)
- **How does bash execute commands and handle signals?** → [references/execution.md](references/execution.md)
- **Which startup files does bash read?** → [references/startup.md](references/startup.md)
- **What are the COMP_* variables?** → [references/completion.md](references/completion.md)
- **How does the bash-completion project work?** → [references/completion.md](references/completion.md)
- **How do I configure .inputrc for completion?** → [references/readline.md](references/readline.md)
- **How does COMP_WORDBREAKS break words?** → [references/completion.md](references/completion.md) and [references/quoting-expansion.md](references/quoting-expansion.md)
- **How do I control nospace/filenames/noquote?** → [references/completion.md](references/completion.md)
- **How do traps and PROMPT_COMMAND work?** → [references/execution.md](references/execution.md)
- **How do I debug completion functions?** → [references/completion.md](references/completion.md)

## Cross-Project References

For carapace-specific bash integration (snippet, patch phase, value formatting), see the **carapace-dev** skill → `references/shell-bash.md`.
