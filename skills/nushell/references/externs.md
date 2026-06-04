# Nushell Extern Commands

In-depth reference for nushell's `extern` command definitions — how to define external commands with completion signatures, attach completers, define subcommands, and organize externs in modules.

## The `extern` Keyword

`extern` defines a command signature for an external (non-nushell) command. It tells nushell how to parse, type-check, and complete the command's arguments:

```nu
export extern ssh [
  destination?: string@complete_none  # Destination Host
  -p: int                             # Destination Port
  -i: string@complete_ssh_identity    # Identity File
]
```

This provides:
1. **Parse-time type checking** — non-matching types result in errors
2. **Syntax highlighting** — based on argument shapes
3. **Completion suggestions** — flags and arguments from the defined signature
4. **Custom completers** — attached via `@` syntax

## Basic Syntax

```nu
extern <command-name> [
  <parameter-definitions>
]
```

### Parameter Types

| Syntax | Description | Example |
|--------|-------------|---------|
| `name: type` | Required positional | `destination: string` |
| `name?: type` | Optional positional | `destination?: string` |
| `...name: type` | Rest parameter (variadic) | `...pathspecs: path` |
| `--name: type` | Long flag with value | `--port: int` |
| `--name` | Boolean flag (no type) | `--verbose` |
| `-s: type` | Short flag with value | `-p: int` |
| `-s` | Boolean short flag | `-v` |

### Flag Definitions

```nu
export extern git [
  --verbose(-v)     # Long flag with short alias
  --port: int       # Flag that takes a value
  -q                # Short-only boolean flag
]
```

- A space is required before `#` for inline documentation
- Flags without explicit type default to `bool` (present/not present)
- The `bool` type cannot be used as a flag annotation — that is the same as the existence of the flag

## Attaching Completers

### Positional Argument Completions

```nu
def complete_ssh_identity [] {
    ls ~/.ssh/id_*
    | where {|f|
        ($f.name | path parse | get extension) != "pub"
      }
    | get name
}

export extern ssh [
  destination?: string@complete_none
  -i: string@complete_ssh_identity
]
```

The `@` syntax attaches a completer function to the parameter. The completer is invoked when the user presses TAB at that argument position.

### Flag Argument Completions

```nu
export extern ssh [
  -i: string@complete_ssh_identity    # Flag with custom completer
]
```

Flag arguments support the same `@` syntax as positional arguments.

### Subcommand Completions

Subcommands are defined using quoted command strings:

```nu
export extern "git add" [
  ...pathspecs: path
]

export extern "git push" [
  remote?: string@"nu-complete git remotes"
  refspec?: string@"nu-complete git branches"
]
```

The quoted string `"git add"` tells nushell this is a subcommand of `git`.

## Module-Based Externs

Extern commands are typically defined inside a module that also contains their completion functions:

```nu
module "ssh extern" {
    def complete_none [] { [] }

    def complete_ssh_identity [] {
        ls ~/.ssh/id_*
        | where {|f|
            ($f.name | path parse | get extension) != "pub"
          }
        | get name
    }

    export extern ssh [
        destination?: string@complete_none
        -p: int
        -i: string@complete_ssh_identity
    ]
}
use "ssh extern" ssh
```

Only the `extern` is exported — the completion functions remain private to the module.

### Directory-Form Modules

A directory with a `mod.nu` file:

```
ssh_extern/
  mod.nu
```

```nu
# ssh_extern/mod.nu
def complete_none [] { [] }

def complete_ssh_identity [] {
    ls ~/.ssh/id_* | where {|f| ($f.name | path parse | get extension) != "pub" } | get name
}

export extern ssh [
    destination?: string@complete_none
    -p: int
    -i: string@complete_ssh_identity
]
```

## Completer Return Types

### Simple List

```nu
def animals [] { ["cat", "dog", "eel"] }
```

### Records with Descriptions and Styles

```nu
def my_commits [] {
    [
        { value: "5c2464", description: "Add .gitignore", style: red },
        { value: "f3a377", description: "Initial commit", style: { fg: green, bg: "#66078c", attr: ub } }
    ]
}
```

### With Options Override

```nu
def animals [] {
    {
        options: {
            case_sensitive: false,
            completion_algorithm: substring,
            sort: false,
        },
        completions: [cat, rat, bat]
    }
}
```

### Context-Aware

```nu
def animal-names [context: string] {
    match ($context | split words | last) {
        cat => ["Missy", "Phoebe"]
        dog => ["Lulu", "Enzo"]
        eel => ["Eww", "Slippy"]
    }
}
```

The completer can also receive cursor position:

```nu
def completer [context: string, position: int] {}
```

### Fallback Behavior

- Return `null` to fall back to nushell's built-in file completion
- Return `[]` (empty list) to suppress completions entirely

## Externs and the External Completer

Externs and the external completer closure work together:

1. **Externs are checked first** — if an extern is defined for the command, nushell uses its signature for flag/argument completion
2. **External completer is fallback** — if no extern is defined, the external completer closure is invoked
3. **Custom completers on externs** — the `@` syntax on extern parameters takes priority over the external completer

This means you can:
- Define externs for commands where you want fine-grained control (specific flags, typed arguments)
- Use the external completer (e.g., carapace) for all other commands

## Limitations

1. **Cannot require flag ordering** — nushell doesn't enforce a particular order of flags and positional arguments
2. **Cannot require `=` syntax** — nushell doesn't support `--flag=value` as a required syntax (space-separated is always valid)
3. **Caret sigil not recognized** — externals called via `^ssh` are not matched against `extern ssh`
4. **Single-dash long flags** — cannot represent `--long` arguments using a single leading hyphen (non-POSIX mode)
5. **No dynamic signatures** — extern signatures are static; they cannot change at runtime

## Externs vs. Custom Commands

| Feature | `extern` | `def` |
|---------|----------|-------|
| Defines body | No | Yes |
| Type checking | Yes | Yes |
| Flag completion | Yes | Yes |
| Argument completion | Yes (via `@`) | Yes (via `@`) |
| Syntax highlighting | Yes | Yes |
| Used for external commands | Yes | No |
| Used for nushell commands | No | Yes |

## Common Patterns

### Git Subcommands

```nu
module "git-completions" {
    def "nu-complete git remotes" [] {
        ^git remote 2>/dev/null | lines | where {|l| ($l | str trim) != ""}
    }

    def "nu-complete git branches" [] {
        ^git branch 2>/dev/null | lines | each {|l| $l | str replace -r '^\*?\s+' '' } | where {|l| ($l | str trim) != ""}
    }

    export extern "git push" [
        remote?: string@"nu-complete git remotes"
        refspec?: string@"nu-complete git branches"
        --force(-f)
        --dry-run(-n)
    ]

    export extern "git checkout" [
        branch?: string@"nu-complete git branches"
        -b                      # Create new branch
    ]
}
use "git-completions"
```

### Suppressing File Completion

When a parameter should not trigger file completion (e.g., an integer argument), attach a completer that returns an empty list:

```nu
def complete_none [] { [] }

export extern myapp [
    count?: int@complete_none    # No file completion for this int arg
]
```

## References

- [Nushell Book: Externs](https://www.nushell.sh/book/externs.html)
- [Nushell Book: Custom Completions](https://www.nushell.sh/book/custom_completions.html)
- [Nushell Book: Modules](https://www.nushell.sh/book/modules.html)

## Related Skills

- [references/completion.md](references/completion.md) — How completers are dispatched
- [references/types.md](references/types.md) — SyntaxShape and type annotations
- [references/configuration.md](references/configuration.md) — Where to place extern definitions
