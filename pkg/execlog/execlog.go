package execlog

import (
	"github.com/rsteube/carapace/internal/log"
	"github.com/rsteube/carapace/pkg/util"
	"github.com/rsteube/carapace/third_party/golang.org/x/sys/execabs"
)

type Cmd struct {
	*execabs.Cmd
}

// Command is like execabs.Command but logs args on execution.
func Command(name string, arg ...string) *Cmd {
	cmd := &Cmd{
		execabs.Command(name, arg...),
	}
	return cmd
}

func (c *Cmd) Run() error {
	log.LOG.Printf("executing %#v", util.FormatCmd(c.Args...))
	return c.Cmd.Run()
}

func (c *Cmd) Start() error {
	log.LOG.Printf("executing %#v", util.FormatCmd(c.Args...))
	return c.Cmd.Start()
}

// Command is the same as execabs.Command.
func LookPath(file string) (string, error) {
	return execabs.LookPath(file)
}
