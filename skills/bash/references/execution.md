# Bash Execution Model

In-depth reference for how bash executes commands, manages processes, handles signals and traps, and interacts with the terminal.

## Command Search and Execution

Bash searches for commands in this order:

1. **Shell functions** — executed in current shell context
2. **Shell builtins** — invoked as builtin within the shell process
3. **PATH lookup** — searches `$PATH` directories, uses hash table for caching
4. **`command_not_found_handle`** — if defined, executed in subshell
5. **External commands** — fork/exec in separate process

If the command name contains a `/`, bash executes it directly as a program (no PATH search).

## Simple Command Expansion Order

When executing a simple command, bash performs:

1. **Variable assignments and redirections** — saved for later processing
2. **Word expansion** — tilde, parameter, command substitution, arithmetic, quote removal
3. **Redirections** — performed
4. **Variable assignment values** — undergo tilde, parameter, command substitution, arithmetic, quote removal

If no command name remains after expansion, variable assignments affect the current shell. Otherwise, they only affect the executed command's environment.

## Execution Environments

### Current Shell Environment

Contains:
- Open files inherited at invocation (modified by `exec` redirections)
- Current working directory
- File creation mode mask (`umask`)
- Current traps
- Shell parameters (set by assignment, `set`, or inherited)
- Shell functions (defined or inherited)
- Shell options (enabled at invocation or by `set`/`shopt`)
- Shell aliases
- Process IDs (background jobs, `$$`, `$PPID`)

### Subshell Environment (for External Commands)

Contains:
- Open files (plus redirections)
- Current working directory
- File creation mode mask
- Exported variables and functions
- Traps reset to inherited values from parent

**Key distinction**: Changes in a subshell cannot affect the parent's environment.

### What Creates a Subshell

| Construct | Subshell? |
|-----------|-----------|
| Command substitution `$(cmd)` | Yes |
| Grouped commands `(cmd)` | Yes |
| Pipeline elements (except last with `lastpipe`) | Yes |
| Background commands `cmd &` | Yes |
| Process substitution `<(cmd)` | Yes |
| Shell functions | **No** — current shell |
| `source`/`.` script | **No** — current shell |
| `exec cmd` | **No** — replaces current process |

## External Commands vs Builtins vs Functions

### External Commands

- Executed in separate process (fork/exec)
- Inherit signal handlers from shell (unless set to ignore)
- Cannot affect shell's execution environment
- `$_` set to full pathname and passed in environment

### Shell Functions

- **Executed in current shell context** — no new process
- Arguments become positional parameters (`$1`, `$2`, etc.)
- `$0` unchanged; `FUNCNAME[0]` set to function name
- Use `local` for local variables
- Dynamic scoping (variables visible to called functions)
- `DEBUG` and `RETURN` traps **not inherited** by default (need `trace` attribute or `-o functrace`)
- `ERR` trap **not inherited** by default (need `-o errtrace`)
- `return` exits function and resumes after call site

### Shell Builtins

- Part of the shell itself
- Execute within the shell process
- In POSIX mode, special builtins affect the shell environment
- Some builtins in pipelines execute in subshell (unless `lastpipe`)

## Pipelines

```bash
[time [-p]] [!] command1 [ | command2 ] ...
```

- Output of each command connected via pipe to input of next
- Each command executed in its **own subshell** (separate process)
- `|&` connects stderr to stdout through pipe
- Exit status: last command's (or last non-zero with `pipefail`)
- With `!`, exit status is logical negation
- With `lastpipe` option (and no job control), last element runs in current shell

```bash
# Enable lastpipe for performance
shopt -s lastpipe
echo "hello" | read var   # var is set in current shell
```

## Job Control

### Process Groups

- Each process has a **process group ID** (PGID)
- **Foreground process group** — PGID equals terminal's PGID; receives keyboard signals
- **Background process groups** — different PGID; immune to keyboard signals
- Background processes attempting terminal I/O receive `SIGTTIN`/`SIGTTOU`

### Job Control Features

