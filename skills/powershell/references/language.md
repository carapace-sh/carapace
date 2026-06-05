# PowerShell Language and Argument Passing

In-depth reference for PowerShell's quoting rules, argument passing to native commands, the parsing model, and how these affect completion and external tool integration.

## Two Parsing Modes

PowerShell parses input in two distinct modes:

| Mode | Purpose | Behavior |
|------|---------|----------|
| **Expression mode** | Value manipulation | Numbers treated as values, operators evaluated |
| **Argument mode** | Command arguments | Input treated as expandable strings |

PowerShell first tries to interpret input as an expression before treating it as an argument. This is fundamentally different from bash, which always treats command-line tokens as arguments.

## Quoting Rules

### Single-Quoted Strings (Verbatim)

- No substitution is performed — the string is passed exactly as typed
- Variable names like `$i` are NOT replaced
- Expressions like `$(2+3)` are NOT evaluated
- The backtick (`` ` ``) is treated as a literal character, not an escape

```powershell
$i = 5
'The value of $i is $i.'
# Output: The value $i is $i.
```

### Double-Quoted Strings (Expandable)

- Variable substitution is performed — `$variable` is replaced with its value
- Expressions are evaluated — `$(2+3)` is replaced with the result
- The backtick (`` ` ``) is the escape character

```powershell
$i = 5
"The value of $i is $i."
# Output: The value 5 is 5.
```

### Variable Expansion in Double Quotes

- Variable names preceded by `$` are replaced with their values
- Use braces `${}` to separate variable names from subsequent characters
- A colon after `$variable` is treated as a scope specifier — use `${variable}` to avoid this

```powershell
"${HOME}: where the heart is."  # Works correctly
"$HOME: where the heart is."   # Error — PowerShell treats $HOME: as scoped variable
```

### Subexpression Expansion

Use `$( )` for expressions within double-quoted strings:

```powershell
"The value of $(2+3) is 5."
"PS version: $($PSVersionTable.PSVersion)"
```

Array indexing or member access **must** be enclosed in a subexpression:

```powershell
"First element: $($array[0])"    # Correct
"First element: $array[0]"       # Wrong — includes literal [0]
```

### Backtick Escaping in Double-Quoted Strings

