package cli

import (
	"sort"
	"testing"
)

func TestSortCommands(t *testing.T) {
	var commands Commands
	for _, name := range []string{"make", "test", "clean"} {
		command := &Command{
			Name: name,
		}

		commands = append(commands, command)
	}

	sort.Sort(commands)
	expected := []string{"clean", "make", "test"}

	for i, command := range commands {
		if command.Name != expected[i] {
			t.Errorf("Sorted order did not match expected")
		}
	}
}
