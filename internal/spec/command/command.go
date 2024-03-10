package command

type Command struct {
	Name        string   `yaml:"name" json:"name" jsonschema_description:"Name of the command"`
	Aliases     []string `yaml:"aliases,omitempty" json:"aliases,omitempty" jsonschema_description:"Aliases of the command"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty" jsonschema_description:"Description of the command"`
	Group       string   `yaml:"group,omitempty" json:"group,omitempty" jsonschema_description:"Group of the command"`
	Hidden      bool     `yaml:"hidden,omitempty" json:"hidden,omitempty" jsonschema_description:"Hidden state of the command"`
	Parsing     Parsing  `yaml:"parsing,omitempty" json:"parsing,omitempty" jsonschema_description:"Flag parsing mode of the command" jsonschema:"enum=interspersed,enum=non-interspersed,enum=disabled"`

	Flags           map[string]string `yaml:"flags,omitempty" json:"flags,omitempty" jsonschema_description:"Flags of the command with their description"`
	PersistentFlags map[string]string `yaml:"persistentflags,omitempty" json:"persistentflags,omitempty" jsonschema_description:"Persistent flags of the command with their description"`
	ExclusiveFlags  [][]string        `yaml:"exclusiveflags,omitempty" json:"exclusiveflags,omitempty" jsonschema_description:"Flags that are mutually exclusive"`
	Run             string            `yaml:"run,omitempty" json:"run,omitempty" jsonschema_description:"Command or script to execute in runnable mode"`
	Completion      struct {
		Flag          map[string][]string `yaml:"flag,omitempty" json:"flag,omitempty" jsonschema_description:"Flag completion"`
		Positional    [][]string          `yaml:"positional,omitempty" json:"positional,omitempty" jsonschema_description:"Positional completion"`
		PositionalAny []string            `yaml:"positionalany,omitempty" json:"positionalany,omitempty" jsonschema_description:"Positional completion for every other position"`
		Dash          [][]string          `yaml:"dash,omitempty" json:"dash,omitempty" jsonschema_description:"Dash completion"`
		DashAny       []string            `yaml:"dashany,omitempty" json:"dashany,omitempty" jsonschema_description:"Dash completion of every other position"`
	} `yaml:"completion,omitempty" json:"completion,omitempty" jsonschema_description:"Completion definition"`
	Commands []Command `yaml:"commands,omitempty" json:"commands,omitempty" jsonschema_description:"Subcommands of the command"`

	Documentation struct {
		Command       string            `yaml:"command,omitempty" json:"command,omitempty" jsonschema_description:"Documentation of the command"`
		Flag          map[string]string `yaml:"flag,omitempty" json:"flag,omitempty" jsonschema_description:"Documentation of flags"`
		Positional    []string          `yaml:"positional,omitempty" json:"positional,omitempty" jsonschema_description:"Documentation of positional arguments"`
		PositionalAny string            `yaml:"positionalany,omitempty" json:"positionalany,omitempty" jsonschema_description:"Documentation of other positional arguments"`
		Dash          []string          `yaml:"dash,omitempty" json:"dash,omitempty" jsonschema_description:"Documentation of dash arguments"`
		DashAny       string            `yaml:"dashany,omitempty" json:"dashany,omitempty" jsonschema_description:"Documentation of other dash arguments"`
	} `yaml:"documentation,omitempty" json:"documentation,omitempty" jsonschema_description:"Documentation"`
	Examples map[string]string `yaml:"examples,omitempty" json:"examples,omitempty" jsonschema_description:"Examples"`
}

func (c *Command) AddFlag(f Flag) {
	switch {
	case f.Persistent:
		if c.PersistentFlags == nil {
			c.PersistentFlags = make(map[string]string)
		}
		c.PersistentFlags[f.format()] = f.Usage

	default:
		if c.Flags == nil {
			c.Flags = make(map[string]string)
		}
		c.Flags[f.format()] = f.Usage
	}
}