| Feature | Description |
|---------|-------------|
| `^Z` (suspend) | Stops current foreground process, returns to shell |
| `^Y` (delayed suspend) | Stops on next read attempt |
| `bg %n` | Resume job in background |
| `fg %n` | Bring job to foreground |
| `jobs` | List jobs |
| `disown %n` | Remove from jobs table (avoid SIGHUP) |
| `wait %n` | Wait for job to complete |

### Job Specifiers

| Specifier | Meaning |
|-----------|---------|
| `%n` | Job number n |
| `%%` or `%+` | Current job |
| `%-` | Previous job |
| `%string` | Job whose command starts with `string` |
| `%?string` | Job whose command contains `string` |

### Job Control Disabled (Non-Interactive)

- Shell and foreground command are in same process group
- `^C` sends `SIGINT` to all processes in that group
- Shell waits for command, then decides action based on whether command handled the signal

## Signals

### Default Signal Handling in Interactive Bash

| Signal | Default Action |
|--------|---------------|
| `SIGTERM` | Ignored |
| `SIGQUIT` | Ignored |
| `SIGINT` | Caught (interruptible `wait`) |
| `SIGTTIN` | Ignored (with job control) |
| `SIGTTOU` | Ignored (with job control) |
| `SIGTSTP` | Ignored (with job control) |

### SIGHUP Handling

- Shell exits by default on `SIGHUP`
- Before exit, interactive shell resends `SIGHUP` to all jobs
- Sends `SIGCONT` to stopped jobs
- `disown` removes from jobs table or marks to avoid `SIGHUP`
- `huponexit` option sends `SIGHUP` to all jobs on interactive login shell exit

### Signal Inheritance for Child Processes

- Non-builtin commands inherit signal handlers from shell (unless set to ignore)
- Async commands (without job control) ignore `SIGINT` and `SIGQUIT`
- Command substitutions ignore `SIGTTIN`, `SIGTTOU`, `SIGTSTP`

## Traps

```bash
trap [-lpP] [action] [sigspec ...]
```

### Trap Types

| Signal | When Triggered |
|--------|---------------|
| `DEBUG` | Before every simple command, `for`/`case`/`select`/`((`/`[[`, and first command in function |
| `EXIT` (0) | When shell exits (any reason) |
| `ERR` | When command returns non-zero (subject to conditions) |
| `RETURN` | When shell function or sourced script finishes |

### ERR Trap Conditions

`ERR` is **not** executed when:

- Command after `while`/`until` or `if`/`elif`
- Command after final `&&`/`||`
- All but last command in pipeline (unless `pipefail`)
- Return status inverted with `!`

### Trap Inheritance

| Trap | Inherited by Functions? | Inherited by Subshells? |
|------|------------------------|------------------------|
| `DEBUG` | Only with `declare -t` or `-o functrace` | Reset to parent's inherited values |
| `RETURN` | Only with `declare -t` or `-o functrace` | Reset to parent's inherited values |
| `ERR` | Only with `-o errtrace` | Reset to parent's inherited values |
| Signal traps | No (unless `trace` attribute) | Reset to parent's inherited values |

### Trap Timing

- When waiting for a foreground command: trap executes after command finishes
- When waiting via `wait` for async command: trap executes immediately after `wait` returns

### Debug Trap and Completion

The `DEBUG` trap fires before every simple command. This can interfere with completion functions if they trigger commands. Completion functions run in the current shell, so `DEBUG` traps are active during their execution.

## PROMPT_COMMAND

```bash
# Single command
PROMPT_COMMAND="history -a"

# Array (each element executed in order)
PROMPT_COMMAND=("history -a" "history -c" "history -r")
```

- Executed in current shell environment before displaying `PS1`
- If set as array, each element is executed sequentially
- If set as string, value is executed as a command
- Useful for: updating window title, syncing history, setting dynamic variables

## Prompt Variables

| Variable | Purpose | Default |
|----------|---------|---------|
| `PS0` | Displayed after reading command, before execution | |
| `PS1` | Primary prompt | `\s-\v\$` |
| `PS2` | Continuation prompt | `> ` |
| `PS3` | Prompt for `select` command | `#? ` |
| `PS4` | Trace prompt (with `set -x`) | `+ ` |

