# Elvish Styling System

In-depth reference for elvish's text styling system — the `styled` and `styled-segment` builtins, style transformers, color models, `ui.Text` and `ui.Segment` types, `ParseStyling` validation, Styledown, and how styling integrates with the completion UI and prompts.

## The `styled` Builtin

```elvish
styled $object $style-transformer...
```

Constructs styled text by applying one or more style transformers. The result is a `ui.Text` value — a list of `ui.Segment` objects, each carrying text content and style information.

### Style Transformers

Transformers use **kebab-case** names and are applied left-to-right as separate arguments:

```elvish
styled $value red bold          # red bold text
styled $value fg-red bg-blue    # red foreground, blue background
styled $value 'bg-#778899'      # background using hex RGB (must be quoted — # starts comments)
```

### Boolean Attribute Transformers

| Transformer | Effect |
|-------------|--------|
| `bold` | Bold text |
| `dim` | Dim/faint text |
| `italic` | Italic text |
| `underlined` | Underlined text |
| `blink` | Blinking text |
| `inverse` | Reverse video (swap foreground/background) |

### Negation Transformers

| Transformer | Effect |
|-------------|--------|
| `no-bold` | Remove bold |
| `no-dim` | Remove dim |
| `no-italic` | Remove italic |
| `no-underlined` | Remove underline |
| `no-blink` | Remove blink |
| `no-inverse` | Remove inverse |

### Toggle Transformers

| Transformer | Effect |
|-------------|--------|
| `toggle-bold` | Toggle bold |
| `toggle-dim` | Toggle dim |
| `toggle-italic` | Toggle italic |
| `toggle-underlined` | Toggle underline |
| `toggle-blink` | Toggle blink |
| `toggle-inverse` | Toggle inverse |

### ANSI Color Names

| Category | Names |
|----------|-------|
| Standard | `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white` |
| Bright | `bright-black`, `bright-red`, `bright-green`, `bright-yellow`, `bright-blue`, `bright-magenta`, `bright-cyan`, `bright-white` |

When used without a prefix, these set the **foreground** color.

### XTerm 256-Color

```
color0  color1  color2  ...  color255
```

256-color palette indices. When used without a prefix, sets the foreground color.

### TrueColor (24-bit RGB)

```
#RRGGBB
```

Must be **quoted** in elvish because `#` starts a comment:

```elvish
styled $value '#ff0000'       # red foreground
styled $value 'fg-#ff0000'    # explicit foreground
styled $value 'bg-#778899'   # background
```

### Foreground and Background Prefixes

| Prefix | Example | Effect |
|--------|---------|--------|
| `fg-` | `fg-red`, `fg-#ff0000`, `fg-color123` | Set foreground color |
| `bg-` | `bg-blue`, `bg-#778899`, `bg-color12` | Set background color |
| `bg-default` | `bg-default` | Reset background to default |

### Concatenation

Styled text values can be concatenated via juxtaposition (no space):

```elvish
styled checkout blue)(styled ' (Switch branches)' dim)
```

This produces a `ui.Text` with two segments: "checkout" in blue and " (Switch branches)" in dim.

### Styled Text in Completion Display

The `&display` parameter of `edit:complex-candidate` accepts `styled` output:

```elvish
edit:complex-candidate checkout &display=(styled checkout blue)(styled ' (Switch branches)' dim) &code-suffix=' '
```

### Styled Text in Prompts

```elvish
set edit:prompt = { styled (tilde-abbr $pwd) cyan; put '> ' }
set edit:rprompt = (constantly (styled (whoami)@(hostname) inverse))
```

### Styled Text in Notifications

```elvish
edit:notify (styled "error: " red)$message
```

## The `styled-segment` Builtin

```elvish
styled-segment $object &fg-color=default &bg-color=default &bold=$false &dim=$false &italic=$false &underlined=$false &blink=$false &inverse=$false
```

Creates a styled segment with explicit option-based styling instead of transformer strings. This is a lower-level API that gives precise control over each style attribute.

### Options

| Option | Type | Default | Description |
|-------|------|---------|-------------|
| `&fg-color` | string | `default` | Foreground color (name, `colorN`, `#RRGGBB`, or `default`) |
| `&bg-color` | string | `default` | Background color (name, `colorN`, `#RRGGBB`, or `default`) |
| `&bold` | boolean | `$false` | Bold text |
| `&dim` | boolean | `$false` | Dim/faint text |
| `&italic` | boolean | `$false` | Italic text |
| `&underlined` | boolean | `$false` | Underlined text |
| `&blink` | boolean | `$false` | Blinking text |
| `&inverse` | boolean | `$false` | Reverse video |

### Example

```elvish
styled-segment checkout &fg-color=blue &bold=$true
```

### When to Use `styled-segment` vs `styled`

- Use `styled` for most cases — it's more concise and readable
- Use `styled-segment` when you need to compute style options programmatically (e.g., from a map or variable)

## The `ui.Text` and `ui.Segment` Types

Internally, styled text is represented as `ui.Text` — a slice of `ui.Segment` structs.

### `ui.Segment`

```go
type Segment struct {
    Text      string
    Style     Style
}

type Style struct {
    Foreground Color
    Background Color
    Bold       bool
    Dim        bool
    Italic     bool
    Underlined bool
    Blink      bool
    Inverse    bool
}
```

### `ui.Text`

```go
type Text []Segment
```

`ui.Text` implements `vals.List` — it can be indexed, sliced, and iterated like an elvish list.

### `ui.Color`

