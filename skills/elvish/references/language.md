# Elvish Language Fundamentals

In-depth reference for the elvish programming language — value types, expressions, quoting, variables, pipelines, control flow, namespaces, modules, and how the language model differs from POSIX shells.

## Design Philosophy

Elvish is intentionally incompatible with bash. Key design differences:

- **Structured data** — values are not just strings; lists, maps, functions, and booleans are first-class
- **Value pipelines** — pipelines carry structured values, not just byte streams
- **No word splitting** — variables are not split on whitespace; `$var` is a single value
- **No implicit string interpolation** — use concatenation (juxtaposition) instead
- **Lexical scoping** — variables are resolved at compile time, not runtime
- **Clean argument passing** — completers receive properly parsed arguments, not raw shell input

## Value Types

### String

Sequence of bytes (possibly empty). Three syntaxes:

| Syntax | Name | Rules |
|--------|------|-------|
| `bareword` | Bareword | ASCII letters, numbers, `!%+,-./:@\_`, non-ASCII printable. No quoting needed. |
| `'single'` | Single-quoted | All characters literal except `'` → `''` for one quote |
| `"double"` | Double-quoted | Supports escape sequences: `\n`, `\t`, `\x41`, `\u0041`, `\U00000041`, `\cX` |

**No string interpolation** — use concatenation:

```elvish
var name = world
echo "hello "$name  # concatenation, not interpolation
```

### Number

| Type | Syntax | Examples |
|------|--------|---------|
| Integer | decimal, hex, octal, binary | `10`, `0xA`, `0o12`, `0b1010` |
| Rational | two integers with `/` | `1/2`, `3/4` |
| Float | decimal point or scientific | `1.0`, `1e1`, `1.0e1` |
| Special | `+Inf`, `-Inf`, `NaN` | |

Digits can be separated by underscores: `1_000_000`.

**Exactness**: integers and rationals are exact; floats are inexact (IEEE 754 double-precision).

### List

```elvish
var li = [a b c]
```

- Indexed with non-negative (from start) or negative (from end) integers
- Slices: `$li[a..b]` (excludes end), `$li[a..=b]` (includes end)

```elvish
$li[0]    # → a
$li[-1]   # → c
$li[0..2] # → [a b]
```

### Map

```elvish
var m = [&key1=value1 &key2=value2]
var empty = [&]
```

- Keys without `=` default to `$true`
- Keys with `=` but no value default to empty string
- Indexed by any key: `$m[key1]`

### Pseudo-map

Values that can be indexed like maps but don't support full map operations. Examples: exceptions, functions.

Printed as `[^tag &key=value]`.

### Nil

`$nil` — initial value of declared but unassigned variables.

### Boolean

`$true` and `$false`. Boolean conversion: `$nil` and exceptions → `$false`; all others → `$true`.

### Exception

No literal syntax. Carries error information as a pseudo-map with `reason` field:

| `reason.type` | Additional fields |
|----------------|-------------------|
| `fail` | `content` |
| `flow` | `name` (break, continue, return) |
| `pipeline` | `exceptions`, `cmd-name`, `pid`, `exit-status` |
| `external-cmd/*` | `cmd-name`, `pid`, `exit-status`, `signal-name`, `signal-number`, `core-dumped` |

### File

Returned by `file:open` and `path:temp-file`. Pseudo-map with `fd` (int) and `name` (string) fields.

### Function

```elvish
{|a b @rest| body }  # lambda with positional, rest, and option parameters
```

Functions don't have return values — they **output** values. Functions are pseudo-maps with fields: `arg-names`, `rest-arg`, `opt-names`, `opt-defaults`, `def`, `body`, `src`.

## Expressions

### Order of Evaluation (Highest to Lowest)

1. **Literals, variables, captures, braced lists** — primary expressions
2. **Indexing** — `expr[index]`
3. **Compounding** — concatenation via juxtaposition (tilde/wildcards unevaluated)
4. **Tilde expansion** — `~` at beginning of compound expression
5. **Wildcard expansion** — `*`, `?`, `**`

### Variable Use

```elvish
$varname       # use variable
$@listvar      # explode list elements
$@stringvar    # explode string codepoints
```

### Output Capture

```elvish
var result = (some command)  # captures all output values
```

- Byte output is split by newlines
- Value output is collected as-is
- Does not introduce new scopes

### Exception Capture

```elvish
var err = ?(some command)  # returns exception or $ok
```

Exceptions are booleanly false; `$ok` is booleanly true.

### Braced List

```elvish
{a b c}  # evaluates to all constituent values
```

### Indexing

```elvish
$li[0]        # list index
$m[key]       # map index
$li[0][0]     # chained indexing
```

Multiple indices applied to each value.

### Compounding (Concatenation)

Multiple expressions with no space between them concatenate:

```elvish
var prefix = hello
echo $prefix-world  # → helloworld
```

