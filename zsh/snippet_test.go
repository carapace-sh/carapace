package zsh

import (
	"testing"

	"github.com/rsteube/carapace/assert"
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

	assert.Equal(t, `"--simple[simple flag]"`, snippetFlagCompletion(root.Flag("simple"), nil))

	valuesAction := ActionValues("a", "b", "c")
	assert.Equal(t, `"--values[values action flag]: :_values '' a b c"`, snippetFlagCompletion(root.Flag("values"), &valuesAction))

	assert.Equal(t, `"(-s --shorthand)"{-s,--shorthand}"[shorthand flag]"`, snippetFlagCompletion(root.Flag("shorthand"), nil))

	assert.Equal(t, `"(*-a *--stringarray)"{\*-a,\*--stringarray}"[stringarray flag]: :"`, snippetFlagCompletion(root.Flag("stringarray"), nil))
	assert.Equal(t, `"*--stringslice[stringslice flag]: :"`, snippetFlagCompletion(root.Flag("stringslice"), nil))
}

func TestSnippetPositionalCompletion(t *testing.T) {
	pos1 := snippetPositionalCompletion(1, ActionValues("a", "b", "c"))
	assert.Equal(t, `"1:: :_values '' a b c"`, pos1)

	pos2 := snippetPositionalCompletion(2, ActionMessage("test"))
	assert.Equal(t, `"2:: : _message -r 'test'"`, pos2)
}

func TestSnippetSubcommands(t *testing.T) {
	root := &cobra.Command{
		Use: "root",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	sub1 := &cobra.Command{
		Use:     "sub1",
		Run:     func(cmd *cobra.Command, args []string) {},
		Aliases: []string{"alias1", "alias2"},
	}
	sub2 := &cobra.Command{
		Use:   "sub2",
		Run:   func(cmd *cobra.Command, args []string) {},
		Short: "short description",
	}
	hidden := &cobra.Command{
		Use:    "hidden",
		Run:    func(cmd *cobra.Command, args []string) {},
		Hidden: true,
	}
	root.AddCommand(sub1)
	root.AddCommand(sub2)
	root.AddCommand(hidden)

	expected := `

  # shellcheck disable=SC2154
  case $state in
    cmnds)
      # shellcheck disable=SC2034
      commands=(
        "sub1:"
        "alias1:"
        "alias2:"
        "sub2:short description"
      )
      _describe "command" commands
      ;;
  esac
  
  case "${words[1]}" in
    sub1)
      _root__sub1
      ;;
    alias1)
      _root__sub1
      ;;
    alias2)
      _root__sub1
      ;;
    sub2)
      _root__sub2
      ;;
  esac`
	assert.Equal(t, expected, snippetSubcommands(root))
}
