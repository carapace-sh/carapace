# BLE Completion System

In-depth reference for ble.sh's completion system — the multi-stage pipeline, sources, actions, `ble/complete/cand/yield`, progcomp integration, menu styles, menu-filter, auto-complete, sabbrev, and dabbrev.

## The Completion Pipeline

When the user presses TAB (or auto-complete triggers), ble.sh executes a multi-stage pipeline:

1. **Trigger Detection** — TAB key, auto-complete timer, or menu navigation activates completion
2. **Context Analysis** — The syntax parser analyzes the command line to determine what kind of completion is needed
3. **Source Generation** — Multiple sources contribute candidates independently
4. **Candidate Processing** — Actions transform candidates with proper shell quoting, suffixes, and closing quotes
5. **Menu Rendering** — Candidates are displayed in the configured menu style
6. **User Interaction** — Navigation, filtering, and selection

### Context Variables

During context analysis, ble.sh sets these variables for sources and actions:

| Variable | Description |
|----------|-------------|
| `COMP1` | Start index of the word being completed |
| `COMP2` | End index (cursor position) of the word being completed |
| `COMPS` | The word being completed (raw text from `COMP1` to `COMP2`) |
| `COMPV` | The value of the word being completed (after quote removal) |
| `comp_type` | Type of completion attempt (e.g., `auto`, `menu`, empty for manual) |
| `comps_flags` | Flags describing the quoting context (see Quote-Insert System) |
| `COMP_PREFIX` | Set by sources to the common prefix of yielded candidates |

## Completion Sources

Each source is a function `ble/complete/source:<name>/generate` that populates the `cand_pack` array.

### Built-in Sources

| Source | Function | Purpose |
|--------|----------|---------|
| `command` | `ble/complete/source:command` | Commands from PATH, builtins, functions, aliases |
| `file` | `ble/complete/source:file` | File/directory paths with glob support |
| `variable` | `ble/complete/source:variable` | Shell variable names |
| `argument` | `ble/complete/source:argument` | Context-aware argument completion (e.g., `cd` directories) |
| `wordlist` | `ble/complete/source:wordlist` | Custom word lists |
| `sabbrev` | `ble/complete/source:sabbrev` | Shell abbreviations |
| `progcomp` | `ble/complete/source:progcomp` | Bash programmable completion (`complete -F`, `complete -C`) |
| `dabbrev` | `ble/complete/source:dabbrev` | Dynamic abbreviation (history words) |
| `mandb` | `ble/complete/source:mandb` | Man page options with descriptions |

### Source Function Template

```bash
## @fn ble/complete/source:$name args...
##   @param[in] args...
##   @var[in] COMP1 COMP2 COMPS COMPV comp_type
##   @var[in] comp_filter_type
##   @var[out] COMP_PREFIX
##   @var[in,out] cand_count cand_cand cand_word cand_pack
##   @var[in,out] cand_limit_reached
function ble/complete/source:my-source {
  local cand "${_ble_complete_yield_varnames[@]/%/=}"
  ble/complete/cand/yield.initialize "word"
  for cand in "${my_candidates[@]}"; do
    [[ $cand == "$COMPV"* ]] && ble/complete/cand/yield "word" "$cand"
  done
}
```

## ble/complete/cand/yield

The primary function for registering completion candidates.

### Function Signature

```bash
## @fn ble/complete/cand/yield ACTION CAND DATA
##   @param[in] ACTION    Action type (plain, word, file, command, etc.)
##   @param[in] CAND      Raw candidate value
##   @param[in] DATA      Additional metadata (e.g., compopt flags)
##   @var[in] COMP_PREFIX
##   @var[in] flag_force_fignore
##   @var[in] flag_source_filter
```

### Internal Processing

1. Checks `flag_force_fignore` — skips if candidate matches `FIGNORE` patterns
2. Checks `flag_source_filter` — runs `ble/complete/candidates/filter#test` for prefix matching
3. Calls `ble/complete/action:"$ACTION"/initialize` to convert `CAND` into `INSERT` (the properly-quoted insertion string)
4. Computes `PREFIX_LEN` — length of common prefix between `CAND` and `COMP_PREFIX`
5. Stores candidate in parallel arrays:

| Array | Content |
|-------|---------|
| `cand_cand[icand]` | Raw candidate value |
| `cand_word[icand]` | Processed insertion value |
| `cand_pack[icand]` | Packed format: `ACTION:ncand,ninsert,PREFIX_LEN:$CAND$INSERT$DATA` |

