# PSReadLine and the Completion UI

In-depth reference for PSReadLine — PowerShell's line-editing module — covering key bindings, edit modes, the completion menu, menu rendering, prediction system, and the rendering pipeline.

## How PSReadLine Works

PSReadLine replaces PowerShell's default line editing with a rich, customizable editing experience. It sits between the user and the terminal:

1. **Initialization** — loaded automatically in PowerShell 7+; reads options from profile
2. **Character input** — reads key sequences from the terminal, maps them to editing commands
3. **Buffer management** — text stored in an internal buffer; characters inserted at cursor
4. **Rendering** — differential rendering updates only changed lines on screen
5. **Completion** — integrates with PowerShell's `TabExpansion2` for tab completion
6. **Prediction** — queries prediction sources for inline or list suggestions

### Module Information

| Property | Value |
|----------|-------|
| Current version | 2.4.5 |
| Shipped with PS7 | Yes (bundled) |
| GitHub | https://github.com/PowerShell/PSReadLine |
| Namespace | `Microsoft.PowerShell.PSConsoleReadLine` |

## Edit Modes

PSReadLine supports three edit modes, each with different default key bindings:

| Mode | Default On | Key Binding Style |
|------|-----------|-------------------|
| **Windows** | Windows | Windows conventions (Ctrl+C copy, Ctrl+V paste) |
| **Emacs** | Non-Windows | Emacs-style (Ctrl+A beginning, Ctrl+E end) |
| **Vi** | — | Vi modal editing (insert/command modes) |

```powershell
Set-PSReadLineOption -EditMode Emacs
```

## Key Bindings

### Managing Key Bindings

```powershell
# List all bound keys
Get-PSReadLineKeyHandler

# List all keys (bound + unbound)
Get-PSReadLineKeyHandler -Bound -Unbound

# Check specific chord
Get-PSReadLineKeyHandler -Chord Enter

# Remove custom binding
Remove-PSReadLineKeyHandler -Key UpArrow, DownArrow
```

### Binding Keys to Functions

```powershell
# Built-in function
Set-PSReadLineKeyHandler -Chord Tab -Function MenuComplete

# Custom script block
Set-PSReadLineKeyHandler -Chord Ctrl+b -ScriptBlock {
    param($key, $arg)
    [Microsoft.PowerShell.PSConsoleReadLine]::RevertLine()
    [Microsoft.PowerShell.PSConsoleReadLine]::Insert('build')
    [Microsoft.PowerShell.PSConsoleReadLine]::AcceptLine()
}
```

### Custom Script Block Parameters

Custom key handlers receive:

| Parameter | Type | Description |
|-----------|------|-------------|
| `$key` | `ConsoleKeyInfo` | The key that was pressed (`.KeyChar`, `.Key`, `.Modifiers`) |
| `$arg` | `object` | Numeric argument prefix (if any) |

### PSConsoleReadLine API for Custom Handlers

```csharp
void Insert(string s)              // Insert text at cursor
void Delete(int start, int len)    // Delete characters
void Replace(int start, int len, string replacement)
void GetBufferState([ref]$line, [ref]$cursor)  // Get current buffer
void SetCursorPosition(int pos)    // Move cursor
void Ding()                        // Sound bell
void AddToHistory(string cmd)      // Add to history
void AcceptLine()                  // Execute current line
void RevertLine()                  // Discard all edits
```

### Key Binding Reference (Selected)

| Function | Windows | Emacs | Vi Insert | Vi Command |
|----------|---------|-------|----------|------------|
| `AcceptLine` | Enter | Enter | Enter | Enter |
| `TabCompleteNext` | Tab | — | — | — |
| `TabCompletePrevious` | Shift+Tab | — | — | — |
| `Complete` | — | Tab | — | — |
| `MenuComplete` | Ctrl+Space | Ctrl+Space | — | — |
| `PossibleCompletions` | — | Alt+= | Ctrl+Space | — |
| `ReverseSearchHistory` | Ctrl+r | Ctrl+r | Ctrl+r | Ctrl+r |
| `HistorySearchBackward` | F8 | — | — | — |
| `BackwardDeleteChar` | Backspace | Backspace, Ctrl+h | Backspace | X, d,h |
| `KillLine` | — | Ctrl+k | — | — |
| `Undo` | Ctrl+z | Ctrl+_, Ctrl+x Ctrl+u | Ctrl+z | Ctrl+z, u |
| `Paste` | Ctrl+v | Ctrl+v | Ctrl+v | Ctrl+v |
| `GotoBrace` | Ctrl+] | Ctrl+] | Ctrl+] | Ctrl+] |

