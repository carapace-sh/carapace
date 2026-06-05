# Oil Shell Quoting and Expansion

In-depth reference for Oil shell's quoting rules, word splitting, expansion, globbing, and how these differ from bash. Covers both OSH (bash-compatible) and YSH (simple word evaluation) modes.

## Quoting Rules

### Single Quotes

Preserve the literal value of every character except `'`. No escape mechanism inside single quotes — to include a single quote, end the string, add an escaped quote, and start a new string:

```bash
echo 'It'\''s a test'    # It's a test
```

### Double Quotes

Preserve the literal value except for `$`, `` ` ``, `\`, `!` (when `histexpand` is on):

```bash
echo "Hello $USER"       # Variable expansion
echo "Result: $(cmd)"    # Command substitution
echo "Value: \$5"        # Escaped dollar sign
```

### ANSI-C Quoting (`$'...'`)

Characters are expanded as in ANSI C:

| Escape | Meaning |
|--------|---------|
| `\\` | Backslash |
| `\'` | Single quote |
| `\"` | Double quote |
| `\a` | Alert (bell) |
| `\b` | Backspace |
| `\e` | Escape character |
| `\f` | Form feed |
| `\n` | Newline |
| `\r` | Carriage return |
| `\t` | Horizontal tab |
| `\v` | Vertical tab |
| `\xHH` | Hex value |
| `\uHHHH` | Unicode 4-digit |
| `\UHHHHHHHH` | Unicode 8-digit |
| `\nnn` | Octal value |

```bash
echo $'tab:\there'       # tab:    here
echo $'line1\nline2'    # Two lines
```

### YSH String Types

YSH adds additional string types:

| Type | Syntax | Description |
|------|--------|-------------|
| Raw string | `r'...'` | No backslash escaping (e.g., `r'C:\path'`) |
| Unicode string | `u'...'` | Unicode-aware string |
| Byte string | `b'...'` | Byte-oriented string |
| Triple-quoted | `'''...'''` | Multi-line, strips leading whitespace |

```bash
echo r'C:\Program Files\'    # Raw: backslashes are literal
echo r'\' u'\\' b'\\'       # Different escape rules per type
```

### Array Indices Must Be Quoted

In OSH, strings inside array indices **must be quoted** to avoid ambiguity:

```bash
# Invalid in OSH
"${SETUP_STATE[$err.cmd]}"

# Valid
"${SETUP_STATE["$err.cmd"]}"
```

### Character Classes in Globs

Brackets within character classes should be escaped:

```bash
echo [\[]    # rather than echo [[]
echo [\]]    # rather than echo []]
```

### Here Doc Terminators

Delimiters must be on their own line:

```bash
# Invalid — delimiter not on its own line
a=$(cat <<EOF
abc
EOF)

# Valid
a=$(cat <<EOF
abc
EOF
)
```

## Word Splitting

### OSH Word Splitting

OSH has **stricter, more consistent** word splitting than bash:

#### Arrays Inside `${}` Are Not Automatically Split

```bash
# Bash: splits "$@" inside ${undef:-"$@"}
echo ${undef:-"$@"}     # bash: splits; OSH: does NOT split

# OSH: use unquoted $@ for splitting
echo ${undef:-$@}       # OSH: splits as expected
```

#### Assignment Builtins Don't Split or Glob

```bash
# Bash: assigns 4+ variables with glob expansion
declare $vars *.py

# OSH: no splitting/globbing occurs
declare $vars *.py      # literal string, not expanded
```

### YSH Simple Word Evaluation

YSH (with `shopt -s simple_word_eval`) eliminates implicit splitting, globbing, and empty elision:

| Construct | Bash | YSH |
|-----------|------|-----|
| `$x` | Split by IFS, glob | Always **one argument** |
| `$(cmd)` | Split by IFS, glob | Always **one argument** |
| `$((1+2))` | One argument | Always **one argument** |
| Empty `$x` | Elided (no argument) | Empty string argument |
| `$pat` (where pat=`*.py`) | Glob-expanded | Literal `*.py` |

```bash
# YSH: predictable behavior
var pic = 'my pic.jpg'
var empty = ''
var pat = '*.py'
argv ${pic} $empty $pat $(cat foo.txt) $((1 + 3))
# Result: ['my pic.jpg', '', '*.py', 'contents of foo.txt', '4']

# Bash: requires quoting for same result
argv "${pic}" "$empty" "$pat" "$(cat foo.txt)" "$((1 + 3))"
```

