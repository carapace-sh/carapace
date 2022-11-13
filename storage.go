package carapace

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// TODO storage needs better naming and structure

type entry struct {
	flag          ActionMap
	positional    []Action
	positionalAny Action
	dash          []Action
	dashAny       Action
	preinvoke     func(cmd *cobra.Command, flag *pflag.Flag, action Action) Action
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
	} else {
		return s.preinvoke(cmd, flag, s.get(cmd).flag[name])
	}
}

func (s _storage) preinvoke(cmd *cobra.Command, flag *pflag.Flag, action Action) Action {
	// TODO yuck - clean this up
	entry := s.get(cmd)

	// The flag might be a slice or a map, in which case it can accept more
	// than one value. If it is, this call wraps its completer in an ActionMultiParts.
	a := getRepeatableFlag(flag, action)

	if entry.preinvoke != nil {
		a = ActionCallback(func(c Context) Action {
			return entry.preinvoke(cmd, flag, action)
		})
	}
	if cmd.HasParent() {
		// TODO from cmd passed to preinvoke function
		return s.preinvoke(cmd.Parent(), flag, a)
	}
	return a
}

// If the given flag is a repeatable one (slice or map), build a special
// completer that allows completions of comma-separated values for this flag.
func getRepeatableFlag(flag *pflag.Flag, action Action) Action {
	if flag == nil {
		return action
	}

	flagType := flag.Value.Type()
	flagTypeRepeatable := strings.HasPrefix(flagType, "[]") || strings.HasPrefix(flagType, "map[")

	if !flagTypeRepeatable {
		return action
	}

	// WARN: There must be a better way to check parse the values than this.
	values := strings.TrimPrefix(flag.Value.String(), "map")
	values = strings.TrimPrefix(values, "[")
	values = strings.TrimSuffix(values, "]")

	listAction := ActionMultiParts(",", func(c Context) Action {
		// First filter out all the values that have been found
		// with other invocations of this flag. Note that this
		// also includes values set at runtime, or with NoOptDefVal.
		var alreadySet []string

		if flagTypeRepeatable {
			// This might have the unintended effect or pulling out
			// the last comma, not sure if this is dangerous.
			alreadySet = append(alreadySet, strings.Split(values, " ")...)
		} else {
			alreadySet = append(alreadySet, flag.Value.String())
		}

		// Finally remove the parts found in the current string.
		return action.Invoke(c).Filter(alreadySet).Filter(c.Parts).ToA().noSpace(false)
	})

	return listAction
}

func (s _storage) getPositional(cmd *cobra.Command, pos int) Action {
	entry := s.get(cmd)

	// TODO nil check?
	if !common.IsDash(cmd) {
		if len(entry.positional) > pos {
			return s.preinvoke(cmd, nil, entry.positional[pos])
		}
		return s.preinvoke(cmd, nil, entry.positionalAny)
	} else {
		if len(entry.dash) > pos {
			return s.preinvoke(cmd, nil, entry.dash[pos])
		}
		return s.preinvoke(cmd, nil, entry.dashAny)
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
