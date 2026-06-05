# Oil Shell Execution Model

In-depth reference for how Oil shell (OSH/YSH) executes commands, manages processes, handles signals and traps, and how it differs from bash.

## Command Search and Execution

OSH searches for commands in this order:

1. **Shell functions** — executed in current shell context
2. **Shell builtins** — invoked as builtin within the shell process
3. **PATH lookup** — searches `$PATH` directories, uses hash table for caching
4. **External commands** — fork/exec in separate process

If the command name contains a `/`, OSH executes it directly as a program (no PATH search).

## Simple Command Expansion Order

When executing a simple command, OSH performs:

1. **Variable assignments and redirections** — saved for later processing
2. **Word expansion** — tilde, parameter, command substitution, arithmetic, quote removal
3. **Redirections** — performed
4. **Variable assignment values** — undergo tilde, parameter, command substitution, arithmetic, quote removal

If no command name remains after expansion, variable assignments affect the current shell. Otherwise, they only affect the executed command's environment.

## Shell Constructs That Start Processes

| Construct | Example | Notes |
|-----------|---------|-------|
| Simple command | `ls /tmp` | External commands start a process |
| Pipeline | `myproc \| wc -l` | Each part starts a process (affected by `lastpipe`, `pipefail`) |
| Command substitution | `d=$(date)` | Process for the command |
| Process substitution | `<(sort left.txt)` | Process for the sub-expression |
| Async (fork) | `sleep 2 &` | Background process |
| Explicit subshell | `( echo hi )` or `forkwait` | Rarely needed; prefer `pushd`/`popd` or `cd { }` in YSH |

## Pipelines: Key Difference from Bash

### lastpipe Semantics

OSH uses `lastpipe` semantics **by default** (like zsh, unlike bash):

```bash
# Bash: read runs in subshell, variable is lost
bash -c 'echo hi | read x; echo x=$x'
x=

# OSH: last part runs in current shell
osh -c 'echo hi | read x; echo x=$x'
x=hi
```

In bash, `shopt -s lastpipe` enables this, but bash **ignores** it in interactive shells. OSH always uses lastpipe semantics.

### Pipeline Suspension Incompatibility

Because the last part of a pipeline may run in the current shell process, **pipelines cannot be suspended with Ctrl-Z** in OSH:

| Shell | Pipeline Suspension |
|-------|---------------------|
| bash | Possible (ignores `lastpipe` interactively) |
| zsh | Doesn't allow suspending some pipelines |
| OSH | Not possible (due to `lastpipe` default) |

### PIPESTATUS

`${PIPESTATUS[@]}` is **only set after actual pipelines** in OSH, not for every command (unlike bash):

```bash
# OSH: PIPESTATUS only set for pipelines
ls | grep foo | wc -l
echo "${PIPESTATUS[@]}"    # 0 1 0

# Single command — no PIPESTATUS
ls
echo "${PIPESTATUS[@]}"    # Not set (unlike bash which sets it)
```

### Pipeline Options

| Option | Description |
|--------|-------------|
| `lastpipe` | Run last pipeline part in current shell (default in OSH) |
| `pipefail` | Pipeline returns failure if any part fails |
| `sigpipe_status_ok` | Status 141 (SIGPIPE) → 0 in pipelines (useful for `yes \| head`) |

## Subshells

### Explicit Subshell

```bash
( echo hi )          # Traditional syntax
forkwait { echo hi } # YSH syntax
```

Explicit subshells are **rarely needed**. Prefer `pushd`/`popd`, or `cd { }` in YSH.

### Implicit Subshells

Subshells can occur without explicit syntax:

- Pipelines (with certain options)
- Process substitution
- Command substitution

### Nested Subshell Syntax

`$((` always starts arithmetic, not a subshell. Use a space:

```bash
$( (cd / && ls) )    # Subshell in command substitution
$({ cd / && ls; })   # Grouping (preferred)
```

