# Tcsh Quoting, Expansion, and Substitution

In-depth reference for tcsh's quoting rules, variable expansion, command substitution, history expansion, alias expansion, and filename substitution (globbing).

## Quoting

Quoting prevents the shell from interpreting special characters. tcsh supports three quoting mechanisms.

### Single Quotes (`'...'`)

Single quotes prevent **all** substitutions:

- No variable expansion (`$var` is literal)
- No command substitution (backticks are literal)
- No history expansion (except `!` — see below)
- The quoted string never yields multiple words

```tcsh
echo '$HOME'       # prints: $HOME
echo '`cmd`'       # prints: `cmd`
echo 'hello world' # prints: hello world
```

**Important**: Single quotes do **not** prevent history expansion (`!`). The `!` character is processed before quoting in tcsh. To prevent history expansion, use backslash: `\!`.

### Double Quotes (`"..."`)

Double quotes allow **variable and command substitution** but prevent other substitutions:

- Variable expansion occurs: `"$HOME"` → the value of `HOME`
- Command substitution occurs: `` "`cmd`" `` → output of `cmd`
- Filename substitution (globbing) is prevented
- History expansion is **not** prevented (use `\!`)
- Blanks and tabs inside command substitution are retained
- Only newlines force new words

```tcsh
echo "$HOME"      # prints: /home/user
echo "`date`"     # prints: current date
echo "hello world" # prints: hello world (one word)
```

### Backslash (`\`)

Backslash quotes a single character, preventing its special meaning:

```tcsh
echo \$HOME    # prints: $HOME
echo \*        # prints: *
echo \"hello\" # prints: "hello"
```

Special cases:

- `\` followed by newline: equivalent to a blank (line continuation)
- Inside quotes: `\` followed by newline results in a literal newline
- `\!` prevents history expansion (unlike single quotes)

### The `backslash_quote` Variable

When set, backslashes always quote `\`, `'`, and `"`:

```tcsh
set backslash_quote
```

This can simplify complex quoting but may cause csh(1) compatibility errors. Useful when embedding quotes in alias definitions.

### C-Style Escape Sequences

tcsh supports `$'...'` syntax for C-style escape sequences:

```tcsh
echo $'hello\nworld'  # prints: hello
                      #         world
echo $'tab\there'     # prints: tab	here
```

Supported sequences: `\a` (bell), `\b` (backspace), `\e` (escape), `\f` (form feed), `\n` (newline), `\r` (carriage return), `\t` (tab), `\v` (vertical tab), `\\` (backslash), `\'` (single quote), `\nnn` (octal).

### Quoting Summary

| Context | `$var` | `` `cmd` `` | `!` | `*?[]` | Space |
|---------|--------|------------|-----|---------|-------|
| Unquoted | Expanded | Expanded | Expanded | Globbed | Word split |
| Single-quoted | Literal | Literal | **Expanded** | Literal | Literal |
| Double-quoted | Expanded | Expanded | **Expanded** | Literal | Literal |
| Backslash-escaped | Literal | Literal | Literal | Literal | Literal |

## Variable Expansion

Variable substitution is triggered by `$` and occurs after alias substitution and parsing, but before command execution.

### Basic Forms

| Syntax | Description |
|--------|-------------|
| `$name` | Words of variable value, separated by blanks |
| `${name}` | Same, but unambiguous (e.g., `${name}suffix`) |
| `$name[selector]` | Select specific words; selector is a number or range |
| `${name[selector]}` | Same with braces |
| `$0` | Name of the file from which input is being read |
| `$number` | Equivalent to `$argv[number]` |
| `$*` | Equivalent to `$argv` |

### Selector Syntax

| Selector | Meaning |
|----------|---------|
| `n` | nth word (1-based) |
| `n-m` | Words n through m |
| `n-` | Words n through end |
| `-m` | Words 1 through m |
| `*` | All words |
| `$` | Last word |

### Unmodified Substitutions

| Syntax | Description |
|--------|-------------|
| `$?name` | Returns `1` if variable is set, `0` if not |
| `$#name` | Number of words in variable |
| `$$` | Process ID of parent shell |
| `$!` | Process ID of last background process |
| `$<` | Read a line from stdin (no further interpretation) |

### Modifiers

History modifiers can be applied to variable substitutions:

| Modifier | Effect |
|----------|--------|
| `:h` | Remove trailing pathname component |
| `:t` | Remove leading pathname components |
| `:r` | Remove filename extension |
| `:e` | Keep only extension |
| `:u` | Uppercase first lowercase letter |
| `:l` | Lowercase first uppercase letter |
| `:s/l/r/` | Substitute l for r |
| `:gs/l/r/` | Global substitution |
| `:q` | Quote substituted words (prevent further expansion) |
| `:x` | Quote as separate words |

