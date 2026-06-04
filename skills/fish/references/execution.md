# Fish Execution Model

In-depth reference for how fish executes commands, manages processes, handles job control, signals, and the internal threading model.

## Command Search and Execution

Fish searches for commands in this order:

1. **Shell functions** — executed in current shell context (no new process)
2. **Shell builtins** — invoked within the shell process
3. **PATH lookup** — searches `$PATH` directories for external commands
4. **External commands** — fork/exec in separate process

If the command name contains a `/`, fish executes it directly as a program (no PATH search).

## Process Types

Fish distinguishes five process types:

| Type | Description | Execution Method |
|------|-------------|------------------|
| `External` | Regular external commands | `fork()` or `posix_spawn()` |
| `Builtin` | Built-in commands | Internal, in background thread |
| `Function` | Fish script functions | Parser evaluation |
| `BlockNode` | AST statement blocks | AST traversal |
| `Exec` | `exec` builtin | `exec()` without fork |

### Fork vs Posix Spawn

Fish prefers `posix_spawn()` when possible for efficiency, but uses `fork()` when:

- Foreground job (needs `tcsetpgrp()` before exec)
- Self-fd redirections (complex FD management)
- `use_posix_spawn()` disabled

### Internal Process Pattern

Builtins and functions that need to write to pipes use the **internal process pattern**:

1. Output captured in memory (`StringOutputStream`)
2. `InternalProc` created with unique ID
3. Background write via `exec_thread_pool()`
4. `mark_exited()` called when write completes

## Execution Flow

```
Command String → Parser → AST → ExecutionContext → exec_job() → Process
```

### ExecutionContext

Central structure maintaining execution state:

- `pstree` — parsed source code and AST
- `cancel_signal` — set when execution should be cancelled (SIGINT/SIGQUIT)
- `pipeline_node` — shared with Parser to track current pipeline
- `block_io` — redirections inherited from containing blocks
- `test_only_suppress_stderr` — internal flag for unit tests

### Statement Execution Methods

| Method | AST Node Type | Description |
|--------|---------------|-------------|
| `populate_job_from_job_node()` | `ast::JobPipeline` | Creates job with processes |
| `run_block_statement()` | `ast::BlockStatement` | Executes for/while/function/begin blocks |
| `run_begin_statement()` | `ast::BraceStatement` | Executes brace groups |
| `run_if_statement()` | `ast::IfStatement` | Evaluates conditions |
| `run_switch_statement()` | `ast::SwitchStatement` | Pattern matching |
| `eval_job_list()` | `ast::JobList` | Sequentially executes jobs |

## Job and Process Hierarchy

```
Job
├── job_id
├── processes: Vec<Process>
├── group: Option<JobGroup>
├── command_str
└── flags
    └── Process
        ├── typ (External/Builtin/Function/BlockNode/Exec)
        ├── argv
        ├── actual_cmd
        ├── pid (external only)
        ├── internal_proc
        └── status
```

### Job Lifecycle

1. **Population** — `populate_job_from_job_node()` traverses AST `JobPipeline` node
2. **Addition** — Job added to parser's job list via `job_add()`
3. **Execution** — `exec_job()` launches all processes
4. **Monitoring** — `job_reap()` checks for completed processes

## Pipelines

```fish
command1 | command2 | command3
```

- Output of each command connected via pipe to input of next
- Each command executed in its **own process** (for external commands)
- `|&` connects stderr to stdout through pipe (same as `2>&1 |`)
- Exit status: last command's exit status (stored in `$status`)
- `$pipestatus` contains exit statuses of all pipe elements

```fish
echo "hello" | cat | tr 'a-z' 'A-Z'
echo $pipestatus  # 0 0 0
```

## Job Control

### Job Control Modes

| Mode | Effect |
|------|--------|
| `JobControl::All` | All jobs get process groups, can be stopped |
| `JobControl::Interactive` | Only interactive jobs get job control |
| `JobControl::None` | No job control |

Set via `status job-control full|interactive|none`.

### Process Groups and Terminal Control

Fish uses process groups for signal distribution and terminal control:

- Only the **foreground process group** can read/write to the terminal
- The `TtyHandoff` utility manages transferring the TTY to child process groups
- Background processes attempting terminal I/O receive `SIGTTIN`/`SIGTTOU`

### Job Control Features

| Feature | Description |
|---------|-------------|
| `^Z` (suspend) | Stops current foreground process, returns to shell |
| `bg %n` | Resume job in background |
| `fg %n` | Bring job to foreground |
| `jobs` | List active jobs |
| `wait` | Wait for job completion |
| `disown` | Remove from jobs table |

### Job Specifiers

| Specifier | Meaning |
|-----------|---------|
| `%n` | Job number n |
| `%%` or `%+` | Current job |
| `%-` | Previous job |
| `%string` | Job whose command starts with `string` |
| `%?string` | Job whose command contains `string` |

## Signals

### Default Signal Handling in Interactive Fish

