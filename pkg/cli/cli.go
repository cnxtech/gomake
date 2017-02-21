package cli

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"strings"
	"text/tabwriter"
)

var (
	ErrIncorrectUsage = errors.New("incorrect usage")
)

// App is a simple cli application
type App struct {
	Name     string
	Usage    string
	Version  string
	Action   func(ctx *Context) error
	Commands Commands
	Flags    []*Flag
}

// Run runs the App with the given args and shows help on errors
func (a *App) Run(args []string) error {
	a.setup()

	context, err := NewContext(a, args[1:])
	if err != nil {
		a.ShowHelp()
		return err
	}

	if context.IsSet(helpFlag.Name) {
		a.ShowHelp()
		return nil
	}

	if context.IsSet(versionFlag.Name) {
		a.ShowVersion()
		return nil
	}

	return context.Action()
}

// ShowHelp displays the help text for the App
func (a *App) ShowHelp() error {
	src := `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [options]{{if .Commands}} command{{end}}

VERSION:
   {{.Version}}

COMMANDS:{{range .Commands}}
   {{.Name}}{{if .Description}}{{"\t"}}{{.Description}}{{end}}{{end}}

OPTIONS:{{range .Flags}}
   --{{.Name}}{{if .Aliases}}, {{join .Aliases ", "}}{{end}}{{"\t"}}{{.Usage}}{{end}}
`
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	helpTemplate := template.Must(template.New("help").Funcs(funcMap).Parse(src))

	writer := tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', 0)
	err := helpTemplate.Execute(writer, a)
	if err != nil {
		return err
	}

	writer.Flush()
	return nil
}

// ShowHelp displays the version text for the App
func (a *App) ShowVersion() {
	fmt.Printf("%v version %v\n", a.Name, a.Version)
}

func (a *App) setup() {
	a.Flags = append(a.Flags, helpFlag, versionFlag)

	if a.Version == "" {
		a.Version = "0.0.0"
	}

	if a.Action == nil {
		a.Action = func(ctx *Context) error {
			return nil
		}
	}
}

// Command is a subcommand for a App
type Command struct {
	Name        string
	Description string
	Action      func(ctx *Context) error
}

type Commands []*Command

// Len returns the length of commands
func (c Commands) Len() int {
	return len(c)
}

// Len returns whether the command at index i is less than at index j
func (c Commands) Less(i, j int) bool {
	return c[i].Name < c[j].Name
}

// Swap swaps the commands at index i and j
func (c Commands) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Context is the context is which an Action is ran
type Context struct {
	Action  func() error
	flagSet map[string]struct{}
}

// NewContext parses arguments and returns a context for what action to run
// and what flags are enabled
func NewContext(app *App, args []string) (*Context, error) {
	// Parse the flags first
	flagSet := make(map[string]struct{})
	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			break
		}

		flagName := strings.TrimLeft(arg, "--")
		for _, flag := range app.Flags {
			aliases := append([]string{flag.Name}, flag.Aliases...)
			for _, alias := range aliases {
				if flagName == alias {
					flagSet[flag.Name] = struct{}{}
					continue
				}
			}
		}
	}

	// Parse the commands
	commands := args[len(flagSet):]
	var action func(ctx *Context) error
	if len(commands) == 0 {
		// Default to application's default action if no commands are found
		action = app.Action
	} else if len(commands) == 1 {
		// Subcommands not supported, so only if there is one command, take that
		// command's Action
		for _, command := range app.Commands {
			if commands[0] == command.Name {
				action = command.Action
			}
		}
	}

	// No appropriate action found, so we return ErrIncorrectUsage
	if action == nil {
		return nil, ErrIncorrectUsage
	}

	context := &Context{
		flagSet: flagSet,
	}

	// Wrap the action call with context
	context.Action = func() error {
		return action(context)
	}

	return context, nil
}

// IsSet returns whether flag with name is enabled
func (c *Context) IsSet(name string) bool {
	_, ok := c.flagSet[name]
	return ok
}

// Flag is a boolean flag that gets passed down to the action called
type Flag struct {
	Name    string
	Aliases []string
	Usage   string
}

var versionFlag = &Flag{
	Name:    "version",
	Aliases: []string{"v"},
	Usage:   "print the version",
}

var helpFlag = &Flag{
	Name:    "help",
	Aliases: []string{"h"},
	Usage:   "show help",
}
