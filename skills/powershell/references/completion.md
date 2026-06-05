# PowerShell Tab Completion

In-depth reference for PowerShell's tab completion system — `Register-ArgumentCompleter`, `TabExpansion2`, `CompletionResult`, the AST-based completion pipeline, and how external tools hook into the completion menu.

## The Completion Flow

When the user presses Tab (or Ctrl+Space for menu completion), PowerShell:

1. Parses the command line into an **Abstract Syntax Tree (AST)** via `Parser.ParseInput()`
2. Calls `TabExpansion2` which wraps `CommandCompletion.CompleteInput()`
3. `CompletionAnalysis.ExtractAstContext()` builds a `CompletionContext` from the AST and tokens
4. `CompletionCompleters` routes to specialized completion methods
5. Custom completers registered via `Register-ArgumentCompleter` are consulted
6. Returns `CompletionResult` objects to PSReadLine for display

### System Flow

```
Parser.ParseInput()
  → CompletionAnalysis.ExtractAstContext()
    → CompletionContext (RelatedAsts, TokenAtCursor, WordToComplete, CustomArgumentCompleters, NativeArgumentCompleters)
      → CompletionCompleters routing
        → CompleteCommand() | CompleteCommandParameter() | CompleteCommandArgument()
        → CompleteVariable() | CompleteType() | CompleteFilename() | CompleteMember() | CompleteOperator()
          → List<CompletionResult>
            → CommandCompletion (matches, replacementIndex, replacementLength)
              → PSReadLine renders results
```

## TabExpansion2

`TabExpansion2` is the built-in PowerShell function that wraps `CommandCompletion.CompleteInput()`. It is the primary entry point for tab completion.

### Syntax

```powershell
# String-based (default)
TabExpansion2 [-inputScript] <String> [[-cursorColumn] <Int32>] [[-options] <Hashtable>]

# AST-based
TabExpansion2 [-ast] <Ast> [-tokens] <Token[]> [-positionOfCursor] <IScriptPosition> [[-options] <Hashtable>]
```

### Options Hashtable

| Key | Description |
|-----|-------------|
| `IgnoreHiddenShares` | Skip hidden UNC shares (e.g., `\\COMPUTER\ADMIN$`) |
| `RelativePaths` | Force relative paths instead of absolute |
| `LiteralPaths` | Prevent replacement of special file characters (wildcards) |

### CommandCompletion Class

`System.Management.Automation.CommandCompletion` is the core class:

| Property | Type | Description |
|----------|------|-------------|
| `CompletionMatches` | `Collection<CompletionResult>` | All completion candidates |
| `CurrentMatchIndex` | `int` | Index of currently selected match |
| `ReplacementIndex` | `int` | Start index for replacement in the input |
| `ReplacementLength` | `int` | Length of text to replace |

Key method: `CompleteInput(input, cursorIndex, options)` — static, parses input and returns `CommandCompletion`.

### Overriding TabExpansion2

You can replace the default `TabExpansion2` function to customize completion behavior globally:

```powershell
function TabExpansion2 {
    param([string]$inputScript, [int]$cursorColumn, [hashtable]$options)
    # Custom completion logic
    # Must return CommandCompletion object
}
```

This is how modules like `Microsoft.PowerShell.UnixTabCompletion` integrate — they replace `TabExpansion2` to provide completions for native Linux/macOS commands.

## Register-ArgumentCompleter

The primary mechanism for registering custom completers. Available since PowerShell 5.0.

### Two Modes

**Native mode** (`-Native`): For external executables where PowerShell can't infer parameter names. Script block receives 3 parameters:

