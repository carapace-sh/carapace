package carapace

import (
	"strings"

	"github.com/rsteube/carapace/internal/config"
)

type carapaceConfig struct {
	DescriptionLength int
}

func (c *carapaceConfig) Completion() ActionMap {
	return ActionMap{
		"DescriptionLength": ActionValues("40", "30"),
	}
}

var conf = carapaceConfig{
	DescriptionLength: 40,
}

func init() {
	config.RegisterConfig("carapace", &conf)
}

type configI interface {
	Completion() ActionMap
}

func ActionConfigs() Action {
	return ActionMultiParts("=", func(c Context) Action {
		switch len(c.Parts) {
		case 0:
			return ActionMultiParts(".", func(c Context) Action {
				switch len(c.Parts) {
				case 0:
					return ActionValues(config.GetConfigs()...).Invoke(c).Suffix(".").ToA()
				case 1:
					fields, err := config.GetConfigFields(c.Parts[0])
					if err != nil {
						return ActionMessage(err.Error())
					}

					vals := make([]string, 0)
					for _, field := range fields {
						vals = append(vals, field.Name, field.Description, field.Style)
					}
					return ActionStyledValuesDescribed(vals...).Invoke(c).Suffix("=").ToA()
				default:
					return ActionValues()
				}
			})
		case 1:
			if m := config.GetConfigMap(strings.Split(c.Parts[0], ".")[0]); m != nil {
				if i, ok := m.(configI); ok {
					// TODO check splitted length
					return i.Completion()[strings.Split(c.Parts[0], ".")[1]]
				}
			}
			return ActionValues()
		default:
			return ActionValues()
		}
	})
}
