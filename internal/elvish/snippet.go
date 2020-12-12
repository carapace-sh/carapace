package elvish

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/uid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var replacer = strings.NewReplacer( // TODO
	`:`, `\:`,
	"\n", ``,
	`"`, `\"`,
	`[`, `\[`,
	`]`, `\]`,
	`'`, `\"`,
)

func snippetLazy(cmd *cobra.Command) string {
	return fmt.Sprintf(`edit:completion:arg-completer[%v] = [@arg]{
    edit:completion:arg-completer[%v] = [@arg]{}
    eval (%v _carapace elvish | slurp)
    $edit:completion:arg-completer[%v] $@arg
}
`, cmd.Name(), cmd.Name(), uid.Executable(), cmd.Name())
}

func Snippet(cmd *cobra.Command, actions map[string]string, lazy bool) string {
	if lazy {
		return snippetLazy(cmd)
	}

	result := fmt.Sprintf(`use str
edit:completion:arg-completer[%v] = [@arg]{
  fn _%v_callback [uid]{
    if (eq $arg[-1] "") {
        arg[-1] = "''"
    }
    eval (echo (str:join "\001" $arg) | xargs --delimiter="\001" %v _carapace elvish $uid | slurp) &ns=(ns [&arg=$arg])
  }

  fn subindex [subcommand]{
    # TODO 'edit:complete-getopt' needs the arguments shortened for subcommmands - pretty optimistic here
    index=1
    for x $arg { if (eq $x $subcommand) { break } else { index = (+ $index 1) } } 
    echo $index
  }
  
  state=(echo (str:join ' ' $arg) | xargs %v _carapace elvish state)
  if (eq 1 0) {
  } %v
}
`, cmd.Name(), cmd.Name(), uid.Executable(), uid.Executable(), snippetFunctions(cmd, actions))

	return result
}

func snippetFunctions(cmd *cobra.Command, actions map[string]string) string {
	function_pattern := ` elif (eq $state '%v') {
    opt-specs = [
%v
    ]
    arg-handlers = [
%v
    ]
    subargs = $arg[(subindex %v):] 
    if (> (count $subargs) 0) {
      edit:complete-getopt $subargs $opt-specs $arg-handlers
    }
  }`

	flags := make([]string, 0)
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if !f.Hidden {
			var s string
			if action, ok := actions[uid.Flag(cmd, f)]; ok {
				s = snippetFlagCompletion(f, action)
			} else {
				s = snippetFlagCompletion(f, "")
			}
			flags = append(flags, s)
		}
	})

	var positionals []string
	if cmd.HasAvailableSubCommands() {
		subcommands := make([]common.Candidate, 0)
		for _, c := range cmd.Commands() {
			if !c.Hidden {
				subcommands = append(subcommands, common.Candidate{Value: c.Name(), Display: c.Name(), Description: c.Short})
				for _, alias := range c.Aliases {
					subcommands = append(subcommands, common.Candidate{Value: alias, Display: alias, Description: c.Short})
				}
			}
		}
		positionals = []string{"        " + snippetPositionalCompletion(ActionCandidates(subcommands...))}
	} else {
		pos := 1
		for {
			if action, ok := actions[uid.Positional(cmd, pos)]; ok {
				positionals = append(positionals, "      "+snippetPositionalCompletion(action))
				pos++
			} else {
				break // TODO only consistent entries for now
			}
		}
		if action, ok := actions[uid.Positional(cmd, 0)]; ok {
			positionals = append(positionals, "      "+snippetPositionalCompletion(action))
			positionals = append(positionals, "      "+"...")
		}
		if len(positionals) == 0 {
			if cmd.ValidArgs != nil {
				positionals = []string{"        " + snippetPositionalCompletion(ActionCandidates(common.CandidateFromValues(cmd.ValidArgs...)...))}
			}
		}
	}

	result := make([]string, 0)
	result = append(result, fmt.Sprintf(function_pattern, uid.Command(cmd), strings.Join(flags, "\n"), strings.Join(positionals, "\n"), cmd.Name()))
	for _, subcmd := range cmd.Commands() {
		if !subcmd.Hidden {
			result = append(result, snippetFunctions(subcmd, actions))
		}
	}
	return strings.Join(result, " ")
}

func snippetPositionalCompletion(action string) string {
	return fmt.Sprintf(`[_]{ %v }`, action)
}

func snippetFlagCompletion(flag *pflag.Flag, action string) (snippet string) {
	spec := []string{}
	if !common.IsShorthandOnly(flag) {
		spec = append(spec, fmt.Sprintf(`&long='%v'`, flag.Name))
	}
	if flag.Shorthand != "" {
		spec = append(spec, fmt.Sprintf(`&short='%v'`, flag.Shorthand))
	}

	spec = append(spec, fmt.Sprintf(`&desc='%v'`, replacer.Replace(flag.Usage)))

	if flag.NoOptDefVal == "" {
		spec = append(spec, `&arg-required=$true`, fmt.Sprintf(`&completer=[_]{ %v }`, action))
	} else {
		spec = append(spec, `&arg-optional=$true`, fmt.Sprintf(`&completer=[_]{ %v }`, action))
	}
	return fmt.Sprintf(`        [%v]`, strings.Join(spec, " "))
}
