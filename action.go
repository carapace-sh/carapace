package carapace

import (
	"bytes"
	exec "golang.org/x/sys/execabs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/rsteube/carapace/internal/bash"
	"github.com/rsteube/carapace/internal/cache"
	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/elvish"
	"github.com/rsteube/carapace/internal/fish"
	"github.com/rsteube/carapace/internal/ion"
	"github.com/rsteube/carapace/internal/nushell"
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
	rawValues []common.RawValue
	callback  CompletionCallback
	nospace   bool
	skipcache bool
}

// ActionMap maps Actions to an identifier
type ActionMap map[string]Action

// Context provides information during completion
type Context struct {
	// CallbackValue contains the (partial) value (or part of it during an ActionMultiParts) currently being completed
	CallbackValue string
	// Args contains the positional arguments of current (sub)command (exclusive the one currently being completed)
	Args []string
	// Parts contains the splitted CallbackValue during an ActionMultiParts (exclusive the part currently being completed)
	Parts []string
}

// CompletionCallback is executed during completion of associated flag or positional argument
type CompletionCallback func(c Context) Action

// Cache cashes values of a CompletionCallback for given duration and keys
func (a Action) Cache(timeout time.Duration, keys ...pkgcache.Key) Action {
	// TODO static actions are using callback now as well (for performance) - probably best to add a `static` bool to Action for this and check that here
	if a.callback != nil { // only relevant for callback actions
		cachedCallback := a.callback
		_, file, line, _ := runtime.Caller(1) // generate uid from wherever Cache() was called
		a.callback = func(c Context) Action {
			if cacheFile, err := cache.File(file, line, keys...); err == nil {
				if rawValues, err := cache.Load(cacheFile, timeout); err == nil {
					return actionRawValues(rawValues...)
				}
				invokedAction := (Action{callback: cachedCallback}).Invoke(c)
				if !invokedAction.skipcache {
					_ = cache.Write(cacheFile, invokedAction.rawValues)
				}
				return invokedAction.ToA()
			}
			return cachedCallback(c)
		}
	}
	return a
}

// InvokedAction is a logical alias for an Action whose (nested) callback was invoked
type InvokedAction Action

// Invoke executes the callback of an action if it exists (supports nesting)
func (a Action) Invoke(c Context) InvokedAction {
	if c.Args == nil {
		c.Args = []string{}
	}
	if c.Parts == nil {
		c.Parts = []string{}
	}
	return InvokedAction(a.nestedAction(c, 10))
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

func (a Action) noSpace(state bool) Action {
	a.nospace = a.nospace || state
	return a
}

func (a Action) skipCache(state bool) Action {
	a.skipcache = a.skipcache || state
	return a
}

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

func (a Action) nestedAction(c Context, maxDepth int) Action {
	if a.rawValues == nil && a.callback != nil && maxDepth > 0 {
		return a.callback(c).nestedAction(c, maxDepth-1).noSpace(a.nospace).skipCache(a.skipcache)
	}
	return a
}

func (a InvokedAction) value(shell string, callbackValue string) string { // TODO use context instead?
	switch shell {
	case "bash":
		return bash.ActionRawValues(callbackValue, a.nospace, a.rawValues...)
	case "fish":
		return fish.ActionRawValues(callbackValue, a.nospace, a.rawValues...)
	case "elvish":
		return elvish.ActionRawValues(callbackValue, a.nospace, a.rawValues...)
	case "ion":
		return ion.ActionRawValues(callbackValue, a.nospace, a.rawValues)
	case "nushell":
		return nushell.ActionRawValues(callbackValue, a.nospace, a.rawValues)
	case "oil":
		return oil.ActionRawValues(callbackValue, a.nospace, a.rawValues...)
	case "powershell":
		return powershell.ActionRawValues(callbackValue, a.nospace, a.rawValues...)
	case "xonsh":
		return xonsh.ActionRawValues(callbackValue, a.nospace, a.rawValues...)
	case "zsh":
		return zsh.ActionRawValues(callbackValue, a.nospace, a.rawValues...)
	default:
		return ""
	}
}

// ActionCallback invokes a go function during completion
func ActionCallback(callback CompletionCallback) Action {
	return Action{callback: callback}
}

// ActionExecCommand invokes given command and transforms its output using given function on success or returns ActionMessage with the first line of stderr if available.
//   carapace.ActionExecCommand("git", "remote")(func(output []byte) carapace.Action {
//     lines := strings.Split(string(output), "\n")
//     return carapace.ActionValues(lines[:len(lines)-1]...)
//   })
func ActionExecCommand(name string, arg ...string) func(f func(output []byte) Action) Action {
	return func(f func(output []byte) Action) Action {
		return ActionCallback(func(c Context) Action {
			var stdout, stderr bytes.Buffer
			cmd := exec.Command(name, arg...)
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				if firstLine := strings.SplitN(stderr.String(), "\n", 2)[0]; strings.TrimSpace(firstLine) != "" {
					return ActionMessage(firstLine)
				}
				return ActionMessage(err.Error())
			}
			return f(stdout.Bytes())
		})
	}
}

