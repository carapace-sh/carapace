package zsh

import (
	"fmt"
	"strings"
	"github.com/spf13/pflag"
)

var replacer = strings.NewReplacer(
	`:`, `\:`,
	`"`, `\"`,
	`[`, `\[`,
	`]`, `\]`,
)

func snippetFlagCompletion(flag *pflag.Flag, action *Action) (snippet string) {
	var suffix, multimark, multimarkEscaped string
    if action == nil {
		if flag.NoOptDefVal != "" {
			suffix = "" // no argument required for flag
		} else {
			suffix = ": :" // require a value
		}
	} else {
		suffix = fmt.Sprintf(": :%v", action.Value)
	}

    if(zshCompFlagCouldBeSpecifiedMoreThenOnce(flag)){
      multimark = "*"
      multimarkEscaped = "\\*"
    }

	// TODO flag without value (without ": :%v") -> flag.NoOptDefaultVal - when nothing configured annd empty flag requires an value
	if flag.Shorthand == "" { // no shorthannd
		snippet = fmt.Sprintf(`"%v--%v[%v]%v"`, multimark, flag.Name, replacer.Replace(flag.Usage), suffix)
	} else {
		snippet = fmt.Sprintf(`"(%v-%v %v--%v)"{%v-%v,%v--%v}"[%v]%v"`, multimark, flag.Shorthand, multimark, flag.Name, multimarkEscaped,flag.Shorthand, multimarkEscaped,flag.Name, replacer.Replace(flag.Usage), suffix)
	}
	return
}

func snippetPositionalCompletion(position int, action Action) string {
	return fmt.Sprintf(`"%v:: :%v" \`+"\n", position, action.Value)
}

func zshCompFlagCouldBeSpecifiedMoreThenOnce(f *pflag.Flag) bool {
	return strings.Contains(f.Value.Type(), "Slice") ||
		strings.Contains(f.Value.Type(), "Array")
}