```powershell
param($wordToComplete, $commandAst, $cursorPosition)
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `$wordToComplete` | `string` | Text the user typed before pressing Tab |
| `$commandAst` | `CommandAst` | AST object for the parsed command line |
| `$cursorPosition` | `int` | 0-based offset of the cursor in the input string |

**Parameter mode** (`-ParameterName`): For PowerShell cmdlets/functions. Script block receives 5 parameters:

```powershell
param($commandName, $parameterName, $wordToComplete, $commandAst, $fakeBoundParameters)
```

| Parameter | Type | Description |
|-----------|------|-------------|
| `$commandName` | `string` | Name of the command being completed |
| `$parameterName` | `string` | Parameter whose value needs completion |
| `$wordToComplete` | `string` | Text before Tab was pressed |
| `$commandAst` | `CommandAst` | AST for the current input line |
| `$fakeBoundParameters` | `IDictionary` | Hashtable of already-bound parameters (`$PSBoundParameters`) |

### NativeFallback Mode

PowerShell 7.6+ adds `-NativeFallback` — a cover-all completer for native commands without specific completers. Only one can be registered at a time. This enables modules like `Microsoft.PowerShell.UnixTabCompletion` to provide completions for many native commands at once.

### Registration Examples

**Native command completer:**

```powershell
$scriptblock = {
    param($wordToComplete, $commandAst, $cursorPosition)
    dotnet complete --position $cursorPosition $commandAst.ToString() | ForEach-Object {
        [System.Management.Automation.CompletionResult]::new(
            $_,               # completionText
            $_,               # listItemText
            'ParameterValue', # resultType
            $_                # toolTip
        )
    }
}
Register-ArgumentCompleter -Native -CommandName dotnet -ScriptBlock $scriptblock
```

**PowerShell parameter completer:**

```powershell
Register-ArgumentCompleter -CommandName Set-TimeZone -ParameterName Id -ScriptBlock {
    param($commandName, $parameterName, $wordToComplete, $commandAst, $fakeBoundParameters)
    (Get-TimeZone -ListAvailable).Id | Where-Object {
        $_ -like "$wordToComplete*"
    }
}
```

**Multiple command names:**

```powershell
Register-ArgumentCompleter -Native -CommandName 'git','git.exe' -ScriptBlock $scriptblock
```

### Critical: Pipeline Unrolling

The script block **must unroll values using the pipeline** (`ForEach-Object`, `Where-Object`, or direct pipeline). Returning an array directly causes PowerShell to treat the entire array as **one** completion value.

```powershell
# WRONG - entire array becomes one completion
return @("a", "b", "c")

# CORRECT - each value becomes a separate completion
"a", "b", "c" | ForEach-Object { [CompletionResult]::new($_, $_, 'ParameterValue', $_) }
```

## The CommandAst Object

`CommandAst` (from `System.Management.Automation.Language`) is the AST node for the current command. It provides properly tokenized, dequoted arguments — unlike bash where `COMP_WORDBREAKS` causes word-splitting problems.

### Key Properties

| Property | Type | Description |
|----------|------|-------------|
| `CommandElements` | `CommandElementAstCollection` | Ordered collection of command name + arguments |
| `Extent` | `IScriptExtent` | The span of text this AST covers in the input |
| `Parent` | `Ast` | Parent AST node |
| `Redirections` | `ReadOnlyCollection<RedirectionAst>` | Any redirections in the command |

### CommandElementAst Properties

Each element in `CommandElements` has:

| Property | Type | Description |
|----------|------|-------------|
| `Extent.Text` | `string` | Raw text of the element including any quotes |
| `Extent.StartOffset` | `int` | 0-based start position in the input |
| `Extent.EndOffset` | `int` | 0-based end position in the input |

### AST Tokenization

PowerShell's AST tokenizer properly handles:
- Double-quoted strings (with variable expansion) as single elements
- Single-quoted strings (literal) as single elements
- Subexpressions (`$(...)`) as single elements
- Redirections are separated from command elements

This means `CommandElements` gives properly parsed tokens without the word-splitting issues that plague bash's `COMP_WORDS`.

### Accessing Arguments from CommandAst

```powershell
param($wordToComplete, $commandAst, $cursorPosition)
$commandElements = $commandAst.CommandElements

# Skip first element (command name), get arguments
$arguments = $commandElements | Select-Object -Skip 1

