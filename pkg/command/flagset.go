package command

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type FlagSet map[string]Flag

type Extended struct {
	Description string `yaml:"description,omitempty" json:"description,omitempty" jsonschema_description:"Description of the flag"`
	Nargs       int    `yaml:"nargs,omitempty" json:"nargs,omitempty" jsonschema_description:"Amount of arguments consumed"`
	Default     string `yaml:"default,omitempty" json:"default,omitempty" jsonschema_description:"Default value of the flag"`
}

func (fs FlagSet) MarshalYAML() (any, error) {
	m := make(map[string]any)

	for _, f := range fs {
		switch {
		case f.Nargs != 0 || f.Default != "":
			m[f.format()] = Extended{
				Description: f.Description,
				Nargs:       f.Nargs,
				Default:     f.Default,
			}
		default:
			m[f.format()] = f.Description
		}
	}
	return m, nil
}

func (fs *FlagSet) UnmarshalYAML(value *yaml.Node) error {
	m := make(map[string]any)
	if err := value.Decode(&m); err != nil {
		return err
	}

	flagSet := make(FlagSet)
	for k, v := range m {
		switch v := v.(type) {
		case string:
			f, err := parseFlag(k, v)
			if err != nil {
				return err
			}
			flagSet[f.Name()] = *f // TODO ref?

		case map[string]any:
			f, err := parseFlag(k, "")
			if err != nil {
				return err
			}
			f.Description, _ = v["description"].(string)
			f.Nargs, _ = v["nargs"].(int)
			switch d := v["default"].(type) {
			case string:
				f.Default = d
			case int:
				f.Default = fmt.Sprint(d)
			case bool:
				f.Default = fmt.Sprint(d)
			case nil:
			default:
				f.Default = fmt.Sprint(d)
			}

			flagSet[f.Name()] = *f // TODO ref?

		default:
			return errors.New("invalid type for FlagSet")
		}
	}
	*fs = flagSet
	return nil
}
