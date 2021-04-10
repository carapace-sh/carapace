package carapace

import (
	"fmt"
	"testing"

	"github.com/rsteube/carapace/internal/assert"
	"github.com/rsteube/carapace/internal/common"
)

func assertEqual(t *testing.T, expected, actual InvokedAction) {
	assert.Equal(t, fmt.Sprintf("%+v\n", expected), fmt.Sprintf("%+v\n", actual))
}

func TestActionCallback(t *testing.T) {
	a := ActionCallback(func(c Context) Action {
		return ActionCallback(func(c Context) Action {
			return ActionCallback(func(c Context) Action {
				return ActionValues("a", "b", "c")
			})
		})
	})
	expected := InvokedAction{
		rawValues: common.RawValuesFrom("a", "b", "c"),
		nospace:   false,
		skipcache: false,
	}
	actual := a.Invoke(Context{})
	assertEqual(t, expected, actual)
}

func TestSkipCache(t *testing.T) {
	a := ActionCallback(func(c Context) Action {
		return ActionValues().Invoke(c).Merge(
			ActionCallback(func(c Context) Action {
				return ActionMessage("skipcache")
			}).Invoke(c)).
			Filter([]string{""}).
			Prefix("").
			Suffix("").
			ToA()
	})
	if a.skipcache {
		t.Fatal("uninvoked skipcache should be false")
	}
	if !a.Invoke(Context{}).skipcache {
		t.Fatal("invoked skipcache should be true")
	}
}

func TestNoSpace(t *testing.T) {
	a := ActionCallback(func(c Context) Action {
		return ActionValues().Invoke(c).Merge(
			ActionMultiParts("", func(c Context) Action {
				return ActionMessage("nospace")
			}).Invoke(c)).
			Filter([]string{""}).
			Prefix("").
			Suffix("").
			ToA()
	})
	if a.nospace {
		t.Fatal("uninvoked nospace should be false")
	}
	if !a.Invoke(Context{}).nospace {
		t.Fatal("invoked nospace should be true")
	}
}