// ActionDirectories completes directories
func ActionDirectories() Action {
	return ActionCallback(func(c Context) Action {
		return actionPath([]string{""}, true).Invoke(c).ToMultiPartsA("/").noSpace(true)
	})
}

// ActionFiles completes files with optional suffix filtering
func ActionFiles(suffix ...string) Action {
	return ActionCallback(func(c Context) Action {
		return actionPath(suffix, false).Invoke(c).ToMultiPartsA("/").noSpace(true)
	})
}

func actionPath(fileSuffixes []string, dirOnly bool) Action {
	return ActionCallback(func(c Context) Action {
		folder := filepath.Dir(c.CallbackValue)
		expandedFolder := folder
		if strings.HasPrefix(c.CallbackValue, "~") {
			homedir, err := os.UserHomeDir()
			if err != nil {
				return ActionMessage(err.Error())
			}
			expandedFolder = filepath.Dir(homedir + "/" + c.CallbackValue[1:])
		}

		files, err := ioutil.ReadDir(expandedFolder)
		if err != nil {
			return ActionMessage(err.Error())
		}
		if folder == "." {
			folder = ""
		} else if !strings.HasSuffix(folder, "/") {
			folder = folder + "/"
		}

		showHidden := c.CallbackValue != "" &&
			!strings.HasSuffix(c.CallbackValue, "/") &&
			strings.HasPrefix(filepath.Base(c.CallbackValue), ".")

		vals := make([]string, 0, len(files))
		for _, file := range files {
			if !showHidden && strings.HasPrefix(file.Name(), ".") {
				continue
			}

			if file.IsDir() {
				vals = append(vals, folder+file.Name()+"/")
			} else if !dirOnly {
				if len(fileSuffixes) == 0 {
					fileSuffixes = []string{""}
				}
				for _, suffix := range fileSuffixes {
					if strings.HasSuffix(file.Name(), suffix) {
						vals = append(vals, folder+file.Name())
						break
					}
				}
			}
		}
		if strings.HasPrefix(c.CallbackValue, "./") {
			return ActionValues(vals...).Invoke(Context{}).Prefix("./").ToA()
		}
		return ActionValues(vals...)
	})
}

// ActionValues completes arbitrary keywords (values)
func ActionValues(values ...string) Action {
	return ActionCallback(func(c Context) Action {
		vals := make([]string, len(values)*2)
		for index, val := range values {
			vals[index*2] = val
			vals[(index*2)+1] = ""
		}
		return ActionValuesDescribed(vals...)
	})
}