# Get specific argument by position
$firstArg = $commandElements[1].Extent.Text
```

## CompletionResult

`System.Management.Automation.CompletionResult` represents a single completion candidate.

### Constructor

```powershell
[CompletionResult]::new($completionText, $listItemText, $resultType, $toolTip)
```

### Properties

| Property | Description |
|----------|-------------|
| `CompletionText` | Text inserted into the command line when selected |
| `ListItemText` | Text displayed in the completion list (Ctrl+Space). Can differ from `CompletionText` — enables showing human-readable text while inserting machine-readable values |
| `ResultType` | `CompletionResultType` enum value controlling ending-key behavior |
| `ToolTip` | Tooltip text shown when the item is selected in the menu. Also used for description display when `ShowToolTips` is enabled |

### Simple String Returns

Script blocks can return plain strings instead of `CompletionResult` objects. PowerShell wraps them automatically with `ResultType = ParameterValue`. However, `CompletionResult` objects are needed for:
- Different display text vs. insertion text
- Custom tooltips/descriptions
- Specific `ResultType` values that control ending-key behavior

## CompletionResultType

The `CompletionResultType` enum controls what happens after a completion is selected — specifically which keys terminate completion mode and get inserted.

### All Values

| Value | Name | Description | Ending Keys |
|-------|------|-------------|-------------|
| 0 | `Text` | Unknown result type | None specific |
| 1 | `History` | History item | None specific |
| 2 | `Command` | Command name | None specific |
| 3 | `ProviderItem` | Provider item (file) | None specific |
| 4 | `ProviderContainer` | Container (directory) | `\`, `/` |
| 5 | `Property` | Object property | `.` |
| 6 | `Method` | Object method | `(`, `)` |
| 7 | `ParameterName` | Cmdlet parameter | `:` |
| 8 | `ParameterValue` | Parameter value | `,` |
| 9 | `Variable` | Variable name | `.` |
| 10 | `Namespace` | .NET namespace | `.` |
| 11 | `Type` | .NET type name | `]` |
| 12 | `Keyword` | Language keyword | None specific |
| 13 | `DynamicKeyword` | Dynamic keyword | None specific |

### How Ending Keys Work

When a completion is accepted, PSReadLine checks if the key that triggered acceptance is an "ending key" for the completion's `ResultType`. If it is, the key is inserted after the completion text. This enables natural continuation:

- `ProviderContainer` + `\` → `C:\Windows\` (continue into directory)
- `Property` + `.` → `$obj.Property.` (trigger member completion)
- `Method` + `(` → `$obj.Method(` (start arguments)
- `Variable` + `.` → `$var.` (trigger member completion)
- `Type` + `]` → `[System.String]` (close type literal)
- `ParameterName` + `:` → `-Path:` (colon syntax)
- `ParameterValue` + `,` → `Value1,` (multiple values)

### ProviderContainer Special Behavior

For `CompletionResultType.ProviderContainer`, PSReadLine **automatically appends a directory separator** if the completion text doesn't already end with one. This allows the user to continue typing into subdirectories without manually adding the separator.

### Space After Completion

Unlike POSIX shells, PowerShell does **not** automatically append a space after most completion types. This is a known difference — POSIX shells insert a space when the completion is a self-contained argument. PowerShell's behavior means the user may need to press Space manually after accepting a completion.

## CompletionCompleters Routing

`CompletionCompleters` is the internal class that routes to specialized completion methods:

| Method | Purpose |
|--------|---------|
| `CompleteCommand()` | Command name completion (cmdlets, functions, aliases, executables) |
| `CompleteCommandParameter()` | Parameter name completion for cmdlets |
| `CompleteCommandArgument()` | Parameter value completion (invokes registered completers) |
| `CompleteVariable()` | Variable name completion (`$var`) |
| `CompleteType()` | Type and namespace completion |
| `CompleteFilename()` | File and provider path completion |
| `CompleteMember()` | Property/method completion after `.` or `::` |
| `CompleteOperator()` | Operator completion |

### PseudoParameterBinder

For parameter completion, `PseudoParameterBinder` performs lightweight parameter binding:
- Determines which parameters are available for a command
- Identifies already-bound parameters
- Returns `UnboundParameters` for completion candidates

### Type Inference for Member Completion

Member completion (after `.` or `::`) uses `TypeInferenceVisitor`:
1. Extracts the expression AST before the operator
2. Calls `AstTypeInference.InferTypeOf()` to determine types
3. Returns matching members for the inferred types

## Other Completion Registration Mechanisms

### ArgumentCompleter Attribute

Applied directly to function parameters:

```powershell
function Get-Data {
    param(
        [ArgumentCompleter({
            param($commandName, $parameterName, $wordToComplete, $commandAst, $fakeBoundParameters)
            'One', 'Two', 'Three' | Where-Object { $_ -like "$wordToComplete*" }
        })]
        [string]$Stage
    )
}
```

### ArgumentCompletions Attribute

Static completion values without validation (PowerShell 7.5+):

```powershell
function Get-Data {
    param(
        [ArgumentCompletions('One', 'Two', 'Three')]
        [string]$Stage
    )
}
```

Unlike `ValidateSet`, `ArgumentCompletions` provides completions without restricting valid values.

### ValidateSet Attribute

Restricts valid values and provides completions:

```powershell
function Get-Data {
    param(
        [ValidateSet('One', 'Two', 'Three')]
        [string]$Stage
    )
}
```

### Enum Parameters

PowerShell automatically provides completions for enum-typed parameters:

```powershell
function Set-Color {
    param([System.ConsoleColor]$Color)
}
# Tab completion automatically lists: Black, DarkBlue, DarkGreen, ...
```

**Limitation**: `Register-ArgumentCompleter` cannot override completion for enum-typed parameters — the type's enum values always take precedence.

## How External Tools Hook Into Completion

### Pattern 1: External CLI with Completion Subcommand

The most common pattern for external tools (used by `dotnet`, `carapace`, etc.):

```powershell
Register-ArgumentCompleter -Native -CommandName mytool -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    mytool complete --position $cursorPosition $commandAst.ToString() | ForEach-Object {
        [CompletionResult]::new($_, $_, 'ParameterValue', $_)
    }
}
```

The external tool provides a subcommand (e.g., `complete`, `_carapace`) that:
- Receives the command line and cursor position
- Returns completion candidates (typically as lines or JSON)
- The script block converts them to `CompletionResult` objects

### Pattern 2: PowerShell Module with Custom Logic

Modules can register completers in their `.psm1` files:

```powershell
Register-ArgumentCompleter -Native -CommandName docker -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    # Parse CommandAst to determine subcommand context
    $elements = $commandAst.CommandElements
    # ... context-aware completion logic
}
```

### Pattern 3: NativeFallback Cover-All Completer

PowerShell 7.6+ allows a single completer for all native commands:

```powershell
Register-ArgumentCompleter -NativeFallback -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    # Provide completions for any native command
}
```

### Pattern 4: TabExpansion2 Replacement

For deep integration, replace `TabExpansion2` entirely:

```powershell
function TabExpansion2 {
    param([string]$inputScript, [int]$cursorColumn, [hashtable]$options)
    # Custom completion logic
    # Return CommandCompletion object
}
```

This is how `Microsoft.PowerShell.UnixTabCompletion` works — it replaces the entire completion pipeline for native commands on Linux/macOS.

## Debugging Completions

### Testing Completions Programmatically

```powershell
# Test a completion directly
TabExpansion2 -inputScript 'Get-Process -' -cursorColumn 15 |
    Select-Object -ExpandProperty CompletionMatches

