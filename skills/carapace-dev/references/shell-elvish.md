# Carapace Library: Elvish Shell Integration Deep Dive

Reference for [carapace](https://github.com/carapace-sh/carapace)'s elvish completion integration — how the snippet works, how completion output is formatted, and how carapace handles elvish-specific edge cases including JSON output via `from-json`, `edit:complex-candidate` with `CodeSuffix` for nospace, `styled` for per-candidate colors, and `edit:notify` for error messages. For cross-shell comparisons, see the **references/shell.md**.

## Source Files

| File | Purpose |
|------|---------|
| `internal/shell/elvish/snippet.go` | Elvish completion script generation (`edit:completion:arg-completer`) |
| `internal/shell/elvish/action.go` | Value formatting, JSON serialization, `complexCandidate` struct, nospace via `CodeSuffix` |
| `third_party/github.com/elves/elvish/pkg/ui/` | Vendored elvish `ui` package — `ParseStyling`, `Style`, `Color` types |
| `internal/shell/shell.go` | Shared dispatch — message handling (elvish uses native `edit:notify`), nospace propagation |
| `complete.go` | Entry point — no elvish-specific patching (unlike bash/nushell/cmd-clink) |

## Elvish Completion System Background

### The `arg-completer` Mechanism

Elvish registers completion functions via the `$edit:completion:arg-completer` map variable. When the user presses TAB, elvish:

1. Looks up `$edit:completion:arg-completer[$command-name]`
2. Calls the function with all command-line words as arguments (including the command name as the first argument)
3. The completer outputs candidates via `put`, byte output, or `edit:complex-candidate`
4. Elvish's matcher system filters candidates against what the user typed
5. Matched candidates are displayed in the completion UI

```elvish
# Registration
set edit:completion:arg-completer[example] = {|@arg|
    # $arg is [example subcommand --flag value...]
    # output completions via put, echo, or edit:complex-candidate
}
```

When completing a new argument (cursor after space), elvish passes a trailing empty string:

```elvish
# User types: example action <Space> ⇥
# Completer called: $edit:completion:arg-completer[example] example action ""
```

### The `edit:complex-candidate` Command

`edit:complex-candidate` builds a rich candidate object for argument completers:

```elvish
edit:complex-candidate $stem &display='' &code-suffix=''
```

| Parameter | Type | Purpose |
|-----------|------|---------|
| `$stem` | positional string | The actual value to insert into the command line |
| `&display` | string or styled text | How the candidate appears in the completion UI. If empty, `$stem` is used. Accepts output of `styled` builtin. |
| `&code-suffix` | string | Suffix appended to the quoted stem on insertion. **Not quoted** — added verbatim. Empty = no suffix. Space `' '` = trailing space after insertion. |

**How insertion works**:

1. When a candidate is accepted, elvish inserts a **quoted version** of `$stem`
2. If `&code-suffix` is non-empty, it is **appended verbatim** (not quoted)
3. This gives the completer precise control over what appears after the inserted value

**Examples from elvish's own file completion**:

```elvish
# Directory — no trailing space (the / is part of stem)
edit:complex-candidate d/ &code-suffix='' &display=(styled-segment d/ &fg-color=blue &bold)

# Regular file — trailing space for next argument
edit:complex-candidate bar &code-suffix=' ' &display=[^styled bar]
```

### The `styled` Builtin

```elvish
styled $object $style-transformer...
```

Constructs styled text by applying transformers. Style transformers use **kebab-case** names:

| Category | Examples |
|----------|---------|
| Boolean attributes | `bold`, `dim`, `italic`, `underlined`, `blink`, `inverse` |
| Negation | `no-bold`, `no-italic`, `no-underlined` |
| Toggle | `toggle-bold`, `toggle-inverse` |
| ANSI colors | `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white` |
| Bright colors | `bright-black`, `bright-red`, etc. |
| XTerm 256-color | `color0` through `color255` |
| TrueColor RGB | `#RRGGBB` (must be quoted — `#` starts comments) |
| Foreground prefix | `fg-red`, `fg-#ff0000`, `fg-color123` |
| Background prefix | `bg-blue`, `bg-default`, `bg-color12` |

Multiple transformers are applied left-to-right as space-separated arguments:

```elvish
styled $value red bold          # red bold text
styled $value fg-red bg-blue    # red foreground, blue background
styled $value 'bg-#778899'      # background using hex RGB
```

### The `edit:notify` Command

```elvish
edit:notify $message
```

Displays a notification above the editor. Accepts strings or styled text. Used by carapace for error messages and usage hints.

### The `from-json` Builtin

```elvish
from-json
```

Reads JSON from stdin, parses it, and outputs structured elvish values on stdout. Carapace pipes its JSON output through `from-json` to produce elvish map/list values that the snippet code can index and iterate over.

```elvish
echo '{"key": "value"}' | from-json  # ▶ [&key=value]
echo '["a","b"]' | from-json         # ▶ [a b]
```

### Elvish's Matcher System

After the completer outputs candidates, elvish matches them against the user's typed text (the **seed**). The matcher table `$edit:completion:matcher` maps completion types to matcher functions:

```elvish
# Default: prefix matching for all types
set edit:completion:matcher[''] = $edit:match-prefix~

# Substring matching for arguments
set edit:completion:matcher[argument] = $edit:match-substr~
```

Built-in matchers: `edit:match-prefix`, `edit:match-substr`, `edit:match-subseq`. All accept `&ignore-case` and `&smart-case` options.

### Deprecated Options

- **`&display-suffix`** — deprecated in elvish 0.14.0. Use `&display` instead. Carapace never used this.
- **`&style`** — removed from `edit:complex-candidate`. Styling is now done via `styled` text in `&display`. Fixed in elvish 0.18.0 ([issue #1011](https://github.com/elves/elvish/issues/1011)).

## The Elvish Snippet

Generated by `elvish.Snippet(cmd)` in `internal/shell/elvish/snippet.go`:

```elvish
set edit:completion:arg-completer[example] = {|@arg|
    example _carapace elvish (all $arg) | from-json | each {|completion|
        put $completion[Messages] | all (one) | each {|m|
            edit:notify (styled "error: " red)$m
        }
        if (not-eq $completion[Usage] "") {
            edit:notify (styled "usage: " $completion[DescriptionStyle])$completion[Usage]
        }
        put $completion[Candidates] | all (one) | peach {|c|
            if (eq $c[Description] "") {
                edit:complex-candidate $c[Value] &display=(styled $c[Display] $c[Style]) &code-suffix=$c[CodeSuffix]
            } else {
                edit:complex-candidate $c[Value] &display=(styled $c[Display] $c[Style])(styled " " $completion[DescriptionStyle]" bg-default")(styled "("$c[Description]")" $completion[DescriptionStyle]) &code-suffix=$c[CodeSuffix]
            }
        }
    }
}
```

On Windows, an additional registration is appended:

```elvish
set edit:completion:arg-completer[example.exe] = $edit:completion:arg-completer[example]
```

### Snippet Walkthrough

**1. Register the completer**

```elvish
set edit:completion:arg-completer[example] = {|@arg|
```

The `set` command assigns the completion function to the `arg-completer` map. The `@arg` rest parameter receives all command-line words.

**2. Invoke carapace and parse JSON**

```elvish
example _carapace elvish (all $arg) | from-json | each {|completion|
```

- `example` — the CLI being completed
- `_carapace` — carapace's hidden completion subcommand
- `elvish` — the shell name, selecting the elvish formatter
- `(all $arg)` — spreads the arg list into individual arguments (equivalent to `...$spans` in nushell)
- `| from-json` — pipes the JSON output through elvish's JSON parser
- `| each {|completion|` — iterates over each completion object (there is typically one)

**3. Display error messages**

```elvish
put $completion[Messages] | all (one) | each {|m|
    edit:notify (styled "error: " red)$m
}
```

- `put $completion[Messages]` — outputs the Messages array
- `| all (one)` — handles the case where the array has one element (elvish `put` unwraps single-element lists, `all (one)` re-wraps it for iteration)
- `| each {|m|` — iterates over each message
- `edit:notify` — displays the message above the editor
- `(styled "error: " red)` — prefixes with red "error: " text
- `$m` — the message string (concatenated directly via juxtaposition)

**4. Display usage hint**

```elvish
if (not-eq $completion[Usage] "") {
    edit:notify (styled "usage: " $completion[DescriptionStyle])$completion[Usage]
}
```

Only shown when there are no candidates (carapace suppresses usage when values exist — see Usage Suppression below). The `DescriptionStyle` from the completion object is used for the "usage: " prefix styling.

**5. Emit completion candidates**

```elvish
put $completion[Candidates] | all (one) | peach {|c|
```

- `peach` — parallel `each` (candidates are emitted concurrently for performance)
- `all (one)` — same single-element unwrapping fix as messages

**6. Build `edit:complex-candidate` per candidate**

```elvish
if (eq $c[Description] "") {
    edit:complex-candidate $c[Value] &display=(styled $c[Display] $c[Style]) &code-suffix=$c[CodeSuffix]
} else {
    edit:complex-candidate $c[Value] &display=(styled $c[Display] $c[Style])(styled " " $completion[DescriptionStyle]" bg-default")(styled "("$c[Description]")" $completion[DescriptionStyle]) &code-suffix=$c[CodeSuffix]
}
```

Two cases based on whether a description exists:

**No description**: Simple display — just the styled value:
```
&display=(styled $c[Display] $c[Style])
```

**With description**: Three concatenated styled segments:
1. `(styled $c[Display] $c[Style])` — the value in its style
2. `(styled " " $completion[DescriptionStyle]" bg-default")` — a space separator with description style on default background (resets any value background color)
3. `(styled "("$c[Description]")" $completion[DescriptionStyle])` — the description in parentheses with description style

The `&code-suffix=$c[CodeSuffix]` field controls trailing space:
- `" "` (space) — normal insertion, add trailing space
- `""` (empty) — nospace, no trailing space (e.g., for directories, `--flag=`)

**7. `.exe` registration on Windows**

```elvish
set edit:completion:arg-completer[example.exe] = $edit:completion:arg-completer[example]
```

On Windows, commands may be invoked with the `.exe` suffix. The snippet registers the same completer for both names.


## No Elvish-Side Patching

Unlike bash (redirect handling via `bash.Patch()`) and nushell (open-quote stripping via `nushell.Patch()`), **elvish has no `Patch()` function**. The `complete.go` dispatch has no `case "elvish"` branch for argument patching — elvish arguments are passed directly to `traverse()` unmodified.

This is because elvish:
- Passes arguments as clean, already-parsed strings via `(all $arg)` — no raw shell input to patch
- Handles quoting natively in its parser — the completer receives properly parsed arguments
- Does not leak redirect operators into arguments — elvish's parser separates redirections before calling completers
- Has no `COMP_WORDBREAKS` equivalent — elvish doesn't split on `:` or `=`

The open-quote problem (present in bash/zsh/fish) also does not apply to elvish because elvish's `(all $arg)` receives already-tokenized arguments from its parser, not raw shell input.

## Value Formatting: `ActionRawValues()`

`elvish.ActionRawValues(currentWord, meta, values)` in `internal/shell/elvish/action.go` formats completion candidates as a JSON object containing metadata and an array of `complexCandidate` records. It is one of the richer shell formatters in carapace — handling nospace via `CodeSuffix`, style validation via elvish's `ParseStyling`, and message/usage integration.

### Output Format

```json
{
  "Usage": "some usage text",
  "Messages": ["error message"],
  "DescriptionStyle": "dim",
  "Candidates": [
    {"Value": "checkout", "Display": "checkout", "Description": "Switch branches", "CodeSuffix": " ", "Style": "blue"},
    {"Value": "--flag=", "Display": "--flag", "Description": "Set flag", "CodeSuffix": "", "Style": "default"}
  ]
}
```

Top-level fields:

| Field | Type | Purpose |
|-------|------|---------|
| `Usage` | `string` | Usage hint shown via `edit:notify` (empty when candidates exist) |
| `Messages` | `[]string` | Error messages shown via `edit:notify` |
| `DescriptionStyle` | `string` | Elvish style string for description text (e.g., `"dim"`, `"fg-cyan"`) |
| `Candidates` | `[]complexCandidate` | Array of completion candidates |

Each `complexCandidate` record:

| Field | Type | Purpose |
|-------|------|---------|
| `Value` | `string` | Text inserted into the command line (the `$stem` of `edit:complex-candidate`) |
| `Display` | `string` | Text shown in the completion menu (passed to `styled` builtin) |
| `Description` | `string` | Help text shown in parentheses after the display |
| `CodeSuffix` | `string` | `" "` for normal trailing space, `""` for nospace |
| `Style` | `string` | Elvish style string for the display text (e.g., `"blue"`, `"bold red"`, `"fg-#ff0000"`) |

### The `completion` and `complexCandidate` Structs

```go
type completion struct {
    Usage            string
    Messages         common.Messages
    DescriptionStyle string
    Candidates       []complexCandidate
}

type complexCandidate struct {
    Value       string
    Display     string
    Description string
    CodeSuffix  string
    Style       string
}
```

### Step 1: Determine Default Styles

```go
valueStyle := "default"
if s := style.Carapace.Value; s != "" && ui.ParseStyling(s) != nil {
    valueStyle = s
}

descriptionStyle := "default"
if s := style.Carapace.Description; s != "" && ui.ParseStyling(s) != nil {
    descriptionStyle = s
}
```

Two default styles are resolved using carapace's `style.Carapace` config and elvish's `ui.ParseStyling`:

- **`valueStyle`** — fallback for candidates whose `Style` field is empty or invalid. Defaults to `"default"`, overridden by `style.Carapace.Value` (default: `style.Default = ""` which resolves to `"default"`).
- **`descriptionStyle`** — style for the description text and "usage: " prefix in the snippet. Defaults to `"default"`, overridden by `style.Carapace.Description` (default: `style.Dim`).

`ui.ParseStyling` is from the **vendored elvish `pkg/ui` package** at `third_party/github.com/elves/elvish/pkg/ui/`. It parses elvish-style kebab-case style strings (e.g., `"bold red"`, `"fg-#ff0000"`) and returns `nil` for invalid strings. This validation ensures that only valid elvish style strings reach the `styled` builtin in the snippet — invalid strings would cause runtime errors in elvish.

### Step 2: Sanitize Values

```go
var sanitizer = strings.NewReplacer(
    "\n", ``,
    "\r", ``,
    "\t", ``,
)

func sanitize(values []common.RawValue) []common.RawValue {
    for index, v := range values {
        (&values[index]).Value = sanitizer.Replace(v.Value)
        (&values[index]).Display = sanitizer.Replace(v.Display)
        (&values[index]).Description = sanitizer.Replace(v.TrimmedDescription())
    }
    return values
}
```

The sanitizer strips control characters that would break the JSON output or elvish's display:

- `\n` → empty — newlines in JSON string values are fine (JSON encodes them as `\\n`), but newlines in completion display text break the UI
- `\r` → empty — carriage returns would corrupt display
- `\t` → empty — tabs in display/description text would create visual misalignment

Unlike the fish formatter (which also strips tabs because they are field delimiters), elvish strips tabs purely for display cleanliness — JSON handles tabs fine, but the completion UI would show them as whitespace gaps.

Note that `Description` uses `TrimmedDescription()` (truncated to `CARAPACE_DESCRIPTION_LENGTH`, default 80 characters) rather than the full `Description`. This prevents overly long descriptions from cluttering the completion menu.

### Step 3: Build Candidates with CodeSuffix (Nospace)

```go
vals := make([]complexCandidate, len(values))
for index, val := range sanitize(values) {
    suffix := " "
    if meta.Nospace.Matches(val.Value) {
        suffix = ""
    }

    if val.Style == "" || ui.ParseStyling(val.Style) == nil {
        val.Style = valueStyle
    }
    vals[index] = complexCandidate{
        Value:       val.Value,
        Display:     val.Display,
        Description: val.Description,
        CodeSuffix:  suffix,
        Style:       val.Style,
    }
}
```

For each candidate:

1. **Determine `CodeSuffix`** based on `meta.Nospace`:
   - If `Nospace.Matches(val.Value)` returns `true` (the value ends with a nospace character like `/`, `:`, `=`, `@`, `,`, or `*` is in the set): `suffix = ""` — no trailing space
   - Otherwise: `suffix = " "` — add trailing space after insertion

   This maps directly to `edit:complex-candidate`'s `&code-suffix` option. A space `" "` means elvish adds a trailing space after the candidate is accepted; an empty string `""` means no trailing space.

2. **Validate and default `Style`**:
   - If `val.Style` is empty or fails `ui.ParseStyling` validation, fall back to `valueStyle` (resolved in Step 1)
   - This prevents runtime errors in elvish — the `styled` builtin would throw an exception for invalid style strings

3. **Build the `complexCandidate` struct** with all five fields

Elvish is one of three carapace shells with **per-candidate nospace support** (alongside zsh and nushell). The `CodeSuffix` approach is elvish-specific — other shells handle nospace differently:

| Shell | Nospace Mechanism |
|-------|-------------------|
| elvish | `CodeSuffix` in `complexCandidate` (`""` vs `" "`) |
| zsh | Per-candidate space suffix in `_describe` values |
| nushell | Trailing space baked into JSON `value` field |
| bash | Global `compopt -o nospace` (all-or-nothing) |
| fish | Not supported |

### Step 4: Usage Suppression

```go
if len(values) > 0 {
    meta.Usage = "" // TODO edit:notify is persistent, so avoid spamming the user for now
}
```

When there are completion candidates, the `Usage` field is cleared. This is a deliberate workaround for a known elvish behavior: `edit:notify` messages are **persistent** — they remain visible until explicitly dismissed or overwritten. If usage hints were shown on every completion, they would clutter the editor.

The TODO comment indicates this is a temporary measure. The ideal behavior would be to show usage only when there are no candidates, but since `edit:notify` doesn't auto-dismiss, showing it at all is considered too noisy when candidates are also present.

When `Usage` is empty, the snippet's `if (not-eq $completion[Usage] "")` check prevents the `edit:notify` call entirely.

### Step 5: Marshal JSON

```go
m, _ := json.Marshal(completion{
    Usage:            meta.Usage,
    Messages:         meta.Messages,
    DescriptionStyle: descriptionStyle,
    Candidates:       vals,
})
return string(m)
```

The entire `completion` struct is marshaled to JSON. The `Messages` type has a custom `MarshalJSON` that serializes the internal `map[string]bool` as a sorted `[]string`.

The JSON output is piped through `from-json` in the snippet, producing elvish map/list values that the snippet code indexes with `$completion[field]` and `$c[field]` notation.

## Style Validation: `ui.ParseStyling`

Carapace vendors a subset of elvish's `pkg/ui` package at `third_party/github.com/elves/elvish/pkg/ui/`. This package provides `ParseStyling`, which is the **same function** elvish uses internally to parse style strings.

### Why Vendored?

The `ParseStyling` function is used in `ActionRawValues()` to validate style strings before they reach the elvish snippet. Without validation, an invalid style string (e.g., `"not-a-real-color"`) passed to the `styled` builtin would cause a runtime error in elvish, crashing the completion.

### How `ParseStyling` Works

```go
func ParseStyling(s string) Styling {
    if !strings.ContainsRune(s, ' ') {
        return parseOneStyling(s)
    }
    var joint jointStyling
    for subs := range strings.SplitSeq(s, " ") {
        parsed := parseOneStyling(subs)
        if parsed == nil {
            return nil
        }
        joint = append(joint, parseOneStyling(subs))
    }
    return joint
}
```

It splits the style string on spaces and tries to parse each token:

| Token pattern | Resolves to |
|---------------|-------------|
| `"default"` or `"fg-default"` | `FgDefault` |
| `"fg-<color>"` | `setForeground{parseColor(color)}` |
| `"bg-default"` | `BgDefault` |
| `"bg-<color>"` | `setBackground{parseColor(color)}` |
| `"no-<bool>"` | `boolOff` (e.g., `no-bold`) |
| `"toggle-<bool>"` | `boolToggle` |
| `"<bool>"` | `boolOn` (e.g., `bold`, `italic`) |
| `"<color>"` | `setForeground{parseColor(color)}` (no prefix = foreground) |

If **any** token fails to parse, `ParseStyling` returns `nil`, and carapace falls back to `valueStyle`.

### Color Parsing

```go
func parseColor(name string) Color {
    if color, ok := colorByName[name]; ok {
        return color  // "red", "blue", "bright-green", etc.
    }
    if strings.HasPrefix(name, "color") {
        // XTerm 256-color: "color0" through "color255"
        i, _ := strconv.Atoi(name[5:])
        return XTerm256Color(uint8(i))
    }
    if strings.HasPrefix(name, "#") && len(name) == 7 {
        // TrueColor RGB: "#ff0000"
        r, _ := strconv.ParseUint(name[1:3], 16, 8)
        g, _ := strconv.ParseUint(name[3:5], 16, 8)
        b, _ := strconv.ParseUint(name[5:7], 16, 8)
        return TrueColor(uint8(r), uint8(g), uint8(b))
    }
    return nil  // invalid color
}
```

Supported color formats:
- **Named ANSI colors**: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`
- **Bright variants**: `bright-black` through `bright-white`
- **XTerm 256-color**: `color0` through `color255`
- **TrueColor RGB**: `#RRGGBB` (6-digit hex, must be quoted in elvish since `#` starts comments)

### Style String Examples

These carapace style strings are valid elvish style strings:

| Carapace style | Elvish interpretation |
|----------------|----------------------|
| `"red"` | Foreground red |
| `"bold red"` | Bold + foreground red |
| `"fg-blue bg-yellow"` | Blue foreground on yellow background |
| `"italic underlined"` | Italic + underlined |
| `"color123"` | XTerm 256-color index 123 |
| `"#ff5500"` | TrueColor RGB (must be quoted in snippet) |
| `"dim"` | Dim/faint text |
| `"inverse"` | Reverse video (swap fg/bg) |

## Message Handling: Native `edit:notify`

Elvish is one of three shells (alongside zsh and the export format) that has **native message support** in carapace's dispatch:

```go
// shell.go
switch shell {
case "elvish", "export", "zsh": // shells with support for showing messages
default:
    values = meta.Messages.Integrate(values, value)
}
```

For elvish, messages are **not** integrated as synthetic `ERR` values. Instead, they are passed through the `Messages` field in the JSON output and displayed via `edit:notify` in the snippet.

This is superior to the fallback approach (used by bash/fish/oil/powershell/tcsh/ion) because:

1. **Proper formatting** — messages appear as distinct notifications above the editor, not mixed into the completion list
2. **Styled text** — the `styled` builtin allows error styling: `(styled "error: " red)$m`
3. **No value collision** — messages don't need fake `ERR` values that could conflict with real completions
4. **Persistent visibility** — `edit:notify` messages remain visible until dismissed, unlike transient completion candidates

## Usage Suppression

Carapace deliberately suppresses the `Usage` field when completion candidates exist:

```go
if len(values) > 0 {
    meta.Usage = "" // TODO edit:notify is persistent, so avoid spamming the user for now
}
```

The comment explains the rationale: `edit:notify` messages are **persistent** — they remain visible until dismissed. If usage hints were shown alongside every completion, the user would be spammed with notifications on every TAB press. By clearing `Usage` when there are candidates, the snippet only shows usage when the completion result is empty (e.g., wrong subcommand or flag), where the hint is actually useful.

This is different from zsh's `_message -r` which is less intrusive (shown in the completion area only during the current completion session).

## Edge Cases and Known Issues

### 1. The `all (one)` Pattern for Single-Element Lists

Elvish has a design choice where `put` with a single-element list "unwraps" the element. This causes problems when iterating:

```elvish
put $completion[Messages] | all (one) | each {|m| ... }
put $completion[Candidates] | all (one) | peach {|c| ... }
```

Without `all (one)`, if `Messages` or `Candidates` contains exactly one element, elvish would unwrap it from a list into a scalar, breaking the `| each` pipeline. The `all (one)` command ensures the value is always a list, even with a single element.

This is a common elvish pattern when consuming JSON arrays via `from-json`.

### 2. Style Validation via Vendored `ui.ParseStyling`

Carapace validates every style string using elvish's own `ui.ParseStyling` function, vendored from the elvish source at `third_party/github.com/elves/elvish/pkg/ui/`. If a style string fails validation, it falls back to `valueStyle`:

```go
if val.Style == "" || ui.ParseStyling(val.Style) == nil {
    val.Style = valueStyle
}
```

This is critical because elvish's `styled` builtin will throw a runtime error if given an invalid style string. Since style strings come from multiple sources (carapace's style system, user configs, YAML specs), validation prevents crashes.

The vendored package ensures that carapace's style validation always matches the exact version of elvish's parser — no version skew.

### 3. Description Style Reset with `bg-default`

In the snippet's description rendering:

```elvish
(styled " " $completion[DescriptionStyle]" bg-default")
```

The `" bg-default"` suffix explicitly resets the background color to default. This is necessary because the value's style may set a background color, and without the reset, the space between the value and the description would inherit that background color, creating a visual artifact.

### 4. `peach` (Parallel Each) for Candidate Emission

The snippet uses `peach` instead of `each` to emit candidates:

```elvish
put $completion[Candidates] | all (one) | peach {|c|
```

`peach` processes elements concurrently, which can improve performance for large candidate sets. However, this means the order of candidate emission is not guaranteed — carapace handles sorting before output, so the order is correct in the data; `peach` just parallelizes the `edit:complex-candidate` calls.

### 5. No Per-Candidate Nospace in Display — Only in CodeSuffix

Elvish's `edit:complex-candidate` controls trailing space via `&code-suffix`, not via a display property. This means:

- The **display** in the completion menu always looks the same regardless of nospace
- The **insertion behavior** differs: `&code-suffix=" "` adds a trailing space, `&code-suffix=""` does not
- This is the correct elvish idiom — directories get `&code-suffix=""` (no trailing space, user continues typing the path), files get `&code-suffix=" "` (trailing space, ready for next argument)

### 6. Values Containing Special Characters

Elvish automatically quotes the `$stem` argument of `edit:complex-candidate` when inserting it into the command line. Carapace does **not** perform any shell-specific escaping of values in the elvish formatter — it relies on elvish's built-in quoting.

This means:
- Values like `hello world` (with spaces) are correctly quoted by elvish upon insertion
- Values like `--flag=value` are inserted as-is
- Values like `$HOME` are inserted literally (elvish quotes prevent expansion)

This is simpler than bash's two-mode quoting or zsh's 5-state machine, because elvish handles quoting at the insertion layer, not at the output formatting layer.

### 7. JSON Output and `from-json` Robustness

The `from-json` builtin in elvish is tolerant of extra whitespace between JSON values and handles multiple JSON objects. Carapace's output is always a single JSON object, so there's no ambiguity.

However, if the carapace subprocess fails or produces invalid JSON, `from-json` will throw an error. The snippet does not have error handling for this case — a failed `_carapace` invocation simply produces no completions.

### 8. No Tag Grouping

Unlike zsh (which uses `_describe -t` for tag groups), elvish has no tag grouping mechanism. Carapace's `EachTag()` iteration is not used in the elvish formatter — all candidates are in a single flat list. Tags from `RawValue.Tag` are used for deduplication in `shell.go`'s `mergeFlags()` logic but are not surfaced in the elvish UI.

### 9. The `DescriptionStyle` as a Shared Field

The `completion.DescriptionStyle` field is set once from `style.Carapace.Description` and used for both:
- The "usage: " prefix in `edit:notify`
- The description text in the completion display `(styled "("$c[Description]")" $completion[DescriptionStyle])`

This means changing `style.Carapace.Description` in carapace's config affects both usage hints and completion descriptions simultaneously — they are always visually consistent.

### 10. `all $arg` vs `...$arg` — Spread Operator History

Older versions of the snippet may use `(all $arg)` while newer elvish versions support `...$arg` as a spread operator. Carapace uses `(all $arg)` which is compatible with a wider range of elvish versions. The `all` command outputs all values from a list, and `$arg` is the rest parameter that captures all arguments to the completer function.

### 11. Windows `.exe` Registration

On Windows, the snippet registers the same completer for both `example` and `example.exe`:

```go
if runtime.GOOS == "windows" {
    result += fmt.Sprintf("set edit:completion:arg-completer[%v.exe] = $edit:completion:arg-completer[%v]\n", cmd.Name(), cmd.Name())
}
```

This is needed because Windows users may invoke commands with or without the `.exe` suffix, and elvish on Windows treats them as different command names for completer lookup.

### 12. No Integration Tests

Unlike bash and zsh which have integration tests in `example/main_test.go`, the elvish integration tests are **commented out**:

```go
// func TestElvish(t *testing.T) {
//     if err := exec.Command("elvish", "--version").Run(); err != nil {
//         t.Skip("skipping elvish")
//     }
//     for cmdline, text := range tests {
//         doComplete(t, "elvish", cmdline, text)
//     }
// }
```

Elvish completions are only tested via the golden file snippet test (`example/cmd/_test/elvish.elv`) and the `testScript` function in `example/cmd/root_test.go`, which invokes `example _carapace elvish` with the test script.

There are **no unit tests** in `internal/shell/elvish/` — no `_test.go` files exist.

### 13. Vendored Elvish Code Must Not Be Modified

The `third_party/github.com/elves/elvish/` directory is excluded from linting (per `.golangci.yml`) and must not be modified. It is vendored from the elvish project for:
- `ui.ParseStyling` — style string validation
- `ui.Style` — style structure definition
- `ui.Color` — color types (ANSI, XTerm256, TrueColor)
- `cli/lscolors` — `LS_COLORS` integration for file styling

If elvish changes its style string syntax, the vendored code must be updated, not patched.

## Completion Dispatch Flow for Elvish

```
User presses TAB
  → Elvish calls $edit:completion:arg-completer[cmd] with command-line args
  → Snippet: example _carapace elvish (all $arg)
  → carapace _carapace elvish subcommand receives args
  → complete() in complete.go detects shell
  → No patching needed (elvish not in switch cases)
  → traverse() classifies args (flag, positional, subcommand, dash)
  → storage.get(cmd) resolves Action
  → action.Invoke(context) resolves callback chain
  → shell.Value("elvish", value, meta, values) dispatches to elvish.ActionRawValues
  → shell.go: messages NOT integrated (native support)
  → shell.go: nospace.Add('*') when messages exist
  → elvish.ActionRawValues formats JSON with completion struct
  → JSON piped through from-json in elvish
  → Snippet: edit:notify for messages, edit:complex-candidate for each candidate
  → Elvish displays candidates in completion UI
```

## Related Skills

- **references/shell.md** — cross-shell comparison table and per-shell feature matrix
- **references/shell-bash.md** — bash integration deep dive
- **references/shell-fish.md** — fish integration deep dive
- **references/shell-zsh.md** — zsh integration deep dive
- **references/shell-nushell.md** — nushell integration deep dive
- **references/shell-xonsh.md** — xonsh integration deep dive
- **references/traverse.md** — the completion engine that produces Actions before formatting
- **references/style.md** — how styles are resolved before shell rendering