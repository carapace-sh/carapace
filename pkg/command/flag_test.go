package command

import (
	"reflect"
	"testing"
)

func TestParseFlag(t *testing.T) {
	test := func(
		description string,
		id string,
		expected *Flag,
	) {
		t.Run(description, func(t *testing.T) {
			expected.Description = description // skip usage test
			f, err := parseFlag(id, description)
			if err != nil {
				t.Error(err.Error())
			}
			if !reflect.DeepEqual(f, expected) {
				t.Error(description)
				t.Logf("expected: %#v", expected)
				t.Logf("actual:   %#v", f)
			}
		})
	}

	test("shorthand-only", "-s", &Flag{
		Shorthand: "s",
	})

	test("shorthand-only slice", "-s*", &Flag{
		Shorthand:  "s",
		Repeatable: true,
	})

	test("shorthand-only slice value", "-s*=", &Flag{
		Shorthand:  "s",
		Repeatable: true,
		Value:      true,
	})

	test("shorthand-only value", "-s=", &Flag{
		Shorthand: "s",
		Value:     true,
	})

	test("shorthand-only optarg", "-s?", &Flag{
		Shorthand: "s",
		Value:     true,
		Optarg:    true,
	})

	test("longhand-only", "--long", &Flag{
		Longhand: "long",
	})

	test("longhand-only value", "--long=", &Flag{
		Longhand: "long",
		Value:    true,
	})

	test("longhand-only slice optarg", "--long*?", &Flag{
		Longhand:   "long",
		Optarg:     true,
		Repeatable: true,
		Value:      true,
	})

	test("both", "-s, --long", &Flag{
		Shorthand: "s",
		Longhand:  "long",
	})

	test("both value", "-s, --long=", &Flag{
		Shorthand: "s",
		Longhand:  "long",
		Value:     true,
	})

	test("both optarg", "-s, --long?", &Flag{
		Shorthand: "s",
		Longhand:  "long",
		Value:     true,
		Optarg:    true,
	})

	test("both slice optarg", "-s, --long*?", &Flag{
		Shorthand:  "s",
		Longhand:   "long",
		Value:      true,
		Repeatable: true,
		Optarg:     true,
	})

	test("nonposix shorthand", "-short", &Flag{
		Shorthand: "short",
	})

	test("nonposix shorthand optarg", "-short?", &Flag{
		Shorthand: "short",
		Value:     true,
		Optarg:    true,
	})

	test("nonposix both", "-short, -long*", &Flag{
		Shorthand:       "short",
		Longhand:        "long",
		Repeatable:      true,
		NameAsShorthand: true,
	})

	test("nonposix mixed", "-short, --long", &Flag{
		Shorthand:       "short",
		Longhand:        "long",
		NameAsShorthand: false,
	})
}
