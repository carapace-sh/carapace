package zsh

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"text/template"
)

type Completions struct {
	actions map[string]Action
}

func (c Completions) invokeCallback(uid string, args []string) Action {
	if action, ok := c.actions[uid]; ok {
		if action.Callback != nil {
			if len(args) > 1 {
				return action.Callback(args[1:]) // TODO 1:0 with len check
			} else {
				return action.Callback([]string{})
			}
		}
	}
	return ActionMessage(fmt.Sprintf("callback %v unknown", uid))
}

func (c Completions) Generate(cmd *cobra.Command) string {
	result := fmt.Sprintf("#compdef _%v %v\n", cmd.Name(), cmd.Name())
	result += c.GenerateFunctions(cmd)

	result += fmt.Sprintf("compdef _%v %v\n", cmd.Name(), cmd.Name())
	return result
}

func (c Completions) GenerateFunctions(cmd *cobra.Command) string {
	if !cmd.HasSubCommands() && !cmd.HasFlags() {
		return fmt.Sprintf("function %v {\n  true\n}\n", uidCommand(cmd))
	}

	function_pattern := `function %v {
  %v
  %v

  _arguments -C \
%v%v
    %v
}
`

	commandsVar := ""
	if cmd.HasSubCommands() {
		commandsVar = "local -a commands"
	}

	inheritedArgs := ""
	if !cmd.HasParent() {
		inheritedArgs = "# shellcheck disable=SC2206\n  local -a -x os_args=(${words})"
	}

	flags := make([]string, 0)
	for _, flag := range zshCompExtractFlag(cmd) {
		var s string
		if action, ok := c.actions[uidFlag(cmd, flag)]; ok {
			s = "    " + snippetFlagCompletion(flag, &action) + " \\\n"
		} else {
			s = "    " + snippetFlagCompletion(flag, nil) + " \\\n"
		}

		flags = append(flags, s)
	}

	positionals := make([]string, 0)
	if cmd.HasSubCommands() {
		positionals = []string{`    "1: :->cmnds" \` + "\n", `    "*::arg:->args"`}
	} else {
		pos := 1
		for {
			if action, ok := c.actions[uidPositional(cmd, pos)]; ok {
				positionals = append(positionals, "    "+snippetPositionalCompletion(pos, action))
				pos++
			} else {
				break // TODO only consisten entriess for now
			}
		}
	}

	result := make([]string, 0)
	result = append(result, fmt.Sprintf(function_pattern, uidCommand(cmd), commandsVar, inheritedArgs, strings.Join(flags, ""), strings.Join(positionals, ""), subcommands(cmd)))
	for _, subcmd := range cmd.Commands() {
		if !subcmd.Hidden {
			result = append(result, c.GenerateFunctions(subcmd))
		}
	}

	return strings.Join(result, "\n")
}

func subcommands(cmd *cobra.Command) string {
	if !cmd.HasSubCommands() {
		return ""
	}
	templ := `
  # shellcheck disable=SC2154
  case $state in
    cmnds)
      # shellcheck disable=SC2034
      commands=(
{{range .Commands}}{{if not .Hidden}}        "{{.Name}}:{{if .Short}}{{.Short}}{{end}}"
{{end}}{{end}}      )
      _describe "command" commands
      ;;
  esac
  
  case "${words[1]}" in
{{range .Commands}}{{if not .Hidden}}    {{.Name}})
      {{uid .}}
      ;;
{{end}}{{end}}  esac
`
	cmd.Usage()
	buf := bytes.NewBufferString("")
	t, _ := template.New("subcommands").Funcs(template.FuncMap{"uid": uidCommand}).Parse(templ)
	t.Execute(buf, cmd)
	return buf.String()
}

func zshCompExtractFlag(c *cobra.Command) []*pflag.Flag {
	var flags []*pflag.Flag
	c.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			flags = append(flags, f)
		}
	})
	c.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			flags = append(flags, f)
		}
	})
	return flags
}

type ZshCompletion struct {
	cmd *cobra.Command
}

func Gen(cmd *cobra.Command) *ZshCompletion {
	addCompletionCommand(cmd)
	return &ZshCompletion{
		cmd: cmd,
	}
}

func (zsh ZshCompletion) PositionalCompletion(action ...Action) {
	for index, a := range action {
		completions.actions[uidPositional(zsh.cmd, index+1)] = a.finalize(uidPositional(zsh.cmd, index+1))
	}
}

func (zsh ZshCompletion) FlagCompletion(actions ActionMap) {
	for name, action := range actions {
		flag := zsh.cmd.Flag(name) // TODO only allowed for local flags
		completions.actions[uidFlag(zsh.cmd, flag)] = action.finalize(uidFlag(zsh.cmd, flag))
	}
}

var completions = Completions{
	actions: make(map[string]Action),
}

func addCompletionCommand(cmd *cobra.Command) {
	for _, c := range cmd.Root().Commands() {
		if c.Name() == "_zsh_completion" {
			return
		}
	}
	cmd.Root().AddCommand(&cobra.Command{
		Use:    "_zsh_completion",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 {
				fmt.Println(completions.Generate(cmd.Root()))
			} else {
				c, _, _ := cmd.Root().Find(parse(args[0])[1:]) // TODO check for errors

				cb := args[0]
				// TODO error on callback for flag without value (./example callback --callback)
				c.Run = func(cmd *cobra.Command, args []string) { // TODO flag parsing should be possible using Command.Traverse() instead of replacing the Run function
					if len(args) > 1 {
						fmt.Println(completions.invokeCallback(cb, args[1:]).Value)
					} else {
						fmt.Println(completions.invokeCallback(cb, []string{}).Value)
					}
				}

				origArg := parse(args[0]) // TODO messy
				if len(os.Args) > 3 {
					origArg = append(origArg, os.Args[3:]...)
				}

				os.Args = origArg
				cmd.Root().Execute()
			}
		},
		FParseErrWhitelist: cobra.FParseErrWhitelist{
			UnknownFlags: true,
		},
	})
}
