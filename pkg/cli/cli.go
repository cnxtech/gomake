/*
Package cli is a very minimal framework for creating command line applications.

cli only supports boolean flags and top level commands, which is all that gomake
needs. We can write a simple greeter like so:

	package main

	import (
		"fmt"
		"os"

		"github.com/hinshun/gomake/pkg/cli"
	)

	func main() {
		app := cli.App{
			Name: "greeter",
			Action: func(ctx *cli.Context) error {
				fmt.Println("hello")
				return nil
			},
		}

		app.Run(os.Args)
	}
*/
package cli
