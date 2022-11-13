package carapace

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/shell/bash"
	"github.com/rsteube/carapace/internal/shell/bash_ble"
	"github.com/rsteube/carapace/internal/shell/elvish"
	"github.com/rsteube/carapace/internal/shell/export"
	"github.com/rsteube/carapace/internal/shell/fish"
	"github.com/rsteube/carapace/internal/shell/ion"
	"github.com/rsteube/carapace/internal/shell/nushell"
	"github.com/rsteube/carapace/internal/shell/oil"
	"github.com/rsteube/carapace/internal/shell/powershell"
	"github.com/rsteube/carapace/internal/shell/tcsh"
	"github.com/rsteube/carapace/internal/shell/xonsh"
	"github.com/rsteube/carapace/internal/shell/zsh"
)

// InvokedAction is a logical alias for an Action whose (nested) callback was invoked.
type InvokedAction struct {
	Action
}

// Filter filters given values (this should be done before any call to Prefix/Suffix as those alter the values being filtered)
//
//	a := carapace.ActionValues("A", "B", "C").Invoke(c)
//	b := a.Filter([]string{"B"}) // ["A", "C"]
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
	return InvokedAction{actionRawValues(filtered...).noSpace(a.nospace).skipCache(a.skipcache)}
}

// Merge merges InvokedActions (existing values are overwritten)
//
//	a := carapace.ActionValues("A", "B").Invoke(c)
//	b := carapace.ActionValues("B", "C").Invoke(c)
//	c := a.Merge(b) // ["A", "B", "C"]
func (a InvokedAction) Merge(others ...InvokedAction) InvokedAction {
	var orderedValues []string
	uniqueRawValues := make(map[string]common.RawValue)
	nospace := a.nospace
	skipcache := a.skipcache
	for _, other := range append([]InvokedAction{a}, others...) {
		for _, c := range other.rawValues {
			_, exists := uniqueRawValues[c.Value]
			uniqueRawValues[c.Value] = c
			if !exists {
				orderedValues = append(orderedValues, c.Value)
			}
		}
		nospace = nospace || other.nospace
		skipcache = skipcache || other.skipcache
	}

	rawValues := make([]common.RawValue, 0, len(uniqueRawValues))
	for _, v := range orderedValues {
		c := uniqueRawValues[v]
		rawValues = append(rawValues, c)
	}
	return InvokedAction{actionRawValues(rawValues...).noSpace(nospace).skipCache(skipcache).withHint(a.hint)}
}

// Prefix adds a prefix to values (only the ones inserted, not the display values)
//
//	a := carapace.ActionValues("melon", "drop", "fall").Invoke(c)
//	b := a.Prefix("water") // ["watermelon", "waterdrop", "waterfall"] but display still ["melon", "drop", "fall"]
func (a InvokedAction) Prefix(prefix string) InvokedAction {
	for index, val := range a.rawValues {
		a.rawValues[index].Value = prefix + val.Value
	}
	return a
}

// Suffix adds a suffx to values (only the ones inserted, not the display values)
//
//	a := carapace.ActionValues("apple", "melon", "orange").Invoke(c)
//	b := a.Suffix("juice") // ["applejuice", "melonjuice", "orangejuice"] but display still ["apple", "melon", "orange"]
func (a InvokedAction) Suffix(suffix string) InvokedAction {
	for index, val := range a.rawValues {
		a.rawValues[index].Value = val.Value + suffix
	}
	return a
}

// ToA casts an InvokedAction to Action.
func (a InvokedAction) ToA() Action {
	return a.Action
}

// ToMultiPartsA create an ActionMultiParts from values with given dividers
//
//	a := carapace.ActionValues("A/B/C", "A/C", "B/C", "C").Invoke(c)
//	b := a.ToMultiPartsA("/") // completes segments separately (first one is ["A/", "B/", "C"])
func (a InvokedAction) ToMultiPartsA(dividers ...string) Action {
	return ActionCallback(func(ctx Context) Action {
		_split := func() func(s string) []string {
			quotedDividiers := make([]string, 0)
			for _, d := range dividers {
				quotedDividiers = append(quotedDividiers, regexp.QuoteMeta(d))
			}
			f := fmt.Sprintf("([^%v]*(%v)?)", strings.Join(quotedDividiers, "|"), strings.Join(quotedDividiers, "|")) // TODO quickfix - this is wrong (fails for dividers longer than one character) an might need a reverse lookahead for character sequence
			r := regexp.MustCompile(f)

			return func(s string) []string {
				if matches := r.FindAllString(s, -1); matches != nil {
					return matches
				}

				return []string{}
			}
		}()

		splittedCV := _split(ctx.CallbackValue)
		for _, d := range dividers {
			if strings.HasSuffix(ctx.CallbackValue, d) {
				splittedCV = append(splittedCV, "")

				break
			}
		}

		uniqueVals := make(map[string]common.RawValue)
		for _, val := range a.rawValues {
			if strings.HasPrefix(val.Value, ctx.CallbackValue) {
				if splitted := _split(val.Value); len(splitted) >= len(splittedCV) {
					v := strings.Join(splitted[:len(splittedCV)], "")
					d := splitted[len(splittedCV)-1]

					if len(splitted) == len(splittedCV) {
						uniqueVals[v] = common.RawValue{
							Value:       v,
							Display:     d,
							Description: val.Description,
							Style:       val.Style,
						}
					} else {
						uniqueVals[v] = common.RawValue{
							Value:       v,
							Display:     d,
							Description: "",
							Style:       "",
						}
					}
				}
			}
		}

		vals := make([]common.RawValue, 0)
		for _, val := range uniqueVals {
			vals = append(vals, val)
		}

		return actionRawValues(vals...).noSpace(true)
	})
}

// onIninalize can take some steps to make everything read for all shells.
func (a InvokedAction) onInitialize(callbackValue string) InvokedAction {
	// The final hint present is used by some shells
	common.CompletionHint = a.hint
	common.CompletionMessage = a.message
	return a
}

func (a InvokedAction) value(shell string, callbackValue string) string { // TODO use context instead?
	a = a.onInitialize(callbackValue)

	shellFuncs := map[string]func(currentWord string, nospace bool, values common.RawValues) string{
		"bash":       bash.ActionRawValues,
		"bash-ble":   bash_ble.ActionRawValues,
		"fish":       fish.ActionRawValues,
		"elvish":     elvish.ActionRawValues,
		"export":     export.ActionRawValues,
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
