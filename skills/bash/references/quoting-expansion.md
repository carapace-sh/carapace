# Bash Quoting and Expansion

In-depth reference for bash's quoting rules, expansion order, word splitting, and how they affect completion.

## Quoting Mechanisms

### 1. Escape Character (`\`)

A non-quoted backslash preserves the literal value of the next character, with the exception of newline. A `\newline` pair is treated as line continuation (the pair is removed from the input).

### 2. Single Quotes (`'...'`)

Preserves the literal value of **every** character. No escape is possible — a single quote cannot appear between single quotes, even when preceded by backslash.

```bash
echo 'hello world'       # hello world
echo 'it'\''s'           # it's (end quote, escaped quote, reopen quote)
echo '$HOME \n \t'       # $HOME \n \t  (all literal)
```

### 3. Double Quotes (`"..."`)

Preserves literal value of all characters except: `$`, `` ` ``, `\`, and `!` (when history expansion enabled; in POSIX mode, `!` has no special meaning).

**Backslash in double quotes** — retains special meaning only when followed by: `$`, `` ` ``, `"`, `\`, or `newline`. These pairs are **removed**. Backslashes before other characters are left as-is.

```bash
echo "$HOME"             # /home/user  (parameter expansion)
echo "hello\nworld"      # hello\nworld  (\n is literal)
echo "say \"hello\""     # say "hello"  (\" is removed)
echo "cost: \$5"         # cost: $5  (\$ is removed)
```

### 4. ANSI-C Quoting (`$'...'`)

Expands backslash escape sequences per ANSI C standard. The result is treated as single-quoted.

| Sequence | Meaning |
|----------|---------|
| `\a` | Alert (bell) |
| `\b` | Backspace |
| `\e` / `\E` | Escape character (not in ANSI C) |
| `\f` | Form feed |
| `\n` | Newline |
| `\r` | Carriage return |
| `\t` | Horizontal tab |
| `\v` | Vertical tab |
| `\\` | Backslash |
| `\'` | Single quote |
| `\"` | Double quote |
| `\?` | Question mark |
| `\nnn` | Octal character (1-3 digits) |
| `\xHH` | Hex character (1-2 digits) |
| `\uHHHH` | Unicode character (1-4 hex digits) |
| `\UHHHHHHHH` | Unicode character (1-8 hex digits) |
| `\cx` | Control-x character |

```bash
echo $'hello\nworld'     # Two lines
echo $'\x41'             # A
echo $'\t'               # Tab character
echo $'it\'s'            # it's
```

### 5. Locale Translation (`$"..."`)

Prefixing a double-quoted string with `$` causes translation via `gettext`. If locale is `C`/`POSIX` or no translation exists, the `$` is ignored and the string is treated as double-quoted.

## Expansion Order

Bash performs expansions in this specific order:

1. **Brace expansion** — `a{d,c,b}e` → `ade ace abe`
2. **Tilde expansion** — `~` → `$HOME`
3. **Parameter and variable expansion** — `$var`, `${var}`, `${var:-default}`
4. **Command substitution** — `$(cmd)` or `` `cmd` ``
5. **Arithmetic expansion** — `$((1+2))`
6. **Process substitution** — `<(cmd)`, `>(cmd)` (simultaneous with steps 2-5)
7. **Word splitting** — based on `$IFS`
8. **Filename expansion** (globbing) — `*`, `?`, `[...]`
9. **Quote removal** — always last

Only brace expansion, word splitting, and filename expansion can increase the number of words.

## Brace Expansion

Strictly textual — performed before other expansions. No syntactic interpretation.

```bash
echo a{d,c,b}e           # ade ace abe
echo {1..5}              # 1 2 3 4 5
echo {a..f}              # a b c d e f
echo {01..05}            # 01 02 03 04 05 (zero-padded)
echo {1..10..2}          # 1 3 5 7 9 (step)
mkdir /usr/local/src/{old,new,dist}
```

`${` inhibits brace expansion — it's treated as parameter expansion instead.

## Tilde Expansion

| Syntax | Result |
|--------|--------|
| `~` | `$HOME` |
| `~/foo` | `$HOME/foo` |
| `~fred/foo` | Home directory of user `fred` |
| `~+/foo` | `$PWD/foo` |
| `~-/foo` | `${OLDPWD}/foo` |
| `~N` | Directory stack entry (like `dirs +N`) |
| `~-N` | Directory stack entry (like `dirs -N`) |

Results are treated as quoted — no word splitting or filename expansion on tilde-expanded values.

## Parameter Expansion

### Basic

```bash
${parameter}              # Value of parameter
${#parameter}             # Length in characters
${#array[@]}              # Number of array elements
${!prefix*}               # Variable names beginning with prefix
${!name[@]}               # Array indices/keys
${!parameter}             # Indirect expansion — expand value as new parameter name
```

