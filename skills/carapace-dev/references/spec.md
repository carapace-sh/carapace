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
    ExclusiveFlags  [][]string        `yaml:"exclusiveflags,omitempty"`
    Flags           map[string]string `yaml:"flags,omitempty"`
    PersistentFlags map[string]string `yaml:"persistentflags,omitempty"`
    Completion      struct {
        Flag          map[string][]string `yaml:"flag,omitempty"`
        Positional    [][]string          `yaml:"positional,omitempty"`
        PositionalAny []string            `yaml:"positionalany,omitempty"`
        Dash          [][]string          `yaml:"dash,omitempty"`
        DashAny       []string            `yaml:"dashany,omitempty"`
    } `yaml:"completion,omitempty"`
    Commands []Command `yaml:"commands,omitempty"`
}
```

Flags map key is the flag definition string (e.g., `-v, --verbose!*?`) and value is the usage string.

## Recursive Command Tree

```go
func command(cmd *cobra.Command) Command {
    c := Command{
        Name:            cmd.Use,
        Description:     cmd.Short,
        Aliases:         cmd.Aliases,
        Group:           cmd.GroupID,
        Hidden:          cmd.Hidden,
        Flags:           make(map[string]string),
        PersistentFlags: make(map[string]string),
        Commands:        make([]Command, 0),
    }

    cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
        if cmd.PersistentFlags().Lookup(flag.Name) != nil {
            return
        }
        f := pflagfork.Flag{Flag: flag}
        c.Flags[f.Definition()] = f.Usage
    })

    cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
        f := pflagfork.Flag{Flag: flag}
        c.PersistentFlags[f.Definition()] = f.Usage
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

The `Definition()` method produces a human-readable flag string used as the map key in the YAML:

```
f.Definition() // e.g., "-v, --verbose!*?" or "-f, --file"
```

Format: `-shorthand, --name<optarg>?` where:
- `!` = required (no `NoOptDefVal`)
- `*` = optarg (`NoOptDefVal != ""`)
- `?` = hidden flag

## Flag Field Extraction

Uses `pflagfork.Flag` to wrap the raw `*pflag.Flag` and read unexported fields via reflection:

```go
f := pflagfork.Flag{Flag: flag}
c.Flags[f.Definition()] = f.Usage
```

The `pflagfork.Flag` wrapper is the same type used by `traverse()` — it provides `Mode()`, `Nargs()`, `OptargDelimiter()`, `Definition()`, and other methods that read unexported pflag fields.

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

The `carapace-spec` module (in carapace-bin) has a separate, slim `internal/pflagfork` for code generation — it only reads flag metadata, no parsing logic. The spec generation reads the same fields that `traverse()` uses at runtime.

## YAML Schema

The generated YAML conforms to the JSON schema at `https://carapace.sh/schemas/command.json`. The YAML-LSP comment at the top enables IDE validation in editors like VS Code with the YAML extension.

## Gotchas

- **Completion fields are empty**: The `Completion` struct in the YAML type is for carapace-bin spec files, not for generated specs from running commands. Generated specs only include `Name`, `Aliases`, `Description`, `Group`, `Hidden`, `Flags`, `PersistentFlags`, and `Commands`.
- **ExclusiveFlags not populated**: The `ExclusiveFlags` field exists in the struct but is always empty in the generated output (TODO).
- **Hidden flags via `?`**: The `Definition()` suffix `?` indicates a hidden flag. Not all flag types support hidden.
- **Skip _carapace**: The spec command itself is excluded from the tree to avoid polluting the generated spec.
- **No annotation of completion actions**: The generated YAML does not include completion actions — it's a structural skeleton only. carapace-bin uses this as a starting point for manual spec authoring.

## Related Skills

- **references/traverse.md** — pflagfork used during runtime traversal
- **references/pflag.md** — pflagfork flag metadata fields
- **carapace skill** — user-facing spec authoring (in carapace-bin)