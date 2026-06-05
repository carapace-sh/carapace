# Zsh Expansion and Quoting

In-depth reference for zsh's quoting rules, parameter expansion, globbing/filename generation, and expansion order — with emphasis on how these affect completion.

## Quoting Rules

### Single Quotes (`'...'`)

All characters inside single quotes are treated literally. No parameter expansion, command substitution, history expansion, or globbing occurs.

```zsh
echo 'hello $USER'    # prints: hello $USER
echo 'it'\''s here'   # literal single quote via end-quote, escaped quote, start-quote
```

A single quote cannot appear inside single quotes even with backslash. Use `$'\''` or the `'\''` idiom.

### Double Quotes (`"..."`)

Parameter expansion, command substitution, and arithmetic expansion occur. History expansion does NOT work in double quotes.

```zsh
echo "hello $USER"           # parameter expansion
echo "result: $(cmd)"        # command substitution
echo "value: $((1+2))"       # arithmetic expansion
```

Backslash escapes only these characters inside double quotes: `$`, `` ` ``, `"`, `\`, and (with `RC_QUOTES`) `'`.

### ANSI-C Quoting (`$'...'`)

Supports escape sequences like C strings:

```zsh
echo $'hello\nworld'         # newline
echo $'\t'                   # tab
echo $'\x1b'                 # escape character
echo $'\u0041'               # Unicode A
echo $'\''                   # literal single quote
```

| Escape | Meaning |
|--------|---------|
| `\a` | Alert/bell |
| `\b` | Backspace |
| `\e` | Escape |
| `\f` | Form feed |
| `\n` | Newline |
| `\r` | Carriage return |
| `\t` | Tab |
| `\v` | Vertical tab |
| `\\` | Backslash |
| `\'` | Single quote |
| `\"` | Double quote |
| `\NNN` | Octal character (1-3 digits) |
| `\xNN` | Hex character |
| `\uXXXX` | Unicode (4 hex digits) |
| `\UXXXXXXXX` | Unicode (8 hex digits) |
| `\cX` | Control character |

### Backslash Escaping

A backslash before a character prevents special interpretation. Outside quotes, backslash escapes any character. Inside double quotes, only `$`, `` ` ``, `"`, `\`, and newline are escaped.

### RC_QUOTES Option

When set, allows single quotes inside double-quoted strings to represent a literal single quote:

```zsh
setopt rc_quotes
echo "it's here"    # works without escaping
```

### Quoting and Completion

Quoting state affects how completion candidates are inserted:

- Inside single quotes: only `'` needs escaping (via `'\''`)
- Inside double quotes: `\`, `"`, `$`, `` ` `` need escaping
- Unquoted: all shell metacharacters need escaping

External completion tools must detect the quoting state to properly escape candidates. Carapace uses a 5-state quoting machine for zsh (see carapace-dev skill → `references/shell-zsh.md`).

## Parameter Expansion

### Basic Forms

| Syntax | Description |
|--------|-------------|
| `${name}` | Substitute value of parameter |
| `${+name}` | 1 if parameter is set, 0 otherwise |
| `${name:-word}` | Use default if unset or null |
| `${name-word}` | Use default if unset |
| `${name:=word}` | Set default if unset or null |
| `${name=word}` | Set default if unset |
| `${name::=word}` | Unconditional assignment |
| `${name:+word}` | Use alternate if set and non-null |
| `${name+word}` | Use alternate if set |
| `${name:?word}` | Error if unset or null |
| `${name?word}` | Error if unset |

### Length and Pattern Removal

| Syntax | Description |
|--------|-------------|
| `${#name}` | Length of value (or element count for arrays) |
| `${name#pattern}` | Remove shortest match from beginning |
| `${name##pattern}` | Remove longest match from beginning |
| `${name%pattern}` | Remove shortest match from end |
| `${name%%pattern}` | Remove longest match from end |
| `${name:#pattern}` | Remove array elements matching pattern |
| `${name/pattern/repl}` | Replace first match |
| `${name//pattern/repl}` | Replace all matches |

