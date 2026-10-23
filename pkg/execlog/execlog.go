package execlog

import (
	"context"

	shlex "github.com/carapace-sh/carapace-shlex/v2"
	"github.com/carapace-sh/carapace/internal/log"
	"github.com/carapace-sh/carapace/third_party/golang.org/x/sys/execabs"
)

type Cmd struct {
	*execabs.Cmd
	mock []byte // if set, execution returns this output instead of running the command
}

// Command is like execabs.Command but logs args on execution.
func Command(name string, arg ...string) *Cmd {
	cmd := &Cmd{
		Cmd: execabs.Command(name, arg...),
	}
	return cmd
}

// CommandContext is like Command but includes a context.
func CommandContext(ctx context.Context, name string, arg ...string) *Cmd {
	cmd := &Cmd{
		Cmd: execabs.CommandContext(ctx, name, arg...),
	}
	return cmd
}

// Mock returns a command simulating invocation of `name` without executing anything.
//
// Context.Command uses this to replay mocked replies: spawning `echo` is
// unreliable on windows (it is no real executable and arguments may not
// contain newlines).
func Mock(output string) *Cmd {
	// mimic `echo` semantics: a trailing newline is always present
	return &Cmd{
		Cmd:  execabs.Command(""),
		mock: []byte(output + "\n"),
	}
}

func (c *Cmd) CombinedOutput() ([]byte, error) {
	log.LOG.Printf("executing %#v", shlex.Join(c.Args, shlex.Default))
	if c.mock != nil {
		return c.mock, nil
	}
	return c.Cmd.CombinedOutput()
}

func (c *Cmd) Output() ([]byte, error) {
	log.LOG.Printf("executing %#v", shlex.Join(c.Args, shlex.Default))
	if c.mock != nil {
		return c.mock, nil
	}
	return c.Cmd.Output()
}

func (c *Cmd) Run() error {
	log.LOG.Printf("executing %#v", shlex.Join(c.Args, shlex.Default))
	if c.mock != nil {
		if c.Stdout != nil {
			c.Stdout.Write(c.mock) //nolint:errcheck // mocked reply, nothing executes
		}
		return nil
	}
	return c.Cmd.Run()
}

func (c *Cmd) Start() error {
	log.LOG.Printf("executing %#v", shlex.Join(c.Args, shlex.Default))
	if c.mock != nil {
		if c.Stdout != nil {
			c.Stdout.Write(c.mock) //nolint:errcheck // mocked reply, nothing executes
		}
		return nil
	}
	return c.Cmd.Start()
}

// Command is the same as execabs.Command.
func LookPath(file string) (string, error) {
	return execabs.LookPath(file)
}
