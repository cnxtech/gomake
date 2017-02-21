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
	// ErrIncorrectUsage is returned when the App is ran with bad arguments.
	ErrIncorrectUsage = errors.New("incorrect usage")
)

// Action is a function to call with the context of arguments to the App.
type Action func(ctx *Context) error

// App is a simple cli application.
type App struct {
	// Name is the name of the program.
	Name string
	// Usage is a brief description of the program.
	Usage string
	// Version is the version of the program.
	Version string
	// Action is the default action to execute when no subcommands are specified.
	Action Action
	// Commands is the list of subcommands the program can run.
	Commands Commands
	// Flags is the list of boolean flags that can be enabled.
	Flags []*Flag
}

// Run runs the App with the given args and shows help on errors.
func (a *App) Run(args []string) error {
	// Initializes default variables for the App.
	a.initialize()

	context, err := NewContext(a, args[1:])
	if err != nil {
		a.ShowHelp()
		return err
	}

	if context.IsSet(HelpFlag.Name) {
		a.ShowHelp()
		return nil
	}

	if context.IsSet(VersionFlag.Name) {
		a.ShowVersion()
		return nil
	}

	return context.Action()
}

func (a *App) initialize() {
	a.Flags = append(a.Flags, HelpFlag, VersionFlag)

	if a.Version == "" {
		a.Version = "0.0.0"
	}

	if a.Action == nil {
		a.Action = func(ctx *Context) error {
			return nil
		}
	}
}

// ShowHelp displays the help text for the App.
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

// ShowHelp displays the version text for the App.
func (a *App) ShowVersion() {
	fmt.Printf("%v version %v\n", a.Name, a.Version)
}