### Parameter Expansion Flags

Flags are applied as `${(flags)name}`:

| Flag | Description |
|------|-------------|
| `@` | In double quotes, array elements become separate words |
| `A` | Create array parameter (`AA` for associative) |
| `C` | Capitalize resulting words |
| `f` | Split result into lines |
| `F` | Join array words with newlines |
| `i` | Sort case-insensitively |
| `k` | Substitute keys of associative array |
| `L` | Convert to lowercase |
| `n` | Sort numerically |
| `o` | Sort ascending |
| `O` | Sort descending |
| `P` | Indirect reference (value is parameter name) |
| `q` | Quote words with backslashes |
| `Q` | Remove one level of quotes |
| `U` | Convert to uppercase |
| `u` | Expand only first occurrence of unique words |
| `z` | Split into words using shell parsing |
| `0` | Split on null bytes |
| `~` | Force pattern interpretation |
| `j:string:` | Join array with string |
| `s:string:` | Field split on string |
| `l:expr::str1::str2:` | Pad left |
| `r:expr::str1::str2:` | Pad right |

### Common Completion-Related Expansions

```zsh
${(k)opt_args}          # Keys of parsed option arguments
${words[@]:1:$CURRENT-1} # All words except command and current
${(j: :)words}          # Join words array with spaces
${(q)value}             # Quote value for safe insertion
```

## Globbing and Filename Generation

### Basic Glob Operators

| Pattern | Description |
|---------|-------------|
| `*` | Match any string (including null) |
| `?` | Match any single character |
| `[...]` | Match any character in set |
| `[^...]` / `[!...]` | Match any character NOT in set |
| `<x-y>` | Match any number from x to y |
| `(pat)` | Group pattern |
| `x\|y` | Match x OR y (requires `EXTENDED_GLOB`) |
| `x#` | Zero or more of x (`EXTENDED_GLOB`) |
| `x##` | One or more of x (`EXTENDED_GLOB`) |
| `^x` | Match anything except x (`EXTENDED_GLOB`) |
| `x~y` | Match x but NOT y (`EXTENDED_GLOB`) |

### Recursive Globbing

| Pattern | Description |
|---------|-------------|
| `**/` | Zero or more directories (recursive) |
| `***/` | Like `**/` but follows symlinks |
| `(*/)#` | Zero or more directory matches |

### Glob Qualifiers (trailing parentheses)

Qualifiers filter and sort matches:

```zsh
ls *(.x0)       # executable files
ls -d *(/)      # directories only
ls *(L+1k)      # files larger than 1KB
ls *(om[1,5])   # 5 most recently modified files
ls *(^@)        # everything except symlinks
```

| Qualifier | Description |
|-----------|-------------|
| `/` | Directories |
| `.` | Plain files |
| `@` | Symbolic links |
| `=` | Sockets |
| `p` | Named pipes (FIFOs) |
| `*` | Executable files |
| `%` | Device files |
| `r`/`w`/`x` | Owner read/write/execute |
| `A`/`I`/`E` | Group read/write/execute |
| `R`/`W`/`X` | World read/write/execute |
| `s` | Setuid files |
| `S` | Setgid files |
| `t` | Sticky bit |
| `U` | Owned by effective user |
| `G` | Owned by effective group |
| `u:id:` | Owned by user |
| `g:id:` | Owned by group |
| `a[Mwhms][-/+]n` | Accessed n days ago |
| `m[Mwhms][-/+]n` | Modified n days ago |
| `c[Mwhms][-/+]n` | Changed n days ago |
| `L[+/-]n` | Size comparison (bytes) |
| `^` | Negate following qualifiers |
| `-` | Toggle symlink following |
| `M` | Mark directories |
| `T` | Append type mark |
| `N` | Null glob for this pattern |
| `D` | Glob dots (include hidden files) |
| `n` | Numeric sort |
| `o:c:` | Sort by criteria |
| `O:c:` | Sort descending |
| `[beg,end]` | Range of results |
| `:modifier` | Apply history modifier |

