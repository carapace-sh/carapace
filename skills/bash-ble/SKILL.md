---
name: bash-ble
description: >
  Use when working with ble.sh (Bash Line Editor) internals — completion system, cand/yield,
  menu styles, progcomp integration, syntax highlighting, faces, keymaps, widgets, blerc,
  bleopt, ble-bind, blehook, sabbrev, dabbrev, auto-complete, menu-filter. Triggers on:
  "ble.sh", "bash-ble", "ble", "BLE", "ble/complete", "cand/yield", "bleopt", "ble-face",
  "ble-bind", "blehook", "blerc", "sabbrev", "dabbrev", "ble-import", "menu-filter",
  "complete_menu_style", "ble/no-default", "ble/filter-by-prefix".
user-invocable: true
---

# Bash BLE (ble.sh) In-Depth Reference

Comprehensive reference for [ble.sh](https://github.com/akinomyoga/ble.sh) (Bash Line Editor) internals, with emphasis on the completion system and how external tools hook into it.

## Sub-Resources

Load the reference that matches your task. When in doubt, load multiple references.

| Keywords | Reference |
|----------|----------|
| ble/complete, cand/yield, completion source, completion action, progcomp, COMPREPLY, compopt, menu style, menu-filter, auto-complete, sabbrev, dabbrev, complete_menu_style, ble/no-default, ble/filter-by-prefix, ble/syntax-raw, ble/no-mark-directories, ble/prog-trim, ble/default, COMPV, COMP1, COMP2, COMPS, comp_type, cand_pack, cand_cand, cand_word, quote-insert, action:plain, action:word, action:file, action:command, action:progcomp, action:mandb, complete_load, complete_limit, complete_ambiguous, complete_polling_cycle, complete_timeout_auto, complete_auto_complete, complete_auto_delay, complete_auto_menu, complete_menu_complete, complete_menu_filter, complete_menu_maxlines, menu_align_min, menu_align_max, menu_prefix, menu_desc_multicolumn_width, menu_complete_selected, menu_filter_fixed, menu_filter_input, menu_complete_match | [references/completion.md](references/completion.md) |
| syntax highlighting, syntax parser, ble-face, gspec, face, syntax_default, syntax_command, syntax_quoted, syntax_quotation, syntax_escape, syntax_varname, syntax_comment, syntax_error, syntax_delimiter, syntax_param_expansion, syntax_history_expansion, syntax_function_name, syntax_glob, syntax_brace, syntax_tilde, syntax_document, command_builtin, command_alias, command_function, command_file, command_keyword, filename_directory, filename_executable, filename_link, varname_unset, varname_export, varname_array, _ble_syntax_tree, _ble_syntax_attr, _ble_syntax_stat, _ble_syntax_nest, highlight_syntax, highlight_filename, highlight_variable, color_scheme, LS_COLORS | [references/syntax-highlighting.md](references/syntax-highlighting.md) |
| editor, widget, keymap, ble-bind, emacs, vi_imap, vi_nmap, vi_omap, vi_xmap, key binding, _ble_edit_str, _ble_edit_ind, _ble_edit_mark, kill ring, yank, undo, overwrite mode, numeric argument, keymap stack, decode, input processing, C-x, M-x, S-x, __default__, __defchar__, __before_widget__, __after_widget__, canvas, terminal rendering, DRAW_BUFF, ble/canvas, ble/textarea, sigwinch, tab_width, emoji_width, grapheme_cluster | [references/editor.md](references/editor.md) |
| installation, startup, blerc, bleopt, blehook, ble-import, XDG, _ble_base, _ble_base_run, _ble_base_cache, _ble_base_state, --attach, --rcfile, ble-attach, PROMPT_COMMAND, complete_load hook, keymap_emacs hook, keymap_vi hook, PRECMD, PREEXEC, POSTEXEC, import_path, config/readline, integration/fzf | [references/startup-config.md](references/startup-config.md) |

## Quick Guide

- **How do I write a BLE completion function?** → [references/completion.md](references/completion.md)
- **How does ble/complete/cand/yield work?** → [references/completion.md](references/completion.md)
- **How does BLE integrate with bash programmable completion?** → [references/completion.md](references/completion.md)
- **What are the BLE menu styles?** → [references/completion.md](references/completion.md)
- **How does menu-filter work?** → [references/completion.md](references/completion.md)
- **How does auto-complete work?** → [references/completion.md](references/completion.md)
- **What are sabbrev and dabbrev?** → [references/completion.md](references/completion.md)
- **What BLE compopt extensions exist?** → [references/completion.md](references/completion.md)
- **How does syntax highlighting work?** → [references/syntax-highlighting.md](references/syntax-highlighting.md)
- **How do I configure faces/colors?** → [references/syntax-highlighting.md](references/syntax-highlighting.md)
- **How do I bind keys?** → [references/editor.md](references/editor.md)
- **What widgets are available?** → [references/editor.md](references/editor.md)
- **How does the editor replace Readline?** → [references/editor.md](references/editor.md)
- **How do I install and configure ble.sh?** → [references/startup-config.md](references/startup-config.md)
- **What is the blerc file?** → [references/startup-config.md](references/startup-config.md)
- **How do bleopt/blehook/ble-import work?** → [references/startup-config.md](references/startup-config.md)

## Cross-Project References

- For carapace-specific BLE integration (snippet, value format, per-candidate suffix), see the **carapace-dev** skill → `references/shell-bash-ble.md`.
- For standard bash completion (COMP_* variables, complete builtin, compgen, compopt, Readline), see the **bash** skill.
