# Carapace Library: Spec Generation

Reference for [carapace](https://github.com/carapace-sh/carapace)'s YAML spec generation from cobra commands in `internal/spec/`.

## spec.Spec — Generate YAML from a Command

```go
func Spec(cmd *cobra.Command) string {
    m := &bytes.Buffer{}
    enc := yaml.NewEncoder(m)
    enc.SetIndent(2)
    _ = enc.Encode(commandLine(cmd))
    return "# yaml-language-server: $schema=https://carapace.sh/schemas/command.json\n" + m.String()
}
```

Writes the command tree as YAML. Prepends a YAML language server comment for IDE schema validation.

## Command Struct

The canonical `Command` struct lives in `pkg/command/` and is imported by `internal/spec/`:

```go
import "github.com/carapace-sh/carapace/pkg/command"

type Command struct {
    Name            string            `yaml:"name"`
    Aliases         []string          `yaml:"aliases,omitempty"`
    Description     string            `yaml:"description,omitempty"`
    Group           string            `yaml:"group,omitempty"`
    Hidden          bool              `yaml:"hidden,omitempty"`
    Parsing         Parsing           `yaml:"parsing,omitempty"`
    Flags           FlagSet           `yaml:"flags,omitempty"`
    PersistentFlags FlagSet           `yaml:"persistentflags,omitempty"`
    ExclusiveFlags  [][]string        `yaml:"exclusiveflags,omitempty"`
    Run             Run               `yaml:"run,omitempty"`
    Completion      struct { ... }   `yaml:"completion,omitempty"`
    Commands      []Command          `yaml:"commands,omitempty"`
    Documentation struct { ... }     `yaml:"documentation,omitempty"`
    Examples      map[string]string  `yaml:"examples,omitempty"`
}
```

These types are the single source of truth — `carapace-spec` imports them from `carapace/pkg/command` rather than maintaining its own copy.

## FlagSet Marshaling

`FlagSet` (`map[string]Flag`) has a custom `MarshalYAML` that produces output compatible with `carapace-spec`:

- **Simple flags** (no `Nargs`, no `Default`): value is a plain string
  ```yaml
  --string=: some string flag
  ```

- **Flags with `Nargs` or `Default`**: value is an `Extended` struct
  ```yaml
  --string*!:
    description: some string flag
    default: default.txt
  --complex=:
    description: some complex flag
    nargs: 2
    default: default.txt
  ```

## Supported Fields

| Field | YAML Key | Source |
|-------|----------|--------|
| `Parsing` | `parsing` | (not yet populated) |
| `Run` | `run` | (not yet populated) |
| `Documentation` | `documentation` | `cmd.Long` → `documentation.command` |
| `Examples` | `examples` | `cmd.Example` (not yet parsed) |
| `ExclusiveFlags` | `exclusiveflags` | `cmd.Flags()` annotations → `exclusiveFlags()` |

## Recursive Command Tree

```go
func commandLine(cmd *cobra.Command) command.Command {
    c := command.Command{
        Name:            cmd.Use,
        Description:     cmd.Short,
        Aliases:         cmd.Aliases,
        Group:           cmd.GroupID,
        Hidden:          cmd.Hidden,
        Flags:           make(command.FlagSet),
        PersistentFlags: make(command.FlagSet),
        Commands:        make([]command.Command, 0),
    }

    if cmd.Long != "" {
        c.Documentation.Command = cmd.Long
    }

    c.ExclusiveFlags = exclusiveFlags(cmd)

    cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
        if cmd.PersistentFlags().Lookup(flag.Name) != nil {
            return
        }
        c.Flags[flagDefinition(flag, cmd.Flags())] = toFlag(flag)
    })

    cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
        c.PersistentFlags[flagDefinition(flag, cmd.Flags())] = toFlag(flag)
    })

    for _, subcmd := range cmd.Commands() {
        if subcmd.Name() != "_carapace" && subcmd.Deprecated == "" {
            c.Commands = append(c.Commands, commandLine(subcmd))
        }
    }
    return c
}
```

Recursively processes the command tree, skipping:
- The `_carapace` subcommand (auto-generated)
- Deprecated commands (`cmd.Deprecated != ""`)

## Helper Functions

### flagDefinition

```go
func flagDefinition(flag *pflag.Flag, flagSet *pflag.FlagSet) string {
    f := pflagfork.Flag{Flag: flag}
    prefix := pflagfork.FlagSet{FlagSet: flagSet}.Prefix()
    return f.Definition(prefix)
}
```

Produces the human-readable flag definition string used as the map key (e.g., `-v, --verbose!*?`).

### toFlag

```go
func toFlag(flag *pflag.Flag) command.Flag {
    f := pflagfork.Flag{Flag: flag}
    // extracts Longhand, Shorthand, NameAsShorthand, Repeatable,
    // Optarg, Value, Hidden, Required, Nargs, Default from pflag metadata
}
```

Converts a `pflag.Flag` to a `command.Flag` struct, populating all fields from the pflag metadata via `pflagfork.Flag` methods.

### exclusiveFlags

```go
func exclusiveFlags(cmd *cobra.Command) [][]string {
    // collects mutually exclusive flag groups from
    // cobra's "cobra_annotation_mutually_exclusive" annotations
}
```

Reads the `cobra_annotation_mutually_exclusive` annotation from each flag, deduplicates groups, and returns them as `[][]string`.

## pflagfork.Flag.Definition

The `Definition(prefix rune)` method produces a human-readable flag string used as the map key in the YAML:

```
f.Definition(prefix) // e.g., "-v, --verbose!*?" or "&f, &&flag!*?"
```

Format: `<p>shorthand, <p><p>name<suffixes>` where `<p>` is the prefix char. Suffixes (in order):
- `&` = hidden flag
- `!` = required (cobra `BashCompOneRequiredFlag` annotation)
- `*` = repeatable (Slice/Array/count type)
- `?` = optarg (`NoOptDefVal != ""`, non-bool types only)
- `=` = takes value (non-bool, non-optarg)

## Usage in command.go

`carapace.Gen(cmd)` registers the `_carapace spec` subcommand:

```go
specCmd := &cobra.Command{
    Use: "spec",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Fprint(cmd.OutOrStdout(), spec.Spec(targetCmd))
    },
}
carapaceCmd.AddCommand(specCmd)
```

Running `myapp _carapace spec` outputs the YAML spec for the entire command tree.

## Type Ownership

The canonical spec types (`Command`, `FlagSet`, `Flag`, `Run`, `Parsing`) live in `carapace/pkg/command/`. The `carapace-spec` module imports them from `carapace/pkg/command` rather than maintaining its own copy, making carapace the single source of truth for the spec format.

## YAML Schema

The generated YAML conforms to the JSON schema at `https://carapace.sh/schemas/command.json`. The YAML-LSP comment at the top enables IDE validation in editors like VS Code with the YAML extension.

## Gotchas

- **Completion fields are empty**: The `Completion` struct in the YAML type is for carapace-bin spec files, not for generated specs from running commands. Generated specs only include `Name`, `Aliases`, `Description`, `Group`, `Hidden`, `Parsing`, `Flags`, `PersistentFlags`, `ExclusiveFlags`, `Run`, `Documentation`, `Examples`, and `Commands`.
- **Hidden flags via `&`**: The `Definition()` suffix `&` indicates a hidden flag.
- **Skip _carapace**: The spec command itself is excluded from the tree to avoid polluting the generated spec.
- **No annotation of completion actions**: The generated YAML does not include completion actions — it's a structural skeleton only. carapace-bin uses this as a starting point for manual spec authoring.
- **Flag defaults are populated**: `flag.DefValue` is used for the `default` field in the `Extended` struct. Bool flags default to `"false"`, slice/array flags default to `'[]'`.
- **Flag nargs is populated**: `pflagfork.Flag.Nargs()` is used for the `nargs` field. Only flags with non-zero `Nargs` or non-empty `Default` use the `Extended` struct format.
- **ExclusiveFlags are populated**: Mutually exclusive flag groups are read from cobra's `cobra_annotation_mutually_exclusive` annotations on flags.

## Related Skills

- **references/traverse.md** — pflagfork used during runtime traversal
- **references/pflag.md** — pflagfork flag metadata fields
- **carapace skill** — user-facing spec authoring (in carapace-bin)