package carapace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rsteube/carapace/internal/shell/zsh"
	"github.com/rsteube/carapace/third_party/github.com/drone/envsubst"
	"github.com/rsteube/carapace/third_party/golang.org/x/sys/execabs"
)

// Context provides information during completion.
type Context struct {
	// CallbackValue contains the (partial) value (or part of it during an ActionMultiParts) currently being completed
	CallbackValue string
	// Args contains the positional arguments of current (sub)command (exclusive the one currently being completed)
	Args []string
	// Parts contains the splitted CallbackValue during an ActionMultiParts (exclusive the part currently being completed)
	Parts []string
	// Env contains environment variables for current context
	Env []string
	// Dir contains the working directory for current context
	Dir string
}

func newContext(args []string) Context {
	context := Context{
		CallbackValue: args[len(args)-1],
		Args:          args[:len(args)-1],
		Env:           os.Environ(),
	}

	if wd, err := os.Getwd(); err == nil {
		context.Dir = wd
	}
	return context
}

// LookupEnv retrieves the value of the environment variable named by the key.
func (c *Context) LookupEnv(key string) (string, bool) {
	prefix := key + "="
	for i := len(c.Env) - 1; i >= 0; i-- {
		if env := c.Env[i]; strings.HasPrefix(env, prefix) {
			return strings.SplitN(env, "=", 2)[1], true
		}
	}
	return "", false
}

// Getenv retrieves the value of the environment variable named by the key.
func (c *Context) Getenv(key string) string {
	v, _ := c.LookupEnv(key)
	return v
}

// Setenv sets the value of the environment variable named by the key.
func (c *Context) Setenv(key, value string) {
	if c.Env == nil {
		c.Env = []string{}
	}
	c.Env = append(c.Env, fmt.Sprintf("%v=%v", key, value))
}

func (c *Context) Envsubst(s string) (string, error) {
	return envsubst.Eval(s, c.Getenv)
}

// Command returns the Cmd struct to execute the named program with the given arguments.
// Env and Dir are set using the Context.
// See exec.Command for most details.
func (c Context) Command(name string, arg ...string) *execabs.Cmd {
	cmd := execabs.Command(name, arg...)
	cmd.Env = c.Env
	cmd.Dir = c.Dir
	return cmd
}

func expandHome(s string) (string, error) {
	if strings.HasPrefix(s, "~") {
		if zsh.NamedDirectories.Matches(s) {
			return zsh.NamedDirectories.Replace(s), nil
		}

		home, err := os.UserHomeDir() // TODO duplicated code
		if err != nil {
			return "", err
		}
		s = strings.Replace(s, "~/", home+"/", 1)
	}
	return s, nil
}

func (c Context) Abs(s string) (string, error) {
	var path string
	if strings.HasPrefix(s, "/") || strings.HasPrefix(s, "~") {
		path = s // path is absolute
	} else {
		expanded, err := expandHome(c.Dir)
		if err != nil {
			return "", err
		}
		abs, err := filepath.Abs(expanded)
		if err != nil {
			return "", err
		}
		path = abs + "/" + s
	}

	expanded, err := expandHome(path)
	if err != nil {
		return "", err
	}
	path = expanded

	result, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if strings.HasSuffix(path, "/") && !strings.HasSuffix(result, "/") {
		result += "/"
	} else if strings.HasSuffix(path, "/.") && !strings.HasSuffix(result, "/.") {
		result += "/."
	}
	return result, nil
}
