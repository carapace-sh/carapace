# Carapace Library: Spec Generation

Reference for [carapace](https://github.com/carapace-sh/carapace)'s YAML spec generation from cobra commands in `internal/spec/`.

## spec.Spec — Generate YAML from a Command

```go
func Spec(cmd *cobra.Command) string {
    m, _ := yaml.Marshal(command(cmd))
    return "# yaml-language-server: $schema=https://carapace.sh/schemas/command.json\n" + string(m)
}
```

Writes the command tree as YAML. Prepends a YAML language server comment for IDE schema validation.

## Command Struct

```go
type Command struct {
    Name            string            `yaml:"name"`
    Aliases         []string          `yaml:"aliases,omitempty"`
    Description     string            `yaml:"description,omitempty"`
    Group           string            `yaml:"group,omitempty"`
    Hidden          bool              `yaml:"hidden,omitempty"`
    Parsing         string            `yaml:"parsing,omitempty"`
    Flags           FlagSet           `yaml:"flags,omitempty"`
    PersistentFlags FlagSet           `yaml:"persistentflags,omitempty"`
    ExclusiveFlags  [][]string        `yaml:"exclusiveflags,omitempty"`
    Run             string            `yaml:"run,omitempty"`
    Completion      struct {
        Flag          map[string][]string `yaml:"flag,omitempty"`
        Positional    [][]string          `yaml:"positional,omitempty"`
        PositionalAny []string            `yaml:"positionalany,omitempty"`
        Dash          [][]string          `yaml:"dash,omitempty"`
        DashAny       []string            `yaml:"dashany,omitempty"`
    } `yaml:"completion,omitempty"`
    Commands      []Command `yaml:"commands,omitempty"`
    Documentation struct {
        Command       string            `yaml:"command,omitempty"`
        Flag          map[string]string `yaml:"flag,omitempty"`
        Positional    []string          `yaml:"positional,omitempty"`
        PositionalAny string            `yaml:"positionalany,omitempty"`
        Dash          []string          `yaml:"dash,omitempty"`
        DashAny       string            `yaml:"dashany,omitempty"`
    } `yaml:"documentation,omitempty"`
    Examples map[string]string `yaml:"examples,omitempty"`
}
```

Flags map key is the flag definition string (e.g., `-v, --verbose!*?`). The value is either a plain string (description) for simple flags, or an `Extended` struct with `description`, `nargs`, and `default` for flags with non-zero `Nargs` or non-empty `Default`.

## FlagSet Marshaling

The `FlagSet` type (`map[string]Flag`) has a custom `MarshalYAML` that produces output compatible with `carapace-spec`:

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

## New Fields

The `Command` struct includes additional fields matching `carapace-spec`:

| Field | YAML Key | Source |
|-------|----------|--------|
| `Parsing` | `parsing` | (not yet populated) |
| `Run` | `run` | (not yet populated) |
| `Documentation` | `documentation` | `cmd.Long` → `documentation.command` |
| `Examples` | `examples` | `cmd.Example` (not yet parsed) |

## Recursive Command Tree

```go
func command(cmd *cobra.Command) Command {
    c := Command{
        Name:            cmd.Use,
        Description:     cmd.Short,
        Aliases:         cmd.Aliases,
        Group:           cmd.GroupID,
        Hidden:          cmd.Hidden,
        Flags:           make(FlagSet),
        PersistentFlags: make(FlagSet),
        Commands:        make([]Command, 0),
    }

    if cmd.Long != "" {
        c.Documentation.Command = cmd.Long
    }

    cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
        if cmd.PersistentFlags().Lookup(flag.Name) != nil {
            return
        }
        f := pflagfork.Flag{Flag: flag}
        prefix := pflagfork.FlagSet{FlagSet: cmd.Flags()}.Prefix()
        c.Flags[f.Definition(prefix)] = Flag{
            Description: f.Usage,
            Nargs:       f.Nargs(),
            Default:     flag.DefValue,
        }
    })

    cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
        f := pflagfork.Flag{Flag: flag}
        prefix := pflagfork.FlagSet{FlagSet: cmd.Flags()}.Prefix()
        c.PersistentFlags[f.Definition(prefix)] = Flag{
            Description: f.Usage,
            Nargs:       f.Nargs(),
            Default:     flag.DefValue,
        }
    })

    for _, subcmd := range cmd.Commands() {
        if subcmd.Name() != "_carapace" && subcmd.Deprecated == "" {
            c.Commands = append(c.Commands, command(subcmd))
        }
    }
    return c
}
```

Recursively processes the command tree, skipping:
- The `_carapace` subcommand (auto-generated)
- Deprecated commands (`cmd.Deprecated != ""`)

## pflagfork.Flag.Definition

The `Definition(prefix rune)` method produces a human-readable flag string used as the map key in the YAML. It now takes a prefix rune parameter (obtained from `FlagSet.Prefix()`) to support custom flag prefixes:

```
f.Definition(prefix) // e.g., "-v, --verbose!*?" or "&f, &&flag!*?"
```

Format: `<p>shorthand, <p><p>name<suffixes>` where `<p>` is the prefix char. Suffixes (in order):
- `&` = hidden flag
- `!` = required (cobra `BashCompOneRequiredFlag` annotation)
- `*` = repeatable (Slice/Array/count type)
- `?` = optarg (`NoOptDefVal != ""`, non-bool types only)
- `=` = takes value (non-bool, non-optarg)

## Flag Field Extraction

Uses `pflagfork.Flag` to wrap the raw `*pflag.Flag` and read unexported fields via reflection:

```go
f := pflagfork.Flag{Flag: flag}
prefix := pflagfork.FlagSet{FlagSet: cmd.Flags()}.Prefix()
c.Flags[f.Definition(prefix)] = f.Usage
```

The `pflagfork.Flag` wrapper is the same type used by `traverse()` — it provides `GetMode()`, `Nargs()`, `OptargDelimiter()`, `ArgumentStyle()`, `Definition(prefix)`, and other methods that read unexported pflag fields.

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

## carapace-spec (Separate Module)

The `carapace-spec` module (at `github.com/carapace-sh/carapace-spec`) defines the canonical `Command` struct and `FlagSet` type with custom `MarshalYAML`/`UnmarshalYAML` methods. The spec generation in carapace mirrors these types to produce output compatible with `carapace-spec`'s schema:

- `FlagSet` (map of definition string → `Flag`) with `MarshalYAML` producing either a string or `Extended{description, nargs, default}`
- Same `format()` string syntax for flag definitions: `-s, --long=*!&`
- Same YAML tags on `Command` struct for `Parsing`, `Run`, `Documentation`, `Examples`

The `pflagfork.Flag` wrapper provides `Nargs()`, `Definition(prefix)`, and other methods that read unexported pflag fields via reflection — matching what `carapace-spec`'s own `internal/pflagfork` does for code generation.

## YAML Schema

The generated YAML conforms to the JSON schema at `https://carapace.sh/schemas/command.json`. The YAML-LSP comment at the top enables IDE validation in editors like VS Code with the YAML extension.

## Gotchas

- **Completion fields are empty**: The `Completion` struct in the YAML type is for carapace-bin spec files, not for generated specs from running commands. Generated specs only include `Name`, `Aliases`, `Description`, `Group`, `Hidden`, `Parsing`, `Flags`, `PersistentFlags`, `Run`, `Documentation`, `Examples`, and `Commands`.
- **ExclusiveFlags not populated**: The `ExclusiveFlags` field exists in the struct but is always empty in the generated output (TODO).
- **Hidden flags via `&`**: The `Definition()` suffix `&` indicates a hidden flag.
- **Skip _carapace**: The spec command itself is excluded from the tree to avoid polluting the generated spec.
- **No annotation of completion actions**: The generated YAML does not include completion actions — it's a structural skeleton only. carapace-bin uses this as a starting point for manual spec authoring.
- **Flag defaults are populated**: `flag.DefValue` is used for the `default` field in the `Extended` struct. Bool flags default to `"false"`, slice/array flags default to `'[]'`.
- **Flag nargs is populated**: `pflagfork.Flag.Nargs()` is used for the `nargs` field. Only flags with non-zero `Nargs` or non-empty `Default` use the `Extended` struct format.

## Related Skills

- **references/traverse.md** — pflagfork used during runtime traversal
- **references/pflag.md** — pflagfork flag metadata fields
- **carapace skill** — user-facing spec authoring (in carapace-bin)