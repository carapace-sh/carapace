# Elvish Editor API

In-depth reference for elvish's interactive editor — the `edit:` module, modes, keybindings, prompts, hooks, navigation, history, and the internal editor architecture.

## Overview

The `edit:` module is the interface to elvish's interactive editor. It provides:
- **Modes** — different UI states (insert, completion, navigation, history, etc.)
- **Keybindings** — per-mode key binding tables
- **Prompts** — customizable left and right prompts
- **Hooks** — callback lists triggered at editor events
- **Completion API** — arg-completers, matchers, complex-candidate (see [references/completion.md](completion.md))
- **Navigation** — file browser mode
- **History** — command history search and listing

## Editor Architecture

### Source Code Structure

| Package | Purpose |
|---------|---------|
| `pkg/edit/editor.go` | Main `Editor` struct, initialization, `ReadCode()`, `Notify()` |
| `pkg/edit/completion.go` | Completion API, arg-completer adapter, `complexCandidate` |
| `pkg/edit/highlight.go` | Syntax highlighting, autofix |
| `pkg/edit/prompt.go` | Prompt management, stale prompt handling |
| `pkg/edit/histwalk.go` | History walking mode |
| `pkg/edit/navigation.go` | Navigation mode |
| `pkg/edit/insert.go` | Insert mode API |
| `pkg/edit/listing.go` | Listing mode (shared by histlist, location, lastcmd) |
| `pkg/edit/instant.go` | Instant mode (evaluate as you type) |
| `pkg/edit/minibuf.go` | Minibuffer |
| `pkg/cli/` | Generic CLI framework (App, TTY, widgets) |
| `pkg/cli/tk/` | Widget toolkit (CodeArea, ComboBox, ListBox, etc.) |
| `pkg/cli/modes/` | Mode implementations (completion, histwalk, etc.) |

### The Editor Struct

```go
type Editor struct {
    app          cli.App           // Underlying CLI application
    ns           *eval.Ns          // Elvish namespace for editor commands
    excMutex     sync.RWMutex      // Protects exception list
    excList      vals.List         // Stores exceptions from callbacks
    autofix      atomic.Value      // Stores autofix state
    applyAutofix func()            // Callback to apply autofix
    AfterCommand []func(...)       // Hooks run after command execution
}
```

### Initialization Sequence

The editor is initialized via `NewEditor()` in this order:

1. **Base setup** — create `Editor` struct, build `edit` namespace, configure TTY
2. **App configuration** — init max height, readline hooks, add-cmd filters, global bindings, insert API, highlighter, prompts
3. **Create CLI App** — `cli.NewApp(appSpec)`
4. **API initialization** — exceptions, variables, command API, listings, navigation, completion, histwalk, instant, minibuf, REPL
5. **Builtin initialization** — buffer builtins, TTY builtins, misc builtins, state API, store API
6. **Namespace finalization** — build namespace, execute embedded `init.elv`

### The CLI App Layer

The editor is built on a generic `cli.App` from `pkg/cli/`:

- `cli.App` manages the TTY, event loop, and addon widgets
- `tk.CodeArea` — the main input area (buffer, cursor, pending code)
- `tk.ComboBox` — combined filter + list (used by completion, histlist, location)
- `tk.ListBox` — scrollable list with horizontal/vertical layout
- Addons are pushed/popped from the app (e.g., completion mode pushes a ComboBox addon)

## Modes

Exactly one mode is active at any time. Each mode has its own UI and keybindings.

### Mode Overview

| Mode | Submodule | Description | Entry Function |
|------|-----------|-------------|---------------|
| Insert | `edit:insert:` | Default mode for typing commands | (default) |
| Command | `edit:command:` | Vi command mode (incomplete) | `edit:command:start` |
| Completion | `edit:completion:` | Shows completion candidates | `edit:completion:start` / `edit:completion:smart-start` |
| Navigation | `edit:navigation:` | File browser | `edit:navigation:start` |
| History | `edit:history:` | History walking | `edit:history:start` |
| Histlist | `edit:histlist:` | History listing | `edit:histlist:start` |
| Location | `edit:location:` | Directory browser | `edit:location:start` |
| Lastcmd | `edit:lastcmd:` | Last command arguments | `edit:lastcmd:start` |
| Instant | `edit:-instant:` | Evaluate as you type | `edit:-instant:start` |

### Listing Modes

`histlist`, `location`, and `lastcmd` are "listing modes" — they share a common UI pattern:
- A filter text input at the top
- A scrollable list below
- Support for the filter DSL (see [references/completion.md](completion.md))
- Common keybindings via `$edit:listing:binding`

### Mode Transitions

Keybindings to **start** modes live in the insert mode binding table:

