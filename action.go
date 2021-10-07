package carapace

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/rsteube/carapace/internal/cache"
	"github.com/rsteube/carapace/internal/common"
	pkgcache "github.com/rsteube/carapace/pkg/cache"
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

func (a Action) nestedAction(c Context, maxDepth int) Action {
	if a.rawValues == nil && a.callback != nil && maxDepth > 0 {
		return a.callback(c).nestedAction(c, maxDepth-1).noSpace(a.nospace).skipCache(a.skipcache)
	}
	return a
}

// NoSpace disables space suffix
func (a Action) NoSpace() Action {
	return a.noSpace(true)
}

// Chdir changes the current working directory to the named directory during invocation.
func (a Action) Chdir(dir string) Action {
	return ActionCallback(func(c Context) Action {
		if dir == "" || dir == "." {
			return a // do nothing on current dir
		}

		if strings.HasPrefix(dir, "~") {
			home, err := os.UserHomeDir()
			if err != nil {
				return ActionMessage(err.Error())
			}
			dir = strings.Replace(dir, "~", home, 1)
		}

		file, err := os.Stat(dir)
		if err != nil {
			return ActionMessage(err.Error())
		}
		if !file.IsDir() {
			return ActionMessage(fmt.Sprintf("%v is not a directory", dir))
		}

		current, err := os.Getwd()
		if err != nil {
			return ActionMessage(err.Error())
		}

		if err := os.Chdir(dir); err != nil {
			return ActionMessage(err.Error())
		}

		a := a.Invoke(c).ToA()

		if err := os.Chdir(current); err != nil {
			return ActionMessage(err.Error())
		}
		return a
	})
}

func (a Action) noSpace(state bool) Action {
	a.nospace = a.nospace || state
	return a
}

func (a Action) skipCache(state bool) Action {
	a.skipcache = a.skipcache || state
	return a
}
