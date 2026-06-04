# Carapace Library: PowerShell Shell Integration Deep Dive

Reference for [carapace](https://github.com/carapace-sh/carapace)'s PowerShell completion integration — how the snippet works, how completion output is formatted, and how carapace handles PowerShell-specific edge cases including the AST-based argument parsing, CompletionResult construction, single-quote stripping, empty string rejection, SGR color rendering, and nospace handling. For cross-shell comparisons, see the **references/shell.md**.

## Source Files

| File | Purpose |
|------|---------|
| `internal/shell/powershell/snippet.go` | PowerShell completion script generation (`Register-ArgumentCompleter -Native`) |
| `internal/shell/powershell/action.go` | Value formatting, JSON `completionResult` serialization, SGR styling, nospace |
| `internal/shell/shell.go` | Shared dispatch — message integration, nospace propagation, filtering |
| `complete.go` | Entry point — no PowerShell-specific patching (unlike bash/nushell/cmd-clink) |
| `pkg/ps/ps.go` | Shell detection — maps `powershell`/`pwsh` process names to `"powershell"` |

## PowerShell Completion System Background

### The `Register-ArgumentCompleter` Mechanism

PowerShell registers completion functions via `Register-ArgumentCompleter`. When the user presses Tab (or Ctrl+Space for menu completion), PowerShell:

1. Parses the command line into an **Abstract Syntax Tree (AST)**
2. Determines the command being completed
3. Looks up the registered argument completer for that command
4. Invokes the completer's script block with context parameters
5. Reads the returned `CompletionResult` objects for display and insertion

### Two Modes of `Register-ArgumentCompleter`

**PowerShell command mode** (`-ParameterName`): For PowerShell cmdlets/functions where parameter names are known. The script block receives 5 parameters:

```powershell
param($commandName, $parameterName, $wordToComplete, $commandAst, $fakeBoundParameters)
```

**Native command mode** (`-Native`): For external executables where PowerShell can't infer parameter names. The script block receives 3 parameters:

```powershell
param($wordToComplete, $commandAst, $cursorPosition)
```

Carapace uses the **native command mode** since it provides completions for external CLI tools, not PowerShell cmdlets.

### Script Block Parameters (Native Mode)

| Parameter | Type | Description |
|-----------|------|-------------|
| `$wordToComplete` | `string` | The text the user has typed before pressing Tab |
| `$commandAst` | `CommandAst` | The AST object representing the parsed command line |
| `$cursorPosition` | `int` | The 0-based offset of the cursor in the input string |

### The `CommandAst` Object

`CommandAst` (from `System.Management.Automation.Language`) is the AST node for the current command. Key properties:

| Property | Type | Description |
|----------|------|-------------|
| `CommandElements` | `CommandElementAstCollection` | Ordered collection of command elements (command name + arguments) |
| `Extent` | `IScriptExtent` | The span of text this AST covers in the input |

Each `CommandElementAst` in `CommandElements` has:

| Property | Type | Description |
|----------|------|-------------|
| `Extent.Text` | `string` | The raw text of the element including any quotes |
| `Extent.StartOffset` | `int` | 0-based start position in the input |
| `Extent.EndOffset` | `int` | 0-based end position in the input |

PowerShell's AST tokenizer properly handles quoting — double-quoted strings, single-quoted strings, and subexpressions are all parsed as single command elements. This means `CommandElements` gives carapace properly tokenized arguments, unlike bash where `COMP_WORDBREAKS` causes word-splitting problems.

### The `CompletionResult` Class

`CompletionResult` (from `System.Management.Automation`) represents a single completion candidate:

```csharp
public class CompletionResult {
    public string CompletionText { get; }   // Text inserted into the command line
    public string ListItemText { get; }     // Text displayed in the completion list
    public CompletionResultType ResultType { get; }  // Type of completion
    public string ToolTip { get; }          // Tooltip shown when item is selected
}
```

Constructor:

```powershell
[CompletionResult]::new($completionText, $listItemText, $resultType, $toolTip)
```

| Parameter | Description |
|-----------|-------------|
| `CompletionText` | The text actually inserted into the command line when selected |
| `ListItemText` | The text displayed in the tab-completion list (Ctrl+Space). Can differ from `CompletionText` — enables showing human-readable text while inserting machine-readable values |
| `ResultType` | The `CompletionResultType` enum value. Carapace always uses `[CompletionResultType]::ParameterValue` |
| `ToolTip` | The tooltip text shown when the item is selected in the completion list. Also used for description display when `CARAPACE_TOOLTIP` is enabled |

### PSReadLine and Completion Display

PSReadLine is PowerShell's line-editing module that provides the completion UI:

- **Tab** cycles through completions inline
- **Ctrl+Space** shows a menu list of all completions with descriptions
- The menu list uses `ListItemText` for display
- Tooltips (from `ToolTip` field) appear when a completion is selected in the menu
- PSReadLine handles common-prefix auto-insertion when multiple candidates share a prefix

### The `` `e `` Escape Character

PowerShell 6+ supports `` `e `` (backtick-e) as the **Escape character** (ASCII 27 / `\x1b`). This is used to embed ANSI/VT100 escape sequences (SGR codes) in PowerShell strings:

```powershell
"`e[31mRed text`e[0m"   # SGR code for red foreground, then reset
```

`` `e[ `` is the CSI (Control Sequence Introducer) — equivalent to `\x1b[` or `ESC[` in other contexts. This is critical for carapace's color rendering in `ListItemText` and `ToolTip` fields.

## The PowerShell Snippet

Generated by `powershell.Snippet(cmd)` in `internal/shell/powershell/snippet.go`:

```powershell
using namespace System.Management.Automation
using namespace System.Management.Automation.Language
Function _example_completer {
    [System.Diagnostics.CodeAnalysis.SuppressMessageAttribute("PSAvoidUsingInvokeExpression", "", Scope="Function", Target="*")]
    param($wordToComplete, $commandAst, $cursorPosition)
    $commandElements = $commandAst.CommandElements

    # double quoted value works but seems single quoted needs some fixing (e.g. "example 'acti" -> "example acti")
    $elems = @()
    foreach ($_ in $commandElements) {
      if ($_.Extent.StartOffset -gt $cursorPosition) {
          break
      }
      $t = $_.Extent.Text
      if ($_.Extent.EndOffset -gt $cursorPosition) {
          $t = $t.Substring(0, $_.Extent.Text.get_Length() - ($_.Extent.EndOffset - $cursorPosition))
      }

      if ($t.Substring(0,1) -eq "'"){
        $t = $t.Substring(1)
      }
      if ($t.get_Length() -gt 0 -and $t.Substring($t.get_Length()-1) -eq "'"){
        $t = $t.Substring(0,$t.get_Length()-1)
      }
      if ($t.get_Length() -eq 0){
        $t = '""'
      }
      $elems += $t.replace('`',', ',')
    }

    $completions = @(
      if (!$wordToComplete) {
        example _carapace powershell $($elems| ForEach-Object {$_}) '' | ConvertFrom-Json | ForEach-Object { [CompletionResult]::new($_.CompletionText, $_.ListItemText.replace("`e[", "`e["), [CompletionResultType]::ParameterValue, $_.ToolTip.replace("`e[", "`e[")) }
      } else {
        example _carapace powershell $($elems| ForEach-Object {$_}) | ConvertFrom-Json | ForEach-Object { [CompletionResult]::new($_.CompletionText, $_.ListItemText.replace("`e[", "`e["), [CompletionResultType]::ParameterValue, $_.ToolTip.replace("`e[", "`e[")) }
      }
    )

    if ($completions.count -eq 0) {
      return "" # prevent default file completion
    }

    $completions
}
Register-ArgumentCompleter -Native -ScriptBlock (Get-Item "Function:_example_completer").ScriptBlock -CommandName 'example','example.exe'
```

### Snippet Walkthrough

**1. Namespace imports**

```powershell
using namespace System.Management.Automation
using namespace System.Management.Automation.Language
```

Enables short-form type references like `[CompletionResult]` and `[CompletionResultType]::ParameterValue` without fully-qualified names.

**2. Suppress `PSAvoidUsingInvokeExpression` warning**

```powershell
[System.Diagnostics.CodeAnalysis.SuppressMessageAttribute("PSAvoidUsingInvokeExpression", "", Scope="Function", Target="*")]
```

PowerShell's PSScriptAnalyzer flags functions that might use `Invoke-Expression`. Carapace's function pipes output from an external command through `ConvertFrom-Json`, which some analyzer configurations misidentify. The suppression attribute silences the warning.

**3. Parse `CommandElements` with cursor-aware truncation**

```powershell
$commandElements = $commandAst.CommandElements
$elems = @()
foreach ($_ in $commandElements) {
  if ($_.Extent.StartOffset -gt $cursorPosition) {
      break
  }
  $t = $_.Extent.Text
  if ($_.Extent.EndOffset -gt $cursorPosition) {
      $t = $t.Substring(0, $_.Extent.Text.get_Length() - ($_.Extent.EndOffset - $cursorPosition))
  }
  ...
  $elems += $t.replace('`',', ',')
}
```

The loop iterates over all command elements (command name + arguments) in AST order:

- **Skip future elements**: Elements that start after the cursor position are irrelevant — `break` stops iteration
- **Truncate the current element**: If the element extends past the cursor (the user is mid-word), truncate it to just the portion before the cursor. This is how carapace gets the partial word being completed
- **Comma fix**: `` $t.replace('`',', ',') `` replaces PowerShell's backtick-comma escape (`` `, ``) with a literal comma. PowerShell uses `` `, `` to escape commas inside argument lists; the replace restores the original comma character

