package zsh

import (
	"testing"
)

func TestActionCallback(t *testing.T) {
	a := ActionCallback(func(args []string) Action {
		return ActionMessage("ActionCallback test")
	}).finalize("someId")

	assertEqual(t, ` eval \$(${os_args[1]} _zsh_completion 'someId' ${${os_args:1:gs/\"/\\\"}:gs/\'/\\\"})`, a.Value)
}

func TestActionExecute(t *testing.T) {
	assertEqual(t, ` eval \$(echo test)`, ActionExecute("echo test").Value)
}

func TestActionBool(t *testing.T) {
	assertEqual(t, `_values '' true false`, ActionBool().Value)
}

func TestActionPathFiles(t *testing.T) {
	assertEqual(t, `_path_files`, ActionPathFiles("").Value)
	assertEqual(t, `_path_files -g '*.go'`, ActionPathFiles("*.go").Value)
}

func TestActionFiles(t *testing.T) {
	assertEqual(t, `_files`, ActionFiles("").Value)
	assertEqual(t, `_files -g '*.go'`, ActionFiles("*.go").Value)
}

func TestActionNetInterfaces(t *testing.T) {
	assertEqual(t, `_net_interfaces`, ActionNetInterfaces().Value)
}

func TestActionUsers(t *testing.T) {
	assertEqual(t, `_users`, ActionUsers().Value)
}

func TestActionGroups(t *testing.T) {
	assertEqual(t, `_groups`, ActionGroups().Value)
}

func TestActionHosts(t *testing.T) {
	assertEqual(t, `_hosts`, ActionHosts().Value)
}

func TestActionOptions(t *testing.T) {
	assertEqual(t, `_options`, ActionOptions().Value)
}

func TestActionValues(t *testing.T) {
	assertEqual(t, ` _message -r 'no values to complete'`, ActionValues().Value)
	assertEqual(t, `_values '' a b`, ActionValues("a", "b").Value)
}

func TestActionValuesDescribed(t *testing.T) {
	assertEqual(t, ` _message -r 'no values to complete'`, ActionValuesDescribed().Value)
	assertEqual(t, `_values '' 'a[aDescription]' 'b[bDescription]'  `, ActionValuesDescribed("a", "aDescription", "b", "bDescription").Value)
}

func TestActionMessage(t *testing.T) {
	assertEqual(t, ` _message -r 'test'`, ActionMessage("test").Value)
}

func TestMultiParts(t *testing.T) {
	assertEqual(t, `_multi_parts % '(a%b a%b%c)'`, ActionMultiParts('%', "a%b", "a%b%c").Value)
}