```elvish
# Start location mode from insert mode
set edit:insert:binding[Alt-l] = { edit:location:start }

# History mode uses both insert and history bindings
set edit:insert:binding[Ctrl-R] = { edit:history:start }
set edit:history:binding[Ctrl-R] = { edit:history:up }
```

To close a mode and return to insert: `edit:close-mode` or `Ctrl-[`.

## Keybindings

### Binding Tables

Each mode has its own binding table:

| Variable | Mode |
|----------|------|
| `$edit:insert:binding` | Insert mode (default) |
| `$edit:command:binding` | Command mode |
| `$edit:completion:binding` | Completion mode |
| `$edit:navigation:binding` | Navigation mode |
| `$edit:history:binding` | History walking mode |
| `$edit:histlist:binding` | History listing mode |
| `$edit:location:binding` | Location mode |
| `$edit:lastcmd:binding` | Last command mode |
| `$edit:-instant:binding` | Instant mode |
| `$edit:listing:binding` | Shared listing mode bindings |
| `$edit:global-binding` | Consulted when active mode doesn't handle a key |

### Setting Keybindings

```elvish
set edit:insert:binding[Alt-x] = { exit }
set edit:insert:binding[Ctrl-P] = { edit:history:start }
set edit:global-binding[Ctrl-C] = { edit:close-mode }
```

### Key Format

**Simple characters**: lowercase for Alt, uppercase for Ctrl:
- `x` — the letter x
- `Alt-x` — Alt+x
- `X` — Shift+x (also Ctrl+X in some contexts)

**Function keys**: `F1`–`F12`, `Up`, `Down`, `Right`, `Left`, `Home`, `Insert`, `Delete`, `End`, `PageUp`, `PageDown`, `Tab`, `Enter`, `Backspace`

**Modifiers**: `A` (Alt), `C` (Ctrl), `M` (Meta), `S` (Shift)
- Stackable: `C+A-X`, `Alt+Enter`, `S-F1`

### Default Insert Mode Bindings

| Key | Action |
|-----|--------|
| `Tab` | `edit:completion:smart-start` |
| `Enter` | `edit:smart-enter` |
| `Up` | `edit:history:start` (then `edit:history:up`) |
| `Ctrl-R` | `edit:histlist:start` |
| `Ctrl-L` | `edit:location:start` |
| `Ctrl-N` | `edit:navigation:start` |
| `Ctrl-A` | Apply autofix suggestion |
| `Ctrl-[` | `edit:close-mode` |

### Parsing Keys Programmatically

```elvish
edit:key Alt-Enter  # Returns the key value for Alt+Enter
```

### Binding Table Conversion

```elvish
edit:binding-table $map  # Convert a map to a binding table
```

## Prompts

### Left-hand Prompt (`$edit:prompt`)

```elvish
set edit:prompt = { tilde-abbr $pwd; put '> ' }
```

The prompt function may write:
- **Value outputs** — strings or `styled` values, joined with no spaces
- **Byte outputs** — output as-is, including newlines

### Right-side Prompt (`$edit:rprompt`)

```elvish
set edit:rprompt = (constantly (styled (whoami)@(hostname) inverse))
```

### Stale Prompt Handling

Prompt functions run on a separate thread. If they don't finish within the threshold, the prompt is marked as **stale**:

```elvish
set edit:prompt-stale-threshold = 0.5  # seconds
set edit:prompt-stale-transform = {|x| styled $x inverse }
```

Default stale behavior: reverse-video styling.

### Prompt Eagerness

Controls when the prompt updates:

| Value | Behavior |
|-------|----------|
| `< 10` | Updated when working directory changes (default: 5) |
| `≥ 10` | Updated on each keystroke |
| Always | Updated when editor becomes active |

```elvish
set edit:-prompt-eagerness = 10  # Update on every keystroke
set edit:-rprompt-eagerness = 5  # Update on directory change
```

### RPrompt Persistency

```elvish
set edit:rprompt-persistent = $true
```

By default, the right prompt is only shown while the editor is active. Setting this to `$true` keeps it visible after accepting a command.

## Hooks

Hook variables are lists of functions executed at certain editor events.

### Available Hooks

| Variable | When Called | Arguments |
|----------|-----------|-----------|
| `$edit:before-readline` | Before editor runs | None |
| `$edit:after-readline` | After editor accepts a command | The accepted line |
| `$edit:after-command` | After shell executes a command | Map with `src`, `duration`, `error` keys |

### Adding Hooks

Append to hook variables rather than replacing them:

```elvish
set edit:before-readline = [ $@edit:before-readline { echo 'going to read' } ]
set edit:after-readline = [ $@edit:after-readline {|line| echo 'just read '$line } ]
set edit:after-command = [ $@edit:after-command {|m| echo 'command took '$m[duration]' seconds' } ]
```