The backtick (`` ` ``) is the PowerShell escape character in double-quoted strings:

| Sequence | Result |
|----------|--------|
| `` `$ `` | Literal `$` (prevent variable expansion) |
| `` `" `` | Literal `"` |
| `` `` `` | Literal `` ` `` |
| `` `0 `` | Null character |
| `` `a `` | Alert (bell) |
| `` `b `` | Backspace |
| `` `e `` | Escape (PS 6+) |
| `` `f `` | Form feed |
| `` `n `` | Newline |
| `` `r `` | Carriage return |
| `` `t `` | Horizontal tab |
| `` `u{xxxx} `` | Unicode character (PS 6+) |
| `` `v `` | Vertical tab |

**Important**: In single-quoted strings, the backtick is treated as a literal character — it does NOT escape anything.

### Here-Strings

Multi-line strings delimited by `@`:

```powershell
# Expandable here-string (variables replaced)
@"
Line 1: $HOME
Line 2: $(Get-Date)
"@

# Literal here-string (no substitution)
@'
Line 1: $HOME
Line 2: $(Get-Date)
'@
```

Rules:
- Opening `@"` or `@'` must be followed by a newline
- Closing `"@` or `'@` must be on a line by itself (at the start of the line)
- Quotation marks within here-strings are interpreted literally

### Including Quote Characters

| Goal | Method |
|------|--------|
| Include `"` in a string | Use single quotes: `'As they say, "live and learn."'` |
| Include `'` in a string | Use double quotes: `"As they say, 'live and learn.'"` |
| Include `'` in single-quoted string | Double the quote: `'don''t'` → `don't` |
| Include `"` in double-quoted string | Double the quote: `"As they say, ""live and learn."""` |
| Include `"` literally in double-quoted string | Use backtick: `` "Use a `" to begin a string." `` |

## Argument Passing to Native Commands

This is one of the most important topics for external tool integration. PowerShell handles argument passing very differently from bash.

### The Problem

PowerShell interprets arguments **before** passing them to native commands. This interpretation:
- Removes outer quote characters
- Expands variables
- Interprets special metacharacters

### PowerShell Metacharacters

These characters have special meaning in argument mode:

```
<space>  '  "  `  ,  ;  (  )  {  }  |  &  <  >  @  #
```

- `< > @ #` are only special at the **start** of a token
- Parentheses `()` begin new expressions — not passed literally
- Commas create arrays for PowerShell cmdlets but are part of expandable strings for native apps

### PSNativeCommandArgumentPassing (PS 7.3+)

PowerShell 7.3 introduced a major change in how arguments are passed to native commands, controlled by `$PSNativeCommandArgumentPassing`:

| Mode | Platform Default | Behavior |
|------|-----------------|----------|
| `Legacy` | — | Pre-7.3 behavior: quotes stripped, empty strings lost |
| `Standard` | Non-Windows | Quotes preserved, empty strings preserved |
| `Windows` | Windows | Same as Standard, but auto-uses Legacy for specific files |

**Windows mode** automatically uses `Legacy` style for:
- `cmd.exe`, `cscript.exe`, `wscript.exe`
- Files ending with `.bat`, `.cmd`, `.js`, `.vbs`, `.wsf`

### Key Behavior Differences

**Quotes preserved (Standard/Windows):**

```powershell
$a = 'a" "b'
TestExe -echoargs $a 'c" "d' e" "f
# Standard: Arg 0 is <a" "b>, Arg 1 is <c" "d>, Arg 2 is <e f>
# Legacy:   Arg 0 is <a b>, Arg 1 is <c d>, Arg 2 is <e f>
```

**Empty strings preserved (Standard/Windows):**

```powershell
TestExe -echoargs '' a b ''
# Standard: Arg 0 is <>, Arg 1 is <a>, Arg 2 is <b>, Arg 3 is <>
# Legacy:   Arg 1 is <a>, Arg 2 is <b>  (empty strings lost)
```

### The Stop-Parsing Token (`--%`)

PowerShell 3.0+ provides `--%` to stop PowerShell from interpreting subsequent input:

```powershell
icacls X:\VMS --% /grant Dom\HVAdmin:(CI)(OI)F
```

Characteristics:
- Only intended for native commands on Windows
- Only substitution performed is `%variable%` environment variable notation
- Effective until the next newline or pipeline character
- Cannot use line continuation (`` ` ``) to extend
- Cannot use command delimiter (`;`) to terminate
- Stream redirection (`>file.txt`) is passed verbatim, not interpreted
- `%<name>%` tokens are always expanded; undefined variables pass through as-is

### Common Pitfalls

| Pitfall | Example | Solution |
|---------|---------|----------|
| Variable expansion before passing | `cmd /c type $HOME\file.txt` | Use quotes: `'$HOME\file.txt'` or `` "`$HOME\file.txt" `` |
| Parentheses misinterpreted | `icacls X:\VMS /grant Dom\HVAdmin:(CI)(OI)F` | Use `--%` or escape: `` `(CI`) ``(PS 2.0) |
| Pipe in arguments | `cmd /c echo "a\|b"` | Use `--%`: `cmd /c --% echo "a\|b"` |
| Empty strings lost (pre-7.3) | `TestExe -echoargs '' a b ''` | Upgrade to PS 7.3+ or use `--%` |
| Tilde not expanded for native commands | `native-cmd ~/file` | Use `$HOME/file` or explicit path |
| Comma creates array for cmdlets | `Write-Output 1,2,3` | For native: quote or use `--%` |
| Batch file `%%` not escapable | `--% echo %%PATH%%` | No workaround — `%VAR%` always expands in `--%` mode |

### Passing Arguments to External Completion Tools

For tools like carapace that receive arguments from PowerShell completers, the key concern is ensuring arguments arrive intact:

```powershell
# The carapace snippet uses this pattern:
$elems = @()
foreach ($_ in $commandElements) {
    # ... process each element
    $elems += $t.replace('`',', ',)  # Fix backtick-comma escaping
}
# Pass as individual arguments:
example _carapace powershell $($elems| ForEach-Object {$_})
```

The `$($elems| ForEach-Object {$_})` pattern is PowerShell's equivalent of bash's `"${array[@]}"` — it expands the array into individual arguments, preventing the array from being passed as a single comma-separated string.

## The Pipeline

### Object Pipeline

Unlike bash's text pipeline, PowerShell passes **.NET objects** between commands:

```powershell
Get-ChildItem -Path *.txt |
  Where-Object {$_.Length -gt 10000} |
    Sort-Object -Property Length |
      Format-Table -Property Name, Length
