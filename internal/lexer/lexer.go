package lexer

import (
	"strings"

	"github.com/rsteube/carapace/third_party/github.com/google/shlex"
)

type State int

const (
	UNQUOTED = iota
	OPEN_DOUBLE
	OPEN_SINGLE
)

type Tokenset struct {
	Tokens []string
	Prefix string
	State  State
}

func Split(s string, pipelines bool) (*Tokenset, error) {
	tokenset, err := split(s, pipelines)
	if err != nil && err.Error() == "EOF found when expecting closing quote" {
		tokenset, err = split(s+`"`, pipelines)
		if err == nil {
			tokenset.State = OPEN_DOUBLE
		}
	}
	if err != nil && err.Error() == "EOF found when expecting closing quote" {
		tokenset, err = split(s+`'`, pipelines)
		if err == nil {
			tokenset.State = OPEN_SINGLE
		}
	}
	return tokenset, err
}

func split(s string, pipelines bool) (*Tokenset, error) {
	splitted, prefix, err := shlex.SplitP(s, pipelines)
	if strings.HasSuffix(s, " ") {
		splitted = append(splitted, "")
	}
	if err != nil {
		return nil, err
	}

	if len(splitted) == 0 {
		splitted = []string{""}
	}

	if len(splitted[len(splitted)-1]) == 0 {
		prefix = s
	}

	t := &Tokenset{
		Tokens: splitted,
		Prefix: prefix,
	}
	return t, nil
}
