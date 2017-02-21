package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hinshun/gomake"
)

const (
	gomakeName = "hinshun/gomake"
)

func main() {
	err := gomake.Gomake(NewGomakefile()).Run(os.Args)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(-1)
	}
}

func NewGomakefile() *gomake.Gomakefile {
	gomakefile := gomake.NewGomakefile()

	rebuild := gomakefile.AddRule("gomake", nil, func() error {
		cmd := exec.Command("go", "build", "cmd/gomake/gomake.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	})
	rebuild.Description = "Rebuilds gomake"

	test := gomakefile.AddRule("test", nil, func() error {
		cmd := exec.Command("go", "test", "./...")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	})
	test.Description = "Tests all the packages"

	clean := gomakefile.AddRule("clean", nil, func() error {
		err := os.Remove("gomake")
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		return nil
	})
	clean.Description = "Removes gomake"

	gomakefile.Targets[""] = rebuild

	return gomakefile
}
