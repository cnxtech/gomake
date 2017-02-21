package cli

import "strings"

// Context is the context is which an Action is ran.
type Context struct {
	// Action is the context wrapped function to be evaluated.
	Action func() error

	flagSet map[string]struct{}
}

// NewContext initializes a new context for the Action to run in.
func NewContext(app *App, args []string) (*Context, error) {
	// Parse the flags first
	flagSet := ParseFlags(app.Flags, args)

	// Parse the commands
	action := ParseCommands(app.Action, app.Commands, args[len(flagSet):])

	// No appropriate action found, so we return ErrIncorrectUsage
	if action == nil {
		return nil, ErrIncorrectUsage
	}

	context := &Context{
		flagSet: flagSet,
	}

	context.Action = func() error {
		// Wrap the action call with context
		return action(context)
	}

	return context, nil
}

// IsSet returns whether flag with name is enabled.
func (c *Context) IsSet(name string) bool {
	_, ok := c.flagSet[name]
	return ok
}

// ParseFlags parses the args and returns a map of flags set.
func ParseFlags(flags Flags, args []string) map[string]struct{} {
	flagSet := make(map[string]struct{})

	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			break
		}

		alias := strings.TrimLeft(arg, "--")
		name := flags.NameForAlias(alias)
		if name != "" {
			flagSet[name] = struct{}{}
		}
	}

	return flagSet
}

// ParseCommands parses the args and returns the Action to invoke.
func ParseCommands(defaultAction Action, commands Commands, args []string) Action {
	var action func(ctx *Context) error

	if len(args) == 0 {
		// Default to application's default action if no commands are found
		action = defaultAction
	} else if len(args) == 1 {
		// Subcommands not supported, so only if there is one command, take that
		// command's Action
		action = commands.ActionForName(args[0])
	}

	return action
}
