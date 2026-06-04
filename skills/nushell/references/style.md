# Nushell Style and Color System

In-depth reference for nushell's color and styling system — how colors work in completions, menus, syntax highlighting, and how carapace converts its style format to nushell's `{fg, bg, attr}` records.

## Color Formats

Nushell supports multiple color formats:

| Format | Example | Description |
|--------|---------|-------------|
| Abbreviation | `r` | Single letter for color |
| Abbreviation + attribute | `rb` | Color abbreviation with style attribute |
| Full name | `red` | Complete color name |
| Full name + attribute | `red_bold` | Color name with underscore and attribute |
| Hex format | `"#ff0000"` | 24-bit truecolor (quotes required) |
| Full hex record | `{ fg: "#ff0000" bg: "#0000ff" attr: b }` | Full record with fg, bg, and attr |
| Closure | `{\|x\| 'yellow' }` | Closure returning a color string (table output only) |

## Attribute Codes

| Code | Meaning |
|------|---------|
| `l` | Blink |
| `b` | Bold |
| `d` | Dimmed |
| `h` | Hidden |
| `i` | Italic |
| `r` | Reverse (swap fg/bg) |
| `s` | Strikethrough |
| `u` | Underline |
| `n` | Nothing (default) |

Multiple attributes are concatenated: e.g., bold+underline → `"bu"`.

## Named Colors

### Standard Colors

| Code | Name | Code | Name |
|------|------|------|------|
| `g` | green | `gb` | green_bold |
| `r` | red | `rb` | red_bold |
| `u` | blue | `ub` | blue_bold |
| `b` | black | `bb` | black_bold |
| `ligr` | light_gray | `ligrb` | light_gray_bold |
| `y` | yellow | `yb` | yellow_bold |
| `p` | purple | `pb` | purple_bold |
| `m` | magenta | `mb` | magenta_bold |
| `c` | cyan | `cb` | cyan_bold |
| `w` | white | `wb` | white_bold |
| `dgr` | dark_gray | `dgrb` | dark_gray_bold |
| `def` | default | `defb` | default_bold |

### Light Variants

| Code | Name | Code | Name |
|------|------|------|------|
| `lg` | light_green | `lgb` | light_green_bold |
| `lr` | light_red | `lrb` | light_red_bold |
| `lu` | light_blue | `lub` | light_blue_bold |
| `ly` | light_yellow | `lyb` | light_yellow_bold |
| `lp` | light_purple | `lpb` | light_purple_bold |
| `lm` | light_magenta | `lmb` | light_magenta_bold |
| `lc` | light_cyan | `lcb` | light_cyan_bold |

### Background Colors

Background colors use `bg_` prefix: `bg_g`, `bg_lg`, `bg_r`, `bg_lr`, `bg_u`, `bg_lu`, `bg_b`, `bg_ligr`, `bg_y`, `bg_ly`, `bg_p`, `bg_lp`, `bg_m`, `bg_lm`, `bg_c`, `bg_lc`, `bg_w`, `bg_dgr`, `bg_def`.

## Style in Completions

### Completion Record Style Field

The `style` field in completion records accepts:

1. **String with foreground color** (hex code or color name):
   ```nu
   { value: "main", style: red }
   { value: "main", style: "#ff0000" }
   ```

2. **Record with `fg`, `bg`, and `attr` fields**:
   ```nu
   { value: "main", style: { fg: green, bg: "#66078c", attr: ub } }
   ```

3. **JSON string** version of the record (for external completers):
   ```json
   {"value": "main", "style": {"fg": "green", "bg": "#66078c", "attr": "ub"}}
   ```

### Carapace Style Conversion

Carapace converts its internal style format to nushell's `{fg, bg, attr}` record format:

#### Attribute Mapping

| Carapace Attribute | Nushell Flag | Meaning |
|-------------------|--------------|---------|
| `Blink` | `l` | Blinking text |
| `Bold` | `b` | Bold text |
| `Dim` | `d` | Dim/faint text |
| `Italic` | `i` | Italic text |
| `Inverse` | `r` | Reverse video (swap fg/bg) |
| `Underlined` | `u` | Underlined text |

Multiple attributes are concatenated: e.g., bold+underline → `"bu"`.

#### Color Name Mapping

