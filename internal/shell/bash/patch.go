package bash

import (
	"os"

	shlex "github.com/rsteube/carapace-shlex"
)

// RedirectError current position is a redirect like `echo test >[TAB]`.
type RedirectError struct{}

func (r RedirectError) Error() string {
	return "current position is a redirect like `echo test >[TAB]`"
}

// Patch patches args if `COMP_LINE` environment variable is set.
//
// Bash passes redirects to the completion function so these need to be filtered out.
//
//	`example action >/tmp/stdout.txt --values 2>/tmp/stderr.txt fi[TAB]`
//	["example", "action", ">", "/tmp/stdout.txt", "--values", "2", ">", "/tmp/stderr.txt", "fi"]
//	["example", "action", "--values", "fi"]
func Patch(args []string) ([]string, error) { // TODO document and fix wordbreak splitting (e.g. `:`)
	compline, ok := os.LookupEnv("COMP_LINE")
	if !ok {
		return args, nil
	}

	if err := os.Unsetenv("COMP_LINE"); err != nil { // prevent it being passes along to embedded completions
		return nil, err
	}

	if compline == "" {
		return args, nil
	}

	tokens, err := shlex.Split(compline)
	if err != nil {
		return nil, err
	}

	if len(tokens) > 1 {
		if previous := tokens[len(tokens)-2]; previous.WordbreakType.IsRedirect() {
			return append(args[:1], tokens[len(tokens)-1].Value), RedirectError{}
		}
	}
	args = append(args[:1], tokens.CurrentPipeline().FilterRedirects().Words().Strings()...)
	return args, nil
}
