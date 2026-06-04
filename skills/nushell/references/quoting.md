# Nushell Quoting and String Types

In-depth reference for nushell's string types, quoting rules, escape sequences, and how quoting affects the completion system — particularly how spans are passed to external completers.

## String Types

Nushell provides six string types:

| Type | Syntax | Escapes | Use in completions |
|------|--------|---------|-------------------|
| Single-quoted | `'text'` | None | Values with special chars |
| Double-quoted | `"text"` | C-style (`\"`, `\\`, etc.) | Values needing escape sequences |
| Raw string | `r#'text'#` | None | Values containing single quotes |
| Bare word | `word` | None | Simple alphanumeric values only |
| Backtick | `` `text` `` | None | Paths with spaces |
| Interpolated single | `$'...'` | None | Interpolation without escapes |
| Interpolated double | `$"..."` | C-style | Interpolation with escapes |

## Single-Quoted Strings (`'...'`)

- Simplest string type
- No escape processing — text is passed through unchanged
- Cannot contain single quotes within the string
- Can span multiple lines

```nu
'hello world'       # hello world
'$HOME \n \t'       # $HOME \n \t  (all literal)
```

## Double-Quoted Strings (`"..."`)

- Support C-style backslash escape characters
- Backslash retains special meaning only when followed by specific characters

### Escape Sequences

| Sequence | Meaning |
|----------|---------|
| `\"` | Double-quote character |
| `\'` | Single-quote character |
| `\\` | Backslash |
| `\/` | Forward slash |
| `\b` | Backspace |
| `\f` | Form feed |
| `\r` | Carriage return |
| `\n` | Newline (line feed) |
| `\t` | Tab |
| `\u{X...}` | Unicode character (1-6 hex digits) |
| `\0` | NUL character (via `\u{0}` or `(char nul)`) |

```nu
"hello\nworld"     # Two lines
"say \"hello\""    # say "hello"
"cost: \$5"        # cost: \$5  (\$ is NOT a special sequence in nushell)
```

**Important**: Unlike bash, `\$` is NOT a special escape in nushell double-quoted strings. `$` is only special in interpolated strings (`$"..."`).

## Raw Strings (`r#'...'#`)

- Behaves like single-quoted strings (no escapes)
- Can contain single quotes without escaping
- Additional `#` symbols allow nesting raw strings containing `'#`

```nu
r#'Raw strings can contain 'quoted' text.'#
# => Raw strings can contain 'quoted' text.

r###'r##'This is an example of a raw string.'##'###
# => r##'This is an example of a raw string.'##
```

## Bare Word Strings

- Unquoted strings consisting of "word" characters only
- Cannot be used in command position (first position)
- Cannot include spaces or special characters
- Many bare words have special meaning in nushell and won't be interpreted as strings

```nu
print hello        # hello (bare word as argument)
true | describe    # bool (not a string!)
[trueX] | describe # list<string> (not a keyword)
```

## Backtick-Quoted Strings (`` `...` ``)

- Alternative to bare word strings
- Can include whitespace
- Cannot contain unmatched backticks
- In first position, still interpreted as command or path

