package command

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Flag struct {
	Longhand    string
	Shorthand   string
	Description string

	NameAsShorthand bool
	Repeatable      bool
	Optarg          bool
	Value           bool
	Hidden          bool
	Required        bool
	Persistent      bool

	Nargs   int
	Default string
}

func (f Flag) Name() string {
	if f.Longhand != "" {
		return f.Longhand
	}
	return f.Shorthand
}

func (f Flag) format() string {
	var s string

	if f.Shorthand != "" {
		s += "-" + f.Shorthand
		if f.Longhand != "" {
			s += ", "
		}
	}

	if f.Longhand != "" {
		switch {
		case f.NameAsShorthand:
			s += "-" + f.Shorthand
		default:
			s += "--" + f.Longhand
		}
	}

	switch {
	case f.Optarg:
		s += "?"
	case f.Value:
		s += "="
	}

	if f.Repeatable {
		s += "*"
	}

	if f.Required {
		s += "!"
	}

	if f.Hidden {
		s += "&"
	}

	return s
}

func parseFlag(s, description string) (*Flag, error) {
	r := regexp.MustCompile(`^(?P<shorthand>-[^-][^ =*?&!]*)?(, )?(?P<longhand>-[-]?[^ =*?&!]*)?(?P<modifier>[=*?&!]*)$`)
	if !r.MatchString(s) {
		return nil, fmt.Errorf("flag syntax invalid: %v", s)
	}

	matches := findNamedMatches(r, s)

	f := &Flag{}
	f.Longhand = strings.TrimLeft(matches["longhand"], "-")
	f.Shorthand = strings.TrimPrefix(matches["shorthand"], "-")
	f.NameAsShorthand = (matches["longhand"] != "" && !strings.HasPrefix(matches["longhand"], "--"))
	f.Description = description
	f.Repeatable = strings.Contains(matches["modifier"], "*")
	f.Optarg = strings.Contains(matches["modifier"], "?")
	f.Value = f.Optarg || strings.Contains(matches["modifier"], "=")
	f.Hidden = strings.Contains(matches["modifier"], "&")
	f.Required = strings.Contains(matches["modifier"], "!")
	if matches["nargs"] != "" {
		var err error
		if f.Nargs, err = strconv.Atoi(matches["nargs"]); err != nil {
			return nil, err
		}
	}

	if f.Longhand == "" && f.Shorthand == "" {
		return nil, fmt.Errorf("malformed flag: '%v'", s)
	}
	return f, nil
}

func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}
