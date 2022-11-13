package carapace

import (
	"regexp"
	"runtime"
	"time"

	"github.com/rsteube/carapace/internal/cache"
	"github.com/rsteube/carapace/internal/common"
	pkgcache "github.com/rsteube/carapace/pkg/cache"
)

// Action indicates how to complete a flag or positional argument.
type Action struct {
	rawValues []common.RawValue
	callback  CompletionCallback
	nospace   bool
	skipcache bool
	hint      string // Non-error message,generally printed by shell caller, not as comp.
	message   string // A message is an error raised by the caller of this action.
}

// ActionMap maps Actions to an identifier.
type ActionMap map[string]Action

// CompletionCallback is executed during completion of associated flag or positional argument.
type CompletionCallback func(c Context) Action

// Cache cashes values of a CompletionCallback for given duration and keys.
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

// Invoke executes the callback of an action if it exists (supports nesting).
func (a Action) Invoke(c Context) InvokedAction {
	if c.Args == nil {
		c.Args = []string{}
	}
	if c.Env == nil {
		c.Env = []string{}
	}
	if c.Parts == nil {
		c.Parts = []string{}
	}
	return InvokedAction{a.nestedAction(c, 10)}
}

func (a Action) nestedAction(c Context, maxDepth int) Action {
	if maxDepth < 0 {
		return ActionMessage("maximum recursion depth exceeded")
	}
	if a.rawValues == nil && a.callback != nil {
		return a.callback(c).nestedAction(c, maxDepth-1).noSpace(a.nospace).skipCache(a.skipcache).withHint(a.hint)
	}
	return a
}

// NoSpace disables space suffix.
func (a Action) NoSpace() Action {
	return a.noSpace(true)
}

// Style sets the style
//
//	ActionValues("yes").Style(style.Green)
//	ActionValues("no").Style(style.Red)
func (a Action) Style(style string) Action {
	return a.StyleF(func(s string) string {
		return style
	})
}

// Style sets the style using a reference
//
//	ActionValues("value").StyleR(&style.Carapace.Value)
//	ActionValues("description").StyleR(&style.Carapace.Value)
func (a Action) StyleR(style *string) Action {
	return ActionCallback(func(c Context) Action {
		if style != nil {
			return a.Style(*style)
		}
		return a
	})
}

// Style sets the style using a function
//
//	ActionValues("dir/", "test.txt").StyleF(style.ForPathExt)
//	ActionValues("true", "false").StyleF(style.ForKeyword)
func (a Action) StyleF(f func(s string) string) Action {
	return ActionCallback(func(c Context) Action {
		invoked := a.Invoke(c)
		for index, v := range invoked.rawValues {
			if v.Value != "ERR" && v.Value != "_" {
				invoked.rawValues[index].Style = f(v.Value)
			}
		}

		return invoked.ToA()
	})
}

// Group gathers sets the group under which to print completions,
// for shells supporting this feature, like ZSH.
//
// ActionValue("192.168.1.1", "127.0.0.1").Group("IPv4 addresses").
func (a Action) Group(group string) Action {
	return a.GroupF(func(value string) string {
		return group
	})
}

// Group gathers sets the group under which to print completions,
// for shells supporting this feature, like ZSH, with a function.
func (a Action) GroupF(f func(value string) (group string)) Action {
	return ActionCallback(func(c Context) Action {
		invoked := a.Invoke(c)
		for index, v := range invoked.rawValues {
			if v.Value != "ERR" && v.Value != "_" {
				invoked.rawValues[index].Group = f(v.Value)
			}
		}

		return invoked.ToA()
	})
}

// Tag marks completions with a tag (which is different from the group).
// This function only has an effect for ZSH, which makes heavy use of tags.
// In most cases, this function is not needed, and for simple gathering of
// some completions under a group description, action.Group() is preferred.
//
// ActionValue("192.168.1.1", "127.0.0.1").Tag("interfaces").
func (a Action) Tag(tag string) Action {
	return a.TagF(func(value string) string {
		return tag
	})
}

// Tag marks completions with a tag (which is different from the group),
// using a function.
//
// This function only has an effect for ZSH, which makes heavy use of tags.
// In most cases, this function is not needed, and for simple gathering of
// some completions under a group description, action.Group() is preferred.
func (a Action) TagF(f func(value string) (tag string)) Action {
	return ActionCallback(func(c Context) Action {
		invoked := a.Invoke(c)
		for index, v := range invoked.rawValues {
			if v.Value != "ERR" && v.Value != "_" {
				invoked.rawValues[index].Tag = f(v.Value)
			}
		}

		return invoked.ToA()
	})
}

// Suffix adds a suffix to all raw values contained in the action.
// This suffix can automatically be removed upon entering a space
// character, or some other special characters in some context.
func (a Action) Suffix(suffix string, removable bool) Action {
	return ActionCallback(func(c Context) Action {
		invoked := a.Invoke(c)
		for index := range invoked.rawValues {
			invoked.rawValues[index].SuffixRemovable = suffix
		}

		return invoked.ToA()
	})
}

// SuffixValues adds a a suffix to all values currently contained by
// the Action, and which are found in the values array passed as parameter.
func (a Action) SuffixValues(values []string, suffix string) Action {
	return ActionCallback(func(c Context) Action {
		invoked := a.Invoke(c)
		for index, val := range invoked.rawValues {
			for _, value := range values {
				if val.Value == value {
					invoked.rawValues[index].SuffixRemovable = suffix
				}
			}
		}

		return invoked.ToA()
	})
}

// Chdir changes the current working directory to the named directory during invocation.
func (a Action) Chdir(dir string) Action {
	return ActionCallback(func(c Context) Action {
		abs, err := c.Abs(dir)
		if err != nil {
			return ActionMessage(err.Error())
		}
		c.Dir = abs
		return a.Invoke(c).ToA()
	})
}

// Suppress suppresses specific error messages using regular expressions.
func (a Action) Suppress(expr ...string) Action {
	return ActionCallback(func(c Context) Action {
		invoked := a.Invoke(c)

		// First filter the message field we use for storing errors.
		filter := false
		if a.message != "" {
			for _, e := range expr {
				r, err := regexp.Compile(e)
				if err != nil {
					return ActionMessage(err.Error())
				}
				if r.MatchString(a.message) {
					a.message = err.Error()
					filter = true
					break
				}
			}
		} else {
			for _, rawValue := range invoked.rawValues {
				if rawValue.Display == "ERR" {
					for _, e := range expr {
						r, err := regexp.Compile(e)
						if err != nil {
							return ActionMessage(err.Error())
						}
						if r.MatchString(rawValue.Description) {
							filter = true
							break
						}
					}
				}
			}
		}

		if filter {
			filtered := make([]common.RawValue, 0)
			for _, r := range invoked.rawValues {
				if r.Display != "ERR" && r.Display != "_" {
					filtered = append(filtered, r)
				}
			}
			invoked.rawValues = filtered
		}
		return invoked.ToA()
	})
}

func (a Action) withHint(s string) Action {
	if s != "" {
		a.hint = s
	}
	return a
}

func (a Action) noSpace(state bool) Action {
	a.nospace = a.nospace || state
	return a
}

func (a Action) skipCache(state bool) Action {
	a.skipcache = a.skipcache || state
	return a
}