| Signal | Default Action |
|--------|---------------|
| `SIGINT` | Caught (interrupts current command) |
| `SIGTERM` | Shell exits |
| `SIGHUP` | Shell exits, sends SIGHUP to all jobs |
| `SIGTTIN` | Ignored (with job control) |
| `SIGTTOU` | Ignored (with job control) |
| `SIGTSTP` | Ignored (with job control) |

### Signal Handling for Child Processes

- External commands inherit signal handlers from shell (unless set to ignore)
- Async commands ignore `SIGINT` and `SIGQUIT`
- Command substitutions ignore `SIGTTIN`, `SIGTTOU`, `SIGTSTP`

### Event Handlers

Fish supports event handlers triggered by signals:

```fish
function on_sigusr1 --on-signal SIGUSR1
    echo "Received SIGUSR1"
end
```

### Named Events

| Event | When Triggered |
|-------|---------------|
| `fish_prompt` | Before displaying prompt |
| `fish_preexec` | Before executing a command |
| `fish_postexec` | After executing a command |
| `fish_exit` | When shell exits |
| `fish_cancel` | When command line is cancelled |
| `fish_focus_in` | When terminal gains focus |
| `fish_focus_out` | When terminal loses focus |

```fish
function on_preexec --on-event fish_preexec
    echo "About to execute: $argv[1]"
end

function on_postexec --on-event fish_postexec
    echo "Command exited with status: $status"
end
```

## Threading and Async I/O

### Topic Monitor

The `TopicMonitor` is the core synchronization primitive for cross-thread signaling:

- **Topic**: Enum representing events (`SigHupInt`, `SigChld`, `InternalExit`)
- **GenerationsList**: Current generation counts for all topics
- **BinarySemaphore**: Cross-platform semaphore (uses `sem_init` on Linux, self-pipe elsewhere)

### Thread Pool (iothread)

Global `ThreadPool` for offloading blocking operations:

- History saving
- Autosuggestion path validation
- Universal variable synchronization
- Completion generation for external commands

### FD Monitor

Provides thread-safe watching of multiple file descriptors using `select` or `poll`:

- `FdEventSignaller`: Wrapper around `eventfd` (Linux) or pipe for async-signal-safe signaling
- `FdMonitorItem`: Contains `OwnedFd` and callback

### Universal Variable Synchronization

Universal variables are shared across all fish instances:

1. Stored in `~/.config/fish/fish_variables`
2. `EnvUniversal::sync()` handles reading and merging
3. Uses filesystem-based notification for cross-process sync
4. Atomic writes via `rewrite_via_temporary_file()`
5. Changes propagate to all running fish instances

### Debouncing and Throttling

- **Autosuggestions**: Path validation debounced
- **History Vacuuming**: Reorganized periodically (`VACUUM_FREQUENCY = 25`)
- **Universal Variables**: Syncing throttled

## Environment Stack

The `EnvStack` manages variables with scoping using an internal `EnvStackImpl` protected by a mutex for thread-safe access across scopes (local, global, universal).

### Variable Scope Resolution

When reading a variable, fish searches scopes in this order:

1. Local (innermost block first)
2. Function
3. Global
4. Universal

### What Creates a Subprocess

| Construct | Subprocess? |
|-----------|-------------|
| External command | Yes (fork/exec or posix_spawn) |
| Pipeline elements | Yes (for external commands) |
| Background commands `cmd &` | Yes |
| Shell functions | **No** — current shell |
| `source` / `.` script | **No** — current shell |
| `exec cmd` | **No** — replaces current process |
| Builtins | **No** — current shell |

**Key distinction**: Changes in a subprocess cannot affect the parent's environment. Functions and builtins execute in the current shell context.

## Comparison with Bash's Execution Model

| Feature | Fish | Bash |
|---------|------|------|
| Process creation | `posix_spawn()` preferred | `fork()` always |
| Job control | Three modes (all/interactive/none) | On/off |
| Pipeline execution | Each element in own process | Last element in current shell with `lastpipe` |
| Variable scoping | Local/function/global/universal | Local (function-only)/global/exported |
| Universal variables | Yes (cross-session, file-backed) | No equivalent |
| Signal handlers | Event-based (`--on-signal`) | `trap` builtin |
| Process substitution | `(cmd \| psub)` | `<(cmd)` |
| Thread model | Thread pool for async I/O | Single-threaded |

## References

- [Core Architecture — DeepWiki](https://deepwiki.com/fish-shell/fish-shell/2-core-architecture)
- [Command Execution Engine — DeepWiki](https://deepwiki.com/fish-shell/fish-shell/4.4-command-execution-engine)
- [Job Control and Process Management — DeepWiki](https://deepwiki.com/fish-shell/fish-shell/4.5-job-control-and-process-management)
- [Threading and Async I/O — DeepWiki](https://deepwiki.com/fish-shell/fish-shell/11.2-threading-and-async-io)
- [Fish shell source: src/proc.rs](https://github.com/fish-shell/fish-shell/blob/main/src/proc.rs)
- [Fish shell source: src/parser.rs](https://github.com/fish-shell/fish-shell/blob/main/src/parser.rs)

## Related Skills

- **bash** → `references/execution.md` — bash execution model for comparison
- **fish** → `references/completion.md` — how execution affects completion context
