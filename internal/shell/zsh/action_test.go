package zsh

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/carapace-sh/carapace/internal/common"
)

// block mirrors one `tag\003suffix\003displays\003values` entry of the payload.
type block struct {
	tag     string
	suffix  string
	values  []string
	display []string
}

// parsePayload decodes the delimited string produced by ActionRawValues into
// its blocks so tests can assert on the per-group suffix handling.
func parsePayload(t *testing.T, payload string) []block {
	t.Helper()
	fields := strings.SplitN(payload, "\001", 4)
	if len(fields) != 4 {
		t.Fatalf("expected 4 top-level fields, got %d in %q", len(fields), payload)
	}
	data := strings.TrimSuffix(fields[3], "\001")
	blocks := make([]block, 0)
	for raw := range strings.SplitSeq(data, "\002") {
		if raw == "" {
			continue
		}
		parts := strings.SplitN(raw, "\003", 4)
		if len(parts) != 4 {
			t.Fatalf("expected 4 block parts, got %d in %q", len(parts), raw)
		}
		blocks = append(blocks, block{
			tag:     parts[0],
			suffix:  parts[1],
			display: strings.Split(parts[2], "\n"),
			values:  strings.Split(parts[3], "\n"),
		})
	}
	return blocks
}

func nospaceMeta(suffixes ...rune) common.Meta {
	var meta common.Meta
	meta.Nospace.Add(suffixes...)
	return meta
}

// TestSuffixGroupingDefaultState verifies that in an unquoted context a nospace
// value ending in a type-separator is placed in its own suffix group with the
// separator stripped from the value.
func TestSuffixGroupingDefaultState(t *testing.T) {
	t.Setenv("CARAPACE_COMPLINE", "example ")

	values := common.RawValues{
		{Value: "src/", Display: "src/", Tag: "dirs"},
		{Value: "pkg/", Display: "pkg/", Tag: "dirs"},
		{Value: "main.go", Display: "main.go", Tag: "dirs"}, // no separator -> empty suffix
	}
	blocks := parsePayload(t, ActionRawValues("", nospaceMeta('/'), values))

	bySuffix := map[string]block{}
	for _, b := range blocks {
		bySuffix[b.suffix] = b
	}

	slash, ok := bySuffix["/"]
	if !ok {
		t.Fatalf("expected a '/' suffix group, got %+v", blocks)
	}
	for _, v := range slash.values {
		if strings.HasSuffix(v, "/") {
			t.Errorf("separator not stripped from value %q", v)
		}
	}
	if got := strings.Join(slash.values, ","); got != "pkg,src" && got != "src,pkg" {
		t.Errorf("unexpected '/' group values: %q", got)
	}

	// main.go has no separator: it must land in the empty-suffix group.
	empty, ok := bySuffix[""]
	if !ok {
		t.Fatalf("expected an empty-suffix group for main.go, got %+v", blocks)
	}
	// main.go is not nospace-matched, so it keeps the normal trailing space.
	if !slices.Contains(empty.values, "main.go ") {
		t.Errorf("main.go (with trailing space) missing from empty-suffix group: %+v", empty.values)
	}
}

// TestSuffixGroupingQuotedState verifies scrutiny point #2: in a quoted context
// no separator suffix group is emitted (the value would otherwise gain a
// closing quote and the separator would land outside it).
func TestSuffixGroupingQuotedState(t *testing.T) {
	t.Setenv("CARAPACE_COMPLINE", "example 'src")

	values := common.RawValues{
		{Value: "src/", Display: "src/", Tag: "dirs"},
	}
	blocks := parsePayload(t, ActionRawValues("", nospaceMeta('/'), values))
	for _, b := range blocks {
		if b.suffix == "/" {
			t.Errorf("quoted state must not produce a '/' suffix group: %+v", b)
		}
		for _, v := range b.values {
			// value keeps its slash (not stripped) and gains a closing quote.
			if !strings.HasSuffix(v, "'") {
				t.Errorf("quoted value missing closing quote: %q", v)
			}
		}
	}
}

// TestSuffixGroupingNarrowing verifies scrutiny point #1: a nospace value ending
// in punctuation that is NOT a type-separator (e.g. '!') is not treated as a
// removable suffix.
func TestSuffixGroupingNarrowing(t *testing.T) {
	t.Setenv("CARAPACE_COMPLINE", "example ")

	values := common.RawValues{
		{Value: "foo!", Display: "foo!", Tag: "vals"},
	}
	blocks := parsePayload(t, ActionRawValues("", nospaceMeta('!'), values))
	for _, b := range blocks {
		if b.suffix == "!" {
			t.Errorf("'!' must not be treated as a removable suffix: %+v", b)
		}
	}
}

// BenchmarkActionRawValues measures rendering cost on a large mixed set of
// files, directories and optarg flags (the worst case for suffix grouping).
func BenchmarkActionRawValues(b *testing.B) {
	b.Setenv("CARAPACE_COMPLINE", "example ")

	const n = 5000
	values := make(common.RawValues, 0, n)
	for i := range n {
		switch i % 3 {
		case 0:
			values = append(values, common.RawValue{Value: fmt.Sprintf("file%d.go", i), Display: fmt.Sprintf("file%d.go", i), Tag: "files"})
		case 1:
			values = append(values, common.RawValue{Value: fmt.Sprintf("dir%d/", i), Display: fmt.Sprintf("dir%d/", i), Tag: "files"})
		case 2:
			values = append(values, common.RawValue{Value: fmt.Sprintf("--flag%d=", i), Display: fmt.Sprintf("--flag%d=", i), Tag: "flags"})
		}
	}
	meta := nospaceMeta('/', '=')

	b.ReportAllocs()
	for b.Loop() {
		_ = ActionRawValues("", meta, values)
	}
}
