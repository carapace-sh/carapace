---
name: oil
description: >
  Use when working with Oil shell (OSH/YSH) internals — programmable completion, Readline,
  quoting/expansion, execution model, startup files, COMP_ARGV, compadjust, compexport,
  OILS_COMP_UI, headless mode, simple word evaluation, strict options, or parser-as-library
  architecture. Triggers on: "oil", "osh", "ysh", "oils", "oil shell", "OSH", "YSH",
  "oil completion", "COMP_ARGV", "compadjust", "compexport", "OILS_COMP_UI",
  "simple word eval", "strict:all", "ysh:all", "lastpipe", "headless mode", "FANOS",
  "NiceDisplay", "parser as library".
user-invocable: true
---

# Oil Shell In-Depth Reference

Comprehensive reference for Oil shell (OSH/YSH) internals, with emphasis on the completion system and how external tools hook into it.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| complete, compgen, compopt, compadjust, compexport, COMPREPLY, COMP_ARGV, COMP_WORDS, COMP_CWORD, COMP_LINE, COMP_POINT, COMP_WORDBREAKS, compspec, completion function, exit status 124, bash-completion, quoting responsibility, argv entries, parser as library, OILS_COMP_UI, NiceDisplay, MinimalDisplay, headless mode, FANOS, glob pattern specs, shell language completion | [references/completion.md](references/completion.md) |
| Readline, .inputrc, key binding, editing mode, emacs, vi, completion display, NiceDisplay, MinimalDisplay, horizontal-scroll-mode, prompt, PS1, PS2, InteractiveLineReader, history expansion, alias handling, parser re-invocation, ANSI escape codes, description display | [references/line-editing.md](references/line-editing.md) |
| quoting, single quotes, double quotes, ANSI-C quoting, YSH strings, r'...', word splitting, IFS, simple word evaluation, parameter expansion, type model, tilde expansion, command substitution, arithmetic expansion, globbing, extended globs, brace expansion, dashglob, nullglob, strict_tilde, process substitution | [references/quoting-expansion.md](references/quoting-expansion.md) |
| execution, pipeline, lastpipe, subshell, fork, exec, process group, job control, signal, trap, DEBUG trap, ERR trap, EXIT trap, RETURN trap, PIPESTATUS, pipefail, errexit, strict:all, xtrace, xtrace_rich, SHX variables, interpreter state, dynamic scope, noforklast | [references/execution.md](references/execution.md) |
| startup, oshrc, yshrc, oshrc.d, yshrc.d, configuration, OILS_HIJACK_SHEBANG, OILS_COMP_UI, OSH_DEBUG_DIR, shell options, shopt, strict:all, ysh:upgrade, ysh:all, --norc, --rcfile, --rcdir, --headless, --eval, --debug-file, PS1 inheritance | [references/startup-config.md](references/startup-config.md) |

## Quick Guide

- **How do I write a completion function?** → [references/completion.md](references/completion.md)
- **How does Readline handle completion display?** → [references/line-editing.md](references/line-editing.md)
- **What is COMP_ARGV and how does it differ from COMP_WORDS?** → [references/completion.md](references/completion.md)
- **How does quoting affect completions?** → [references/completion.md](references/completion.md) and [references/quoting-expansion.md](references/quoting-expansion.md)
- **What is the compadjust builtin?** → [references/completion.md](references/completion.md)
- **What is the compexport builtin?** → [references/completion.md](references/completion.md)
- **How does OILS_COMP_UI affect the completion display?** → [references/completion.md](references/completion.md) and [references/line-editing.md](references/line-editing.md)
- **How does Oil handle quoting differently from bash?** → [references/quoting-expansion.md](references/quoting-expansion.md)
- **What is simple word evaluation in YSH?** → [references/quoting-expansion.md](references/quoting-expansion.md)
- **How do pipelines differ from bash (lastpipe)?** → [references/execution.md](references/execution.md)
- **What are the strict:all and ysh:all option groups?** → [references/startup-config.md](references/startup-config.md)
- **Which startup files does Oil read?** → [references/startup-config.md](references/startup-config.md)
- **How does headless mode work?** → [references/completion.md](references/completion.md)
- **How does the parser-as-library architecture work?** → [references/line-editing.md](references/line-editing.md)
- **How do I debug completion functions?** → [references/completion.md](references/completion.md)
- **How does xtrace_rich work?** → [references/execution.md](references/execution.md)

## Cross-Project References

For carapace-specific Oil integration (snippet, patch phase, value formatting, nospace indicator), see the **carapace-dev** skill → `references/shell-oil.md`.
