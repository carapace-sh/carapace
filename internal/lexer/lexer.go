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
		tokenset, err = split(s+`_"`, pipelines)
		if err == nil {
			last := tokenset.Tokens[len(tokenset.Tokens)-1]
			tokenset.Tokens[len(tokenset.Tokens)-1] = last[:len(last)-1]
			tokenset.Prefix = tokenset.Prefix[:len(tokenset.Prefix)-1]
			tokenset.State = OPEN_DOUBLE
		}
	}
	if err != nil && err.Error() == "EOF found when expecting closing quote" {
		tokenset, err = split(s+`_'`, pipelines)
		if err == nil {
			last := tokenset.Tokens[len(tokenset.Tokens)-1]
			tokenset.Tokens[len(tokenset.Tokens)-1] = last[:len(last)-1]
			tokenset.Prefix = tokenset.Prefix[:len(tokenset.Prefix)-1]
			tokenset.State = OPEN_SINGLE
		}
	}
	return tokenset, err
}

func split(s string, pipelines bool) (*Tokenset, error) {
	f := shlex.Split
	if pipelines {
		f = shlex.SplitP
	}
	splitted, err := f(s)
	if strings.HasSuffix(s, " ") {
		splitted = append(splitted, "")
	}
	if err != nil {
		return nil, err
	}

	if len(splitted) == 0 {
		splitted = []string{""}
	}
	return &Tokenset{
		Tokens: splitted,
		Prefix: s[:strings.LastIndex(s, splitted[len(splitted)-1])],
	}, nil
}