// ActionValuesDescribed completes arbitrary key (values) with an additional description (value, description pairs)
func ActionValuesDescribed(values ...string) Action {
	return ActionCallback(func(c Context) Action {
		vals := make([]common.RawValue, len(values)/2)
		for index, val := range values {
			if index%2 == 0 {
				vals[index/2] = common.RawValue{Value: val, Display: val, Description: values[index+1]}
			}
		}
		return actionRawValues(vals...)
	})
}

func actionRawValues(rawValues ...common.RawValue) Action {
	return Action{
		rawValues: rawValues,
	}
}

// ActionMessage displays a help messages in places where no completions can be generated
func ActionMessage(msg string) Action {
	return ActionCallback(func(c Context) Action {
		return ActionValuesDescribed("_", "", "ERR", msg).
			Invoke(c).Prefix(c.CallbackValue).ToA(). // needs to be prefixed with current callback value to not be filtered out
			noSpace(true).skipCache(true)
	})
}

// ActionMultiParts completes multiple parts of words separately where each part is separated by some char (CallbackValue is set to the currently completed part during invocation)
func ActionMultiParts(divider string, callback func(c Context) Action) Action {
	return ActionCallback(func(c Context) Action {
		index := strings.LastIndex(c.CallbackValue, string(divider))
		prefix := ""
		if len(divider) == 0 {
			prefix = c.CallbackValue
			c.CallbackValue = ""
		} else if index != -1 {
			prefix = c.CallbackValue[0 : index+len(divider)]
			c.CallbackValue = c.CallbackValue[index+len(divider):] // update CallbackValue to only contain the currently completed part
		}
		parts := strings.Split(prefix, string(divider))
		if len(parts) > 0 && len(divider) > 0 {
			parts = parts[0 : len(parts)-1]
		}
		c.Parts = parts

		return callback(c).Invoke(c).Prefix(prefix).ToA().noSpace(true)
	})
}

func actionSubcommands(cmd *cobra.Command) Action {
	vals := make([]string, 0)
	for _, subcommand := range cmd.Commands() {
		if !subcommand.Hidden && subcommand.Deprecated == "" {
			vals = append(vals, subcommand.Name(), subcommand.Short)
			for _, alias := range subcommand.Aliases {
				vals = append(vals, alias, subcommand.Short)
			}
		}
	}
	return ActionValuesDescribed(vals...)
}

func actionFlags(cmd *cobra.Command) Action {
	return ActionCallback(func(c Context) Action {
		re := regexp.MustCompile("^-(?P<shorthand>[^-=]+)")
		isShorthandSeries := re.MatchString(c.CallbackValue)

		vals := make([]string, 0)
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if f.Deprecated != "" {
				return // skip deprecated flags
			}

			if f.Changed &&
				!strings.Contains(f.Value.Type(), "Slice") &&
				!strings.Contains(f.Value.Type(), "Array") {
				return // don't repeat flag
			}

			if isShorthandSeries {
				if f.Shorthand != "" && f.ShorthandDeprecated == "" {
					for _, shorthand := range c.CallbackValue[1:] {
						if shorthandFlag := cmd.Flags().ShorthandLookup(string(shorthand)); shorthandFlag != nil && shorthandFlag.Value.Type() != "bool" && shorthandFlag.NoOptDefVal == "" {
							return // abort shorthand flag series if a previous one is not bool and requires an argument (no default value)
						}
					}
					vals = append(vals, f.Shorthand, f.Usage)
				}
			} else {
				if !common.IsShorthandOnly(f) {
					vals = append(vals, "--"+f.Name, f.Usage)
				}
				if f.Shorthand != "" && f.ShorthandDeprecated == "" {
					vals = append(vals, "-"+f.Shorthand, f.Usage)
				}
			}
		})

		if isShorthandSeries {
			matches := re.FindStringSubmatch(c.CallbackValue)
			parts := strings.Split(matches[1], "")
			return ActionValuesDescribed(vals...).Invoke(c).Filter(parts).Prefix(c.CallbackValue).ToA().noSpace(true)
		}
		return ActionValuesDescribed(vals...)
	})
}
