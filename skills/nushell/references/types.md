# Nushell Type System and SyntaxShape

In-depth reference for nushell's three-tier type system — SyntaxShape (parse-time), Type (semantic), and Value (runtime) — and how types drive the completion system through the `@` syntax and `CompleterWrapper`.

## The Three-Tier Type System

Nushell maintains three distinct representations of types at different stages of execution:

| Representation | Stage | Purpose |
|---------------|-------|---------|
| `SyntaxShape` | Parse-time | Guides the parser on what syntactic form to expect; provides completion hints |
| `Type` | Semantic | Type checking and inference during parsing; lives in the AST |
| `Value` | Runtime | Actual data values during evaluation |

## SyntaxShape — Parse-Time Type Descriptors

`SyntaxShape` describes the syntactic forms the parser should recognize. It extends beyond runtime values to include parser-specific concepts.

### Basic Types

| SyntaxShape | Description | Completion behavior |
|-------------|-------------|---------------------|
| `Any` | Accept any syntactic form | No specific completion |
| `Int` | Integer literal | No specific completion |
| `Float` | Float literal | No specific completion |
| `Number` | Int or Float | No specific completion |
| `String` | String literal or bare word | File completion fallback |
| `Boolean` | `true` or `false` | Boolean completion |
| `Binary` | Binary literal | No specific completion |
| `DateTime` | Date/time literal | No specific completion |
| `Duration` | Duration literal | No specific completion |
| `Filesize` | Filesize literal | No specific completion |

### Structural Types

| SyntaxShape | Description | Completion behavior |
|-------------|-------------|---------------------|
| `List(Box<SyntaxShape>)` | List with inner type | No specific completion |
| `Record(Vec<(String, SyntaxShape)>)` | Record with typed fields | No specific completion |
| `Table(Vec<(String, SyntaxShape)>)` | Table with typed columns | No specific completion |
| `Range` | Range expression | No specific completion |

### Code Types

| SyntaxShape | Description |
|-------------|-------------|
| `Block` | Code block `{ ... }` |
| `Closure(Option<Vec<SyntaxShape>>)` | Closure with optional parameter types |
| `Signature` | Function signature |

### File System Shapes

| SyntaxShape | Description | Completion behavior |
|-------------|-------------|---------------------|
| `Filepath` | File path | FileCompletion |
| `Directory` | Directory path | DirectoryCompletion |
| `GlobPattern` | Glob pattern | FileCompletion |

### Parser-Specific Shapes

| SyntaxShape | Description |
|-------------|-------------|
| `Expression` | General expression |
| `MathExpression` | Math expression |
| `RowCondition` | Row filtering condition |
| `Keyword(Vec<u8>, Box<SyntaxShape>)` | Specific keyword followed by value |
| `OneOf(Vec<SyntaxShape>)` | Try multiple shapes in order |

### CompleterWrapper — The Critical Bridge

```rust
SyntaxShape::CompleterWrapper(Box<SyntaxShape>, DeclId)
```

This is the shape created by the `@` syntax. It wraps:
1. An inner `SyntaxShape` (e.g., `String`) — used for type-checking
2. A `DeclId` — the declaration ID of the completer function

When the parser encounters `string@animals`, it creates `CompleterWrapper(String, animals_decl_id)`. The parser uses the inner shape for type-checking, and the completion system uses the `DeclId` to invoke the completer.

## Type — Semantic Type System

The `Type` enum represents the semantic type system used for type checking and inference during parsing. Each `Expression` in the AST carries a `ty` field.

### Basic Types

| Type | Annotation | Description |
|------|-----------|-------------|
| `Int` | `int` | Whole numbers |
| `Float` | `float` | Numbers with fractional component |
| `String` | `string` | Text |
| `Bool` | `bool` | True or False |
| `Nothing` | `nothing` | Absence of a value |
| `Any` | `any` | Superset of all types |
| `Number` | — | Supertype of Int and Float |
| `Duration` | `duration` | Time passage (e.g., `2min`) |
| `Filesize` | `filesize` | File size (e.g., `64mb`) |
| `Date` | `datetime` | Point in time |
| `Range` | `range` | Value range |
| `Binary` | `binary` | Raw bytes |
| `CellPath` | `cell-path` | Navigation path |
| `Closure` | `closure` | Anonymous function |
| `Block` | — | Code block (syntactic form) |
| `ListStream` | — | Stream of values |

### Structured Types

| Type | Annotation | Description |
|------|-----------|-------------|
| `List` | `list<T>` | Ordered sequence |
| `Record` | `record<key: type, ...>` | Key-value pairs |
| `Table` | `table` | Two-dimensional container (list of records) |

### Subtyping Rules

- `Type::Any` is the **supertype** of all other types
- `Type::Number` is a **supertype** of both `Type::Int` and `Type::Float`
- Lists are **covariant**: `List<T>` is a subtype of `List<U>` if `T` is a subtype of `U`
- Records/Tables use **structural subtyping**: a record with extra fields is a subtype
- `Type::Closure` is compatible with `Type::Block`
- `Type::List(T)` is compatible with `Type::ListStream`
- `OneOf`: A type `T` is a subtype of `OneOf<U1, U2, ...>` if `T` is a subtype of any `Ui`

### Type Utility Methods

| Method | Purpose |
|--------|---------|
| `is_subtype_of()` | Check subtyping relationship |
| `is_numeric()` | Check if type is numeric |
| `is_list()` | Check if type is a list variant |
| `accepts_cell_paths()` | Check if type accepts cell paths |
| `to_shape()` | Convert to SyntaxShape |