### Globbing Flags (Extended Glob)

Flags appear before the pattern in `(#flag)` syntax:

| Flag | Description |
|------|-------------|
| `(#i)` | Case insensitive |
| `(#l)` | Lowercase matches both cases |
| `(#I)` | Case sensitive (negate i/l) |
| `(#b)` | Activate backreferences (capture groups) |
| `(#B)` | Deactivate backreferences |
| `(#cN,M)` | Repeat N to M times |
| `(#m)` | Set `$MATCH`, `$MBEGIN`, `$MEND` |
| `(#aN)` | Approximate matching (N errors) |
| `(#s)` | Match at start of string |
| `(#e)` | Match at end of string |
| `(#q)` | Ignore for pattern matching (use qualifiers only) |
| `(#u)` | Respect locale for multibyte |
| `(#U)` | Single byte mode |

The `(#b)` flag is particularly important for completion — it enables capture groups in `zstyle list-colors` patterns, allowing different parts of a completion display to be colored differently.

### ksh-like Glob Operators (with `KSH_GLOB`)

| Pattern | Description |
|---------|-------------|
| `@(pat)` | Match pattern exactly once |
| `*(pat)` | Match zero or more times |
| `+(pat)` | Match one or more times |
| `?(pat)` | Match zero or one time |
| `!(pat)` | Match anything except |

## Expansion Order

Zsh performs expansions in this order:

1. **History expansion** (interactive shells only) — `!` references
2. **Alias expansion** — expanded before command parsing
3. **Combined five expansions** (left-to-right):
   - Process substitution: `<(...)`, `>(...)`, `=(...)`
   - Parameter expansion: `$var`, `${var}`
   - Command substitution: `$(cmd)`, `` `cmd` ``
   - Arithmetic expansion: `$((expr))`, `$[expr]`
   - Brace expansion: `foo{a,b,c}`
4. **Filename expansion** (tilde): `~user`, `~nameddir`
5. **Filename generation** (globbing): `*`, `?`, `[...]`

After all expansions, unquoted `\`, `'`, and `"` are removed.

### SH_WORD_SPLIT

Unlike bash, zsh does **not** perform word splitting on unquoted parameter expansions by default. Enable with:

```zsh
setopt sh_word_split    # bash-like behavior
```

Without this, `"$var"` and `$var` produce the same result (no splitting on spaces).

## Special Characters and Completion

### WORDCHARS

The `WORDCHARS` parameter controls which characters are treated as part of a word by ZLE word-movement commands. Default: `*?_-.[]~=/&;!#$%^(){}<>`

This affects how `PREFIX` and `SUFFIX` are determined during completion — characters in `WORDCHARS` are not treated as word boundaries.

### FIGNORE / fignore

Array of suffixes to ignore during completion:

```zsh
fignore=(.o .class .pyc)
```

Files ending with these suffixes are excluded from file completion.

### Nomatch vs Nullglob

By default, zsh prints an error if a glob pattern matches nothing:

```zsh
setopt nomatch     # default: error on no match
setopt null_glob   # silently remove non-matching patterns
setopt csh_null_glob  # remove pattern, error if ALL patterns fail
```

## References

### Documentation

- [Zsh Manual: Expansion](https://zsh.sourceforge.io/Doc/Release/Expansion.html) — official expansion reference
- [Zsh Manual: Substitution](https://zsh.sourceforge.io/Doc/Release/Substitution.html) — parameter expansion details
- [Zsh Manual: Filename Generation](https://zsh.sourceforge.io/Doc/Release/Filename-Generation.html) — globbing reference
- [zshexpn(1) man page](https://linux.die.net/man/1/zshexpn) — expansion and substitution reference
- [zshparam(1) man page](https://linux.die.net/man/1/zshparam) — parameter reference
