# BLE Startup and Configuration

In-depth reference for ble.sh's installation, initialization, configuration files, and the bleopt/blehook/ble-import systems.

## Installation

### Build from Source (Recommended)

```bash
git clone https://github.com/akinomyoga/ble.sh
cd ble.sh
make install
```

Creates `~/.local/share/blesh/ble.sh`, `~/.local/share/blesh/lib/`, and `~/.local/share/blesh/keymap/`.

### Package Managers

- **AUR** (Arch Linux): `blesh-git`
- **NixOS**: `ble.sh`
- **Guix**: `ble.sh`

### Pre-built Tarball

```bash
curl -Lo ble.sh 'https://github.com/akinomyoga/ble.sh/releases/download/nightly/ble.sh-0-nightly.tar.xz'
mkdir -p ~/.local/share/blesh; tar -xJf ble.sh -C ~/.local/share/blesh
```

### Requirements

- Bash 3.0+ (recommended: Bash 4.0+)
- POSIX utilities (`stty`, `gawk`, `tar`, `xz` for acceleration)
- UTF-8 encoding
- License: BSD 3-clause

## .bashrc Setup

### Recommended Pattern

```bash
# At top of .bashrc — load ble.sh but defer attachment
[[ $- == *i* ]] && source ~/.local/share/blesh/ble.sh --attach=none

# Your bashrc settings here...

# At end of .bashrc — attach to terminal
[[ ! ${BLE_VERSION-} ]] || ble-attach
```

The `--attach=none` flag defers attachment, which is completed by calling `ble-attach` later. This ensures your `.bashrc` settings (PROMPT_COMMAND, aliases, etc.) are in place before ble.sh takes over the terminal.

### Attachment Strategies

| Strategy | Description |
|----------|-------------|
| `attach` | Immediately after sourcing |
| `prompt` | Deferred via PROMPT_COMMAND |
| `none` | Manual `ble-attach` required |

## Initialization Sequence

1. **Argument Parsing** — Parse `--help`, `--version`, `--attach`, `--rcfile`, `-o BLEOPT=VALUE`
2. **Environment Adjustment** — Set safe defaults:
   - `POSIXLY_CORRECT` — saved and unset
   - `FUNCNEST` — saved and unset
   - `set -e`, `set -u`, `set -x` — disabled
   - `expand_aliases` — enabled (needed for alias expansion)
   - `LC_COLLATE=C` — consistent pattern matching
   - `IFS=$' \t\n'` — consistent word splitting
3. **Directory Initialization** — Set up XDG paths
4. **Module Loading** — Load in dependency order:
   - `util.sh` → `decode.sh` → `color.sh` → `canvas.sh` → `history.sh` → `edit.sh` → `core-syntax.sh` → `core-complete.sh` → keymaps
5. **Configuration Loading** — Source `~/.blerc`
6. **Attachment** — Attach to terminal

## XDG Directory Layout

| Variable | Path | Purpose |
|----------|------|---------|
| `_ble_base` | `~/.local/share/blesh` | Installation directory |
| `_ble_base_run` | `${XDG_RUNTIME_DIR:-/tmp}/blesh/$UID` | Ephemeral runtime |
| `_ble_base_cache` | `${XDG_CACHE_HOME:-$HOME/.cache}/blesh` | Persistent cache |
| `_ble_base_state` | `${XDG_STATE_HOME:-$HOME/.local/state}/blesh` | Persistent state |

## Configuration Files

### Search Order

| Priority | Path | Description |
|----------|------|-------------|
| 1 | `--rcfile FILE` argument | Explicit config file |
| 2 | `~/.blerc` | User config in home directory |
| 3 | `${XDG_CONFIG_HOME:-$HOME/.config}/blesh/init.sh` | XDG-compliant path |

### blerc Template

```bash
# ~/.blerc — ble.sh configuration

# --- Options ---
bleopt complete_menu_style=desc
bleopt complete_auto_delay=100
bleopt highlight_syntax=1
bleopt highlight_filename=1

# --- Faces ---
ble-face -s syntax_command fg=brown
ble-face -s syntax_quoted fg=green
ble-face -s syntax_comment fg=242
ble-face -s auto_complete fg=238,bg=254

# --- Key Bindings ---
ble-bind -f 'C-x h' 'insert-string "Hello!"'

# --- Sabbrevs ---
ble-sabbrev gs='git status'
ble-sabbrev gc='git commit'
ble-sabbrev gp='git push'

# --- Completion Hook ---
function my/complete-load-hook {
  bleopt complete_auto_delay=300
}
blehook/eval-after-load complete my/complete-load-hook
```

## bleopt — Option Configuration

The `bleopt` function manages behavioral options stored in `bleopt_*` shell variables.

### Syntax

```bash
bleopt                          # Display all options
bleopt option                   # Display specific option value
bleopt option=value             # Set option (with validation)
bleopt option:=value           # Set or create option unconditionally
bleopt -r | --reset            # Reset options to defaults
bleopt -u | --changed          # Show only modified options
bleopt -I | --initialize       # Re-run validation functions
bleopt opt+=value              # Append to colon-separated list
bleopt opt-=value              # Remove from colon-separated list
```

### Internal Structure

- `bleopt_NAME` — current value
- `_ble_opt_def_NAME` — default value
- `bleopt/check:NAME` — validation function

### Common Options

