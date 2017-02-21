package cli

import (
	"errors"
	"testing"
)

var (
	defaultErr = errors.New("default")
	gomakeErr  = errors.New("gomake")
)

func NewTestApp() *App {
	return &App{
		Action: func(ctx *Context) error {
			return defaultErr
		},
		Flags: []*Flag{
			{
				Name:    "help",
				Aliases: []string{"h"},
			},
		},
		Commands: Commands{
			{
				Name: "gomake",
				Action: func(ctx *Context) error {
					return gomakeErr
				},
			},
		},
	}
}

func TestNewContext(t *testing.T) {
	app := NewTestApp()

	// Test that empty args will set the default action
	context, err := NewContext(app, []string{})
	if err != nil {
		t.Errorf("Unexpected err: %s", err)
	}

	err = context.Action()
	if err != defaultErr {
		t.Errorf("Expected %s but got %s", defaultErr, err)
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
}

func TestParseFlags(t *testing.T) {
	app := NewTestApp()

	// Test that a unknown flag will not set flag
	flagSet := ParseFlags(app.Flags, []string{"--unknown"})
	_, ok := flagSet["unknown"]
	if ok {
		t.Errorf("Expected unknown to be not set")
	}

	// Test that a known flag will set flag
	flagSet = ParseFlags(app.Flags, []string{"--help"})
	_, ok = flagSet["help"]
	if !ok {
		t.Errorf("Expected help to be set")
	}

	// Test that a known flag's alias will set flag
	flagSet = ParseFlags(app.Flags, []string{"--h"})
	_, ok = flagSet["help"]
	if !ok {
		t.Errorf("Expected help to be set")
	}
}

func TestParseCommands(t *testing.T) {
	app := NewTestApp()
	ctx := &Context{}

	// Test that no args return the default Action
	action := ParseCommands(app.Action, app.Commands, []string{})
	err := action(ctx)
	if err != defaultErr {
		t.Errorf("Expected %s but got %s", defaultErr, err)
	}

	// Test that an unknown command returns nil
	action = ParseCommands(app.Action, app.Commands, []string{"unknown"})
	if action != nil {
		t.Errorf("Expected action to be nil")
	}

	// Test that an known flag returns nil
	action = ParseCommands(app.Action, app.Commands, []string{"--help"})
	if action != nil {
		t.Errorf("Expected action to be nil")
	}

	// Test that a known command return the command's Action
	action = ParseCommands(app.Action, app.Commands, []string{"gomake"})
	err = action(ctx)
	if err != gomakeErr {
		t.Errorf("Expected %s but got %s", gomakeErr, err)
	}
}