### After-Command Hook Map

The `$edit:after-command` callback receives a map with:

| Key | Type | Description |
|-----|------|-------------|
| `src` | string | Source code of the command |
| `duration` | number | Execution duration in seconds |
| `error` | exception or `$nil` | Exception if command failed |

## Navigation Mode

A file browser integrated into the editor.

### Starting Navigation

```elvish
edit:navigation:start  # or Ctrl-N
```

### Navigation Layout

Three columns with configurable width ratio:

```elvish
set edit:navigation:width-ratio = [1 3 4]
```

### Navigation Variables

| Variable | Description |
|----------|-------------|
| `$edit:navigation:selected-file` | Currently selected file |
| `$edit:navigation:width-ratio` | Width ratio of the 3 columns |

### Navigation Functions

| Function | Description |
|----------|-------------|
| `edit:navigation:insert-selected` | Insert selected filename |
| `edit:navigation:insert-selected-and-quit` | Insert and close navigation |
| `edit:navigation:trigger-filter` | Toggle filtering |
| `edit:navigation:trigger-shown-hidden` | Toggle showing hidden files |

## History

### History Walking Mode

Navigate through command history with Up/Down keys:

```elvish
edit:history:start   # Enter history walking mode
edit:history:up      # Walk to previous entry
edit:history:down    # Walk to next entry
edit:history:accept  # Accept current entry
```

### History Listing Mode

Browse and search command history:

```elvish
edit:histlist:start           # Open history listing
edit:histlist:toggle-dedup    # Toggle deduplication
```

### History Functions

| Function | Description |
|----------|-------------|
| `edit:command-history &cmd-only=$false &dedup=$false &newest-first` | Output command history entries |
| `edit:end-of-history` | Add "End of history" notification |
| `edit:history:fast-forward` | Import new history entries from other sessions |
| `edit:history:down-or-quit` | Walk to next or quit if at newest |
| `edit:insert-last-word` | Insert last word of last command |

### History Variables

| Variable | Description |
|----------|-------------|
| `$edit:add-cmd-filters` | Filters applied before adding command to history |

### Last Command Mode

Quick access to the last command's arguments:

```elvish
edit:lastcmd:start
```

Shows the last command and its individual words. Select a word to insert it.

## Location Mode

Directory browser that shows frequently visited directories:

```elvish
edit:location:start  # or Ctrl-L
```

### Location Variables

| Variable | Description |
|----------|-------------|
| `$edit:location:hidden` | Directories to hide from location list |
| `$edit:location:pinned` | Directories always shown at top |
| `$edit:location:workspaces` | Map of workspace types to patterns |

### Workspace Configuration

```elvish
set edit:location:workspaces = [&git=$E:HOME/repos &go=$E:HOME/go/src]
```

## Instant Mode

Evaluates code immediately as you type, showing the output:

```elvish
edit:-instant:start
```

Useful for quick calculations or testing expressions. The code is evaluated on every keystroke.

## Custom Listing Mode

Create custom listing UIs:

```elvish
edit:listing:start-custom $items &binding=$nil &caption='' &keep-bottom=$false &accept=$nil &auto-accept=$false
```

| Option | Type | Description |
|--------|------|-------------|
| `&binding` | binding table | Key bindings for the listing |
| `&caption` | string | Caption text for the filter area |
| `&keep-bottom` | boolean | Keep selection at bottom when filtering |
| `&accept` | function | Called when an item is accepted |
| `&auto-accept` | boolean | Auto-accept when only one item matches |

## Abbreviations

### Simple Abbreviations (`$edit:abbr`)

Expand when space is pressed after the abbreviation:

```elvish
set edit:abbr[gs] = 'git status'
set edit:abbr[dc] = 'docker-compose'
```

### Command Abbreviations (`$edit:command-abbr`)

Expand command names (first word on the line):

```elvish
set edit:command-abbr[g] = 'git'
```

### Small-Word Abbreviations (`$edit:small-word-abbr`)

Abbreviations that expand at small-word boundaries:

```elvish
set edit:small-word-abbr[gc] = 'git commit'
```

## Autofix

The editor identifies commands to fix errors automatically (e.g., suggesting `use str` for missing imports).

Applied by:
- `edit:completion:smart-start` (Tab binding)
- `edit:smart-enter` (Enter binding)

```elvish
edit:apply-autofix  # Manually apply autofix
```

## Word Types

Three word definitions used by word-movement and word-killing commands:

