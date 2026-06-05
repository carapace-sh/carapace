# PowerShell Terminal Styling and ANSI Rendering

In-depth reference for PowerShell's ANSI/VT100 escape sequence support, the `$PSStyle` variable, console color APIs, and how styled text is rendered in the terminal.

## The `` `e `` Escape Character

PowerShell 6+ supports `` `e `` (backtick-e) as the **Escape character** (ASCII 27 / `\x1b`). This is the foundation for all ANSI/VT100 escape sequences in PowerShell:

```powershell
"`e[31mRed text`e[0m"   # SGR code for red foreground, then reset
```

`` `e[ `` is the CSI (Control Sequence Introducer) — equivalent to `\x1b[` or `ESC[` in other contexts. This is critical for embedding SGR color codes in PowerShell strings.

### Alternative: `[char]0x1b`

For PowerShell 5.1 (which lacks `` `e ``), use:

```powershell
$esc = [char]0x1b
"$esc[31mRed text$esc[0m"
```

### Alternative: `[char]27`

```powershell
$esc = [char]27
"$esc[31mRed text$esc[0m"
```

## The $PSStyle Automatic Variable

PowerShell 7.2+ provides `$PSStyle` as the primary interface for styled output. It is an instance of `System.Management.Automation.PSStyle`.

### Text Decorations

| Property | Effect |
|----------|--------|
| `$PSStyle.Reset` | Turn off all decorations |
| `$PSStyle.Blink` / `.BlinkOff` | Blink on/off |
| `$PSStyle.Bold` / `.BoldOff` | Bold on/off |
| `$PSStyle.Dim` / `.DimOff` | Dim on/off (PS 7.4+) |
| `$PSStyle.Hidden` / `.HiddenOff` | Hidden on/off |
| `$PSStyle.Reverse` / `.ReverseOff` | Reverse video on/off |
| `$PSStyle.Italic` / `.ItalicOff` | Italic on/off |
| `$PSStyle.Underline` / `.UnderlineOff` | Underline on/off |
| `$PSStyle.Strikethrough` / `.StrikethroughOff` | Strikethrough on/off |

### Foreground Colors

Standard 16 console colors via `$PSStyle.Foreground`:

| Property | Color |
|----------|-------|
| `Black` | Black |
| `BrightBlack` | Bright black (dark gray) |
| `White` | White |
| `BrightWhite` | Bright white |
| `Red` / `BrightRed` | Red variants |
| `Green` / `BrightGreen` | Green variants |
| `Blue` / `BrightBlue` | Blue variants |
| `Yellow` / `BrightYellow` | Yellow variants |
| `Cyan` / `BrightCyan` | Cyan variants |
| `Magenta` / `BrightMagenta` | Magenta variants |

### 24-bit True Color

```powershell
$PSStyle.Foreground.FromRgb(245, 245, 220)    # RGB values
$PSStyle.Foreground.FromRgb(0xf5f5dc)          # Hex value
$PSStyle.Background.FromRgb(0x2f6aff)          # Background hex
```

### Usage Example

```powershell
"$($PSStyle.Foreground.Red)$($PSStyle.Bold)Error:$($PSStyle.Reset) Something went wrong"
"$($PSStyle.Background.BrightCyan)Power$($PSStyle.Underline)$($PSStyle.Bold)Shell$($PSStyle.Reset)"
```

### Static Methods (PS 7.4+)

```powershell
[System.Management.Automation.PSStyle]::MapForegroundColorToEscapeSequence('Red')
[System.Management.Automation.PSStyle]::MapBackgroundColorToEscapeSequence('Black')
[System.Management.Automation.PSStyle]::MapColorPairToEscapeSequence('Red','Black')
```

## Output Rendering Control

`$PSStyle.OutputRendering` controls how PowerShell handles ANSI sequences in output:

| Value | Behavior |
|-------|----------|
| `ANSI` | ANSI sequences always passed through as-is |
| `PlainText` | ANSI sequences always stripped to plain text |
| `Host` | **Default** — ANSI removed from redirected/piped output, kept for console |

```powershell
$PSStyle.OutputRendering = 'ANSI'       # Keep ANSI in all output
$PSStyle.OutputRendering = 'PlainText'  # Strip all ANSI
$PSStyle.OutputRendering = 'Host'       # Default behavior
```

**Important**: Use `ANSI` mode when redirecting output to files or pipelines intended for downstream execution that expects ANSI codes.

## Formatting Control

`$PSStyle.Formatting` controls default formatting for output streams:

| Property | Controls |
|----------|----------|
| `FormatAccent` | List item formatting |
| `ErrorAccent` | Error metadata formatting |
| `Error` | Error message formatting |
| `Warning` | Warning message formatting |
| `Verbose` | Verbose message formatting |
| `Debug` | Debug message formatting |
| `TableHeader` | Table header formatting |
| `CustomTableHeaderLabel` | Custom table header formatting |
| `FeedbackName` / `FeedbackText` / `FeedbackAction` | Feedback provider formatting (PS 7.4+) |

## FileInfo Coloring

`$PSStyle.FileInfo` controls colors for `Get-ChildItem` output:

| Property | Description |
|----------|-------------|
| `Directory` | Directory color |
| `SymbolicLink` | Symbolic link color |
| `Executable` | Executable file color |
| `Extension` | Dictionary for per-extension colors |

```powershell
$PSStyle.FileInfo.Directory = $PSStyle.Foreground.BrightCyan
$PSStyle.FileInfo.Extension['.ps1'] = $PSStyle.Foreground.Cyan
$PSStyle.FileInfo.Extension['.py'] = $PSStyle.Foreground.Yellow
```

Predefined extension colors: `.ps1`, `.ps1xml`, `.psd1`, `.psm1`.

## Progress Bar Control

`$PSStyle.Progress` controls progress bar rendering:

| Property | Default | Description |
|----------|---------|-------------|
| `Style` | — | ANSI string for rendering style |
| `MaxWidth` | 120 | Maximum width (min: 18) |
| `View` | `Minimal` | `Minimal` or `Classic` |
| `UseOSCIndicator` | `$false` | Use OSC indicator for OSC-capable terminals |

```powershell
$PSStyle.Progress.View = 'Minimal'
$PSStyle.Progress.MaxWidth = 80
```

If the host doesn't support Virtual Terminal, `View` is automatically set to `Classic`.

## Write-Host Colors

`Write-Host` supports `-ForegroundColor` and `-BackgroundColor` parameters using `[ConsoleColor]` enum values:

```powershell
Write-Host "Error" -ForegroundColor Red -BackgroundColor Black
Write-Host "Success" -ForegroundColor Green
```

### ConsoleColor Values

`Black`, `DarkBlue`, `DarkGreen`, `DarkCyan`, `DarkRed`, `DarkMagenta`, `DarkYellow`, `Gray`, `DarkGray`, `Blue`, `Green`, `Cyan`, `Red`, `Magenta`, `Yellow`, `White`

### Limitations

`Write-Host` only supports the 16 `ConsoleColor` values. For 24-bit true color or text decorations (bold, underline, etc.), use `$PSStyle` or raw ANSI escape sequences.

## SGR (Select Graphic Rendition) Codes

SGR codes are the standard ANSI escape sequences for text styling. Format: `ESC[<code>m`

### Common SGR Codes

| Code | Effect |
|------|--------|
| `0` | Reset all attributes |
| `1` | Bold |
| `2` | Dim |
| `3` | Italic |
| `4` | Underline |
| `5` | Blink (slow) |
| `7` | Reverse video |
| `8` | Hidden |
| `9` | Strikethrough |
| `22` | Bold off |
| `23` | Italic off |
| `24` | Underline off |
| `27` | Reverse off |
| `28` | Hidden off |
| `29` | Strikethrough off |
| `30-37` | Foreground (standard 8 colors) |
| `38;5;<n>` | 256-color foreground |
| `38;2;<r>;<g>;<b>` | 24-bit true color foreground |
| `39` | Default foreground |
| `40-47` | Background (standard 8 colors) |
| `48;5;<n>` | 256-color background |
| `48;2;<r>;<g>;<b>` | 24-bit true color background |
| `49` | Default background |
| `90-97` | Bright foreground |
| `100-107` | Bright background |

### PowerShell SGR Examples

```powershell
# 256-color
"`e[38;5;208mOrange text`e[0m"

# 24-bit true color
"`e[38;2;255;165;0mOrange text`e[0m"

# Bold + red
"`e[1;31mBold red`e[0m"

# Multiple attributes
"`e[3;4;36mItalic underlined cyan`e[0m"
```

## StringDecorated Type

PowerShell internally uses `System.Management.Automation.StringDecorated` for ANSI-aware string handling:

| Property/Method | Description |
|-----------------|-------------|
| `IsDecorated` | `true` if string contains ESC or C1 CSI sequences |
| `Length` | Length **without** ANSI escape sequences |
| `Substring(int contentLength)` | Substring preserving ANSI sequences |
| `ToString()` | Plaintext version |
| `ToString(bool Ansi)` | Raw ANSI string if `true`, plaintext if `false` |
| `FormatHyperlink(string text, uri link)` | Creates clickable terminal hyperlink (OSC 8) |

### Hyperlink Support

```powershell
# OSC 8 hyperlinks (supported in modern terminals)
$PSStyle.FormatHyperlink("Click here", "https://example.com")
```

## Terminal Compatibility

| Platform | Terminal | VT Support |
|----------|----------|------------|
| Windows 10+ | Windows Console Host | Yes |
| Windows | Windows Terminal | Yes (full) |
| macOS | Terminal.app | Yes (xterm compatible) |
| Linux | Varies | Check terminal docs |

### Disabling ANSI Output

| Variable | Effect |
|----------|--------|
| `$Env:TERM=dumb` | Sets `$Host.UI.SupportsVirtualTerminal = $false` |
| `$Env:TERM=xterm-mono` | Sets `$PSStyle.OutputRendering = PlainText` |
| `$Env:NO_COLOR` exists | Sets `$PSStyle.OutputRendering = PlainText` |

### Checking VT Support

```powershell
$Host.UI.SupportsVirtualTerminal  # Returns $true if VT is supported
```

## PSReadLine Color Integration

PSReadLine applies its own syntax coloring independently from `$PSStyle`. Colors are configured via `Set-PSReadLineOption -Colors`:

```powershell
Set-PSReadLineOption -Colors @{
    "Command"          = "`e[36m"       # Cyan commands
    "String"           = "`e[38;5;100m"  # 256-color strings
    "InlinePrediction" = "`e[38;5;238m"  # Dim predictions
    "Selection"        = "`e[7;36m"     # Reverse + cyan selection
}
```

PSReadLine injects ANSI escape sequences directly into its render buffer during `GenerateRender()`. The rendering pipeline uses `VTColorUtils.AnsiReset` to reset colors between token types.

## C# Integration

```csharp
using System.Management.Automation;

string output = $"{PSStyle.Instance.Foreground.Red}{PSStyle.Instance.Bold}Hello{PSStyle.Instance.Reset}";
```

## Cmdlets Using ANSI Output

| Cmdlet | ANSI Usage |
|--------|------------|
| `Show-Markdown` | Renders markdown with ANSI styling |
| `Get-Error` | Formatted error view with colors |
| `Select-String` | Highlights matching patterns (PS 7.0+) |
| `Write-Progress` | Progress bar rendering |
| PSReadLine | Syntax coloring, completion menu, predictions |

## References

- [about_ANSI_Terminals](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_ansi_terminals)
- [PSStyle Class](https://learn.microsoft.com/en-us/dotnet/api/system.management.automation.psstyle)
- [ConsoleColor Enum](https://learn.microsoft.com/en-us/dotnet/api/system.consolecolor)
- [Write-Host](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/write-host)
- [PSReadLine Colors](https://learn.microsoft.com/en-us/powershell/module/psreadline/set-psreadlineoption)

## Related Skills

- **powershell** → `references/editor.md` — PSReadLine rendering pipeline, token coloring
- **powershell** → `references/completion.md` — CompletionResult, how SGR codes appear in completions
- **carapace-dev** → `references/shell-powershell.md` — Carapace SGR styling in completion output
