# Nushell Execution Model and Pipeline

In-depth reference for nushell's execution model — the two-stage parse/evaluate architecture, pipeline data flow, external command invocation, and how the module system works.

## Two-Stage Execution Model

Nushell follows a strict two-stage execution model that differs fundamentally from traditional shells:

### Stage 1: Parsing (Static Analysis)

The parser processes the **entire** source code upfront, similar to compilation in static languages:

- All code that will be evaluated must be known and available
- The AST (Abstract Syntax Tree) is constructed
- Type checking occurs
- Parse-time constant evaluation happens
- Syntax errors are detected

### Stage 2: Evaluation (Runtime)

Only after successful parsing does evaluation begin:

- The entire source is evaluated sequentially
- Line 1 is evaluated, then Line 2, and so on

**Key limitation**: Nushell cannot support `eval` functionality since all code must be known at parse-time. This is why completions work well — the parser can determine context from the AST.

## PipelineData: The Core Abstraction

PipelineData is the foundational abstraction for how data flows through command pipelines. It has four variants:

| Variant | Purpose | Contains |
|---------|---------|----------|
| `Empty` | No data to pass | Nothing |
| `Value` | Single value | `Value` + optional `PipelineMetadata` |
| `ListStream` | Stream of values | `ListStream` + optional `PipelineMetadata` |
| `ExternalStream` | External process output | stdout, stderr, exit_code streams |

## Three-Part Pipeline Structure

Every pipeline has three parts:

```
input | filter | output
```

1. **Input (Source/Producer)**: Creates or loads data (e.g., `open`, `ls`, `curl`)
2. **Filter**: Transforms data (e.g., `where`, `update`, `each`, `get`)
3. **Output (Sink)**: Does something with final data (e.g., `save`, `table`, `find`)

## The `$in` Variable

The `$in` variable holds the current pipeline input:

| Rule | Behavior |
|------|----------|
| First position in a pipeline in a closure | `$in` refers to the pipeline input to that closure |
| Elsewhere in a pipeline | `$in` refers to the previous expression's result |
| No input | `$in` is `null` |
| Multi-statement lines (semicolons) | `$in` cannot capture previous statement results |

```nu
date now                    # 1: today
| $in + 1day                # 2: tomorrow
| format date '%F'          # 3: Format as YYYY-MM-DD
| $'($in) Report'           # 4: Format the directory name
| mkdir $in                 # 5: Create the directory
```

## Command Execution Flow

### Command Search Order

1. **Shell functions / custom commands** — executed in current scope
2. **Shell builtins** — invoked within the nushell process
3. **External commands** — searched in `$env.PATH`, executed in separate process

If the command name contains a `/`, nushell executes it directly as a program (no PATH search).

### External Command Invocation

External commands are invoked using the caret prefix (`^`):

```nu
^external_command arg1 arg2
```

Without `^`, nushell first checks if the command is a built-in or custom command.

### Data Flow with External Commands

| Direction | Behavior |
|-----------|----------|
| `internal \| external` | Data converted to string, sent to external's stdin |
| `external \| internal` | Output converted to UTF-8 text stream; if unsuccessful, binary stream |
| `external \| external` | Standard Unix behavior — stdout connects to stdin |

### ExternalStream Handling

The `ExternalStream` variant captures:
- Raw byte streams from stdout and stderr
- Exit code tracking
- Binary vs. text mode detection

### Structured Data Formatting

When piping to external commands, nushell may render structured data with borders. To avoid this, explicitly convert data:

```nu
# Wrong — table borders leak into external command
ls /usr/share/nvim/runtime/ | get name | ^grep tutor | ^ls -la $in

# Right — convert to plain text first
ls /usr/share/nvim/runtime/ | get name | to text | ^grep tutor | ^ls -la $in
```

## Single Return Value Per Expression

An expression can only return a **single value** — the **last** value in the expression:

```nu
def latest-file [] {
    echo "Returning the last file"   # This value is discarded
    ls | sort-by modified | last     # Only this is returned
}
```

## Input/Output Type Declarations

Commands can declare their input and output types:

```nu
def my-filter []: nothing -> list { }
def my-filter []: [
  nothing -> list
  range -> list
] { }
```

Check with `help <command>`:

```nu
help ls
# => Input/output types:
# =>   ╭───┬─────────┬────────╮
# =>   │ # │  input  │ output │
# =>   ├───┼─────────┼────────┤
# =>   │ 0 │ nothing │ table  │
# =>   ╰───┴─────────┴────────╯
```

## Module System

### Module Forms

1. **File-form**: A file named `<module_name>.nu`
2. **Directory-form**: A directory with a `mod.nu` file inside (`<module_name>/mod.nu`)

### Export Types

| Export | Syntax | Description |
|--------|--------|-------------|
| Command | `export def` | Custom commands |
| Alias | `export alias` | Aliases |
| Constant | `export const` | Constants |
| Extern | `export extern` | External command definitions |
| Submodule | `export module` | Nested modules |
| Re-export | `export use` | Imported symbols from other modules |
| Environment | `export-env` | Environment variables |

### The `main` Export

A command named `main` takes on the name of the module when imported:

```nu
# inc.nu
export def main []: int -> int { $in + 1 }
```

```nu
use inc.nu
inc 5  # calls inc.nu's main
```

### Submodules

```nu
# Via export module (keeps submodule namespace)
export module ./increment.nu

# Via export use (flattens into parent namespace)
export use ./increment.nu
```

### Module Search Path

Modules are searched in order:

1. Current directory (relative paths)
2. `$NU_LIB_DIRS` constant
3. `$env.NU_LIB_DIRS` (deprecated)

### Import Patterns

```nu
use <module>                    # Import entire module
use <module> *                  # Import all definitions
use <module> [def1, def2]       # Import specific definitions
use <module> ['sub command']    # Import subcommands
```

### Hiding Imports

```nu
use std/assert
hide assert          # Remove imported definition
hide assert main     # Remove specific definition
```

### Environment from Modules

`export-env` runs only when the `use` call is **evaluated**, not just parsed:

```nu
export-env {
    $env.NU_MODULES_DIR = ($nu.default-config-dir | path join "scripts")
}
```

## How Execution Affects Completions

### Static Parsing Enables Completions

Since nushell parses the entire source code upfront, completions can be context-aware:

1. `NuCompleter` bridges Reedline's completion API with nushell's engine state
2. Uses `EngineState` and `Stack` to provide context-aware suggestions
3. Respects command signatures and their defined completions

### External Command Completion

- Commands not known to nushell (no `extern` definition) trigger the external completer
- The `^` prefix forces external command lookup but is not matched against `extern` definitions
- PATH-based command completion is handled by `CommandCompletion`

### Pipeline Context

The completion system uses the AST to determine:
- Whether the cursor is in command position, argument position, or flag position
- Which command's signature to use for flag/argument completion
- Whether the expression is inside a closure, subexpression, or block

## References

- [Nushell Book: Thinking in Nu](https://www.nushell.sh/book/thinking_in_nu.html)
- [Nushell Book: Modules](https://www.nushell.sh/book/modules.html)
- [DeepWiki: PipelineData](https://deepwiki.com/nushell/nushell/3.2-pipelinedata)

## Related Skills

- [references/completion.md](references/completion.md) — How the completion system uses the AST
- [references/types.md](references/types.md) — Type system and SyntaxShape
- [references/externs.md](references/externs.md) — Extern command definitions
