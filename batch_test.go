package carapace

import (
	"testing"

	"github.com/rsteube/carapace/internal/common"
)

func TestBatch(t *testing.T) {
	b := Batch(
		ActionValues("A", "B"),
		ActionValues("B", "C"),
		ActionValues("C", "D"),
	)
	expected := InvokedAction{
		Action{
			rawValues: common.RawValuesFrom("A", "B", "C", "D"),
		},
	}
	actual := b.Invoke(Context{}).Merge()
	assertEqual(t, expected, actual)
}

func TestBatchSingle(t *testing.T) {
	b := Batch(
		ActionValues("A", "B"),
	)
	expected := InvokedAction{
		Action{
			rawValues: common.RawValuesFrom("A", "B"),
		},
	}
	actual := b.Invoke(Context{}).Merge()
	assertEqual(t, expected, actual)
}

func TestBatchNone(t *testing.T) {
	b := Batch()
	expected := InvokedAction{
		Action{
			rawValues: common.RawValuesFrom(),
		},
	}
	actual := b.Invoke(Context{}).Merge()
	assertEqual(t, expected, actual)
}

func TestBatchToA(t *testing.T) {
	b := Batch(
		ActionValues("A", "B"),
		ActionValues("B", "C"),
		ActionValues("C", "D"),
	)
	expected := InvokedAction{
		Action{
			rawValues: common.RawValuesFrom("A", "B", "C", "D"),
		},
	}
	actual := b.ToA().Invoke(Context{})
	assertEqual(t, expected, actual)
}