| Type | Definition | Example in `abc++ /* xyz` |
|------|-----------|--------------------------|
| Big word | Sequence of non-whitespace characters (vi "WORD") | `abc++`, `/*`, `xyz` |
| Small word | Sequence of alphanumerics OR sequence of non-alphanumeric, non-whitespace (vi/zsh "word") | `abc`, `++`, `/*`, `xyz` |
| Alphanumerical word | Sequence of alphanumerical characters (bash "word") | `abc`, `xyz` |

### Word-Related Functions

| Function | Description |
|----------|-------------|
| `edit:kill-word-left/right` | Kill big word |
| `edit:kill-small-word-left/right` | Kill small word |
| `edit:kill-alnum-word-left/right` | Kill alnum word |
| `edit:move-dot-left-word/right-word` | Move by big word |
| `edit:move-dot-left-small-word/right-small-word` | Move by small word |
| `edit:move-dot-left-alnum-word/right-alnum-word` | Move by alnum word |
| `edit:transpose-word/small-word/alnum-word` | Swap words |

## Editor Variables (Complete Reference)

| Variable | Type | Description |
|----------|------|-------------|
| `$edit:-dot` | int | Current cursor position (byte position) |
| `$edit:abbr` | map | Simple abbreviation map |
| `$edit:add-cmd-filters` | list | Filters before adding command to history |
| `$edit:after-command` | list | Functions called after command completes |
| `$edit:after-readline` | list | Functions called after each readline cycle |
| `$edit:before-readline` | list | Functions called before each readline cycle |
| `$edit:command-abbr` | map | Command abbreviation map |
| `$edit:command-duration` | number | Duration of most recent interactive command |
| `$edit:current-command` | string | Content of current input |
| `$edit:exceptions` | list | List of exceptions from callbacks |
| `$edit:insert:quote-paste` | boolean | Whether to quote bracketed paste content |
| `$edit:max-height` | int | Maximum height editor can use |
| `$edit:prompt` | function | Left-hand prompt function |
| `$edit:prompt-stale-threshold` | number | Seconds before prompt considered stale |
| `$edit:prompt-stale-transform` | function | Stale prompt transformer |
| `$edit:rprompt` | function | Right-side prompt function |
| `$edit:rprompt-persistent` | boolean | Keep rprompt visible after accepting |
| `$edit:rprompt-stale-threshold` | number | RPrompt stale threshold |
| `$edit:rprompt-stale-transform` | function | RPrompt stale transformer |
| `$edit:small-word-abbr` | map | Small-word abbreviation map |

## Editor Functions (Complete Reference)

### Mode Control

| Function | Description |
|----------|-------------|
| `edit:close-mode` | Close current active mode |
| `edit:clear` | Clear screen |
| `edit:redraw &full=$false` | Trigger redraw |
| `edit:return-line` | End read iteration, evaluate code |
| `edit:return-eof` | Terminate REPL |

### Text Editing

| Function | Description |
|----------|-------------|
| `edit:insert-at-dot $text` | Insert text at cursor |
| `edit:replace-input $text` | Replace current input |
| `edit:insert-last-word` | Insert last word of last command |
| `edit:insert-raw` | Insert next input uninterpreted |
| `edit:kill-rune-left/right` | Delete one character |
| `edit:kill-line-left/right` | Delete to start/end of line |
| `edit:move-dot-sol/eol` | Move to start/end of line |
| `edit:move-dot-left/right` | Move one character |
| `edit:move-dot-up/down` | Move one line |
| `edit:transpose-rune` | Swap adjacent characters |

### Notifications

| Function | Description |
|----------|-------------|
| `edit:notify $message` | Display notification above editor |
| `edit:apply-autofix` | Execute suggested autofix |

### Variable Management

| Function | Description |
|----------|-------------|
| `edit:add-var $name $init` | Add variable to REPL |
| `edit:add-vars $map` | Add multiple variables |
| `edit:del-var $name` | Delete variable from REPL |
| `edit:del-vars $list` | Delete multiple variables |

### Utility

| Function | Description |
|----------|-------------|
| `edit:key $string` | Parse string into key |
| `edit:binding-table $map` | Convert map to binding table |
| `edit:wordify $code` | Break elvish code into words |
| `edit:toggle-quote-paste` | Toggle bracketed paste quoting |

## References

- [Elvish Editor API Reference](https://elv.sh/ref/edit.html) — official documentation
- [Elvish Source: pkg/edit/](https://github.com/elves/elvish/tree/master/pkg/edit) — editor source code
- [Elvish Source: pkg/cli/](https://github.com/elves/elvish/tree/master/pkg/cli) — CLI framework

## Related Skills

- For completion system details (arg-completer, complex-candidate, matchers), see [references/completion.md](completion.md).
- For styling system details (styled, styled-segment, colors), see [references/styling.md](styling.md).
- For carapace-specific elvish integration, see the **carapace-dev** skill → `references/shell-elvish.md`.