## Completion Functions

### TabCompleteNext / TabCompletePrevious

Cycle through completions one at a time (inline). First press queries `TabExpansion2` and caches results. Subsequent presses cycle through the cache.

### Complete (Emacs mode)

Shows the common prefix if multiple completions share one, or completes the single match. If ambiguous, first Tab shows common prefix, second Tab lists all completions.

### MenuComplete

Interactive menu — shows all completions in a columnar layout. Navigate with arrow keys. Type characters to filter. See **The Completion Menu** below.

### PossibleCompletions

Lists all completions below the command line (non-interactive display).

## The Completion Menu

### How the Menu is Drawn

When `MenuComplete` is invoked, PSReadLine enters an interactive event loop:

1. Queries `TabExpansion2` for completion candidates
2. Calculates column layout (see below)
3. Hides cursor, saves position
4. Draws menu below the command line
5. Enters key-reading loop for navigation

### Column Layout

Items fill **top-to-bottom, then left-to-right**:

```
Item 0    Item 3    Item 6    Item 9
Item 1    Item 4    Item 7
Item 2    Item 5    Item 8
```

Layout calculation:

```csharp
colWidth = Math.Min(matches.Max(c => LengthInBufferCells(c.ListItemText)) + 2, bufferWidth);
columns  = Math.Max(1, bufferWidth / colWidth);
rows     = (matches.Count + columns - 1) / columns;
```

Items are padded to `columnWidth` with spaces, or truncated with `...` if too long.

### Menu Navigation Keys

