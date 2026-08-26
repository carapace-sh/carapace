package multi

import (
	"strings"
	"testing"
)

func TestSnippetBash(t *testing.T) {
	names := []string{"carapace-test", "sub1", "sub2"}
	s, err := Snippet("bash", names, "sub1", nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, "#!/bin/bash") {
		t.Error("missing bash shebang")
	}
	if !strings.Contains(s, `_sub1_completer`) {
		t.Error("missing completer function name")
	}
	if !strings.Contains(s, `"${command}"`) {
		t.Error("missing command extraction")
	}
	if !strings.Contains(s, `sub1`) {
		t.Error("missing sub1 in registration")
	}
	if !strings.Contains(s, `sub2`) {
		t.Error("missing sub2 in registration")
	}
}

func TestSnippetZsh(t *testing.T) {
	names := []string{"carapace-test", "sub1", "sub2"}
	s, err := Snippet("zsh", names, "sub1", nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, "#compdef") {
		t.Error("missing compdef header")
	}
	if !strings.Contains(s, `_sub1_completer`) {
		t.Error("missing completer function name")
	}
	if !strings.Contains(s, `"${command}"`) {
		t.Error("missing command extraction")
	}
}

func TestSnippetFish(t *testing.T) {
	names := []string{"carapace-test", "sub1", "sub2"}
	s, err := Snippet("fish", names, "sub1", nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, `function _sub1_completer`) {
		t.Error("missing completer function name")
	}
	if !strings.Contains(s, `$argv[1]`) {
		t.Error("missing argv reference")
	}
	// Verify each name gets its own -c registration (not all using defaultName)
	if !strings.Contains(s, `complete -c "sub2"`) {
		t.Error("sub2 should have its own -c registration")
	}
	if !strings.Contains(s, `complete -c "carapace-test"`) {
		t.Error("carapace-test should have its own -c registration")
	}
}

func TestSnippetElvish(t *testing.T) {
	names := []string{"carapace-test", "sub1", "sub2"}
	s, err := Snippet("elvish", names, "sub1", nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, "edit:completion:arg-completer") {
		t.Error("missing arg-completer")
	}
	if !strings.Contains(s, `$c _carapace elvish`) {
		t.Error("missing invocation with $c")
	}
}

func TestSnippetNushell(t *testing.T) {
	names := []string{"carapace-test", "sub1", "sub2"}
	s, err := Snippet("nushell", names, "sub1", nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, `sub1_completer`) {
		t.Error("missing completer name")
	}
	if !strings.Contains(s, `$spans.0`) {
		t.Error("missing spans.0")
	}
}

func TestSnippetPowershell(t *testing.T) {
	names := []string{"carapace-test", "sub1", "sub2"}
	s, err := Snippet("powershell", names, "sub1", nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, "System.Management.Automation") {
		t.Error("missing namespace")
	}
	if !strings.Contains(s, `$_sub1_completer`) {
		t.Error("missing completer variable")
	}
	if !strings.Contains(s, `($elems[0] -replace ('\.exe$', ''))`) {
		t.Error("missing exe replacement")
	}
}

func TestSnippetXonsh(t *testing.T) {
	names := []string{"carapace-test", "sub1", "sub2"}
	s, err := Snippet("xonsh", names, "sub1", nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, "contextual_command_completer") {
		t.Error("missing contextual_command_completer")
	}
	if !strings.Contains(s, `context.command`) {
		t.Error("missing context.command")
	}
}

func TestSnippetOil(t *testing.T) {
	names := []string{"carapace-test", "sub1", "sub2"}
	s, err := Snippet("oil", names, "sub1", nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, "#!/bin/osh") {
		t.Error("missing osh shebang")
	}
	if !strings.Contains(s, `_sub1_completer`) {
		t.Error("missing completer function name")
	}
}

func TestSnippetTcsh(t *testing.T) {
	names := []string{"carapace-test", "sub1", "sub2"}
	s, err := Snippet("tcsh", names, "sub1", nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, "complete \"carapace-test\"") {
		t.Error("missing carapace-test registration")
	}
	if !strings.Contains(s, "complete \"sub1\"") {
		t.Error("missing sub1 registration")
	}
}

func TestSnippetBashBle(t *testing.T) {
	names := []string{"carapace-test", "sub1", "sub2"}
	s, err := Snippet("bash-ble", names, "sub1", nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, "BLE_ATTACHED") {
		t.Error("missing BLE_ATTACHED check")
	}
	if !strings.Contains(s, "ble/complete/cand/yield") {
		t.Error("missing ble/complete/cand/yield")
	}
}

func TestSnippetUnsupported(t *testing.T) {
	_, err := Snippet("cmd-clink", []string{"test"}, "test", nil)
	if err == nil {
		t.Error("expected error for cmd-clink")
	}
	_, err = Snippet("ion", []string{"test"}, "test", nil)
	if err == nil {
		t.Error("expected error for ion")
	}
	_, err = Snippet("unknown", []string{"test"}, "test", nil)
	if err == nil {
		t.Error("expected error for unknown shell")
	}
}

