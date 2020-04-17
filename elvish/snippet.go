package elvish

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace/uid"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var replacer = strings.NewReplacer( // TODO
	`:`, `\:`,
	`"`, `\"`,
	`[`, `\[`,
	`]`, `\]`,
	`'`, `\"`,
)

func Snippet(cmd *cobra.Command, actions map[string]string) string {
	result := fmt.Sprintf(`edit:completion:arg-completer[%v] = [@arg]{
  fn _%v_callback [uid]{
    # TODO there is no 'eval' in elvish and '-source' needs a file so use a tempary one for callback 
    mkdir -p ~/.elvish/_carapace
    echo (joins ' ' $arg) | xargs %v _carapace elvish $uid > ~/.elvish/_carapace/_%v_callback.elv
    -source ~/.elvish/_carapace/_%v_callback.elv
    rm ~/.elvish/_carapace/_%v_callback.elv
  }

  fn subindex [subcommand]{
    # TODO 'edit:complete-getopt' needs the arguments shortened for subcommmands - pretty optimistic here
    index=1
    for x $arg { if (eq $x $subcommand) { break } else { index = (+ $index 1) } } 
    echo $index
  }
  
  state=(echo (joins ' ' $arg) | xargs %v _carapace elvish state)
  if (eq 1 0) {
  } %v
}
`, cmd.Name(), cmd.Name(), cmd.Name(), cmd.Name(), cmd.Name(), cmd.Name(), cmd.Name(), snippetFunctions(cmd, actions))

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
	if cmd.HasSubCommands() {
		subcommands := make([]string, 0)
		for _, c := range cmd.Commands() {
			if !c.Hidden {
				subcommands = append(subcommands, c.Name(), c.Short)
				for _, alias := range c.Aliases {
					subcommands = append(subcommands, alias, c.Short)
				}
			}
		}
		positionals = []string{"        " + snippetPositionalCompletion(ActionValuesDescribed(subcommands...))}
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
		if len(positionals) == 0 {
			if cmd.ValidArgs != nil {
				positionals = []string{"        " + snippetPositionalCompletion(ActionValues(cmd.ValidArgs...))}
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
	spec := []string{
		fmt.Sprintf(`&long=%v`, flag.Name),
		fmt.Sprintf(`&desc='%v'`, flag.Usage),
	}
	if flag.Shorthand != "" {
		spec = append(spec, fmt.Sprintf(`&short=%v`, flag.Shorthand))
	}
	if flag.NoOptDefVal == "" {
		spec = append(spec, `&arg-required=$true`, fmt.Sprintf(`&completer=[_]{ %v }`, action))
	}
	return fmt.Sprintf(`        [%v]`, strings.Join(spec, " "))
}