Example:

```tcsh
set file = /usr/local/bin/script.sh
echo $file:h    # /usr/local/bin
echo $file:t    # script.sh
echo $file:r    # /usr/local/bin/script
echo $file:e    # sh
```

### Variable Expansion in Completions

When a completion uses the `$varname` list type, the variable is expanded at completion time:

```tcsh
set hosts = (server1 server2 server3)
complete ssh 'p/1/$hosts/'
```

Changes to the variable take effect immediately for subsequent completions.

## Command Substitution

### Syntax

`` `command` ``

### Behavior

- Output is broken into separate words at blanks, tabs, and newlines
- Null words are discarded
- Output undergoes variable and command substitution
- The single final newline does not force a new word

### Inside Double Quotes

- Blanks and tabs are retained
- Only newlines force new words
- The single final newline does not force a new word

### The `csubstnonl` Variable

Since tcsh 6.12, all newline and carriage return characters in command output are replaced by spaces by default. This prevents multi-line output from being split into multiple words unexpectedly.

To restore the old behavior (newlines create new words):

```tcsh
unset csubstnonl
```

### Command Substitution in Completions

The backtick list type in the `complete` builtin uses command substitution:

```tcsh
complete kill 'p/*/`ps -e -o pid=`/'
```

When used for completions, `COMMAND_LINE` is set in the environment before execution.

## History Expansion

History substitutions introduce words from the history list into the input stream. They begin with `!` (or `^` at the start of a line).

### The `histchars` Variable

Changes the history substitution characters:

```tcsh
set histchars = "!^#"   # default: ! ^ #
```

### Event Specification

| Pattern | Meaning |
|---------|---------|
| `!n` | Event number n |
| `!-n` | n events before current |
| `!!` | Previous event (same as `!-1`) |
| `!string` | Most recent event beginning with `string` |
| `!?string?` | Most recent event containing `string` |
| `!#` | Current event (the line being typed) |

### Word Designators

| Designator | Meaning |
|------------|---------|
| `:0` | First (command) word |
| `:n` | nth argument |
| `:^` | First argument |
| `:$` | Last argument |
| `:%` | Word matched by `?string?` search |
| `:x-y` | Range of words (x through y) |
| `:*` | All arguments (same as `:^-$`) |
| `:-` | Command and all arguments except last |
| `:x-` | Words x through second-to-last |

### Word Modifiers

| Modifier | Effect |
|----------|--------|
| `:h` | Remove trailing pathname component |
| `:t` | Remove leading pathname components |
| `:r` | Remove filename extension |
| `:e` | Keep only extension |
| `:u` | Uppercase first lowercase letter |
| `:l` | Lowercase first uppercase letter |
| `:s/l/r/` | Substitute l for r (first occurrence) |
| `:gs/l/r/` | Global substitution |
| `:p` | Print the result but do not execute |
| `:q` | Quote the substituted words |
| `:x` | Quote as separate words |

### Preventing History Expansion

| Method | Works In |
|--------|----------|
| `\!` | Anywhere (backslash) |
| Single quotes | **Does NOT prevent** `!` expansion |
| Double quotes | **Does NOT prevent** `!` expansion |
| `set noglob` | Does not affect history expansion |
| Space after `!` | `! ` is not history expansion |

This is a critical difference from bash: in tcsh, single quotes do **not** prevent history expansion. Only backslash works.

### History Expansion and Completions

History expansion can interfere with completion. The `autoexpand` variable runs `expand-history` before each completion attempt, which can cause unexpected behavior if the current word contains `!`.

## Alias Expansion

After a command line is parsed, the first word of each command is checked for an alias. If found, the first word is replaced by the alias definition.

### Alias Expansion Rules

- If the alias contains a history reference, it undergoes history substitution as though the original command were the previous input line
- Alias substitution is repeated until the first word has no alias
- Loops are detected and cause an error (except self-referential aliases with `\!*`)
- Any quoting of any character in a word for which an alias is defined prevents substitution

### The `\!*` Pattern in Aliases

The `\!*` pattern in an alias is replaced by the arguments to the command:

```tcsh
alias print 'pr \!* | lpr'
```

When `print file1 file2` is executed, `\!*` is replaced with `file1 file2`.

### Preventing Alias Expansion

Quote any character of the aliased word:

```tcsh
\ls    # Prevents alias expansion of 'ls'
'ls'   # Also prevents alias expansion
```

## Filename Substitution (Globbing)