Numbers are implicitly converted to strings. Multiple values generate all combinations (Cartesian product).

### Tilde Expansion

```elvish
~           # home directory
~user       # user's home directory
```

Only at the beginning of a compound expression.

### Wildcard Expansion

| Pattern | Matches |
|---------|---------|
| `?` | One character (except `/`) |
| `*` | Any characters (except `/`) |
| `**` | Any characters including `/` |

**Modifiers** (appended with `:`):

| Modifier | Effect |
|----------|--------|
| `nomatch-ok` | No error on no match |
| `but:xxx` | Exclude matches |
| `type:dir` | Only directories |
| `type:regular` | Only regular files |
| `match-hidden` | Match `.` at start of filename |
| `set:abc` | Match characters in set |
| `range:a-z` | Match character range |

Default: error on no match, don't match `.` at filename start.

## Metacharacters

Characters with special meaning in elvish:

| Char | Meaning |
|------|---------|
| `$` | Introduces variable use |
| `*`, `?` | Wildcards |
| `(`, `)` | Output capture |
| `[`, `]` | List/map literal, indexing |
| `{`, `}` | Lambda literal, braced list |
| `<`, `>` | IO redirection |
| `;` | Pipeline separator |
| `\|` | Pipeline connector, function signature |
| `&` | Background pipeline, key-value pair |

## Quoting

### Single-Quoted Strings

All characters are literal. Two consecutive single quotes produce one:

```elvish
'it''s'  # → it's
```

### Double-Quoted Strings

Support escape sequences:

| Escape | Meaning |
|--------|---------|
| `\a`, `\b`, `\t`, `\n`, `\v`, `\f`, `\r` | Control characters |
| `\e` | Escape (0x1B) |
| `\"`, `\\` | Literal quote and backslash |
| `\xHH` | Hex byte (2 digits) |
| `\ooo` | Octal byte (3 digits) |
| `\uHHHH` | Unicode codepoint (4 hex digits) |
| `\UHHHHHHHH` | Unicode codepoint (8 hex digits) |
| `\cX` | ASCII control character |

### Barewords

Characters: ASCII letters, numbers, symbols `!%+,-./:@\_`, non-ASCII printable. `~` and `=` allowed when not parsed as metacharacters.

### Line Continuation

```elvish
echo hello ^  # continues on next line
     world
```

`^` followed by newline continues the current line.

## Variables

### Declaration and Assignment

```elvish
var a b = foo bar     # declare and assign
set x = value         # assign existing variable
set x[0] = newvalue   # assign list element
set x @rest = a b c d # rest assignment
```

### Temporary Assignment

```elvish
tmp x = value         # restored when function finishes
with x = value { ... } # run block with temporary value
with [x = new-x] [y = new-y] { ... } # multiple temporaries
```

### Deletion

```elvish
del x          # delete variable
del m[key]     # delete map entry
```

### Variable Suffixes

| Suffix | Constraint | Default |
|--------|-----------|---------|
| `~` | Only callable values (functions, external commands) | `nop` |
| `:` | Only namespaces | — |

### Scoping

**Lexical scoping** — variable lookup goes from current scope → parent scope → ... → builtin namespace. Nonexistent variables cause compilation errors.

### Qualified Names

```elvish
namespace:var       # namespace-qualified variable
a:b:c              # chain: a → b → c
```

### Closure Semantics

Functions capture variables from outer scopes (**upvalues**). Elvish only captures variables that are actually used.

## Pipelines

### Byte Pipelines

```elvish
echo hello | tr h H   # stdout → stdin (byte channel)
```

### Value Pipelines

```elvish
put [a b c] | each {|x| echo $x }  # value output → value input
```

Each port has both a **file channel** (bytes) and a **value channel** (structured values). Pipelines connect both channels.

### Pipeline Execution

Commands in a pipeline run **in parallel**. The pipeline terminates when all commands finish.

### Pipeline Exceptions

- Single exception → rethrown
- Multiple exceptions → composite exception
- SIGPIPE exceptions suppressed if next command terminated first

### Background Pipeline

```elvish
command &  # runs without waiting
```

Exceptions don't affect the parent.

## IO Ports and Redirection

| Port | Purpose |
|------|---------|
| 0 | stdin |
| 1 | stdout |
| 2 | stderr |

### Redirection Operators

| Operator | Meaning |
|----------|---------|
| `<` | Read from file |
| `>` | Write to file |
| `>>` | Append to file |
| `<>` | Read+write |

### Redirection Sources

| Source | Example |
|--------|---------|
| Filename | `> output.txt` |
| File object | `> $file-obj` |
| Duplicate port | `>&2` (redirect stdout to stderr) |
| Close port | `>&-` (close the port) |

Redirections are applied in the order they appear.

## Special Commands

### `var` — Declare Variables

```elvish
var a b = foo bar
var x              # declared but unassigned ($nil)
```