`((` always starts arithmetic, not grouping. Use braces:

```bash
if { test -f a || test -f b; } && grep foo c
```

## Process Optimizations

Oil implements `noforklast` optimization — avoiding unnecessary forks when the last command in a pipeline or subshell is a simple external command. However, forking is required in certain cases:

| Case | Why Forking Required |
|------|---------------------|
| Job control | Restoring process state after shell runs |
| Traps | Not run in optimized process (issue #1853) |
| `set -o pipefail` | Bug with optimization |
| `! false` expressions | Negation requires forking |
| `verbose_errexit` | YSH option for detailed errors |
| Crash dump | Needs process isolation |
| Stats/tracing | Counting exit codes |

## Signals

### Signal Handling

| Signal | Typical Source | OSH Behavior |
|--------|---------------|---------------|
| `SIGINT` | Ctrl-C | Interrupts and returns to prompt; `trap` builtin |
| `SIGTERM` | `kill` command | `trap` builtin |
| `SIGQUIT` | Ctrl-\\ | Default behavior |
| `SIGTTIN` | Job control | Background process reading from terminal |
| `SIGTTOU` | Job control | Background process writing to terminal |
| `SIGWINCH` | Terminal resize | Window size change |

## Traps

| Trap | Description |
|------|-------------|
| `DEBUG` | Runs before "leaf" commands (`echo hi`, `a=b`, `[[ x -eq y ]]`, `(( a = 42 ))`), but **not** before compound commands |
| `ERR` | Runs when a command fails (where `set -o errexit` would exit). Enable in functions/subshells with `set -o errtrace` (`set -E`) |
| `EXIT` | Runs on shell exit |
| `RETURN` | Runs on function return |

### DEBUG Trap + Pipelines + Job Control

There is a fundamental incompatibility between the DEBUG trap, `lastpipe`, and job control:

1. OSH emulates bash's `DEBUG` trap (runs before "leaf" commands)
2. Running the trap before the last pipeline part (with `lastpipe`) exacerbates a race condition
3. In `echo hi | cat`, `echo hi` could finish before `cat` starts, so `cat` can't join the process group

**OSH solution**: Disables the `DEBUG` trap for the last part of a pipeline **only when job control is enabled**. This doesn't affect debugging batch programs.

## Job Control

### Builtins

| Builtin | Description |
|---------|-------------|
| `fg` | Bring background job to foreground |
| `bg` | Run job in background |
| `wait` | Wait for background jobs |
| `jobs` | List active jobs |
| `disown` | Remove job from active list |

### Job Control and lastpipe

Job control requires putting a pipeline in a process group so it can be suspended and cancelled all at once. But `lastpipe` means the last part runs in the current shell, which can't be in a separate process group.

| Shell | Behavior |
|-------|----------|
| bash | Ignores `shopt -s lastpipe` in job control shells |
| zsh | Doesn't allow suspending some pipelines |
| OSH | Disables DEBUG trap for last pipeline part when job control is on |

## Error Handling

### Exit Codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Runtime error (file can't be opened, no job to foreground) |
| `2` | Parse error (didn't attempt to do anything) |
| `3` | Expression error (YSH: divide by zero, index out of bounds) |
| `126` | Permission denied (POSIX) |
| `127` | Command not found (POSIX) |

### Boolean Exit Codes

| Code | Meaning |
|------|---------|
| `0` | True |
| `1` | False |
| `2` | Error |

### strict:all Option Group

`shopt -s strict:all` disallows problematic shell constructs:

| Option | Description |
|--------|-------------|
| `strict_arith` | Strings that don't look like integers cause fatal errors |
| `strict_argv` | Empty `argv` arrays disallowed |
| `strict_array` | No implicit string↔array conversions |
| `strict_control_flow` | `break`/`continue` outside loops are fatal |
| `strict_errexit` | Warns when errors would be lost |
| `strict_word_eval` | More word evaluation errors are fatal |
| `strict_nameref` | Invalid variable names in namerefs are fatal |
| `strict_tilde` | Failed tilde expansions cause errors |
| `strict_glob` | Parse the sublanguage more strictly |

### YSH Error Handling

YSH provides a cleaner error handling model:

```bash
try {
  echo hi
  badcmd
}
if failed {
  echo "error: $_error"
}
```

The `_error` register contains the last error. `io.captureAll()` gets stdout, stderr, and status all at once — other shells can't do this.

## Xtrace (Tracing)

### Traditional `set -x`

```bash
$ sh -x -c 'echo 1; echo 2'
+ echo 1
1
+ echo 2
2
```

### Rich Xtrace (YSH)

With `shopt -s xtrace_rich` (included in `ysh:upgrade`):

```
. builtin set '+e'
> pipeline
  | part 103
    . 103 exec ls
  | part 104
    . 104 exec grep OOPS
  | command 105: wc -l
  ; process 103: status 0
  ; process 104: status 1
  ; process 105: status 0
< pipeline
. builtin echo end
```

| Symbol | Meaning |
|--------|---------|
| `.` | Builtin or exec |
| `>` `<` | Synchronous constructs (pipeline, proc, eval, source, wait, trap handlers) |
| `\|` | Process start |
| `;` | Process end (with exit status) |

### SHX Variables (YSH PS4)

Default YSH PS4: `${SHX_indent}${SHX_punct}${SHX_pid_str} `

| Variable | Purpose |
|----------|---------|
| `SHX_indent` | Indentation for synchronous constructs |
| `SHX_punct` | Prefix character (`.`, `>`, `<`, `|`, `;`) |
| `SHX_pid_str` | PID annotation for child processes |
| `SHX_descriptor` | Alias for `BASH_XTRACEFD` |

## Interpreter State

### Memory Model

- Shell has a **stack but no heap**
- Stack stores: local function variables, `$@` argument array
- Environment variables become global variables with `export` flag
- Functions and variables have **separate namespaces**

### Dynamic Scope

OSH uses dynamic scope (like bash). YSH limits this capability — procs have truly local variables (like Python/JavaScript), no dynamic scope rule.

### Value Types

| Type | Description |
|------|-------------|
| `Str` | String value |
| `BashArray` / `List` | Indexed array (OSH / YSH) |
| `AssocArray` / `Dict` | Associative array (OSH / YSH) |
| `Undef` | Unset variable |

## Differences from Bash: Summary

| Feature | OSH | Bash |
|---------|-----|------|
| Pipeline last part | Runs in current shell by default | Runs in subshell by default |
| `echo x \| read` result | Variable is preserved | Variable is lost |
| Pipeline suspension | Not possible (due to `lastpipe`) | Possible (ignores `lastpipe` interactively) |
| `PIPESTATUS` | Only set after actual pipeline | Set on every command |
| `declare -i` | No-op | Supported |
| DEBUG trap + lastpipe | Disabled for last pipeline part when job control on | N/A |
| Type model | Values tagged with types | Locations tagged with types |
| Xtrace | Rich hierarchical tracing available | Flat `+` prefix only |
| Error handling | `strict:all` + YSH `try`/`failed` | `set -e` only |

## References

- [Oil Shell Process Model](https://oils.pub/release/latest/doc/process-model.html)
- [Oil Shell Interpreter State](https://oils.pub/release/latest/doc/interpreter-state.html)
- [Known Differences Between OSH and Other Shells](https://oils.pub/release/latest/doc/known-differences.html)
- [OSH Quirks](https://oils.pub/release/latest/doc/quirks.html)
- [Tracing Execution in Oils (xtrace)](https://oils.pub/release/latest/doc/xtrace.html)
- [YSH Error Handling](https://oils.pub/release/latest/doc/ysh-error.html)

## Related Skills

- **bash skill → references/execution.md** — bash execution model and signals
