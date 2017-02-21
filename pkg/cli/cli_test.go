package cli

import (
	"errors"
	"sort"
	"testing"
)

func TestRun(t *testing.T) {
	app := &App{}
	err := app.Run([]string{"gomake"})
	if err != nil {
		t.Errorf("Unexpected err")
	}
}

func TestShowHelp(t *testing.T) {
	app := &App{}
	err := app.ShowHelp()
	if err != nil {
		t.Errorf("Unexpected err")
	}
}

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

func TestNewContext(t *testing.T) {
	expected := errors.New("expected")
	gomakeErr := errors.New("gomake")
	app := &App{
		Action: func(ctx *Context) error {
			return expected
		},
		Flags: []*Flag{helpFlag},
		Commands: Commands{
			{
				Name: "gomake",
				Action: func(ctx *Context) error {
					return gomakeErr
				},
			},
		},
	}

	// Test that empty args will set the default action
	context, err := NewContext(app, []string{})
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
	}

	err = context.Action()
	if err != expected {
		t.Errorf("Expected %s but got %s", expected, err)
	}

	// Test that an unknown flag will return ErrIncorrectUsage
	context, err = NewContext(app, []string{"--unknown"})
	if err != ErrIncorrectUsage {
		t.Errorf("Expected %s but got %s", ErrIncorrectUsage, err)
	}

	// Test that an unknown command will return ErrIncorrectUsage
	context, err = NewContext(app, []string{"unknown"})
	if err != ErrIncorrectUsage {
		t.Errorf("Expected %s but got %s", ErrIncorrectUsage, err)
	}

	// Test that an unknown command  after a known command will return ErrIncorrectUsage
	context, err = NewContext(app, []string{"gomake", "unknown"})
	if err != ErrIncorrectUsage {
		t.Errorf("Expected %s but got %s", ErrIncorrectUsage, err)
	}

	// Test that a known flag after a known command will return ErrIncorrectUsage
	context, err = NewContext(app, []string{"gomake", "--help"})
	if err != ErrIncorrectUsage {
		t.Errorf("Expected %s but got %s", ErrIncorrectUsage, err)
	}

	// Test that a known flag will set flag on context
	context, err = NewContext(app, []string{"--help"})
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
	}

	if !context.IsSet("help") {
		t.Errorf("Expected help to be set")
	}

	// Test that a known command will set action on context
	context, err = NewContext(app, []string{"gomake"})
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
	}

	err = context.Action()
	if err != gomakeErr {
		t.Errorf("Expected %s but got %s", gomakeErr, err)
	}
}
