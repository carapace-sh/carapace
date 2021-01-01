package carapace

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/rsteube/carapace/internal/bash"
	"github.com/rsteube/carapace/internal/cache"
	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/elvish"
	"github.com/rsteube/carapace/internal/fish"
	"github.com/rsteube/carapace/internal/oil"
	"github.com/rsteube/carapace/internal/powershell"
	"github.com/rsteube/carapace/internal/xonsh"
	"github.com/rsteube/carapace/internal/zsh"
	pkgcache "github.com/rsteube/carapace/pkg/cache"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Action indicates how to complete a flag or positional argument
type Action struct {
	rawValues  []common.RawValue
	bash       func() string
	elvish     func() string
	fish       func() string
	oil        func() string
	powershell func() string
	xonsh      func() string
	zsh        func() string
	callback   CompletionCallback
}

type ActionMap map[string]Action

func (a ActionMap) invokeCallback(uid string, args []string) Action {
	if action, ok := a[uid]; ok {
		if action.callback != nil {
			return action.callback(args)
		}
	}
	return ActionMessage(fmt.Sprintf("callback %v unknown", uid))
}

func (a ActionMap) shell(shell string) map[string]string {
	actions := make(map[string]string, len(a))
	for key, value := range map[string]Action(a) {
		actions[key] = value.Value(shell)
	}
	return actions
}

type CompletionCallback func(args []string) Action

func (a Action) Cache(timeout time.Duration, keys ...pkgcache.CacheKey) Action {
	if a.callback != nil { // only relevant for callback actions
		cachedCallback := a.callback
		_, file, line, _ := runtime.Caller(1) // generate uid from wherever Cache() was called
		a.callback = func(args []string) Action {
			if cacheFile, err := cache.File(file, line, keys...); err == nil {
				if rawValues, err := cache.Load(cacheFile, timeout); err == nil {
					return actionRawValues(rawValues...)
				} else {
					oldState := skipCache
					skipCache = false // TODO find a better solution for this
					invokedAction := (Action{callback: cachedCallback}).Invoke(args)
					if !skipCache {
						_ = cache.Write(cacheFile, invokedAction.rawValues)
					}
					skipCache = oldState || skipCache
					return invokedAction.ToA()
				}
			}
			return cachedCallback(args)
		}
	}
	return a
}

type InvokedAction Action

// Invoke executes the callback of an action if it exists (supports nesting)
func (a Action) Invoke(args []string) InvokedAction {
	return InvokedAction(a.nestedAction(args, 5))
}

func (a InvokedAction) Merge(others ...InvokedAction) InvokedAction {
	uniqueRawValues := make(map[string]common.RawValue)
	for _, other := range append([]InvokedAction{a}, others...) {
		for _, c := range other.rawValues {
			uniqueRawValues[c.Value] = c
		}
	}

	rawValues := make([]common.RawValue, 0, len(uniqueRawValues))
	for _, c := range uniqueRawValues {
		rawValues = append(rawValues, c)
	}
	return InvokedAction(actionRawValues(rawValues...))
}

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
	return InvokedAction(actionRawValues(filtered...))
}

func (a InvokedAction) Prefix(prefix string) InvokedAction {
	for index, val := range a.rawValues {
		a.rawValues[index].Value = prefix + val.Value
	}
	return a
}

func (a InvokedAction) Suffix(suffix string) InvokedAction {
	for index, val := range a.rawValues {
		a.rawValues[index].Value = val.Value + suffix
	}
	return a
}

func (a InvokedAction) ToA() Action {
	return Action(a)
}

func (a InvokedAction) ToMultiPartsA(divider string) Action {
	return ActionMultiParts(divider, func(args, parts []string) Action {
		uniqueVals := make(map[string]string)
		for _, val := range a.rawValues {
			if strings.HasPrefix(val.Value, strings.Join(parts, divider)) {
				if splitted := strings.Split(val.Value, divider); len(splitted) > len(parts) {
					if len(splitted) == len(parts)+1 {
						uniqueVals[splitted[len(parts)]] = val.Description
					} else {
						uniqueVals[splitted[len(parts)]+divider] = ""
					}
				}
			}
		}

		vals := make([]string, 0, len(uniqueVals)*2)
		for val, description := range uniqueVals {
			vals = append(vals, val, description)
		}

		return ActionValuesDescribed(vals...)
	})
}

func (a Action) nestedAction(args []string, maxDepth int) Action {
	if a.rawValues == nil && a.callback != nil && maxDepth > 0 {
		return a.callback(args).nestedAction(args, maxDepth-1)
	} else {
		return a
	}
}

func (a Action) Value(shell string) string {
	var f func() string
	switch shell {
	case "bash":
		f = a.bash
	case "fish":
		f = a.fish
	case "elvish":
		f = a.elvish
	case "oil":
		f = a.oil
	case "powershell":
		f = a.powershell
	case "xonsh":
		f = a.xonsh
	case "zsh":
		f = a.zsh
	}

	if f == nil {
		// TODO "{}" for xonsh?
		return ""
	} else {
		return f()
	}
}