### Constructs That Evaluate to 0-N Arguments

Even in YSH, these produce multiple arguments:

1. **Array splicing**: `"$@"`, `"${myarray[@]}"`, `@myarray` (with `parse_at`)
2. **Static globbing**: `echo *.py` (pattern in source code, not from variable)
3. **Brace expansion**: `{alice,bob}@example.com`

### Explicit Split/Glob in YSH

YSH provides explicit operators for when you need split/glob behavior:

```bash
@[split(mystr, IFS?)]    # Explicit word splitting
@[glob(mypat)]           # Explicit globbing
@[maybe(s)]              # Explicit empty elision
```

## Parameter Expansion

### Standard Expansions (OSH)

| Syntax | Description |
|--------|-------------|
| `${var}` | Value of `var` |
| `${var:-word}` | Use `word` if `var` unset or null |
| `${var:=word}` | Assign `word` if `var` unset or null |
| `${var:+word}` | Use `word` if `var` set and non-null |
| `${var:?msg}` | Error with `msg` if `var` unset or null |
| `${#var}` | String length |
| `${var%pattern}` | Remove shortest suffix match |
| `${var%%pattern}` | Remove longest suffix match |
| `${var#pattern}` | Remove shortest prefix match |
| `${var##pattern}` | Remove longest prefix match |
| `${var/pat/repl}` | Replace first match |
| `${var//pat/repl}` | Replace all matches |
| `${var:offset}` | Substring from offset |
| `${var:offset:length}` | Substring with length |

### Type Model Difference

| Aspect | Bash | Oil |
|--------|------|-----|
| Type tagging | **Locations** (cells) are typed | **Values** are typed (like Python/JS) |
| `declare -i` | Marks location as integer | **No-op** in OSH |
| `declare -A x=()` | Creates empty assoc array | Clears content of assoc array |
| `declare x=(one two)` | Creates indexed array | Creates indexed array |
| `declare x=(['k']=v)` | Creates assoc array | Creates **indexed** array (key is integer) |

### `[[ -v var ]]` Limitations

In OSH, `[[ -v var ]]` treats expressions as strings, returning false. Workaround:

```bash
# OSH workaround
if [[ "${assoc['key']:+exists}" ]]; then
    echo "key exists"
fi

# YSH: use proper syntax
var d = { key: 42 }
if ('key' in d) {
    echo "key exists"
}
```

## Tilde Expansion

| Syntax | Description |
|--------|-------------|
| `~` | Home directory (`$HOME`) |
| `~/path` | Home directory + path |
| `~user` | Home directory of `user` |

### strict_tilde Option

With `shopt -s strict_tilde`, failed tilde expansions cause **hard errors** instead of silently evaluating to `~` or `~bad`:

```bash
# Without strict_tilde
echo ~nonexistent    # ~nonexistent (silent fallback)

# With strict_tilde
echo ~nonexistent    # Error: tilde expansion failed
```

## Command Substitution

### Standard Syntax

```bash
output=$(command)       # Preferred
output=`command`        # Legacy (discouraged)
```

### Nested Subshell in Command Substitution

`$((` **always** starts an arithmetic subexpression in OSH. Use a space for nested subshell:

```bash
# Invalid — $(( always means arithmetic
$((cd / && ls))

# Valid — space after $((
$( (cd / && ls) )

# Preferred — use grouping
$({ cd / && ls; })
```

### `((` Always Starts Arithmetic

```bash
# Invalid in OSH for grouping
if ((test -f a || test -f b) && grep foo c)

# Use braces for grouping
if { test -f a || test -f b; } && grep foo c
```

## Arithmetic Expansion

### Standard Syntax

```bash
echo $((1 + 2))         # 3
((a = 42))              # Assignment
```

### Dynamic Command Subs Disallowed by Default

For security, OSH disallows dynamically parsed arithmetic by default:

```bash
# Disallowed by default — avoids code injection
a[$(echo 42)]=value

# Enable if needed
shopt -s eval_unsafe_arith
```

### printf '%d' Stricter

```bash
# Bash: returns 0 with warning
printf '%d' 'notanumber'    # 0

# OSH: requires valid integer
printf '%d' 'notanumber'    # Error
```