**4. Strip single quotes from token boundaries**

```powershell
if ($t.Substring(0,1) -eq "'"){
  $t = $t.Substring(1)
}
if ($t.get_Length() -gt 0 -and $t.Substring($t.get_Length()-1) -eq "'"){
  $t = $t.Substring(0,$t.get_Length()-1)
}
```

Strips a leading and/or trailing single quote from each token. The comment in the source explains: *"double quoted value works but seems single quoted needs some fixing (e.g. "example 'acti" -> "example acti")"*.

**Why this is needed**: When the user types `example 'acti⇥`, PowerShell's AST includes the single quotes as part of the `Extent.Text` (e.g., `'acti`). But carapace's traversal engine expects unquoted argument values — it matches against the raw value without quotes. If single quotes are not stripped, the completion prefix `'acti` won't match candidates like `action`, `add`, etc.

**Why only single quotes**: Double-quoted values "just work" — PowerShell's `Extent.Text` for double-quoted strings preserves the quotes, but the pipeline `| ForEach-Object {$_}` and `xargs`-style processing handles them correctly. Single-quoted strings are different because PowerShell's single quotes are literal (no interpolation), and the AST text includes them verbatim.

**5. Replace empty tokens with `""`**

```powershell
if ($t.get_Length() -eq 0){
  $t = '""'
}
```

