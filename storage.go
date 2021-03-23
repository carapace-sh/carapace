package carapace

import (
	"fmt"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
)

// TODO storage needs better naming and structure

type entry struct {
	flag          ActionMap
	positional    []Action
	positionalAny Action
}

type _storage map[*cobra.Command]*entry

func (s _storage) get(cmd *cobra.Command) (e *entry) {
	var ok bool
	if e, ok = s[cmd]; !ok {
		e = &entry{}
		s[cmd] = e
	}
	return
}

func (s _storage) getFlag(cmd *cobra.Command, name string) Action {
	if flag := cmd.LocalFlags().Lookup(name); flag == nil && cmd.HasParent() {
		return s.getFlag(cmd.Parent(), name)
	}
	// TODO nil check?
	return s.get(cmd).flag[name]
}

func (s _storage) getPositional(cmd *cobra.Command, pos int) Action {
	entry := s.get(cmd)
	// TODO nil check?
	if len(entry.positional) > pos {
		return entry.positional[pos]
	} else {
		return entry.positionalAny
	}
}

// TODO implicit execution during build - go:generate possible?
func (s _storage) check() []string {
	errors := make([]string, 0)
	for cmd, entry := range s {
		for name := range entry.flag {
			if flag := cmd.LocalFlags().Lookup(name); flag == nil {
				errors = append(errors, fmt.Sprintf("unknown flag for %s: %s\n", uid.Command(cmd), name))
			}
		}
	}
	return errors
}

var storage = make(_storage)