## Globbing

### Standard Globs

| Pattern | Description |
|---------|-------------|
| `*` | Match any string |
| `?` | Match any single character |
| `[...]` | Match any character in class |
| `[!...]` | Match any character NOT in class |

### Extended Globs

| Pattern | Description |
|---------|-------------|
| `?(pattern)` | Match zero or one occurrence |
| `*(pattern)` | Match zero or more occurrences |
| `+(pattern)` | Match one or more occurrences |
| `@(pattern)` | Match exactly one occurrence |
| `!(pattern)` | Match anything except |

### Extended Globs Are Static in OSH

```bash
# Works — pattern is in source code
echo *.@(cc|h)

# Does NOT expand — pattern comes from variable
pat='*.@(cc|h)'
echo $pat    # OSH: literal string, not expanded (like mksh)
```

### Extended Glob vs. Negation

OSH distinguishes with a space:

```bash
[[ !(a == a) ]]    # Extended glob
[[ ! (a == a) ]]   # Negation of equality test
```

### Brace Expansion

```bash
echo {a,b,c}         # a b c
echo {1..5}          # 1 2 3 4 5
echo {a,b}{1,2}      # a1 a2 b1 b2
```

### Brace Expansion Is All-or-Nothing

Invalid syntax causes the entire expansion to abort (unlike bash which does partial expansion):

```bash
# Bash: partial expansion → a{ b{
{a,b}{

# OSH: syntax error
{a,b}{
```

### Globbing Options

| Option | Default | Description |
|--------|---------|-------------|
| `noglob` / `-f` | off | Disable filename expansion |
| `nullglob` | off | Unmatched globs expand to nothing |
| `failglob` | off | Error if glob doesn't match |
| `dotglob` | off | Include dotfiles |
| `dashglob` | on (OSH), off (YSH) | Include files starting with `-` |
| `nocaseglob` | off | Case-insensitive globbing |

### dashglob

YSH disables `dashglob` by default to prevent `rm *` from being confused by files called `-rf`:

```bash
# OSH: dashglob on by default
touch -- -rf
echo *          # -rf file.txt

# YSH: dashglob off by default
echo *          # file.txt (no -rf)
```

## Process Substitution

```bash
diff <(sort left.txt) <(sort right.txt)
```

With `shopt -s process_sub_fail`, process substitution failure is treated like `pipefail` for process subs.

## Quoting and Completion

The quoting model directly affects how completion works:

| Aspect | Bash | Oil |
|--------|------|-----|
| What `$2` contains | Last segment after `COMP_WORDBREAKS` | Full argv entry |
| Dequoting | Plugin must handle | Shell handles internally |
| Quoting completions | Plugin must handle | Shell handles internally |
| `COMPREPLY` entries | May need quoting | Just set values (argv) |

See [references/completion.md](references/completion.md) for details on how quoting affects the completion API.

## Differences from Bash: Summary

| Feature | Bash | Oil |
|---------|------|-----|
| Array splitting in `${}` | Splits `"$@"` inside `${:-}` | Does not split; use unquoted `$@` |
| Assignment builtin splitting | Splits and globs | No splitting or globbing |
| Extended globs from variables | Dynamic expansion | Static only (like mksh) |
| Brace expansion errors | Partial expansion | All-or-nothing (syntax error) |
| `declare -i` | Marks location as integer | No-op |
| `$((` in command sub | Ambiguous | Always arithmetic; use `$( (` for subshell |
| `((` for grouping | Works | Always arithmetic; use `{ }` for grouping |
| `printf '%d'` invalid | Returns 0 | Error |
| Type model | Locations typed | Values typed |
| `strict_tilde` | Not available | Failed tilde causes error |
| `dashglob` | Always on | On in OSH, off in YSH |

## References

- [Known Differences Between OSH and Other Shells](https://oils.pub/release/latest/doc/known-differences.html)
- [YSH vs. Shell](https://oils.pub/release/latest/doc/ysh-vs-shell.html)
- [Simple Word Evaluation](https://oils.pub/release/latest/doc/simple-word-eval.html)
- [Shell Options Reference](https://oils.pub/release/latest/doc/ref/chap-option.html)

## Related Skills

- **bash skill → references/quoting-expansion.md** — bash quoting and expansion rules
- **oil skill → references/completion.md** — how quoting affects the completion API
