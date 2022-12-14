package common

import (
	"fmt"

	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
)

type Group struct {
	Cmd *cobra.Command
}

func (g Group) Tag() string {
	tag := "commands"
	if id := g.Cmd.GroupID; id != "" {
		tag = fmt.Sprintf("%v %v", id, tag)
	} else if len(g.Cmd.Parent().Groups()) != 0 {
		tag = "other commands"
	}
	return tag
}

func (g Group) Style() string {
	if g.Cmd.Parent() == nil || g.Cmd.Parent().Groups() == nil {
		return style.Default
	}

	groupStyles := []string{
		style.Carapace.H1,
		style.Carapace.H2,
		style.Carapace.H3,
		style.Carapace.H4,
		style.Carapace.H5,
		style.Carapace.H6,
		style.Carapace.H7,
		style.Carapace.H8,
		style.Carapace.H9,
		style.Carapace.H10,
		style.Carapace.H11,
		style.Carapace.H12,
	}

	for index, group := range g.Cmd.Parent().Groups() {
		if group.ID == g.Cmd.GroupID && index < len(groupStyles) {
			return groupStyles[index]
		}
	}
	return style.Default
}
