# Tcsh Execution Model

In-depth reference for how tcsh executes commands, manages processes, handles job control, signals, and I/O redirection.

## Command Search and Execution

### Command Search Order

When tcsh encounters a command, it searches in this order:

1. **Shell builtins** — executed within the shell process
2. **PATH lookup** — searches directories in the `path` variable
3. **Hash table** — uses cached locations (unless `-f` flag, `unhash`, or command contains `/`)

If the command name contains a `/`, tcsh executes it directly as a program (no PATH search).

### Builtin Command Execution

Builtin commands are executed within the shell process. If any component of a pipeline **except the last** is a builtin, the pipeline is executed in a subshell. Parenthesized commands are always executed in a subshell.

### External Command Execution

For external commands, tcsh:

1. Forks a child process
2. Sets up I/O redirections
3. Calls `execve(2)` to replace the child with the command
4. Waits for the child to complete (foreground) or continues (background)

### Script Execution

If a file has execute permissions but is not in executable system format, tcsh spawns a new shell to read it. The `#!interpreter arg ...` convention is supported for specifying the interpreter.

## Simple Commands, Pipelines, and Sequences

### Simple Command

A sequence of words; the first word specifies the command to execute. Remaining words are arguments.

### Pipeline

Commands joined by `|`:

```tcsh
cmd1 | cmd2 | cmd3
```

Output of each command is connected to input of the next. The exit status is that of the last command.

`|&` redirects both standard output and diagnostic (stderr) output through the pipe:

```tcsh
cmd1 |& cmd2
```

### Sequences

| Operator | Meaning |
|----------|---------|
| `;` | Sequential execution (cmd1; cmd2) |
| `&&` | Execute cmd2 only if cmd1 succeeds |
| `\|\|` | Execute cmd2 only if cmd1 fails |
| `&` | Execute in background |

### Grouping

Parentheses group commands into a subshell:

```tcsh
(cmd1; cmd2) | cmd3
```

## I/O Redirection

### Basic Redirections

| Syntax | Description |
|--------|-------------|
| `< name` | Open file as standard input |
| `> name` | Redirect standard output to file |
| `>! name` | Force overwrite (bypass `noclobber`) |
| `>>& name` | Append stdout and stderr to file |
| `>>&! name` | Force append (bypass `noclobber`) |
| `<< word` | Here-document: read shell input up to line matching `word` |
| `>& name` | Redirect stdout and stderr to file |
| `>&! name` | Force redirect (bypass `noclobber`) |

### The `noclobber` Variable

When set, prevents accidental file destruction:

```tcsh
set noclobber
```

With `noclobber` set:

- `> file` fails if `file` already exists
- `>! file` forces overwrite
- `>> file` fails if `file` does not exist
- `>>! file` forces creation

### Redirection and Completion

Redirections affect completion context. When tcsh detects a redirection operator, the next word should be completed as a filename. The `tenematch()` function in `tw.parse.c` tracks redirection context to switch to filename completion.

## Job Control

### Job Table

The shell associates a job with each pipeline. Jobs are tracked in a table and assigned small integer numbers.

### Starting Background Jobs

Append `&` to a command:

```tcsh
make all &
[1] 12345
```

The shell prints the job number and process ID.

### Job References

| Reference | Meaning |
|-----------|---------|
| `%n` | Job number n |
| `%string` | Job whose name begins with `string` |
| `%?string` | Job whose name contains `string` |
| `%+`, `%%`, `%` | Current job |
| `%-` | Previous job |

### Foreground and Background

| Command | Description |
|---------|-------------|
| `fg [%job]` | Bring job to foreground |
| `bg [%job]` | Put job in background |
| `jobs [-l]` | List active jobs (with PIDs if `-l`) |
| `wait` | Wait for all background jobs |

### Suspending Jobs

| Key/Command | Signal | Behavior |
|-------------|--------|----------|
| `^Z` | SIGSTOP | Immediately suspends the current foreground job |
| `^]` | Delayed suspend | Generates STOP when the program reads input |
| `suspend` | — | Stops the shell itself (useful after `su`) |
| `stop %job` | — | Stops a background job |

### The `listjobs` Variable

When set, the shell lists all jobs whenever a job is suspended:

```tcsh
set listjobs
```

### The `notify` Variable and Builtin

By default, the shell reports job status changes just before printing the prompt.

- `notify` variable: enables asynchronous job completion announcements
- `notify %job` builtin: marks a specific job for immediate status reporting

### The `nostat` Variable

List of directories to skip during `stat(2)` calls in completion:

```tcsh
set nostat = (/nfs /tmp)
```

This can speed up completion when network filesystems are slow.

## Signals

### Default Signal Handling

