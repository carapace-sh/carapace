package common

import (
	"testing"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use: "test",
}

var testSubCmd = &cobra.Command{
	Use: "sub",
}

func init() {
	testCmd.Flags().StringP("optarg", "o", "", "optarg flag")
	testCmd.Flags().StringP("string", "s", "", "string flag")
	testCmd.Flags().BoolP("bool", "b", false, "bool flag")

	testSubCmd.Flags().StringP("suboptarg", "o", "", "optarg flag")
	testSubCmd.Flags().StringP("substring", "s", "", "string flag")
	testSubCmd.Flags().BoolP("subbool", "b", false, "bool flag")
	testCmd.AddCommand(testSubCmd)
}

func testTraverseLenient(t *testing.T, args ...string) {
	if _, _, err := TraverseLenient(testCmd, args); err != nil {
		t.Error(err.Error())
	}
}

func TestNoArg(t *testing.T) {
	testTraverseLenient(t)
}

func TestEmptyArg(t *testing.T) {
	testTraverseLenient(t, "")
}

func TestLongFlagIncomplete(t *testing.T) {
	testTraverseLenient(t, "test", "--stri")
}

func TestLongFlagNoArg(t *testing.T) {
	testTraverseLenient(t, "test", "--string")
}

func TestLongFlagEmptyArg(t *testing.T) {
	testTraverseLenient(t, "test", "--string", "")
}

func TestShortFlagNoArg(t *testing.T) {
	testTraverseLenient(t, "test", "-s")
}

func TestShortFlagEmptyArg(t *testing.T) {
	testTraverseLenient(t, "test", "-s", "")
}

func TestShortFlagChainNoArg(t *testing.T) {
	testTraverseLenient(t, "test", "-bs")
}

func TestShortFlagChainEmptyArg(t *testing.T) {
	testTraverseLenient(t, "test", "-bs", "")
}
