package gomake

import (
	"sort"

	"github.com/hinshun/gomake/pkg/cli"
)

const (
	Version = "0.1.0"
)

// Gomake creates a cli app for the given Gomakefile
func Gomake(gomakefile *Gomakefile) *cli.App {
	app := &cli.App{
		Name:    "gomake",
		Usage:   "Makefile for gophers",
		Version: Version,
		Action: func(ctx *cli.Context) error {
			_, ok := gomakefile.Targets[""]
			if !ok {
				return nil
			}

			results := gomakefile.Make("")
			return HandleResults(results)
		},
	}

	for gomakeTarget, rule := range gomakefile.Targets {
		// Skip default
		if gomakeTarget == "" {
			continue
		}

		// Create closure around target for command
		target := gomakeTarget
		command := &cli.Command{
			Name:        target,
			Description: rule.Description,
			Action: func(ctx *cli.Context) error {
				results := gomakefile.Make(target)
				return HandleResults(results)
			},
		}

		app.Commands = append(app.Commands, command)
	}

	sort.Sort(app.Commands)
	return app
}
