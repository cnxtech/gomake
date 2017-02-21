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
	gomake.Gomake(NewGomakefile()).Run(os.Args)
}

func NewGomakefile() *gomake.Gomakefile {
	gomakefile := gomake.NewGomakefile()

	goBuild := gomakefile.AddRule("gomake", "Builds gomake", nil, func() error {
		build := exec.Command("go", "build", "cmd/gomake/gomake.go")
		build.Stdout = os.Stdout
		build.Stderr = os.Stderr

		return build.Run()
	})

	gomakefile.AddRule("test", "Tests all the packages", []*gomake.Rule{goBuild}, func() error {
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

	gomakefile.Targets[""] = goBuild

	return gomakefile
}