### `set` — Assign Variables

```elvish
set x = value
set x[0] = newvalue
set x @rest = a b c d
```

Lists and maps are immutable — `set` creates new values.

### `tmp` — Temporary Assignment

```elvish
tmp x = value  # restored when function finishes
```

### `with` — Run with Temporary Assignment

```elvish
with x = new { echo $x }
with [x = new-x] [y = new-y] { echo $x $y }
```

### `del` — Delete Variables/Elements

```elvish
del x
del m[key]
```

### `and`, `or`, `coalesce` — Logics

```elvish
and $true $false   # → $false (short-circuit)
or $true $false    # → $true (short-circuit)
coalesce $nil a    # → a (first truthy value)
```

### `if` — Condition

```elvish
if <condition> { <body> } elif <condition> { <body> } else { <body> }
```

Condition is a single expression (not a command). Multiple values are AND'ed.

### `while` — Conditional Loop

```elvish
while <condition> { <body> } else { <else-body> }
```

### `for` — Iterative Loop

```elvish
for <var> <container> { <body> } else { <else-body> }
```

### `try` — Exception Control

```elvish
try { <try-block> } catch e { <catch-block> } else { <else-block> } finally { <finally-block> }
```

At least `catch` or `finally` required. `else` block runs only if no exception.

### `fn` — Function Definition

```elvish
fn name { <body> }
fn name {|a b @rest| <body> }
```

Defines a function that captures `return`. Equivalent to `var name~ = {|a b @rest| body }`.

### `pragma` — Language Pragmas

```elvish
pragma unknown-command = disallow
```

Affects command resolution for the current scope.

## Namespaces and Modules

### Namespace Syntax

```elvish
namespace:command
namespace:var
```

### Special Namespaces

| Namespace | Purpose |
|-----------|---------|
| `e:` | External commands — `e:ls` calls the external `ls` |
| `E:` | Environment variables — `$E:PATH` accesses `$PATH` |
| `builtin:` | Builtin commands and variables |
| `edit:` | Editor API (interactive mode only) |
| `local:` | Current scope's local namespace |
| `up:` | Parent scope's namespace |

### Module Import

```elvish
use module-name          # imports into current namespace
use module-name alias    # imports with alias
use str                  # → str:compare, str:contains, etc.
```

### Pre-defined Modules

| Module | Purpose |
|--------|---------|
| `builtin` | Built-in commands and variables |
| `edit` | Editor API (interactive mode only) |
| `epm` | Elvish Package Manager |
| `math` | Mathematical functions and constants |
| `path` | Path manipulation (join, split, etc.) |
| `platform` | Platform information (os, arch) |
| `re` | Regular expressions |
| `readline-binding` | Readline-compatible keybindings |
| `str` | String manipulation |
| `unix` | Unix-specific (rlimits, umask) |
| `md` | Markdown rendering |
| `doc` | Documentation access |
| `file` | File operations (open, seek, tell) |
| `os` | OS operations (mkdir-all, symlink, rename) |
| `runtime` | Runtime paths |
| `store` | Persistent store access |

### User-Defined Modules

`.elv` files in module search directories. Relative imports: `./module` or `../module`.

Modules are cached after first import. Re-importing returns the cached namespace.

### Circular Dependencies

Not supported — elvish detects and reports circular imports.

## Key Differences from POSIX Shells

| Feature | Elvish | Bash/Zsh |
|---------|--------|----------|
| Word splitting | No — `$var` is always one value | Yes — `$var` is split by IFS |
| String interpolation | No — use concatenation | Yes — `"hello $name"` |
| Pipeline data | Both byte and value channels | Byte only |
| Variable scope | Lexical | Dynamic (bash), lexical with quirks (zsh) |
| Arrays | `[a b c]` literal | `(a b c)` or `a b c` (word-split) |
| Maps | `[&k=v]` literal | Associative arrays (bash 4+) |
| Exception handling | `try`/`catch`/`finally` | `trap` / `set -e` |
| Function return | Output values, not return | Exit status + stdout |
| Command resolution | `e:` prefix for externals | PATH search by default |
| Null separation | Value channel | `xargs -0`, `print0` |

## References

- [Elvish Language Specification](https://elv.sh/ref/language.html) — official language reference
- [Elvish Builtin Functions](https://elv.sh/ref/builtin.html) — builtin commands and variables
- [Elvish Learn: Tour](https://elv.sh/learn/tour.html) — quick tour of the language
- [Elvish Philosophy](https://elv.sh/learn/philosophy.html) — design philosophy and FAQ

## Related Skills

- For how the language model affects completion (clean args, no word splitting), see [references/completion.md](completion.md).
- For module and startup configuration, see [references/startup-config.md](startup-config.md).
- For carapace-specific elvish integration, see the **carapace-dev** skill → `references/shell-elvish.md`.