### Batch Yield

```bash
## @fn ble/complete/cand/yield.batch action data
##   @arr[in] cands
function ble/complete/cand/yield.batch {
  local ACTION=$1 DATA=$2
  # Batch processing with polling/cancellation support
}
```

### Filename Yield

```bash
function ble/complete/cand/yield-filenames {
  local action=$1; shift
  # Handles hidden files, FIGNORE, etc.
  ble/complete/cand/yield.initialize "$action"
  ble/complete/cand/yield.batch "$action"
}
```

### Unpack

```bash
## @fn ble/complete/cand/unpack data
##   @param[in] data  # ACTION:ncand,ninsert,PREFIX_LEN:$CAND$INSERT$DATA
##   @var[out] ACTION CAND INSERT DATA PREFIX_LEN
```

## Completion Actions

Actions transform raw candidate values into properly quoted insertion strings. Each action type implements up to five member functions:

| Member Function | Purpose |
|-----------------|---------|
| `initialize` | Convert one candidate `CAND` → `INSERT` (with quoting) |
| `initialize.batch` | Batch convert candidates (for performance) |
| `complete` | Finalize unique match (add suffix, close quotes) |
| `init-menu-item` | Menu item styling (SGR codes for display) |
| `get-desc` | Description text for `desc` menu style |

### Action Types

| Action | Function | Purpose |
|--------|----------|---------|
| `plain` | `ble/complete/action:plain` | Basic text insertion with quoting |
| `word` | `ble/complete/action:word` | Word insertion with space suffix |
| `file` | `ble/complete/action:file` | File path with directory marker (`/`) |
| `file_rhs` | `ble/complete/action:file_rhs` | File right-hand-side for assignments |
| `command` | `ble/complete/action:command` | Command name insertion |
| `variable` | `ble/complete/action:variable` | Variable name with context-specific suffixes |
| `progcomp` | `ble/complete/action:progcomp` | Programmable completion handling |
| `mandb` | `ble/complete/action:mandb` | Man page descriptions |
| `literal-word` | `ble/complete/action:literal-word` | No quoting, space suffix |
| `literal-substr` | `ble/complete/action:literal-substr` | No quoting, no suffix |
| `suffix-sabbrev` | `ble/complete/action:suffix-sabbrev` | Suffix abbreviation expansion |

### action:progcomp Details

The `progcomp` action handles candidates from bash's programmable completion system:

```bash
function ble/complete/action:progcomp/complete {
  if [[ $DATA == *:filenames:* ]]; then
    ble/complete/action:file/complete
  else
    if [[ $DATA != *:ble/no-mark-directories:* && -d $CAND ]]; then
      ble/complete/action/requote-final-insert
      ble/complete/action/complete.mark-directory
    else
      ble/complete/action:word/complete
    fi
  fi

  [[ $DATA == *:nospace:* ]] && suffix=${suffix%' '}
  [[ $DATA == *:ble/no-mark-directories:* && -d $CAND ]] && suffix=${suffix%/}
}
```

The `DATA` field carries compopt flags as colon-separated entries (e.g., `:filenames:`, `:nospace:`, `:ble/no-mark-directories:`).

## Quote-Insert System

Handles shell quoting for completion candidates based on the current quoting context. The `comps_flags` variable encodes the context:

