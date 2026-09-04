package command

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

type Run string

func Alias(s ...string) (Run, error) {
	if len(s) == 0 {
		return "", errors.New("invalid alias")
	}

	var alias = struct { // pseudo-struct to enforce `flow` style
		A []string `yaml:",flow"`
	}{s}

	m, err := yaml.Marshal(alias)
	if err != nil {
		return "", err
	}
	return Run(m[3:]), nil // cut `a: ` prefix
}

func (r Run) Type() string { // TODO return custom type
	switch s := string(r); {
	case strings.HasPrefix(s, "$"):
		return "macro"
	case strings.HasPrefix(s, "#!"):
		return "script"
	case strings.HasPrefix(s, "["):
		return "alias"
	default:
		return ""
	}
}

func (r *Run) UnmarshalYAML(value *yaml.Node) error {
	var script string
	if err := value.Decode(&script); err == nil {
		*r = Run(script)
		return nil
	}

	var alias []string
	if err := value.Decode(&alias); err != nil {
		return err
	}

	var err error
	*r, err = Alias(alias...)
	return err
}