### Default Values

| Syntax | Meaning |
|--------|---------|
| `${parameter:-word}` | Use `word` if unset or null |
| `${parameter:=word}` | Assign `word` if unset or null |
| `${parameter:+word}` | Use `word` if set and non-null |
| `${parameter:?word}` | Error if unset or null |

With colon: tests for both existence and non-null. Without colon: tests only for existence.

### Substring

```bash
${parameter:offset}           # From offset to end
${parameter:offset:length}   # Length characters from offset
${parameter: -7}              # Last 7 characters (space before - is required)
```

### Pattern Removal

```bash
${parameter#pattern}    # Remove shortest match from beginning
${parameter##pattern}   # Remove longest match from beginning
${parameter%pattern}    # Remove shortest match from end
${parameter%%pattern}   # Remove longest match from end
```

Patterns use glob syntax (not regex). Common uses:

```bash
file=archive.tar.gz
${file#*.}              # tar.gz
${file##*.}             # gz
${file%.*}              # archive.tar
${file%%.*}             # archive
```

### Pattern Substitution

```bash
${parameter/pattern/string}     # Replace first match
${parameter//pattern/string}    # Replace all matches
${parameter/#pattern/string}    # Match at beginning only
${parameter/%pattern/string}    # Match at end only
```

If `string` is null, matches are deleted. `&` in the replacement refers to the matched text (escape with `\&`).

### Case Modification

```bash
${parameter^pattern}    # Uppercase first matching char
${parameter^^pattern}   # Uppercase all matching chars
${parameter,pattern}    # Lowercase first matching char
${parameter,,pattern}   # Lowercase all matching chars
```

Without pattern, all characters are matched:

```bash
${var^}                 # First char to uppercase
${var^^}                # All to uppercase
${var,}                 # First char to lowercase
${var,,}                # All to lowercase
```

### Parameter Transformation

```bash
${parameter@U}          # Uppercase
${parameter@u}          # Uppercase first character
${parameter@L}          # Lowercase
${parameter@Q}          # Quote for reuse as input
${parameter@E}          # Expand backslash escapes
${parameter@P}          # Expand as prompt string
${parameter@A}          # Assignment statement
${parameter@K}          # Quoted key-value pairs
${parameter@a}          # Attribute flags
${parameter@k}          # Separate words
```

## Command Substitution

```bash
$(command)              # Standard form
`command`               # Deprecated backtick form
$(< file)               # Faster than $(cat file)
${| command; }          # Execute in current environment, capture to REPLY
```

- Executes in subshell environment (except `${| ...; }` form)
- Trailing newlines deleted
- Within double quotes: no word splitting or filename expansion

## Arithmetic Expansion

```bash
$(( expression ))
```

- Tokens undergo parameter expansion, command substitution, quote removal
- Empty strings evaluate to 0
- May be nested
- Supports: `+`, `-`, `*`, `/`, `%`, `**`, `++`, `--`, `==`, `!=`, `<`, `>`, `<=`, `>=`, `&&`, `||`, `!`, `~`, `&`, `|`, `^`, `<<`, `>>`, assignment operators, ternary `?:`

## Process Substitution

```bash
<(list)                 # Read output of list from file
>(list)                 # Write input to list via file
```

- Process runs asynchronously
- Filename passed as argument to the current command
- Requires FIFOs or `/dev/fd` support

## Word Splitting

After parameter expansion, command substitution, and arithmetic expansion, bash scans results for word splitting — **only if not within double quotes**.

### IFS (Internal Field Separator)

| IFS State | Behavior |
|-----------|----------|
| Unset | Defaults to `<space><tab><newline>` |
| Null (`IFS=''`) | No word splitting occurs |
| Set | Each character is a delimiter |

**IFS whitespace** (space, tab, newline) is always treated as whitespace regardless of IFS value:
- Leading/trailing IFS whitespace is removed
- Sequences of IFS whitespace delimiters are treated as a single delimiter
- Non-whitespace IFS characters delimit individually

**Null arguments:**
- Explicit null arguments (`""` or `''`) are retained
- Unquoted implicit null arguments (from unset parameters) are removed

### Special Parameter Splitting

| Syntax | Behavior |
|--------|----------|
| `"$*"` | Single word: `"$1c$2c..."` where `c` is first IFS character |
| `"$@"` | Separate words: `"$1" "$2" ...` |
| `$*` (unquoted) | Subject to word splitting and globbing |
| `$@` (unquoted) | Subject to word splitting and globbing |

### Arrays

```bash
"${array[@]}"           # Each element as separate word
"${array[*]}"           # Single word with IFS separator
```

## Filename Expansion (Globbing)

After word splitting, words containing `*`, `?`, or `[` (that are not quoted) are treated as patterns.