| Carapace Name | Nushell Name |
|---------------|-------------|
| `black` | `black` |
| `red` | `red` |
| `green` | `green` |
| `yellow` | `yellow` |
| `blue` | `blue` |
| `magenta` | `magenta` |
| `cyan` | `cyan` |
| `white` | `white` |
| `bright-black` | `dark_gray` |
| `bright-red` | `light_red` |
| `bright-green` | `light_green` |
| `bright-yellow` | `light_yellow` |
| `bright-blue` | `light_blue` |
| `bright-magenta` | `light_magenta` |
| `bright-cyan` | `light_cyan` |
| `bright-white` | `white` |

#### XTerm256 Colors

XTerm256 colors (`color0` through `color255`) are converted to hex RGB strings. For example, `color196` → `#ff0000`.

#### Hex Colors

Hex colors (`#RRGGBB`) are passed through unchanged.

#### JSON Output Examples

```json
[
    {"value": "branch ", "display": "branch", "description": "Switch branches", "style": {"fg": "green"}},
    {"value": "checkout ", "display": "checkout", "description": "Switch working tree"},
    {"value": "main", "display": "main", "description": "Default branch", "style": {"fg": "red"}}
]
```

When style is empty, the `style` field is omitted from JSON output (via `omitempty`).

## Menu Styles

Reedline menus use separate style configuration from `color_config`:

```nu
$env.config.menus ++= [{
    name: completion_menu
    style: {
        text: green                    # Unselected suggestion
        selected_text: green_reverse   # Selected suggestion
        description_text: yellow       # Description text
    }
}]
```

### MenuTextStyle

| Style | Default | Purpose |
|-------|---------|---------|
| `text` | DarkGray.normal() | Unselected suggestions |
| `selected_text` | Green.bold().reverse() | Selected suggestion |
| `description_text` | Yellow.normal() | Description text |
| `selected_match_style` | — | Matched text in selected item (fuzzy highlighting) |
| `match_style` | — | Matched text in unselected items (fuzzy highlighting) |

## Syntax Coloring (shape_*)

Nushell uses `shape_*` settings in `$env.config.color_config` for syntax highlighting:

```nu
$env.config.color_config.shape_garbage = { fg: "#FFFFFF" bg: "#FF0000" attr: b }
$env.config.color_config.shape_int = { fg: "#0000ff" attr: b }
$env.config.color_config.shape_string = green
$env.config.color_config.shape_filepath = cyan
$env.config.color_config.shape_flag = blue_bold
$env.config.color_config.shape_internalcall = cyan_bold
$env.config.color_config.shape_external = cyan
```

See [references/types.md](references/types.md) for the full list of `shape_*` settings.

## Color Config (Primitive Values)

```nu
$env.config.color_config.separator = purple
$env.config.color_config.header = gb
$env.config.color_config.bool = red
$env.config.color_config.int = green
$env.config.color_config.hints = dark_gray
```

| Primitive | Default |
|-----------|---------|
| `binary` | White.normal() |
| `bool` | White.normal() |
| `cell-path` | White.normal() |
| `datetime` | White.normal() |
| `duration` | White.normal() |
| `filesize` | White.normal() |
| `float` | White.normal() |
| `int` | White.normal() |
| `list` | White.normal() |
| `nothing` | White.normal() |
| `range` | White.normal() |
| `record` | White.normal() |
| `string` | White.normal() |
| `leading_trailing_space_bg` | Rgb(128,128,128) |
| `header` | Green.bold() |
| `empty` | Blue.normal() |
| `row_index` | Green.bold() |
| `hints` | DarkGray.normal() |

### Closures for Dynamic Colors

Closures are only executed for table output and do not work for `shape_*` configurations:

```nu
$env.config.color_config.filesize = {|x|
    if $x == 0b { 'dark_gray' }
    else if $x < 1mb { 'cyan' }
    else { 'blue' }
}
```

## LS_COLORS

For the `ls` command, nushell uses the standard `LS_COLORS` environment variable:

```nu
$env.LS_COLORS = "di=1;34:*.nu=3;33;46"
```

Can use [vivid](https://github.com/sharkdp/vivid) for themes:

```nu
$env.LS_COLORS = (vivid generate molokai)
```

## Built-in Themes

```nu
use std/config light-theme
$env.config.color_config: (light-theme)
```

## References

- [Nushell Book: Coloring and Theming](https://www.nushell.sh/book/coloring_and_theming.html)
- [Carapace Nushell Style Conversion](https://github.com/carapace-sh/carapace/blob/main/internal/shell/nushell/style.go)

## Related Skills

- [references/completion.md](references/completion.md) — How styles are used in completion records
- [references/reedline.md](references/reedline.md) — Menu style configuration
- [references/types.md](references/types.md) — Shape_* syntax coloring
