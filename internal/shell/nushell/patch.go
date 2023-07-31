package nushell

import (
	"strings"

	"github.com/rsteube/carapace/internal/lexer"
)

// Patch uses the lexer to parse and patch given arguments which
// are currently passed unprocessed to the completion function.
//
// see https://www.nushell.sh/book/working_with_strings.html
func Patch(args []string) []string {
	for index, arg := range args {
		if len(arg) == 0 {
			continue
		}

		switch arg[0] {
		case '"', "'"[0]:
			if tokenset, err := lexer.Split(arg, false); err == nil {
				args[index] = tokenset.Tokens[0]
			}
		case '`':
			args[index] = strings.Trim(arg, "`")
		}
	}
	return args
}
