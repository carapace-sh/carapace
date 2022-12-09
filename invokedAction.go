package carapace

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	_shell "github.com/rsteube/carapace/internal/shell"
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
	a.rawValues = common.RawValues(a.rawValues).Filter(values...)
	return a
}

// Merge merges InvokedActions (existing values are overwritten)
//
//	a := carapace.ActionValues("A", "B").Invoke(c)
//	b := carapace.ActionValues("B", "C").Invoke(c)
//	c := a.Merge(b) // ["A", "B", "C"]
func (a InvokedAction) Merge(others ...InvokedAction) InvokedAction {
	uniqueRawValues := make(map[string]common.RawValue)
	var meta common.Meta
	for _, other := range append([]InvokedAction{a}, others...) {
		for _, c := range other.rawValues {
			uniqueRawValues[c.Value] = c
		}
		meta.Merge(other.meta)
	}

	rawValues := make([]common.RawValue, 0, len(uniqueRawValues))
	for _, c := range uniqueRawValues {
		rawValues = append(rawValues, c)
	}

	invoked := InvokedAction{Action{rawValues: rawValues}}
	invoked.meta.Merge(meta)
	return invoked
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
	return ActionCallback(func(c Context) Action {
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

		splittedCV := _split(c.CallbackValue)
		for _, d := range dividers {
			if strings.HasSuffix(c.CallbackValue, d) {
				splittedCV = append(splittedCV, "")
				break
			}

		}

		uniqueVals := make(map[string]common.RawValue)
		for _, val := range a.rawValues {
			if strings.HasPrefix(val.Value, c.CallbackValue) {
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

		a := Action{rawValues: vals}
		for _, divider := range dividers {
			if runes := []rune(divider); len(runes) == 0 {
				a.meta.Nospace.Add('*')
				break
			} else {
				a.meta.Nospace.Add(runes[len(runes)-1])
			}
		}
		return a
	})
}

func (a InvokedAction) value(shell string, callbackValue string) string {
	return _shell.Value(shell, callbackValue, a.meta, a.rawValues)
}
