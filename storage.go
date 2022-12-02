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
		a := s.preinvoke(cmd, flag, s.get(cmd).flag[name])

		return ActionCallback(func(c Context) Action { // TODO verify order of execution is correct
			invoked := a.Invoke(c)
			if invoked.meta.Usage == "" {
				invoked.meta.Usage = flag.Usage
			}
			return invoked.ToA()
		})
	}
}

func (s _storage) preinvoke(cmd *cobra.Command, flag *pflag.Flag, action Action) Action {
	// TODO yuck - clean this up
	entry := s.get(cmd)
	a := action
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

func (s _storage) getPositional(cmd *cobra.Command, pos int) Action {
	entry := s.get(cmd)

	var a Action
	// TODO nil check?
	if !common.IsDash(cmd) {
		if len(entry.positional) > pos {
			a = s.preinvoke(cmd, nil, entry.positional[pos])
		} else {
			a = s.preinvoke(cmd, nil, entry.positionalAny)
		}
	} else {
		if len(entry.dash) > pos {
			a = s.preinvoke(cmd, nil, entry.dash[pos])
		} else {
			a = s.preinvoke(cmd, nil, entry.dashAny)
		}
	}

	return ActionCallback(func(c Context) Action { // TODO verify order of execution is correct
		invoked := a.Invoke(c)
		if invoked.meta.Usage == "" && len(strings.Fields(cmd.Use)) > 1 {
			invoked.meta.Usage = cmd.Use
		}
		return invoked.ToA()
	})
}

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