func TestSnippetExport(t *testing.T) {
	s, err := Snippet("export", []string{"test"}, "test", nil)
	if err != nil {
		t.Error(err)
	}
	if s != "" {
		t.Error("export snippet should be empty")
	}
}

func TestSnippetFuncs(t *testing.T) {
	names := []string{"carapace-test", "sub1"}
	funcs := map[string][]string{
		"bash": {"# bash-enrichment"},
		"zsh":  {"# zsh-enrichment"},
	}
	s, err := Snippet("bash", names, "sub1", funcs)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, "# bash-enrichment") {
		t.Error("bash snippet missing enrichment code")
	}
	// zsh enrichment should not appear in bash snippet
	if strings.Contains(s, "# zsh-enrichment") {
		t.Error("bash snippet should not contain zsh enrichment")
	}
}

func TestSingleSnippetBash(t *testing.T) {
	s, err := SingleSnippet("bash", "sub1", []string{"carapace-test", "sub1"}, "sub1", nil)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(s, "#!/bin/bash") {
		t.Error("missing bash shebang")
	}
	if !strings.Contains(s, `_sub1_completion`) {
		t.Error("missing function name")
	}
	if !strings.Contains(s, `sub1 _carapace bash`) {
		t.Error("missing invocation")
	}
	if !strings.Contains(s, `complete -o noquote -F _sub1_completion sub1`) {
		t.Error("missing complete registration")
	}
}

func TestSingleSnippetAllShells(t *testing.T) {
	shells := []string{"bash", "bash-ble", "zsh", "fish", "elvish", "nushell", "powershell", "xonsh", "oil", "tcsh"}
	names := []string{"carapace-test", "sub1"}
	for _, sh := range shells {
		s, err := SingleSnippet(sh, "sub1", names, "sub1", nil)
		if err != nil {
			t.Errorf("%v: %v", sh, err)
		}
		if s == "" {
			t.Errorf("%v: empty snippet", sh)
		}
	}
}

func TestSnippetFuncsForShell(t *testing.T) {
	if s := snippetFuncsForShell("bash", nil); s != "" {
		t.Error("expected empty for nil map")
	}
	if s := snippetFuncsForShell("bash", map[string][]string{}); s != "" {
		t.Error("expected empty for empty map")
	}
	funcs := map[string][]string{
		"bash": {"line1", "line2"},
	}
	if s := snippetFuncsForShell("bash", funcs); s != "line1\nline2\n" {
		t.Errorf("unexpected: %q", s)
	}
	if s := snippetFuncsForShell("zsh", funcs); s != "" {
		t.Error("expected empty for missing shell key")
	}
}

func TestSnippetHyphenatedDefaultName(t *testing.T) {
	names := []string{"example-multi", "sub1"}
	// PowerShell and xonsh use underscore-prefixed identifiers
	shells := []string{"powershell", "xonsh"}
	for _, sh := range shells {
		s, err := Snippet(sh, names, "example-multi", nil)
		if err != nil {
			t.Errorf("%v: %v", sh, err)
		}
		if strings.Contains(s, `_example-multi_`) {
			t.Errorf("%v: hyphenated identifier not sanitized in snippet: %s", sh, s)
		}
		if !strings.Contains(s, `_example__multi_`) {
			t.Errorf("%v: missing sanitized identifier in snippet", sh)
		}
		if !strings.Contains(s, `example-multi`) {
			t.Errorf("%v: missing raw command name in registration", sh)
		}
	}

	// nushell uses variable names without underscore prefix
	s, err := Snippet("nushell", names, "example-multi", nil)
	if err != nil {
		t.Error(err)
	}
	if strings.Contains(s, `example-multi_completer`) {
		t.Error("nushell: hyphenated identifier not sanitized")
	}
	if !strings.Contains(s, `example__multi_completer`) {
		t.Error("nushell: missing sanitized identifier")
	}
}

func TestSingleSnippetHyphenatedCommand(t *testing.T) {
	shells := []string{"powershell", "xonsh"}
	names := []string{"example-multi", "sub1"}
	for _, sh := range shells {
		s, err := SingleSnippet(sh, "example-multi", names, "example-multi", nil)
		if err != nil {
			t.Errorf("%v: %v", sh, err)
		}
		if strings.Contains(s, `_example-multi_`) {
			t.Errorf("%v: hyphenated identifier not sanitized in single snippet", sh)
		}
		if !strings.Contains(s, `_example__multi_`) {
			t.Errorf("%v: missing sanitized identifier in single snippet", sh)
		}
		if !strings.Contains(s, `example-multi`) {
			t.Errorf("%v: missing raw command name in single snippet", sh)
		}
	}

	// nushell uses variable names without underscore prefix
	s, err := SingleSnippet("nushell", "example-multi", names, "example-multi", nil)
	if err != nil {
		t.Error(err)
	}
	if strings.Contains(s, `example-multi_completer`) {
		t.Error("nushell: hyphenated identifier not sanitized")
	}
	if !strings.Contains(s, `example__multi_completer`) {
		t.Error("nushell: missing sanitized identifier")
	}
}
