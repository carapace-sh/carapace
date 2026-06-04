# Fish Language Fundamentals

In-depth reference for fish's language — quoting, expansion, variables, lists, functions, control flow, command substitution, builtins, and how they differ from POSIX shells.

## Quoting Mechanisms

### 1. Escape Character (`\`)

A backslash preserves the literal value of the next character. Special characters that need escaping:

| Escape | Meaning |
|--------|---------|
| `\$` | Dollar sign |
| `\\` | Backslash |
| `\*` | Star |
| `\?` | Question mark |
| `\~` | Tilde |
| `\#` | Hash |
| `\(` / `\)` | Parentheses |
| `\{` / `\}` | Curly brackets |
| `\[` / `\]` | Square brackets |
| `\<` / `\>` | Angle brackets |
| `\&` | Ampersand |
| `\|` | Pipe |
| `\;` | Semicolon |
| `\"` | Double quote |
| `\'` | Single quote |
| `\ ` | Space |

### Unicode Escapes

| Escape | Meaning |
|--------|---------|
| `\xHH` | Hex byte |
| `\ooo` | Octal character |
| `\uXXXX` | 16-bit Unicode |
| `\UXXXXXXXX` | 32-bit Unicode |
| `\cX` | Control sequence |

### 2. Single Quotes (`'...'`)

No expansions are performed. Only meaningful escapes: `\'` and `\\`.

```fish
echo 'The value is $HOME'     # The value is $HOME (literal)
echo 'it'\''s a test'         # it's a test (end quote, escaped quote, reopen)
```

### 3. Double Quotes (`"..."`)

Only variable expansion (`$VAR`) and command substitution (`$(cmd)`) are performed. Escape sequences like `\n` are **not** interpreted in double quotes.

Meaningful escapes in double quotes: `\"`, `\$`, `\\`, `\` followed by newline (continuation).

```fish
echo "The value is $HOME"     # The value is /home/user (expanded)
echo "say \"hello\""          # say "hello"
echo "cost: \$5"              # cost: $5
echo "hello\nworld"           # hello\nworld (literal \n, not newline)
```

### Key Difference from Bash

Fish does **not** perform word splitting on unquoted variable expansions. `"$var"` and `$var` behave the same way when the variable contains spaces — both produce a single argument. This eliminates the most common quoting bug in POSIX shells.

## Expansion Order

Fish performs expansions in this specific order:

1. **Command substitution** — `$(cmd)` or `(cmd)`
2. **Variable expansion** — `$var`, `$var[1]`
3. **Bracket expansion** — `{a,b,c}`
4. **Wildcard expansion** — `*`, `**`, `?`

Unlike bash, fish does **not** have brace expansion, tilde expansion, parameter expansion, arithmetic expansion, or process substitution as separate expansion phases. Fish's design is simpler and more predictable.

## Variable Expansion

### Basic

```fish
echo $HOME                    # /home/user
set WORD cat
echo "The plural of $WORD is ${WORD}s"  # The plural of cat is cats
```

### Dereferencing (Double `$`)

```fish
set foo a b c
set a 10; set b 20; set c 30
echo $$foo[1]                 # 10 (expand $foo[1] → "a", then expand $a → 10)
```

### Variable Scope

Four scopes available:

| Scope | Flag | Description |
|-------|------|-------------|
| Universal | `-U` / `--universal` | Shared across all Fish sessions, persisted to disk |
| Global | `-g` / `--global` | Specific to current session |
| Function | `-f` / `--function` | Specific to currently executing function |
| Local | `-l` / `--local` | Specific to current block of commands |

```fish
set --global name Patrick
set --local place "at the Krusty Krab"
set -U fish_color_command blue  # Persists across sessions
```

### Scope Resolution

When reading a variable, fish searches scopes in this order:

1. Local (innermost block first)
2. Function
3. Global
4. Universal

When writing a variable with `set` (no scope flag), fish writes to the innermost existing scope. If the variable doesn't exist in any scope, it's created in the function scope (if inside a function) or global scope.

### Exporting Variables

```fish
set -gx EDITOR emacs          # Global + exported
set -ux PATH /custom/bin $PATH  # Universal + exported
```

### Overriding for Single Command

```fish
GIT_DIR=somerepo git status   # Environment variable only for this command
```

## Lists (Arrays)

Fish variables are inherently **lists** (arrays). There is no distinction between a scalar and an array — a scalar is just a list of one element.

### Setting Lists

```fish
set mylist first second third
set fruit apple orange banana
```

### List Indices

- Start at **1** (not 0)
- Negative indices count from end
- Invalid indices are silently ignored

```fish
set fruit apple orange banana
echo $fruit[1]        # apple
echo $fruit[-1]       # banana (last element)
echo $fruit[2..3]     # orange banana (range)
echo $fruit[-1..1]    # banana orange apple (reverse range)
```

### List Operations

```fish
set smurf blue small
set smurf[2] evil     # Change element: blue evil
set -e smurf[1]       # Erase first element: evil
count $smurf          # Count elements: 1
contains blue $smurf  # Check membership: exits 1 (not found)
set -a smurf extra    # Append: blue evil extra
set -p smurf first    # Prepend: first blue evil extra
```

### PATH Variables

Variables ending in `PATH` (like `$PATH`, `$GOPATH`) are special — they're automatically colon-delimited when exported:

```fish
set -gx MYPATH /bin /usr/bin
echo "$MYPATH"        # /bin:/usr/bin
```

### Key Difference from Bash

- Fish lists are 1-indexed; bash arrays are 0-indexed
- Fish has no separate "scalar" type — everything is a list
- Fish does not word-split on `$IFS` — list elements are always separate
- Fish uses `set` instead of `=` for assignment

## Command Substitution

### Syntax

```fish
echo (pwd)           # Preferred: parentheses
echo $(pwd)          # Also valid: dollar-parentheses
```

### Behavior

- Command output is split on **newlines only** (not spaces)
- Each line becomes a separate argument
- Trailing newlines are stripped
- Double-quoted command substitution prevents line splitting

```fish
# Unquoted: each line becomes a separate argument
for line in (cat file.txt)
    echo "Line: $line"
