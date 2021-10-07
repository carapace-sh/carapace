package carapace

import (
	"strings"

	"github.com/rsteube/carapace/internal/bash"
	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/elvish"
	"github.com/rsteube/carapace/internal/fish"
	"github.com/rsteube/carapace/internal/ion"
	"github.com/rsteube/carapace/internal/nushell"
	"github.com/rsteube/carapace/internal/oil"
	"github.com/rsteube/carapace/internal/powershell"
	"github.com/rsteube/carapace/internal/tcsh"
	"github.com/rsteube/carapace/internal/xonsh"
	"github.com/rsteube/carapace/internal/zsh"
)

// InvokedAction is a logical alias for an Action whose (nested) callback was invoked
type InvokedAction Action

// Filter filters given values (this should be done before any call to Prefix/Suffix as those alter the values being filtered)
//   a := carapace.ActionValues("A", "B", "C").Invoke(c)
//   b := a.Filter([]string{"B"}) // ["A", "C"]
func (a InvokedAction) Filter(values []string) InvokedAction {
	toremove := make(map[string]bool)
	for _, v := range values {
		toremove[v] = true
	}
	filtered := make([]common.RawValue, 0)
	for _, rawValue := range a.rawValues {
		if _, ok := toremove[rawValue.Value]; !ok {
			filtered = append(filtered, rawValue)
		}
	}
	return InvokedAction(actionRawValues(filtered...).noSpace(a.nospace).skipCache(a.skipcache))
}

// Merge merges InvokedActions (existing values are overwritten)
//   a := carapace.ActionValues("A", "B").Invoke(c)
//   b := carapace.ActionValues("B", "C").Invoke(c)
//   c := a.Merge(b) // ["A", "B", "C"]
func (a InvokedAction) Merge(others ...InvokedAction) InvokedAction {
	uniqueRawValues := make(map[string]common.RawValue)
	nospace := a.nospace
	skipcache := a.skipcache
	for _, other := range append([]InvokedAction{a}, others...) {
		for _, c := range other.rawValues {
			uniqueRawValues[c.Value] = c
		}
		nospace = a.nospace || other.nospace
		skipcache = a.skipcache || other.skipcache
	}

	rawValues := make([]common.RawValue, 0, len(uniqueRawValues))
	for _, c := range uniqueRawValues {
		rawValues = append(rawValues, c)
	}
	return InvokedAction(actionRawValues(rawValues...).noSpace(nospace).skipCache(skipcache))
}

// Prefix adds a prefix to values (only the ones inserted, not the display values)
//   a := carapace.ActionValues("melon", "drop", "fall").Invoke(c)
//   b := a.Prefix("water") // ["watermelon", "waterdrop", "waterfall"] but display still ["melon", "drop", "fall"]
func (a InvokedAction) Prefix(prefix string) InvokedAction {
	for index, val := range a.rawValues {
		a.rawValues[index].Value = prefix + val.Value
	}
	return a
}

// Suffix adds a suffx to values (only the ones inserted, not the display values)
//   a := carapace.ActionValues("apple", "melon", "orange").Invoke(c)
//   b := a.Suffix("juice") // ["applejuice", "melonjuice", "orangejuice"] but display still ["apple", "melon", "orange"]
func (a InvokedAction) Suffix(suffix string) InvokedAction {
	for index, val := range a.rawValues {
		a.rawValues[index].Value = val.Value + suffix
	}
	return a
}

// ToA casts an InvokedAction to Action
func (a InvokedAction) ToA() Action {
	return Action(a)
}

// ToMultiPartsA create an ActionMultiParts from values with given divider
//   a := carapace.ActionValues("A/B/C", "A/C", "B/C", "C").Invoke(c)
//   b := a.ToMultiPartsA("/") // completes segments separately (first one is ["A/", "B/", "C"])
func (a InvokedAction) ToMultiPartsA(divider string) Action {
	return ActionMultiParts(divider, func(c Context) Action {
		uniqueVals := make(map[string]string)
		for _, val := range a.rawValues {
			if strings.HasPrefix(val.Value, strings.Join(c.Parts, divider)) {
				if splitted := strings.Split(val.Value, divider); len(splitted) > len(c.Parts) {
					if len(splitted) == len(c.Parts)+1 {
						uniqueVals[splitted[len(c.Parts)]] = val.Description
					} else {
						uniqueVals[splitted[len(c.Parts)]+divider] = ""
					}
				}
			}
		}

		vals := make([]string, 0, len(uniqueVals)*2)
		for val, description := range uniqueVals {
			vals = append(vals, val, description)
		}
		return ActionValuesDescribed(vals...).noSpace(true)
	})
}

func (a InvokedAction) value(shell string, callbackValue string) string { // TODO use context instead?
	shellFuncs := map[string]func(currentWord string, nospace bool, values common.RawValues) string{
		"bash":       bash.ActionRawValues,
		"fish":       fish.ActionRawValues,
		"elvish":     elvish.ActionRawValues,
		"ion":        ion.ActionRawValues,
		"nushell":    nushell.ActionRawValues,
		"oil":        oil.ActionRawValues,
		"powershell": powershell.ActionRawValues,
		"tcsh":       tcsh.ActionRawValues,
		"xonsh":      xonsh.ActionRawValues,
		"zsh":        zsh.ActionRawValues,
	}
	if f, ok := shellFuncs[shell]; ok {
		return f(callbackValue, a.nospace, a.rawValues)
	}
	return ""
}