| Key | Action |
|-----|--------|
| **Right Arrow** | Next column (wraps to next row's leftmost) |
| **Left Arrow** | Previous column (wraps to previous row's rightmost) |
| **Down Arrow** | Next item (wraps to 0 from last) |
| **Up Arrow** | Previous item (wraps to last from 0) |
| **Page Down** | Jump to last item in current column |
| **Page Up** | Jump to first item in current column |
| **Tab** | Show common prefix or move to next item |
| **Shift+Tab** | Move to previous item |
| **Backspace** | Remove last filter character, restore previous menu |
| **Escape / Ctrl+G** | Cancel menu, undo all changes |
| **Space / Enter** | Accept current selection |
| **Ending keys** | Accept + insert key (see CompletionResultType) |

### Filtering

Typing characters during menu completion **filters** the menu incrementally. PSReadLine maintains a **stack of menus** — each filter creates a new filtered menu. Backspace pops the stack to restore the previous (less filtered) state.

### Tooltips

Tooltips are displayed **below the menu**, separated by one blank line:

- Tooltip text comes from `CompletionResult.ToolTip`
- Only shown if it adds information (not identical to `ListItemText` or `CompletionText`)
- Tooltip lines are calculated by scanning for newlines and buffer-width wrapping
- Tooltips are hidden if they would scroll the command line off screen
- Tooltip color uses `Options._emphasisColor`

### Selection Colors

| Color | Used For | Source |
|-------|----------|--------|
| `Options._selectionColor` | Highlighted/selected menu item | `Set-PSReadLineOption -Colors @{ Selection = "..." }` |
| `Options._emphasisColor` | Tooltip text | `Set-PSReadLineOption -Colors @{ Emphasis = "..." }` |

### Undo Integration

Before entering the menu, PSReadLine records an undo point. All intermediate edits from cycling completions are consolidated into a single `GroupedEdit` on exit — enabling one-step undo of the entire completion session.

## Predictive IntelliSense

PSReadLine 2.1+ supports predictive suggestions. 2.2+ adds plugin support. 2.2.6+ enables predictions by default.

### Prediction Sources

```powershell
Set-PSReadLineOption -PredictionSource None           # Disabled
Set-PSReadLineOption -PredictionSource History         # History only
Set-PSReadLineOption -PredictionSource Plugin          # External plugins only
Set-PSReadLineOption -PredictionSource HistoryAndPlugin # Both (default for PS 7.2+)
```

### View Styles

**InlineView** (default): Shows prediction as dimmed ghost text after the cursor. Press `RightArrow` to accept the full suggestion, or `Ctrl+f` / `AcceptNextSuggestionWord` to accept one word at a time.

**ListView**: Shows a scrollable dropdown list below the command line. Predictions are grouped by source (History first, then plugins in load order). Navigate with arrow keys.

```powershell
Set-PSReadLineOption -PredictionViewStyle InlineView  # Ghost text
Set-PSReadLineOption -PredictionViewStyle ListView     # Dropdown
```

### Prediction Key Bindings

| Function | Description | Default |
|----------|-------------|---------|
| `AcceptSuggestion` | Accept inline suggestion | Built into `RightArrow` |
| `AcceptNextSuggestionWord` | Accept next word | Unbound (bind to `Ctrl+f`) |
| `NextSuggestion` | Next in ListView | Unbound |
| `PreviousSuggestion` | Previous in ListView | Unbound |
| `SwitchPredictionView` | Toggle Inline/List | F2 |
| `ShowFullPredictionTooltip` | Show tooltip | F4 |

### Creating a Prediction Plugin

Prediction plugins are C# modules implementing `System.Management.Automation.Subsystem.Prediction.ICommandPredictor`:

1. Create a C# class library project
2. Implement `ICommandPredictor` interface
3. Register the predictor in the module manifest
4. Install the module

The interface declares methods for:
- Querying prediction results based on input text
- Providing feedback (whether a suggestion was accepted)

### CompletionPredictor Module

The `CompletionPredictor` module (shipped with PSReadLine) adds IntelliSense for anything tab-completable:
- With `InlineView`: behaves like normal tab completion
- With `ListView`: full IntelliSense experience with descriptions

### Az.Tools.Predictor

Azure PowerShell prediction plugin using machine learning to predict Azure cmdlet commands and parameters.

## PSReadLine Options

### Set-PSReadLineOption

```powershell
Set-PSReadLineOption [-EditMode <EditMode>] [-ContinuationPrompt <string>]
    [-HistoryNoDuplicates] [-MaximumHistoryCount <int>] [-MaximumKillRingCount <int>]
    [-BellStyle <BellStyle>] [-DingTone <int>] [-DingDuration <int>]
    [-CompletionQueryItems <int>] [-WordDelimiters <string>]
    [-HistorySearchCaseSensitive] [-HistorySaveStyle <HistorySaveStyle>]
    [-HistorySavePath <string>] [-Colors <hashtable>] [-ViModeIndicator <ViModeStyle>]
    [-PredictionSource <PredictionSource>] [-PredictionViewStyle <PredictionViewStyle>]
    [-ExtraPromptLineCount <int>] [-ShowToolTips] [-TerminateOrphanedConsoleApps]
    [-AddToHistoryHandler <Func[string,Object]>] [-CommandValidationHandler <Action[CommandAst]>]
    [-PromptText <string[]>]
```

### Key Options

| Option | Default | Description |
|--------|---------|-------------|
| `EditMode` | Windows | Editing mode (Windows/Emacs/Vi) |
| `HistoryNoDuplicates` | True | Skip duplicates in history recall |
| `MaximumHistoryCount` | 4096 | Max commands in history |
| `MaximumKillRingCount` | 10 | Kill ring size |
| `CompletionQueryItems` | 100 | Max completions before prompting |
| `WordDelimiters` | `;:,.\[\]{}()/\\\|^&*-=+'"–—―` | Word boundary characters |
| `BellStyle` | Audible | Audible/Visual/None |
| `ShowToolTips` | True | Show tooltips in completion menu |
| `PredictionSource` | None/HistoryAndPlugin | Prediction source |
| `PredictionViewStyle` | InlineView | InlineView/ListView |
| `HistorySaveStyle` | SaveIncrementally | How history is saved |
| `HistorySavePath` | Platform-specific | Path to history file |

### Color Configuration

```powershell
Set-PSReadLineOption -Colors @{
    "Comment"              = "DarkGray"
    "Command"             = "#8181f7"
    "String"              = "$([char]0x1b)[38;5;100m"
    "Error"               = [ConsoleColor]::DarkRed
    "InlinePrediction"    = "`e[97;2;3m"
    "ListPrediction"      = "DarkYellow"
    "ListPredictionSelected" = "`e[48;5;238m"
    "Selection"           = "`e[36m"
    "Emphasis"            = "`e[33m"
}
```

Available color keys: `ContinuationPrompt`, `Emphasis`, `Error`, `Selection`, `Default`, `Comment`, `Keyword`, `String`, `Variable`, `Command`, `Parameter`, `Type`, `Number`, `Member`, `InlinePrediction`, `ListPrediction`, `ListPredictionSelected`.

### Handler-Based Options

**AddToHistoryHandler** — filter commands from history:

```powershell
Set-PSReadLineOption -AddToHistoryHandler {
    param([string]$Line)
    if ($Line -match "^git") { return $false } else { return $true }
}
```

**CommandValidationHandler** — validate before execution:

```powershell
using namespace System.Management.Automation.Language
Set-PSReadLineOption -CommandValidationHandler {
    param([CommandAst]$CommandAst)
    # Custom validation logic
}
Set-PSReadLineKeyHandler -Chord Enter -Function ValidateAndAcceptLine
```

**ViModeChangeHandler** — custom Vi mode indicator:

```powershell
function OnViModeChange {
    if ($args[0] -eq 'Command') {
        Write-Host -NoNewline "`e[1 q"  # Blinking block cursor
    } else {
        Write-Host -NoNewline "`e[5 q"  # Blinking line cursor
    }
}
Set-PSReadLineOption -ViModeIndicator Script -ViModeChangeHandler $Function:OnViModeChange
```

## The Rendering Pipeline

PSReadLine uses differential rendering to minimize terminal writes.

### Main Flow

```
Render(force)
  → Screen reader mode? → RenderForScreenReader() (diff-based)
  → Normal mode → ForceRender()
    → GenerateRender() — token coloring + ANSI escape sequences
    → ReallyRender() — write to console, clear leftovers
```

### Token-Based Coloring

`GenerateRender()` uses PowerShell's parser tokens to apply syntax highlighting:

| Token Type | Color Source |
|------------|-------------|
| Comments | `_commentColor` |
| Commands | `_commandColor` |
| Parameters (`-` prefix) | `_parameterColor` |
| Variables | `_variableColor` |
| Strings | `_stringColor` |
| Numbers | `_numberColor` |
| Keywords | `_keywordColor` |
| Operators | `_operatorColor` |
| Types | `_typeColor` |
| Members | `_memberColor` |

Colors are applied by injecting ANSI escape sequences into the render buffer. `UpdateColorsIfNecessary()` only writes color changes when the color actually changes.

### Differential Rendering

`CalculateWhereAndWhatToRender()` finds the first changed logical line between the previous and current render. Only lines from that point onward are rewritten. Leftover content from the previous render is cleared with spaces.

### Performance Optimizations

- Skips rendering when >10 keys queued and <50ms since last render
- Caches physical line counts in `RenderedLineData`
- Reuses `_consoleBufferLines` list across renders
- Skips resize check within 50ms window (avoids PTY round-trips)

## Profile Persistence

Add PSReadLine configuration to your PowerShell profile for persistence:

```powershell
# Linux/macOS: ~/.config/powershell/Microsoft.PowerShell_profile.ps1
# Windows: ~/Documents/PowerShell/Microsoft.PowerShell_profile.ps1

Set-PSReadLineOption -EditMode Emacs
Set-PSReadLineOption -PredictionSource HistoryAndPlugin
Set-PSReadLineOption -ShowToolTips $true
Set-PSReadLineOption -Colors @{
    "Command" = "`e[36m"
    "InlinePrediction" = "`e[38;5;238m"
}

Set-PSReadLineKeyHandler -Chord UpArrow -Function HistorySearchBackward
Set-PSReadLineKeyHandler -Chord DownArrow -Function HistorySearchForward
Set-PSReadLineKeyHandler -Chord Tab -Function MenuComplete
```

## References

- [PSReadLine Module](https://learn.microsoft.com/en-us/powershell/module/psreadline/)
- [Set-PSReadLineKeyHandler](https://learn.microsoft.com/en-us/powershell/module/psreadline/set-psreadlinekeyhandler)
- [Set-PSReadLineOption](https://learn.microsoft.com/en-us/powershell/module/psreadline/set-psreadlineoption)
- [PSReadLine GitHub](https://github.com/PowerShell/PSReadLine)
- [about_PSReadLine_Functions](https://learn.microsoft.com/en-us/powershell/module/psreadline/about/about_psreadline_functions)
- [Using Predictors](https://learn.microsoft.com/en-us/powershell/scripting/learn/shell/using-predictors)
- [PSReadLine Rendering (DeepWiki)](https://deepwiki.com/PowerShell/PSReadLine/2.2-rendering-system)
- [PSReadLine Completion System (DeepWiki)](https://deepwiki.com/PowerShell/PSReadLine/3.4-completion-system)

## Related Skills

- **powershell** → `references/completion.md` — Register-ArgumentCompleter, CompletionResult, TabExpansion2
- **powershell** → `references/styling.md` — ANSI escape sequences, $PSStyle, terminal rendering
- **carapace-dev** → `references/shell-powershell.md` — Carapace-specific PowerShell integration