A word containing `*`, `?`, `[`, or `{`, or beginning with `~`, is a candidate for filename substitution.

### Pattern Matching

| Pattern | Matches |
|---------|---------|
| `*` | Any string (including empty) |
| `?` | Any single character |
| `[...]` | Any one character in brackets |
| `[^...]` | Any character NOT in brackets |
| `[a-z]` | Any character in range a-z |
| `{a,b,c}` | Brace expansion: `a`, `b`, or `c` |

### Dot File Rules

A `.` at the beginning of a filename or after `/` must be matched explicitly:

- `*` does **not** match `.bashrc`
- `.*` does match `.bashrc`
- Unless `globdot` or `globstar` is set

### Tilde Expansion

| Pattern | Expansion |
|---------|-----------|
| `~` | User's home directory |
| `~user` | Home directory of `user` |
| `~+` | Current working directory (same as `cwd`) |
| `~-` | Previous directory (same as `owd`) |

Tilde must be at the beginning of a word or after `/`.

### Globstar (`**`)

When `globstar` is set:

- `**` matches any string including `/` (recursive directory descent)
- `***` also descends into symbolic link directories

```tcsh
set globstar
ls **/*.c    # All .c files in current directory and subdirectories
```

### Error Handling

| Variable | Behavior |
|----------|---------|
| (default) | Error if no files match the pattern |
| `nonomatch` | Unmatched patterns left unchanged (no error) |
| `noglob` | Disable filename substitution entirely |

### The `fignore` Variable

List of file suffixes to ignore during completion:

```tcsh
set fignore = (.o .a .class)
```

This affects both interactive completion and the `complete` builtin (unless a select pattern is used, which overrides `fignore`).

### Globbing and Completions

Filename completion uses the same globbing rules, but:

- Completion is prefix-based: the partial word is the prefix
- `fignore` filters out unwanted suffixes
- `nostat` prevents `stat(2)` calls on specified directories
- `recognize_only_executables` limits command completions to executable files

## Metacharacters

### Word Separators

These characters always separate words, regardless of whitespace:

| Character | Meaning |
|-----------|---------|
| `&` | Background |
| `\|` | Pipe |
| `;` | Command separator |
| `<` | Redirect input |
| `>` | Redirect output |
| `(` | Subshell start |
| `)` | Subshell end |

Doubled forms: `&&`, `\|\|`, `<<`, `>>`, `<&`, `>&`.

### Quoting Metacharacters

To use a metacharacter literally:

```tcsh
echo "hello;world"   # double quotes prevent ; as separator
echo hello\;world     # backslash prevents ; as separator
echo 'hello;world'    # single quotes prevent ; as separator
```

## Quoting and Completion

### How Quoting Affects Completion Candidates

When tcsh generates completion candidates, it must handle quoting in the current word:

1. **Unquoted word**: The partial word is matched literally against candidates
2. **Quoted word**: The quote context is tracked; candidates are inserted with appropriate quoting
3. **Backtick context**: Inside backticks, completion switches to command/variable context

### The `quoter` in Carapace's Tcsh Formatter

Carapace's tcsh action formatter (`internal/shell/tcsh/action.go`) uses a `quoter` replacer to escape special characters in completion candidates:

```go
var quoter = strings.NewReplacer(
    `&`, `\&`,
    `<`, `\<`,
    `>`, `\>`,
    "`", "\\`",
    `'`, `\'`,
    `"`, `\"`,
    `{`, ``,   // TODO: escaping not working
    `}`, ``,   // TODO: escaping not working
    `$`, `\$`,
    `#`, `\#`,
    `|`, `\|`,
    `?`, `\?`,
    `(`, `\(`,
    `)`, `\)`,
    `;`, `\;`,
    ` `, `\ `,
    `[`, `\[`,
    `]`, `\]`,
    `*`, `\*`,
    `\`, `\\`,
)
```

This ensures that completion candidates containing metacharacters are properly escaped when inserted into the command line.

### Description Embedding

Since tcsh has no native description support, carapace embeds descriptions using the `_(description)` format:

```
value_(description)
```

Spaces in descriptions are replaced with underscores. This is a workaround — the user must delete the description portion manually if they don't want it.

## References

- tcsh(1) man page — Lexical structure, substitutions, quoting
- `sh.lex.c` in the tcsh source — lexical analyzer
- `sh.glob.c` in the tcsh source — globbing implementation

## Related Skills

- [references/completion.md](references/completion.md) — how completions use variable expansion and command substitution
- [references/editor.md](references/editor.md) — how the editor handles word boundaries
- [references/execution.md](references/execution.md) — how the shell processes commands after expansion
