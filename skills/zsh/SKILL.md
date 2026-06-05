---
name: zsh
description: >
  Use when working with zsh shell internals — completion system (compsys), compinit, compdef,
  compadd, _describe, _arguments, _alternative, _values, _multi_parts, _sep_parts,
  _regex_arguments, zstyle, matcher specification, tag system, ZLE widgets, keymaps, bindkey,
  compstate, quoting, parameter expansion, globbing, startup files, fpath, autoload, or
  compinstall. Triggers on: "zsh", "zsh completion", "zsh compinit", "zsh compdef",
  "zsh compadd", "zsh _describe", "zsh _arguments", "zsh zstyle", "zsh matcher",
  "zsh ZLE", "zsh widget", "zsh bindkey", "zsh compstate", "zsh quoting",
  "zsh expansion", "zsh globbing", "zsh startup", "zsh fpath", "zsh autoload",
  "zsh compinstall", "zsh list-colors", "zsh tag-order", "zsh group-name",
  "zsh menu-select", "zsh menu-complete", "curcontext", "words", "CURRENT",
  "PREFIX", "SUFFIX", "compquote", "comptilde", "compset".
user-invocable: true
---

# Zsh Shell In-Depth Reference

Comprehensive reference for [zsh](https://www.zsh.org/) shell internals, with emphasis on the completion system and how external tools hook into it.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| compinit, compdef, compadd, _describe, _arguments, _alternative, _values, _multi_parts, _sep_parts, _regex_arguments, _regex_words, _combination, _sequence, _path_files, _files, _gnu_generic, _wanted, _requested, _next_label, _all_labels, _tags, _message, _setup, _pick_variant, _call_function, _dispatch, _normal, _generic, _guard, _store_cache, _retrieve_cache, _cache_invalid, tag system, tag-order, group-name, group-order, zstyle, list-colors, matcher, matcher-list, completer, curcontext, context format, special parameters (words, CURRENT, PREFIX, SUFFIX, IPREFIX, ISUFFIX, QIPREFIX, QISUFFIX), opt_args, val_args, line, state, context, NORMARG, service, _compskip, expl, compinstall, bashcompinit, #compdef directive, _complete, _approximate, _correct, _expand, _history, _match, _prefix, _list, _menu, _oldlist, _ignored, _user_expand, _all_matches, how external tools hook in, completion dispatch flow, completion caching | [references/completion.md](references/completion.md) |
| ZLE, zle builtin, widget, keymap, bindkey, emacs, viins, vicmd, main, .safe, isearch, complete-word, menu-complete, menu-select, reverse-menu-complete, delete-char-or-list, list-choices, expand-or-complete, accept-and-menu-complete, compstate, compwidget, compquote, comptilde, compset, compcall, BUFFER, CURSOR, LBUFFER, RBUFFER, KEYMAP, WIDGET, KEYS, NUMERIC, CONTEXT, vared, zle_highlight, KEYTIMEOUT, user-defined widget, completion widget, zle -C, zle -N | [references/zle.md](references/zle.md) |
| quoting, single quotes, double quotes, $'...' ANSI-C quoting, backslash escaping, RC_QUOTES, parameter expansion, ${name}, ${name:-word}, ${name:=word}, ${name:+word}, ${name:?word}, ${#name}, ${name#pattern}, ${name##pattern}, ${name%pattern}, ${name%%pattern}, ${(flags)name}, expansion flags, globbing, filename generation, *, ?, [...], <m-n>, **/, glob qualifiers, extended glob, (#b), (#i), ksh glob, expansion order, tilde expansion, brace expansion, process substitution, command substitution, arithmetic expansion, WORDCHARS, FIGNORE, nullglob, nomatch | [references/expansion-quoting.md](references/expansion-quoting.md) |
| startup, .zshenv, .zshprofile, .zshrc, .zlogin, .zlogout, ZDOTDIR, login shell, interactive shell, fpath, FPATH, autoload, compinit security, compaudit, zcompdump, compinit -C, compinit -D, compinit -u, compinit -i, compinstall, RCS, GLOBAL_RCS, completion function discovery, #compdef, _compdir, completion directories, Base, Zsh, Unix, X, platform-specific completions | [references/startup-config.md](references/startup-config.md) |

## Quick Guide

- **How do I write a completion function?** → [references/completion.md](references/completion.md)
- **How does the completion dispatch flow work?** → [references/completion.md](references/completion.md)
- **How do I use _arguments with the state machine?** → [references/completion.md](references/completion.md)
- **How do I use _describe for simple completions?** → [references/completion.md](references/completion.md)
- **How do I configure zstyle for completion?** → [references/completion.md](references/completion.md)
- **How does the matcher specification work?** → [references/completion.md](references/completion.md)
- **How does the tag system work?** → [references/completion.md](references/completion.md)
- **How do I use list-colors for colored completions?** → [references/completion.md](references/completion.md)
- **How do external tools hook into zsh completion?** → [references/completion.md](references/completion.md)
- **How do I cache expensive completions?** → [references/completion.md](references/completion.md)
- **How do ZLE widgets and keymaps work?** → [references/zle.md](references/zle.md)
- **How do I create a custom ZLE widget?** → [references/zle.md](references/zle.md)
- **How do completion widgets interact with ZLE?** → [references/zle.md](references/zle.md)
- **What is compstate and how do I use it?** → [references/zle.md](references/zle.md)
- **How does zsh quoting work?** → [references/expansion-quoting.md](references/expansion-quoting.md)
- **How does parameter expansion work?** → [references/expansion-quoting.md](references/expansion-quoting.md)
- **How does zsh globbing work?** → [references/expansion-quoting.md](references/expansion-quoting.md)
- **Which startup files does zsh read?** → [references/startup-config.md](references/startup-config.md)
- **How does compinit work?** → [references/startup-config.md](references/startup-config.md)
- **How do I add custom completion functions?** → [references/startup-config.md](references/startup-config.md)

## Cross-Project References

- For carapace-specific zsh integration (snippet, quoting state machine, zstyle generation, _describe protocol, named directories), see the **carapace-dev** skill → `references/shell-zsh.md`.
