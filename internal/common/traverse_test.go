package common

import (
	"testing"

	"github.com/rsteube/carapace/internal/assert"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use: "test",
}

var testSubCmd = &cobra.Command{
	Use: "sub",
}

var testSub2Cmd = &cobra.Command{
	Use:                "sub2",
	DisableFlagParsing: true,
}

var intVar int

func init() {
	testCmd.Flags().StringP("optarg", "o", "", "optarg flag")
	testCmd.Flags().StringP("string", "s", "", "string flag")
	testCmd.Flags().BoolP("bool", "b", false, "bool flag")
	testCmd.Flags().IntP("int", "i", 0, "int flag")
	testCmd.Flags().IntVar(&intVar, "intvar", 0, "int flag")
	testCmd.Flags().IntSlice("intslice", []int{}, "intlice flag")

	testSubCmd.Flags().StringP("suboptarg", "o", "", "optarg flag")
	testSubCmd.Flags().StringP("substring", "s", "", "string flag")
	testSubCmd.Flags().BoolP("subbool", "b", false, "bool flag")
	testCmd.AddCommand(testSubCmd)

	testCmd.AddCommand(testSub2Cmd)
}

func testTraverseLenient(t *testing.T, args ...string) (targetCmd *cobra.Command, targetArgs []string, err error) {
	if targetCmd, targetArgs, err = TraverseLenient(testCmd, args); err != nil {
		t.Error(err.Error())
	}
	return
}

func TestNoArg(t *testing.T) {
	testTraverseLenient(t)
}

func TestEmptyArg(t *testing.T) {
	testTraverseLenient(t, "")
}

func TestLongFlagIncomplete(t *testing.T) {
	testTraverseLenient(t, "--stri")
}

func TestLongFlagNoArg(t *testing.T) {
	testTraverseLenient(t, "--string")
}

func TestLongFlagEmptyArg(t *testing.T) {
	testTraverseLenient(t, "--string", "")
}

func TestShortFlagNoArg(t *testing.T) {
	testTraverseLenient(t, "-s")
}

func TestShortFlagEmptyArg(t *testing.T) {
	testTraverseLenient(t, "-s", "")
}

func TestShortFlagChainNoArg(t *testing.T) {
	testTraverseLenient(t, "-bs")
}

func TestShortFlagChainEmptyArg(t *testing.T) {
	testTraverseLenient(t, "-bs", "")
}

func TestIntEmptyArg(t *testing.T) {
	testTraverseLenient(t, "--int", "")
}

func TestIntVarEmptyArg(t *testing.T) {
	testTraverseLenient(t, "--intvar", "")
}

func TestIntSliceEmptyArg(t *testing.T) {
	testTraverseLenient(t, "--intslice", "")
}

func TestIntSliceIncompleteArg(t *testing.T) {
	testTraverseLenient(t, "--intslice", "1,")
}

func TestDisabledFlagParsing(t *testing.T) {
	_, args, _ := testTraverseLenient(t, "sub2", "--arg")
	assert.Equal(t, args[0], "--arg") // TODO test whole slice
}

func TestFlagValueMatchesCmd(t *testing.T) {
	cmd, args, _ := testTraverseLenient(t, "--string", testCmd.Name())
	if cmd != testCmd {
		t.Errorf("should be testCmd [was: '%v']", cmd)
	}
	if len(args) != 0 {
		t.Errorf("args should be empty [was: '%v']", len(args))
	}
	if v := cmd.Flag("string").Value.String(); v != testCmd.Use {
		t.Errorf("string flag should be set [expected: '%v', was: '%v']", testCmd.Name(), v)
	}
}

func TestFlagValueMatchesSubCmd(t *testing.T) {
	cmd, args, _ := testTraverseLenient(t, "sub", "--substring", testSubCmd.Name())
	if cmd != testSubCmd {
		t.Errorf("should be testSubCmd [was: '%v']", cmd.Name())
	}
	if len(args) != 0 {
		t.Errorf("args should be empty [was: '%v']", len(args))
	}
	if v := testSubCmd.Flag("substring").Value.String(); v != testSubCmd.Use {
		t.Errorf("string flag should be set [expected: '%v', was: '%v']", testSubCmd.Name(), v)
	}
}
