package carapace

import (
	"strings"

	"github.com/rsteube/carapace/internal/common"
	"github.com/rsteube/carapace/internal/config"
	"github.com/rsteube/carapace/internal/pflagfork"
	"github.com/spf13/cobra"
)

type InFlag struct {
	*pflagfork.Flag
	// currently consumed args since encountered flag
	Args []string
}

func (f InFlag) Consumes(arg string) bool {
	switch {
	case f.Flag == nil:
		return false
	case !f.TakesValue():
		return false
	case f.IsOptarg():
		return false
	case len(f.Args) == 0:
		return true
	case f.Nargs() > 1 && len(f.Args) < f.Nargs():
		return true
	case f.Nargs() < 0 && !strings.HasPrefix(arg, "-"):
		return true
	default:
		return false
	}
}

func traverse(c *cobra.Command, args []string) (Action, Context) {
	logger.Printf("traverse called for %#v with args %#v\n", c.Name(), args)
	storage.preRun(c, args)

	if config.IsLenient() {
		logger.Printf("allowing unknown flags")
		c.FParseErrWhitelist.UnknownFlags = true
	}

	inArgs := []string{} // args consumed by current command
	var inFlag *InFlag   // last encountered flag that still expects arguments
	c.LocalFlags()       // TODO force  c.mergePersistentFlags() which is missing from c.Flags()
	fs := pflagfork.FlagSet{FlagSet: c.Flags()}

	context := NewContext(args...)
loop:
	for i, arg := range context.Args {
		switch {
		// flag argument
		case inFlag != nil && inFlag.Consumes(arg):
			logger.Printf("arg %#v is a flag argument\n", arg)
			inArgs = append(inArgs, arg)
			inFlag.Args = append(inFlag.Args, arg)

			if !inFlag.Consumes("") {
				inFlag = nil // no more args expected
			}
			continue

		// dash
		case arg == "--":
			logger.Printf("arg %#v is dash\n", arg)
			inArgs = append(inArgs, context.Args[i:]...)
			break loop

		// flag
		case !c.DisableFlagParsing && strings.HasPrefix(arg, "-"):
			logger.Printf("arg %#v is a flag\n", arg)
			inArgs = append(inArgs, arg)
			inFlag = &InFlag{
				Flag: fs.LookupArg(arg),
				Args: []string{},
			}

			if inFlag.Flag == nil {
				logger.Printf("flag %#v is unknown", arg)
			}
			continue

		// subcommand
		case subcommand(c, arg) != nil:
			logger.Printf("arg %#v is a subcommand\n", arg)

			switch {
			case c.DisableFlagParsing:
				logger.Printf("flag parsing disabled for %#v\n", c.Name())

			default:
				logger.Printf("parsing flags for %#v with args %#v\n", c.Name(), inArgs)
				if err := c.ParseFlags(inArgs); err != nil {
					return ActionMessage(err.Error()), context
				}
				context.Args = c.Flags().Args()
			}

			return traverse(subcommand(c, arg), args[i+1:])

		// positional
		default:
			logger.Printf("arg %#v is a positional\n", arg)
			inArgs = append(inArgs, arg)
		}
	}

	toParse := inArgs
	if inFlag != nil && len(inFlag.Args) == 0 && inFlag.Consumes("") {
		logger.Printf("removing arg %#v since it is a flag missing its argument\n", toParse[len(toParse)-1])
		toParse = toParse[:len(toParse)-1]
	} else if fs.IsShorthandSeries(context.CallbackValue) {
		logger.Printf("arg %#v is a shorthand flag series", context.CallbackValue)
		localInFlag := &InFlag{
			Flag: fs.LookupArg(context.CallbackValue),
			Args: []string{},
		}
		if localInFlag.Consumes("") && len(context.CallbackValue) > 2 {
			logger.Printf("removing shorthand %#v from flag series since it is missing its argument\n", localInFlag.Shorthand)
			toParse = append(toParse, strings.TrimSuffix(context.CallbackValue, localInFlag.Shorthand))
		} else {
			toParse = append(toParse, context.CallbackValue)
		}

	}

	// TODO duplicated code
	switch {
	case c.DisableFlagParsing:
		logger.Printf("flag parsing is disabled for %#v\n", c.Name())

	default:
		logger.Printf("parsing flags for %#v with args %#v\n", c.Name(), toParse)
		if err := c.ParseFlags(toParse); err != nil {
			return ActionMessage(err.Error()), context
		}
		context.Args = c.Flags().Args()
	}

	switch {
	// dash argument
	case common.IsDash(c):
		logger.Printf("completing dash for arg %#v\n", context.CallbackValue)
		context.Args = c.Flags().Args()[c.ArgsLenAtDash():]
		logger.Printf("context: %#v\n", context.Args)

		return storage.getPositional(c, len(context.Args)), context

	// flag argument
	case inFlag != nil && inFlag.Consumes(context.CallbackValue):
		logger.Printf("completing flag argument of %#v for arg %#v\n", inFlag.Name, context.CallbackValue)
		context.Parts = inFlag.Args
		return storage.getFlag(c, inFlag.Name), context

	// flag
	case !c.DisableFlagParsing && strings.HasPrefix(context.CallbackValue, "-"):
		if f := fs.LookupArg(context.CallbackValue); f != nil && f.IsOptarg() && strings.Contains(context.CallbackValue, string(f.OptargDelimiter())) {
			logger.Printf("completing optional flag argument for arg %#v\n", context.CallbackValue)
			prefix, optarg := f.Split(context.CallbackValue)
			context.CallbackValue = optarg
			return storage.getFlag(c, f.Name).Prefix(prefix), context
		}
		logger.Printf("completing flags for arg %#v\n", context.CallbackValue)
		return actionFlags(c), context

	// positional or subcommand
	default:
		logger.Printf("completing positionals and subcommands for arg %#v\n", context.CallbackValue)
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
