package carapace

import (
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/env"
	"github.com/rsteube/carapace/internal/pflagfork"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type inArgsX []inArgX
type inArgX struct {
	Value       string  `yaml:",omitempty"`
	Name        string  `yaml:",omitempty"`
	Description string  `yaml:",omitempty"`
	Type        string  `yaml:",omitempty"`
	OptArg      *inArgX `yaml:",omitempty"`
}

func traverse(c *cobra.Command, args []string) (Action, Context) {
	LOG.Printf("traverse called for %#v with args %#v\n", c.Name(), args)
	storage.preRun(c, args)

	if env.Lenient() {
		LOG.Printf("allowing unknown flags")
		c.FParseErrWhitelist.UnknownFlags = true
	}

	inArgsX := make(inArgsX, 0)
	inArgs := []string{}        // args consumed by current command
	inPositionals := []string{} // positionals consumed by current command
	var inFlag *pflagfork.Flag  // last encountered flag that still expects arguments
	c.LocalFlags()              // TODO force  c.mergePersistentFlags() which is missing from c.Flags()
	fs := pflagfork.FlagSet{FlagSet: c.Flags()}

	context := NewContext(args...)
loop:
	for i, arg := range context.Args {
		switch {
		// flag argument
		case inFlag != nil && inFlag.Consumes(arg):
			LOG.Printf("arg %#v is a flag argument\n", arg)
			inArgsX = append(inArgsX, inArgX{
				Value: arg,
				Type:  "flag argument",
			})
			inArgs = append(inArgs, arg)
			inFlag.Args = append(inFlag.Args, arg)

			if !inFlag.Consumes("") {
				inFlag = nil // no more args expected
			}
			continue

		// dash
		case arg == "--":
			LOG.Printf("arg %#v is dash\n", arg)
			inArgsX = append(inArgsX, inArgX{
				Type:  "dash",
				Value: arg,
			})
			inArgs = append(inArgs, context.Args[i:]...)
			break loop

		// flag
		case !c.DisableFlagParsing && strings.HasPrefix(arg, "-") && (fs.IsInterspersed() || len(inPositionals) == 0):
			LOG.Printf("arg %#v is a flag\n", arg)
			inArgs = append(inArgs, arg)
			inFlag = fs.LookupArg(arg)

			if inFlag == nil {
				LOG.Printf("flag %#v is unknown", arg)
				inArgsX = append(inArgsX, inArgX{
					Type:  "unknown flag",
					Value: arg,
				})
			} else {
				inArgsX = append(inArgsX, inArgX{
					Name:        inFlag.Name,
					Description: inFlag.Usage,
					Type:        "flag", // TODO flagtype
					Value:       arg,
				})
			}
			continue

		// subcommand
		case subcommand(c, arg) != nil:
			LOG.Printf("arg %#v is a subcommand\n", arg)

			switch {
			case c.DisableFlagParsing:
				LOG.Printf("flag parsing disabled for %#v\n", c.Name())

			default:
				LOG.Printf("parsing flags for %#v with args %#v\n", c.Name(), inArgs)
				if err := c.ParseFlags(inArgs); err != nil {
					return ActionMessage(err.Error()), context
				}
				context.Args = c.Flags().Args()
			}

			subCmd := subcommand(c, arg)
			inArgsX = append(inArgsX, inArgX{
				Name:        subCmd.Name(),
				Description: subCmd.Short,
				Type:        "subcommand",
				Value:       arg,
			})
			m, _ := yaml.Marshal(inArgsX) // TODO pass on to next traverse invocation
			LOG.Println(string(m))

			return traverse(subCmd, args[i+1:])

		// positional
		default:
			LOG.Printf("arg %#v is a positional\n", arg)
			inArgsX = append(inArgsX, inArgX{
				Type:  "positional argument",
				Value: arg,
			})
			inArgs = append(inArgs, arg)
			inPositionals = append(inPositionals, arg)
		}
	}

	m, _ := yaml.Marshal(inArgsX)
	LOG.Println(string(m))

	toParse := inArgs
	if inFlag != nil && len(inFlag.Args) == 0 && inFlag.Consumes("") {
		LOG.Printf("removing arg %#v since it is a flag missing its argument\n", toParse[len(toParse)-1])
		toParse = toParse[:len(toParse)-1]
	} else if (fs.IsInterspersed() || len(inPositionals) == 0) && fs.IsShorthandSeries(context.Value) { // TODO shorthand series isn't correct anymore (can have value attached)
		LOG.Printf("arg %#v is a shorthand flag series", context.Value) // TODO not aways correct
		localInFlag := fs.LookupArg(context.Value)

		if localInFlag != nil && (len(localInFlag.Args) == 0 || localInFlag.Args[0] == "") && (!localInFlag.IsOptarg() || strings.HasSuffix(localInFlag.Prefix, string(localInFlag.OptargDelimiter()))) { // TODO && len(context.Value) > 2 {
			// TODO check if empty prefix
			suffix := localInFlag.Prefix[strings.LastIndex(localInFlag.Prefix, localInFlag.Shorthand):]
			LOG.Printf("removing suffix %#v since it is a flag missing its argument\n", suffix)
			toParse = append(toParse, strings.TrimSuffix(localInFlag.Prefix, suffix))
		} else {
			LOG.Printf("adding shorthand flag %#v", context.Value)
			toParse = append(toParse, context.Value)
		}

	}

	// TODO duplicated code
	switch {
	case c.DisableFlagParsing:
		LOG.Printf("flag parsing is disabled for %#v\n", c.Name())

	default:
		LOG.Printf("parsing flags for %#v with args %#v\n", c.Name(), toParse)
		if err := c.ParseFlags(toParse); err != nil {
			return ActionMessage(err.Error()), context
		}
		context.Args = c.Flags().Args()
	}

	switch {
	// dash argument
	case common.IsDash(c):
		LOG.Printf("completing dash for arg %#v\n", context.Value)
		context.Args = c.Flags().Args()[c.ArgsLenAtDash():]
		LOG.Printf("context: %#v\n", context.Args)

		return storage.getPositional(c, len(context.Args)), context

	// flag argument
	case inFlag != nil && inFlag.Consumes(context.Value):
		LOG.Printf("completing flag argument of %#v for arg %#v\n", inFlag.Name, context.Value)
		context.Parts = inFlag.Args
		return storage.getFlag(c, inFlag.Name), context

	// flag
	case !c.DisableFlagParsing && strings.HasPrefix(context.Value, "-") && (fs.IsInterspersed() || len(inPositionals) == 0):
		if f := fs.LookupArg(context.Value); f != nil && len(f.Args) > 0 {
			LOG.Printf("completing optional flag argument for arg %#v with prefix %#v\n", context.Value, f.Prefix)

			switch f.Value.Type() {
			case "bool":
				return ActionValues("true", "false").StyleF(style.ForKeyword).Usage(f.Usage).Prefix(f.Prefix), context
			default:
				return storage.getFlag(c, f.Name).Prefix(f.Prefix), context
			}
		} else if f != nil && fs.IsPosix() && !strings.HasPrefix(context.Value, "--") && !f.IsOptarg() && f.Prefix == context.Value {
			LOG.Printf("completing attached flag argument for arg %#v with prefix %#v\n", context.Value, f.Prefix)
			return storage.getFlag(c, f.Name).Prefix(f.Prefix), context
		}
		LOG.Printf("completing flags for arg %#v\n", context.Value)
		return actionFlags(c), context

	// positional or subcommand
	default:
		LOG.Printf("completing positionals and subcommands for arg %#v\n", context.Value)
		batch := Batch(storage.getPositional(c, len(context.Args)))
		if c.HasAvailableSubCommands() && len(context.Args) == 0 {
			batch = append(batch, actionSubcommands(c))
		}
		return batch.ToA(), context
	}
}

func subcommand(cmd *cobra.Command, arg string) *cobra.Command {
	if subcommand, _, _ := cmd.Find([]string{arg}); subcommand != cmd {
		return subcommand
	}
	return nil
}