end

# Quoted: all output is one argument
set data "$(cat data.txt)"
```

### Key Difference from Bash

Fish splits command substitution output on **newlines only**, not on `$IFS`. This is fundamentally different from bash and eliminates the COMP_WORDBREAKS problem. `origin/main`, `--option=value`, and `user@host` stay as single tokens.

### Process Substitution

```fish
diff -u (grep fish file1 | psub) (grep fish file2 | psub)
```

`psub` creates a temporary file or named pipe and outputs its path, enabling process output to be used where a filename is expected.

## Brace Expansion

```fish
echo input.{c,h,txt}    # input.c input.h input.txt
cp file{,.bak}           # cp file file.bak
mkdir /usr/local/src/{old,new,dist}
```

Unlike bash, fish's brace expansion is simpler — no sequence expressions like `{1..5}`. Use `(seq 1 5)` instead.

## Wildcard Expansion (Globbing)

| Pattern | Meaning |
|---------|---------|
| `*` | Matches any characters (not `/`) |
| `**` | Matches any characters, descends into subdirectories |
| `?` | Matches any single character except `/` (deprecated) |

```fish
ls *.txt              # All .txt files
ls **/*.jpg           # All .jpg files recursively
```

Fish does not have bash's `extglob`, `globstar` (it's always on), `nullglob` (unmatched patterns produce no matches), or `dotglob` (dotfiles are not matched by `*`).

## Functions

### Defining Functions

```fish
function ll
    ls -l $argv
end
```

### Function Arguments

- `$argv` contains all arguments passed to the function
- Arguments are 1-indexed: `$argv[1]`, `$argv[2]`, etc.

```fish
function greet
    echo "Hello, $argv[1]!"
end
greet World  # Hello, World!
```

### Autoloading Functions

Functions are autoloaded from directories in `$fish_function_path`:

1. `~/.config/fish/functions/` (user functions)
2. `/etc/fish/functions/` (system-wide)
3. Vendor directories under `$XDG_DATA_DIRS`

Function files must be named `<function-name>.fish`.

### Defining Aliases

Fish does not have a separate `alias` command. Use functions instead:

```fish
function ls
    command ls --color=auto $argv
end
```

Or use `alias` (which creates a function under the hood):

```fish
alias ll 'ls -la'  # Creates a function called ll
```

### Function Events

Functions can be triggered by events:

```fish
function on_exit --on-event fish_exit
    echo fish is now exiting
end

function on_variable_change --on-variable fish_mode
    echo "Mode changed to $fish_mode"
end
```

## Control Flow

### If Statement

```fish
if test "$number" -gt 10
    echo "Greater than 10"
else if test "$number" -gt 5
    echo "Greater than 5"
else
    echo "Smaller or equal"
end
```

### Switch Statement

```fish
switch (uname)
case Linux
    echo "Hi Tux!"
case Darwin
    echo "Hi Hexley!"
case '*'
    echo "Hi, stranger!"
end
```

### Combiners (and/or)

```fish
# Run second command if first succeeded
test -e /etc/file && echo "exists"

# Run second command if first failed
test -e /etc/file || echo "doesn't exist"

# Using keywords
set -q XDG_CONFIG_HOME; and set -l configdir $XDG_CONFIG_HOME
or set -l configdir ~/.config
```

### While Loop

```fish
while true
    echo "Still running"
    sleep 1
end
```

### For Loop

```fish
for file in *
    echo "file: $file"
end

for i in (seq 1 5)
    echo $i
end
```

### Begin Block

```fish
begin
    set -l foo bar
end
# $foo is not accessible here
```

### Break and Continue

```fish
break       # Exit loop
continue    # Skip to next iteration
```

## Key Builtins

### `set` — Display and Change Shell Variables

```fish
set [SCOPE] [EXPORT] NAME [VALUE ...]
set -a NAME VALUE ...     # Append
set -p NAME VALUE ...     # Prepend
set -e NAME               # Erase
set -q NAME               # Query (test if defined)
set -n                    # List variable names only
set -S NAME               # Show detailed info
set -L NAME               # Don't abbreviate long values
```

| Scope | Flag | Description |
|-------|------|-------------|
| Universal | `-U` | Persists across all fish instances |
| Global | `-g` | Current session |
| Function | `-f` | Current function |
| Local | `-l` | Current block |

| Export | Flag | Description |
|--------|------|-------------|
| Export | `-x` | Export to child processes |
| Unexport | `-u` | Don't export |

### `test` — Perform Tests

```fish
test EXPRESSION
[ EXPRESSION ]  # Also valid
```

**File operators:** `-d`, `-e`, `-f`, `-r`, `-w`, `-x`, `-L`, `-s`, `-S`, `-O`, `-G`

**String operators:** `=`, `!=`, `-n` (non-empty), `-z` (empty)

**Numeric operators:** `-eq`, `-ne`, `-gt`, `-ge`, `-lt`, `-le`

**Combining:** `-a` (and), `-o` (or), `!` (not), `()` (grouping)

### `string` — Manipulate Strings

| Subcommand | Description |
|------------|-------------|
| `string collect` | Collect multiline input into single output |
| `string escape` | Escape for various contexts (`--style=script/var/url/regex`) |
| `string join SEP` | Join strings with separator |
| `string join0` | Join with NUL separator |
| `string length` | Report string length (`-V` for visible width) |
| `string lower` | Convert to lowercase |
| `string match` | Match patterns (`-r` for regex, `-a` for all, `-i` for case-insensitive) |
| `string pad` | Extend to given width |
| `string repeat` | Repeat string |
| `string replace` | Replace matching substrings (`-a` for all, `-r` for regex) |
| `string shorten` | Truncate with ellipsis |
| `string split SEP` | Split on separator |
| `string split0` | Split on NUL |
| `string sub` | Extract substring |
| `string trim` | Remove leading/trailing whitespace |
| `string unescape` | Unescape strings |
| `string upper` | Convert to uppercase |

Common patterns:

```fish
string length 'hello, world'          # 12
string match -r '(\d+):(\d+)' 2:34   # 2:34, 2, 34
string replace is was 'blue is favorite'  # blue was favorite
string split . example.com            # example, com
string trim '  abc  '                # abc
string upper 'hello'                  # HELLO
```

### `math` — Perform Mathematics

```fish
math [(-s | --scale) N] [(-b | --base) BASE] EXPRESSION
```

**Operators:** `+`, `-`, `*`/`x`, `/`, `^`, `%`

**Constants:** `e`, `pi`, `tau`

**Functions:** `abs`, `ceil`, `floor`, `round`, `sqrt`, `pow`, `sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `atan2`, `sinh`, `cosh`, `tanh`, `log`, `log2`, `log10`, `max`, `min`, `fac`, `ncr`, `npr`, `bitand`, `bitor`, `bitxor`

```fish
math 1+1                    # 2
math -s0 10.0 / 6.0         # 1
math "sin(pi)"              # 0
math --base=hex 192         # 0xc0
```

### `argparse` — Parse Options

```fish
argparse [OPTIONS] OPTION_SPEC ... -- [ARG ...]
```

**Option specification format:** `short/long=`

- `=` — requires a value (last instance saved)
- `=?` — optional value (last instance saved)
- `=+` — requires value (each instance saved)
- `=*` — optional value (each instance saved)
- `&` — don't save in `$argv`
- `#` — integer flag
- `!SCRIPT` — validate value by running SCRIPT

```fish
argparse 'h/help' 'n/name=' -- $argv
or return

if set -q _flag_help
    echo "Usage: my_function [-h | --help]" >&2
    return 1
end

echo "Name: $_flag_name"
```

### `read` — Read Line of Input

```fish
read [OPTIONS] [VARIABLE ...]
```

| Option | Description |
|--------|-------------|
| `-d DELIM` | Split on DELIM |
| `-n NCHARS` | Return after NCHARS characters |
| `-t` | Split using shell tokenization rules |
| `-a` | Store as list in single variable |
| `-z` | Use NUL as line terminator |
| `-s` | Silent mode (passwords) |
| `-p CMD` | Prompt command |
| `-P STR` | Literal prompt string |
| `-S` | Enable syntax highlighting and completions |

```fish
echo hello | read foo              # Store 'hello' in $foo
echo a==b==c | read -d == -l a b c # Split on ==
```

### Other Important Builtins

| Builtin | Description |
|---------|-------------|
| `count` | Count elements (arguments + stdin newlines) |
| `contains` | Test if word is in list (`-i` for index) |
| `status` | Query runtime info (`is-login`, `is-interactive`, `current-command`, `filename`) |
| `echo` | Display text (`-n` no newline, `-s` no spaces, `-e` escape sequences) |
| `printf` | Formatted output |
| `source` / `.` | Execute commands from file in current shell |

## Special Variables

| Variable | Description |
|----------|-------------|
| `$argv` | Arguments to function/script |
| `$status` | Exit status of last foreground job |
| `$pipestatus` | Exit statuses of all pipe elements |
| `$PWD` | Current directory |
| `$HOME` | User's home directory |
| `$PATH` | Command search path (list, auto-colon-joined on export) |
| `$fish_pid` | Shell's process ID |
| `$CMD_DURATION` | Last command duration in milliseconds |
| `$history` | Command history list |
| `$EUID` | Effective user ID (set at startup) |
| `$fish_version` | Fish version string |

## Input/Output Redirection

```fish
command > file          # Redirect stdout
command 2> file         # Redirect stderr
command >> file         # Append stdout
command < file         # Redirect stdin
command &> file         # Both stdout and stderr
command >? file         # Noclobber (don't overwrite existing)
```

### Piping

```fish
cat foo.txt | head
make fish 2>| less      # Pipe stderr
print &| less           # Both stdout and stderr
```

## Comparison with POSIX Shells

| Feature | Fish | Bash |
|---------|------|------|
| Variable assignment | `set var value` | `var=value` |
| Word splitting | No (lists only) | Yes (on `$IFS`) |
| Command substitution splitting | Newlines only | `$IFS` |
| Array indexing | 1-based | 0-based |
| Array syntax | `$var[1]` | `${var[0]}` |
| String quoting | No `$'...'` | `$'...'` ANSI-C quoting |
| Brace expansion | Simple `{a,b}` | Full `{1..5..2}` |
| Arithmetic | `math 1+2` | `$((1+2))` |
| Function args | `$argv` | `$1`, `$2`, `$@` |
| Alias | Function-based | `alias` builtin |
| Test | `test` / `[` | `test` / `[` / `[[` |
| Case matching | `switch`/`case` | `case` in `case`/`select` |
| Nullglob | Always on | `shopt -s nullglob` |
| Process substitution | `(cmd \| psub)` | `<(cmd)` |

## References

- [The fish language — fish-shell docs](https://fishshell.com/docs/current/language.html)
- [set builtin — fish-shell docs](https://fishshell.com/docs/current/cmds/set.html)
- [test builtin — fish-shell docs](https://fishshell.com/docs/current/cmds/test.html)
- [string builtin — fish-shell docs](https://fishshell.com/docs/current/cmds/string.html)
- [math builtin — fish-shell docs](https://fishshell.com/docs/current/cmds/math.html)
- [argparse builtin — fish-shell docs](https://fishshell.com/docs/current/cmds/argparse.html)
- [read builtin — fish-shell docs](https://fishshell.com/docs/current/cmds/read.html)
- [Fish for bash users — fish-shell docs](https://fishshell.com/docs/current/fish_for_bash_users.html)

## Related Skills

- **bash** → `references/quoting-expansion.md` — bash quoting/expansion for comparison
- **fish** → `references/completion.md` — how language features affect completions
