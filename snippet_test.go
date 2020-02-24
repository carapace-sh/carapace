package zsh

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestSnippetFlagCompletion(t *testing.T) {
	root := &cobra.Command{
		Use: "root",
	}
	root.Flags().Bool("simple", false, "simple flag")
	root.Flags().String("values", "b", "values action flag")
	root.Flags().BoolP("shorthand", "s", false, "shorthand flag")
	root.Flags().StringArrayP("stringarray", "a", []string{"a"}, "stringarray flag")
	root.Flags().StringSlice("stringslice", []string{"a"}, "stringslice flag")

	assertEqual(t, `"--simple[simple flag]"`, snippetFlagCompletion(root.Flag("simple"), nil))

	valuesAction := ActionValues("a", "b", "c")
	assertEqual(t, `"--values[values action flag]: :_values '' a b c"`, snippetFlagCompletion(root.Flag("values"), &valuesAction))

	assertEqual(t, `"(-s --shorthand)"{-s,--shorthand}"[shorthand flag]"`, snippetFlagCompletion(root.Flag("shorthand"), nil))

	assertEqual(t, `"(*-a *--stringarray)"{\*-a,\*--stringarray}"[stringarray flag]: :"`, snippetFlagCompletion(root.Flag("stringarray"), nil))
	assertEqual(t, `"*--stringslice[stringslice flag]: :"`, snippetFlagCompletion(root.Flag("stringslice"), nil))
}

func TestSnippetPositionalCompletion(t *testing.T) {
	pos1 := snippetPositionalCompletion(1, ActionValues("a", "b", "c"))
	assertEqual(t, `"1:: :_values '' a b c"`, pos1)

	pos2 := snippetPositionalCompletion(2, ActionMessage("test"))
	assertEqual(t, `"2:: : _message -r 'test'"`, pos2)
}
