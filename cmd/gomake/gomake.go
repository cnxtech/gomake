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

	rebuild := gomakefile.AddRule("gomake", "Rebuilds gomake", nil, func() error {
		build := exec.Command("go", "build", "cmd/gomake/gomake.go")
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr

		return build.Run()
	})

	gomakefile.AddRule("test", "Tests all the packages", nil, func() error {
		test := exec.Command("go", "test", "./...")
		test.Stdout = os.Stdout
		test.Stderr = os.Stderr

		return test.Run()
	})

	gomakefile.AddRule("clean", "Removes gomake", nil, func() error {
		err := os.Remove("gomake")
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		return nil
	})

	gomakefile.Targets[""] = rebuild

	return gomakefile
}