# Test with cursor position
$s = 'Get-Process -Na'
TabExpansion2 -inputScript $s -cursorColumn $s.Length
```

### Trace-Command for Parameter Binding

```powershell
Trace-Command -Name ParameterBinding -Expression { Get-Process -Name a } -PSHost
```

### Common Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| Array returned as one completion | Script block returns array directly | Use pipeline (`ForEach-Object`) to unroll |
| Empty CompletionResult throws error | .NET constructor rejects empty strings | Ensure all fields are non-empty (use space for empty) |
| Enum parameter completions can't be overridden | Type system takes precedence | Change parameter type or use `ArgumentCompleter` on non-enum type |
| Completions not showing | Execution policy blocks script | Check `Get-ExecutionPolicy` |
| Single-quoted prefix doesn't match | AST includes quotes in `Extent.Text` | Strip quotes before matching |

## References

- [Register-ArgumentCompleter](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/register-argumentcompleter)
- [about_Tab_Expansion](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_tab_expansion)
- [TabExpansion2](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/tabexpansion2)
- [CompletionResult Class](https://learn.microsoft.com/en-us/dotnet/api/system.management.automation.completionresult)
- [CompletionResultType Enum](https://learn.microsoft.com/en-us/dotnet/api/system.management.automation.completionresulttype)
- [CommandAst Class](https://learn.microsoft.com/en-us/dotnet/api/system.management.automation.language.commandast)
- [CommandCompletion Source](https://github.com/PowerShell/PowerShell/blob/master/src/System.Management.Automation/engine/CommandCompletion/CommandCompletion.cs)

## Related Skills

- **carapace-dev** → `references/shell-powershell.md` — Carapace-specific PowerShell integration (snippet, value formatting, SGR styling, nospace)
- **powershell** → `references/editor.md` — PSReadLine, menu completion, prediction system
- **powershell** → `references/styling.md` — ANSI escape sequences, $PSStyle, terminal rendering