| Flag | Context | Quoting Method |
|------|---------|----------------|
| `S` | Inside `'...'` | Escape `'` as `'\''` |
| `E` | Inside `$'...'` | Escape special chars as `\x` |
| `D` | Inside `"..."` | Escape `\`, `"`, `` ` ``, `$` |
| `I` | Inside `$"..."` | Same as `D` |
| `B` | After `\` | Remove leading `\` from insert |
| `x` | Brace expansion | Special handling |
| `p` | Param expansion | Prevent identifier continuation |
| `v` | COMPV available | Use append mode |

### Quote-Insert Modes

1. **Append Mode** — When `CAND` starts with `COMPV`, only the differing suffix is inserted
2. **Preserve Mode** — Match `comps_fixed` parts of the original text
3. **Rewrite Mode** — Complete replacement of the word

## Progcomp Integration

ble.sh integrates with bash's `complete`, `compgen`, and `compopt` builtins, providing a compatibility layer that translates between bash's programmable completion and ble.sh's internal candidate system.

### COMP_* Variables Set by ble.sh

When invoking a user-defined completion function (`complete -F`), ble.sh sets:

| Variable | Description |
|----------|-------------|
| `COMP_LINE` | Full command line text |
| `COMP_POINT` | Cursor position |
| `COMP_CWORD` | Current word index |
| `COMP_WORDS` | Array of words |
| `COMP_TYPE` | Completion type (9=TAB, 33=`!`, 37=`%`, 63=`?`, 64=`@`) |
| `COMP_KEY` | Invoking key (decimal ASCII) |
| `COMPREPLY` | Return array for user functions |

### Wrapped compopt

ble.sh replaces the `compopt` builtin with its own implementation:

```bash
function ble/complete/progcomp/compopt {
  local ospec has_cmd compopt_args
  ble/complete/progcomp/compopt/.read-arguments "$@" || return 2

  if ((${#ospec[@]})); then
    local s
    for s in "${ospec[@]}"; do
      case $s in
      (-*) comp_opts=${comp_opts//:"${s:1}":/:}${s:1}: ;;
      (+*) comp_opts=${comp_opts//:"${s:1}":/:} ;;
      esac
    done
  elif [[ $has_cmd ]]; then
    builtin compopt "${compopt_args[@]}"
  fi
}
```

Options are recorded in the `comp_opts` variable (colon-separated) and passed through the `DATA` field of `cand_pack`.

### Standard Bash Options Supported

| Option | Handler |
|--------|---------|
| `-o filenames` | Maps to `action:file`, adds `:filenames:` flag |
| `-o noquote` | Sets `flag_noquote=1`, minimal quoting |
| `-o nospace` | No space suffix (`:nospace:` in DATA) |
| `-o dirnames` | Directory completion |
| `-o plusdirs` | Add directories after other completions |
| `-A type` | Action type (command, file, function, etc.) |
| `-F func` | User completion function |
| `-C cmd` | Command completion |
| `-D` / `-E` / `-I` | Default/empty/initial completion |

### ble.sh compopt Extensions

| Option | Description |
|--------|-------------|
| `ble/syntax-raw` | Insert literal value, operate on literal word not expansion |
| `ble/default` | Enable ble.sh's default completion fallback (default ON) |
| `ble/no-default` | Disable `-o ble/default` — skip fallback completion |
| `ble/no-mark-directories` | Suppress `/` suffix for directories |
| `ble/prog-trim` | Trim trailing space from `complete -C` output |
| `ble/filter-by-prefix` | Filter candidates by `$COMPV` prefix |

### COMPREPLY Processing

When a user completion function populates `COMPREPLY`:

```bash
function ble/complete/progcomp/.compgen-helper-func {
  local -a COMP_WORDS
  local COMP_LINE COMP_POINT COMP_CWORD COMP_TYPE COMP_KEY cmd cur prev
  ble/complete/progcomp/.compvar-initialize

  # Push wrapped compopt
  ble/function#push compopt 'ble/complete/progcomp/compopt "$@"'

  # Execute the user's completion function
  builtin eval '"$comp_func" "$cmd" "$cur" "$prev"' < /dev/null

  ble/function#pop compopt
}
```

After the function returns, `COMPREPLY` entries are yielded via `ble/complete/cand/yield` with the `progcomp` action type.

### Key Differences from Standard Bash

1. **Asynchronous cancellation** — Uses `_ble_builtin_read_hook` to check for user input during completion
2. **Timeout protection** — `ble/complete/progcomp/.check-limits` prevents runaway completions
3. **Batch processing** — Uses AWK for efficient bulk quote-insert operations
4. **Nested completion** — Completion functions can yield ble.sh candidates directly via `ble/complete/cand/yield`
5. **Workarounds** — Special handling for git completion functions that append trailing spaces
6. **mandb enrichment** — Command candidates are enriched with man page descriptions when available

## Menu Styles

Controlled by `bleopt complete_menu_style`:

| Style | Description |
|-------|-------------|
| `align` | Fixed-width columns with word wrapping |
| `align-nowrap` | Fixed-width columns, no wrapping (default) |
| `dense` | Single-space separated, wrapping allowed |
| `dense-nowrap` | Single-space separated, no wrapping |
| `linewise` | One candidate per line |
| `desc` | Two-column with ANSI escape descriptions |
| `desc-text` | Two-column with plain text descriptions |
| `desc-raw` | Raw escape sequences in descriptions |

### Menu Style Configuration

```bash
bleopt complete_menu_style=align-nowrap
bleopt menu_align_min=4          # Minimum column width
bleopt menu_align_max=20         # Maximum column width
bleopt menu_prefix=              # Prefix format (printf with %d for index)
bleopt menu_desc_multicolumn_width=65
```

### Switching Styles Dynamically

| Key | Action |
|-----|--------|
| `C-x right` | Next style |
| `C-x left` | Previous style |
| `C-x a` | `align-nowrap` |
| `C-x c` | `dense-nowrap` |
| `C-x d` | `desc` |
| `C-x l` | `linewise` |

### Menu State Variables

| Variable | Purpose |
|----------|---------|
| `_ble_complete_menu_items` | Complete list of menu items |
| `_ble_complete_menu_page_style` | Current style |
| `_ble_complete_menu_selected` | Selected item index (-1 = none) |
| `_ble_complete_menu_page_icons` | Rendered icon data for visible items |
| `_ble_complete_menu_page_index` | Current page number |
| `_ble_complete_menu_page_offset` | Index of first item on current page |

Each entry in `_ble_complete_menu_page_icons` has the format:
```
x0,y0,x1,y1,${#pack},${#esc1}[,bbox]:$pack$esc1
```

### Menu Construction Pipeline

1. **Initialization** — Copy menu items and class info to global state
2. **Size Calculation** — Get terminal dimensions, apply `complete_menu_maxlines`
3. **Cache Management** — Hash of `nitem,lines,cols:menu_style` for reuse
4. **Page Selection** — Style-specific `guess` function predicts target page
5. **Page Construction** — `ble/complete/menu-style:$style/construct-page` renders the page
6. **Result Storage** — Save page data to global state variables

### Menu Navigation Keys

| Key | Widget | Description |
|-----|--------|-------------|
| `RET` | `menu_complete/accept` | Accept selected |
| `C-g` | `menu_complete/cancel` | Cancel menu |
| `C-f`, `→` | `menu/forward-column` | Next in row |
| `C-b`, `←` | `menu/backward-column` | Previous in row |
| `C-n`, `↓` | `menu/forward-line` | Next line |
| `C-p`, `↑` | `menu/backward-line` | Previous line |
| `prior` | `menu/backward-page` | Previous page |
| `next` | `menu/forward-page` | Next page |
| `C-x right` | `menu_complete/switch-style +` | Next style |

### Menu Style Implementation Functions

- `ble/complete/menu-style:$style/construct-page` — Render page
- `ble/complete/menu-style:$style/guess` — Predict page for item
- `ble/complete/menu-style:desc/locate` — Navigate within layout

## Menu-Filter

The menu-filter feature allows users to narrow down visible menu items by typing while the completion menu is active.

### Configuration

```bash
bleopt complete_menu_filter=1    # Enable menu filtering (default: 1)
```

### Filter Faces

| Face | Default | Purpose |
|------|---------|---------|
| `menu_filter_fixed` | `bold` | Style for fixed filter portion |
| `menu_filter_input` | `fg=16,bg=229` | Style for actively typed filter text |

### Filter Types

Candidates are tested against multiple filter types in order:

1. `head` — Prefix match
2. `substr` — Substring match
3. `hsubseq` — Fuzzy/hierarchical subsequence
4. `subseq` — Subsequence match

### Implementation

- **Idle handler** (`ble/complete/menu-filter.idle`) waits for user input and applies filtering
- **Highlight layer** (`ble/highlight/layer:menu_filter`) renders filter text styling
- **Cancellation** — Filter can be interrupted via polling mechanism

## Auto-Complete

Enabled with `bleopt complete_auto_complete=1` (Bash 4.0+). Automatically shows completion suggestions as you type.

### Configuration

| Option | Default | Description |
|--------|---------|-------------|
| `complete_auto_complete` | `1` | Enable auto-completion |
| `complete_auto_delay` | `100` | Delay in milliseconds before triggering |
| `complete_auto_history` | `1` | Enable history-based suggestions |
| `complete_auto_menu` | empty | Auto-display menu after delay |
| `complete_auto_wordbreaks` | `$' \t\n'` | Word break characters |
| `complete_auto_complete_opts` | empty | Behavior options (see below) |

### Auto-Complete Options

Comma-separated or colon-separated list:

| Option | Description |
|--------|-------------|
| `history` | Enable history suggestions |
| `syntax` | Enable syntax-based suggestions |
| `history-disabled` | Disable history suggestions |
| `syntax-disabled` | Disable syntax suggestions |
| `syntax-unique` | Only suggest unique match |
| `suppress-inside-line` | Don't suggest when not at end of line |
| `suppress-inside-word` | Don't suggest when inside a word |

### Auto-Complete Keybindings

| Key | Widget | Description |
|-----|--------|-------------|
| Suggestion shown | `auto_complete/accept` | Accept suggestion (default: `→` or `End`) |
| Suggestion shown | `auto_complete/reject` | Reject suggestion (default: any other key) |

## Sabbrev (Static Abbreviations)

Zsh-style abbreviation expansion. Abbreviations expand to longer text when Space is pressed.

### Definition

```bash
ble-sabbrev KEY=VALUE            # Word sabbrev (default)
ble-sabbrev -w KEY=VALUE        # Word sabbrev (explicit)
ble-sabbrev -m KEY=COMMAND       # Dynamic sabbrev (generates completions via COMPREPLY)
ble-sabbrev -i KEY=VALUE        # Inline sabbrev (triggered anywhere in word)
ble-sabbrev -l KEY=VALUE        # Line sabbrev (triggered at line beginning)
ble-sabbrev -s KEY=VALUE        # Suffix sabbrev (triggered when followed by `*=`)
```

### Sabbrev Types

| Type | Flag | Trigger |
|------|------|---------|
| Wordwise | `-w` | At word boundaries |
| Inline | `-i` | Anywhere within the word |
| Line | `-l` | At line beginning |
| Suffix | `-s` | When followed by `*=` |
| Menu | `-m` | Generates completion candidates via `COMPREPLY` |

### Examples

```bash
ble-sabbrev L='| less'
ble-sabbrev '\date'='date +%F'
ble-sabbrev "~mybin=$HOME/bin"
ble-sabbrev -s g='| grep'
ble-sabbrev -m cmd='compgen -W "$(commands)"'
```

### Trigger Keys

- `SP` — `magic-space` widget performs sabbrev expansion
- `/` — `magic-slash` for named directory expansion
- `M-'` / `C-x '` (Emacs) or `C-]` (Vim) — `sabbrev-expand`

### Configuration

```bash
bleopt edit_magic_expand=history:sabbrev  # Enable expansions
bleopt edit_magic_accept=sabbrev          # Expand on Enter
bleopt complete_source_sabbrev_opts=      # Options like no-empty-completion
bleopt complete_source_sabbrev_ignore=    # Colon-separated patterns to ignore
```

## Dabbrev (Dynamic Abbreviations)

Searches command history for words matching the current prefix. Each invocation cycles to the next match.

### Keybindings

| Key | Action |
|-----|--------|
| `C-r` | Next match |
| `C-s` | Previous match |
| `C-g`, `C-x C-g`, `C-M-g` | Cancel |
| `RET`, `C-m` | Accept and exit |
| `C-RET`, `C-j` | Accept line and execute |

### Widget

```bash
ble/widget/dabbrev-expand    # Trigger dabbrev expansion
```

### State Variables

| Variable | Purpose |
|----------|---------|
| `_ble_complete_dabbrev_original` | Original search word |
| `_ble_complete_dabbrev_regex1` | Regex pattern for matching |
| `_ble_complete_dabbrev_regex2` | Extended regex for word boundaries |
| `_ble_complete_dabbrev_index` | Current history entry index |
| `_ble_complete_dabbrev_pos` | Position in current entry |
| `_ble_complete_dabbrev_stack` | Stack for previous states (undo) |

## Completion Key Bindings

| Key | Widget | Description |
|-----|--------|-------------|
| `C-i`, `TAB` | `complete` | Start completion |
| `M-?` | `complete show_menu` | Show completion list |
| `M-*` | `complete insert_all` | Insert all completions |
| `C-TAB` | `menu-complete` | Start menu-complete |
| `M-/` | `complete context=filename` | Complete filenames |
| `M-~` | `complete context=username` | Complete users |
| `M-$` | `complete context=variable` | Complete variables |
| `M-@` | `complete context=hostname` | Complete hosts |
| `M-!` | `complete context=command` | Complete commands |

## bleopt Completion Options Reference

### Menu Display

| Option | Default | Description |
|--------|---------|-------------|
| `complete_menu_style` | `align-nowrap` | Menu style |
| `complete_menu_maxlines` | `-1` | Max menu height (-1 = unlimited) |
| `complete_menu_complete` | `1` | Enable menu-complete on TAB |
| `complete_menu_complete_opts` | `insert-selection` | Options: `hidden`, `insert-selection`, `enter_menu` |
| `complete_menu_filter` | `1` | Enable menu filtering |
| `complete_menu_color` | `on` | Enable menu item coloring |
| `complete_menu_color_match` | `on` | Highlight matching prefix |

### Style-Specific

| Option | Default | Description |
|--------|---------|-------------|
| `menu_align_min` | `4` | Min column width for align style |
| `menu_align_max` | `20` | Max column width for align style |
| `menu_prefix` | empty | Prefix format (`%d` = 1-indexed number) |
| `menu_desc_multicolumn_width` | `65` | Width threshold for desc multicolumn mode |

### Behavior

| Option | Default | Description |
|--------|---------|-------------|
| `complete_limit` | empty | Max candidates (empty = unlimited) |
| `complete_limit_auto` | `2000` | Max for auto-completion |
| `complete_limit_auto_menu` | `100` | Max for auto menu display |
| `complete_timeout_auto` | `5000` | Timeout for auto-completion (ms) |
| `complete_ambiguous` | `1` | Enable ambiguous completion |
| `complete_contract_function_names` | `1` | Contract function name prefixes |
| `complete_requote_threshold` | `0` | Min char savings for requote |
| `complete_skip_matched` | `on` | Skip already matched candidates |
| `complete_allow_reduction` | empty | Allow text rewriting |
| `complete_polling_cycle` | `50` | Check for user input every N iterations |
| `complete_stdin_frequency` | linked | How often to check stdin |

### Auto-Complete

| Option | Default | Description |
|--------|---------|-------------|
| `complete_auto_complete` | `1` | Enable auto-completion |
| `complete_auto_delay` | `1` | Delay in ms |
| `complete_auto_history` | `1` | From history |
| `complete_auto_menu` | empty | Auto-menu delay |
| `complete_auto_wordbreaks` | `$' \t\n'` | Word breaks |
| `complete_auto_complete_opts` | empty | Behavior options |

### Sabbrev

| Option | Default | Description |
|--------|---------|-------------|
| `complete_source_sabbrev_opts` | empty | Options like `no-empty-completion` |
| `complete_source_sabbrev_ignore` | empty | Colon-separated patterns to ignore |

## Writing Custom Completion Functions

### Simple Wordlist

```bash
# ~/.blerc
function my/complete-mycmd {
  local -a my_words=(foo bar baz qux)
  ble/complete/source:wordlist "${my_words[@]}"
}
blehook/eval-after-load complete my/complete-mycmd
```

### Custom Source with cand/yield

```bash
function ble/complete/source:mycustom/generate {
  local action=word
  local cand
  ble/complete/cand/yield.initialize "$action"

  for cand in alpha beta gamma delta; do
    [[ $cand == "$COMPV"* ]] && ble/complete/cand/yield "$action" "$cand"
  done
}
```

### Using Progcomp

```bash
function ble/complete/source:myservice/generate {
  ble/complete/source:progcomp/generate
}
```

### Integration with External Frameworks

ble.sh integrates with:
- **bash-completion** — Loads completion specs via `_comp_load`
- **fzf-completion** — Fuzzy finder integration
- **carapace** — Shell completion framework (uses `ble/complete/cand/yield` with `mandb` action)
- **zoxide** — Smart directory jumping
- **atuin** — Shell history database

The progcomp integration layer (`ble/complete/source:progcomp`) handles:
- `_comp_command_offset` for command wrappers
- `complete -D`, `-E`, `-I` for default/empty/initial completions
- Dynamic loading of completion functions

## References

- [ble.sh Completion Manual](https://github.com/akinomyoga/ble.sh/wiki/Manual-%C2%A77-Completion)
- [ble.sh Completion System (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/5-completion-system)
- [Completion Architecture (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/5.1-completion-architecture)
- [Completion Sources and Actions (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/5.2-completion-sources-and-actions)
- [Menu System (DeepWiki)](https://deepwiki.com/akinomyoga/ble.sh/5.3-menu-system)