After stripping quotes, a token that was just `''` (empty single-quoted string) becomes empty. The empty string is replaced with `""` (empty double-quoted string) so it's passed as an empty argument to the carapace subprocess rather than being swallowed by PowerShell's argument handling.

**6. Build the carapace command with two branches**

```powershell
$completions = @(
  if (!$wordToComplete) {
    example _carapace powershell $($elems| ForEach-Object {$_}) '' | ConvertFrom-Json | ...
  } else {
    example _carapace powershell $($elems| ForEach-Object {$_}) | ConvertFrom-Json | ...
  }
)
```

Two branches based on `$wordToComplete`:

- **Empty word (`!$wordToComplete`)**: Appends an empty string `''` as the last argument. This tells carapace that the cursor is at a new argument position (nothing has been typed yet), so carapace should suggest all valid values for that position
- **Non-empty word**: Passes only the parsed elements. The last element (truncated by cursor position) is the partial word being completed, which carapace uses for prefix filtering

`$($elems| ForEach-Object {$_})` expands the `$elems` array into individual arguments on the command line. This is PowerShell's equivalent of `"${array[@]}"` in bash — it prevents the array from being passed as a single comma-separated string.

**7. Construct `CompletionResult` from JSON**

```powershell
... | ConvertFrom-Json | ForEach-Object {
  [CompletionResult]::new(
    $_.CompletionText,
    $_.ListItemText.replace("`e[", "`e["),
    [CompletionResultType]::ParameterValue,
    $_.ToolTip.replace("`e[", "`e[")
  )
}
```

The carapace subprocess outputs a JSON array of `completionResult` objects. The pipeline:

1. `ConvertFrom-Json` parses the JSON into PowerShell `PSCustomObject` array
2. `ForEach-Object` iterates over each result
3. `[CompletionResult]::new(...)` constructs the completion result with 4 arguments

**The `` .replace("`e[", "`e[") `` mystery**: This replace looks like a no-op (replacing a string with itself). But it's actually critical — PowerShell's `ConvertFrom-Json` may mangle or unescape the `` `e `` escape character in JSON strings. The `.replace()` call ensures the `` `e[ `` escape sequences survive the JSON deserialization intact. Without this, SGR color codes in `ListItemText` and `ToolTip` would be broken.

**8. Prevent default file completion on empty results**

```powershell
if ($completions.count -eq 0) {
  return "" # prevent default file completion
}
```

When no completions are found, PowerShell would normally fall back to file/path completion. Returning an empty string prevents this — the user sees no completions rather than irrelevant file suggestions.

**9. Register with `Register-ArgumentCompleter -Native`**

```powershell
Register-ArgumentCompleter -Native -ScriptBlock (Get-Item "Function:_example_completer").ScriptBlock -CommandName 'example','example.exe'
```

- **`-Native`**: Registers for a native (external) command, not a PowerShell cmdlet
- **`-ScriptBlock`**: Uses `(Get-Item "Function:_name").ScriptBlock` to retrieve the function's script block object. This is necessary because `Register-ArgumentCompleter` requires a `[scriptblock]` type, not a function name
- **`-CommandName`**: Registers for both `example` and `example.exe`. On Windows, native commands may be invoked with or without the `.exe` extension. Both must be registered

**Windows-specific**: On Windows, the snippet adds a `# ` comment prefix before `,'example.exe'` to effectively skip the `.exe` registration (since the comma-separated list syntax works differently). On non-Windows, `example.exe` is included as a fallback. The `prefix` variable in Go code controls this:

```go
prefix := " # "
if runtime.GOOS == "windows" {
    prefix = ""
}
```

## No PowerShell-Side Patching

Unlike bash (which needs `bash.Patch()` for redirect filtering and COMP_WORDBREAKS handling), PowerShell requires **no Go-side argument patching** in `complete.go`. The PowerShell snippet handles all argument pre-processing client-side:

| Problem | Bash | PowerShell |
|---------|------|------------|
| Redirect tokens in args | `bash.Patch()` strips `>`, `2>`, etc. | Not needed — PowerShell's AST parser handles redirects |
| Word splitting (`COMP_WORDBREAKS`) | `bash.Patch()` re-lexes `COMP_LINE` to reconstruct split words | Not needed — PowerShell's AST preserves quoted strings intact |
| Cursor position | `COMP_POINT` env var | `$cursorPosition` parameter + `Extent` offsets |
| Open quotes | Snippet retry with `''`, `'"`, `"` | Not needed — PowerShell's AST tokenizer handles quotes natively |
| Single-quote boundary | N/A | Snippet strips leading/trailing `'` from `Extent.Text` |

PowerShell's AST-based parsing gives carapace properly tokenized arguments without the word-splitting issues that plague bash. The `CommandElements` collection provides parsed, dequoted tokens with offset information — eliminating the need for `shlex` re-lexing or redirect filtering.

## Edge Cases and How Carapace Handles Them

### Edge Case 1: Empty `CompletionResult` Fields Cause Errors

**Problem**: PowerShell's `[CompletionResult]::new()` constructor throws an error if any parameter is an empty string. This is a .NET API limitation — the constructor validates that `CompletionText`, `ListItemText`, and `ToolTip` are non-empty.

**Carapace's solution**:
- Skip values with empty `Value` fields entirely (`continue` in the loop)
- Replace empty `ListItemText` and `ToolTip` with a single space via `ensureNotEmpty()`

```go
func ensureNotEmpty(s string) string {
    if s == "" {
        return " "
    }
    return s
}
```

This is unique to PowerShell — no other shell has this constraint.

### Edge Case 2: Single-Quoted Arguments Need Quote Stripping

**Problem**: When the user types `example 'partial_wor⇥`, PowerShell's AST includes the single quotes in `Extent.Text` (e.g., `'partial_wor`). If carapace receives `'partial_wor` as the completion prefix, it won't match candidates that start with `partial_wor` — the quote breaks prefix matching.

**Carapace's solution**: The snippet strips leading and trailing single quotes from each token:

```powershell
if ($t.Substring(0,1) -eq "'"){
  $t = $t.Substring(1)
}
if ($t.get_Length() -gt 0 -and $t.Substring($t.get_Length()-1) -eq "'"){
  $t = $t.Substring(0,$t.get_Length()-1)
}
```

Double-quoted strings don't need this treatment — they "just work" through the AST parser.

### Edge Case 3: Empty Arguments After Quote Stripping

**Problem**: After stripping quotes from `''` (an empty single-quoted string), the token becomes empty. Passing an empty string as an argument through PowerShell's pipeline can cause it to be swallowed entirely — the carapace subprocess wouldn't see it as an argument position.

**Carapace's solution**: Replace empty tokens with `""` (double-quoted empty string):

```powershell
if ($t.get_Length() -eq 0){
  $t = '""'
}
```

This ensures the empty argument position is preserved when passed to the carapace subprocess.

### Edge Case 4: Comma Escaping in PowerShell

**Problem**: PowerShell uses commas as array separators. When a completion value like `file,name.txt` is passed through the pipeline, PowerShell may interpret it as two separate arguments. PowerShell also uses the backtick-comma sequence (`` `, ``) as an escape for commas inside argument lists.

**Carapace's solution**: The snippet replaces PowerShell's comma escape (`` `, ``) with a literal comma:

```powershell
$t.replace('`',', ',')
```

This handles the case where PowerShell has escaped a comma with a backtick in the `Extent.Text`.

### Edge Case 5: Special Characters in Completion Values

**Problem**: When a completion value contains characters that are special in PowerShell (spaces, `$`, `{}`, `()`, etc.), inserting the raw value would cause PowerShell to interpret it as syntax rather than a literal string.

**Carapace's solution**: Wrap values containing special characters in single quotes:

```go
if strings.ContainsAny(val.Value, ` {}()[]*$?\"|<>&(),;#`+"`") {
    val.Value = fmt.Sprintf("'%v'", val.Value)
}
```

Single quotes in PowerShell treat content literally — no variable expansion, no escape interpretation (except `''` for a literal single quote). This is similar to how bash uses double-quoting and xonsh uses Python string quoting.

**Known limitation**: If a completion value itself contains a single quote (e.g., `it's`), the current quoting would produce `'it's'` which is invalid in PowerShell. The sanitizer does not escape single quotes within values. This is noted by the `// TODO` comment on the sanitizer.

### Edge Case 6: Nospace is Embedded in CompletionText

**Problem**: PowerShell has no native per-candidate "don't add trailing space" mechanism like zsh's `CodeSuffix` or elvish's `suffix` field. `CompletionResult` doesn't have a nospace property.

**Carapace's solution**: Bake the trailing space directly into `CompletionText`:

- Normal completion: `val.Value = val.Value + " "` — trailing space separates from next argument
- Nospace completion: leave `CompletionText` as-is — cursor stays right after the value

This is simple and effective. PSReadLine respects the full `CompletionText` when inserting, so the embedded space works correctly.

### Edge Case 7: SGR Escape Sequences Through JSON Deserialization

**Problem**: Carapace embeds SGR escape sequences in `ListItemText` and `ToolTip` using PowerShell's `` `e[ `` notation (backtick-e-bracket). But when these strings pass through `ConvertFrom-Json`, the JSON deserializer may mangle or unescape the `` `e `` character.

**Carapace's solution**: The snippet applies `.replace("`e[", "`e[")` after `ConvertFrom-Json`:

```powershell
$_.ListItemText.replace("`e[", "`e[")
$_.ToolTip.replace("`e[", "`e[")
```

While this looks like a no-op (replacing a string with itself), it forces PowerShell to re-evaluate the `` `e[ `` escape sequence after JSON deserialization, ensuring the escape character is properly interpreted as the ESC character (ASCII 27).

### Edge Case 8: Default File Completion Prevention

**Problem**: When no completions are found, PowerShell's default behavior is to fall back to file/path completion. This is often undesirable — if carapace returns no completions for a flag value, the user doesn't want to see random file suggestions.

**Carapace's solution**: Return an empty string when `completions.count` is zero:

```powershell
if ($completions.count -eq 0) {
    return "" # prevent default file completion
}
```

This tells PSReadLine "I handled this completion but found nothing" rather than "I didn't handle this, fall back to defaults."

### Edge Case 9: `.exe` Extension on Windows

**Problem**: On Windows, native commands may be invoked as either `example` or `example.exe`. PowerShell treats these as different command names, so a completer registered for `example` won't activate for `example.exe`.

**Carapace's solution**: The `Register-ArgumentCompleter` call includes both names:

```powershell
Register-ArgumentCompleter -Native -ScriptBlock ... -CommandName 'example','example.exe'
```

On non-Windows platforms, the `.exe` extension is commented out (using `# ` prefix) since it's irrelevant:

```go
prefix := " # "
if runtime.GOOS == "windows" {
    prefix = ""
}
```

### Edge Case 10: Cursor-Aware Token Truncation

**Problem**: When the cursor is in the middle of a token (e.g., `example --fla⇥g`), the full `Extent.Text` would include characters after the cursor. Passing the entire token would cause carapace to filter on the wrong prefix.

**Carapace's solution**: Truncate the current token at the cursor position:

```powershell
if ($_.Extent.EndOffset -gt $cursorPosition) {
    $t = $t.Substring(0, $_.Extent.Text.get_Length() - ($_.Extent.EndOffset - $cursorPosition))
}
```

This calculates the portion of the token before the cursor and passes only that prefix to carapace for completion matching. This is similar to how bash uses `COMP_POINT` and fish uses `commandline -cp` — but PowerShell's AST provides this information directly via `Extent` offsets.

## Related Skills

- **references/shell.md** — cross-shell feature comparison and shared dispatch
- **references/traverse.md** — the completion engine that produces Actions before shell formatting
- **references/style.md** — how styles are resolved before SGR rendering
