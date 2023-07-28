package lexer

import (
	"encoding/json"
	"testing"

	"github.com/rsteube/carapace/internal/assert"
)

func TestSplit(t *testing.T) {
	_test := func(s string, expected Tokenset) {
		t.Run(s, func(t *testing.T) {
			tokenset, err := Split(s, true)
			if err != nil {
				t.Error(err.Error())
			}

			expected, _ := json.MarshalIndent(expected, "", "  ")
			actual, _ := json.MarshalIndent(tokenset, "", "  ")
			assert.Equal(t, string(expected), string(actual))
		})
	}

	_test(``, Tokenset{
		Tokens: []string{""},
	})

	_test(` `, Tokenset{
		Tokens: []string{""},
		Prefix: ` `,
	})

	_test(`"example`, Tokenset{
		Tokens: []string{"example"},
		State:  OPEN_DOUBLE,
	})

	_test(`'example`, Tokenset{
		Tokens: []string{"example"},
		State:  OPEN_SINGLE,
	})

	_test(`example a`, Tokenset{
		Tokens: []string{"example", "a"},
		Prefix: `example `,
	})

	_test(`example  a`, Tokenset{
		Tokens: []string{"example", "a"},
		Prefix: `example  `,
	})

	_test(`example "a`, Tokenset{
		Tokens: []string{"example", "a"},
		Prefix: `example `,
		State:  OPEN_DOUBLE,
	})

	_test(`example  'a`, Tokenset{
		Tokens: []string{"example", "a"},
		Prefix: `example  `,
		State:  OPEN_SINGLE,
	})

	_test(`example action `, Tokenset{
		Tokens: []string{"example", "action", ""},
		Prefix: `example action `,
	})

	_test(`example action -`, Tokenset{
		Tokens: []string{"example", "action", "-"},
		Prefix: `example action `,
	})

	_test(`example action --`, Tokenset{
		Tokens: []string{"example", "action", "--"},
		Prefix: `example action `,
	})

	_test(`example action - `, Tokenset{
		Tokens: []string{"example", "action", "-", ""},
		Prefix: `example action - `,
	})

	_test(`example action -- `, Tokenset{
		Tokens: []string{"example", "action", "--", ""},
		Prefix: `example action -- `,
	})

	_test(`example "action" -- `, Tokenset{
		Tokens: []string{"example", "action", "--", ""},
		Prefix: `example "action" -- `,
	})

	_test(`example 'action' -- `, Tokenset{
		Tokens: []string{"example", "action", "--", ""},
		Prefix: `example 'action' -- `,
	})

	_test(`example 'action' -- | echo `, Tokenset{
		Tokens: []string{"echo", ""},
		Prefix: `example 'action' -- | echo `,
	})

	_test(`example 'action' -- || echo `, Tokenset{
		Tokens: []string{"echo", ""},
		Prefix: `example 'action' -- || echo `,
	})

	_test(`example 'action' -- && echo `, Tokenset{
		Tokens: []string{"echo", ""},
		Prefix: `example 'action' -- && echo `,
	})

	_test(`example 'action' -- ; echo `, Tokenset{
		Tokens: []string{"echo", ""},
		Prefix: `example 'action' -- ; echo `,
	})

	_test(`example 'action' -- & echo `, Tokenset{
		Tokens: []string{"echo", ""},
		Prefix: `example 'action' -- & echo `,
	})
}