```go
type Color struct {
    Type  ColorType  // Default, ANSI8Bit, TrueColor
    Name  string     // For ANSI colors: "red", "blue", etc.
    Color uint8      // For 256-color: 0-255
    R, G, B uint8    // For TrueColor: 0-255 each
}
```

## `ParseStyling` Validation

The `ui.ParseStyling` function (from `pkg/ui/`) parses elvish-style kebab-case style strings and returns a `ui.Style` value, or `nil` for invalid strings.

```go
func ParseStyling(s string) *Style
```

### Valid Style Strings

| Input | Result |
|-------|--------|
| `"bold"` | Bold style |
| `"red"` | Red foreground |
| `"bold red"` | Bold + red foreground |
| `"fg-red bg-blue"` | Red foreground, blue background |
| `"fg-#ff0000"` | TrueColor red foreground |
| `"color123"` | 256-color index 123 |
| `""` | Empty style (no attributes) |
| `"invalid"` | `nil` (invalid) |
| `"bold invalid"` | `nil` (entire string invalid if any part fails) |

### Carapace Integration

Carapace uses `ui.ParseStyling` to validate style strings before passing them to the elvish `styled` builtin in the snippet. Invalid style strings would cause runtime errors in elvish, so carapace falls back to a default style when `ParseStyling` returns `nil`:

```go
if val.Style == "" || ui.ParseStyling(val.Style) == nil {
    val.Style = valueStyle  // fallback to configured default
}
```

The vendored elvish `pkg/ui` package is at `third_party/github.com/elves/elvish/pkg/ui/` in the carapace source tree.

## Styledown

Styledown is elvish's system for serializing and deserializing styled text to/from a plain text markup format.

### `render-styledown`

```elvish
render-styledown $s
```

Converts a `ui.Text` value to Styledown markup string. Useful for debugging or logging styled text.

### `derender-styledown`

```elvish
derender-styledown $s &style-defs=''
```

Converts a Styledown markup string back to `ui.Text`. The `&style-defs` option allows defining custom style mappings.

## Color and Terminal Support

### NO_COLOR

Elvish respects the [NO_COLOR](https://no-color.org) environment variable (since 0.20.0). If set and non-empty, builtin UI elements and styled texts will not have colors.

```elvish
set E:NO_COLOR = 1  # Disable colors
```

### LSCOLORS

Elvish's file completion uses `$E:LSCOLORS` (on macOS) or `$E:LS_COLORS` (on Linux) for file type coloring, similar to `ls --color`.

### Terminal Capability

Elvish auto-detects terminal color support:
- **8-color** — basic ANSI colors
- **256-color** — XTerm 256-color palette
- **TrueColor** — 24-bit RGB (detected via `COLORTERM` env var or terminal type)

## Style String Reference (Quick Lookup)

| Category | Format | Examples |
|----------|--------|---------|
| Boolean attributes | `bold`, `dim`, `italic`, `underlined`, `blink`, `inverse` | `styled x bold` |
| Negation | `no-bold`, `no-italic`, etc. | `styled x no-bold` |
| Toggle | `toggle-bold`, `toggle-inverse`, etc. | `styled x toggle-bold` |
| ANSI foreground | Color name | `styled x red` |
| ANSI background | `bg-` prefix | `styled x bg-blue` |
| Bright foreground | `bright-` prefix | `styled x bright-red` |
| 256-color foreground | `colorN` | `styled x color123` |
| 256-color background | `bg-colorN` | `styled x bg-color45` |
| TrueColor foreground | `#RRGGBB` (quoted) | `styled x '#ff0000'` |
| TrueColor background | `bg-#RRGGBB` (quoted) | `styled x 'bg-#778899'` |
| Explicit foreground | `fg-` prefix | `styled x fg-red` |
| Reset background | `bg-default` | `styled x bg-default` |
| Multiple transforms | Space-separated | `styled x bold fg-red bg-blue` |

## Comparison with Other Shell Styling Systems

| Feature | Elvish | Bash | Zsh | Fish |
|---------|--------|------|-----|------|
| Styling API | `styled` builtin | ANSI escape codes | `%F{color}`, `$fg` | `set_color` |
| Style format | Kebab-case strings | `\e[1;31m` | `fg-red`, `bg-blue` | `--bold`, `--red` |
| TrueColor | `#RRGGBB` | `\e[38;2;R;G;Bm` | `%F{#RRGGBB}` | `--color=RRGGBB` |
| 256-color | `colorN` | `\e[38;5;Nm` | `%F{N}` | `--color=N` |
| Completion styling | `styled` in `&display` | Readline variables | `zstyle` colors | `--color` option |
| Prompt styling | `styled` in prompt fn | PS1 escape codes | `PROMPT` with `%F` | `fish_color_*` |
| Validation | `ParseStyling` | None | None | None |

## References

- [Elvish Builtin: styled](https://elv.sh/ref/builtin.html#styled) — official documentation
- [Elvish Builtin: styled-segment](https://elv.sh/ref/builtin.html#styled-segment) — official documentation
- [Elvish Source: pkg/ui/](https://github.com/elves/elvish/tree/master/pkg/ui) — UI types and styling
- [NO_COLOR](https://no-color.org/) — standard for disabling colors

## Related Skills

- For how styling integrates with completion (complex-candidate display), see [references/completion.md](completion.md).
- For how styling is used in prompts and notifications, see [references/editor.md](editor.md).
- For carapace's style validation via `ParseStyling`, see the **carapace-dev** skill → `references/shell-elvish.md`.