| Pattern | Meaning |
|---------|---------|
| `*` | Any string |
| `?` | Any single character |
| `[abc]` | Character class |
| `[!abc]` or `[^abc]` | Negated character class |
| `[[:class:]]` | POSIX character class (`alnum`, `alpha`, `digit`, `lower`, `upper`, etc.) |

### Glob Options

| Option | Effect |
|--------|--------|
| `nullglob` | Remove unmatched patterns (instead of keeping literal) |
| `failglob` | Error if no matches |
| `nocaseglob` | Case-insensitive matching |
| `dotglob` | Match files starting with `.` |
| `globstar` | `**` matches directories recursively |
| `extglob` | Extended patterns: `?(pattern)`, `*(pattern)`, `+(pattern)`, `@(pattern)`, `!(pattern)` |

### FIGNORE

The `FIGNORE` variable filters filename completions:

```bash
export FIGNORE=".bak:~:.old:.swp"   # Colon-separated suffixes to ignore
```

### GLOBIGNORE

Colon-separated patterns to exclude from glob matches. Setting `GLOBIGNORE` implicitly enables `dotglob`.

## Special Parameters

| Parameter | Description |
|-----------|-------------|
| `$*` | All positional parameters, IFS-joined |
| `$@` | All positional parameters, separate words |
| `$#` | Number of positional parameters |
| `$?` | Exit status of most recent command |
| `$-` | Current option flags |
| `$$` | Process ID of shell |
| `$!` | Process ID of most recent background job |
| `$0` | Name of shell or shell script |
| `$_` | Last argument of previous command (at shell start, absolute path of shell) |

## Arrays

### Declaration

```bash
declare -a name          # Indexed array
declare -A name          # Associative array
name=(value1 value2)     # Compound assignment
name=([key]=value)       # With subscripts
name+=(value)            # Append
```

### Referencing

```bash
${name[subscript]}       # Single element
${name[@]}               # All elements (separate words)
${name[*]}               # All elements (single word)
${#name[@]}              # Number of elements
${!name[@]}              # Array indices/keys
```

Negative indices count from end (`-1` = last element).

## Quoting in Completion Functions

### COMP_WORDBREAKS and Completion

`COMP_WORDBREAKS` (default: `" \t\n\"'@><=;|&(:`) determines how Readline splits the command line into words for completion. This creates problems:

- `--option=value` is split at `=`, so `$2` is only `value`
- `user@host:port` is split at `@` and `:`
- `File::Basename` is split at `:`

### Reassembling Words

The bash-completion project provides helpers:

```bash
# Reassemble words excluding specified characters from word breaks
_comp__reassemble_words : =   # Reassemble on : and =

# Get current word at cursor position
_comp__get_cword_at_cursor : =  # Handles wordbreak chars
```

### Dequoting in Completion Functions

When a completion function receives the current word, it may contain shell quoting. The bash-completion project provides:

```bash
# Safely expand a quoted word
_comp_dequote "$cur"     # Result in REPLY

# Shell-quote a word for COMPREPLY
_comp_quote "$value"    # Result in REPLY
```

### Quoting and COMPREPLY

When `complete -o noquote` is used (as carapace does), the completion function is responsible for proper quoting of COMPREPLY entries. When `complete -o filenames` is used, Readline handles quoting automatically.

### The `-o filenames` Option

Tells Readline that completions are filenames. Readline then:
- Adds `/` suffix to directory names
- Quotes special characters
- Suppresses trailing space for directories (when `mark-directories` is on)

### The `-o noquote` Option

Tells Readline not to quote completions. The completion function must handle quoting itself. This is used by external completion frameworks (like carapace) that manage their own quoting.

### The `-o fullquote` Option

Tells Readline to quote all completed words, even if they are not filenames.

## References

- [GNU Bash Manual: Quoting](https://www.gnu.org/software/bash/manual/html_node/Quoting.html)
- [GNU Bash Manual: Shell Expansions](https://www.gnu.org/software/bash/manual/html_node/Shell-Expansions.html)
- [GNU Bash Manual: Shell Parameter Expansion](https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html)
- [GNU Bash Manual: Word Splitting](https://www.gnu.org/software/bash/manual/html_node/Word-Splitting.html)
- [GNU Bash Manual: Filename Expansion](https://www.gnu.org/software/bash/manual/html_node/Filename-Expansion.html)
- [GNU Bash Manual: Special Parameters](https://www.gnu.org/software/bash/manual/html_node/Special-Parameters.html)
- [GNU Bash Manual: Arrays](https://www.gnu.org/software/bash/manual/html_node/Arrays.html)

## Related Skills

- **references/completion.md** — how quoting affects COMP_WORDBREAKS and completion functions
- **references/readline.md** — how Readline handles quoting during completion display