## SyntaxShape ↔ Type Conversions

| SyntaxShape | → Type |
|-------------|--------|
| `SyntaxShape::Int` | `Type::Int` |
| `SyntaxShape::String` | `Type::String` |
| `SyntaxShape::Expression` | `Type::Any` |
| `SyntaxShape::Block` | `Type::Block` |
| `SyntaxShape::Closure` | `Type::Closure` |
| `SyntaxShape::RowCondition` | `Type::Bool` |
| `SyntaxShape::List(x)` | `Type::List` (recursive) |
| `SyntaxShape::Filepath` | `Type::String` |
| `SyntaxShape::Directory` | `Type::String` |
| `SyntaxShape::CompleterWrapper(inner, _)` | `inner.to_type()` |

Parser-specific shapes typically map to `Type::Any`.

## Value — Runtime Type Representation

The `Value` enum represents actual data at runtime. Each variant contains the value data and a span for error reporting.

Each `Value` can report its `Type` via `get_type()`:
- For lists, it infers the most specific common type of elements
- If a list contains both `Int` and `Float`, it returns `Type::Number`
- If a list has incompatible types, it returns `Type::Any`
- Lists of records become `Type::Table` with inferred columns

## How Types Drive Completions

### The `@` Syntax

```nu
def animals [] { ["cat", "dog", "eel"] }
def my-command [animal: string@animals] { print $animal }
```

`string@animals` tells nushell:
1. **Type-checking**: The parameter has `SyntaxShape::String` (parsed as `CompleterWrapper(String, animals_decl_id)`)
2. **Completion**: When completing this parameter, invoke the `animals` command

### Type-Driven Built-in Completions

The `SyntaxShape` of a parameter determines which built-in completer is used when no custom completer is attached:

| SyntaxShape | Built-in Completer |
|-------------|-------------------|
| `Filepath` | FileCompletion |
| `Directory` | DirectoryCompletion |
| `String` (no @) | FileCompletion (fallback) |
| `Any` | No specific completion |
| Other | No specific completion |

### Command Signatures and Type Annotations

```nu
# Simple parameter
def my-command [x: int] { }

# Optional parameter with default
def fully [some?: int = 9] { $some }

# Input/Output type declarations
def my-filter []: nothing -> list { }
def my-filter []: [
  nothing -> list
  range -> list
] { }
```

The `:` infix operator connects a parameter name to its type annotation. Flags with explicit type annotations expect an argument of that type. Flags without type annotations default to `bool` (present/not present).

### Context-Aware Completions

Custom completions can be context-aware by receiving the command-line context:

```nu
def animal-names [context: string] {
    match ($context | split words | last) {
        cat => ["Missy", "Phoebe"]
        dog => ["Lulu", "Enzo"]
        eel => ["Eww", "Slippy"]
    }
}

def my-command [
    animal: string@animals
    name: string@animal-names
] { print $"The ($animal) is named ($name)." }
```

The completer receives the full command-line context as a string, enabling dependent completions.

### Completion Options Override

A completer can return a record with `options` to override global settings:

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

### Suppressing Completions

- Return `null` to fall back to nushell's built-in file completion
- Return `[]` (empty list) to suppress completions entirely

## Shape_* Configuration (Syntax Coloring)

Nushell uses `shape_*` settings in `$env.config.color_config` for syntax highlighting, which also affects how completions are displayed:

| Shape | Default Style | Affects |
|-------|--------------|---------|
| `shape_block` | Blue.bold() | Block expressions |
| `shape_bool` | LightCyan | Boolean values |
| `shape_custom` | bold() | Custom command calls |
| `shape_external` | Cyan | External command names |
| `shape_externalarg` | Green.bold() | External command arguments |
| `shape_filepath` | Cyan | File paths |
| `shape_flag` | Blue.bold() | Command flags |
| `shape_float` | Purple.bold() | Float values |
| `shape_garbage` | White.on(Red).bold() | Parse errors |
| `shape_globpattern` | Cyan.bold() | Glob patterns |
| `shape_int` | Purple.bold() | Integer values |
| `shape_internalcall` | Cyan.bold() | Internal command calls |
| `shape_list` | Cyan.bold() | List literals |
| `shape_literal` | Blue | Literal values |
| `shape_operator` | Yellow | Operators |
| `shape_pipe` | Purple.bold() | Pipe symbol |
| `shape_range` | Yellow.bold() | Range expressions |
| `shape_record` | Cyan.bold() | Record literals |
| `shape_signature` | Green.bold() | Function signatures |
| `shape_string` | Green | String values |
| `shape_string_interpolation` | Cyan.bold() | String interpolation |
| `shape_table` | Blue | Table literals |
| `shape_variable` | Purple | Variable names |

## References

- [Nushell Book: Types of Data](https://www.nushell.sh/book/types_of_data.html)
- [Nushell Language Guide: Types Overview](https://www.nushell.sh/lang-guide/chapters/types/00_types_overview.html)
- [DeepWiki: Type System and Type Checking](https://deepwiki.com/nushell/nushell/3.3-type-system-and-type-checking)
- [DeepWiki: Completion System](https://deepwiki.com/nushell/nushell/6.3-completion-system)

## Related Skills

- [references/completion.md](references/completion.md) — How completers are dispatched based on types
- [references/externs.md](references/externs.md) — Extern command type annotations
- [references/quoting.md](references/quoting.md) — String types and quoting