### PS1 Escape Sequences

| Sequence | Meaning |
|----------|---------|
| `\a` | Bell |
| `\d` | Date ("weekday month date") |
| `\D{format}` | `strftime` format |
| `\e` | Escape character |
| `\h` | Hostname (short) |
| `\H` | Hostname (full) |
| `\j` | Number of jobs |
| `\l` | Terminal device name |
| `\n` | Newline |
| `\r` | Carriage return |
| `\s` | Shell name |
| `\t` | Time (24-hour HH:MM:SS) |
| `\T` | Time (12-hour HH:MM:SS) |
| `\@` | Time (12-hour am/pm) |
| `\A` | Time (24-hour HH:MM) |
| `\u` | Username |
| `\v` | Bash version (short) |
| `\V` | Bash version (full) |
| `\w` | Working directory (with `~` abbreviation) |
| `\W` | Basename of working directory |
| `\!` | History number |
| `\#` | Command number |
| `\$` | `#` if UID=0, else `$` |
| `\nnn` | Octal character |
| `\\` | Backslash |
| `\[` | Begin non-printing sequence (for terminal escapes) |
| `\]` | End non-printing sequence |

### Prompt Expansion

After decoding escape sequences, prompt strings undergo: parameter expansion, command substitution, arithmetic expansion, and quote removal (subject to `promptvars` option).

### Common PS1 Patterns

```bash
# User@host:directory$
PS1='\u@\h:\w\$ '

# With git branch (via PROMPT_COMMAND)
PS1='\u@\h:\w$(__git_ps1)\$ '

# Colored prompt
PS1='\[\e[1;32m\]\u@\h\[\e[0m\]:\[\e[1;34m\]\w\[\e[0m\]\$ '

# With exit code indicator
PS1='$? \u@\h:\w\$ '
```

## Terminal Interaction

### Foreground/Background Process Groups

- **Foreground**: PGID equals terminal's PGID; receives keyboard signals (`SIGINT`, `SIGTSTP`)
- **Background**: Different PGID; immune to keyboard signals
- Background processes attempting terminal read receive `SIGTTIN`
- Background processes attempting terminal write receive `SIGTTOU` (if `stty tostop`)

### How Bash Arranges Process Groups

1. Shell creates a new process group for each pipeline
2. Sets the pipeline's PGID using `setpgid()`
3. Places the pipeline in the foreground with `tcsetpgrp()` (if not background)
4. The terminal driver delivers keyboard signals to the foreground group

### Bracketed Paste

Bash 5.1+ supports bracketed paste mode (enabled by default in Readline):

- Terminal wraps pasted text in `\e[200~` ... `\e[201~` sequences
- Readline treats the entire paste as a single unit
- Prevents accidental execution of pasted commands
- Controlled by `enable-bracketed-paste` Readline variable

## References

- [GNU Bash Manual: Executing Commands](https://www.gnu.org/software/bash/manual/html_node/Executing-Commands.html)
- [GNU Bash Manual: Command Search and Execution](https://www.gnu.org/software/bash/manual/html_node/Command-Search-and-Execution.html)
- [GNU Bash Manual: Command Execution Environment](https://www.gnu.org/software/bash/manual/html_node/Command-Execution-Environment.html)
- [GNU Bash Manual: Signals](https://www.gnu.org/software/bash/manual/html_node/Signals.html)
- [GNU Bash Manual: Job Control](https://www.gnu.org/software/bash/manual/html_node/Job-Control.html)
- [GNU Bash Manual: Pipelines](https://www.gnu.org/software/bash/manual/html_node/Pipelines.html)
- [GNU Bash Manual: Controlling the Prompt](https://www.gnu.org/software/bash/manual/html_node/Controlling-the-Prompt.html)
- [GNU Bash Manual: Bourne Shell Builtins (trap)](https://www.gnu.org/software/bash/manual/html_node/Bourne-Shell-Builtins.html)

## Related Skills

- **references/completion.md** — how completion functions execute (current shell vs subshell)
- **references/startup.md** — when traps and PROMPT_COMMAND are set up
