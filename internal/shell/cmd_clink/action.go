package cmd_clink

import (
	"github.com/carapace-sh/carapace/internal/common"
	"github.com/carapace-sh/carapace/internal/shell/fish"
)

func ActionRawValues(currentWord string, meta common.Meta, values common.RawValues) string {
	return fish.ActionRawValues(currentWord, meta, values) // TODO
}