| Category | Option | Default | Description |
|----------|--------|---------|-------------|
| Input | `input_encoding` | `UTF-8` | Character encoding |
| Highlight | `highlight_syntax` | `1` | Enable syntax highlighting |
| Highlight | `highlight_filename` | `1` | Highlight filenames |
| Highlight | `highlight_variable` | `1` | Highlight variable types |
| Complete | `complete_auto_complete` | `1` | Enable auto-completion |
| Complete | `complete_menu_complete` | `1` | Enable TAB menu completion |
| Complete | `complete_menu_style` | `align-nowrap` | Menu rendering style |
| Edit | `edit_line_type` | `graphical` | Line type for beginning/end |
| Edit | `edit_vim_mode` | `emacs` | Default editing mode |
| History | `history_share` | `''` | Share history between sessions |
| Display | `prompt_eol_mark` | `↵` | End-of-line marker |
| Display | `exec_errexit_mark` | (set) | Error exit marker |

### Disabling Features

```bash
bleopt highlight_syntax=       # Disable syntax highlighting
bleopt highlight_filename=     # Disable filename highlighting
bleopt highlight_variable=     # Disable variable highlighting
bleopt complete_auto_complete= # Disable auto-complete
bleopt complete_menu_complete= # Disable menu-complete by TAB
bleopt prompt_eol_mark=''     # Disable EOF marker
bleopt exec_errexit_mark=     # Disable error exit marker
```

## blehook — Lifecycle Hooks

The `blehook` function allows registering callbacks for lifecycle events.

### Syntax

```bash
blehook HOOK_NAME function_name       # Register hook
blehook HOOK_NAME+=function_name      # Append hook
blehook -D HOOK_NAME                  # Remove hook
blehook/invoke HOOK_NAME              # Invoke hook
blehook/eval-after-load MODULE FN     # Deferred registration
```

### Core Hooks

| Hook | When Triggered |
|------|----------------|
| `PRECMD` | Before each prompt |
| `PREEXEC` | Before command execution |
| `POSTEXEC` | After command execution |
| `complete_load` | When completion system loads |
| `keymap_emacs` | When Emacs keymap loads |
| `keymap_vi` | When Vi keymap loads |
| `EDIT_BEFORE_HOOK` | Before editing operations |
| `EDIT_AFTER_HOOK` | After editing operations |
| `DECODE_HOOK` | During key decoding |

### Usage Patterns

**Hook-based configuration (ble-0.3 and earlier):**

```bash
function my/complete-load-hook {
  bleopt complete_auto_delay=300
}
blehook/eval-after-load complete my/complete-load-hook
```

**Callback-based configuration (ble >= 0.4):**

```bash
function my/set-up-completion {
  bleopt complete_auto_delay=300
}
ble-import core-complete -C 'my/set-up-completion'
```

## ble-import — Module System

The `ble-import` function loads ble.sh modules and extensions.

### Syntax

```bash
ble-import module_name           # Load module
ble-import -d module_name        # Delayed loading (background)
ble-import -C 'code' module_name # Run code after module loads
```

### Available Modules

| Module | Description |
|--------|-------------|
| `core-complete` | Core completion definitions |
| `core-syntax` | Core syntax parsing |
| `config/readline` | Readline compatibility layer |
| `integration/fzf-completion` | fzf fuzzy finder completion |
| `integration/fzf-key-bindings` | fzf key bindings |
| `vim-surround` | Vim surround plugin |
| `vim-arpeggio` | Vim arpeggio plugin |
| `vim-airline` | Vim airline plugin |

### Custom Import Paths

```bash
bleopt import_path="${XDG_DATA_HOME:-$HOME/.local/share}/blesh/local"
```

## Module Organization

| Directory | Purpose | Key Files |
|-----------|---------|-----------|
| `ble.pp` | Bootstrap/preproc | Main entry with module includes |
| `src/` | Core runtime | `util.sh`, `decode.sh`, `edit.sh`, `color.sh`, `canvas.sh`, `history.sh` |
| `lib/` | Feature modules | `core-syntax.sh`, `core-complete.sh`, `keymap.vi.sh`, `keymap.emacs.sh` |
| `contrib/` | Optional integrations | Third-party tool integrations |

## Integration with External Tools

ble.sh integrates with several external tools:

| Tool | Integration |
|------|-------------|
| **bash-completion** | Loads completion specs via `_comp_load` |
| **fzf** | `ble-import integration/fzf-completion` |
| **carapace** | Uses `ble/complete/cand/yield` with `mandb` action for descriptions |
| **zoxide** | Smart directory jumping |
| **atuin** | Shell history database |
| **Readline** | `ble-import config/readline` for `.inputrc` compatibility |

## Version Differences

### ble.sh 0.3 (Stable)

- Hook-based configuration via `blehook/eval-after-load`
- Standard completion pipeline
- Emacs and Vi modes

### ble.sh 0.4 (Development)

- Callback-based configuration via `ble-import -C`
- `ble-face` type prefixes (`gspec`, `ref`, `copy`, `sgrspec`, `ansi`)
- `command_suffix` and `command_suffix_new` faces
- Variable type highlighting faces (`varname_unset`, `varname_export`, etc.)
- `ble-face FACE:=` syntax for defining new faces
- `ble-face -d` for defining faces

## References

- [ble.sh README](https://github.com/akinomyoga/ble.sh)
- [ble.sh Wiki](https://github.com/akinomyoga/ble.sh/wiki)
- [Installation and Getting Started (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/1.1-installation-and-getting-started)
- [Configuration System (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/1.2-configuration-system)
- [Initialization and Lifecycle (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/2.1-initialization-and-lifecycle)