// ActionCallback invokes a go function during completion
func ActionCallback(callback CompletionCallback) Action {
	return Action{callback: callback}
}

// ActionBool completes true/false
func ActionBool() Action {
	return ActionValues("true", "false")
}

// ActionDirectories completes directories
func ActionDirectories() Action {
	return ActionCallback(func(args []string) Action {
		return actionPath("", true).Invoke(args).ToMultiPartsA("/")
	})
}

// ActionFiles completes files with optional suffix filtering
func ActionFiles(suffix string) Action {
	return ActionCallback(func(args []string) Action {
		return actionPath(suffix, false).Invoke(args).ToMultiPartsA("/")
	})
}

func actionPath(fileSuffix string, dirOnly bool) Action {
	folder := filepath.Dir(CallbackValue)
	if files, err := ioutil.ReadDir(folder); err != nil {
		return ActionMessage(err.Error())
	} else {
		if folder == "." {
			folder = ""
		} else if !strings.HasSuffix(folder, "/") {
			folder = folder + "/"
		}

		vals := make([]string, len(files))
		for index, file := range files {
			if file.IsDir() {
				vals[index] = folder + file.Name() + "/"
			} else if !dirOnly && strings.HasSuffix(file.Name(), fileSuffix) {
				vals[index] = folder + file.Name()
			}
		}
		return ActionValues(vals...)
	}
}

// ActionValues completes arbitrary keywords (values)
func ActionValues(values ...string) Action {
	vals := make([]string, len(values)*2)
	for index, val := range values {
		vals[index*2] = val
		vals[(index*2)+1] = ""
	}
	return ActionValuesDescribed(vals...)
}

// ActionValuesDescribed completes arbitrary key (values) with an additional description (value, description pairs)
func ActionValuesDescribed(values ...string) Action {
	vals := make([]common.RawValue, len(values)/2)
	for index, val := range values {
		if index%2 == 0 {
			vals[index/2] = common.RawValue{Value: val, Display: val, Description: values[index+1]}
		}
	}
	return actionRawValues(vals...)
}

func actionRawValues(rawValues ...common.RawValue) Action {
	return Action{
		rawValues:  rawValues,
		bash:       func() string { return bash.ActionRawValues(rawValues...) },
		elvish:     func() string { return elvish.ActionRawValues(rawValues...) },
		fish:       func() string { return fish.ActionRawValues(rawValues...) },
		oil:        func() string { return oil.ActionRawValues(rawValues...) },
		powershell: func() string { return powershell.ActionRawValues(rawValues...) },
		xonsh:      func() string { return xonsh.ActionRawValues(rawValues...) },
		zsh:        func() string { return zsh.ActionRawValues(rawValues...) },
	}
}

var skipCache bool

// ActionMessage displays a help messages in places where no completions can be generated
func ActionMessage(msg string) Action {
	skipCache = true // TODO find a better solution - any call to ActionMessage i assumed to be an error for now
	return ActionValuesDescribed("_", "", "ERR", msg)
}

// CallbackValue is set to the currently completed flag/positional value during callback (note that this is updated during ActionMultiParts)
var CallbackValue string

// ActionMultiParts completes multiple parts of words separately where each part is separated by some char (CallbackValue is set to the currently completed part during invocation)
func ActionMultiParts(divider string, callback func(args []string, parts []string) Action) Action {
	return ActionCallback(func(args []string) Action {
		oldValue := CallbackValue
		defer func() { CallbackValue = oldValue }() // TODO verify this

		index := strings.LastIndex(CallbackValue, string(divider))
		prefix := ""
		if len(divider) == 0 {
			prefix = CallbackValue
		} else if index != -1 {
			prefix = CallbackValue[0 : index+len(divider)]
			CallbackValue = CallbackValue[index+len(divider):] // update CallbackValue to only contain the currently completed part
		}
		parts := strings.Split(prefix, string(divider))
		if len(parts) > 0 {
			parts = parts[0 : len(parts)-1]
		}

		return callback(args, parts).Invoke(args).Prefix(prefix).ToA()
	})
}

func actionSubcommands(cmd *cobra.Command) Action {
	vals := make([]string, 0)
	for _, subcommand := range cmd.Commands() {
		if !subcommand.Hidden {
			vals = append(vals, subcommand.Name(), subcommand.Short)
			for _, alias := range subcommand.Aliases {
				vals = append(vals, alias, subcommand.Short)
			}
		}
	}
	return ActionValuesDescribed(vals...)
}

func actionFlags(cmd *cobra.Command) Action {
	vals := make([]string, 0)
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		if !common.IsShorthandOnly(f) {
			vals = append(vals, "--"+f.Name, f.Usage)
		}
		if f.Shorthand != "" {
			vals = append(vals, "-"+f.Shorthand, f.Usage)
		}
	})

	return ActionValuesDescribed(vals...)
}
