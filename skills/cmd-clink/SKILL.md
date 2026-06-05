---
name: cmd-clink
description: >
  Use when working with clink or cmd.exe shell internals — clink completion system,
  argmatcher Lua API, match generators, Readline integration, input line coloring,
  auto-suggestions, prompt filtering, DLL injection, cmd.exe integration, doskey aliases,
  or clink configuration. Triggers on: "clink", "cmd-clink", "cmd.exe completion",
  "clink argmatcher", "clink generator", "clink prompt", "clink suggester",
  "clink classifier", "match_builder", "line_state", "clink inject", "doskey",
  "clink-select-complete", "clink popup", ".inputrc clink", "clink settings".
user-invocable: true
---

# Clink (cmd-clink) In-Depth Reference

Comprehensive reference for [clink](https://chrisant996.github.io/clink/) — the tool that injects GNU Readline into Windows `cmd.exe` to provide rich completion, history, and line-editing capabilities. Emphasis on the completion system and how external tools hook into it.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| argmatcher, clink.argmatcher(), addarg, addflags, chaincommand, loop, setclassifier, setdelayinit, onarg, onadvance, onlink, onalias, linking parsers, loopchars, nowordbreakchars, hideflags, nofiles, adddescriptions, fromhistory, delayinit, user_data, shared_user_data, shorthand form, flag prefix, end of flags | [references/argmatcher.md](references/argmatcher.md) |
| generator, clink.generator(), generate, match_builder, addmatch, addmatches, setnosort, setvolatile, setsuppressappend, setappendcharacter, setforcequoting, setfullyqualify, line_state, getline, getword, getwordinfo, getcursor, getcommandoffset, getwordbreakinfo, onfiltermatches, ondisplaymatches, completion pipeline, match generation order, match filtering, match display | [references/completion.md](references/completion.md) |
| Readline, .inputrc, key binding, editing mode, emacs, vi, bind, completion commands, clink-select-complete, old-menu-complete, menu-complete, popup, clink-show-help, clink-reload, clink-diagnostics, clink-copy-line, Readline variables, clink-specific commands, luafunc, keyseq, escape codes, keymap | [references/line-editing.md](references/line-editing.md) |
| cmd.exe, injection, DLL, clink inject, autorun, clink.bat, profile directory, settings, clink set, environment variables, doskey, PROMPT, command separator, redirection, built-in commands, clink_start.cmd, startup order, CLINK_PROFILE, CLINK_NOAUTORUN, console, terminal, ANSI, prompt filter, promptfilter, right-side prompt, transient prompt, promptcoroutine, .clinkprompt, prompt theme | [references/cmd-integration.md](references/cmd-integration.md) |
| classifier, word_classifications, classifyword, applycolor, color.arg, color.flag, color.cmd, color.doskey, color.argmatcher, color.executable, color.unrecognized, color.unexpected, color.cmdsep, color.cmdredir, color.suggestion, color.input, match.coloring_rules, LS_COLORS, color theme, .clinktheme, color settings | [references/coloring.md](references/coloring.md) |
| autosuggest, suggester, clink.suggester(), suggestion list, inline suggestion, F2, autosuggest.strategy, autosuggest.enable, autosuggest.inline, rl_buffer:hassuggestion, rl_buffer:insertsuggestion, match_prev_cmd, history suggestion, completion suggestion | [references/suggestions.md](references/suggestions.md) |
| history, clink_history, history.save, history.shared, history.max_lines, history.dupe_mode, history.ignore_space, history.expand_mode, history.auto_expand, history.sticky_search, history.time_stamp, rl.gethistorycount, rl.gethistoryitems, clink.onhistory, history compact, history expansion, event designators, word designators, F7, clink-popup-history | [references/history.md](references/history.md) |
| Lua script, script directory, completions directory, lazy loading, clink.path, CLINK_PATH, installscripts, onbeginedit, onendedit, onfilterinput, oncommand, onaftercommand, oninputlinechanged, oninject, onprovideline, clink.reload, require, Lua 5.2, clink.info | [references/scripting.md](references/scripting.md) |

## Quick Guide

- **How do I write an argmatcher?** → [references/argmatcher.md](references/argmatcher.md)
- **How do completions and match generators work?** → [references/completion.md](references/completion.md)
- **How do I configure key bindings and .inputrc?** → [references/line-editing.md](references/line-editing.md)
- **How does clink inject into cmd.exe?** → [references/cmd-integration.md](references/cmd-integration.md)
- **How do I color the input line?** → [references/coloring.md](references/coloring.md)
- **How do auto-suggestions work?** → [references/suggestions.md](references/suggestions.md)
- **How does history work?** → [references/history.md](references/history.md)
- **How do I write and load Lua scripts?** → [references/scripting.md](references/scripting.md)
- **How do I create a custom prompt?** → [references/cmd-integration.md](references/cmd-integration.md)
- **What are the match_builder methods?** → [references/completion.md](references/completion.md)
- **How do I filter or modify completions?** → [references/completion.md](references/completion.md)
- **How does clink-select-complete work?** → [references/line-editing.md](references/line-editing.md)
- **How do I set up clink autorun?** → [references/cmd-integration.md](references/cmd-integration.md)
- **What color settings are available?** → [references/coloring.md](references/coloring.md)

## Cross-Project References

- For carapace-specific cmd-clink integration (snippet, patch phase, value formatting), see the **carapace-dev** skill → `references/shell.md` (cmd-clink row in the supported shells table).
