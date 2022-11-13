package carapace

import (
	"regexp"
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/pflagfork"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func actionRawValues(rawValues ...common.RawValue) Action {
	return Action{
		rawValues: rawValues,
	}
}

func actionSubcommands(cmd *cobra.Command) Action {
	vals := make([]string, 0)
	groups := make(map[string]string, len(cmd.Commands()))

	// Order the completions in the order of groups, when they have one
	for _, group := range cmd.Groups() {
		for _, subcommand := range cmd.Commands() {
			if subcommand.Hidden || subcommand.Deprecated != "" {
				continue
			}
			if subcommand.GroupID != group.ID {
				continue
			}

			vals = append(vals, subcommand.Name(), subcommand.Short)
			for _, alias := range subcommand.Aliases {
				vals = append(vals, alias, subcommand.Short)
			}

			for _, use := range append(subcommand.Aliases, subcommand.Name()) {
				groups[use] = subcommand.GroupID
			}
		}
	}

	// Then add commands not belonging to a group.
	for _, subcommand := range cmd.Commands() {
		if subcommand.Hidden || subcommand.Deprecated != "" {
			continue
		}
		if subcommand.GroupID != "" {
			continue
		}

		vals = append(vals, subcommand.Name(), subcommand.Short)
		for _, alias := range subcommand.Aliases {
			vals = append(vals, alias, subcommand.Short)
		}
	}

	// And gather them under their group as description
	grouper := groupCommands(cmd, groups)

	return ActionValuesDescribed(vals...).GroupF(grouper).Tag(string(common.Command))
}

// Generates a function to tag command completions with their corresponding group.
func groupCommands(cmd *cobra.Command, groups map[string]string) func(string) string {
	groupCommands := func(command string) string {
		cmdGroupID := groups[command]
		for _, group := range cmd.Groups() {
			if group.ID == cmdGroupID {
				title := strings.TrimSpace(group.Title)
				if !strings.HasSuffix(title, "commands") {
					title += " commands"
				}
				return title
			}
		}

		if len(cmd.Groups()) > 0 {
			return "other commands"
		}

		return "commands"
	}

	return groupCommands
}

func actionFlags(cmd *cobra.Command) Action {
	return ActionCallback(func(ctx Context) Action {
		re := regexp.MustCompile("^-(?P<shorthand>[^-=]+)")
		isShorthandSeries := re.MatchString(ctx.CallbackValue) && pflagfork.IsPosix(cmd.Flags())

		vals := make([]string, 0)

		var noOptDefValues []string

		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			// Build all raw completion values
			yesL, long, yesS, short := buildflagValues(cmd, &ctx, flag, isShorthandSeries)
			if !yesL {
				return
			}

			// Add the completion values we've built
			if yesS {
				vals = append(vals, short, flag.Usage)
			}

			vals = append(vals, long, flag.Usage)

			// And modify completions with any potential spec.
			suffixedVals := addflagSpecSuffix(flag, long, short, yesS)
			noOptDefValues = append(noOptDefValues, suffixedVals...)
		})

		// First generate the values without prefixing them...
		flagAction := ActionValuesDescribed(vals...).Tag(string(common.Flag))

		// Because we might have modifier functions that
		// need to match their unprefixed/modified values.
		flagAction = flagAction.SuffixValues(noOptDefValues, "=")

		// Apply any stacked flags prefix.
		if isShorthandSeries {
			return flagAction.Invoke(ctx).Prefix(ctx.CallbackValue).ToA()
		}

		return flagAction
	})
}

// addflagSpecSuffix.
func addflagSpecSuffix(flag *pflag.Flag, long, short string, isshort bool) []string {
	var noOptDefValues []string

	// Having a specified noOptDefValue adds an equal sign in most cases
	if flag.NoOptDefVal != "" && flag.Value.Type() != "bool" {
		if isshort {
			noOptDefValues = append(noOptDefValues, short)
		}
		noOptDefValues = append(noOptDefValues, long)
	}

	return noOptDefValues
}

func buildflagValues(cmd *cobra.Command, c *Context, f *pflag.Flag, series bool) (bool, string, bool, string) {
	var long, short string

	yesL, yesS := false, false

	if f.Deprecated != "" {
		return yesL, long, yesS, short
	}

	if f.Changed &&
		!strings.Contains(f.Value.Type(), "Slice") &&
		!strings.Contains(f.Value.Type(), "Array") &&
		!strings.Contains(f.Value.Type(), "map") &&
		f.Value.Type() != "count" {
		return yesL, long, yesS, short
	}

	// Here we know we have a flag to be completed.
	// We always populate the long variable, even if it's with a shorhand
	yesL = true

	if series {
		if f.Shorthand != "" && f.ShorthandDeprecated == "" {
			for _, shorthand := range c.CallbackValue[1:] {
				// abort shorthand flag series if a previous one is not bool or count and requires an argument (no default value)
				sflag := cmd.Flags().ShorthandLookup(string(shorthand))
				if sflag != nil {
					if sflag.Value.Type() != "bool" &&
						sflag.Value.Type() != "count" &&
						sflag.NoOptDefVal == "" {

						yesL = false

						return yesL, long, yesS, short
					}
				}
			}

			// Else we have a valid stacked (shorthand) flag
			long = f.Shorthand
		}

		return yesL, long, yesS, short
	}

	if flagstyle := pflagfork.Style(f); flagstyle != pflagfork.ShorthandOnly {
		if flagstyle == pflagfork.NameAsShorthand {
			long = "-" + f.Name
		} else {
			long = "--" + f.Name
		}
	}
	if f.Shorthand != "" && f.ShorthandDeprecated == "" {
		short = "-" + f.Shorthand
		yesS = true
	}

	return yesL, long, yesS, short
}