```nu
`ls`               # Run the external ls binary
`..`                # Move up one directory
`./my dir`          # Change to directory with spaces
ls `./my dir/*`    # Combine globs with spaces
```

## String Interpolation

### Single-Quoted Interpolation (`$'...'`)

- No escape support
- Expressions in `()` are evaluated
- Cannot contain `'` or unmatched `()`

```nu
let name = "Alice"
$'Hello, ($name)'   # Hello, Alice
```

### Double-Quoted Interpolation (`$"..."`)

- Supports C-style backslash escapes
- Expressions in `()` are evaluated
- Escape parentheses with `\(` and `)`

```nu
let name = "Alice"
$"greetings, ($name)"                    # greetings, Alice
$"2 + 2 is (2 + 2) \(you guessed it!)"   # 2 + 2 is 4 (you guessed it!)
```

### Parse-Time Evaluation

Interpolated strings are evaluated at parse time. Values with formatting dependent on configuration may use default config before `config.nu` loads:

```nu
const x = $"(2kb)"
# x will be "2.0 KB" regardless of config settings
```

## Metacharacters Requiring Quoting

Characters that require quoting in nushell:

```
(space) { } ( ) [ ] $ " ' ` < > & | ; # \
```

When a completion value contains any of these characters, it must be quoted for nushell to accept it.

## How Quoting Affects Completions

### Spans Passed to External Completers

Nushell's parser passes arguments to the completer closure **with quoting intact**. When a user is completing inside quotes, nushell includes the quote characters in the span:

| User types | Span received by completer | After quote stripping |
|-----------|---------------------------|----------------------|
| `example action` | `action` | `action` |
| `example "act` | `"act` or `'"act'` | `act` |
| `example 'act` | `'act` or `"'act"` | `act` |
| `example `` `act` `` | `` `act` `` | `act` |

This is fundamentally different from bash, which strips quotes before passing words to completion functions.

### Carapace's Patch Phase

Carapace handles quote stripping in `nushell.Patch()`:

```go
func Patch(args []string) []string {
    for index, arg := range args {
        if len(arg) == 0 { continue }
        switch arg[0] {
        case '"', "'"[0]:
            if tokens, err := shlex.Split(arg); err == nil {
                args[index] = tokens[0].Value
            }
        case '`':
            args[index] = strings.Trim(arg, "`")
        }
    }
    return args
}
```

- **Double-quoted and single-quoted args** — parsed via `shlex.Split` to extract the inner value
- **Backtick-quoted args** — trimmed of backticks directly via `strings.Trim`
- **Empty args** — skipped (common for cursor-at-end positions)

### Quoting Completion Values

When returning completion values that contain metacharacters, external completers must quote them. Carapace's nushell formatter does this:

```go
if strings.ContainsAny(val.Value, ` {}()[]<>$&"'|;#\`+"`") {
    switch {
    case strings.HasPrefix(val.Value, "~"):
        val.Value = fmt.Sprintf(`~"%v"`, escaper.Replace(val.Value[1:]))
    default:
        val.Value = fmt.Sprintf(`"%v"`, escaper.Replace(val.Value))
    }
}
```

Two quoting modes:

- **Tilde-prefixed values** (`~"rest"`): Tilde is kept outside double quotes to preserve nushell's tilde expansion. Example: `~/Documents/My File.txt` → `~"/Documents/My File.txt"`.
- **All other values** (`"value"`): Entire value wrapped in double quotes. Example: `hello world` → `"hello world"`.

Inside double quotes, only backslash and double-quote need escaping:

```go
var escaper = strings.NewReplacer(
    `\`, `\\`,
    `"`, `\"`,
)
```

### The `^` Sigil for External Commands

The `^` sigil executes any string as an external command:

```nu
^'C:\Program Files\exiftool.exe'

let foo = 'C:\Program Files\exiftool.exe'
^$foo
```

This is relevant for completions because nushell distinguishes between internal and external commands. The `^` prefix forces external command lookup.

## String Operations

### Appending/Prepending

```nu
['foo', 'bar'] | each {|s| '~/' ++ $s}  # ~/foo, ~/bar
['foo', 'bar'] | each {|s| '~/' + $s}   # ~/foo, ~/bar
```

### String Comparison

| Operator | Description |
|----------|-------------|
| `=~` | Regex match |
| `!~` | Regex no match |
| `starts-with` | Prefix check |
| `ends-with` | Suffix check |

## Edge Cases

1. **Open quoting** — An open quote like `"partial` is handled by shlex (which closes the quote internally), but the closing quote character is not present in the completion output. The nushell formatter's quoting logic adds quotes back if needed.

2. **Mixed quoting** — If a span contains mixed quote types (e.g., `"it's"`), `shlex.Split` handles the outer quotes correctly but the inner quote is preserved as-is.

3. **No COMP_WORDBREAKS** — Nushell's parser doesn't split on `:` or `=`, so no wordbreak prefix handling is needed (unlike bash).

4. **Backtick strings are not POSIX** — Nushell's backtick syntax is not POSIX shell quoting, so carapace handles them separately with `strings.Trim` rather than `shlex.Split`.

5. **Bare words in command position** — Bare words in the first position are interpreted as commands, not strings. This affects how completions are triggered.

6. **`$` in double quotes** — Unlike bash, `$` is NOT special in nushell double-quoted strings (only in interpolated `$"..."` strings). This means `\$` is literal `\$` in `"..."` but `\$` is an escaped `$` in `$"..."`.

## References

- [Nushell Book: Working with Strings](https://www.nushell.sh/book/working_with_strings.html)
- [Carapace Nushell Patch](https://github.com/carapace-sh/carapace/blob/main/internal/shell/nushell/patch.go)

## Related Skills

- [references/completion.md](references/completion.md) — How completions interact with quoting
- [references/types.md](references/types.md) — SyntaxShape and string types
