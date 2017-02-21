/*
Package gomake provides a makefile-like syntax to define rules and their
dependencies.

gomake is designed to be imported and used to create a main package with rules
to make your program. We can start off by writing a Gomakefile that will build
itself:

	package main

	import (
		"fmt"
		"os"
		"os/exec"

		"github.com/hinshun/gomake"
	)

	func main() {
		gomakefile := gomake.NewGomakefile()

		gomakeItself := gomakefile.AddRule("itself", "Builds gomake", nil, func() error {
			build := exec.Command("go", "build")
			build.Stdout = os.Stdout
			build.Stderr = os.Stderr

			return build.Run()
		})

		// Sets the default target
		gomakefile.Targets[""] = gomakeItself

		err := gomake.Gomake(gomakefile).Run(os.Args)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
*/
package gomake

import (
	"sort"

	"github.com/hinshun/gomake/pkg/cli"
)

const (
	// Version is the current version of gomake.
	Version = "0.1.0"
)

// Gomake creates a cli app for the given Gomakefile.
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