```

Objects flow through with full type information — no text parsing needed.

### Parameter Binding

Pipeline input binds to parameters in two ways:

| Mode | Description |
|------|-------------|
| **ByValue** | Parameter accepts values matching (or convertible to) the expected .NET type |
| **ByPropertyName** | Parameter accepts input when the input object has a matching property name |

### One-at-a-Time Processing

When piping objects to a command, PowerShell sends them **one at a time**. When passing as a collection via a parameter, they are sent as a **single array object**. This difference has significant consequences for cmdlet behavior.

### Native Commands in the Pipeline

```powershell
ipconfig.exe | Select-String -Pattern 'IPv4'
```

PowerShell 7.4+ added `PSNativeCommandPreserveBytePipe` to preserve raw byte-stream data when redirecting or piping native command I/O.

## Script Blocks

Script blocks are **first-class objects** in PowerShell (`System.Management.Automation.ScriptBlock`):

```powershell
$block = { param($p1, $p2); "p1: $p1"; "p2: $p2" }
&$block "hello" "world"
```

### Delay-Bind Script Blocks

A unique PowerShell feature — script blocks used as parameter values that receive pipeline input:

```powerspowershell
dir config.log | Rename-Item -NewName { "old_$($_.Name)" }
```

The script block runs once per pipeline object, with `$_` bound to the current object.

## Operators

PowerShell operators are prefixed with `-` (dash), making them distinct from other shells:

### Comparison Operators

| Operator | Description |
|----------|-------------|
| `-eq` / `-ne` | Equal / Not equal (case-insensitive) |
| `-ceq` / `-cne` | Case-sensitive equal / not equal |
| `-ieq` / `-ine` | Explicitly case-insensitive |
| `-gt` / `-lt` / `-le` / `-ge` | Greater/less than |
| `-like` / `-notlike` | Wildcard pattern matching |
| `-match` / `-notmatch` | Regex matching |
| `-replace` | Regex replacement |
| `-in` / `-notin` | Value in collection |
| `-contains` / `-notcontains` | Collection contains value |

### Special Operators

| Operator | Description |
|----------|-------------|
| `&` | Call/invocation operator |
| `.` | Member access / dot-source |
| `::` | Static member access |
| `$( )` | Subexpression |
| `@( )` | Array subexpression (always returns array) |
| `[type]` | Cast operator |
| `|` | Pipeline |
| `??` | Null-coalescing (PS 7+) |
| `??=` | Null-coalescing assignment (PS 7+) |
| `?.` / `?[]` | Null-conditional (PS 7+) |
| `&&` / `||` | Pipeline chain (PS 7+) |
| `..` | Range operator |
| `-f` | Format operator |

## Variables

### Automatic Variables (Selected)

| Variable | Description |
|----------|-------------|
| `$_` / `$PSItem` | Current pipeline object |
| `$null` | Null/empty value |
| `$true` / `$false` | Boolean literals |
| `$PSVersionTable` | PowerShell version info |
| `$PROFILE` | Path to profile file |
| `$HOME` | User's home directory |
| `$PWD` | Current working directory |
| `$LASTEXITCODE` | Exit code of last native command |
| `$Error` | Array of recent errors |
| `$PSBoundParameters` | Dictionary of passed parameters |
| `$IsWindows` / `$IsLinux` / `$IsMacOS` | Platform detection |
| `$args` | Undeclared parameters |

### Type Constraints

```powershell
[int]$number = 42
[string[]]$computers = "Server1", "Server2"
[datetime]$date = "09/12/91"
```

### Scope Modifiers

```powershell
$Global:Var = "global"
$Script:Var = "script scope"
$Local:Var = "local scope"
```

## Functions

### Verb-Noun Convention

PowerShell functions follow `Verb-Noun` naming with approved verbs (`Get-Verb`):

```powershell
function Get-MyData {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory, ValueFromPipeline)]
        [string]$InputData
    )
    process {
        $InputData | ForEach-Object { $_.ToUpper() }
    }
}
```

### Common Parameters

Advanced functions automatically get: `-Verbose`, `-Debug`, `-ErrorAction`, `-WarningAction`, `-InformationAction`, `-ErrorVariable`, `-WhatIf`, `-Confirm` (with `SupportsShouldProcess`).

## References

- [about_Quoting_Rules](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_quoting_rules)
- [about_Parsing](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_parsing)
- [about_Script_Blocks](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_script_blocks)
- [about_Operators](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_operators)
- [about_Variables](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_variables)
- [about_Automatic_Variables](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_automatic_variables)
- [about_Pipelines](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_pipelines)
- [Native Commands in PowerShell](https://devblogs.microsoft.com/powershell/native-commands-in-powershell-a-new-approach/)

## Related Skills

- **powershell** → `references/completion.md` — How argument passing affects completers
- **powershell** → `references/startup-config.md` — Profiles, execution policy
- **carapace-dev** → `references/shell-powershell.md` — Carapace snippet argument handling
