package command

import (
	"testing"

	"github.com/carapace-sh/carapace/pkg/assert"
	"gopkg.in/yaml.v3"
)

func TestFlagSet(t *testing.T) {
	fs := FlagSet{
		"string": Flag{
			Longhand:    "string",
			Description: "some string flag",
			Required:    true,
			Repeatable:  true,
		},
		"complex": Flag{
			Longhand:    "complex",
			Shorthand:   "c",
			Description: "some complex flag",
			Value:       true,
			Nargs:       2,
			Default:     "default.txt",
		},
		"defaultonly": Flag{
			Longhand:    "defaultonly",
			Description: "flag with default only",
			Value:       true,
			Default:     "/tmp/out.txt",
		},
	}

	expected := `--defaultonly=:
    description: flag with default only
    default: /tmp/out.txt
--string*!: some string flag
-c, --complex=:
    description: some complex flag
    nargs: 2
    default: default.txt
`
	m, err := yaml.Marshal(fs)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, string(m), expected)

	var actual FlagSet
	if err := yaml.Unmarshal(m, &actual); err != nil {
		t.Error(err)
	}
	assert.Equal(t, actual, fs)
}
