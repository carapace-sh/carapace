package sandbox

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/internal/assert"
	"github.com/rsteube/carapace/internal/common"
	"github.com/spf13/cobra"
)

type Sandbox struct {
	t    *testing.T
	f    func() *cobra.Command
	env  map[string]string
	keep bool
	mock *common.Mock
}

func newSandbox(t *testing.T, f func() *cobra.Command) Sandbox {
	tempDir, err := os.MkdirTemp(os.TempDir(), "carapace-sandbox_"+t.Name()+"_")
	if err != nil {
		t.Fatal("failed to create sandbox dir: " + err.Error())
	}
	return Sandbox{
		t:   t,
		f:   f,
		env: make(map[string]string),
		mock: &common.Mock{
			Dir:     tempDir,
			Replies: make(map[string]string),
		},
	}
}

func (s *Sandbox) Keep() {
	s.keep = true
}

func (s *Sandbox) Env(key, value string) {
	s.env[key] = value
}

func (s *Sandbox) remove() {
	if !s.keep && strings.HasPrefix(s.mock.Dir, os.TempDir()) {
		os.RemoveAll(s.mock.Dir)
	}
}

func (s *Sandbox) Files(args ...string) {
	if len(args)%2 != 0 {
		s.t.Errorf("invalid amount of arguments: %v", len(args))
	}

	if !strings.HasPrefix(s.mock.Dir, os.TempDir()) {
		s.t.Errorf("sandbox dir not in os.TempDir: %v", s.mock.Dir)
	}

	for i := 0; i < len(args); i += 2 {
		file := args[i]
		content := args[i+1]

		if strings.HasPrefix(file, "../") {
			s.t.Fatalf("invalid filename: %v", file)
		}

		path := fmt.Sprintf("%v/%v", s.mock.Dir, file)

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil && !os.IsExist(err) {
			s.t.Fatal(err.Error())
		}

		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			s.t.Fatal(err.Error())
		}
	}

}

// Reply mocks a command for given arguments (Only works for `(Context).Command`).
func (s *Sandbox) Reply(args ...string) reply {
	m, _ := json.Marshal(args)
	return reply{s, string(m)}
}

type reply struct {
	*Sandbox
	call string
}

// With sets the output for the mocked command.
func (r reply) With(s string) {
	r.mock.Replies[r.call] = s
}

// Run invokes `go run` on given package for sandbox tests.
func Run(t *testing.T, pkg string) (f func(func(s *Sandbox))) {
	cmdF := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:                "integration",
			Run:                func(cmd *cobra.Command, args []string) {},
			DisableFlagParsing: true,
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		}

		carapace.Gen(cmd).PositionalAnyCompletion(
			carapace.ActionCallback(func(c carapace.Context) carapace.Action {
				args := []string{"run", pkg, "_carapace", "export", ""}
				args = append(args, c.Args...)
				args = append(args, c.CallbackValue)
				return carapace.ActionExecCommand("go", args...)(func(output []byte) carapace.Action {
					return carapace.ActionImport(output)
				})
			}),
		)
		return cmd
	}

	return func(f func(s *Sandbox)) {
		s := newSandbox(t, cmdF)
		defer s.remove()
		f(&s)
	}
}

// Run executes the sandbox with given arguments.
func (s *Sandbox) Run(args ...string) run {
	c := carapace.NewContext(args)
	for key, value := range s.env {
		c.Setenv(key, value)
	}
	// TODO actually invoke it here instead of Expect
	m, _ := json.Marshal(args)

	return run{
		s.t,
		string(m),
		s.mock.Dir,
		c,
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			b, _ := json.Marshal(s.mock)
			c.Setenv("CARAPACE_SANDBOX", string(b))
			return carapace.ActionExecute(s.f()).Invoke(c).ToA()
		}),
	}
}

type run struct {
	t       *testing.T
	id      string
	dir     string
	context carapace.Context
	actual  carapace.Action
}

func (r run) invoke(a carapace.Action) string {
	invoked := a.Invoke(r.context)
	meta, rawValues := common.FromInvokedAction(invoked) // TODO nil check
	rawValues = rawValues.FilterPrefix(r.context.CallbackValue)
	sort.Sort(common.ByValue(rawValues))
	// TODO improve this
	m, err := json.MarshalIndent(struct {
		Meta      common.Meta
		RawValues common.RawValues
	}{
		Meta:      meta,
		RawValues: rawValues,
	}, "", "  ")
	if err != nil {
		r.t.Fatal(err.Error())
	}
	return string(m)
}

// Expects validates output of Run with given Action.
func (r run) Expect(expected carapace.Action) {
	r.t.Run(r.id, func(t *testing.T) {
		// t.Parallel() TODO prevent concurrent map write for this (storage.go)
		assert.Equal(r.t, r.invoke(expected.Chdir(r.dir)), r.invoke(r.actual))
	})
}