| Signal | Default Behavior |
|--------|-----------------|
| SIGINT (interrupt) | Terminate foreground command; in scripts, controlled by `onintr` |
| SIGQUIT | Ignored unless shell started with `-q` |
| SIGTERM | Caught by login shells; inherited by non-login shells |
| SIGHUP | Shell exits; children receive SIGHUP (controlled by `hup`) |
| SIGTSTP | `^Z` sends to foreground job |
| SIGTTIN | Background job stops if it tries to read from terminal |
| SIGTTOU | Background job stops if `stty tostop` is set and it tries to write |

### The `onintr` Builtin

Controls interrupt handling in scripts:

```tcsh
onintr          # Restore default behavior
onintr -        # Ignore all interrupts
onintr label    # Execute 'goto label' on interrupt
```

### The `hup` Builtin

Controls whether the shell sends SIGHUP to children on exit:

```tcsh
hup [command]   # Run command; send SIGHUP to children on exit
```

Without arguments, non-interactive shells will send SIGHUP to children on exit.

### The `nohup` Builtin

```tcsh
nohup [command] # Run command immune to hangup signals
```

Without arguments, non-interactive shells ignore hangup signals for the remainder of the script.

## Control Flow

### `if` / `else if` / `else` / `endif`

```tcsh
if (expr) command

if (expr) then
  ...
else if (expr) then
  ...
else
  ...
endif
```

The simple form (`if (expr) command`) must be a single simple command.

### `foreach` / `end`

```tcsh
foreach name (wordlist)
  ...
end
```

Sets `name` to each member of `wordlist` in turn. `break` and `continue` control loop flow.

### `while` / `end`

```tcsh
while (expr)
  ...
end
```

Executes commands while `expr` evaluates non-zero. `break` and `continue` work here too.

### `switch` / `case` / `breaksw` / `endsw`

```tcsh
switch (string)
case pattern1:
  ...
  breaksw
case pattern2:
  ...
  breaksw
default:
  ...
  breaksw
endsw
```

Each case label is matched against `string` using glob pattern matching. `breaksw` continues after `endsw`. Without `breaksw`, execution falls through to the next case.

### `goto`

```tcsh
goto label
```

Rewinds input and searches for `label:`. Continues execution after that line.

### `repeat`

```tcsh
repeat count command
```

Executes `command` `count` times. I/O redirections occur only once.

## Exit Status

The `status` variable contains the exit status of the last command:

```tcsh
echo $status
```

Unlike bash's `$?`, tcsh uses `$status`. There is no `$?` in tcsh.

### Exit Status in Completions

External completion commands can use exit status to signal errors, but tcsh does not interpret the exit status of completion commands. The output is used regardless of exit status.

## The `source` Builtin

```tcsh
source [-h] name [args ...]
```

Reads and executes commands from `name` in the current shell. Commands are not placed in the history list unless `-h` is used. Arguments are available in `argv`.

This is the tcsh equivalent of bash's `.` (dot) builtin.

## The `exec` Builtin

```tcsh
exec command
```

Replaces the current shell with `command`. No new process is created; the shell's process is used directly.

## The `exit` Builtin

```tcsh
exit [expr]
```

Exits the shell with the value of `expr` (or status of last command if omitted). For login shells, the logout files are executed.

## The `limit` / `unlimit` Builtins

```tcsh
limit [resource [maximum-use]]
unlimit [-hf] [resource]
```

Controls resource limits. Without arguments, `limit` prints all current limits. `unlimited` removes a limit.

Common resources: `cputime`, `filesize`, `datasize`, `stacksize`, `coredumpsize`, `memoryuse`, `descriptors`.

## The `nice` Builtin

```tcsh
nice [+number] [command]
```

Adjusts scheduling priority. Default increment is 4. Higher numbers mean lower priority. Superusers can use negative numbers to increase priority.

## The `time` Builtin

```tcsh
time [command]
```

Executes `command` and prints a time summary. Without a command, prints a summary for the current shell. The `time` shell variable controls the format and automatic timing threshold.

## The `sched` Builtin

```tcsh
sched [+]hh:mm command
sched
sched -n
```

Schedules a command to execute at a specified time. Without arguments, lists scheduled items. `-n` removes item n.

## Autologout

The `autologout` variable controls automatic logout:

```tcsh
set autologout = 60       # Logout after 60 minutes
set autologout = (60 30)  # Logout after 60 min, lock after 30 min
set autologout = 0        # Disable
```

Default: 60 minutes for login and superuser shells. Disabled if `DISPLAY` is set, the tty is a pty, or the shell was not compiled with the feature.

## References

- tcsh(1) man page — execution, job control, signals, builtins
- `sh.c` in the tcsh source — main shell loop
- `sh.exec.c` in the tcsh source — command execution
- `sh.proc.c` in the tcsh source — process and job control

## Related Skills

- [references/completion.md](references/completion.md) — how completions interact with the execution model
- [references/startup-config.md](references/startup-config.md) — startup file execution order
- [references/quoting-expansion.md](references/quoting-expansion.md) — how expansion transforms commands before execution
